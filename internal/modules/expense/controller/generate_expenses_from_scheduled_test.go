package controller_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/shared/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateExpensesFromScheduledHandler(t *testing.T) {
	t.Parallel()

	// Definição dos casos de teste
	testCases := []struct {
		name             string
		mockSetup        func(usecase *mocks.MockusecaseGenerateExpensesFromScheduledUseCase)
		expectedStatus   int
		expectedResponse string
		customAssertions func(t *testing.T, body []byte)
	}{
		{
			name: "should return 200 and count of expenses created",
			mockSetup: func(usecase *mocks.MockusecaseGenerateExpensesFromScheduledUseCase) {
				usecase.EXPECT().Execute(mock.Anything).Return(5, nil).Once()
			},
			expectedStatus: 200,
			customAssertions: func(t *testing.T, body []byte) {
				var response fiber.Map
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, float64(5), response["expenses_created"])
			},
		},
		{
			name: "should return 200 with zero expenses created",
			mockSetup: func(usecase *mocks.MockusecaseGenerateExpensesFromScheduledUseCase) {
				usecase.EXPECT().Execute(mock.Anything).Return(0, nil).Once()
			},
			expectedStatus: 200,
			customAssertions: func(t *testing.T, body []byte) {
				var response fiber.Map
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, float64(0), response["expenses_created"])
			},
		},
		{
			name: "should return 500 if usecase fails",
			mockSetup: func(usecase *mocks.MockusecaseGenerateExpensesFromScheduledUseCase) {
				usecase.EXPECT().Execute(mock.Anything).Return(0, errors.New("database error")).Once()
			},
			expectedStatus:   500,
			expectedResponse: `{"status_code":500,"message":"Internal Server Error","error":"GenerateExpensesFromScheduled: database error"}`,
		},
	}

	// Setup comum para todos os testes
	app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})

	generateExpensesFromScheduled := mocks.NewMockusecaseGenerateExpensesFromScheduledUseCase(t)
	generateExpensesFromScheduledHandler := controller.NewGenerateExpensesFromScheduled(generateExpensesFromScheduled.Execute)
	app.Post("/generate-expenses", generateExpensesFromScheduledHandler)

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock para o caso específico
			tc.mockSetup(generateExpensesFromScheduled)

			// Preparação da requisição
			req := httptest.NewRequest("POST", "http://localhost:8080/generate-expenses", nil)

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
