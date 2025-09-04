package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	eventcounter "github.com/reb-felipe/eventcounter/pkg"
)

var (
	conn    *amqp091.Connection
	channel *amqp091.Channel
)

type MessageBody struct {
	MessageID string `json:"id"`
}

func getChannel() (*amqp091.Channel, error) {
	if conn != nil && !conn.IsClosed() && channel != nil && !channel.IsClosed() {
		return channel, nil
	}

	var err error
	conn, err = amqp091.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}

	channel, err = conn.Channel()
	if err != nil {
		return nil, err
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

func Receive(consumer *eventcounter.ConsumerWrapper) error {
	var wg sync.WaitGroup

	createdChan := make(chan EventMessage)
	updatedChan := make(chan EventMessage)
	deletedChan := make(chan EventMessage)

	channel, err := getChannel()
	if err != nil {
		return err
	}

	deliveries, err := channel.Consume("eventcountertest", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	processedMessages := make(map[string]bool)

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	wg.Add(3)
	go processCreatedEvents(createdChan, consumer, &wg)
	go processUpdatedEvents(updatedChan, consumer, &wg)
	go processDeletedEvents(deletedChan, consumer, &wg)

	for {
		select {
		case msg := <-deliveries:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(5 * time.Second)

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
				createdChan <- EventMessage{Context: ctx, UserID: userId}
			case "updated":
				updatedChan <- EventMessage{Context: ctx, UserID: userId}
			case "deleted":
				deletedChan <- EventMessage{Context: ctx, UserID: userId}
			}

			processedMessages[mb.MessageID] = true
			msg.Ack(false)
		case <-timer.C:
			close(createdChan)
			close(updatedChan)
			close(deletedChan)

			wg.Wait()
			return nil
		}
	}
}
