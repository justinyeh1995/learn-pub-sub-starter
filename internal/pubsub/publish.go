package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	/**
	    Publish JSON Bytes to RabbitMQ channel
	**/
	valBytes, err := json.Marshal(val)
	if err != nil {
		log.Fatalf("Cannot marshal %v", val)
		return err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        valBytes,
	}
	ctx := context.Background()
	return ch.PublishWithContext(ctx, exchange, key, false, false, msg)
}

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(val)
	if err != nil {
		return err
	}
	valBytes := buffer.Bytes()

	msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        valBytes,
	}
	ctx := context.Background()
	return ch.PublishWithContext(ctx, exchange, key, false, false, msg)
}

// The closet way to do enum in go
type SimpleQueueType string

const (
	Durable   SimpleQueueType = "durable"
	Transient SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("%w", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		queueType == Durable,
		queueType != Durable,
		queueType != Durable,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "peril_dlx",
		})
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("%w", err)
	}

	if err := ch.QueueBind(
		q.Name,
		key,
		exchange,
		false,
		nil); err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("%w", err)
	}

	return ch, q, nil
}
