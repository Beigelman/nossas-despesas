package email

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type MailTrapEmailProvider struct {
	client *resty.Client
	apiKey string
}

type email struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type MailTrapBody struct {
	From    email
	To      []email
	Subject string
	Html    string
}

func NewMailTrapEmailProvider(apiKey string) *MailTrapEmailProvider {
	return &MailTrapEmailProvider{client: resty.New(), apiKey: apiKey}
}

func (p *MailTrapEmailProvider) Send(ctx context.Context, email email) error {
	var to []email
	for _, e := range email.To {
		to = append(to, email{Email: e, Name: "test"})
	}

	req, err := p.client.R().
		SetContext(ctx).
		SetHeader("Api-Token", p.apiKey).
		SetBody(&MailTrapBody{
			From: email{
				Email: email.From,
				Name:  "Nossas Despesas",
			},
			To:      to,
			Subject: email.Subject,
			Html:    email.Html,
		}).Post("https://sandbox.api.mailtrap.io/api/send/2649128")

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if req.StatusCode() != 200 {
		return fmt.Errorf("failed to send email: %s", req.Body())
	}

	return nil
}
