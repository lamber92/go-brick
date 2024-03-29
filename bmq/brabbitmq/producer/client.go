package producer

import (
	"go-brick/berror"
	"go-brick/blog/logger"
	"go-brick/bmq/brabbitmq/config"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	defaultConfirmChanSize = 1000
)

type Client struct {
	conf    *config.Config
	subConf *config.ProducerConfig

	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	confirms   chan amqp.Confirmation

	idx uint
}

func newClient(conf *config.Config, id uint) (client *Client, err error) {
	client = &Client{
		conf:    conf,
		subConf: conf.Extra.(*config.ProducerConfig),
		idx:     id,
	}
	if len(client.subConf.Queue) == 0 {
		return nil, berror.NewInvalidArgument(nil, client.buildLogPrefix()+"queue name cannot be empty")
	}
	if err = client.initChannel(); err != nil {
		return
	}
	if err = client.initExchange(); err != nil {
		return
	}
	if err = client.initQueue(); err != nil {
		return
	}
	if err = client.initConfirms(); err != nil {
		return
	}
	return
}

// recover reestablish the underlying connection and reinitialize the queue
func (cli *Client) recover() (err error) {
	logger.Infra.Infof(cli.buildLogPrefix() + "try to recover now")
	if err = cli.initChannel(); err != nil {
		return
	}
	if err = cli.initExchange(); err != nil {
		return
	}
	if err = cli.initQueue(); err != nil {
		return
	}
	if err = cli.initConfirms(); err != nil {
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

// initExchange
func (cli *Client) initExchange() error {
	if err := cli.channel.ExchangeDeclare(
		cli.subConf.Exchange,                // name of the exchange
		cli.subConf.ExchangeType.ToString(), // type
		true,                                // durable
		false,                               // delete when complete
		false,                               // internal
		false,                               // noWait
		nil,                                 // arguments
	); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"fail to bind exchange")
	}
	return nil
}

// initQueue
func (cli *Client) initQueue() error {
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

func (cli *Client) initConfirms() error {
	if cli.subConf.NoConfirm {
		return nil
	}
	if err := cli.channel.Confirm(false); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"channel could not be put into confirm mode")
	}
	cli.confirms = cli.channel.NotifyPublish(make(chan amqp.Confirmation, defaultConfirmChanSize))
	return nil
}

func (cli *Client) close() error {
	if err := cli.connection.Close(); err != nil {
		return berror.Convert(err, cli.buildLogPrefix()+"close connection fail")
	}
	return nil
}

func (cli *Client) buildLogPrefix() string {
	return "[rabbitmq-producer][client][" + cli.conf.Key + "][" + strconv.FormatUint(uint64(cli.idx), 10) + "] "
}
