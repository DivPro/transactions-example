package consumer

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"golang.org/x/exp/slog"
)

type Handler[T any] interface {
	Handle(v T) error
}

type Consumer[T any] struct {
	logger  *slog.Logger
	handler Handler[T]
}

func NewConsumer[T any](
	logger *slog.Logger,
	handler Handler[T],
) *Consumer[T] {
	return &Consumer[T]{
		logger:  logger,
		handler: handler,
	}
}

func (c *Consumer[_]) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer[_]) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			session.MarkMessage(message, "")
			var v T
			if err := json.Unmarshal(message.Value, &v); err != nil {
				c.logger.Error("unmarshal msg", err, string(message.Value))
			}
			if err := c.handler.Handle(v); err != nil {
				c.logger.Error("process msg with %T", err, c.handler, v)
			}
		case <-session.Context().Done():
			return nil
		}
	}
}
