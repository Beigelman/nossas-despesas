package controller

import (
	"fmt"
	"net/http"

	"github.com/Beigelman/nossas-despesas/internal/modules/expense/usecase"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type (
	PredictExpenseCategoryRequest struct {
		Name   string `json:"name" validate:"required"`
		Amount int    `json:"amount" validate:"required"`
	}

	PredictExpenseCategoryResponse struct {
		CategoryID int    `json:"category_id"`
		Name       string `json:"name"`
		Amount     int    `json:"amount"`
	}
)

type PredictExpenseCategory func(ctx *fiber.Ctx) error

func NewPredictExpenseCategory(predictExpenseCategory usecase.PredictExpenseCategory) PredictExpenseCategory {
	valid := validator.New()

	return func(ctx *fiber.Ctx) error {
		var req PredictExpenseCategoryRequest
		if err := ctx.BodyParser(&req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		if err := valid.Validate(req); err != nil {
			return except.BadRequestError("invalid request body").SetInternal(err)
		}

		categoryID, err := predictExpenseCategory(ctx.Context(), usecase.PredictExpenseCategoryInput{
			Name:   req.Name,
			Amount: req.Amount,
		})
		if err != nil {
			return fmt.Errorf("predictExpenseCategory: %w", err)
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, PredictExpenseCategoryResponse{
			CategoryID: categoryID,
			Name:       req.Name,
			Amount:     req.Amount,
		}))
	}
}
