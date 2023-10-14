package repository

import (
	"github.com/Beigelman/ludaapi/internal/domain/entity"
	"github.com/Beigelman/ludaapi/internal/pkg/ddd"
)

type ExpenseRepository interface {
	ddd.Repository[entity.ExpenseID, entity.Expense]
}
