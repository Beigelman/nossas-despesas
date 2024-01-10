package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/api"
	"github.com/Beigelman/ludaapi/internal/pkg/except"
	"github.com/Beigelman/ludaapi/internal/query"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type (
	GetGroupExpenses func(ctx *fiber.Ctx) error

	GetGroupExpensesCursor struct {
		LastExpenseID int `json:"last_expense_id"`
	}

	GetGroupExpensesResponse struct {
		Expenses  []query.ExpenseDetails `json:"expenses"`
		NextToken string                 `json:"next_token"`
	}
)

func NewGetGroupExpenses(getGroupExpenses query.GetGroupExpenses) GetGroupExpenses {
	const defaultLimit = 100

	return func(ctx *fiber.Ctx) error {
		groupID, err := strconv.Atoi(ctx.Params("group_id"))
		if err != nil {
			return except.BadRequestError("invalid group id").SetInternal(err)
		}

		token, err := decodeCursor(ctx.Query("next_token", ""))
		if err != nil {
			return except.BadRequestError("invalid next jwt").SetInternal(err)
		}

		expenses, err := getGroupExpenses(ctx.Context(), query.GetGroupExpensesInput{
			GroupID:       groupID,
			LastExpenseID: token.LastExpenseID,
			Limit:         defaultLimit,
		})
		if err != nil {
			return fmt.Errorf("query.GetGroupExpenses: %w", err)
		}

		nextToken := ""
		if expenses != nil && len(expenses) == defaultLimit {
			nextToken, err = encodeCursor(&GetGroupExpensesCursor{
				LastExpenseID: expenses[len(expenses)-1].ID,
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
			LastExpenseID: 0,
		}, nil
	}

	var cur *GetGroupExpensesCursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}

	return cur, nil
}
