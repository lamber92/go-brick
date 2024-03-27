package producer

import (
	"context"
	"fmt"
	"go-brick/berror"
	"go-brick/berror/bcode"
	"go-brick/bmq/brabbitmq/config"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	url   = "amqp://root:123456@localhost:5672/"
	vhost = "test" // must be created manually by the administrator

	key1 = "test_rabbitmq_producer_1"
	key2 = "test_rabbitmq_producer_2"
)

func TestPublishWithConfirm(t *testing.T) {
	p1, err := _testNewProducer1()
	if err != nil {
		t.Fatal(err)
	}
	p1.SetTraceFunc(_testTraceFunc)

	p2, err := _testNewProducer2()
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
				if err := p1.Publish(context.Background(), msg); err != nil {
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
				if err := p2.Publish(context.Background(), msg); err != nil {
					t.Error(err)
					return
				}
				g.Done()
			}()
		}
	}()
	g.Wait()
}

func TestPublishWithoutConfirm(t *testing.T) {
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
			NoConfirm:    true, // should not wait confirmation
		},
		Key: key1,
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	p1.DisableTrace()

	messagesCount := 500

	g := sync.WaitGroup{}
	g.Add(messagesCount)
	go func() {
		for i := 0; i < messagesCount; i++ {
			idx := i
			go func() {
				msg := _testBuildMsg([]byte(fmt.Sprintf("{\"_id\":%d}", idx)))
				if err := p1.Publish(context.Background(), msg); err != nil {
					t.Error(err)
					return
				}
				g.Done()
			}()
		}
	}()
	g.Wait()
	_ = p1.Close()
}

func TestClose(t *testing.T) {
	p1, err := _testNewProducer1()
	if err != nil {
		t.Fatal(err)
	}
	p1.SetTraceFunc(_testTraceFunc)

	g := sync.WaitGroup{}

	testFunc := func(i int) error {
		msg := _testBuildMsg([]byte(fmt.Sprintf("{\"_id\":%d}", i)))
		return p1.Publish(context.Background(), msg)
	}

	// build 3 goroutines loop call Publish()
	for j := 0; j < 3; j++ {
		g.Add(1)
		go func() {
			defer func() {
				g.Done()
			}()
			for i := 0; i < 100; i++ {
				if err := testFunc(i); err != nil {
					// only not client-close error need to mark
					if !berror.IsCode(err, bcode.ClientClosed) {
						t.Error(err)
					}
					return
				}
				time.Sleep(time.Second)
			}
		}()
	}

	time.Sleep(time.Second * 5)
	if err = p1.Close(); err != nil {
		t.Fatal(err)
	}
	g.Wait()
}

func TestReconnect(t *testing.T) {
	p1, err := _testNewProducer1()
	if err != nil {
		t.Fatal(err)
	}
	p1.SetTraceFunc(_testTraceFunc)

	g := sync.WaitGroup{}
	goroutinesCount := 3
	exit := make(chan struct{}, 3)

	testFunc := func(i int) error {
		msg := _testBuildMsg([]byte(fmt.Sprintf("{\"_id\":%d}", i)))
		return p1.Publish(context.Background(), msg)
	}

	// build 3 goroutines loop call Publish()
	for j := 0; j < goroutinesCount; j++ {
		jj := j
		g.Add(1)
		go func() {
			defer func() {
				t.Logf("[%d] goroutine exit", jj)
				g.Done()
			}()
			for i := 0; i < 500; i++ {
				select {
				case <-exit:
					return
				default:
					t.Logf("[%d] push before", jj)
					if err := testFunc(i); err != nil {
						// only not client-close error need to mark
						if !berror.IsCode(err, bcode.ClientClosed) {
							t.Error(err)
						} else {
							t.Log(err)
						}
					} else {
						t.Logf("[%d] push done", jj)
					}
					time.Sleep(time.Second)
				}
			}
		}()
	}

	time.Sleep(time.Second * 5)

	// simulate network disconnection
	if err = p1.client.connection.Close(); err != nil {
		t.Fatal(err)
	}

	// disconnect the network by various means 20s

	time.Sleep(time.Minute * 2)

	// stop loops
	for i := 0; i < goroutinesCount; i++ {
		exit <- struct{}{}
	}
}

func _testNewProducer1() (*Producer, error) {
	return New(&config.Config{
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
	}, 1)
}

func _testNewProducer2() (*Producer, error) {
	return New(&config.Config{
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
		Key: key2,
	}, 2)
}

func _testTraceFunc(ctx context.Context, err error, data *amqp.Publishing, since time.Duration) {
	if err == nil {
		fmt.Printf("body: %s | cost: %d\n", string(data.Body), since.Milliseconds())
	} else {
		fmt.Printf("body: %s | cost: %d | err: %s\n", string(data.Body), since.Milliseconds(), err.Error())
	}
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
