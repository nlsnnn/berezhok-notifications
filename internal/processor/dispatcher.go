package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

type Processor interface {
	Process(ctx context.Context, msg NotificationMessage) error
}

type Dispatcher struct {
	processors map[string]Processor
	logger     *slog.Logger
}

func NewDispatcher(logger *slog.Logger) *Dispatcher {
	return &Dispatcher{
		processors: make(map[string]Processor),
		logger:     logger,
	}
}

func (d *Dispatcher) Register(msgType string, p Processor) {
	d.processors[msgType] = p
}

func (d *Dispatcher) Dispatch(ctx context.Context, body []byte) error {
	var msg NotificationMessage
	if err := json.Unmarshal(body, &msg); err != nil {
		return fmt.Errorf("unmarshal message: %w", err)
	}

	p, ok := d.processors[msg.Type]
	if !ok {
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}

	return p.Process(ctx, msg)
}
