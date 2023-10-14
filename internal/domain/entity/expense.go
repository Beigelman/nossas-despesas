package entity

import (
	"fmt"
	"time"

	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type SplitRatio struct {
	Payer    int
	Receiver int
}

type ExpenseID struct{ Value int }

type Expense struct {
	ddd.Entity[ExpenseID]
	Name        string
	Amount      int // Value in cents
	Description string
	GroupID     GroupID
	CategoryID  CategoryID
	SplitRatio  SplitRatio
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
	SplitRatio  SplitRatio
	PayerID     UserID
	ReceiverID  UserID
}

func NewExpense(p ExpenseParams) (*Expense, error) {
	expense := Expense{
		Entity: ddd.Entity[ExpenseID]{
			ID:        p.ID,
			CreatedAt: time.Now(),
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

func (e *Expense) validate() error {
	if e.SplitRatio.Payer+e.SplitRatio.Receiver != 100 {
		return ErrInvalidSplitRatio
	}

	return nil
}
