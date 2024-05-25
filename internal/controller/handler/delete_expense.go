package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	DeleteExpense func(ctx *fiber.Ctx) error

	DeleteExpenseRequest struct {
		GroupID     int    `json:"group_id" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Amount      int    `json:"amount" validate:"required"`
		Description string `json:"description"`
		CategoryID  int    `json:"category_id" validate:"required"`
		SplitRatio  struct {
			Payer    int `json:"payer" validate:"required"`
			Receiver int `json:"receiver" validate:"required"`
		} `json:"split_ratio" validate:"required"`
		PayerID    int `json:"payer_id" validate:"required"`
		ReceiverID int `json:"receiver_id" validate:"required"`
	}

	DeleteExpenseResponse struct {
		ID int `json:"id"`
	}
)

func NewDeleteExpense(deleteExpense usecase.DeleteExpense) DeleteExpense {
	return func(ctx *fiber.Ctx) error {
		expenseID, err := strconv.Atoi(ctx.Params("expense_id"))
		if err != nil {
			return except.BadRequestError("invalid expense id")
		}

		expense, err := deleteExpense(ctx.Context(), entity.ExpenseID{Value: expenseID})
		if err != nil {
			return fmt.Errorf("DeleteExpense: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(
			api.NewResponse(http.StatusOK, DeleteExpenseResponse{
				ID: expense.ID.Value,
			}),
		)
	}
}
