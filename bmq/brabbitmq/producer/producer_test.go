package producer

import (
	"fmt"
	"sync"
	"testing"
)

func TestProducer(t *testing.T) {
	var (
		insertKey = "insert"
		updateKey = "update"
	)

	if err := updateHub(insertKey, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "modi",
		Type:  "producer",
		Extra: &Config{
			Exchange:     "e.direct.erp",
			ExchangeType: "direct",
			RoutingKey:   "insert_opensearch_goods_lib",
			Reliable:     true,
			Persistent:   true,
		},
		//
		Key: insertKey,
	}); err != nil {
		panic(err)
	}

	if err := updateHub(updateKey, &mqConfig{
		Url:   "amqp://guest:guest@192.168.1.245:5672/",
		VHost: "modi",
		Type:  "producer",
		Extra: &Config{
			Exchange:     "e.direct.erp",
			ExchangeType: "direct",
			RoutingKey:   "update_opensearch_goods_lib",
			Reliable:     true,
			Persistent:   true,
		},
		//
		Key: updateKey,
	}); err != nil {
		panic(err)
	}

	ps, err := GetProducerGroup([]string{insertKey, updateKey})
	if err != nil {
		panic(err)
	}

	g := sync.WaitGroup{}
	for _, p := range ps {
		p := p
		g.Add(1)
		go func() {
			for i := 0; i < 10000; i++ {
				msg := BuildSimpleTextMsg([]byte(fmt.Sprintf("{\"_id\":%d}", i)), 2, 1)
				if err := p.PushWithConfirm(msg); err != nil {
					panic(err)
				}
			}
			g.Done()
		}()
	}
	g.Wait()
}
