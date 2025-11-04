package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
)

func TestPredictExpenseCategoryHandler(t *testing.T) {
	t.Parallel()

	validBodyReq := controller.PredictExpenseCategoryRequest{
		Name:   "Farmácia",
		Amount: 5000,
	}

	invalidBodyReq := map[string]any{
		"name": "Farmácia",
		// amount is missing
	}

	// Definição dos casos de teste
	testCases := []struct {
		name             string
		body             any
		mockSetup        func(usecase *mocks.MockusecasePredictExpenseCategory)
		expectedStatus   int
		expectedResponse string
		customAssertions func(t *testing.T, body []byte)
	}{
		{
			name: "should return 200 and predict category successfully",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecasePredictExpenseCategory) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(1, nil).Once()
			},
			expectedStatus: 200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[controller.PredictExpenseCategoryResponse]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Equal(t, 1, response.Data.CategoryID)
				assert.Equal(t, "Farmácia", response.Data.Name)
				assert.Equal(t, 5000, response.Data.Amount)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:             "should return 400 if request body is invalid JSON",
			body:             "invalid json",
			mockSetup:        func(usecase *mocks.MockusecasePredictExpenseCategory) {},
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid request body","error":"invalid request body: internal=invalid character 'i' looking for beginning of value"}`,
		},
		{
			name:             "should return 400 if required fields are missing",
			body:             invalidBodyReq,
			mockSetup:        func(usecase *mocks.MockusecasePredictExpenseCategory) {},
			expectedStatus:   400,
			expectedResponse: `{"status_code":400,"message":"invalid request body","error":"invalid request body: internal=validation errors: [Amount]: '0' | Needs to implement 'required'"}`,
		},
		{
			name: "should return 500 if usecase returns error",
			body: validBodyReq,
			mockSetup: func(usecase *mocks.MockusecasePredictExpenseCategory) {
				usecase.EXPECT().Execute(mock.Anything, mock.Anything).Return(0, errors.New("unexpected error")).Once()
			},
			expectedStatus:   500,
			expectedResponse: `{"status_code":500,"message":"Internal Server Error","error":"predictExpenseCategory: unexpected error"}`,
		},
	}

	// Setup comum para todos os testes
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

	predictExpenseCategory := mocks.NewMockusecasePredictExpenseCategory(t)
	predictExpenseCategoryHandler := controller.NewPredictExpenseCategory(predictExpenseCategory.Execute)
	app.Post("/predict", predictExpenseCategoryHandler)

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock para o caso específico
			tc.mockSetup(predictExpenseCategory)

			// Preparação da requisição
			var bodyBytes []byte
			var err error
			if str, ok := tc.body.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tc.body)
				assert.Nil(t, err)
			}
			req := httptest.NewRequest("POST", "http://localhost:8080/predict", bytes.NewBuffer(bodyBytes))
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

