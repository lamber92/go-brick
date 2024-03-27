package config

import (
	"context"
	"go-brick/bconfig"
	"go-brick/bconfig/bstorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYaml(t *testing.T) {
	bconfig.Init(bconfig.Option{
		Type:      bstorage.YAML,
		ConfigDir: "./test_config",
	})

	var (
		testProducerKey = "RabbitMQ.PublishSMS"
		testConsumerKey = "RabbitMQ.SubscribeSMS"
		testNamespace   = "test_yaml_config"
	)

	conf, err := Load(context.Background(), testProducerKey, testNamespace)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &Config{
		Url:   "amqp://root:123456@localhost:5672/",
		VHost: "test",
		Type:  "producer",
		Extra: &ProducerConfig{
			Queue:        "test_queue_1",
			Exchange:     "e.direct.test",
			ExchangeType: "direct",
			RoutingKey:   "test_rabbitmq_producer_key_1",
			Persistent:   true,
		},
		Key: testProducerKey,
	}, conf)

	conf, err = Load(context.Background(), testConsumerKey, testNamespace)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &Config{
		Url:   "amqp://root:123456@localhost:5672/",
		VHost: "test",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "test_queue_1",
			Consumer:      "Consumer.SubscribeSMS.1",
			Exchange:      "e.direct.test",
			ExchangeType:  "direct",
			BindingKey:    "test_rabbitmq_producer_key_1",
			PrefetchCount: 10,
			ConsumerCount: 1,
		},
		Key: testConsumerKey,
	}, conf)
}
