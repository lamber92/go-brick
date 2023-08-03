package apollo_test

import (
	"context"
	"go-brick/bconfig/apollo"
	"go-brick/bconfig/yaml"
	"go-brick/berror"
	"go-brick/berror/bcode"
	"sync"
	"testing"
	"time"
)

func TestLoadDefaultNSConfig(t *testing.T) {
	// Pre-set configuration items: the value of Server.Access is `{"IpWhiteList":["127.0.0.1","192.168.1.1"]}`
	yaml.InitRootDir("./config_test")
	v, err := yaml.NewStatic().Load(context.Background(), "Apollo", "test")
	if err != nil {
		t.Fatal(err)
	}
	conf := &apollo.Config{}
	if err = v.Unmarshal(conf); err != nil {
		t.Fatal(err)
	}
	config, err := apollo.New(conf)
	if err != nil {
		t.Fatal(err)
	}
	value, err := config.Load(context.Background(), "Server.Access")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value.String())                      // {"IpWhiteList":["127.0.0.1","192.168.1.1"]}
	t.Log(value.GetStringSlice("IpWhiteList")) // [127.0.0.1 192.168.1.1]

	wg := sync.WaitGroup{}
	wg.Add(4)
	config.RegisterOnChange(func(event string) {
		t.Logf("event: %s", event) // event: {"Changes":{"Server.Access":{"OldValue":"{\"IpWhiteList\":[\"127.0.0.1\",\"192.168.1.1\"]}","NewValue":"{\"IpWhiteList\":[\"127.0.0.1\",\"192.168.1.1\",\"0.0.0.0\"]}","ChangeType":"MODIFY"}}}
		wg.Done()
	})
	wg.Wait()

	// 	modify the value of the configuration item `Server.Access` on Apollo background WEBUI to `{"IpWhiteList":["127.0.0.1","192.168.1.1","0.0.0.0"]}`
	for i := 0; i < 10; i++ {
		value, err = config.Load(context.Background(), "Server.Access")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(value.GetStringSlice("IpWhiteList")) // [127.0.0.1 192.168.1.1 0.0.0.0]
		time.Sleep(time.Millisecond * 200)
	}
}

func TestLoadOtherNSConfig(t *testing.T) {
	yaml.InitRootDir("./config_test")
	v, err := yaml.NewStatic().Load(context.Background(), "Apollo", "test")
	if err != nil {
		t.Fatal(err)
	}
	conf := &apollo.Config{}
	if err = v.Unmarshal(conf); err != nil {
		t.Fatal(err)
	}
	config, err := apollo.New(conf)
	if err != nil {
		t.Fatal(err)
	}

	// listen change first.
	// the "gateway" namespace has not been listened now, wait for the trigger below.
	wg := sync.WaitGroup{}
	// An `event` is triggered when the configuration is loaded for the first time.
	wg.Add(3)
	config.RegisterOnChange(func(event string) {
		t.Logf("event: %s", event)
		wg.Done()
	})

	time.Sleep(time.Second * 5)

	// get config first
	value, err := config.Load(context.Background(), "xxxx", "other_namespace")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("fisrt: %s", value.String()) // {"xxx":["127.0.0.1","0.0.0.0"]}

	// wait the config modified-notification
	// need to modify config on Apollo background WEBUI
	wg.Wait()

	// get config again
	value, err = config.Load(context.Background(), "xxxx", "other_namespace")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("second: %s", value.String()) // {"xxx":["127.0.0.1"]}

	value, err = config.Load(context.Background(), "Server.Access")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(value.GetStringSlice("IpWhiteList")) // [127.0.0.1 192.168.1.1 0.0.0.0]
}

func TestLoadInvalidNamespace(t *testing.T) {
	yaml.InitRootDir("./config_test")
	v, err := yaml.NewStatic().Load(context.Background(), "Apollo", "test")
	if err != nil {
		t.Fatal(err)
	}
	conf := &apollo.Config{}
	if err = v.Unmarshal(conf); err != nil {
		t.Fatal(err)
	}
	config, err := apollo.New(conf)
	if err != nil {
		t.Fatal(err)
	}

	value, err := config.Load(context.Background(), "Server.Access", "invalid_namespace")
	if err != nil {
		if err.(berror.Error).Status().Code().Is(bcode.NotFound) {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
		return
	}
	t.Log(value.GetStringSlice("IpWhiteList")) // [127.0.0.1 192.168.1.1 0.0.0.0]
}

func TestLoadInvalidKey(t *testing.T) {
	yaml.InitRootDir("./config_test")
	v, err := yaml.NewStatic().Load(context.Background(), "Apollo", "test")
	if err != nil {
		t.Fatal(err)
	}
	conf := &apollo.Config{}
	if err = v.Unmarshal(conf); err != nil {
		t.Fatal(err)
	}
	config, err := apollo.New(conf)
	if err != nil {
		t.Fatal(err)
	}

	value, err := config.Load(context.Background(), "invalid key")
	if err != nil {
		if err.(berror.Error).Status().Code().Is(bcode.NotFound) {
			t.Log(err)
		} else {
			t.Fatal(err)
		}
		return
	}
	t.Log(value.String()) // [127.0.0.1 192.168.1.1 0.0.0.0]
}
