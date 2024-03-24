package producer

import (
	"fmt"
	"go-brick/berror"
	"go-brick/bmq/brabbitmq/config"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conf    *config.Config
	subConf *config.ProducerConfig

	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	reader     <-chan amqp.Delivery

	idx uint
}

func newClient(conf *config.Config, id uint) (client *Client, err error) {
	client = &Client{
		conf:    conf,
		subConf: conf.Extra.(*config.ProducerConfig),
		idx:     id,
	}
	if err = client.initChannel(); err != nil {
		return
	}
	if err = client.bindExchange(); err != nil {
		return
	}
	if err = client.bindQueue(); err != nil {
		return
	}
	return
}

// initChannel
func (cli *Client) initChannel() (err error) {
	cli.connection, err = amqp.Dial(cli.conf.Url + cli.conf.VHost)
	if err != nil {
		return
	}
	if cli.channel, err = cli.connection.Channel(); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"fail to init channel")
	}
	return
}

// bindExchange
func (cli *Client) bindExchange() error {
	if err := cli.channel.ExchangeDeclare(
		cli.subConf.Exchange,     // name of the exchange
		cli.subConf.ExchangeType, // type
		true,                     // durable
		false,                    // delete when complete
		false,                    // internal
		false,                    // noWait
		nil,                      // arguments
	); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"fail to bind exchange")
	}
	return nil
}

// bindQueue
func (cli *Client) bindQueue() error {
	queue, err := cli.channel.QueueDeclare(
		cli.subConf.Queue,     // name of the queue
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // noWait
		cli.subConf.QueueArgs, // arguments
	)
	if err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"fail to bind queue")
	}
	cli.queue = &queue

	if err = cli.channel.QueueBind(
		queue.Name,             // name of the queue
		cli.subConf.RoutingKey, // routing key
		cli.subConf.Exchange,   // source exchange
		false,                  // noWait
		nil,                    // arguments
	); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"find to bind queue")
	}
	return nil
}

func (cli *Client) reconnectForConsumer() error {
	// 退出并关闭现有的所有goroutine
	if forceQuit := cli.consumer.notifyGoroutineToExit(); forceQuit {
		return nil
	}

	logger.Common.Infof("[%d] RabbitMQ-Client[%s] reconnect...", cli.ID, cli.Conf.Key)

	consumerConf := cli.Conf.Extra.(*ConsumerConfig)
	// 重连
	if err := cli.initConnect(); err != nil {
		return err
	}
	// 重新初始化channel
	if err := cli.initChannel(consumerConf); err != nil {
		return err
	}
	// 重新初始化消费者
	if err := cli.reInitConsumers(consumerConf); err != nil {
		return err
	}

	if _, ok := consumerHub[cli.Conf.Key]; ok {
		consumerHub[cli.Conf.Key][cli.ID] = cli
	}
	// 重新拉起检测协程
	go handlerNotifyClose(cli)
	// 重新拉起消费协程
	cli.consumer.RunOnce(nil)
	return nil
}

func (cli *Client) reconnectForProducer() error {
	logger.Common.Infof("[%d] RabbitMQ-Client[%s] reconnect...", cli.ID, cli.Conf.Key)
	if err := cli.initConnect(); err != nil {
		return err
	}
	producer, err := cli.newSimpleProducer()
	if err != nil {
		return err
	}
	cli.producer = producer
	producerHub[cli.Conf.Key] = cli

	go handlerNotifyClose(cli)
	return nil
}

/*
handlerNotifyClose
有两种情况会触发这里：
1. 显式调用ShutDown()：此时消费协业务逻辑程组必然已经完全退出.无影响
2. 断网：此时消费协业务逻辑程组未退出.需要通知
*/
func handlerNotifyClose(cli *Client) {
	var (
		times     uint64 = 1
		notifyErr *amqp.Error
	)

	select {
	case notifyErr = <-cli.connection.NotifyClose(make(chan *amqp.Error)):
		logger.Common.Infof("[%d] RabbitMQ-Client[%s] connection closed: %v", cli.ID, cli.Conf.Key, notifyErr)
	case notifyErr = <-cli.Channel.NotifyClose(make(chan *amqp.Error)):
		_ = cli.connection.Close()
		logger.Common.Infof("[%d] RabbitMQ-Client[%s] channel closed: %v", cli.ID, cli.Conf.Key, notifyErr)
	}

	reconnect := func() bool {
		cli.operationLock.Lock()
		defer cli.operationLock.Unlock()

		switch cli.Conf.Type {
		case Consumer:
			if err := cli.reconnectForConsumer(); err != nil {
				logger.Common.Errorf(err, "[%d] RabbitMQ-Client[%s] reconnect failed...[retry-times: %d]", cli.ID, cli.Conf.Key, times)
				return false
			}
			return true
		case TypeProducer:
			if err := cli.reconnectForProducer(); err != nil {
				logger.Common.Errorf(err, "[%d] RabbitMQ-Client[%s] reconnect failed...[retry-times: %d]", cli.ID, cli.Conf.Key, times)
				return false
			}
			return true
		default:
			logger.Common.Errorf(fmt.Errorf("invalid RabbitMQ-Clien[%s]...Type[%s]", cli.Conf.Key, cli.Conf.Type), "")
			return false
		}
	}

	for {
		if reconnect() {
			break
		}
		times++
		time.Sleep(5 * time.Second)
	}
}

// ShutdownConsumers 待当前所有正在执行的消费协程完成最后一次完整消费后退出
func (cli *Client) ShutdownConsumers() error {
	cli.operationLock.Lock()
	defer cli.operationLock.Unlock()

	_ = cli.consumer.notifyGoroutineToExit()

	conf := cli.Conf.Extra.(*ConsumerConfig)
	// will close() the deliver channel
	if err := cli.Channel.Cancel(conf.Consumer, true); err != nil {
		return fmt.Errorf("[%d] RabbitMQ-Consumer[%s] cancel failed: %v", cli.ID, cli.Conf.Key, err)
	}
	if err := cli.connection.Close(); err != nil {
		return fmt.Errorf("[%d] RabbitMQ-Consumer[%s] connection close. err: %v", cli.ID, cli.Conf.Key, err)
	}
	defer logger.Common.Infof("[%d] RabbitMQ-Consumer[%s] shutdown OK", cli.ID, cli.Conf.Key)
	// 标识主动退出
	cli.ForceExit = true
	return nil
}

// ShutdownProducers 退出所有生产者连接
func (cli *Client) ShutdownProducers() error {
	cli.operationLock.Lock()
	defer cli.operationLock.Unlock()

	if err := cli.connection.Close(); err != nil {
		return fmt.Errorf("[%d] RabbitMQ-Producer[%s] connection close. err: %v", cli.ID, cli.Conf.Key, err)
	}
	defer logger.Common.Infof("[%d] RabbitMQ-Producer[%s] shutdown OK", cli.ID, cli.Conf.Key)
	// 标识主动退出
	cli.ForceExit = true
	return nil
}

func (cli *Client) buildLogPrefix() string {
	return "[rabbitmq-producer][client][" + cli.conf.Key + "][" + strconv.FormatUint(uint64(cli.idx), 10) + "] "
}
