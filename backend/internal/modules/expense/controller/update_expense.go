package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
)

type (
	UpdateExpense func(ctx *fiber.Ctx) error

	UpdateExpenseRequest struct {
		Name         *string    `json:"name"`
		Amount       *int       `json:"amount"`
		RefundAmount *int       `json:"refund_amount"`
		Description  *string    `json:"description"`
		CategoryID   *int       `json:"category_id"`
		SplitType    *string    `json:"split_type" validate:"omitempty,oneof=equal proportional transfer"`
		PayerID      *int       `json:"payer_id"`
		ReceiverID   *int       `json:"receiver_id"`
		CreatedAt    *time.Time `json:"created_at"`
	}

	UpdateExpenseResponse struct {
		ID         int     `json:"id"`
		Name       string  `json:"name"`
		Amount     float32 `json:"amount"`
		PayerID    int     `json:"payer_id"`
		ReceiverID int     `json:"receiver_id"`
	}
)

func NewUpdateExpense(updateExpense usecase.UpdateExpense) UpdateExpense {
	valid := validator.New()
	return func(ctx *fiber.Ctx) error {
		expenseID, err := strconv.Atoi(ctx.Params("expense_id"))
		if err != nil {
			return except.BadRequestError("invalid expense id")
		}

		var req UpdateExpenseRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		expns, err := updateExpense(ctx.Context(), usecase.UpdateExpenseParams{
			ID:           expense.ID{Value: expenseID},
			Name:         req.Name,
			Amount:       req.Amount,
			RefundAmount: req.RefundAmount,
			Description:  req.Description,
			CategoryID: func() *category.ID {
				if req.CategoryID != nil {
					return &category.ID{Value: *req.CategoryID}
				}
				return nil
			}(),
			SplitType: func() *expense.SplitType {
				if req.SplitType != nil {
					splitType := expense.SplitType(*req.SplitType)
					return &splitType
				}
				return nil
			}(),
			PayerID: func() *user.ID {
				if req.PayerID != nil {
					return &user.ID{Value: *req.PayerID}
				}
				return nil
			}(),
			ReceiverID: func() *user.ID {
				if req.ReceiverID != nil {
					return &user.ID{Value: *req.ReceiverID}
				}
				return nil
			}(),
			CreatedAt: req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("UpdateExpense: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse(http.StatusCreated, UpdateExpenseResponse{
				ID:         expns.ID.Value,
				Name:       expns.Name,
				Amount:     float32(expns.Amount) / 100,
				PayerID:    expns.PayerID.Value,
				ReceiverID: expns.ReceiverID.Value,
			}),
		)
	}
}
