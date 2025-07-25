package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	/**
	    Publish JSON Bytes to RabbitMQ channel
	**/
	valBytes, err := json.Marshal(val)
	if err != nil {
		log.Fatalf("Cannot marshal %s", val)
		return err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        valBytes,
	}
	ctx := context.Background()
	ch.PublishWithContext(ctx, exchange, key, false, false, msg)

	return nil
}
