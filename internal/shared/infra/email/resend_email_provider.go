package email

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type ResendEmailProvider struct {
	client *resend.Client
}

func NewResendEmailProvider(apiKey string) *ResendEmailProvider {
	return &ResendEmailProvider{client: resend.NewClient(apiKey)}
}

func (p *ResendEmailProvider) Send(ctx context.Context, email Email) error {
	params := &resend.SendEmailRequest{
		From:    email.From,
		To:      email.To,
		Html:    email.Html,
		Subject: email.Subject,
		Cc:      email.Cc,
		ReplyTo: email.ReplyTo,
	}

	if _, err := p.client.Emails.SendWithContext(ctx, params); err != nil {
		return fmt.Errorf("resend.Emails.SendWithContext: %w", err)
	}

	return nil
}
