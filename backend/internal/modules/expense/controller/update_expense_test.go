package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

func TestUpdateExpenseHandler(t *testing.T) {
	t.Parallel()

	validUpdateRequest := controller.UpdateExpenseRequest{
		Name:        stringPtr("Updated Expense"),
		Amount:      intPtr(200),
		Description: stringPtr("Updated description"),
		CategoryID:  intPtr(2),
		SplitType:   stringPtr("proportional"),
		PayerID:     intPtr(1),
		ReceiverID:  intPtr(2),
	}

	partialUpdateRequest := controller.UpdateExpenseRequest{
		Name:   stringPtr("Partially Updated Expense"),
		Amount: intPtr(150),
	}

	invalidSplitTypeRequest := controller.UpdateExpenseRequest{
		Name:      stringPtr("Updated Expense"),
		SplitType: stringPtr("invalid"),
	}

	updatedExpense, _ := expense.New(expense.Attributes{
		ID:          expense.ID{Value: 1},
		Name:        "Updated Expense",
		Amount:      200,
		Description: "Updated description",
		GroupID:     group.ID{Value: 1},
		CategoryID:  category.ID{Value: 2},
		SplitRatio:  expense.SplitRatio{Payer: 60, Receiver: 40},
		PayerID:     user.ID{Value: 1},
		ReceiverID:  user.ID{Value: 2},
	})

	// Definição dos casos de teste
	testCases := []struct {
		name             string
		expenseID        string
		body             any
		mockSetup        func(usecase *mocks.MockusecaseUpdateExpense)
		expectedStatus   int
		expectedResponse string
		customAssertions func(t *testing.T, body []byte)
	}{
		{
			name:      "should return 201 and update the expense",
			expenseID: "1",
			body:      validUpdateRequest,
			mockSetup: func(usecase *mocks.MockusecaseUpdateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(updatedExpense, nil).Once()
			},
			expectedStatus: 201,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.UpdateExpenseResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 201, response.StatusCode)
				assert.Equal(t, 1, response.Data.ID)
				assert.Equal(t, "Updated Expense", response.Data.Name)
				assert.Equal(t, float32(2.0), response.Data.Amount)
				assert.Equal(t, 1, response.Data.PayerID)
				assert.Equal(t, 2, response.Data.ReceiverID)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:      "should return 201 and update expense with partial data",
			expenseID: "1",
			body:      partialUpdateRequest,
			mockSetup: func(usecase *mocks.MockusecaseUpdateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(updatedExpense, nil).Once()
			},
			expectedStatus: 201,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.UpdateExpenseResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 201, response.StatusCode)
				assert.Equal(t, 1, response.Data.ID)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:             "should return 400 if expense_id is not a valid number",
			expenseID:        "invalid",
			body:             validUpdateRequest,
			mockSetup:        func(usecase *mocks.MockusecaseUpdateExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid expense id","error":"invalid expense id"}`,
		},
		{
			name:             "should return 400 if request body is invalid",
			expenseID:        "1",
			body:             invalidSplitTypeRequest,
			mockSetup:        func(usecase *mocks.MockusecaseUpdateExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid request body","error":"invalid request body: internal=validation errors: [SplitType]: 'invalid' | Needs to implement 'oneof'"}`,
		},
		{
			name:             "should return 422 if request body cannot be parsed",
			expenseID:        "1",
			body:             "invalid json",
			mockSetup:        func(usecase *mocks.MockusecaseUpdateExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   422,
			expectedResponse: `{"status_code":422,"message":"Unprocessable Entity","error":"Unprocessable Entity: internal=invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:      "should return 404 if expense is not found",
			expenseID: "999",
			body:      validUpdateRequest,
			mockSetup: func(usecase *mocks.MockusecaseUpdateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, except.NotFoundError("expense not found")).Once()
			},
			expectedStatus:   404,
			expectedResponse: `{"status_code":404,"message":"expense not found","error":"UpdateExpense: expense not found"}`,
		},
		{
			name:      "should return 422 if request is not processable",
			expenseID: "1",
			body:      validUpdateRequest,
			mockSetup: func(usecase *mocks.MockusecaseUpdateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, except.UnprocessableEntityError()).Once()
			},
			expectedStatus:   422,
			expectedResponse: `{"status_code":422,"message":"Unprocessable Entity","error":"UpdateExpense: Unprocessable Entity"}`,
		},
		{
			name:      "should return 500 if it gets an unexpected error",
			expenseID: "1",
			body:      validUpdateRequest,
			mockSetup: func(usecase *mocks.MockusecaseUpdateExpense) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, errors.New("unexpected error")).Once()
			},
			expectedStatus:   500,
			expectedResponse: `{"status_code":500,"message":"Internal Server Error","error":"UpdateExpense: unexpected error"}`,
		},
	}

	// Setup comum para todos os testes
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

	updateExpense := mocks.NewMockusecaseUpdateExpense(t)
	updateExpenseHandler := controller.NewUpdateExpense(updateExpense.Execute)
	app.Put("/expenses/:expense_id", updateExpenseHandler)

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock para o caso específico
			tc.mockSetup(updateExpense)

			// Preparação da requisição
			var req *http.Request
			if body, ok := tc.body.(string); ok {
				// Para casos de JSON inválido
				req = httptest.NewRequest("PUT", "http://localhost:8080/expenses/"+tc.expenseID, bytes.NewBufferString(body))
			} else {
				// Para casos normais
				bodyBytes, _ := json.Marshal(tc.body)
				req = httptest.NewRequest("PUT", "http://localhost:8080/expenses/"+tc.expenseID, bytes.NewBuffer(bodyBytes))
			}
			req.Header.Set("Content-Type", "application/json")

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

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
