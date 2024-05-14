package main

import (
	"context"
	"log"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/percy/internal"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		err error
	)
	conn, err := internal.ConnectRabbitMQ("ninja", "1234qwer", "localhost:5672", "customer")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// tip:
	// you should recreate channel for each concurrent task, but reuse the connection!
	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// err = client.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to set QoS: %s", err)
	// }

	// the rabbitmq will keep sending the message till it get expired or receive back an ACK.
	// autoAck can be dangerous
	messageBus, err := client.Consume("customer_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	var blocking chan struct{}
	// go func() {
	// 	for message := range messageBus {
	// 		log.Printf("new message: %v\n", message)
	// 		if !message.Redelivered {
	// 			err = message.Nack(false, true)
	// 			if err != nil {
	// 				log.Println(err)
	// 				continue
	// 			}
	// 		}
	// 		err = message.Ack(false)
	// 		if err != nil {
	// 			//panic(err)
	// 			log.Println(err)
	// 			continue
	// 		}

	// 		log.Printf("acknowledge message %s\n", message.MessageId)
	// 	}
	// }()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {
				log.Printf("new message: %s", msg)
				<-time.After(1 * time.Second)
				err = msg.Ack(false)
				if err != nil {
					log.Println("ack failed")
					return err
				}
				log.Printf("acknowledged message %s\n", msg.MessageId)
				return nil
			})
		}
	}()
	log.Println("consuming use CTRL+C")
	<-blocking
}
