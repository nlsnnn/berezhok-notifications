package processor

import (
	"encoding/json"
	"time"
)

type NotificationMessage struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Recipient string          `json:"recipient"`
	Template  string          `json:"template"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

const (
	TypeEmail = "email"
	TypeSMS   = "sms"
	TypePush  = "push"
)
