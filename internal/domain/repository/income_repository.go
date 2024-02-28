package repository

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type IncomeRepository interface {
	ddd.Repository[entity.IncomeID, entity.Income]
	GetUserMonthlyIncomes(ctx context.Context, userID entity.UserID, date *time.Time) ([]entity.Income, error)
}
