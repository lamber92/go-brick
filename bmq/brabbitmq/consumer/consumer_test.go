package consumer

import (
	"go-brick/berror"
	"go-brick/bmq/brabbitmq/config"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	url   = "amqp://root:123456@localhost:5672/"
	vhost = "test" // must be created manually by the administrator

	key1 = "test_rabbitmq_producer_1"
)

func TestConsumerWork(t *testing.T) {
	c1, err := _testNewConsumer1(10)
	if err != nil {
		t.Fatal(err)
	}
	if err = c1.Work(func(context *Context, deliveries []*amqp.Delivery, idx uint) error {
		for _, v := range deliveries {
			t.Logf("[%d] %s", idx, string(v.Body))
		}
		time.Sleep(time.Second)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Second)
	if err = c1.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestConsumerWorkRecover(t *testing.T) {
	c1, err := _testNewConsumer1(10)
	if err != nil {
		t.Fatal(err)
	}

	times := 0
	if err = c1.Work(func(context *Context, deliveries []*amqp.Delivery, idx uint) error {
		if times >= 5 {
			times = 0
			return nil
		} else {
			times++
			time.Sleep(time.Second * 5)
			panic("panic test")
		}
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
	_ = c1.Close()
}

func TestConsumerWorkPlugin(t *testing.T) {
	c1, err := _testNewConsumer1(10)
	if err != nil {
		t.Fatal(err)
	}
	if err = c1.
		Use(func(ctx *Context, deliveries []*amqp.Delivery, idx uint) error {
			for _, v := range deliveries {
				t.Logf("[%d] %s", idx, string(v.Body))
			}
			return ctx.Next(deliveries, idx)
		}).
		Work(func(ctx *Context, deliveries []*amqp.Delivery, idx uint) error {
			t.Log("-------------------")
			time.Sleep(time.Second * 2)
			return nil
		}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 10)
	_ = c1.Close()
}

func TestConsumerRetryInfinitely(t *testing.T) {
	c1, err := _testNewConsumer1(10)
	if err != nil {
		t.Fatal(err)
	}
	if err = c1.Work(func(context *Context, deliveries []*amqp.Delivery, idx uint) error {
		for _, v := range deliveries {
			t.Logf("[%d] %s", idx, string(v.Body))
		}
		t.Log("-------------------")
		time.Sleep(time.Second)
		return EventRetryInfinitely
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
	_ = c1.Close()
}

func TestConsumerRetryExceeded(t *testing.T) {
	c1, err := _testNewConsumer1(10)
	if err != nil {
		t.Fatal(err)
	}
	if err = c1.
		SetMaxRetryTimes(3).
		Work(func(context *Context, deliveries []*amqp.Delivery, idx uint) error {
			for _, v := range deliveries {
				t.Logf("[%d] %s", idx, string(v.Body))
			}
			t.Log("-------------------")
			time.Sleep(time.Second)
			return berror.NewInternalError(nil, "test retry exceeded")
		}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
	_ = c1.Close()
}

func TestConsumerReconnect(t *testing.T) {
	c1, err := _testNewConsumer1(3)
	if err != nil {
		t.Fatal(err)
	}
	if err = c1.
		SetMaxRetryTimes(3).
		Work(func(context *Context, deliveries []*amqp.Delivery, idx uint) error {
			for _, v := range deliveries {
				t.Logf("[%d] %s", idx, string(v.Body))
			}
			t.Log("-------------------")
			time.Sleep(time.Second)
			return nil
		}); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second * 10)
	_ = c1.client.connection.Close()
	time.Sleep(time.Second * 30)

	_ = c1.Close()
}

// _testSimulateDisconnect
// pure testing method, simulating network disconnection scenarios
func (c *Consumer) _testSimulateDisconnect() error {
	conf := c.client.conf.Extra.(*config.ConsumerConfig)
	// will close() the deliver channel
	if err := c.client.channel.Cancel(conf.Consumer, true); err != nil {
		return err
	}
	if err := c.client.connection.Close(); err != nil {
		return err
	}
	return nil
}

// _testSimulateCancelChannel
// pure test method, simulate the channel is deleted
func (c *Consumer) _testSimulateCancelChannel() error {
	// will close() the deliver channel
	if err := c.client.channel.Close(); err != nil {
		return err
	}
	return nil
}

func _testNewConsumer1(prefetchCount uint32) (*Consumer, error) {
	return New(&config.Config{
		Url:   url,
		VHost: vhost,
		Type:  config.TypeProducer,
		Extra: &config.ConsumerConfig{
			Queue:         "test_queue_1",
			Exchange:      "e.direct.test",
			ExchangeType:  config.ExchangeTypeDirect,
			BindingKey:    "test_rabbitmq_producer_key_1",
			PrefetchCount: prefetchCount,
		},
		Key: key1,
	}, 1)
}
