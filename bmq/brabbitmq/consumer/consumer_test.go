package consumer

import (
	"fmt"
	"go-brick/bmq/brabbitmq/config"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

func TestConsumerWork(t *testing.T) {
	if err := updateHub(Default, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "erp",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "q.durable.test.yishou.canal",
			Consumer:      "",
			Exchange:      "e.direct.erp",
			ExchangeType:  "direct",
			BindingKey:    "k.test",
			PrefetchCount: 1,
			ConsumerCount: 1,
		},
		//
		Key: Default,
	}); err != nil {
		panic(err)
	}
	consumers, err := GetConsumer(Default)
	if err != nil {
		panic(err)
	}

	//if err = database.InitMysql(nil, database.Default); err != nil {
	//	panic(err)
	//}

	for _, consumer := range consumers {
		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, ds []*amqp.Delivery, j int) error {
			for _, d := range ds {
				fmt.Println(string(d.Body))
				data := canal.CanalBinlogData{}
				_ = data.UnmarshalFromJson(d.Body)
				fmt.Println(data)

				// 下面注释的代码在测试mysql突然断连的时候使用
				//db, _ := database.GetDB(database.Default)
				//var a int
				//if err = db.Table("test").Select("a").Scan(&a).Error; err != nil {
				//	fmt.Println(err)
				//	return err
				//}
			}

			time.Sleep(time.Second * 2)
			if err = ds[len(ds)-1].Nack(true, true); err != nil {
				return err
			}
			return nil
		})

		time.Sleep(3 * time.Second)
		_ = consumer.testFakeCloseChannel() // 断线重连
	}

	time.Sleep(60 * time.Second)
	_ = Close() // 主动退出
}

func TestConsumerPanicRecovery(t *testing.T) {
	if err := updateHub(Default, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "erp",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "q.durable.erp.yishou.canal",
			Consumer:      "",
			Exchange:      "e.direct.erp",
			ExchangeType:  "direct",
			BindingKey:    "k.yishou.erp.canal",
			PrefetchCount: 1,
			ConsumerCount: 1,
		},
		//
		Key: Default,
	}); err != nil {
		panic(err)
	}
	consumers, err := GetConsumer(Default)
	if err != nil {
		panic(err)
	}

	testRecovery := func(ctx *rabbitmq.Context, delivery []*amqp.Delivery, j int) error {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r.(type) {
				case string:
					err = fmt.Errorf(r.(string))
				case error:
					err = r.(error)
				default:
					err = fmt.Errorf(fmt.Sprintf("%+v", r))
				}
				fmt.Printf("delivery hdrChain recover. err: %v", err)
				if err = delivery[len(delivery)-1].Nack(true, true); err != nil {
					fmt.Printf("delivery Nack failed. err: %v", err)
				}
			}
		}()
		return ctx.Next(delivery, j)
	}

	for _, c := range consumers {
		c.Use(testRecovery)

		b := false
		c.RunOnce(func(ctx *rabbitmq.Context, ds []*amqp.Delivery, j int) error {
			for _, d := range ds {
				if b {
					b = false
					fmt.Println(string(d.Body))
					panic("test panic")
				} else {
					b = true
				}
				fmt.Println(string(d.Body))
			}
			_ = ds[len(ds)-1].Nack(true, true)
			time.Sleep(time.Second)
			return nil
		})

		time.Sleep(3 * time.Second)
		_ = c.testFakeDisconnect() // 断线重连
	}

	time.Sleep(3 * time.Second)
	_ = Close() // 主动退出
}

func TestConsumerContextAndMiddleWare(t *testing.T) {
	if err := updateHub(Default, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "erp",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "q.durable.erp.yishou.canal",
			Consumer:      "",
			Exchange:      "e.direct.erp",
			ExchangeType:  "direct",
			BindingKey:    "k.yishou.erp.canal",
			PrefetchCount: 1,
			ConsumerCount: 1,
		},
		//
		Key: Default,
	}); err != nil {
		panic(err)
	}
	consumers, err := GetConsumer(Default)
	if err != nil {
		panic(err)
	}

	const Flag = "flag"

	for _, c := range consumers {
		c.Use(func(ctx *rabbitmq.Context, delivery []*amqp.Delivery, j int) error {
			flag := []string{"1"}
			ctx.Set(Flag, flag)

			err = ctx.Next(delivery, j)

			m, ok := ctx.Get(Flag)
			if !ok {
				fmt.Println("not ok")
				return fmt.Errorf("not ok")
			}
			fmt.Println(m)
			time.Sleep(time.Second)
			return err
		})
		c.Use(func(ctx *rabbitmq.Context, delivery []*amqp.Delivery, j int) error {
			m, ok := ctx.Get(Flag)
			if !ok {
				fmt.Println("not ok")
				return fmt.Errorf("not ok")
			}
			flag, ok := m.([]string)
			if !ok {
				fmt.Println("not ok")
				return fmt.Errorf("not ok")
			}
			flag = append(flag, "2")
			ctx.Set(Flag, flag)

			return ctx.Next(delivery, j)
		})

		c.RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
			m, ok := ctx.Get(Flag)
			if !ok {
				fmt.Println("not ok")
				return fmt.Errorf("not ok")
			}
			flag, ok := m.([]string)
			if !ok {
				fmt.Println("not ok")
				return fmt.Errorf("not ok")
			}
			flag = append(flag, "3")
			ctx.Set(Flag, flag)

			return d[len(d)-1].Nack(true, true)
		})
		time.Sleep(5 * time.Second)
		_ = c.testFakeDisconnect()
	}
}

func TestConsumerBizErrorRetry(t *testing.T) {
	if err := updateHub(Default, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "erp",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "q.durable.erp.yishou.canal",
			Consumer:      "",
			Exchange:      "e.direct.erp",
			ExchangeType:  "direct",
			BindingKey:    "k.yishou.erp.canal",
			PrefetchCount: 1,
			ConsumerCount: 1,
		},
		//
		Key: Default,
	}); err != nil {
		panic(err)
	}
	consumers, err := GetConsumer(Default)
	if err != nil {
		panic(err)
	}

	for _, consumer := range consumers {
		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
			_ = d[len(d)-1].Nack(true, true)
			return fmt.Errorf("test biz error retry")
		})

		time.Sleep(20 * time.Second)
		_ = consumer.testFakeDisconnect() // 断线重连
	}

	time.Sleep(20 * time.Second)
	_ = Close() // 主动退出
}

func TestConsumerAckFailedRetry(t *testing.T) {
	if err := updateHub(Default, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "erp",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "q.durable.erp.yishou.canal",
			Consumer:      "",
			Exchange:      "e.direct.erp",
			ExchangeType:  "direct",
			BindingKey:    "k.yishou.erp.canal",
			PrefetchCount: 1,
			ConsumerCount: 1,
		},
		//
		Key: Default,
	}); err != nil {
		panic(err)
	}
	consumers, err := GetConsumer(Default)
	if err != nil {
		panic(err)
	}

	for _, consumer := range consumers {
		// 注意！测试时需要手动改handleAckAndNack()里面的ack为失败返回err，否则观察不到整个流程
		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
			return &EventAckFail{}
		})

		time.Sleep(20 * time.Second)
		_ = consumer.testFakeDisconnect() // 断线重连
	}

	time.Sleep(20 * time.Second)
	_ = Close() // 主动退出
}

func TestConsumerPassByRetriesExceeded(t *testing.T) {
	if err := updateHub(Default, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "erp",
		Type:  "consumer",
		Extra: &ConsumerConfig{
			Queue:         "q.durable.erp.yishou.canal",
			Consumer:      "",
			Exchange:      "e.direct.erp",
			ExchangeType:  "direct",
			BindingKey:    "k.yishou.erp.canal",
			PrefetchCount: 1,
			ConsumerCount: 1,
		},
		//
		Key: Default,
	}); err != nil {
		panic(err)
	}
	consumers, err := GetConsumer(Default)
	if err != nil {
		panic(err)
	}

	for _, consumer := range consumers {
		consumer.EnableRetry().RunOnce(func(ctx *rabbitmq.Context, d []*amqp.Delivery, j int) error {
			_ = d[len(d)-1].Nack(true, true)
			//_ = d.ack(false)
			time.Sleep(time.Second)
			return fmt.Errorf("xxxx")
		})
	}

	time.Sleep(300 * time.Second)
	_ = Close() // 主动退出
}
