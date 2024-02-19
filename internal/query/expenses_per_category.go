package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/ludaapi/internal/pkg/db"
	"time"
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

func NewGetExpensesPerCategory(db db.Database) GetExpensesPerCategory {
	dbClient := db.Client()
	return func(ctx context.Context, params GetExpensesPerCategoryInput) ([]ExpensesPerCategory, error) {
		var info []expensesPerCategoryInfo
		if err := dbClient.SelectContext(ctx, &info, `
			with base as (
			    select
			        distinct on (ex.id) ex.id as id,
			        ex.amount_cents amount,
			        ex.category_id  as category_id,
			        ex.deleted_at as deleted_at
			    from expenses ex
			    where ex.group_id = $1
			    and ex.created_at >= $2
			    and ex.created_at <= $3
			    order by ex.id desc, ex.version desc
			)
			select 
				cat.name as category_name, 
				cg.name as category_group_name, 
				sum(amount) as amount 
			from base b
			inner join categories cat on b.category_id = cat.id
			inner join category_groups cg on cg.id = cat.category_group_id
			where b.deleted_at is null
			and cg.name != 'Balanço'
			and cat.deleted_at is null
			group by 1, 2 
			order by 2 desc;
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
