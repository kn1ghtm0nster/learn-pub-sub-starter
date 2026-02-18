package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func (ps routing.PlayingState) {
		defer fmt.Print("> ")
		fmt.Println("pause state: ", ps)
		gs.HandlePause(ps)
	}
}

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {

	return func (am gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		fmt.Printf("move: %v\n", am)
		gs.HandleMove(am)
	}
}