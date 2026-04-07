package processor

import (
	"context"
	"fmt"
)

type EmailProcessor struct {
	// Add any dependencies like email client, templates, etc.
}

func NewEmailProcessor() *EmailProcessor {
	return &EmailProcessor{}
}

func (p *EmailProcessor) Process(ctx context.Context, msg NotificationMessage) error {
	// Here you would implement the logic to send an email based on the msg.Recipient, msg.Template, and msg.Payload.
	// For example, you might use an email client library to send the email.
	fmt.Printf("Sending email to %s using template %s with payload: %s\n", msg.Recipient, msg.Template, string(msg.Payload))

	return nil
}
