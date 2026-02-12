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
	fmt.Println("Starting Peril client...")
	connStr := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to get username: %v", err)
		return
	}

	queueName := fmt.Sprintf("pause.%s", username)

	ch, _, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.SimpleQueueTransient)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
		return
	}

	defer ch.Close()

	gameState := gamelogic.NewGameState(username)

	for {
		consoleInput := gamelogic.GetInput()

		if len(consoleInput) == 0 {
			continue
		}

		userInput := consoleInput[0]

		switch userInput {
			case "spawn":
				err = gameState.CommandSpawn(consoleInput)
				if err != nil {
					fmt.Println(err)
				}
				continue
			case "move":
				message, err := gameState.CommandMove(consoleInput)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(message)
				continue
			case "status":
				gameState.CommandStatus()
				continue
			case "help":
				gamelogic.PrintClientHelp()
				continue
			case "spam":
				fmt.Println("Spamming not allowed yet!")
				continue
			case "quit":
				gamelogic.PrintQuit()
				return
			default:
				fmt.Printf("Unknown command: %s\n", userInput)
				continue
		}
	}
}
