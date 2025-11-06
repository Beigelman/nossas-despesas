package predict

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Beigelman/nossas-despesas/internal/pkg/env"
	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

type Client struct {
	httpClient *resty.Client
}

func NewClient(ctx context.Context, url string, environment env.Environment) (*Client, error) {
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	var tokenSource oauth2.TokenSource
	if environment == env.Production {
		credentials, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			return nil, fmt.Errorf("google.FindDefaultCredentials: %w", err)
		}

		tokenSource, err = idtoken.NewTokenSource(ctx, url, option.WithCredentials(credentials))
		if err != nil {
			return nil, fmt.Errorf("idtoken.NewTokenSource: %w", err)
		}
	} else {
		tokenSource = oauth2.StaticTokenSource(&oauth2.Token{
			TokenType:   "Bearer",
			AccessToken: "fake-token",
		})
	}

	return &Client{
		httpClient: resty.New().SetBaseURL(url).OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
			token, err := tokenSource.Token()
			if err != nil {
				return fmt.Errorf("tokenSource.Token: %w", err)
			}
			r.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
			return nil
		}),
	}, nil
}

type predictResponse struct {
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	AmountCents int    `json:"amount_cents"`
}

func (c *Client) ExpenseCategory(ctx context.Context, name string, amount int) (categoryID int, err error) {
	response, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"name":         name,
			"amount_cents": amount,
		}).
		Post("/predict")
	if err != nil {
		return 0, fmt.Errorf("failed to call predict endpoint: %w", err)
	}

	if response.IsError() {
		return 0, fmt.Errorf("failed to predict expense category with status code %d: %s", response.StatusCode(), response.String())
	}

	var predictResponse predictResponse
	if err := json.Unmarshal(response.Body(), &predictResponse); err != nil {
		return 0, fmt.Errorf("failed to unmarshal predict response: %w", err)
	}

	return predictResponse.CategoryID, nil
}
