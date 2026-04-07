package adapters

import (
	"fmt"

	"github.com/resend/resend-go/v3"
)

type ResendClient struct {
	From   string
	Client *resend.Client
}

func NewResendClient(apiKey string, from string) *ResendClient {
	client := resend.NewClient(apiKey)
	return &ResendClient{From: from, Client: client}
}

func (c *ResendClient) SendEmail(to string, templateId string, variables map[string]interface{}) error {
	params := &resend.SendEmailRequest{
		To: []string{to},
		Template: &resend.EmailTemplate{
			Id:        templateId,
			Variables: variables,
		},
	}

	sent, err := c.Client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	fmt.Printf("email sent: %s", sent.Id)
	return nil
}
