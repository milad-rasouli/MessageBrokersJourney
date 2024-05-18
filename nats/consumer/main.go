package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/nats/helper"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	var (
		natsAddress = nats.DefaultURL // "nats://127.0.0.1:4222"
	)
	log.Println("Consumer is running.")

	nc, err := nats.Connect(natsAddress)
	helper.HandleError(err)
	defer nc.Drain()

	js, err := jetstream.New(nc)
	helper.HandleError(err)

	streamName := "EVENTS"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"events.>"},
	})
	helper.HandleError(err)

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{})
	helper.HandleError(err)

	fetchStart := time.Now()
	msgs, err := cons.Fetch(1, jetstream.FetchMaxWait(time.Second))
	helper.HandleError(err)
	i := 0
	for msg := range msgs.Messages() {
		log.Printf("message: %s", msg.Data())
		msg.Ack()
		i++
	}

	fmt.Printf("got %d messages in %v\n", i, time.Since(fetchStart))

	//<-time.After(50 * time.Second)
}
