package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Mock do postgres.GetExpenses
type MockGetExpenses struct {
	ExpenseDetails []postgres.ExpenseDetails
	Error          error
}

func (m *MockGetExpenses) Execute(ctx context.Context, input postgres.GetExpensesInput) ([]postgres.ExpenseDetails, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ExpenseDetails, nil
}

func TestGetExpensesHandler(t *testing.T) {
	t.Parallel()

	now := time.Now()
	validExpenseDetails := []postgres.ExpenseDetails{
		{
			ID:           1,
			Name:         "Test Expense 1",
			Amount:       100.50,
			RefundAmount: nil,
			Description:  "Test description 1",
			CategoryID:   1,
			PayerID:      1,
			ReceiverID:   2,
			GroupID:      1,
			SplitRatio: postgres.SplitRatio{
				Payer:    50,
				Receiver: 50,
			},
			SplitType: "equal",
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: nil,
		},
		{
			ID:           2,
			Name:         "Test Expense 2",
			Amount:       200.75,
			RefundAmount: nil,
			Description:  "Test description 2",
			CategoryID:   1,
			PayerID:      1,
			ReceiverID:   2,
			GroupID:      1,
			SplitRatio: postgres.SplitRatio{
				Payer:    70,
				Receiver: 30,
			},
			SplitType: "proportional",
			CreatedAt: now.Add(-time.Hour),
			UpdatedAt: now.Add(-time.Hour),
			DeletedAt: nil,
		},
	}

	// Array with exactly 25 items to test pagination
	paginatedExpenses := make([]postgres.ExpenseDetails, 25)
	for i := 0; i < 25; i++ {
		paginatedExpenses[i] = postgres.ExpenseDetails{
			ID:           i + 1,
			Name:         fmt.Sprintf("Test Expense %d", i+1),
			Amount:       float32(100 + i),
			RefundAmount: nil,
			Description:  fmt.Sprintf("Test description %d", i+1),
			CategoryID:   1,
			PayerID:      1,
			ReceiverID:   2,
			GroupID:      1,
			SplitRatio: postgres.SplitRatio{
				Payer:    50,
				Receiver: 50,
			},
			SplitType: "equal",
			CreatedAt: now.Add(-time.Duration(i) * time.Hour),
			UpdatedAt: now.Add(-time.Duration(i) * time.Hour),
			DeletedAt: nil,
		}
	}

	// Definição dos casos de teste
	testCases := []struct {
		name               string
		groupID            int
		nextToken          string
		searchQuery        string
		mockExpenseDetails []postgres.ExpenseDetails
		mockError          error
		expectedStatus     int
		expectedResponse   string
		customAssertions   func(t *testing.T, body []byte)
	}{
		{
			name:               "should return 200 and expenses list",
			groupID:            1,
			nextToken:          "",
			searchQuery:        "",
			mockExpenseDetails: validExpenseDetails,
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.GetExpensesResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data.Expenses, 2)
				assert.Equal(t, 1, response.Data.Expenses[0].ID)
				assert.Equal(t, "Test Expense 1", response.Data.Expenses[0].Name)
				assert.Equal(t, float32(100.50), response.Data.Expenses[0].Amount)
				assert.Empty(t, response.Data.NextToken) // No next token since < 25 items
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 200 and expenses with search",
			groupID:            1,
			nextToken:          "",
			searchQuery:        "grocery",
			mockExpenseDetails: validExpenseDetails,
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.GetExpensesResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data.Expenses, 2)
				assert.Empty(t, response.Data.NextToken)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 200 and expenses with next token when limit reached",
			groupID:            1,
			nextToken:          "",
			searchQuery:        "",
			mockExpenseDetails: paginatedExpenses,
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.GetExpensesResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data.Expenses, 25)
				assert.NotEmpty(t, response.Data.NextToken) // Should have next token
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 200 and empty expenses list",
			groupID:            1,
			nextToken:          "",
			searchQuery:        "",
			mockExpenseDetails: []postgres.ExpenseDetails{},
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.GetExpensesResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data.Expenses, 0)
				assert.Empty(t, response.Data.NextToken)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 400 if invalid next token",
			groupID:            1,
			nextToken:          "invalid-token",
			searchQuery:        "",
			mockExpenseDetails: nil,
			mockError:          nil,
			expectedStatus:     400,
			expectedResponse:   `{"status_code":400,"message":"invalid next token","error":"invalid next token: internal=illegal base64 data at input byte 7"}`,
		},
		{
			name:               "should return 400 if next token contains invalid JSON",
			groupID:            1,
			nextToken:          "aW52YWxpZCBqc29u", // "invalid json" in base64
			searchQuery:        "",
			mockExpenseDetails: nil,
			mockError:          nil,
			expectedStatus:     400,
			expectedResponse:   `{"status_code":400,"message":"invalid next token","error":"invalid next token: internal=invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:               "should return 400 if next token is malformed base64",
			groupID:            1,
			nextToken:          "@@@@",
			searchQuery:        "",
			mockExpenseDetails: nil,
			mockError:          nil,
			expectedStatus:     400,
			expectedResponse:   `{"status_code":400,"message":"invalid next token","error":"invalid next token: internal=illegal base64 data at input byte 0"}`,
		},
		{
			name:               "should handle empty string next token correctly",
			groupID:            1,
			nextToken:          "",
			searchQuery:        "",
			mockExpenseDetails: validExpenseDetails,
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.GetExpensesResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data.Expenses, 2)
				assert.Empty(t, response.Data.NextToken)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 500 if database error",
			groupID:            1,
			nextToken:          "",
			searchQuery:        "",
			mockExpenseDetails: nil,
			mockError:          errors.New("database error"),
			expectedStatus:     500,
			expectedResponse:   `{"status_code":500,"message":"Internal Server Error","error":"query.GetExpenses: database error"}`,
		},
	}

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock
			mockGetExpenses := &MockGetExpenses{
				ExpenseDetails: tc.mockExpenseDetails,
				Error:          tc.mockError,
			}

			// Setup comum para todos os testes
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			getExpensesHandler := controller.NewGetExpenses(mockGetExpenses.Execute)

			app.Get("/expenses", func(c *fiber.Ctx) error {
				c.Locals("group_id", tc.groupID)
				return c.Next()
			}, getExpensesHandler)

			// Preparação da requisição
			url := "http://localhost:8080/expenses"
			params := []string{}

			if tc.nextToken != "" {
				params = append(params, "next_token="+tc.nextToken)
			}
			if tc.searchQuery != "" {
				params = append(params, "search="+tc.searchQuery)
			}

			if len(params) > 0 {
				url += "?" + strings.Join(params, "&")
			}

			req := httptest.NewRequest("GET", url, nil)

			// Execução da requisição
			resp, err := app.Test(req)
			assert.Nil(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			// Assertions
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Assertions específicas do caso
			if tc.customAssertions != nil {
				tc.customAssertions(t, body)
			} else {
				assert.Equal(t, tc.expectedResponse, string(body))
			}
			assert.NoError(t, resp.Body.Close())
		})
	}
}
