package main

import (
	"log"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/helper"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("RabbitMQ consumer is running")
	connStr := "amqp://test:1234qwer@127.0.0.1:5672/cc-dev-vhost"

	conn, err := amqp.Dial(connStr)
	helper.FailOnError(err, "failed to connect to rabbitmq")
	defer conn.Close()

	ch, err := conn.Channel()
	helper.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	helper.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	helper.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
