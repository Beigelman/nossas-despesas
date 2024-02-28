package repository

import (
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type ExpenseRepository interface {
	ddd.Repository[entity.ExpenseID, entity.Expense]
}
