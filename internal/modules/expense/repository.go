package expense

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/modules/group"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Expense]
	GetByGroupDate(ctx context.Context, groupId group.ID, date time.Time) ([]Expense, error)
	BulkStore(ctx context.Context, expenses []Expense) error
}
