package pubsub

import (
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)


func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		amqp.Table{
			"x-dead-letter-exchange": routing.ExchangeDeadLetter,
		},
	)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range deliveryCh {
			target, err := unmarshaller(d.Body)
			if err != nil {
				log.Printf("Could not unmarshal: %v", err)
				d.Ack(false)
				continue
			}

			ackValue := handler(target)

			switch ackValue {
				case Ack:
					if err := d.Ack(false); err != nil {
						log.Fatalf("Failed to ack message: %v", err)
						return
					}
					log.Print("Message processed and acknowledged")
				case NackRequeue:
					if err := d.Nack(false, true); err != nil {
						log.Fatalf("Failed to nack and requeue message: %v", err)
						return
					}
					log.Print("Message processing failed, message requeued")
				case NackDiscard:
					if err := d.Nack(false, false); err != nil {
						log.Fatalf("Failed to nack and discard message: %v", err)
						return
					}
					log.Print("Message processing failed, message discarded")
			}
		}
	}()

	return nil
}