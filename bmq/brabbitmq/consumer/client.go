package consumer

import (
	"go-brick/berror"
	"go-brick/blog/logger"
	"go-brick/bmq/brabbitmq/config"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conf    *config.Config
	subConf *config.ConsumerConfig

	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	reader     <-chan amqp.Delivery

	idx uint
}

func newClient(conf *config.Config, id uint) (client *Client, err error) {
	client = &Client{
		conf:    conf,
		subConf: conf.Extra.(*config.ConsumerConfig),
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
	if err = client.newReader(); err != nil {
		return
	}
	return
}

// recover reestablish the underlying connection and reinitialize the queue
func (cli *Client) recover() (err error) {
	logger.Infra.Infof(cli.buildLogPrefix() + "recover now")
	if err = cli.initChannel(); err != nil {
		return
	}
	if err = cli.bindExchange(); err != nil {
		return
	}
	if err = cli.bindQueue(); err != nil {
		return
	}
	if err = cli.newReader(); err != nil {
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
		cli.subConf.BindingKey, // binding key
		cli.subConf.Exchange,   // source exchange
		false,                  // noWait
		nil,                    // arguments
	); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"find to bind queue")
	}

	if err = cli.channel.Qos(int(cli.subConf.PrefetchCount), 0, true); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"fail to set Qos")
	}
	return nil
}

func (cli *Client) newReader() (err error) {
	cli.reader, err = cli.channel.Consume(
		cli.subConf.Queue,
		cli.subConf.Consumer, // must use different consumer_tag
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"fail to new consume")
	}
	return
}

func (cli *Client) close() error {
	if err := cli.channel.Cancel(cli.subConf.Consumer, true); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"close channel fail")
	}
	if err := cli.connection.Close(); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"close connection fail")
	}
	return nil
}

func (cli *Client) buildLogPrefix() string {
	return "[rabbitmq-consumer][client][" + cli.conf.Key + "][" + strconv.FormatUint(uint64(cli.idx), 10) + "] "
}
