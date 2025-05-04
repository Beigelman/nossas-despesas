package expense

import (
	"context"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/group"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Expense]
	GetByGroupDate(ctx context.Context, groupId group.ID, date time.Time) ([]Expense, error)
	BulkStore(ctx context.Context, expenses []Expense) error
}

type ScheduledExpenseRepository interface {
	ddd.Repository[ScheduledExpenseID, ScheduledExpense]
	GetActiveScheduledExpenses(ctx context.Context) ([]ScheduledExpense, error)
	BulkStore(ctx context.Context, scheduledExpenses []ScheduledExpense) error
}
