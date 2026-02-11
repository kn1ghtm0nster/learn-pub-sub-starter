package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return
	}

	defer ch.Close()
	defer conn.Close()
	fmt.Println("Starting Peril server...")
	fmt.Println("Connected to RabbitMQ successfully!")

	gamelogic.PrintServerHelp()

	// Main loop to read console input and publish messages
	for {
		consoleInput := gamelogic.GetInput()

		if len(consoleInput) == 0 {
			continue
		}

		userInput := consoleInput[0]

		switch userInput {
			case "pause":
				fmt.Println("Pausing game...")
				err = pubsub.PublishJSON(
					ch,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: true,
					},
				)
				if err != nil {
					log.Printf("Failed to publish pause message: %v", err)
				}
			case "resume":
				fmt.Println("Resuming game...")
				err = pubsub.PublishJSON(
					ch,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: false,
					},
				)
				if err != nil {
					log.Printf("Failed to publish resume message: %v", err)
				}
			case "quit":
				fmt.Println("Exiting game...")
				return
			default:
				fmt.Printf("Unknown command: %s\n", userInput)
		}
	}
}
