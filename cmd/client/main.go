package main

import (
	"fmt"
	"log"

	gamelogic "github.com/justinyeh1995/learn-pub-sub-starter/internal/gamelogic"
	pubsub "github.com/justinyeh1995/learn-pub-sub-starter/internal/pubsub"
	routing "github.com/justinyeh1995/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const connStr = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Panicf("Error connecting to RBMQ, %v", err)
	}
	defer conn.Close()
	log.Println("Connection Succeeds")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Panic(err)
	}
	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Panicf("Unable to establish a channel. %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	// create a new game state
	newGameState := gamelogic.NewGameState(username)
	pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+newGameState.GetUsername(),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(newGameState),
	)
	//REPL loop
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			log.Print("Spawning")
			err := newGameState.CommandSpawn(words)
			if err != nil {
				log.Printf("Could not spawn new command")
			}
		case "move":
			log.Print("Moving")
			move, err := newGameState.CommandMove(words)
			log.Printf("move: %v", move)
			if err != nil {
				log.Printf("Could not publish time: %v", err)
			}
		case "status":
			log.Print("Status")
			newGameState.CommandStatus()
		case "help":
			gamelogic.PrintServerHelp()
		case "spam":
			log.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Print("Invalid command")
		}
	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)

	// // Block until a signal is received.
	// s := <-signalChan
	// log.Println("Got signal:", s)
}
