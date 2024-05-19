package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/nats/helper"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	var (
		url    = nats.DefaultURL
		nc     *nats.Conn
		js     jetstream.JetStream
		stream jetstream.Stream
		err    error
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	nc, err = nats.Connect(url)
	helper.HandleError(err)
	defer nc.Drain()
	{
		js, err = jetstream.New(nc)
		helper.HandleError(err)
		cfg := jetstream.StreamConfig{
			Name:      "EVENTS",
			Retention: jetstream.WorkQueuePolicy,
			Subjects:  []string{"events.>"},
		}

		stream, err = js.CreateStream(ctx, cfg)
		helper.HandleError(err)
	}
	fmt.Println("created the stream")

	js.Publish(ctx, "events.us.page_loaded", []byte("#1 message"))

	js.Publish(ctx, "events.eu.mouse_clicked", []byte("#2 message"))
	js.Publish(ctx, "events.us.input_focused", []byte("#3 message"))
	fmt.Println("published 3 messages")

	fmt.Println("# Stream info without any consumers")
	printStreamState(ctx, stream)

}

func printStreamState(ctx context.Context, stream jetstream.Stream) {
	info, _ := stream.Info(ctx)
	b, _ := json.MarshalIndent(info.State, "", " ")
	fmt.Println(string(b))
}
