package producer

import (
	"context"
	"fmt"
	"go-brick/bmq/brabbitmq/config"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestProducer_PushWithConfirm(t *testing.T) {
	const (
		key1  = "test_rabbitmq_producer_1"
		key2  = "test_rabbitmq_producer_2"
		url   = "amqp://root:123456@localhost:5672/"
		vhost = "test"
	)
	p1, err := New(&config.Config{
		Url:   url,
		VHost: vhost,
		Type:  config.TypeProducer,
		Extra: &config.ProducerConfig{
			Queue:        "test_queue_1",
			Exchange:     "e.direct.test",
			ExchangeType: config.ExchangeTypeDirect,
			RoutingKey:   "test_rabbitmq_producer_key_1",
			Persistent:   true,
		},
		Key: key1,
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	p1.SetTraceFunc(_testTraceFunc)

	p2, err := New(&config.Config{
		Url:   url,
		VHost: vhost,
		Type:  config.TypeProducer,
		Extra: &config.ProducerConfig{
			Queue:        "test_queue_2",
			Exchange:     "e.direct.test",
			ExchangeType: config.ExchangeTypeDirect,
			RoutingKey:   "test_rabbitmq_producer_key_2",
			Persistent:   true,
		},
		//
		Key: key2,
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	p2.SetTraceFunc(_testTraceFunc)

	g := sync.WaitGroup{}
	g.Add(200)
	go func() {
		for i := 0; i < 100; i++ {
			idx := i
			go func() {
				msg := _testBuildMsg([]byte(fmt.Sprintf("{\"_id\":%d}", idx)))
				if err := p1.Push(context.Background(), msg); err != nil {
					t.Error(err)
					return
				}
				g.Done()
			}()
		}
	}()
	go func() {
		for i := 100; i < 200; i++ {
			idx := i
			go func() {
				msg := _testBuildMsg([]byte(fmt.Sprintf("{\"_id\":%d}", idx)))
				if err := p2.Push(context.Background(), msg); err != nil {
					t.Error(err)
					return
				}
				g.Done()
			}()
		}
	}()
	g.Wait()
}

func _testTraceFunc(ctx context.Context, data *amqp.Publishing, since time.Duration) {
	fmt.Printf("body: %s | cost: %d\n", string(data.Body), since.Milliseconds())
}

func _testBuildMsg(b []byte) *amqp.Publishing {
	return &amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     "text/plain",
		ContentEncoding: "",
		Body:            b,
		DeliveryMode:    2,
		Priority:        0,
	}
}
