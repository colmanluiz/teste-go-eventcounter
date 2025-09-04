package eventcounter

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Consumer interface {
	Created(ctx context.Context, uid string) error
	Updated(ctx context.Context, uid string) error
	Deleted(ctx context.Context, uid string) error
}

type ConsumerWrapper struct {
	Consumer    Consumer
	EventCounts map[string]map[string]int
	mu          sync.Mutex
}

func NewConsumerWrapper() *ConsumerWrapper {
	return &ConsumerWrapper{
		EventCounts: make(map[string]map[string]int),
	}
}

func (c *ConsumerWrapper) randomSleep() {
	time.Sleep(time.Second * time.Duration(rand.Intn(30)))
}

func (c *ConsumerWrapper) Created(ctx context.Context, uid string) error {
	c.mu.Lock()
	if c.EventCounts[uid] == nil {
		c.EventCounts[uid] = make(map[string]int)
	}

	c.EventCounts[uid]["created"]++
	count := c.EventCounts[uid]["created"]
	c.mu.Unlock()

	log.Printf("User %s now has %d created events", uid, count)
	return nil
}

func (c *ConsumerWrapper) Updated(ctx context.Context, uid string) error {
	c.mu.Lock()
	if c.EventCounts[uid] == nil {
		c.EventCounts[uid] = make(map[string]int)
	}

	c.EventCounts[uid]["updated"]++
	count := c.EventCounts[uid]["updated"]
	c.mu.Unlock()

	log.Printf("User %s now has %d updated events", uid, count)
	return nil
}

func (c *ConsumerWrapper) Deleted(ctx context.Context, uid string) error {
	c.mu.Lock()
	if c.EventCounts[uid] == nil {
		c.EventCounts[uid] = make(map[string]int)
	}

	c.EventCounts[uid]["deleted"]++
	count := c.EventCounts[uid]["deleted"]
	c.mu.Unlock()

	log.Printf("User %s now has %d deleted events", uid, count)
	return nil
}
