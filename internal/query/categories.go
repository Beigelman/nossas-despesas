package query

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

type (
	Category struct {
		ID              int    `db:"id" json:"id"`
		Name            string `db:"name" json:"name"`
		Icon            string `db:"icon" json:"icon"`
		CategoryGroupID int    `db:"category_group_id" json:"-"`
	}

	CategoryGroup struct {
		ID         int        `db:"id" json:"id"`
		Name       string     `db:"name" json:"name"`
		Icon       string     `db:"icon" json:"icon"`
		Categories []Category `json:"categories"`
	}

	GetCategories func(ctx context.Context) ([]CategoryGroup, error)
)

func NewGetCategories(db db.Database) GetCategories {
	dbClient := db.Client()
	return func(ctx context.Context) ([]CategoryGroup, error) {
		var categories []Category
		if err := dbClient.SelectContext(ctx, &categories, `
			select
    			id,
				name,
				icon,
				category_group_id
			from categories c
			where c.deleted_at is null
		`); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		var categoryGroups []CategoryGroup
		if err := dbClient.SelectContext(ctx, &categoryGroups, `
			select
    			id,
				name,
				icon
			from category_groups c
			where c.deleted_at is null
		`); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		for i, categoryGroup := range categoryGroups {
			for _, category := range categories {
				if category.CategoryGroupID == categoryGroup.ID {
					categoryGroups[i].Categories = append(categoryGroups[i].Categories, category)
				}
			}
		}

		return categoryGroups, nil
	}
}
