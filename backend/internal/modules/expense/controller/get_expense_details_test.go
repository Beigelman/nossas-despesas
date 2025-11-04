package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
)

// Mock do postgres.GetExpenseDetails
type MockGetExpenseDetails struct {
	ExpenseDetails []postgres.ExpenseDetails
	Error          error
}

func (m *MockGetExpenseDetails) Execute(ctx context.Context, expenseID int) ([]postgres.ExpenseDetails, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ExpenseDetails, nil
}

func TestGetExpenseDetailsHandler(t *testing.T) {
	t.Parallel()

	validExpenseDetails := []postgres.ExpenseDetails{
		{
			ID:           1,
			Name:         "Test Expense",
			Amount:       100.50,
			RefundAmount: nil,
			Description:  "Test description",
			CategoryID:   1,
			PayerID:      1,
			ReceiverID:   2,
			GroupID:      1,
			SplitRatio: postgres.SplitRatio{
				Payer:    50,
				Receiver: 50,
			},
			SplitType: "equal",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	expenseDetailsWithRefund := []postgres.ExpenseDetails{
		{
			ID:           2,
			Name:         "Test Expense with Refund",
			Amount:       200.75,
			RefundAmount: func() *float32 { v := float32(50.25); return &v }(),
			Description:  "Test description with refund",
			CategoryID:   1,
			PayerID:      1,
			ReceiverID:   2,
			GroupID:      1,
			SplitRatio: postgres.SplitRatio{
				Payer:    70,
				Receiver: 30,
			},
			SplitType: "proportional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	expenseDetailsWrongGroup := []postgres.ExpenseDetails{
		{
			ID:           3,
			Name:         "Test Expense Wrong Group",
			Amount:       300.00,
			RefundAmount: nil,
			Description:  "Test description",
			CategoryID:   1,
			PayerID:      1,
			ReceiverID:   2,
			GroupID:      999, // Grupo diferente
			SplitRatio: postgres.SplitRatio{
				Payer:    50,
				Receiver: 50,
			},
			SplitType: "equal",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	// Definição dos casos de teste
	testCases := []struct {
		name               string
		expenseID          string
		groupID            int
		mockExpenseDetails []postgres.ExpenseDetails
		mockError          error
		expectedStatus     int
		expectedResponse   string
		customAssertions   func(t *testing.T, body []byte)
	}{
		{
			name:               "should return 200 and expense details",
			expenseID:          "1",
			groupID:            1,
			mockExpenseDetails: validExpenseDetails,
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpenseDetails]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 1)
				assert.Equal(t, 1, response.Data[0].ID)
				assert.Equal(t, "Test Expense", response.Data[0].Name)
				assert.Equal(t, float32(100.50), response.Data[0].Amount)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 200 and expense details with refund",
			expenseID:          "2",
			groupID:            1,
			mockExpenseDetails: expenseDetailsWithRefund,
			mockError:          nil,
			expectedStatus:     200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpenseDetails]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 1)
				assert.Equal(t, 2, response.Data[0].ID)
				assert.Equal(t, "Test Expense with Refund", response.Data[0].Name)
				assert.Equal(t, float32(200.75), response.Data[0].Amount)
				assert.NotNil(t, response.Data[0].RefundAmount)
				assert.Equal(t, float32(50.25), *response.Data[0].RefundAmount)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:               "should return 400 if expense_id is not a valid number",
			expenseID:          "invalid",
			groupID:            1,
			mockExpenseDetails: nil,
			mockError:          nil,
			expectedStatus:     400,
			expectedResponse:   `{"status_code":400,"message":"invalid expense id","error":"invalid expense id"}`,
		},
		{
			name:               "should return 404 if expense is not found",
			expenseID:          "999",
			groupID:            1,
			mockExpenseDetails: []postgres.ExpenseDetails{}, // Empty slice
			mockError:          nil,
			expectedStatus:     404,
			expectedResponse:   `{"status_code":404,"message":"expense not found","error":"expense not found"}`,
		},
		{
			name:               "should return 403 if group mismatch",
			expenseID:          "3",
			groupID:            1,
			mockExpenseDetails: expenseDetailsWrongGroup,
			mockError:          nil,
			expectedStatus:     403,
			expectedResponse:   `{"status_code":403,"message":"group mismatch","error":"group mismatch"}`,
		},
		{
			name:               "should return 500 if database error",
			expenseID:          "1",
			groupID:            1,
			mockExpenseDetails: nil,
			mockError:          errors.New("database error"),
			expectedStatus:     500,
			expectedResponse:   `{"status_code":500,"message":"Internal Server Error","error":"query.GetExpenseDetails: database error"}`,
		},
	}

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock
			mockGetExpenseDetails := &MockGetExpenseDetails{
				ExpenseDetails: tc.mockExpenseDetails,
				Error:          tc.mockError,
			}

			// Setup comum para todos os testes
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			getExpenseDetailsHandler := controller.NewGetExpenseDetails(mockGetExpenseDetails.Execute)

			app.Get("/expenses/:expense_id/details", func(c *fiber.Ctx) error {
				c.Locals("group_id", tc.groupID)
				return c.Next()
			}, getExpenseDetailsHandler)

			// Preparação da requisição
			req := httptest.NewRequest("GET", "http://localhost:8080/expenses/"+tc.expenseID+"/details", nil)

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
