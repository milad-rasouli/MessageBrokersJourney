package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/percy/internal"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var (
		err error
	)
	log.Println("the percy producer is running!")

	// all consuming will be done on this connection
	consumeConn, err := internal.ConnectRabbitMQWithTLS(
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
	defer consumeConn.Close()
	// tip:
	// you should recreate channel for each concurrent task, but reuse the connection!
	consumeClient, err := internal.NewRabbitMQClient(consumeConn)
	if err != nil {
		panic(err)
	}
	defer consumeClient.Close()

	producerConn, err := internal.ConnectRabbitMQWithTLS(
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
	defer producerConn.Close()
	// tip:
	// you should recreate channel for each concurrent task, but reuse the connection!
	producerClient, err := internal.NewRabbitMQClient(producerConn)
	if err != nil {
		panic(err)
	}
	defer producerClient.Close()

	// _, err = client.CreateQueue("customer_created", true, false)
	// if err != nil {
	// 	panic(err)
	// }

	// // it's for test
	// _, err = client.CreateQueue("customer_test", false, true)
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.CreateBinding("customer_created", "customer.created.*", "customer_test2")
	// if err != nil {
	// 	panic(err)
	// }

	// // it's for test
	// err = client.CreateBinding("customer_created", "customer.*", "customer_test2")
	// if err != nil {
	// 	panic(err)
	// }

	queue, err := consumeClient.CreateQueue("", true, true)
	if err != nil {
		panic(err)
	}

	err = consumeClient.CreateBinding(queue.Name, queue.Name, "customer_callback")
	if err != nil {
		panic(err)
	}
	messageBus, err := consumeClient.Consume(queue.Name, "customer-api", true)
	if err != nil {
		panic(err)
	}
	go func() {
		for message := range messageBus {
			log.Printf("message callback %s\n", message.CorrelationId)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {

		err = producerClient.Send(ctx, "customer_test2", "customer.created.us", amqp.Publishing{
			ContentType:   "text/plain",
			DeliveryMode:  amqp.Persistent,
			ReplyTo:       queue.Name,
			CorrelationId: fmt.Sprintf("customer_created_%d", i),
			Body:          []byte("An cool message between services"),
		})
		if err != nil {
			panic(err)
		}

		// // sending a transient message
		// err = client.Send(ctx, "customer_test2", "customer.test", amqp.Publishing{
		// 	ContentType:  "text/plain",
		// 	DeliveryMode: amqp.Transient,
		// 	Body:         []byte("An uncool undurable message between services"),
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// <-time.After(10 * time.Second)
	}

	var blocking chan struct{}
	<-blocking
}
