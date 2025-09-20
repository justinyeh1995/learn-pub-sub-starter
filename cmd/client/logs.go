package main

import (
	"time"

	"github.com/justinyeh1995/learn-pub-sub-starter/internal/pubsub"
	"github.com/justinyeh1995/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(
	ch *amqp.Channel,
	username, msg string,
) error {
	err := pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		routing.GameLog{
			Username:    username,
			CurrentTime: time.Now(),
			Message:     msg,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
