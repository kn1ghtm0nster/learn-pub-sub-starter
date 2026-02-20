package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	dat, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		context.Background(), 
		exchange, 
		key, 
		false, 
		false, 
		amqp.Publishing{
			ContentType: "application/json",
			Body: 		 dat,
	})
}

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
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
			var rawBytes T
			if err := json.Unmarshal(d.Body, &rawBytes); err != nil {
				if err := d.Ack(false); err != nil {
					log.Fatalf("Failed to ack message after unmarshal error: %v", err)
					return
				}
				continue
			}

			ackValue := handler(rawBytes)
			
			switch ackValue {
				case Ack:
					if err := d.Ack(false); err != nil {
						log.Fatalf("Failed to ack message: %v", err)
						return
					}
					log.Printf("Message processed and acknowledged: %s", d.Body)
				case NackRequeue:
					if err := d.Nack(false, true); err != nil {
						log.Fatalf("Failed to nack and requeue message: %v", err)
						return
					}
					log.Printf("Message processing failed, message requeued: %s", d.Body)
				case NackDiscard:
					if err := d.Nack(false, false); err != nil {
						log.Fatalf("Failed to nack and discard message: %v", err)
						return
					}
					log.Printf("Message processing failed, message discarded: %s", d.Body)
			}
		}
	}()

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	durable := false
	autoDelete := false
	exclusive := false
	noWait := false

	if queueType == SimpleQueueDurable{
		durable = true
	}

	if queueType == SimpleQueueTransient {
		autoDelete = true
		exclusive = true
	}

	queue, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = ch.QueueBind(queueName, key, exchange, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return ch, queue, nil
}
