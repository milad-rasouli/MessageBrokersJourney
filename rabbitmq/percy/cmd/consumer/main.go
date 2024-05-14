package main

import (
	"log"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/percy/internal"
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

	// the rabbitmq will keep sending the message till it get expired or receive back an ACK.
	// autoAck can be dangerous
	messageBus, err := client.Consume("customer_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	var blocking chan struct{}
	go func() {
		for message := range messageBus {
			log.Printf("new message: %v\n", message)
			if !message.Redelivered {
				err = message.Nack(false, true)
				if err != nil {
					log.Println(err)
					continue
				}
			}
			err = message.Ack(false)
			if err != nil {
				//panic(err)
				log.Println(err)
				continue
			}

			log.Printf("acknowledge message %s\n", message.MessageId)
		}
	}()
	<-blocking
}
