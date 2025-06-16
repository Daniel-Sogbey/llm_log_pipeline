package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type PubSub struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewPubSubConnection(connectionURL string) (*PubSub, error) {
	var connection *amqp.Connection
	for {
		conn, err := amqp.Dial(connectionURL)
		if err == nil {
			connection = conn
			break
		}

		time.Sleep(5 * time.Second)
		log.Println("Retrying in 5 seconds ...")
	}

	channel, err := connection.Channel()

	if err != nil {
		return nil, err
	}

	return &PubSub{connection: connection, channel: channel}, nil
}

func (ps *PubSub) ExchangeDeclare(name, kind string, durable bool) error {
	err := ps.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("Created Exchange %s of Kind %s with Durability set to %t \n", name, kind, durable)

	return nil
}

func (ps *PubSub) QueueDeclare(name string, durable bool) (*amqp.Queue, error) {
	queue, err := ps.channel.QueueDeclare(
		name,
		durable,
		!durable,
		!durable,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	log.Printf("Created Queue %s with Durability set to %t \n", name, durable)

	return &queue, nil
}

func (ps *PubSub) QueueBind(name, key, exchange string) error {
	err := ps.channel.QueueBind(
		name,
		key,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	log.Printf("Binding Queue %s To Exchange %s\n", name, exchange)

	return nil
}
