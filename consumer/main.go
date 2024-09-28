package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "Service1Queue"

func main() {
	// new connection
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Open channel
	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	// Subscribe to get messages from the queue
	messages, err := channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case message := <-messages:
			log.Printf("Message: %s\n", message.Body)
		case <-sigchan:
			log.Println("Interrupt detected!")
			os.Exit(0)
		}
	}
}
