package expense

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/civil"

	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ScheduledExpenseID struct{ Value int }

type ScheduledExpense struct {
	ddd.Entity[ScheduledExpenseID]
	Name            string
	Amount          int
	Description     string
	GroupID         group.ID
	CategoryID      category.ID
	SplitType       SplitType
	PayerID         user.ID
	ReceiverID      user.ID
	FrequencyInDays int
	LastGeneratedAt *civil.Date
	IsActive        bool
}

type ScheduledExpenseAttributes struct {
	ID              ScheduledExpenseID
	Name            string
	Amount          int
	Description     string
	GroupID         group.ID
	CategoryID      category.ID
	SplitType       SplitType
	PayerID         user.ID
	ReceiverID      user.ID
	LastGeneratedAt *civil.Date
	FrequencyInDays int
}

func NewScheduledExpense(attr ScheduledExpenseAttributes) (*ScheduledExpense, error) {
	scheduledExpense := &ScheduledExpense{
		Entity: ddd.Entity[ScheduledExpenseID]{
			ID:        attr.ID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name:            attr.Name,
		Amount:          attr.Amount,
		Description:     attr.Description,
		GroupID:         attr.GroupID,
		CategoryID:      attr.CategoryID,
		SplitType:       attr.SplitType,
		PayerID:         attr.PayerID,
		ReceiverID:      attr.ReceiverID,
		FrequencyInDays: attr.FrequencyInDays,
		LastGeneratedAt: attr.LastGeneratedAt,
		IsActive:        true,
	}

	if err := scheduledExpense.validate(); err != nil {
		return nil, fmt.Errorf("scheduled expense validation failed: %w", err)
	}

	return scheduledExpense, nil
}

func (s *ScheduledExpense) validate() error {
	if s.Name == "" {
		return fmt.Errorf("name is required")
	}

	if s.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}

	return nil
}

func (s *ScheduledExpense) ToExpense() (*Expense, error) {
	createdAt := time.Now()

	return New(Attributes{
		Name:        s.Name,
		Amount:      s.Amount,
		Description: s.Description,
		GroupID:     s.GroupID,
		CategoryID:  s.CategoryID,
		SplitRatio:  NewEqualSplitRatio(), // This is a temporary value, it will be updated when the expense is created
		SplitType:   s.SplitType,
		PayerID:     s.PayerID,
		ReceiverID:  s.ReceiverID,
		CreatedAt:   &createdAt,
	})
}

func (se *ScheduledExpense) ShouldGenerateExpense() bool {
	if !se.IsActive {
		return false
	}

	if se.LastGeneratedAt == nil {
		return true
	}

	today := civil.DateOf(time.Now())
	nextGeneration := se.calculateNextGenerationDate()
	return today == nextGeneration || today.After(nextGeneration)
}

func (se *ScheduledExpense) calculateNextGenerationDate() civil.Date {
	lastGen := se.LastGeneratedAt
	if lastGen == nil {
		return civil.DateOf(se.CreatedAt)
	}

	return lastGen.AddDays(se.FrequencyInDays)
}

func (se *ScheduledExpense) UpdateLastGeneratedAt() {
	today := civil.DateOf(time.Now())
	se.LastGeneratedAt = &today
	se.UpdatedAt = time.Now()
	se.Version++
}

func (se *ScheduledExpense) Deactivate() {
	se.IsActive = false
	se.UpdatedAt = time.Now()
	se.Version++
}

type ScheduledExpenseRepository interface {
	ddd.Repository[ScheduledExpenseID, ScheduledExpense]
	GetActiveScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error)
	BulkStore(ctx context.Context, scheduledExpenses []ScheduledExpense) error
}
