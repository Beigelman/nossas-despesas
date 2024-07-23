package expense

import (
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/modules/category"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ID struct{ Value int }

type Expense struct {
	ddd.Entity[ID]
	Name         string
	Amount       int // Value in cents
	RefundAmount *int
	Description  string
	GroupID      group.ID
	CategoryID   category.ID
	SplitRatio   SplitRatio
	SplitType    SplitType
	PayerID      user.ID
	ReceiverID   user.ID
}

type Attributes struct {
	ID          ID
	Name        string
	Amount      int
	Description string
	GroupID     group.ID
	CategoryID  category.ID
	SplitRatio  SplitRatio
	SplitType   SplitType
	PayerID     user.ID
	ReceiverID  user.ID
	CreatedAt   *time.Time
}

type UpdateAttributes struct {
	Name         *string
	Amount       *int
	RefundAmount *int
	Description  *string
	CategoryID   *category.ID
	SplitRatio   *SplitRatio
	SplitType    *SplitType
	PayerID      *user.ID
	ReceiverID   *user.ID
	CreatedAt    *time.Time
}

func New(attr Attributes) (*Expense, error) {
	createdAt := time.Now()
	if attr.CreatedAt != nil {
		createdAt = *attr.CreatedAt
	}
	expense := Expense{
		Entity: ddd.Entity[ID]{
			ID:        attr.ID,
			CreatedAt: createdAt,
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Name:        attr.Name,
		Amount:      attr.Amount,
		Description: attr.Description,
		GroupID:     attr.GroupID,
		CategoryID:  attr.CategoryID,
		SplitRatio:  attr.SplitRatio,
		SplitType:   attr.SplitType,
		PayerID:     attr.PayerID,
		ReceiverID:  attr.ReceiverID,
	}

	if err := expense.validate(); err != nil {
		return nil, fmt.Errorf("expense.Validate: %w", err)
	}

	return &expense, nil
}

func (e *Expense) Update(p UpdateAttributes) error {
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
	if p.SplitType != nil {
		e.SplitType = *p.SplitType
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

	if e.SplitType == SplitTypes.Equal && e.SplitRatio.Payer != 50 {
		return ErrInvalidSplitRatio
	}

	if e.RefundAmount != nil && *e.RefundAmount > e.Amount {
		return ErrInvalidRefundAmount
	}

	return nil
}
