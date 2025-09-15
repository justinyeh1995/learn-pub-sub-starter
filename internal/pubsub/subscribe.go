package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
) error {
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		log.Printf("Cannot declare a binding to the queue", err)
		return err
	}
	deliveries, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("Cannot create a amqp.Delivery", err)
		return err
	}
	go func() {
		defer ch.Close()
		// Range over the channel
		for delivery := range deliveries {
			// io.ReadAll(body)
			body := delivery.Body
			var data T
			// json.Decoder or json.Unmarshal??
			if err := json.Unmarshal(body, &data); err != nil {
				return
			}
			handler(data)
			// Acknowledge the message with delivery.Ack(false) to remove it from the queue
			delivery.Ack(false)
		}
	}()
	return nil
}
