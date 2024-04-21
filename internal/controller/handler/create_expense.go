package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	vo "github.com/Beigelman/nossas-despesas/internal/domain/valueobject"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateExpense func(ctx *fiber.Ctx) error

	CreateExpenseRequest struct {
		Name        string     `json:"name" validate:"required"`
		Amount      int        `json:"amount" validate:"required"`
		Description string     `json:"description"`
		CategoryID  int        `json:"category_id" validate:"required"`
		SplitType   string     `json:"split_type" validate:"required,oneof=equal proportional transfer"`
		PayerID     int        `json:"payer_id" validate:"required"`
		ReceiverID  int        `json:"receiver_id" validate:"required"`
		CreatedAt   *time.Time `json:"created_at"`
	}

	CreateExpenseResponse struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Amount     float32 `json:"amount"`
		PayerID    int     `json:"payer_id"`
		ReceiverID int     `json:"receiver_id"`
	}
)

func NewCreateExpense(createExpense usecase.CreateExpense) CreateExpense {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		var req CreateExpenseRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.UnprocessableEntityError("group_id not found in context")
		}

		expense, err := createExpense(ctx.Context(), usecase.CreateExpenseParams{
			GroupID:     entity.GroupID{Value: groupID},
			Name:        req.Name,
			Amount:      req.Amount,
			Description: req.Description,
			CategoryID:  entity.CategoryID{Value: req.CategoryID},
			SplitType:   vo.SplitType(req.SplitType),
			PayerID:     entity.UserID{Value: req.PayerID},
			ReceiverID:  entity.UserID{Value: req.ReceiverID},
			CreatedAt:   req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("CreateExpense: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, CreateExpenseResponse{
				ID:         expense.ID.Value,
				Name:       expense.Name,
				Amount:     float32(expense.Amount) / 100,
				PayerID:    expense.PayerID.Value,
				ReceiverID: expense.ReceiverID.Value,
			}),
		)
	}
}
