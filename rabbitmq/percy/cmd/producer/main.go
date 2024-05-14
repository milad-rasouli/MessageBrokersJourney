package main

import (
	"context"
	"log"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/percy/internal"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var (
		err error
	)
	log.Println("percy producer is running!")
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

	err = client.CreateQueue("customer_created", true, false)
	if err != nil {
		panic(err)
	}

	// it's for test
	err = client.CreateQueue("customer_test", false, true)
	if err != nil {
		panic(err)
	}

	err = client.CreateBinding("customer_created", "customer.created.*", "customer_test2")
	if err != nil {
		panic(err)
	}

	// it's for test
	err = client.CreateBinding("customer_created", "customer.*", "customer_test2")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = client.Send(ctx, "customer_test2", "customer.created.us", amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         []byte("An cool message between services"),
	})
	if err != nil {
		panic(err)
	}

	// sending a transient message
	err = client.Send(ctx, "customer_test2", "customer.test", amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Transient,
		Body:         []byte("An uncool undurable message between services"),
	})
	if err != nil {
		panic(err)
	}

	<-time.After(10 * time.Second)

	// To send the data we need to make an queue
	log.Println(client)
}
