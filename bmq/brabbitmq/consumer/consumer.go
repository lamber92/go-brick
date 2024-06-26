package consumer

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lamber92/go-brick/berror"
	"github.com/lamber92/go-brick/berror/bstatus"
	"github.com/lamber92/go-brick/blog/logger"
	"github.com/lamber92/go-brick/bmq/brabbitmq/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultBatchFetchMessagePeriod = time.Second
)

type Consumer struct {
	client *Client
	id     uint

	existing   bool
	exitWorker chan struct{}
	// wait for all worker goroutines to exit
	exitWorkerDone sync.WaitGroup
	exitMonitor    chan struct{}

	// maximum messages count of a single read
	batchCount int
	// maximum duration of a single read
	batchPeriod time.Duration

	// ensure that the worker can only execute one at the same time
	working              *atomic.Bool
	handlerChain         []Handler
	handlerChainReadOnly []Handler
	//
	ackFunc  func([]*amqp.Delivery) error // 自动Re-Ack方法
	nackFunc func([]*amqp.Delivery) error // 自动Re-Nack方法
	//
	retryHdr *RetryHandler
	//
	trace     bool
	traceFunc TraceFunc
	//
	sync.Mutex
}

// New 创建消费者
func New(basicConf *config.Config, idx ...uint) (*Consumer, error) {
	var id uint = 0
	if len(idx) > 0 {
		id = idx[0]
	} else {
		id = counter.Increase(basicConf.Key)
	}

	cli, err := newClient(basicConf, id)
	if err != nil {
		return nil, err
	}
	res := newDefaultConsumer(cli)
	logger.Infra.Infof("[rabbitmq-consumer][%s][%d] init success", basicConf.Key, id)
	// pull up monitor
	go res.monitor()
	return res, nil
}

func newDefaultConsumer(cli *Client) *Consumer {
	res := &Consumer{
		client: cli,
		id:     cli.idx,
		//
		batchCount:  int(cli.subConf.PrefetchCount),
		batchPeriod: defaultBatchFetchMessagePeriod, // Default 2 seconds timeout
		//
		exitWorker:  make(chan struct{}),
		exitMonitor: make(chan struct{}),
		//
		working:              &atomic.Bool{},
		handlerChain:         make([]Handler, 0),
		handlerChainReadOnly: make([]Handler, 0),
		//
		ackFunc:  defaultAck,
		nackFunc: defaultNack,
		//
		trace:     true,
		traceFunc: defaultTraceFunc,
	}
	res.retryHdr = newRetryHandler(res)
	return res
}

// Use 增加中间件，仅执行消费前操作有效
func (c *Consumer) Use(fs ...Handler) *Consumer {
	if len(fs) > 0 {
		c.handlerChain = append(c.handlerChain, fs...)
	}
	return c
}

// DisableRetry 设置消费失败时重试，此操作仅适用于消费业务幂等的情况，否则请在业务逻辑中设计重试规则
func (c *Consumer) DisableRetry() *Consumer {
	c.retryHdr.enable = false
	return c
}

// SetHandleReAck 设置重Ack动作自定义func(在RunXXX前调用)
func (c *Consumer) SetHandleReAck(f func([]*amqp.Delivery) error) *Consumer {
	c.ackFunc = f
	return c
}

// SetHandleReNack 设置重Nack动作自定义func(在RunXXX前调用)
func (c *Consumer) SetHandleReNack(f func([]*amqp.Delivery) error) *Consumer {
	c.nackFunc = f
	return c
}

// SetMaxRetryTimes 设置最大重试次数
func (c *Consumer) SetMaxRetryTimes(times uint) *Consumer {
	c.retryHdr.maxRetryTimes = times
	return c
}

// SetHandleRetryTimeInterval 设置重试时间间隔策略
func (c *Consumer) SetHandleRetryTimeInterval(f func(times uint) time.Duration) *Consumer {
	c.retryHdr.timeIntervalFunc = f
	return c
}

// SetBatchFetchMessageCount set the maximum number of messages to be obtained in a single batch
func (c *Consumer) SetBatchFetchMessageCount(count int) *Consumer {
	c.batchCount = count
	return c
}

// SetBatchFetchMessagePeriod set the maximum waiting time for a single batch retrieval of messages
func (c *Consumer) SetBatchFetchMessagePeriod(duration time.Duration) *Consumer {
	c.batchPeriod = duration
	return c
}

func (c *Consumer) DisableTrace() *Consumer {
	c.trace = false
	return c
}

func (c *Consumer) GetID() uint {
	return c.id
}

func (c *Consumer) GetKey() string {
	return c.client.conf.Key
}

// Work only allowed to be consumed messages in one goroutine
func (c *Consumer) Work(f Handler) error {
	if f == nil {
		return berror.NewInternalError(nil, "work cannot be nil")
	}
	if !c.working.CompareAndSwap(false, true) {
		return berror.NewInternalError(nil, c.buildLogPrefix()+"the worker is working now")
	}
	if len(c.handlerChainReadOnly) == 0 {
		c.handlerChainReadOnly = append(c.handlerChainReadOnly, c.handlerChain...)
		c.handlerChainReadOnly = append(c.handlerChainReadOnly, c.handlerTrace, f)
	}
	go c.startWork()
	return nil
}

func (c *Consumer) handlerTrace(ctx *Context, ds []*amqp.Delivery, idx uint) error {
	if c.trace {
		now := time.Now()
		err := ctx.Next(ds, idx)
		c.traceFunc(ctx, err, ds, time.Since(now), idx)
		return err
	}
	return ctx.Next(ds, idx)
}

func (c *Consumer) startWork() {
	go func() {
		c.exitWorkerDone.Add(1)
		timer := time.NewTimer(0)
	LOOP:
		for {
			select {
			case _, ok := <-c.exitWorker:
				if ok {
					close(c.exitWorker)
				}
				break LOOP
			case tmp := <-c.client.reader:
				// refer to the k8s priority implementation plan:
				// https://github.com/kubernetes/kubernetes/blob/v1.25.4/pkg/controller/nodelifecycle/scheduler/taint_manager.go:274
				// make sure to check the exit-event,
				// and then process the messages in the queue.
			PRIORITY:
				for {
					select {
					case _, ok := <-c.exitWorker:
						if ok {
							close(c.exitWorker)
						}
						break LOOP
					default:
						break PRIORITY
					}
				}

				ds := make([]*amqp.Delivery, 0, c.batchCount)
				ds = append(ds, &tmp)

				// set a timer to obtain the maximum number of messages within the maximum waiting period.
				timer.Reset(c.batchPeriod)
			BATCH:
				for i := 0; i < c.batchCount-1; i++ {
					select {
					case sub := <-c.client.reader:
						ds = append(ds, &sub)
					case <-timer.C:
						break BATCH
					}
				}
				timer.Stop()

				// handle new messages
				if len(ds) > 0 {
					if err := c.handleMessage(ds); err != nil {
						break LOOP
					}
				}
			}
		}

		c.exitWorkerDone.Done()
		c.working.CompareAndSwap(true, false)
		logger.Infra.Info(c.buildLogPrefix() + "the consumer's goroutine is about to be closed")
		return
	}()
}

// handleMessage
// strategies for handling messages after they are successfully consumed or failed to be consumed
// returning non-nil indicates that the outer loop needs to be interrupted
func (c *Consumer) handleMessage(ds []*amqp.Delivery) error {
	err := newContext(c.handlerChainReadOnly).Handle(ds, c.id)
	if err == nil {
		// running to this point indicates that
		// the business side consumes successfully
		// and calls the Ack() method to set the message offset.
		return c.Ack(ds)
	}

	if !c.retryHdr.Enable() {
		return err
	}
	// running to this point indicates
	// that an error was encountered during processing
	// on the business side and needs to be retried.
	if c.retryHdr.InfiniteRetry(err) {
		if err = c.Nack(ds); err != nil {
			return err
		}
		// keep going and wait for the retry cycle
	} else {
		if c.retryHdr.ExceededLimit() {
			c.retryHdr.ClearRetriedTimes()
			logger.Infra.WithError(err).
				Errorf(c.buildLogPrefix()+"[Logic] retries retryTimes have exceeded limit: [%d], discard messages...", c.retryHdr.maxRetryTimes)
			return c.Ack(ds)
		} else {
			if err = c.Nack(ds); err != nil {
				return err
			}
			// keep going and wait for the retry cycle
		}
	}
	return c.retryHdr.waitForNextRetry(err, "Logic")
}

// Ack call rabbitmq.client Ack()
func (c *Consumer) Ack(ds []*amqp.Delivery) error {
	err := c.ackFunc(ds)
	if err == nil {
		c.retryHdr.ClearRetriedTimes()
		return nil
	}
	for {
		if err = c.retryHdr.waitForNextRetry(err, "Ack"); err != nil {
			return err
		}
		if err = c.ackFunc(ds); err == nil {
			c.retryHdr.ClearRetriedTimes()
			return nil
		}
	}
}

// Nack call rabbitmq.client Nack()
func (c *Consumer) Nack(ds []*amqp.Delivery) error {
	err := c.nackFunc(ds)
	if err == nil {
		return nil
	}
	for {
		if err = c.retryHdr.waitForNextRetry(err, "Nack"); err != nil {
			return err
		}
		if err = c.nackFunc(ds); err == nil {
			return nil
		}
	}
}

// monitor listen the underlying connection notify
func (c *Consumer) monitor() {
	logger.Infra.Infof(c.buildLogPrefix() + "monitor start")

	for {
		var (
			times        uint64 = 1
			notification string
		)
		select {
		case notifyErr, ok := <-c.client.connection.NotifyClose(make(chan *amqp.Error)):
			if ok {
				notification = notifyErr.Error()
			}
			logger.Infra.Infof(c.buildLogPrefix()+"connection has been closed: %s", notification)
		case notifyErr, ok := <-c.client.channel.NotifyClose(make(chan *amqp.Error)):
			if ok {
				notification = notifyErr.Error()
			}
			_ = c.client.connection.Close()
			logger.Infra.Infof(c.buildLogPrefix()+"channel has been closed: %s", notification)
		case notification, _ = <-c.client.channel.NotifyCancel(make(chan string)):
			_ = c.client.connection.Close()
			logger.Infra.Infof(c.buildLogPrefix()+"queue has been deleted: %s", notification)
		case _, ok := <-c.exitMonitor:
			if ok {
				close(c.exitMonitor)
			}
			logger.Infra.Infof(c.buildLogPrefix() + "monitor exit")
			return
		}

	LOOP:
		for {
			select {
			case _, ok := <-c.exitMonitor:
				if ok {
					close(c.exitMonitor)
				}
				logger.Infra.Infof(c.buildLogPrefix() + "monitor stop")
				return
			default:
				err := c.recover()
				if err == nil {
					break LOOP
				}
				logger.Infra.WithError(err).Errorf(c.buildLogPrefix()+"recover fail...[retry-times: %d]", times)
				times++
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// recover automatic recovery
// automatically re-establish the underlying connection and retry consumption
func (c *Consumer) recover() error {
	c.Lock()
	defer c.Unlock()
	// if the shutdown process is being executed, return directly
	if c.existing {
		return nil
	}
	// stop the running goroutine
	c.stopWork()
	// try to restore the connection
	if err := c.client.recover(); err != nil {
		return err
	}
	// ensure that there is only one worker currently and it has been executed
	if c.working.CompareAndSwap(false, true) && len(c.handlerChainReadOnly) > 0 {
		c.exitWorker = make(chan struct{})
		go c.startWork()
		logger.Infra.Info(c.buildLogPrefix() + "worker recover to run again")
	}
	logger.Infra.Info(c.buildLogPrefix() + "recover success")
	return nil
}

// stopWork exit all work-goroutine
func (c *Consumer) stopWork() {
	// notify current consumer to close all working goroutines
	select {
	case _, ok := <-c.exitWorker:
		if !ok {
			// chan has been closed
			break
		}
	default:
		c.exitWorker <- struct{}{}
	}
	// wait for all goroutine to operate
	c.exitWorkerDone.Wait()
	return
}

// Close 待当前所有正在执行的消费协程完成最后一次完整消费后退出
func (c *Consumer) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.existing {
		return berror.NewInternalError(nil, c.buildLogPrefix()+"has been closed")
	}
	c.existing = true

	if c.exitMonitor != nil {
		c.exitMonitor <- struct{}{}
	}
	// stop the running goroutine
	c.stopWork()
	// close client
	if err := c.client.close(); err != nil {
		return err
	}
	c.retryHdr.Close()

	logger.Infra.Infof(c.buildLogPrefix() + "shutdown success")
	return nil
}

func (c *Consumer) buildLogPrefix() string {
	return "[rabbitmq-consumer][" + c.client.conf.Key + "][" + strconv.FormatUint(uint64(c.id), 10) + "] "
}

// defaultAck default recall ack()
func defaultAck(ds []*amqp.Delivery) error {
	if len(ds) == 0 {
		return nil
	}
	if err := ds[len(ds)-1].Ack(true); err != nil {
		return berror.New(bstatus.New(EventCodeAckFail, err.Error(), nil))
	}
	return nil
}

// defaultNack default recall Nack()
func defaultNack(ds []*amqp.Delivery) error {
	if len(ds) == 0 {
		return nil
	}
	if err := ds[len(ds)-1].Nack(true, true); err != nil {
		return berror.New(bstatus.New(EventCodeNackFail, err.Error(), nil))
	}
	return nil
}
