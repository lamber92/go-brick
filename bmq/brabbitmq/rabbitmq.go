package brabbitmq

import (
	"context"
	"fmt"
	"go-brick/berrgroup"
	"go-brick/berror"
	"go-brick/blog/logger"
	"go-brick/bmq/brabbitmq/config"
	"go-brick/bmq/brabbitmq/consumer"
	"go-brick/bmq/brabbitmq/producer"
	"go-brick/bstructure/bset"
	"sync"
)

var (
	consumerHub sync.Map // map[string][]*Consumer
	producerHub sync.Map // map[string]*Producer
)

// Init connect to RabbitMQ and init all producers and consumers
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

func updateHub(key string, conf *config.Config) error {
	switch conf.Type {
	case config.TypeConsumer:
		if _, ok := consumerHub.Load(key); ok {
			return nil
		}
		subConf := conf.Extra.(*config.ConsumerConfig)
		consumers := make([]*consumer.Consumer, 0, subConf.ConsumerCount)
		for idx := 0; idx < subConf.ConsumerCount; idx++ {
			newConsumer, err := consumer.New(conf, uint(idx))
			if err != nil {
				return err
			}
			consumers = append(consumers, newConsumer)
		}
		consumerHub.Store(key, consumers)
	case config.TypeProducer:
		if _, ok := producerHub.Load(key); ok {
			return nil
		}
		newProducer, err := producer.New(conf, 0)
		if err != nil {
			return err
		}
		producerHub.Store(key, newProducer)
	default:
		return fmt.Errorf("unsupported config type: %s", conf.Type)
	}
	return nil
}

// GetConsumer get the initialized consumer instance
func GetConsumer(key string) ([]*consumer.Consumer, error) {
	v, ok := consumerHub.Load(key)
	if !ok {
		return nil, berror.NewNotFound(nil, fmt.Sprintf("[rabbitmq-consumer][%s] has not been initialized yet!", key))
	}
	res, ok := v.([]*consumer.Consumer)
	if !ok {
		return nil, berror.NewInternalError(nil, fmt.Sprintf("[rabbitmq-consumer][%s] has invalid instance type", key))
	}
	return res, nil
}

// GetProducer get the initialized producer instance
func GetProducer(key string) (*producer.Producer, error) {
	v, ok := producerHub.Load(key)
	if !ok {
		return nil, berror.NewNotFound(nil, fmt.Sprintf("[rabbitmq-producer][%s] has not been initialized yet!", key))
	}
	res, ok := v.(*producer.Producer)
	if !ok {
		return nil, berror.NewInternalError(nil, fmt.Sprintf("[rabbitmq-producer][%s] has invalid instance type", key))
	}
	return res, nil
}

func CloseConsumer(key string) error {
	v, ok := consumerHub.Load(key)
	if !ok {
		return berror.NewNotFound(nil, fmt.Sprintf("[rabbitmq-consumer][%s] has not been initialized yet!", key))
	}
	list, ok := v.([]*consumer.Consumer)
	if !ok {
		return berror.NewInternalError(nil, fmt.Sprintf("[rabbitmq-consumer][%s] has invalid instance type", key))
	}
	eg, ctx := berrgroup.WithContext(context.Background())
	defer ctx.Cancel()
	for _, tmp := range list {
		cons := tmp
		eg.Go(func() error {
			if err := cons.Close(); err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	consumerHub.Delete(key)
	return nil
}

func CloseProducer(key string) error {
	v, ok := producerHub.Load(key)
	if !ok {
		return berror.NewNotFound(nil, fmt.Sprintf("[rabbitmq-producer][%s] has not been initialized yet!", key))
	}
	prod, ok := v.(*producer.Producer)
	if !ok {
		return berror.NewInternalError(nil, fmt.Sprintf("[rabbitmq-producer][%s] has invalid instance type", key))
	}
	if err := prod.Close(); err != nil {
		return err
	}
	producerHub.Delete(key)
	return nil
}

// Close exit all rabbitmq producers and consumers
func Close() {
	eg, ctx := berrgroup.WithContext(context.Background())
	defer ctx.Cancel()

	closedKeys := bset.NewSafeSet[string]()
	consumerHub.Range(func(key, value any) bool {
		closedKeys.Add(key.(string))
		tmp, ok := value.([]*consumer.Consumer)
		if !ok {
			logger.Infra.Errorf("[rabbitmq-consumer][%s] has invalid instance type", key)
			return true
		}
		for _, v := range tmp {
			cons := v
			eg.Go(func() error {
				if err := cons.Close(); err != nil {
					closedKeys.Delete(key.(string))
					return err
				}
				return nil
			})
		}
		return true
	})
	for _, key := range closedKeys.ToSlice() {
		consumerHub.Delete(key)
	}

	closedKeys.Clear()
	producerHub.Range(func(key, value any) bool {
		closedKeys.Add(key.(string))
		tmp, ok := value.(*producer.Producer)
		if !ok {
			logger.Infra.Errorf("[rabbitmq-producer][%s] has invalid instance type", key)
			return true
		}
		eg.Go(func() error {
			if err := tmp.Close(); err != nil {
				closedKeys.Delete(key.(string))
				return err
			}
			return nil
		})
		return true
	})
	for _, key := range closedKeys.ToSlice() {
		producerHub.Delete(key)
	}

	return
}
