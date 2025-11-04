package service

import "context"

type Predicter interface {
	ExpenseCategory(ctx context.Context, name string, amount int) (categoryID int, err error)
}
