package main

import (
	"context"
	"flag"
)

var (
	amqpUrl string
	amqpExchange string
)

func init() {
	flag.StringVar(&amqpUrl, "amqp-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	flag.StringVar(&amqpExchange, "amqp-exchange", "eventcountertest", "RabbitMQ Exchange")
}

type ConsumerStr struct {}

func (c *ConsumerStr) Created(ctx context.Context, uid string) (context.Context, string) {
	return ctx, uid
}

func (c *ConsumerStr) Updated(ctx context.Context, uid string) (context.Context, string) {
	return ctx, uid
}

func (c *ConsumerStr) Deleted(ctx context.Context, uid string) (context.Context, string) {
	return ctx, uid
}

func main() {

}
