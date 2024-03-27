package consumer

import (
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
	c1, err := _testNewConsumer1()
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
	_ = c1.Close()
}

func TestConsumerWorkRecover(t *testing.T) {
	c1, err := _testNewConsumer1()
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
	c1, err := _testNewConsumer1()
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
			panic("panic test")
		}
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
	_ = c1.Close()
}

//func TestConsumerBizErrorRetry(t *testing.T) {
//	if err := updateHub(Default, &mqConfig{
//		Url:   "amqp://guest:guest@192.168.1.245:5672/",
//		VHost: "erp",
//		Type:  "consumer",
//		Extra: &ConsumerConfig{
//			Queue:         "q.durable.erp.yishou.canal",
//			Consumer:      "",
//			Exchange:      "e.direct.erp",
//			ExchangeType:  "direct",
//			BindingKey:    "k.yishou.erp.canal",
//			PrefetchCount: 1,
//			ConsumerCount: 1,
//		},
//		//
//		Key: Default,
//	}); err != nil {
//		panic(err)
//	}
//	consumers, err := GetConsumer(Default)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, consumer := range consumers {
//		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
//			_ = d[len(d)-1].Nack(true, true)
//			return fmt.Errorf("test biz error retry")
//		})
//
//		time.Sleep(20 * time.Second)
//		_ = consumer.testFakeDisconnect() // 断线重连
//	}
//
//	time.Sleep(20 * time.Second)
//	_ = Close() // 主动退出
//}
//
//func TestConsumerAckFailedRetry(t *testing.T) {
//	if err := updateHub(Default, &mqConfig{
//		Url:   "amqp://guest:guest@192.168.1.245:5672/",
//		VHost: "erp",
//		Type:  "consumer",
//		Extra: &ConsumerConfig{
//			Queue:         "q.durable.erp.yishou.canal",
//			Consumer:      "",
//			Exchange:      "e.direct.erp",
//			ExchangeType:  "direct",
//			BindingKey:    "k.yishou.erp.canal",
//			PrefetchCount: 1,
//			ConsumerCount: 1,
//		},
//		//
//		Key: Default,
//	}); err != nil {
//		panic(err)
//	}
//	consumers, err := GetConsumer(Default)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, consumer := range consumers {
//		// 注意！测试时需要手动改handleAckAndNack()里面的ack为失败返回err，否则观察不到整个流程
//		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
//			return &EventAckFail{}
//		})
//
//		time.Sleep(20 * time.Second)
//		_ = consumer.testFakeDisconnect() // 断线重连
//	}
//
//	time.Sleep(20 * time.Second)
//	_ = Close() // 主动退出
//}
//
//func TestConsumerPassByRetriesExceeded(t *testing.T) {
//	if err := updateHub(Default, &mqConfig{
//		Url:   "amqp://guest:guest@192.168.1.245:5672/",
//		VHost: "erp",
//		Type:  "consumer",
//		Extra: &ConsumerConfig{
//			Queue:         "q.durable.erp.yishou.canal",
//			Consumer:      "",
//			Exchange:      "e.direct.erp",
//			ExchangeType:  "direct",
//			BindingKey:    "k.yishou.erp.canal",
//			PrefetchCount: 1,
//			ConsumerCount: 1,
//		},
//		//
//		Key: Default,
//	}); err != nil {
//		panic(err)
//	}
//	consumers, err := GetConsumer(Default)
//	if err != nil {
//		panic(err)
//	}
//
//	for _, consumer := range consumers {
//		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
//			_ = d[len(d)-1].Nack(true, true)
//			//_ = d.ack(false)
//			time.Sleep(time.Second)
//			return fmt.Errorf("xxxx")
//		})
//	}
//
//	time.Sleep(300 * time.Second)
//	_ = Close() // 主动退出
//}

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

func _testNewConsumer1() (*Consumer, error) {
	return New(&config.Config{
		Url:   url,
		VHost: vhost,
		Type:  config.TypeProducer,
		Extra: &config.ConsumerConfig{
			Queue:         "test_queue_1",
			Exchange:      "e.direct.test",
			ExchangeType:  config.ExchangeTypeDirect,
			BindingKey:    "test_rabbitmq_producer_key_1",
			PrefetchCount: 10,
		},
		Key: key1,
	}, 1)
}
