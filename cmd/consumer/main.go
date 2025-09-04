package main

import (
	"flag"
	"log"

	eventcounter "github.com/reb-felipe/eventcounter/pkg"
)

var (
	amqpUrl      string
	amqpExchange string
)

func init() {
	flag.StringVar(&amqpUrl, "amqp-url", "amqp://guest:guest@localhost:5672/", "RabbitMQ URL")
	flag.StringVar(&amqpExchange, "amqp-exchange", "eventcountertest", "RabbitMQ Exchange")
}

func main() {
	flag.Parse()

	consumer := eventcounter.NewConsumerWrapper()

	if err := Declare(); err != nil {
		log.Printf("can`t declare queue or exchange, err: %s", err.Error())
	}

	if err := Receive(consumer); err != nil {
		log.Printf("can't consume any message, err: %s", err.Error())
	}

	if err := WriteJSONOutput(consumer, "./output"); err != nil {
		log.Printf("failed to write output: %v", err)
	}
}
