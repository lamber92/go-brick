package brabbitmq_test

import (
	"context"
	"fmt"
	"go-brick/bconfig"
	"go-brick/bconfig/bstorage"
	"go-brick/bcontext"
	"go-brick/blog"
	"go-brick/bmq/brabbitmq"
	"os"
	"sync"
	"testing"
)

const (
	testProducerKey = "RabbitMQ.PublishSMS"
	testConsumerKey = "RabbitMQ.SubscribeSMS"
	testNamespace   = "dev_yaml_config"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("GO_ENV_NAME", "dev_apollo_config"); err != nil {
		panic(err)
	}
	bconfig.Init(bconfig.Option{
		Type:      bstorage.YAML,
		ConfigDir: "./config/test_config",
	})
	if err := brabbitmq.Init(
		context.Background(),
		[]string{
			testProducerKey,
			testConsumerKey,
		},
		testNamespace,
	); err != nil {
		panic(err)
	}
	m.Run()
}

func TestProducer(t *testing.T) {
	defer bconfig.Close()
	defer brabbitmq.Close()

	p, err := brabbitmq.GetProducer(testProducerKey)
	if err != nil {
		t.Fatal(err)
	}

	var (
		wg       = sync.WaitGroup{}
		msgCount = 500
	)
	wg.Add(500)
	for i := 0; i < msgCount; i++ {
		ii := i
		go func() {
			defer wg.Done()
			ctx := bcontext.New()
			if err2 := p.Publish(ctx,
				brabbitmq.BuildTextMsg4Publish(ctx, []byte(fmt.Sprintf("value: %d", ii)),
					true)); err2 != nil {
				blog.Error(ctx, err2, "publish message fail")
				return
			}
			blog.Info(ctx, "publish message success")
		}()
	}
	wg.Wait()
}

func TestConsumer(t *testing.T) {

}
