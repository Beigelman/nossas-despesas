package controller

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/gofiber/fiber/v2"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	vo "github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
)

type CreateScheduledExpenseRequest struct {
	Name            string     `json:"name" validate:"required"`
	Amount          int        `json:"amount" validate:"required"`
	Description     string     `json:"description" validate:"required"`
	CategoryID      int        `json:"category_id" validate:"required"`
	SplitType       string     `json:"split_type" validate:"required"`
	PayerID         int        `json:"payer_id" validate:"required"`
	ReceiverID      int        `json:"receiver_id" validate:"required"`
	FrequencyInDays int        `json:"frequency_in_days" validate:"required"`
	LastGeneratedAt *time.Time `json:"last_generated_at"`
}

type CreateScheduledExpense func(ctx *fiber.Ctx) error

func NewCreateScheduledExpense(createScheduledExpense usecase.CreateScheduledExpense) CreateScheduledExpense {
	valid := validator.New()
	return func(c *fiber.Ctx) error {
		var req CreateScheduledExpenseRequest
		if err := c.BodyParser(&req); err != nil {
			return except.UnprocessableEntityError().SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		groupID, ok := c.Locals("group_id").(int)
		if !ok {
			return except.UnprocessableEntityError("group_id not found in context")
		}

		var lastGeneratedAt *civil.Date
		if req.LastGeneratedAt != nil {
			date := civil.DateOf(*req.LastGeneratedAt)
			lastGeneratedAt = &date
		}

		err := createScheduledExpense(c.Context(), usecase.CreateScheduledExpenseInput{
			Name:            req.Name,
			Amount:          req.Amount,
			Description:     req.Description,
			GroupID:         group.ID{Value: groupID},
			CategoryID:      category.ID{Value: req.CategoryID},
			SplitType:       vo.SplitType(req.SplitType),
			PayerID:         user.ID{Value: req.PayerID},
			ReceiverID:      user.ID{Value: req.ReceiverID},
			FrequencyInDays: req.FrequencyInDays,
			LastGeneratedAt: lastGeneratedAt,
		})
		if err != nil {
			return fmt.Errorf("CreateScheduledExpense: %w", err)
		}

		return c.SendStatus(http.StatusCreated)
	}
}
