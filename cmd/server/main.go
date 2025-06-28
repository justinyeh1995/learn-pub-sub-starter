package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    fmt.Println("Starting Peril server...")
    const connStr := "amqp://guest:guest@localhost:5672/"
    conn, err := amqp.Dial(connStr)
    if (err != nil) {
        log.Panicf("Unable to establish a connection.")
    }
    
    defer conn.Close()
    fmt.Println("Connection Succeeds")
    
    // How do you wait for a signal in go?
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt)
    
    // Block until a signal is received. 
    s := <-signalChan
    fmt.Println("Got signal:", s)
}
