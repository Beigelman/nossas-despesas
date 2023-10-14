package expenserepo

import (
	"database/sql"
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"testing"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
	"github.com/stretchr/testify/assert"
)

func TestToEntity(t *testing.T) {
	// Test with valid input
	model := ExpenseModel{
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
	expense := toEntity(model)
	assert.NotNil(t, expense)
	assert.Equal(t, entity.ExpenseID{Value: 1}, expense.ID)
	assert.Equal(t, "Test Expense", expense.Name)
	assert.Equal(t, 1000, expense.Amount)
	assert.Equal(t, "Test Description", expense.Description)
	assert.Equal(t, entity.GroupID{Value: 1}, expense.GroupID)
	assert.Equal(t, entity.CategoryID{Value: 2}, expense.CategoryID)
	assert.Equal(t, entity.SplitRatio{Payer: 1, Receiver: 2}, expense.SplitRatio)
	assert.Equal(t, entity.UserID{Value: 1}, expense.PayerID)
	assert.Equal(t, entity.UserID{Value: 2}, expense.ReceiverID)
	assert.True(t, expense.CreatedAt.Before(time.Now()))
	assert.True(t, expense.UpdatedAt.Before(time.Now()))
	assert.True(t, expense.DeletedAt != nil)
	assert.Equal(t, 1, expense.Version)

	// Test with invalid input
	model = ExpenseModel{
		ID:          2,
		Name:        "Test Expense 2",
		AmountCents: 2000,
		Description: "Test Description 2",
		GroupID:     3,
		CategoryID:  4,
		SplitRatio: SplitRatio{
			Payer:    3,
			Receiver: 4,
		},
		PayerID:    3,
		ReceiverID: 4,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt:  sql.NullTime{Time: time.Time{}, Valid: false},
		Version:    2,
	}
	expense = toEntity(model)
	assert.NotNil(t, expense)
	assert.Equal(t, entity.ExpenseID{Value: 2}, expense.ID)
	assert.Equal(t, "Test Expense 2", expense.Name)
	assert.Equal(t, 2000, expense.Amount)
	assert.Equal(t, "Test Description 2", expense.Description)
	assert.Equal(t, entity.GroupID{Value: 3}, expense.GroupID)
	assert.Equal(t, entity.CategoryID{Value: 4}, expense.CategoryID)
	assert.Equal(t, entity.SplitRatio{Payer: 3, Receiver: 4}, expense.SplitRatio)
	assert.Equal(t, entity.UserID{Value: 3}, expense.PayerID)
	assert.Equal(t, entity.UserID{Value: 4}, expense.ReceiverID)
	assert.True(t, expense.CreatedAt.Before(time.Now()))
	assert.True(t, expense.UpdatedAt.Before(time.Now()))
	assert.True(t, expense.DeletedAt == nil)
	assert.Equal(t, 2, expense.Version)
}

func TestToModel(t *testing.T) {
	// Test with valid input
	deletedAt := time.Now()
	expense := &entity.Expense{
		Entity: ddd.Entity[entity.ExpenseID]{
			ID:        entity.ExpenseID{Value: 1},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: &deletedAt,
			Version:   1,
		},
		Name:        "Test Expense",
		Amount:      1000,
		Description: "Test Description",
		GroupID:     entity.GroupID{Value: 1},
		CategoryID:  entity.CategoryID{Value: 2},
		SplitRatio: entity.SplitRatio{
			Payer:    1,
			Receiver: 2,
		},
		PayerID:    entity.UserID{Value: 1},
		ReceiverID: entity.UserID{Value: 2},
	}
	model := toModel(expense)
	assert.NotNil(t, model)
	assert.Equal(t, 1, model.ID)
	assert.Equal(t, "Test Expense", model.Name)
	assert.Equal(t, 1000, model.AmountCents)
	assert.Equal(t, "Test Description", model.Description)
	assert.Equal(t, 1, model.GroupID)
	assert.Equal(t, 2, model.CategoryID)
	assert.Equal(t, SplitRatio{Payer: 1, Receiver: 2}, model.SplitRatio)
	assert.Equal(t, 1, model.PayerID)
	assert.Equal(t, 2, model.ReceiverID)
	assert.True(t, model.CreatedAt.Before(time.Now()))
	assert.True(t, model.UpdatedAt.Before(time.Now()))
	assert.True(t, model.DeletedAt.Valid)
	assert.Equal(t, 1, model.Version)

	// Test with invalid input
	expense = &entity.Expense{
		Entity: ddd.Entity[entity.ExpenseID]{
			ID:        entity.ExpenseID{Value: 2},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
			Version:   2,
		},
		Name:        "Test Expense 2",
		Amount:      2000,
		Description: "Test Description 2",
		GroupID:     entity.GroupID{Value: 3},
		CategoryID:  entity.CategoryID{Value: 4},
		SplitRatio: entity.SplitRatio{
			Payer:    3,
			Receiver: 4,
		},
		PayerID:    entity.UserID{Value: 3},
		ReceiverID: entity.UserID{Value: 4},
	}
	model = toModel(expense)
	assert.NotNil(t, model)
	assert.Equal(t, 2, model.ID)
	assert.Equal(t, "Test Expense 2", model.Name)
	assert.Equal(t, 2000, model.AmountCents)
	assert.Equal(t, "Test Description 2", model.Description)
	assert.Equal(t, 3, model.GroupID)
	assert.Equal(t, 4, model.CategoryID)
	assert.Equal(t, SplitRatio{Payer: 3, Receiver: 4}, model.SplitRatio)
	assert.Equal(t, 3, model.PayerID)
	assert.Equal(t, 4, model.ReceiverID)
	assert.True(t, model.CreatedAt.Before(time.Now()))
	assert.True(t, model.UpdatedAt.Before(time.Now()))
	assert.False(t, model.DeletedAt.Valid)
	assert.Equal(t, 2, model.Version)
}
