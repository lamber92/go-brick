package config

import (
	"context"
	"fmt"
	"go-brick/bconfig"
	"go-brick/berror"
	"time"

	"github.com/spf13/cast"
)

type Type string

func (t Type) ToString() string {
	return string(t)
}

const (
	TypeConsumer Type = "consumer"
	TypeProducer Type = "producer"
)

type Config struct {
	Url   string
	VHost string
	Type  Type
	Extra any
	Key   string
}

type ProducerConfig struct {
	Queue        string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Persistent   bool // persistent messages
	QueueArgs    map[string]interface{}
}

type ConsumerConfig struct {
	Queue         string
	Consumer      string
	Exchange      string
	ExchangeType  string
	BindingKey    string
	PrefetchCount uint32
	ConsumerCount int
	QueueArgs     map[string]interface{}
}

func LoadConfig(ctx context.Context, key string, namespace ...string) (*Config, error) {
	v, err := bconfig.Static.Load(ctx, key, namespace...)
	if err != nil {
		return nil, berror.Convert(err, buildLogPrefix(key)+"failed to load rabbitmq config")
	}
	conf := &Config{}
	if err = v.Unmarshal(conf); err != nil {
		return nil, berror.Convert(err, buildLogPrefix(key)+"failed to unmarshal rabbitmq config: "+v.String())
	}

	switch conf.Type {
	case TypeConsumer:
		consumer := ConsumerConfig{}
		if err = v.Sub("Extra").Unmarshal(&consumer); err != nil {
			return nil, berror.Convert(err, buildLogPrefix(key)+"failed to unmarshal rabbitmq-consumer config: "+v.String())
		}
		if len(consumer.Queue) == 0 {
			return nil, berror.Convert(err, buildLogPrefix(key)+"'queue name' cannot be empty: "+v.String())
		}
		// at least 1 consumer
		if consumer.ConsumerCount == 0 {
			consumer.ConsumerCount = 1
		}
		// consumer_tag is required
		if len(consumer.Consumer) == 0 {
			consumer.Consumer = fmt.Sprintf("%s_%d", key, time.Now().UnixNano())
		}
		// handle special param
		if consumer.QueueArgs != nil {
			// must be of integer
			if v, ok := consumer.QueueArgs["x-message-ttl"]; ok {
				consumer.QueueArgs["x-message-ttl"] = cast.ToInt(v)
			}
		}
		conf.Extra = &consumer
	case TypeProducer:
		producer := ProducerConfig{}
		if err = v.Sub("Extra").Unmarshal(&producer); err != nil {
			return nil, berror.Convert(err, buildLogPrefix(key)+"failed to unmarshal rabbitmq-producer config: "+v.String())
		}
		if len(producer.Queue) == 0 {
			return nil, berror.Convert(err, buildLogPrefix(key)+"'queue name' cannot be empty: "+v.String())
		}
		// handle special param
		if producer.QueueArgs != nil {
			// must be of integer
			if v, ok := producer.QueueArgs["x-message-ttl"]; ok {
				producer.QueueArgs["x-message-ttl"] = cast.ToInt(v)
			}
		}
		conf.Extra = &producer
	default:
		return nil, berror.NewInternalError(nil, buildLogPrefix(key)+"unsupported config: "+conf.Type.ToString())
	}

	conf.Key = key
	return conf, nil
}

func buildLogPrefix(key string) string {
	return "[rabbitmq][" + key + "] "
}
