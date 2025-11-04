package controller_test

import (
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

func TestDeleteExpenseHandler(t *testing.T) {
	t.Parallel()

	deletedExpense, _ := expense.New(expense.Attributes{
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
		expenseID        string
		mockSetup        func(usecase *mocks.MockusecaseDeleteExpense)
		expectedStatus   int
		expectedResponse string
		customAssertions func(t *testing.T, body []byte)
	}{
		{
			name:      "should return 200 and delete the expense",
			expenseID: "1",
			mockSetup: func(usecase *mocks.MockusecaseDeleteExpense) {
				usecase.EXPECT().Execute(mock.Anything, expense.ID{Value: 1}).Return(deletedExpense, nil).Once()
			},
			expectedStatus: 200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.DeleteExpenseResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Equal(t, 1, response.Data.ID)
				assert.NotEmpty(t, response.Date) // Verifica se a data existe
			},
		},
		{
			name:             "should return 400 if expense_id is not a valid number",
			expenseID:        "invalid",
			mockSetup:        func(usecase *mocks.MockusecaseDeleteExpense) {}, // Não precisa de mock para este caso
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid expense id","error":"invalid expense id"}`,
		},
		{
			name:      "should return 404 if expense is not found",
			expenseID: "999",
			mockSetup: func(usecase *mocks.MockusecaseDeleteExpense) {
				usecase.EXPECT().Execute(mock.Anything, expense.ID{Value: 999}).Return(nil, except.NotFoundError("expense not found")).Once()
			},
			expectedStatus:   404,
			expectedResponse: `{"status_code":404,"message":"expense not found","error":"DeleteExpense: expense not found"}`,
		},
		{
			name:      "should return 422 if request is not processable",
			expenseID: "1",
			mockSetup: func(usecase *mocks.MockusecaseDeleteExpense) {
				usecase.EXPECT().Execute(mock.Anything, expense.ID{Value: 1}).Return(nil, except.UnprocessableEntityError()).Once()
			},
			expectedStatus:   422,
			expectedResponse: `{"status_code":422,"message":"Unprocessable Entity","error":"DeleteExpense: Unprocessable Entity"}`,
		},
		{
			name:      "should return 500 if it gets an unexpected error",
			expenseID: "1",
			mockSetup: func(usecase *mocks.MockusecaseDeleteExpense) {
				usecase.EXPECT().Execute(mock.Anything, expense.ID{Value: 1}).Return(nil, errors.New("unexpected error")).Once()
			},
			expectedStatus:   500,
			expectedResponse: `{"status_code":500,"message":"Internal Server Error","error":"DeleteExpense: unexpected error"}`,
		},
	}

	// Setup comum para todos os testes
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

	deleteExpense := mocks.NewMockusecaseDeleteExpense(t)
	deleteExpenseHandler := controller.NewDeleteExpense(deleteExpense.Execute)
	app.Delete("/expenses/:expense_id", deleteExpenseHandler)

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock para o caso específico
			tc.mockSetup(deleteExpense)

			// Preparação da requisição
			req := httptest.NewRequest("DELETE", "http://localhost:8080/expenses/"+tc.expenseID, nil)

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
