package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func (ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		fmt.Println("pause state: ", ps)
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {

	return func (am gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		fmt.Printf("move: %v\n", am)
		moveOutcome:= gs.HandleMove(am)

		switch moveOutcome {
			case gamelogic.MoveOutComeSafe:
				return pubsub.Ack
			case gamelogic.MoveOutcomeMakeWar:
				username := gs.GetPlayerSnap().Username
				routingKey := fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, username)
				publishData := gamelogic.RecognitionOfWar{
					Attacker: am.Player,
					Defender: gs.GetPlayerSnap(),
				}
				err := pubsub.PublishJSON(
					ch,
					routing.ExchangePerilTopic,
					routingKey,
					publishData,
				)
				if err != nil {
					log.Fatalf("Error publishing data: %v", err)
				}
				return pubsub.NackRequeue
			default:
				return pubsub.NackDiscard
		}
	}
}

func handlerWarMessages(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func (rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, _, _ := gs.HandleWar(rw)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			return pubsub.Ack
		default:
			fmt.Printf("Unknown war outcome: %v\n", outcome)
			return pubsub.NackDiscard
		}
	}
}