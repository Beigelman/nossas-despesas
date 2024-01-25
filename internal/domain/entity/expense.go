package entity

import (
	"fmt"
	vo "github.com/Beigelman/ludaapi/internal/domain/valueobject"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type ExpenseID struct{ Value int }

type Expense struct {
	ddd.Entity[ExpenseID]
	Name        string
	Amount      int // Value in cents
	Description string
	GroupID     GroupID
	CategoryID  CategoryID
	SplitRatio  vo.SplitRatio
	PayerID     UserID
	ReceiverID  UserID
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
	Name        *string
	Amount      *int
	Description *string
	CategoryID  *CategoryID
	SplitRatio  *vo.SplitRatio
	PayerID     *UserID
	ReceiverID  *UserID
	CreatedAt   *time.Time
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
	e.Name = func() string {
		if p.Name != nil {
			return *p.Name
		}
		return e.Name
	}()
	e.Amount = func() int {
		if p.Amount != nil {
			return *p.Amount
		}
		return e.Amount
	}()
	e.Description = func() string {
		if p.Description != nil {
			return *p.Description
		}
		return e.Description
	}()
	e.CategoryID = func() CategoryID {
		if p.CategoryID != nil {
			return *p.CategoryID
		}
		return e.CategoryID
	}()
	e.SplitRatio = func() vo.SplitRatio {
		if p.SplitRatio != nil {
			return *p.SplitRatio
		}
		return e.SplitRatio
	}()
	e.PayerID = func() UserID {
		if p.PayerID != nil {
			return *p.PayerID
		}
		return e.PayerID
	}()
	e.ReceiverID = func() UserID {
		if p.ReceiverID != nil {
			return *p.ReceiverID
		}
		return e.ReceiverID
	}()
	e.CreatedAt = func() time.Time {
		if p.CreatedAt != nil {
			return *p.CreatedAt
		}
		return e.CreatedAt
	}()
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

	return nil
}
