package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Parallel()

	validBodyReq := controller.CreateExpenseRequest{
		Name:        "Test Expense",
		Amount:      100,
		Description: "My first expense",
		CategoryID:  1,
		SplitType:   "equal",
		PayerID:     1,
		ReceiverID:  1,
	}

	invalidBodyReq := map[string]any{
		"group_id":    1,
		"name":        "Test Expense",
		"amount":      100,
		"description": "My first expense",
		"category_id": 1,
		"payer_id":    1,
		"receiver_id": 2,
	}

	newExpense, _ := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "Test Expense",
		Amount:      100,
		Description: "My first expense",
		GroupID:     group.ID{Value: 1},
		CategoryID:  category.ID{Value: 1},
		SplitRatio:  expense.SplitRatio{Payer: 50, Receiver: 50},
		PayerID:     user.ID{Value: 1},
		ReceiverID:  user.ID{Value: 2},
	})

	// Definição dos casos de teste
	testCases := []struct {
		name             string
		body             any
		mockSetup        func(usecase *mocks.MockusecaseCreateExpense)
		expectedStatus   int
		expectedResponse string
		customAssertions func(t *testing.T, body []byte)
	}{
		{
			name: "should return 201 and create a new expense",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecaseCreateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(newExpense, nil).Once()
			},
			expectedStatus: 201,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.CreateExpenseResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 201, response.StatusCode)
				assert.Equal(t, 1, response.Data.ID)
				assert.Equal(t, "Test Expense", response.Data.Name)
				assert.Equal(t, float32(1), response.Data.Amount)
				assert.Equal(t, 1, response.Data.PayerID)
				assert.Equal(t, 2, response.Data.ReceiverID)
				assert.NotEmpty(t, response.Date) // Verifica se a data existe
			},
		},
		{
			name:             "should return 400 if request body is invalid",
			body:             invalidBodyReq,
			mockSetup:        func(usecase *mocks.MockusecaseCreateExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid request body","error":"invalid request body: internal=validation errors: [SplitType]: '' | Needs to implement 'required'"}`,
		},
		{
			name: "should return 422 if request body is not processable",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecaseCreateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, except.UnprocessableEntityError()).Once()
			},
			expectedStatus:   422,
			expectedResponse: `{"status_code":422,"message":"Unprocessable Entity","error":"CreateExpense: Unprocessable Entity"}`,
		},
		{
			name: "should return 500 if it gets an unexpected error",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecaseCreateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, errors.New("unexpected error")).Once()
			},
			expectedStatus:   500,
			expectedResponse: `{"status_code":500,"message":"Internal Server Error","error":"CreateExpense: unexpected error"}`,
		},
	}

	// Setup comum para todos os testes
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

	createExpense := mocks.NewMockusecaseCreateExpense(t)
	createExpenseHandler := controller.NewCreateExpense(createExpense.Execute)
	app.Post("/expenses", func(c *fiber.Ctx) error {
		c.Locals("group_id", 1)
		return c.Next()
	}, createExpenseHandler)

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock para o caso específico
			tc.mockSetup(createExpense)

			// Preparação da requisição
			bodyBytes, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "http://localhost:8080/expenses", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Execução da requisição
			resp, err := app.Test(req)
			assert.Nil(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			// Assertions comuns
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
