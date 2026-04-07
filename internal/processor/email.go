package processor

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	TemplateApprove = "application-approved"
	TemplateReject  = "registration-rejection"
)

type EmailProcessor struct {
	emailSender EmailSender
}

type EmailSender interface {
	SendEmail(to string, templateId string, variables map[string]interface{}) error
}

func NewEmailProcessor(emailSender EmailSender) *EmailProcessor {
	return &EmailProcessor{emailSender: emailSender}
}

func (p *EmailProcessor) Process(ctx context.Context, msg NotificationMessage) error {
	fmt.Printf("Sending email to %s using template %s with payload: %s\n", msg.Recipient, msg.Template, string(msg.Payload))

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("failed to parse payload: %w", err)
	}

	variables := map[string]interface{}{
		"login":    msg.Recipient,
		"password": payload["password"],
	}

	err := p.emailSender.SendEmail(msg.Recipient, msg.Template, variables)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
