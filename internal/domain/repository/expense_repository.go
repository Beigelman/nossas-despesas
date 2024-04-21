package repository

import (
	"context"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ExpenseRepository interface {
	ddd.Repository[entity.ExpenseID, entity.Expense]
	GetByGroupDate(ctx context.Context, groupId entity.GroupID, date time.Time) ([]entity.Expense, error)
	BulkStore(ctx context.Context, expenses []entity.Expense) error
}
