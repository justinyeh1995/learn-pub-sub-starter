package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Block until a signal is received.
	s := <-signalChan
	log.Println("Got signal:", s)
}
