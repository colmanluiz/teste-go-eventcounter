package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

var (
	conn    *amqp091.Connection
	channel *amqp091.Channel
)

type MessageBody struct {
	MessageID string `json:"id"`
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

	processedMessages := make(map[string]bool)

	for msg := range deliveries {
		parts := strings.Split(msg.RoutingKey, ".")
		userId := parts[0]
		eventType := parts[2]
		ctx := context.Background()

		mb := MessageBody{}
		err := json.Unmarshal(msg.Body, &mb)
		if err != nil {
			return err
		}

		if processedMessages[mb.MessageID] {
			log.Printf("message already processed, msg ID: %s", mb.MessageID)
			msg.Ack(false)
			continue
		}

		switch eventType {
		case "created":
			if err := consumer.Created(ctx, userId); err != nil {
				log.Printf("Failed to process message %s: %v", mb.MessageID, err)
				msg.Ack(false)
				continue
			}
		case "updated":
			if err := consumer.Updated(ctx, userId); err != nil {
				log.Printf("Failed to process message %s: %v", mb.MessageID, err)
				msg.Ack(false)
				continue
			}
		case "deleted":
			if err := consumer.Deleted(ctx, userId); err != nil {
				log.Printf("Failed to process message %s: %v", mb.MessageID, err)
				msg.Ack(false)
				continue
			}
		}

		processedMessages[mb.MessageID] = true
		msg.Ack(false)
	}

	return nil
}
