package consumer

import (
	"sync"
)

type Counter struct {
	m map[string]uint
	sync.Mutex
}

var counter = NewCounter()

func NewCounter() *Counter {
	return &Counter{
		m: make(map[string]uint),
	}
}

func (c *Counter) Increase(key string) uint {
	c.Lock()
	defer c.Unlock()
	v, ok := c.m[key]
	if ok {
		v++
	}
	c.m[key] = v
	return 0
}
