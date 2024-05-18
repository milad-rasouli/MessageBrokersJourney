package main

import (
	"context"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/nats/helper"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	var (
		natsAddress = nats.DefaultURL // "nats://127.0.0.1:4222"
	)

	nc, err := nats.Connect(natsAddress)
	helper.HandleError(err)
	defer nc.Drain()

	js, err := jetstream.New(nc)
	helper.HandleError(err)

	streamName := "EVENTS"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	})

	<-time.After(50 * time.Second)
}
