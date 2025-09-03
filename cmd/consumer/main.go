package main

import (
	"context"
	"flag"
)

var (
	amqpUrl      string
	amqpExchange string
)

func init() {
	flag.StringVar(&amqpUrl, "amqp-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	flag.StringVar(&amqpExchange, "amqp-exchange", "eventcountertest", "RabbitMQ Exchange")
}

type ConsumerStr struct {
	eventCounts map[string]map[string]int
	// eventCounts = {
	//     "user123": {
	//         "created": 5,
	//         "updated": 3,
	//         "deleted": 1
	//     },
	//     "user456": {
	//         "created": 2,
	//         "updated": 0,
	//         "deleted": 2
	//     }
	// }
}

func (c *ConsumerStr) Created(ctx context.Context, uid string) error {
	if c.eventCounts[uid] == nil {
		c.eventCounts[uid] = make(map[string]int)
	}

	c.eventCounts[uid]["created"]++
	return nil
}

func (c *ConsumerStr) Updated(ctx context.Context, uid string) error {
	if c.eventCounts[uid] == nil {
		c.eventCounts[uid] = make(map[string]int)
	}

	c.eventCounts[uid]["updated"]++
	return nil
}

func (c *ConsumerStr) Deleted(ctx context.Context, uid string) error {
	if c.eventCounts[uid] == nil {
		c.eventCounts[uid] = make(map[string]int)
	}

	c.eventCounts[uid]["deleted"]++
	return nil
}

func main() {

}
