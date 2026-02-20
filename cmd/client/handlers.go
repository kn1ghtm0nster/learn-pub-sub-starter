package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func (ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		fmt.Println("pause state: ", ps)
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {

	return func (am gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		fmt.Printf("move: %v\n", am)
		moveOutcome:= gs.HandleMove(am)

		switch moveOutcome {
			case gamelogic.MoveOutComeSafe, gamelogic.MoveOutcomeMakeWar:
				return pubsub.Ack
			default:
				return pubsub.NackDiscard
		}
	}
}