package income

import (
	"context"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/modules/user"
	"github.com/Beigelman/nossas-despesas/internal/pkg/ddd"
)

type Repository interface {
	ddd.Repository[ID, Income]
	GetUserMonthlyIncomes(ctx context.Context, userID user.ID, date *time.Time) ([]Income, error)
}
