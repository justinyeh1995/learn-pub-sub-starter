package main

import (
	gamelogic "github.com/justinyeh1995/learn-pub-sub-starter/internal/gamelogic"
	"github.com/justinyeh1995/learn-pub-sub-starter/internal/pubsub"
	"github.com/justinyeh1995/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishLog(
	rec gamelogic.RecognitionOfWar,
	gs *gamelogic.GameState,
	ch *amqp.Channel,
	log routing.GameLog,
) error {
	err := pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+rec.Attacker.Username,
		log,
	)
	if err != nil {
		return err
	}
	return nil
}
