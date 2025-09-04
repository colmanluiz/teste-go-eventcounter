package eventcounter

import (
	"context"
	"log"
	"math/rand"
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
	if c.EventCounts[uid] == nil {
		c.EventCounts[uid] = make(map[string]int)
	}

	c.EventCounts[uid]["created"]++
	log.Printf("User %s now has %d created events", uid,
		c.EventCounts[uid]["created"])
	return nil
}

func (c *ConsumerWrapper) Updated(ctx context.Context, uid string) error {
	if c.EventCounts[uid] == nil {
		c.EventCounts[uid] = make(map[string]int)
	}

	c.EventCounts[uid]["updated"]++
	log.Printf("User %s now has %d updated events", uid,
		c.EventCounts[uid]["updated"])
	return nil
}

func (c *ConsumerWrapper) Deleted(ctx context.Context, uid string) error {
	if c.EventCounts[uid] == nil {
		c.EventCounts[uid] = make(map[string]int)
	}

	c.EventCounts[uid]["deleted"]++
	log.Printf("User %s now has %d deleted events", uid,
		c.EventCounts[uid]["deleted"])
	return nil
}
