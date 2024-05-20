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
		streamName    = "EVENTS"
		consumerName1 = "processor-1"
		url           = "nats://ninja:1234qwer@localhost:4222" //nats.DefaultURL
		nc            *nats.Conn
		stream        jetstream.Stream
		err           error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nc, err = nats.Connect(url)
	helper.HandleError(err)
	defer nc.Drain()
	{
		js, err := jetstream.New(nc)
		helper.HandleError(err)

		cfg := jetstream.StreamConfig{
			Name:      streamName,
			Retention: jetstream.WorkQueuePolicy,
			Subjects:  []string{"events.>"},
		}

		stream, err = js.CreateStream(ctx, cfg)
		helper.HandleError(err)
		fmt.Println("created the stream")
	}

	cons1, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name: consumerName1,
	})
	// msgs, _ := cons1.Fetch(3)
	// for msg := range msgs.Messages() {
	// 	log.Println("Received message: ", string(msg.Data()))
	// 	msg.DoubleAck(ctx)
	// }

	msgs, _ := cons1.Fetch(5, jetstream.FetchMaxWait(time.Second*30))
	for msg := range msgs.Messages() {
		log.Printf("Received message: %s: %s\n", msg.Subject(), string(msg.Data()))
		msg.DoubleAck(ctx)
	}

	stream.DeleteConsumer(ctx, consumerName1)
}
