package main

import (
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)




func CreateGameLog(ch *amqp.Channel, exchange, routingKey, message, username string) error {
	gameLog := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     message,
		Username:    username,
	}

	return pubsub.PublishGob(
		ch,
		exchange,
		routingKey,
		gameLog,
	)
}