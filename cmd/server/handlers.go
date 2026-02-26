package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)


func handleGameLogMessage(log routing.GameLog) pubsub.AckType {
	defer fmt.Print("> ")

	err := gamelogic.WriteLog(log)
	if err != nil {
		fmt.Println("Error writing game log:", err)
		return pubsub.NackDiscard
	}
	return pubsub.Ack
}