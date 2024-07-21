package income

import (
	"context"
	"github.com/Beigelman/nossas-despesas/internal/domain/entity"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Income]
	GetUserMonthlyIncomes(ctx context.Context, userID entity.UserID, date *time.Time) ([]Income, error)
}
