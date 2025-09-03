package main

import (
	"errors"

	"github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp091.Connection
	channel *amqp091.Channel
)

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
