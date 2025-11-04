package predict

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClient *resty.Client
}

func NewClient(url string) *Client {
	return &Client{
		httpClient: resty.New().SetBaseURL(url),
	}
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
