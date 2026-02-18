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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
		return
	}
	defer publishCh.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to get username: %v", err)
		return
	}

	gameState := gamelogic.NewGameState(username)
	moveQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username)
	pauseQueueName := fmt.Sprintf("pause.%s", username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		moveQueueName,
		fmt.Sprintf("%s.*", routing.ArmyMovesPrefix),
		pubsub.SimpleQueueTransient,
		handlerMove(gameState),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to move messages: %v", err)
		return
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		pauseQueueName,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe to pause messages: %v", err)
		return
	}

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
				mv, err := gameState.CommandMove(consoleInput)
				if err != nil {
					fmt.Println(err)
					continue
				}
				err = pubsub.PublishJSON(
					publishCh,
					routing.ExchangePerilTopic,
					fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username),
					mv,
				)
				if err != nil {
					fmt.Println("Failed to publish move message:", err)
					continue
				}
				fmt.Println("move published successfully")
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
