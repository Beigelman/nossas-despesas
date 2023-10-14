package handler_test

import (
	"context"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/usecase"
	"testing"
)

func TestCreateExpenseHandler(t *testing.T) {

	createExpense := func(ctx context.Context, p usecase.CreateExpenseParams) (*entity.Expense, error) {
		return nil, nil
	}

}
