package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Mock do postgres.GetExpensesPerCategory
type MockGetExpensesPerCategory struct {
	ExpensesPerCategory []postgres.ExpensesPerCategory
	Error               error
}

func (m *MockGetExpensesPerCategory) Execute(ctx context.Context, input postgres.GetExpensesPerCategoryInput) ([]postgres.ExpensesPerCategory, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ExpensesPerCategory, nil
}

func TestGetExpensesPerCategoryHandler(t *testing.T) {
	t.Parallel()

	// Criando dados de teste usando raw JSON unmarshaling para contornar limitações de tipos não exportados
	validExpensesPerCategory := []postgres.ExpensesPerCategory{}

	// Simular dados usando JSON para contornar o tipo não exportado
	jsonData := `[
		{
			"name": "Food",
			"amount": 50075,
			"categories": [
				{"name": "Restaurant", "amount": 30000},
				{"name": "Grocery", "amount": 20075}
			]
		},
		{
			"name": "Transport",
			"amount": 25050,
			"categories": [
				{"name": "Bus", "amount": 15000},
				{"name": "Uber", "amount": 10050}
			]
		},
		{
			"name": "Entertainment",
			"amount": 15025,
			"categories": [
				{"name": "Movies", "amount": 8000},
				{"name": "Games", "amount": 7025}
			]
		}
	]`

	err := json.Unmarshal([]byte(jsonData), &validExpensesPerCategory)
	if err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	// Definição dos casos de teste
	testCases := []struct {
		name                    string
		groupID                 int
		startDate               string
		endDate                 string
		mockExpensesPerCategory []postgres.ExpensesPerCategory
		mockError               error
		expectedStatus          int
		expectedResponse        string
		customAssertions        func(t *testing.T, body []byte)
	}{
		{
			name:                    "should return 200 and expenses per category",
			groupID:                 1,
			startDate:               "2024-01-01T00:00:00Z",
			endDate:                 "2024-01-31T23:59:59Z",
			mockExpensesPerCategory: validExpensesPerCategory,
			mockError:               nil,
			expectedStatus:          200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerCategory]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 3)
				assert.NotEmpty(t, response.Date)

				// Verificar que os dados foram retornados corretamente
				assert.Contains(t, string(body), "Food")
				assert.Contains(t, string(body), "Transport")
				assert.Contains(t, string(body), "Entertainment")
			},
		},
		{
			name:                    "should return 200 and empty expenses per category",
			groupID:                 1,
			startDate:               "2024-01-01T00:00:00Z",
			endDate:                 "2024-01-31T23:59:59Z",
			mockExpensesPerCategory: []postgres.ExpensesPerCategory{},
			mockError:               nil,
			expectedStatus:          200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerCategory]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 0)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:                    "should return 200 with no date filters",
			groupID:                 1,
			startDate:               "",
			endDate:                 "",
			mockExpensesPerCategory: validExpensesPerCategory,
			mockError:               nil,
			expectedStatus:          200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerCategory]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 3)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:                    "should return 422 if invalid date format",
			groupID:                 1,
			startDate:               "invalid-date",
			endDate:                 "2024-01-31T23:59:59Z",
			mockExpensesPerCategory: nil,
			mockError:               nil,
			expectedStatus:          422,
			expectedResponse:        `{"status_code":422,"message":"Unprocessable Entity","error":"Unprocessable Entity: internal=failed to decode: schema: error converting value for \"start_date\". Details: parsing time \"invalid-date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid-date\" as \"2006\""}`,
		},
		{
			name:                    "should return 500 if database error",
			groupID:                 1,
			startDate:               "2024-01-01T00:00:00Z",
			endDate:                 "2024-01-31T23:59:59Z",
			mockExpensesPerCategory: nil,
			mockError:               errors.New("database error"),
			expectedStatus:          500,
			expectedResponse:        `{"status_code":500,"message":"Internal Server Error","error":"query.getExpensesPerCategory: database error"}`,
		},
	}

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock
			mockGetExpensesPerCategory := &MockGetExpensesPerCategory{
				ExpensesPerCategory: tc.mockExpensesPerCategory,
				Error:               tc.mockError,
			}

			// Setup comum para todos os testes
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			getExpensesPerCategoryHandler := controller.NewGetExpensesPerCategory(mockGetExpensesPerCategory.Execute)

			app.Get("/expenses-per-category", func(c *fiber.Ctx) error {
				c.Locals("group_id", tc.groupID)
				return c.Next()
			}, getExpensesPerCategoryHandler)

			// Preparação da requisição
			url := "http://localhost:8080/expenses-per-category"
			params := []string{}

			if tc.startDate != "" {
				params = append(params, "start_date="+tc.startDate)
			}
			if tc.endDate != "" {
				params = append(params, "end_date="+tc.endDate)
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
