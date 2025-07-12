package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateScheduledExpenseHandler(t *testing.T) {
	t.Parallel()

	validBodyReq := controller.CreateScheduledExpenseRequest{
		Name:            "Test Scheduled Expense",
		Amount:          200,
		Description:     "My first scheduled expense",
		CategoryID:      1,
		SplitType:       "equal",
		PayerID:         1,
		ReceiverID:      2,
		FrequencyInDays: 30,
		LastGeneratedAt: nil,
	}

	validBodyReqWithDate := controller.CreateScheduledExpenseRequest{
		Name:            "Test Scheduled Expense",
		Amount:          200,
		Description:     "My first scheduled expense",
		CategoryID:      1,
		SplitType:       "equal",
		PayerID:         1,
		ReceiverID:      2,
		FrequencyInDays: 30,
		LastGeneratedAt: func() *time.Time {
			t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			return &t
		}(),
	}

	invalidBodyReq := map[string]any{
		"name":        "Test Scheduled Expense",
		"amount":      200,
		"description": "My first scheduled expense",
		"category_id": 1,
		"payer_id":    1,
		"receiver_id": 2,
		// FrequencyInDays is missing
	}

	// Definição dos casos de teste
	testCases := []struct {
		name             string
		body             any
		mockSetup        func(usecase *mocks.MockusecaseCreateScheduledExpense)
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "should return 201 and create a new scheduled expense",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecaseCreateScheduledExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedStatus:   201,
			expectedResponse: "Created",
		},
		{
			name: "should return 201 and create a new scheduled expense with last_generated_at",
			body: validBodyReqWithDate,
			mockSetup: func(usecase *mocks.MockusecaseCreateScheduledExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedStatus:   201,
			expectedResponse: "Created",
		},
		{
			name:             "should return 400 if request body is invalid",
			body:             invalidBodyReq,
			mockSetup:        func(usecase *mocks.MockusecaseCreateScheduledExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid request body","error":"invalid request body: internal=validation errors: [SplitType]: '' | Needs to implement 'required' and [FrequencyInDays]: '0' | Needs to implement 'required'"}`,
		},
		{
			name:             "should return 422 if request body cannot be parsed",
			body:             "invalid json",
			mockSetup:        func(usecase *mocks.MockusecaseCreateScheduledExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   422,
			expectedResponse: `{"status_code":422,"message":"Unprocessable Entity","error":"Unprocessable Entity: internal=invalid character 'i' looking for beginning of value"}`,
		},
		{
			name: "should return 422 if request body is not processable",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecaseCreateScheduledExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(except.UnprocessableEntityError()).Once()
			},
			expectedStatus:   422,
			expectedResponse: `{"status_code":422,"message":"Unprocessable Entity","error":"CreateScheduledExpense: Unprocessable Entity"}`,
		},
		{
			name: "should return 500 if it gets an unexpected error",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecaseCreateScheduledExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(errors.New("unexpected error")).Once()
			},
			expectedStatus:   500,
			expectedResponse: `{"status_code":500,"message":"Internal Server Error","error":"CreateScheduledExpense: unexpected error"}`,
		},
	}

	// Setup comum para todos os testes
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

	createScheduledExpense := mocks.NewMockusecaseCreateScheduledExpense(t)
	createScheduledExpenseHandler := controller.NewCreateScheduledExpense(createScheduledExpense.Execute)
	app.Post("/scheduled-expenses", func(c *fiber.Ctx) error {
		c.Locals("group_id", 1)
		return c.Next()
	}, createScheduledExpenseHandler)

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock para o caso específico
			tc.mockSetup(createScheduledExpense)

			// Preparação da requisição
			var req *http.Request
			if body, ok := tc.body.(string); ok {
				// Para casos de JSON inválido
				req = httptest.NewRequest("POST", "http://localhost:8080/scheduled-expenses", bytes.NewBufferString(body))
			} else {
				// Para casos normais
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest("POST", "http://localhost:8080/scheduled-expenses", bytes.NewBuffer(bodyBytes))
			}
			req.Header.Set("Content-Type", "application/json")

			// Execução da requisição
			resp, err := app.Test(req)
			assert.Nil(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)

			// Assertions comuns
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			// Verificação da resposta
			if tc.expectedResponse != "" {
				assert.Equal(t, tc.expectedResponse, string(body))
			}
			assert.NoError(t, resp.Body.Close())
		})
	}
}
