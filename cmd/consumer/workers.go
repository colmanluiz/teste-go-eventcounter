package main

import (
	"context"
	"log"
	"sync"

	eventcounter "github.com/reb-felipe/eventcounter/pkg"
)

type EventMessage struct {
	Context context.Context
	UserID  string
}

func processCreatedEvents(createdChan chan EventMessage, consumer *eventcounter.ConsumerWrapper, wg *sync.WaitGroup) {
	defer wg.Done()
	for eventMsg := range createdChan {
		if err := consumer.Created(eventMsg.Context, eventMsg.UserID); err != nil {
			log.Printf("an error occured with the message from user %v: %s", eventMsg.UserID, err)
			continue
		}
	}
}

func processUpdatedEvents(updatedChan chan EventMessage, consumer *eventcounter.ConsumerWrapper, wg *sync.WaitGroup) {
	defer wg.Done()
	for eventMsg := range updatedChan {
		if err := consumer.Updated(eventMsg.Context, eventMsg.UserID); err != nil {
			log.Printf("an error occured with the message from user %v: %s", eventMsg.UserID, err)
			continue
		}
	}
}

func processDeletedEvents(deletedChan chan EventMessage, consumer *eventcounter.ConsumerWrapper, wg *sync.WaitGroup) {
	defer wg.Done()
	for eventMsg := range deletedChan {
		if err := consumer.Deleted(eventMsg.Context, eventMsg.UserID); err != nil {
			log.Printf("an error occured with the message from user %v: %s", eventMsg.UserID, err)
			continue
		}
	}
}
