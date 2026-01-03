package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/o-ga09/web-ya-hime/internal/domain/category"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

type categoryRepository struct{}

func NewCategoryRepository() category.ICategoryRepository {
	return &categoryRepository{}
}

func (r *categoryRepository) Save(ctx context.Context, model *category.Category) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection is not set in context")
	}

	if model.ID == "" {
		model.ID = uuid.GenerateID()
	}

	query := `
		INSERT INTO categories (id, name, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP(6), CURRENT_TIMESTAMP(6))
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			updated_at = CURRENT_TIMESTAMP(6)
	`

	_, err := db.ExecContext(ctx, query, model.ID, model.Name)
	if err != nil {
		return fmt.Errorf("failed to save category: %w", err)
	}

	return nil
}

func (r *categoryRepository) List(ctx context.Context) (category.CategorySlice, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection is not set in context")
	}

	query := `
		SELECT id, name, created_at, updated_at
		FROM categories
		ORDER BY created_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	var categories category.CategorySlice
	for rows.Next() {
		var c category.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) Detail(ctx context.Context, model *category.Category) (*category.Category, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection is not set in context")
	}

	query := `
		SELECT id, name, created_at, updated_at
		FROM categories
		WHERE id = ?
	`

	var c category.Category
	err := db.QueryRowContext(ctx, query, model.ID).Scan(
		&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get category detail: %w", err)
	}

	return &c, nil
}

func (r *categoryRepository) Delete(ctx context.Context, model *category.Category) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection is not set in context")
	}

	query := `DELETE FROM categories WHERE id = ?`

	result, err := db.ExecContext(ctx, query, model.ID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
