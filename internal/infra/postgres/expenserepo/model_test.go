package expenserepo

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpenseModel(t *testing.T) {
	// Test fields
	expense := ExpenseModel{
		ID:          1,
		Name:        "Test Expense",
		AmountCents: 1000,
		Description: "Test Description",
		GroupID:     1,
		CategoryID:  2,
		SplitRatio: SplitRatio{
			Payer:    1,
			Receiver: 2,
		},
		PayerID:    1,
		ReceiverID: 2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		Version:    1,
	}
	assert.Equal(t, 1, expense.ID)
	assert.Equal(t, "Test Expense", expense.Name)
	assert.Equal(t, 1000, expense.AmountCents)
	assert.Equal(t, "Test Description", expense.Description)
	assert.Equal(t, 1, expense.GroupID)
	assert.Equal(t, 2, expense.CategoryID)
	assert.Equal(t, SplitRatio{Payer: 1, Receiver: 2}, expense.SplitRatio)
	assert.Equal(t, 1, expense.PayerID)
	assert.Equal(t, 2, expense.ReceiverID)
	assert.True(t, expense.CreatedAt.Before(time.Now()))
	assert.True(t, expense.UpdatedAt.Before(time.Now()))
	assert.True(t, expense.DeletedAt.Valid)
	assert.Equal(t, 1, expense.Version)
}
