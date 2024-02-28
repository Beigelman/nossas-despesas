package entity

import (
	"fmt"
	vo "github.com/Beigelman/nossas-despesas/internal/domain/valueobject"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ExpenseID struct{ Value int }

type Expense struct {
	ddd.Entity[ExpenseID]
	Name         string
	Amount       int // Value in cents
	RefundAmount *int
	Description  string
	GroupID      GroupID
	CategoryID   CategoryID
	SplitRatio   vo.SplitRatio
	PayerID      UserID
	ReceiverID   UserID
}

type ExpenseParams struct {
	ID          ExpenseID
	Name        string
	Amount      int
	Description string
	GroupID     GroupID
	CategoryID  CategoryID
	SplitRatio  vo.SplitRatio
	PayerID     UserID
	ReceiverID  UserID
	CreatedAt   *time.Time
}

type ExpenseUpdateParams struct {
	Name         *string
	Amount       *int
	RefundAmount *int
	Description  *string
	CategoryID   *CategoryID
	SplitRatio   *vo.SplitRatio
	PayerID      *UserID
	ReceiverID   *UserID
	CreatedAt    *time.Time
}

func NewExpense(p ExpenseParams) (*Expense, error) {
	createdAt := time.Now()
	if p.CreatedAt != nil {
		createdAt = *p.CreatedAt
	}
	expense := Expense{
		Entity: ddd.Entity[ExpenseID]{
			ID:        p.ID,
			CreatedAt: createdAt,
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name:        p.Name,
		Amount:      p.Amount,
		Description: p.Description,
		GroupID:     p.GroupID,
		CategoryID:  p.CategoryID,
		SplitRatio:  p.SplitRatio,
		PayerID:     p.PayerID,
		ReceiverID:  p.ReceiverID,
	}

	if err := expense.validate(); err != nil {
		return nil, fmt.Errorf("expense.Validate: %w", err)
	}

	return &expense, nil
}

func (e *Expense) Update(p ExpenseUpdateParams) error {
	if p.Name != nil {
		e.Name = *p.Name
	}
	if p.Amount != nil {
		e.Amount = *p.Amount
	}
	if p.RefundAmount != nil {
		e.RefundAmount = p.RefundAmount
	}
	if p.Description != nil {
		e.Description = *p.Description
	}
	if p.CategoryID != nil {
		e.CategoryID = *p.CategoryID
	}
	if p.SplitRatio != nil {
		e.SplitRatio = *p.SplitRatio
	}
	if p.PayerID != nil {
		e.PayerID = *p.PayerID
	}
	if p.ReceiverID != nil {
		e.ReceiverID = *p.ReceiverID
	}
	if p.CreatedAt != nil {
		e.CreatedAt = *p.CreatedAt
	}
	e.UpdatedAt = time.Now()
	e.Version++

	if err := e.validate(); err != nil {
		return fmt.Errorf("expense.Validate: %w", err)
	}

	return nil
}

func (e *Expense) Delete() {
	now := time.Now()
	e.DeletedAt = &now
	e.UpdatedAt = now
	e.Version++
}

func (e *Expense) validate() error {
	if e.SplitRatio.Payer+e.SplitRatio.Receiver != 100 || e.SplitRatio.Payer > 99 {
		return ErrInvalidSplitRatio
	}

	if e.RefundAmount != nil && *e.RefundAmount > e.Amount {
		return ErrInvalidRedfundAmount
	}

	return nil
}
