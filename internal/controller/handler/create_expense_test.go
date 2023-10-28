package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Beigelman/ludaapi/internal/controller/handler"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	mockusecase "github.com/Beigelman/ludaapi/internal/tests/mocks/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http/httptest"
	"testing"
)

func TestCreateExpenseHandler(t *testing.T) {
	t.Parallel()
	app := fiber.New(fiber.Config{
		ErrorHandler: api.ErrorHandler,
	})

	createExpense := mockusecase.NewMockCreateExpense(t)
	createExpenseHandler := handler.NewCreateExpense(createExpense.Execute)
	bodyReq := handler.CreateExpenseRequest{
		GroupID:     1,
		Name:        "Test Expense",
		Amount:      100,
		Description: "My first expense",
		CategoryID:  1,
		SplitRatio: struct {
			Payer    int `json:"payer" validate:"required"`
			Receiver int `json:"receiver" validate:"required"`
		}{
			Payer:    50,
			Receiver: 50,
		},
		PayerID:    1,
		ReceiverID: 1,
	}

	newExpense, _ := entity.NewExpense(entity.ExpenseParams{
		ID:          entity.ExpenseID{Value: 1},
		Name:        "Test Expense",
		Amount:      100,
		Description: "My first expense",
		GroupID:     entity.GroupID{Value: 1},
		CategoryID:  entity.CategoryID{Value: 1},
		SplitRatio:  entity.SplitRatio{Payer: 50, Receiver: 50},
		PayerID:     entity.UserID{Value: 1},
		ReceiverID:  entity.UserID{Value: 2},
	})

	app.Post("/expenses", createExpenseHandler)

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
		var response api.Response[handler.CreateExpenseResponse]
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
		assert.Equal(t, `{"status_code":400,"message":"error=invalid request body, internal=validation errors: [Payer]: '0' | Needs to implement 'required' and [Receiver]: '0' | Needs to implement 'required'"}`, string(body))
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
		assert.Equal(t, `{"status_code":422,"message":"error=Unprocessable Entity"}`, string(body))
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
