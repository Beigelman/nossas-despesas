package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"net/http"

	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type CreateExpense func(ctx *fiber.Ctx) error

type CreateExpenseRequest struct {
	GroupID     int    `json:"group_id"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id"`
	SplitRatio  struct {
		Payer    int `json:"payer"`
		Receiver int `json:"receiver"`
	} `json:"split_ratio"`
	PayerID    int `json:"payer_id"`
	ReceiverID int `json:"receiver_id"`
}

type CreateExpenseResponse struct {
	ID         int     `json:"id"`
	Amount     float32 `json:"name"`
	PayerID    int     `json:"payer_id"`
	ReceiverID int     `json:"receiver_id"`
}

func NewCreateExpenseHandler(createExpense usecase.CreateExpense) CreateExpense {
	return func(ctx *fiber.Ctx) error {
		var req CreateExpenseRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		expense, err := createExpense(ctx.Context(), usecase.CreateExpenseParams{
			GroupID:     entity.GroupID{Value: req.GroupID},
			Name:        req.Name,
			Amount:      req.Amount,
			Description: req.Description,
			CategoryID:  entity.CategoryID{Value: req.CategoryID},
			SplitRatio: entity.SplitRatio{
				Payer:    req.SplitRatio.Payer,
				Receiver: req.SplitRatio.Receiver,
			},
			PayerID:    entity.UserID{Value: req.PayerID},
			ReceiverID: entity.UserID{Value: req.ReceiverID},
		})
		if err != nil {
			return fmt.Errorf("CreateExpense: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateExpenseResponse{
				ID:         expense.ID.Value,
				Amount:     float32(expense.Amount) / 100,
				PayerID:    expense.PayerID.Value,
				ReceiverID: expense.ReceiverID.Value,
			}),
		)
	}
}
