package handler

import (
	"fmt"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	vo "github.com/Beigelman/ludaapi/internal/domain/valueobject"
	"github.com/Beigelman/ludaapi/internal/pkg/validator"
	"net/http"
	"strconv"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"github.com/gofiber/fiber/v2"
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

		expense, err := updateExpense(ctx.Context(), usecase.UpdateExpenseParams{
			ID:           entity.ExpenseID{Value: expenseID},
			Name:         req.Name,
			Amount:       req.Amount,
			RefundAmount: req.RefundAmount,
			Description:  req.Description,
			CategoryID: func() *entity.CategoryID {
				if req.CategoryID != nil {
					return &entity.CategoryID{Value: *req.CategoryID}
				}
				return nil
			}(),
			SplitType: func() *vo.SplitType {
				if req.SplitType != nil {
					splitType := vo.SplitType(*req.SplitType)
					return &splitType
				}
				return nil
			}(),
			PayerID: func() *entity.UserID {
				if req.PayerID != nil {
					return &entity.UserID{Value: *req.PayerID}
				}
				return nil
			}(),
			ReceiverID: func() *entity.UserID {
				if req.ReceiverID != nil {
					return &entity.UserID{Value: *req.ReceiverID}
				}
				return nil
			}(),
			CreatedAt: req.CreatedAt,
		})
		if err != nil {
			return fmt.Errorf("UpdateExpense: %w", err)
		}

		return ctx.Status(http.StatusCreated).JSON(
			api.NewResponse[UpdateExpenseResponse](http.StatusCreated, UpdateExpenseResponse{
				ID:         expense.ID.Value,
				Name:       expense.Name,
				Amount:     float32(expense.Amount) / 100,
				PayerID:    expense.PayerID.Value,
				ReceiverID: expense.ReceiverID.Value,
			}),
		)
	}
}
