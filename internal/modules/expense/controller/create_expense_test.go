package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/controller"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	mockusecase "github.com/Beigelman/nossas-despesas/internal/tests/mocks/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Parallel()
	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
	})

	createExpense := mockusecase.NewMockCreateExpense(t)
	createExpenseHandler := controller.NewCreateExpense(createExpense.Execute)
	bodyReq := controller.CreateExpenseRequest{
		Name:        "Test Expense",
		Amount:      100,
		Description: "My first expense",
		CategoryID:  1,
		SplitType:   "equal",
		PayerID:     1,
		ReceiverID:  1,
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

	app.Post("/expenses", func(c *fiber.Ctx) error {
		c.Locals("group_id", 1)
		return c.Next()
	}, createExpenseHandler)

	t.Run("should return 201 and create a new expense", func(t *testing.T) {
		createExpense.EXPECT().Execute(mock.Anything, mock.Anything).Return(newExpense, nil).Once()
		bodyBytes, _ := json.Marshal(bodyReq)
		req := httptest.NewRequest("POST", "http://localhost:8080/expenses", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		// http.Response

		resp, err := app.Test(req)
		assert.Nil(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		var response api.Response[controller.CreateExpenseResponse]
		assert.Nil(t, json.Unmarshal(body, &response))
		assert.Equal(t, 201, response.StatusCode)
		assert.Equal(t, float32(1), response.Data.Amount)
		assert.Equal(t, 1, response.Data.PayerID)
		assert.Equal(t, 2, response.Data.ReceiverID)
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("should return 400 if request body is invalid", func(t *testing.T) {
		bodyBytes, _ := json.Marshal(map[string]any{
			"group_id":    1,
			"name":        "Test Expense",
			"amount":      100,
			"description": "My first expense",
			"category_id": 1,
			"payer_id":    1,
			"receiver_id": 2,
		})
		req := httptest.NewRequest("POST", "http://localhost:8080/expenses", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		// http.Response

		resp, err := app.Test(req)
		assert.Nil(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		// then
		assert.Equal(t, `{"status_code":400,"message":"invalid request body","error":"invalid request body: internal=validation errors: [SplitType]: '' | Needs to implement 'required'"}`, string(body))
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("should return 422 if request body is not processable", func(t *testing.T) {
		createExpense.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, except.UnprocessableEntityError()).Once()
		bodyBytes, _ := json.Marshal(bodyReq)
		req := httptest.NewRequest("POST", "http://localhost:8080/expenses", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		// http.Response

		resp, err := app.Test(req)
		assert.Nil(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		// then
		assert.Equal(t, `{"status_code":422,"message":"Unprocessable Entity","error":"Unprocessable Entity"}`, string(body))
		assert.Equal(t, 422, resp.StatusCode)
	})

	t.Run("should return 500 if it gets an unexpected error", func(t *testing.T) {
		createExpense.EXPECT().Execute(mock.Anything, mock.Anything).Return(nil, errors.New("unexpected error")).Once()
		bodyBytes, _ := json.Marshal(bodyReq)
		req := httptest.NewRequest("POST", "http://localhost:8080/expenses", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		// http.Response

		resp, err := app.Test(req)
		assert.Nil(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		// then
		assert.Equal(t, `{"status_code":500,"message":"Internal Server Error"}`, string(body))
		assert.Equal(t, 500, resp.StatusCode)
	})
}
