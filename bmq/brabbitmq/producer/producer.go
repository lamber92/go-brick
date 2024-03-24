package producer

import (
	"context"
	"go-brick/berror"
	"go-brick/blog/logger"
	"go-brick/bmq/brabbitmq/config"
	"go-brick/btrace"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	client   *Client
	confirms chan amqp.Confirmation

	maxRetryTimes    uint                                 // 最大重试次数（默认为5，设置为0时不重试）
	timeIntervalFunc func(times uint) (sec time.Duration) // 重试时间间隔策略方法

	trace bool
	id    uint
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
	return res, nil
}

func newDefaultProducer(cli *Client) (*Producer, error) {
	res := &Producer{
		client: cli,
		id:     cli.idx,
		//
		maxRetryTimes:    _defaultMaxRetryTimes,
		timeIntervalFunc: defaultRetryTimeInterval,
		// trace
		trace: true,
	}
	// TODO: 看看能不能放到client类里面
	if cli.subConf.Reliable {
		logger.Infra.Info(res.buildLogPrefix() + "enable requires confirmation")
		if err := cli.channel.Confirm(false); err != nil {
			return nil, berror.Convert(err, res.buildLogPrefix()+"channel could not be put into confirm mode")
		}
		res.confirms = cli.channel.NotifyPublish(make(chan amqp.Confirmation, 500))
	}
	return res, nil
}

// SetMaxRetryTimes 设置最大重试次数
func (p *Producer) SetMaxRetryTimes(times uint) {
	p.maxRetryTimes = times
}

// SetHandleRetryINRStrategy 设置重试时间间隔策略
func (p *Producer) SetHandleRetryINRStrategy(f func(times uint) (sec time.Duration)) {
	p.timeIntervalFunc = f
}

func (p *Producer) DisableTrace() {
	p.trace = false
}

// PushWithoutConfirm push message without confirm
func (p *Producer) PushWithoutConfirm(ctx context.Context, data *amqp.Publishing) (err error) {
	var times uint = 0
	if p.trace {
		begin := time.Now()
		defer func() {
			btrace.AppendMDIntoCtx(ctx, newTraceMD(data.Type, string(data.Body), times, time.Since(begin).Milliseconds()))
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
		}
	}
	return
}

// Push push message with confirm
func (p *Producer) Push(ctx context.Context, data *amqp.Publishing) (err error) {
	var times uint = 0
	if p.trace {
		begin := time.Now()
		defer func() {
			btrace.AppendMDIntoCtx(ctx, newTraceMD(data.Type, string(data.Body), times, time.Since(begin).Milliseconds()))
		}()
	}

	for ; times <= p.maxRetryTimes; times++ {
		if times > 0 {
			_ = <-time.After(p.timeIntervalFunc(times))
		}
		if err = p.Push(ctx, data); err != nil {
			err = berror.Convert(err, p.buildLogPrefix()+"push message fail. body: "+string(data.Body))
			continue
		}
		c := <-p.GetConfirm()
		if !c.Ack {
			err = berror.Convert(err, p.buildLogPrefix()+"get message confirm fail. body: "+string(data.Body))
			continue
		}
		return
	}
	return
}

// GetConfirm 获取MQ服务接收确认(Reliable==true时生效)
func (p *Producer) GetConfirm() <-chan amqp.Confirmation {
	return p.confirms
}

func (p *Producer) buildLogPrefix() string {
	return "[rabbitmq-producer][" + p.client.conf.Key + "][" + strconv.FormatUint(uint64(p.id), 10) + "] "
}
