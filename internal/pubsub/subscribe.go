package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(body []byte) (T, error) {
			var data T
			// json.Decoder or json.Unmarshal??
			if err := json.Unmarshal(body, &data); err != nil {
				return data, err
			}
			return data, nil
		},
	)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(body []byte) (T, error) {
			buffer := bytes.NewBuffer(body)
			dec := gob.NewDecoder(buffer)
			var data T
			err := dec.Decode(&data)
			if err != nil {
				return data, err
			}
			return data, nil
		},
	)
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
	)
	if err != nil {
		log.Printf("Cannot declare a binding to the queue %v", err)
		return err
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		log.Printf("Cannot set QoS, %v", err)
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
		log.Printf("Cannot create a amqp.Delivery %v", err)
		return err
	}
	go func() {
		defer ch.Close()
		// Range over the channel
		for delivery := range deliveries {
			// io.ReadAll(body)
			body := delivery.Body
			data, err := unmarshaller(body)
			if err != nil {
				return
			}
			ackType := handler(data)
			// Depending on the returned "ackType", call the appropriate acknowledgment method.
			switch ackType {
			case Ack:
				log.Println("Ack")
				delivery.Ack(false)
			case NackRequeue:
				log.Println("NackRequeue")
				delivery.Nack(false, true)
			case NackDiscard:
				log.Println("NackDiscard")
				delivery.Nack(false, false)
			}
		}
	}()
	return nil
}
