package controller_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/postgres"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
)

// Mock do postgres.GetExpensesPerPeriod
type MockGetExpensesPerPeriod struct {
	ExpensesPerPeriod []postgres.ExpensesPerPeriod
	Error             error
}

func (m *MockGetExpensesPerPeriod) Execute(ctx context.Context, input postgres.GetExpensesPerPeriodInput) ([]postgres.ExpensesPerPeriod, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.ExpensesPerPeriod, nil
}

func TestGetExpensesPerPeriodHandler(t *testing.T) {
	t.Parallel()

	validExpensesPerPeriodDaily := []postgres.ExpensesPerPeriod{
		{
			Date:   "2024-01-01",
			Amount: 10050,
			Count:  3,
		},
		{
			Date:   "2024-01-02",
			Amount: 15075,
			Count:  2,
		},
		{
			Date:   "2024-01-03",
			Amount: 8025,
			Count:  1,
		},
	}

	validExpensesPerPeriodMonthly := []postgres.ExpensesPerPeriod{
		{
			Date:   "2024-01",
			Amount: 35000,
			Count:  10,
		},
		{
			Date:   "2024-02",
			Amount: 42500,
			Count:  12,
		},
	}

	// Definição dos casos de teste
	testCases := []struct {
		name                  string
		groupID               int
		aggregate             string
		startDate             string
		endDate               string
		mockExpensesPerPeriod []postgres.ExpensesPerPeriod
		mockError             error
		expectedStatus        int
		expectedResponse      string
		customAssertions      func(t *testing.T, body []byte)
	}{
		{
			name:                  "should return 200 and expenses per period with daily aggregation",
			groupID:               1,
			aggregate:             "day",
			startDate:             "2024-01-01T00:00:00Z",
			endDate:               "2024-01-31T23:59:59Z",
			mockExpensesPerPeriod: validExpensesPerPeriodDaily,
			mockError:             nil,
			expectedStatus:        200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerPeriod]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 3)
				assert.Equal(t, "2024-01-01", response.Data[0].Date)
				assert.Equal(t, 10050, response.Data[0].Amount)
				assert.Equal(t, 3, response.Data[0].Count)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:                  "should return 200 and expenses per period with monthly aggregation",
			groupID:               1,
			aggregate:             "month",
			startDate:             "2024-01-01T00:00:00Z",
			endDate:               "2024-12-31T23:59:59Z",
			mockExpensesPerPeriod: validExpensesPerPeriodMonthly,
			mockError:             nil,
			expectedStatus:        200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerPeriod]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 2)
				assert.Equal(t, "2024-01", response.Data[0].Date)
				assert.Equal(t, 35000, response.Data[0].Amount)
				assert.Equal(t, 10, response.Data[0].Count)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:                  "should return 200 and empty expenses per period",
			groupID:               1,
			aggregate:             "day",
			startDate:             "2024-01-01T00:00:00Z",
			endDate:               "2024-01-31T23:59:59Z",
			mockExpensesPerPeriod: []postgres.ExpensesPerPeriod{},
			mockError:             nil,
			expectedStatus:        200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerPeriod]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 0)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:                  "should return 200 with default daily aggregation",
			groupID:               1,
			aggregate:             "",
			startDate:             "2024-01-01T00:00:00Z",
			endDate:               "2024-01-31T23:59:59Z",
			mockExpensesPerPeriod: validExpensesPerPeriodDaily,
			mockError:             nil,
			expectedStatus:        200,
			customAssertions: func(t *testing.T, body []byte) {
				var response api.Response[[]postgres.ExpensesPerPeriod]
				assert.Nil(t, json.Unmarshal(body, &response))
				assert.Equal(t, 200, response.StatusCode)
				assert.Len(t, response.Data, 3)
				assert.NotEmpty(t, response.Date)
			},
		},
		{
			name:                  "should return 400 if invalid date format",
			groupID:               1,
			aggregate:             "day",
			startDate:             "invalid-date",
			endDate:               "2024-01-31T23:59:59Z",
			mockExpensesPerPeriod: nil,
			mockError:             nil,
			expectedStatus:        400,
			expectedResponse:      `{"status_code":400,"message":"Bad Request","error":"Bad Request: internal=failed to decode: schema: error converting value for \"start_date\". Details: parsing time \"invalid-date\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid-date\" as \"2006\""}`,
		},
		{
			name:                  "should return 500 if database error",
			groupID:               1,
			aggregate:             "day",
			startDate:             "2024-01-01T00:00:00Z",
			endDate:               "2024-01-31T23:59:59Z",
			mockExpensesPerPeriod: nil,
			mockError:             errors.New("database error"),
			expectedStatus:        500,
			expectedResponse:      `{"status_code":500,"message":"Internal Server Error","error":"query.GetExpensesPerPeriod: database error"}`,
		},
	}

	// Execução dos casos de teste
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup do mock
			mockGetExpensesPerPeriod := &MockGetExpensesPerPeriod{
				ExpensesPerPeriod: tc.mockExpensesPerPeriod,
				Error:             tc.mockError,
			}

			// Setup comum para todos os testes
			app := fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
			getExpensesPerPeriodHandler := controller.NewGetExpensesPerPeriod(mockGetExpensesPerPeriod.Execute)

			app.Get("/expenses-per-period", func(c *fiber.Ctx) error {
				c.Locals("group_id", tc.groupID)
				return c.Next()
			}, getExpensesPerPeriodHandler)

			// Preparação da requisição
			url := "http://localhost:8080/expenses-per-period"
			params := []string{}

			if tc.aggregate != "" {
				params = append(params, "aggregate="+tc.aggregate)
			}
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
