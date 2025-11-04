package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	expensesPerCategoryInfo struct {
		CategoryGroup string `db:"category_group_name"`
		Category      string `db:"category_name" `
		Amount        int    `db:"amount"`
	}

	expensesPerInnerCategory struct {
		Category string `json:"name"`
		Amount   int    `json:"amount"`
	}

	ExpensesPerCategory struct {
		CategoryGroup string                     `json:"name"`
		Amount        int                        `json:"amount"`
		Categories    []expensesPerInnerCategory `json:"categories"`
	}

	GetExpensesPerCategoryInput struct {
		GroupID   int       `json:"group_id"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	GetExpensesPerCategory func(ctx context.Context, params GetExpensesPerCategoryInput) ([]ExpensesPerCategory, error)
)

func NewGetExpensesPerCategory(db *db.Client) GetExpensesPerCategory {
	dbClient := db.Conn()
	return func(ctx context.Context, params GetExpensesPerCategoryInput) ([]ExpensesPerCategory, error) {
		var info []expensesPerCategoryInfo
		if err := dbClient.SelectContext(ctx, &info, `
			SELECT 
				cat.name AS category_name, 
				cg.name AS category_group_name, 
				SUM(ex.amount_cents) AS amount 
			FROM expenses_latest ex
			INNER JOIN categories cat ON ex.category_id = cat.id
			INNER JOIN category_groups cg ON cg.id = cat.category_group_id
			WHERE ex.group_id = $1
			AND ex.created_at >= $2
			AND ex.created_at <= $3
			AND ex.deleted_at IS NULL
			AND cg.deleted_at IS NULL
			AND cat.deleted_at IS NULL
			GROUP BY 1, 2 
			ORDER BY 2 DESC;
		`, params.GroupID, params.StartDate, params.EndDate); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		var categoriesGroups = make(map[string]ExpensesPerCategory)
		for _, e := range info {
			if _, ok := categoriesGroups[e.CategoryGroup]; !ok {
				categoriesGroups[e.CategoryGroup] = ExpensesPerCategory{
					CategoryGroup: e.CategoryGroup,
					Amount:        e.Amount,
					Categories:    []expensesPerInnerCategory{{Category: e.Category, Amount: e.Amount}},
				}
			} else {
				categoriesGroups[e.CategoryGroup] = ExpensesPerCategory{
					CategoryGroup: e.CategoryGroup,
					Amount:        categoriesGroups[e.CategoryGroup].Amount + e.Amount,
					Categories: append(categoriesGroups[e.CategoryGroup].Categories, expensesPerInnerCategory{
						Category: e.Category,
						Amount:   e.Amount,
					}),
				}
			}
		}

		var expensesPerCategory []ExpensesPerCategory
		for _, cg := range categoriesGroups {
			expensesPerCategory = append(expensesPerCategory, cg)
		}

		return expensesPerCategory, nil
	}
}
