package postgres

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

func NewGetCategories(db *db.Client) GetCategories {
	dbClient := db.Conn()
	return func(ctx context.Context) ([]CategoryGroup, error) {
		var categories []Category
		if err := dbClient.SelectContext(ctx, &categories, `
			SELECT
    			id,
				name,
				icon,
				category_group_id
			FROM categories c
			WHERE c.deleted_at IS NULL
		`); err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("db.SelectContext: %w", err)
		}

		var categoryGroups []CategoryGroup
		if err := dbClient.SelectContext(ctx, &categoryGroups, `
			SELECT
    			id,
				name,
				icon
			FROM category_groups c
			WHERE c.deleted_at IS NULL
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
