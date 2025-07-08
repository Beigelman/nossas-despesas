package postgres

import (
	"database/sql"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/expense"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
	"github.com/stretchr/testify/assert"
)

func TestToEntity_CompleteModel(t *testing.T) {
	deletedAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	refundAmount := int64(1000)

	model := ExpenseModel{
		ID:                1,
		Name:              "Test Expense",
		AmountCents:       5000,
		RefundAmountCents: sql.NullInt64{Int64: refundAmount, Valid: true},
		Description:       "Test Description",
		GroupID:           1,
		CategoryID:        1,
		SplitRatio:        SplitRatio{Payer: 70, Receiver: 30},
		SplitType:         "custom",
		PayerID:           1,
		ReceiverID:        2,
		CreatedAt:         time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		DeletedAt:         sql.NullTime{Time: deletedAt, Valid: true},
		Version:           1,
	}

	entity := ToEntity(model)

	assert.Equal(t, expense.ID{Value: 1}, entity.ID)
	assert.Equal(t, "Test Expense", entity.Name)
	assert.Equal(t, 5000, entity.Amount)
	assert.NotNil(t, entity.RefundAmount)
	assert.Equal(t, 1000, *entity.RefundAmount)
	assert.Equal(t, "Test Description", entity.Description)
	assert.Equal(t, group.ID{Value: 1}, entity.GroupID)
	assert.Equal(t, category.ID{Value: 1}, entity.CategoryID)
	assert.Equal(t, 70, entity.SplitRatio.Payer)
	assert.Equal(t, 30, entity.SplitRatio.Receiver)
	assert.Equal(t, expense.SplitType("custom"), entity.SplitType)
	assert.Equal(t, user.ID{Value: 1}, entity.PayerID)
	assert.Equal(t, user.ID{Value: 2}, entity.ReceiverID)
	assert.Equal(t, model.CreatedAt, entity.CreatedAt)
	assert.Equal(t, model.UpdatedAt, entity.UpdatedAt)
	assert.NotNil(t, entity.DeletedAt)
	assert.Equal(t, deletedAt, *entity.DeletedAt)
	assert.Equal(t, 1, entity.Version)
}

func TestToEntity_WithoutOptionalFields(t *testing.T) {
	model := ExpenseModel{
		ID:                2,
		Name:              "Simple Expense",
		AmountCents:       3000,
		RefundAmountCents: sql.NullInt64{Valid: false},
		Description:       "Simple Description",
		GroupID:           1,
		CategoryID:        1,
		SplitRatio:        SplitRatio{Payer: 50, Receiver: 50},
		SplitType:         "equal",
		PayerID:           1,
		ReceiverID:        2,
		CreatedAt:         time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC),
		DeletedAt:         sql.NullTime{Valid: false},
		Version:           0,
	}

	entity := ToEntity(model)

	assert.Equal(t, expense.ID{Value: 2}, entity.ID)
	assert.Equal(t, "Simple Expense", entity.Name)
	assert.Equal(t, 3000, entity.Amount)
	assert.Nil(t, entity.RefundAmount)
	assert.Equal(t, "Simple Description", entity.Description)
	assert.Equal(t, expense.SplitType("equal"), entity.SplitType)
	assert.Nil(t, entity.DeletedAt)
	assert.Equal(t, 0, entity.Version)
}

func TestToModel_CompleteEntity(t *testing.T) {
	deletedAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	refundAmount := 1000

	entity := &expense.Expense{
		Entity: ddd.Entity[expense.ID]{
			ID:        expense.ID{Value: 1},
			CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			DeletedAt: &deletedAt,
			Version:   1,
		},
		Name:         "Test Expense",
		Amount:       5000,
		RefundAmount: &refundAmount,
		Description:  "Test Description",
		GroupID:      group.ID{Value: 1},
		CategoryID:   category.ID{Value: 1},
		SplitRatio:   expense.SplitRatio{Payer: 70, Receiver: 30},
		SplitType:    expense.SplitType("custom"),
		PayerID:      user.ID{Value: 1},
		ReceiverID:   user.ID{Value: 2},
	}

	model := ToModel(entity)

	assert.Equal(t, 1, model.ID)
	assert.Equal(t, "Test Expense", model.Name)
	assert.Equal(t, 5000, model.AmountCents)
	assert.True(t, model.RefundAmountCents.Valid)
	assert.Equal(t, int64(1000), model.RefundAmountCents.Int64)
	assert.Equal(t, "Test Description", model.Description)
	assert.Equal(t, 1, model.GroupID)
	assert.Equal(t, 1, model.CategoryID)
	assert.Equal(t, 70, model.SplitRatio.Payer)
	assert.Equal(t, 30, model.SplitRatio.Receiver)
	assert.Equal(t, "custom", model.SplitType)
	assert.Equal(t, 1, model.PayerID)
	assert.Equal(t, 2, model.ReceiverID)
	assert.Equal(t, entity.CreatedAt, model.CreatedAt)
	assert.Equal(t, entity.UpdatedAt, model.UpdatedAt)
	assert.True(t, model.DeletedAt.Valid)
	assert.Equal(t, deletedAt, model.DeletedAt.Time)
	assert.Equal(t, 1, model.Version)
}

func TestToModel_WithoutOptionalFields(t *testing.T) {
	entity := &expense.Expense{
		Entity: ddd.Entity[expense.ID]{
			ID:        expense.ID{Value: 2},
			CreatedAt: time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC),
			DeletedAt: nil,
			Version:   0,
		},
		Name:         "Simple Expense",
		Amount:       3000,
		RefundAmount: nil,
		Description:  "Simple Description",
		GroupID:      group.ID{Value: 1},
		CategoryID:   category.ID{Value: 1},
		SplitRatio:   expense.SplitRatio{Payer: 50, Receiver: 50},
		SplitType:    expense.SplitType("equal"),
		PayerID:      user.ID{Value: 1},
		ReceiverID:   user.ID{Value: 2},
	}

	model := ToModel(entity)

	assert.Equal(t, 2, model.ID)
	assert.Equal(t, "Simple Expense", model.Name)
	assert.Equal(t, 3000, model.AmountCents)
	assert.False(t, model.RefundAmountCents.Valid)
	assert.Equal(t, "Simple Description", model.Description)
	assert.Equal(t, "equal", model.SplitType)
	assert.False(t, model.DeletedAt.Valid)
	assert.Equal(t, 0, model.Version)
}

func TestToScheduledExpenseModel_CompleteEntity(t *testing.T) {
	lastGeneratedAt := civil.Date{Year: 2024, Month: 1, Day: 1}

	entity := expense.ScheduledExpense{
		Entity: ddd.Entity[expense.ScheduledExpenseID]{
			ID:        expense.ScheduledExpenseID{Value: 1},
			CreatedAt: time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
			Version:   1,
		},
		Name:            "Scheduled Test",
		Amount:          8000,
		Description:     "Scheduled Description",
		GroupID:         group.ID{Value: 1},
		CategoryID:      category.ID{Value: 1},
		SplitType:       expense.SplitType("equal"),
		PayerID:         user.ID{Value: 1},
		ReceiverID:      user.ID{Value: 2},
		FrequencyInDays: 30,
		LastGeneratedAt: &lastGeneratedAt,
		IsActive:        true,
	}

	model := ToScheduledExpenseModel(entity)

	assert.Equal(t, 1, model.ID)
	assert.Equal(t, "Scheduled Test", model.Name)
	assert.Equal(t, 8000, model.AmountCents)
	assert.Equal(t, "Scheduled Description", model.Description)
	assert.Equal(t, 1, model.GroupID)
	assert.Equal(t, 1, model.CategoryID)
	assert.Equal(t, "equal", model.SplitType)
	assert.Equal(t, 1, model.PayerID)
	assert.Equal(t, 2, model.ReceiverID)
	assert.Equal(t, 30, model.FrequencyInDays)
	assert.True(t, model.LastGeneratedAt.Valid)
	assert.Equal(t, lastGeneratedAt, model.LastGeneratedAt.V)
	assert.True(t, model.IsActive)
	assert.Equal(t, 1, model.Version)
}

func TestToScheduledExpenseEntity_CompleteModel(t *testing.T) {
	lastGeneratedAt := civil.Date{Year: 2024, Month: 1, Day: 1}

	model := ScheduledExpenseModel{
		ID:              1,
		Name:            "Scheduled Test",
		AmountCents:     8000,
		Description:     "Scheduled Description",
		GroupID:         1,
		CategoryID:      1,
		SplitType:       "equal",
		PayerID:         1,
		ReceiverID:      2,
		FrequencyInDays: 30,
		LastGeneratedAt: sql.Null[civil.Date]{V: lastGeneratedAt, Valid: true},
		IsActive:        true,
		CreatedAt:       time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC),
		Version:         1,
	}

	entity := ToScheduledExpenseEntity(model)

	assert.Equal(t, expense.ScheduledExpenseID{Value: 1}, entity.ID)
	assert.Equal(t, "Scheduled Test", entity.Name)
	assert.Equal(t, 8000, entity.Amount)
	assert.Equal(t, "Scheduled Description", entity.Description)
	assert.Equal(t, group.ID{Value: 1}, entity.GroupID)
	assert.Equal(t, category.ID{Value: 1}, entity.CategoryID)
	assert.Equal(t, expense.SplitType("equal"), entity.SplitType)
	assert.Equal(t, user.ID{Value: 1}, entity.PayerID)
	assert.Equal(t, user.ID{Value: 2}, entity.ReceiverID)
	assert.Equal(t, 30, entity.FrequencyInDays)
	assert.NotNil(t, entity.LastGeneratedAt)
	assert.Equal(t, lastGeneratedAt, *entity.LastGeneratedAt)
	assert.True(t, entity.IsActive)
	assert.Equal(t, 1, entity.Version)
}

func TestToScheduledExpenseEntity_WithoutLastGeneratedAt(t *testing.T) {
	model := ScheduledExpenseModel{
		ID:              2,
		Name:            "Scheduled Test Basic",
		AmountCents:     6000,
		Description:     "Basic Scheduled Description",
		GroupID:         1,
		CategoryID:      1,
		SplitType:       "equal",
		PayerID:         1,
		ReceiverID:      2,
		FrequencyInDays: 7,
		LastGeneratedAt: sql.Null[civil.Date]{Valid: false},
		IsActive:        true,
		CreatedAt:       time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC),
		UpdatedAt:       time.Date(2024, 1, 2, 11, 0, 0, 0, time.UTC),
		Version:         0,
	}

	entity := ToScheduledExpenseEntity(model)

	assert.Equal(t, expense.ScheduledExpenseID{Value: 2}, entity.ID)
	assert.Equal(t, "Scheduled Test Basic", entity.Name)
	assert.Equal(t, 6000, entity.Amount)
	assert.Equal(t, 7, entity.FrequencyInDays)
	assert.Nil(t, entity.LastGeneratedAt)
	assert.True(t, entity.IsActive)
	assert.Equal(t, 0, entity.Version)
}
