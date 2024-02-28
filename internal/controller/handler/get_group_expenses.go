package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/Beigelman/nossas-despesas/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type (
	GetGroupExpenses func(ctx *fiber.Ctx) error

	GetGroupExpensesCursor struct {
		LastExpenseDate time.Time `json:"last_expense_date"`
	}

	GetGroupExpensesResponse struct {
		Expenses  []query.ExpenseDetails `json:"expenses"`
		NextToken string                 `json:"next_token"`
	}
)

func NewGetGroupExpenses(getGroupExpenses query.GetGroupExpenses) GetGroupExpenses {
	const defaultLimit = 25

	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		token, err := decodeCursor(ctx.Query("next_token", ""))
		if err != nil {
			return except.BadRequestError("invalid next jwt").SetInternal(err)
		}

		expenses, err := getGroupExpenses(ctx.Context(), query.GetGroupExpensesInput{
			GroupID:         groupID,
			LastExpenseDate: token.LastExpenseDate,
			Limit:           defaultLimit,
		})
		if err != nil {
			return fmt.Errorf("query.GetGroupExpenses: %w", err)
		}

		nextToken := ""
		if expenses != nil && len(expenses) == defaultLimit {
			nextToken, err = encodeCursor(&GetGroupExpensesCursor{
				LastExpenseDate: expenses[len(expenses)-1].CreatedAt,
			})
			if err != nil {
				return fmt.Errorf("encodeCursor: %w", err)
			}
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse[GetGroupExpensesResponse](http.StatusOK, GetGroupExpensesResponse{
			Expenses:  expenses,
			NextToken: nextToken,
		}))
	}
}

func encodeCursor(cursor *GetGroupExpensesCursor) (string, error) {
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(serializedCursor), nil
}

func decodeCursor(cursor string) (*GetGroupExpensesCursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	if string(decodedCursor) == "" {
		return &GetGroupExpensesCursor{
			LastExpenseDate: time.Now().AddDate(0, 2, 0),
		}, nil
	}

	var cur *GetGroupExpensesCursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}

	return cur, nil
}
