package producer

import (
	"context"
	"go-brick/berror"
	"go-brick/blog/logger"
	"go-brick/bmq/brabbitmq/config"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	client *Client
	//
	existing    bool
	exitMonitor chan struct{}
	//
	maxRetryTimes    uint                                 // max retry times (default is 5, don't retry when set to 0)
	timeIntervalFunc func(times uint) (sec time.Duration) // retry interval strategy method
	//
	confirm bool // identifies whether it is necessary to wait for rabbitmq to return message confirmation
	//
	trace     bool
	traceFunc TraceFunc
	id        uint
	//
	sync.Mutex
}

func New(basicConf *config.Config, idx uint) (*Producer, error) {
	cli, err := newClient(basicConf, idx)
	if err != nil {
		return nil, err
	}
	res, err := newDefaultProducer(cli)
	if err != nil {
		return nil, err
	}
	logger.Infra.Infof("[rabbitmq-consumer][%s][%d] init success", basicConf.Key, idx)
	// pull up monitor
	go res.monitor()
	return res, nil
}

func newDefaultProducer(cli *Client) (*Producer, error) {
	res := &Producer{
		client: cli,
		id:     cli.idx,
		//
		exitMonitor: make(chan struct{}),
		//
		maxRetryTimes:    _defaultMaxRetryTimes,
		timeIntervalFunc: defaultRetryTimeInterval,
		// trace
		trace:     true,
		traceFunc: defaultTraceFunc,
	}
	return res, nil
}

// SetMaxRetryTimes 设置最大重试次数
func (p *Producer) SetMaxRetryTimes(times uint) {
	p.maxRetryTimes = times
}

// SetRetryTimesInterval 设置重试时间间隔策略
func (p *Producer) SetRetryTimesInterval(f func(times uint) (sec time.Duration)) {
	p.timeIntervalFunc = f
}

func (p *Producer) DisableTrace() {
	p.trace = false
}

func (p *Producer) SetTraceFunc(f TraceFunc) {
	p.traceFunc = f
}

func (p *Producer) GetKey() string {
	return p.client.conf.Key
}

// PushWithoutConfirm push message without confirm
func (p *Producer) PushWithoutConfirm(ctx context.Context, data *amqp.Publishing) (err error) {
	var times uint = 0
	if p.trace {
		begin := time.Now()
		defer func() {
			p.traceFunc(ctx, data, time.Since(begin))
		}()
	}

	for ; times <= p.maxRetryTimes; times++ {
		if times > 0 {
			_ = <-time.After(p.timeIntervalFunc(times))
		}
		if err = p.client.channel.PublishWithContext(
			ctx,
			p.client.subConf.Exchange,   // publish to an exchange
			p.client.subConf.RoutingKey, // routing to 0 or more queues
			false,                       // mandatory
			false,                       // immediate
			*data,
		); err != nil {
			err = berror.Convert(err, p.buildLogPrefix()+"push message fail")
			logger.Infra.WithError(err).
				With(logger.NewField().String("body", string(data.Body))).
				With(logger.NewField().String("message_id", data.MessageId)).
				Warn("")
			continue
		}
		return
	}
	return
}

// Push
// push one message to rabbitmq server
func (p *Producer) Push(ctx context.Context, data *amqp.Publishing) (err error) {
	var times uint = 0
	if p.trace {
		begin := time.Now()
		defer func() {
			p.traceFunc(ctx, data, time.Since(begin))
		}()
	}

	for ; times <= p.maxRetryTimes; times++ {
		if times > 0 {
			_ = <-time.After(p.timeIntervalFunc(times))
		}
		if err = p.client.channel.PublishWithContext(
			ctx,
			p.client.subConf.Exchange,   // publish to an exchange
			p.client.subConf.RoutingKey, // routing to 0 or more queues
			false,                       // mandatory
			false,                       // immediate
			*data,
		); err != nil {
			err = berror.Convert(err, p.buildLogPrefix()+"push message fail. body: "+string(data.Body))
			logger.Infra.WithError(err).
				With(logger.NewField().String("body", string(data.Body))).
				With(logger.NewField().String("message_id", data.MessageId)).
				Warn("")
			continue
		}
		c := <-p.GetConfirm()
		if !c.Ack {
			err = berror.Convert(err, p.buildLogPrefix()+"get message confirm fail. body: "+string(data.Body))
			logger.Infra.WithError(err).
				With(logger.NewField().String("body", string(data.Body))).
				With(logger.NewField().String("message_id", data.MessageId)).
				Warn("")
			continue
		}
		return
	}
	return
}

// GetConfirm 获取MQ服务接收确认(Reliable==true时生效)
func (p *Producer) GetConfirm() <-chan amqp.Confirmation {
	return p.client.confirms
}

// monitor listen the underlying connection notify
func (p *Producer) monitor() {
	logger.Infra.Infof(p.buildLogPrefix() + "monitor start")

	for {
		var (
			times     uint64 = 1
			notifyErr *amqp.Error
		)
		select {
		case notifyErr = <-p.client.connection.NotifyClose(make(chan *amqp.Error)):
			logger.Infra.Infof(p.buildLogPrefix()+"connection closed: %v", notifyErr)
		case notifyErr = <-p.client.channel.NotifyClose(make(chan *amqp.Error)):
			_ = p.client.connection.Close()
			logger.Infra.Infof(p.buildLogPrefix()+"channel closed: %v", notifyErr)
		case _, ok := <-p.exitMonitor:
			if ok {
				close(p.exitMonitor)
			}
			logger.Infra.Infof(p.buildLogPrefix() + "monitor go to exit")
			return
		}

	LOOP:
		for {
			select {
			case _, ok := <-p.exitMonitor:
				if ok {
					close(p.exitMonitor)
				}
				logger.Infra.Infof(p.buildLogPrefix() + "monitor stop")
				return
			default:
				err := p.recover()
				if err == nil {
					break LOOP
				}
				logger.Infra.WithError(err).Errorf(p.buildLogPrefix()+"recover fail...[retry-times: %d]", times)
				times++
				time.Sleep(5 * time.Second)
			}
		}
	}
}

// recover automatic recovery
// automatically re-establish the underlying connection and retry consumption
func (p *Producer) recover() error {
	p.Lock()
	defer p.Unlock()
	// if the shutdown process is being executed, return directly
	if p.existing {
		return nil
	}
	// try to restore the connection
	if err := p.client.recover(); err != nil {
		return err
	}
	logger.Infra.Info(p.buildLogPrefix() + "recover success")
	return nil
}

// Close 待当前所有正在执行的消费协程完成最后一次完整消费后退出
func (p *Producer) Close() error {
	p.Lock()
	defer p.Unlock()
	// notify monitor go to exit
	p.existing = true
	p.exitMonitor <- struct{}{}
	// close client
	if err := p.client.close(); err != nil {
		return err
	}
	logger.Infra.Infof(p.buildLogPrefix() + "shutdown success")
	return nil
}

func (p *Producer) buildLogPrefix() string {
	return "[rabbitmq-producer][" + p.client.conf.Key + "][" + strconv.FormatUint(uint64(p.id), 10) + "] "
}
