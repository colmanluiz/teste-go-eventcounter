package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

var (
	conn    *amqp091.Connection
	channel *amqp091.Channel
)

type MessageBody struct {
	ID string `json:"id"`
}

func getChannel() (*amqp091.Channel, error) {
	var err error
	if channel != nil && !channel.IsClosed() {
		return channel, nil
	}

	conn, err := amqp091.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// TODO: Error handling

	if channel.IsClosed() {
		return nil, errors.New("channel is closed")
	}

	return channel, nil
}

func Declare() error { // declare the queue here because wants to make sure that the queue exists before trying to consume messages from it.
	channel, err := getChannel()
	if err != nil {
		return err
	}

	if err := channel.ExchangeDeclare(amqpExchange, "topic", true, false, false, false, nil); err != nil {
		return err
	}

	if _, err := channel.QueueDeclare(
		"eventcountertest", true, false, false, false, nil,
	); err != nil {
		return err
	}

	if err := channel.QueueBind("eventcountertest", "*.event.*", amqpExchange, false, nil); err != nil {
		return err
	}

	return nil
}

func Receive(consumer *ConsumerStr) error {
	channel, err := getChannel()
	if err != nil {
		return err
	}

	deliveries, err := channel.Consume("eventcountertest", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	processedMessages := make(map[string]bool) // implement later

	for msg := range deliveries {
		parts := strings.Split(msg.RoutingKey, ".")
		userId := parts[0]
		eventType := parts[2]

		mb := MessageBody{}
		err := json.Unmarshal(msg.Body, &mb)
		if err != nil {
			return err
		}

		switch eventType {
		case "created":
			consumer.Created(ctx, userId)
		case "updated":
			consumer.Updated(ctx, userId)
		case "deleted":
			consumer.Deleted(ctx, userId)
		}

		msg.Ack(false)
	}

	return nil
}
