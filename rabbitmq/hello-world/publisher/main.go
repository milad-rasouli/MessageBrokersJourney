package main

import (
	"context"
	"log"
	"time"

	"github.com/Milad75Rasouli/MessageBrokersJourney/rabbitmq/helper"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Println("RabbitMQ publisher is running")
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Greeting!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	helper.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
