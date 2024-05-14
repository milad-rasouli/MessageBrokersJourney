package internal

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// the connection used by the client
	conn *amqp.Connection
	// the channel is used to process/ send message
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}
	// for waiting for ack need to active the confirm mode
	err = ch.Confirm(false)
	if err != nil {
		return RabbitClient{}, err
	}
	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	return err
}

func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving nowait false, will make the channel return an error if it fails.
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	// return rc.ch.PublishWithContext(ctx,
	// 	exchange,
	// 	routingKey,
	// 	true,  // Mandatory is used to determine an error should be returned upon failure
	// 	false, //immediate
	// 	options,
	// )
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx,
		exchange,
		routingKey,
		true,  // Mandatory is used to determine an error should be returned upon failure
		false, //immediate
		options,
	)
	if err != nil {
		return err
	}
	log.Println(confirmation.Wait())
	return nil

}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rc RabbitClient) Qos(prefetchCount, prefetchSize int, global bool) error {
	return rc.ch.Qos(prefetchCount, prefetchSize, global)
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}
