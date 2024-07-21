package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	query2 "github.com/Beigelman/nossas-despesas/internal/modules/expense/query"
	"net/http"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/api"
	"github.com/Beigelman/nossas-despesas/internal/pkg/except"
	"github.com/gofiber/fiber/v2"
)

type (
	GetExpenses func(ctx *fiber.Ctx) error

	GetExpensesCursor struct {
		LastExpenseID   int       `json:"last_expense_id"`
		LastExpenseDate time.Time `json:"last_expense_date"`
	}

	GetExpensesResponse struct {
		Expenses  []query2.ExpenseDetails `json:"expenses"`
		NextToken string                  `json:"next_token"`
	}
)

func NewGetExpenses(getGroupExpenses query2.GetExpenses) GetExpenses {
	const defaultLimit = 25

	return func(ctx *fiber.Ctx) error {
		groupID, ok := ctx.Locals("group_id").(int)
		if !ok {
			return except.BadRequestError("invalid group id")
		}

		token, err := decodeCursor(ctx.Query("next_token", ""))
		if err != nil {
			return except.BadRequestError("invalid next token").SetInternal(err)
		}

		search := ctx.Query("search")

		expenses, err := getGroupExpenses(ctx.Context(), query2.GetExpensesInput{
			GroupID:         groupID,
			LastExpenseDate: token.LastExpenseDate,
			LastExpenseID:   token.LastExpenseID,
			Limit:           defaultLimit,
			Search:          search,
		})
		if err != nil {
			return fmt.Errorf("query.GetExpenses: %w", err)
		}

		nextToken := ""
		if len(expenses) == defaultLimit {
			lastExpense := expenses[len(expenses)-1]
			nextToken, err = encodeCursor(&GetExpensesCursor{
				LastExpenseDate: lastExpense.CreatedAt,
				LastExpenseID:   lastExpense.ID,
			})
			if err != nil {
				return fmt.Errorf("encodeCursor: %w", err)
			}
		}

		return ctx.Status(http.StatusOK).JSON(api.NewResponse(http.StatusOK, GetExpensesResponse{
			Expenses:  expenses,
			NextToken: nextToken,
		}))
	}
}

func encodeCursor(cursor *GetExpensesCursor) (string, error) {
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(serializedCursor), nil
}

func decodeCursor(cursor string) (*GetExpensesCursor, error) {
	decodedCursor, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	if string(decodedCursor) == "" {
		return &GetExpensesCursor{
			LastExpenseDate: time.Now().AddDate(0, 2, 0),
		}, nil
	}

	var cur *GetExpensesCursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}

	return cur, nil
}
