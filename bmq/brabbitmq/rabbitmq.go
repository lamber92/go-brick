package brabbitmq

import (
	"context"
	"fmt"
	"go-brick/bmq/brabbitmq/config"
	"go-brick/bmq/brabbitmq/consumer"
	"sync"
)

var (
	consumerHub sync.Map // map[string][]*Consumer
	producerHub sync.Map // map[string]*Producer
)

// Init init RabbitMQ clients
func Init(ctx context.Context, keys []string, namespace ...string) error {
	f := func(key string) error {
		conf, err := config.LoadConfig(ctx, key, namespace...)
		if err != nil {
			return err
		}
		if err = updateHub(key, conf); err != nil {
			return err
		}
		return nil
	}
	for _, key := range keys {
		if err := f(key); err != nil {
			return err
		}
	}
	return nil
}

func updateHub(key string, conf *config.Config) (err error) {
	switch conf.Type {
	case config.TypeConsumer:
		if _, ok := consumerHub.Load(key); ok {
			return nil
		}
		consumerConf := conf.Extra.(*config.ConsumerConfig)

		consumers := make([]*consumer.Consumer, 0, consumerConf.ConsumerCount)
		for idx := 0; idx < consumerConf.ConsumerCount; idx++ {
			newConsumer, err := consumer.New(conf, uint(idx))
			if err != nil {
				return err
			}
			consumers = append(consumers, newConsumer)
		}
		consumerHub.Store(key, consumers)

	case TypeProducer:
		if _, ok := producerHub[key]; ok {
			return nil
		}
		cli, err := newClient(conf, 0)
		if err != nil {
			return err
		}
		producer, err := cli.newSimpleProducer()
		if err != nil {
			return err
		}
		cli.producer = producer
		producerHub[key] = cli
		// 拉起通知协程
		go handlerNotifyClose(cli)

	default:
		return fmt.Errorf("unsupported config type: %s", conf.Type)
	}

	return nil
}

// GetConsumer 获取已初始化的消费者实例
func GetConsumer(key string) ([]*consumer.Consumer, error) {
	if c, ok := consumerHub[key]; !ok {
		return nil, fmt.Errorf("could not find consumer-Key: %s", key)
	} else {
		consumerSlice := make([]*Consumer, 0, len(c))
		for _, v := range c {
			consumerSlice = append(consumerSlice, v.consumer)
		}
		return consumerSlice, nil
	}
}

// GetProducer 获取已初始化的生产者实例
func GetProducer(key string) (*Producer, error) {
	if c, ok := producerHub[key]; !ok {
		return nil, fmt.Errorf("batchCount not find producer-key: %s", key)
	} else {
		return c.producer, nil
	}
}

// GetProducerGroup 获取已初始化的一组生产者实例
// 如果有一个不成功，返回nil。否则按keys中顺序压入切片中
func GetProducerGroup(keys []string) ([]*Producer, error) {
	r := make([]*Producer, 0, len(keys))
	for _, key := range keys {
		if c, ok := producerHub[key]; !ok {
			return nil, fmt.Errorf("batchCount not find producer-key: %s", key)
		} else {
			r = append(r, c.producer)
		}
	}
	return r, nil
}

// Close 主动关闭RabbitMQ连接
func Close() (err error) {
	for _, cList := range consumerHub {
		for _, v := range cList {
			if err = v.ShutdownConsumers(); err != nil {
				return err
			}
		}
	}
	for _, v := range producerHub {
		if err = v.ShutdownProducers(); err != nil {
			return err
		}
	}
	return
}
