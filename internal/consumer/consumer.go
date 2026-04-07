package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	maxRetries = 5
	retryDelay = time.Second * 2
)

type Dispatcher interface {
	Dispatch(ctx context.Context, body []byte) error
}

type Consumer struct {
	ch         *amqp.Channel
	queue      string
	dlqQueue   string
	dispatcher Dispatcher
	logger     *slog.Logger
}

func New(ch *amqp.Channel, dispatcher Dispatcher, logger *slog.Logger) *Consumer {
	return &Consumer{
		ch:         ch,
		queue:      queueNotificationsAll,
		dlqQueue:   queueDLQ,
		dispatcher: dispatcher,
		logger:     logger,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.ch.Consume(
		c.queue,
		"",    // consumer tag
		false, // ack вручную
		false, false, false, nil,
	)
	if err != nil {
		return err
	}

	c.logger.Info("consumer started", "queue", c.queue)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("consumer channel closed")
			}
			c.handle(ctx, msg)
		}
	}
}

func (c *Consumer) handle(ctx context.Context, msg amqp.Delivery) {
	var lastErr error

	for attempt := range maxRetries {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		if err := c.dispatcher.Dispatch(ctx, msg.Body); err != nil {
			lastErr = err
			c.logger.Warn("processing failed, retrying",
				"attempt", attempt+1,
				"max", maxRetries,
				"err", err,
			)
			continue
		}

		msg.Ack(false)
		return
	}

	c.logger.Error("message failed after all retries, sending to DLQ",
		"err", lastErr,
	)

	if err := c.publishToDLQ(msg); err != nil {
		c.logger.Error("failed to publish to DLQ", "err", err)
	}

	msg.Nack(false, false) // false, false. не возвращаем в очередь
}

func (c *Consumer) publishToDLQ(msg amqp.Delivery) error {
	return c.ch.Publish(
		"",         // default exchange
		c.dlqQueue, // routing key = queue name
		false, false,
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: amqp.Persistent,
			Headers:      msg.Headers, // original headers
		},
	)
}
