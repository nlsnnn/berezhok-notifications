package consumer

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeNotifications = "notifications"
	exchangeTypeTopic     = "topic"

	queueNotificationsAll = "notifications.all"
	queueDLQ              = "dlq.notifications"

	bindingKeyAll = "notification.*"
)

type Conn struct {
	Channel *amqp.Channel
}

// GetConn establishes a connection to RabbitMQ and returns a Conn struct containing the channel
func GetConn(rabbitURL string) (Conn, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return Conn{}, err
	}

	ch, err := conn.Channel()
	return Conn{
		Channel: ch,
	}, err
}

func (c *Conn) Setup() error {
	// Declare main exchange
	err := c.Channel.ExchangeDeclare(
		exchangeNotifications,
		exchangeTypeTopic,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	// Declare DLQ
	_, err = c.Channel.QueueDeclare(
		queueDLQ,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("declare dlq: %w", err)
	}

	// Declare main queue
	_, err = c.Channel.QueueDeclare(
		queueNotificationsAll,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	// Bind queue to exchange
	err = c.Channel.QueueBind(
		queueNotificationsAll,
		bindingKeyAll, // routing key pattern
		exchangeNotifications,
		false, nil,
	)
	if err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	return nil
}

func (c *Conn) Close() error {
	if err := c.Channel.Close(); err != nil {
		return err
	}
	return nil
}
