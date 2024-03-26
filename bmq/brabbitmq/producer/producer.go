package producer

import (
	"context"
	"go-brick/berror"
	"go-brick/berror/bcode"
	"go-brick/blog/logger"
	"go-brick/bmq/brabbitmq/config"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultMaxWaitConfirmTime = time.Second * 2 //
	defaultMaxWaitPublishTime = time.Second * 2
)

func init() {
	if err := berror.RegisterCustomizedMapping(
		amqp.ErrClosed,
		berror.NewClientClose(nil, amqp.ErrClosed.Error()),
	); err != nil {
		panic(err)
	}
}

type Producer struct {
	client *Client
	//
	existing    bool
	exitMonitor chan struct{}
	//
	maxRetryTimes    uint                                 // max retry times (default is 5, don't retry when set to 0)
	timeIntervalFunc func(times uint) (sec time.Duration) // retry interval strategy method
	//
	publishTimer *time.Timer
	exitPublish  chan struct{}
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
		//
		publishTimer: time.NewTimer(0),
		exitPublish:  make(chan struct{}, 1), // because the consumer of this pipe is not resident, a buffer is needed, otherwise the producer will block.
		//
		trace:     true,
		traceFunc: defaultTraceFunc,
	}
	return res, nil
}

// SetMaxRetryTimes 设置最大重试次数
func (p *Producer) SetMaxRetryTimes(times uint) *Producer {
	p.maxRetryTimes = times
	return p
}

// SetRetryTimesInterval 设置重试时间间隔策略
func (p *Producer) SetRetryTimesInterval(f func(times uint) (sec time.Duration)) *Producer {
	p.timeIntervalFunc = f
	return p
}

func (p *Producer) DisableTrace() *Producer {
	p.trace = false
	return p
}

func (p *Producer) SetTraceFunc(f TraceFunc) *Producer {
	p.traceFunc = f
	return p
}

func (p *Producer) GetKey() string {
	return p.client.conf.Key
}

// Publish push one message to rabbitmq server
func (p *Producer) Publish(ctx context.Context, data *amqp.Publishing) (err error) {
	var times uint = 0
	if p.trace {
		begin := time.Now()
		defer func() {
			p.traceFunc(ctx, err, data, time.Since(begin))
		}()
	}

	for ; times <= p.maxRetryTimes; times++ {
		if times > 0 {
			_ = <-time.After(p.timeIntervalFunc(times))
		}
		if p.existing {
			return berror.NewClientClose(nil, p.buildLogPrefix()+"client is going to exit")
		}
		if err = p.publish(ctx, data); err != nil {
			logger.Infra.WithError(err).
				With(logger.NewField().String("body", string(data.Body))).
				With(logger.NewField().String("message_id", data.MessageId)).
				Warn("push message fail")
			continue
		}
		if !p.client.subConf.NoConfirm {
			confirmation, err2 := p.getConfirmation()
			if err2 != nil {
				err = err2
				// at this point, if get an error because the client connection is interrupted,
				// there is no need to return an error because the message has most likely been delivered.
				//
				// if return err and push message again,
				// the message will be duplicated;
				// else, it will cause the message to be lost.
				// the latter is chosen here.
				if berror.IsCode(err, bcode.ClientClosed) {
					return nil
				}
				logger.Infra.WithError(err).
					With(logger.NewField().String("body", string(data.Body))).
					With(logger.NewField().String("message_id", data.MessageId)).
					Warn("get message confirmation fail")
				continue
			} else if !confirmation.Ack {
				logger.Infra.WithError(err).
					With(logger.NewField().String("body", string(data.Body))).
					With(logger.NewField().String("message_id", data.MessageId)).
					Warn("push message fail")
				continue
			}
		}
		return
	}
	return
}

func (p *Producer) publish(ctx context.Context, data *amqp.Publishing) (err error) {
	p.publishTimer.Reset(defaultMaxWaitPublishTime)
	defer p.publishTimer.Stop()
	select {
	case _, ok := <-p.exitPublish:
		if ok {
			close(p.exitPublish)
		}
		return berror.NewClientClose(nil, p.buildLogPrefix()+"client is going to exit")
	case <-p.publishTimer.C:
		return berror.NewGatewayTimeout(nil, p.buildLogPrefix()+"publish timeout")
	default:
		if err = p.client.channel.PublishWithContext(
			ctx,
			p.client.subConf.Exchange,   // publish to an exchange
			p.client.subConf.RoutingKey, // routing to 0 or more queues
			false,                       // mandatory
			false,                       // immediate
			*data,
		); err != nil {
			return berror.Convert(err, p.buildLogPrefix()+"publish message fail | body: "+string(data.Body))
		}
		return
	}
}

// getConfirmation receive the response from rabbitmq-server after receiving the push message
func (p *Producer) getConfirmation() (amqp.Confirmation, error) {
	p.publishTimer.Reset(defaultMaxWaitConfirmTime)
	defer p.publishTimer.Stop()
	select {
	case _, ok := <-p.exitPublish:
		if ok {
			close(p.exitPublish)
		}
		return amqp.Confirmation{}, berror.NewClientClose(nil, p.buildLogPrefix()+"client is going to exit")
	case <-p.publishTimer.C:
		return amqp.Confirmation{}, berror.NewGatewayTimeout(nil, p.buildLogPrefix()+"confirm timeout")
	case info, ok := <-p.client.confirms:
		if !ok {
			return amqp.Confirmation{}, berror.NewClientClose(nil, p.buildLogPrefix()+"confirm channel has been close")
		}
		return info, nil
	}
}

// monitor listen the underlying connection notify
func (p *Producer) monitor() {
	logger.Infra.Infof(p.buildLogPrefix() + "monitor start")

	for {
		var (
			times        uint64 = 1
			notification string
		)
		select {
		case notifyErr, ok := <-p.client.connection.NotifyClose(make(chan *amqp.Error)):
			if ok {
				notification = notifyErr.Error()
			}
			logger.Infra.Infof(p.buildLogPrefix()+"connection has been closed: %s", notification)
		case notifyErr, ok := <-p.client.channel.NotifyClose(make(chan *amqp.Error)):
			_ = p.client.connection.Close()
			if ok {
				notification = notifyErr.Error()
			}
			logger.Infra.Infof(p.buildLogPrefix()+"channel has been closed: %s", notification)
		case notification, _ = <-p.client.channel.NotifyCancel(make(chan string)):
			_ = p.client.connection.Close()
			logger.Infra.Infof(p.buildLogPrefix()+"queue has been deleted: %s", notification)
		case _, ok := <-p.exitMonitor:
			if ok {
				close(p.exitMonitor)
			}
			logger.Infra.Infof(p.buildLogPrefix() + "monitor exit")
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
	if p.existing {
		return berror.NewInternalError(nil, p.buildLogPrefix()+"has been closed")
	}
	p.existing = true

	if p.exitMonitor != nil {
		p.exitMonitor <- struct{}{}
	}
	if p.exitPublish != nil {
		p.exitPublish <- struct{}{}
	}
	// close client
	if err := p.client.close(); err != nil {
		return err
	}

	p.publishTimer.Stop()
	logger.Infra.Infof(p.buildLogPrefix() + "shutdown success")
	return nil
}

func (p *Producer) buildLogPrefix() string {
	return "[rabbitmq-producer][" + p.client.conf.Key + "][" + strconv.FormatUint(uint64(p.id), 10) + "] "
}
