package main

import (
	"fmt"
	"log"

	"github.com/justinyeh1995/learn-pub-sub-starter/internal/gamelogic"
	"github.com/justinyeh1995/learn-pub-sub-starter/internal/pubsub"
	"github.com/justinyeh1995/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connStr = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Panicf("Unable to establish a connection.")
	}

	defer conn.Close()
	fmt.Println("Connection Succeeds")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		"game_logs.*",
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("Cannot bind to queue. Error Message: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			log.Print("Pausing")
			data := routing.PlayingState{
				IsPaused: true,
			}
			err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				data,
			)
			if err != nil {
				log.Printf("Could not publish time: %v", err)
			}
		case "resume":
			log.Print("Resuming")
			data := routing.PlayingState{
				IsPaused: false,
			}
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				data,
			)
			if err != nil {
				log.Printf("Could not publish time: %v", err)
			}
		case "quit":
			log.Print("Quitting")
			return
		default:
			log.Print("Invalid command")
		}
	}

	// How do you wait for a signal in go?
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)

	// // Block until a signal is received.
	// s := <-signalChan
	// fmt.Println("Got signal:", s)
}
