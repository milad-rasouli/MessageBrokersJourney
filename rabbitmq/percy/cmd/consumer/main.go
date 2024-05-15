package main

import (
	"context"
	"log"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/percy/internal"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		err error
	)
	// instead of the hard coded items use env in real product
	conn, err := internal.ConnectRabbitMQWithTLS(
		"ninja",
		"1234qwer",
		"localhost:5671",
		"customer",
		"/home/milad/Documents/MessageBrokersJourney/tls-gen/basic/result/ca_certificate.pem",
		"/home/milad/Documents/MessageBrokersJourney/tls-gen/basic/result/client_milad_certificate.pem",
		"/home/milad/Documents/MessageBrokersJourney/tls-gen/basic/result/client_milad_key.pem",
	)
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

	publishConn, err := internal.ConnectRabbitMQWithTLS(
		"ninja",
		"1234qwer",
		"localhost:5671",
		"customer",
		"/home/milad/Documents/MessageBrokersJourney/tls-gen/basic/result/ca_certificate.pem",
		"/home/milad/Documents/MessageBrokersJourney/tls-gen/basic/result/client_milad_certificate.pem",
		"/home/milad/Documents/MessageBrokersJourney/tls-gen/basic/result/client_milad_key.pem",
	)
	if err != nil {
		panic(err)
	}
	defer publishConn.Close()
	publishClient, err := internal.NewRabbitMQClient(publishConn)
	if err != nil {
		panic(err)
	}
	defer publishClient.Close()
	err = client.Qos(5, 0, false)
	if err != nil {
		log.Fatalf("Failed to set QoS: %s", err)
	}

	// leave the queue name blank, the rabbitmq will generate you a random name
	queue, err := client.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	err = client.CreateBinding(queue.Name, "", "customer_test2")
	if err != nil {
		panic(err)
	}

	// the rabbitmq will keep sending the message till it get expired or receive back an ACK.
	// autoAck can be dangerous
	messageBus, err := client.Consume(queue.Name, "email-service", false)
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
				log.Printf("new message: %v", msg)
				<-time.After(3 * time.Second)
				err = msg.Ack(false)
				if err != nil {
					log.Println("ack failed")
					return err
				}

				err = publishClient.Send(ctx, "customer_callback", msg.ReplyTo, amqp091.Publishing{
					ContentType:   "text/plain",
					DeliveryMode:  amqp091.Persistent,
					Body:          []byte("RPC COMPLETE"),
					CorrelationId: msg.CorrelationId, // the publisher will know what message we got
				})
				if err != nil {
					panic(err)
				}
				log.Printf("acknowledged message %s\n", msg.MessageId)
				return nil
			})
		}
	}()
	log.Println("consuming use CTRL+C")
	<-blocking
}
