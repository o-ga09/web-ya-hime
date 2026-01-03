package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/o-ga09/web-ya-hime/internal/domain/category"
	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

type subcategoryRepository struct{}

func NewSubcategoryRepository() subcategory.ISubcategoryRepository {
	return &subcategoryRepository{}
}

func (r *subcategoryRepository) Save(ctx context.Context, model *subcategory.Subcategory) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection is not set in context")
	}

	if model.ID == "" {
		model.ID = uuid.GenerateID()
	}

	query := `
		INSERT INTO subcategories (id, category_id, name, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP(6), CURRENT_TIMESTAMP(6))
		ON DUPLICATE KEY UPDATE
			name = VALUES(name),
			updated_at = CURRENT_TIMESTAMP(6)
	`

	_, err := db.ExecContext(ctx, query, model.ID, model.CategoryID, model.Name)
	if err != nil {
		return fmt.Errorf("failed to save subcategory: %w", err)
	}

	return nil
}

func (r *subcategoryRepository) List(ctx context.Context, categoryID string) (subcategory.SubcategorySlice, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection is not set in context")
	}

	var query string
	var args []interface{}

	if categoryID != "" {
		query = `
			SELECT s.id, s.category_id, s.name, s.created_at, s.updated_at,
				   c.id, c.name, c.created_at, c.updated_at
			FROM subcategories s
			INNER JOIN categories c ON s.category_id = c.id
			WHERE s.category_id = ?
			ORDER BY s.created_at DESC
		`
		args = append(args, categoryID)
	} else {
		query = `
			SELECT s.id, s.category_id, s.name, s.created_at, s.updated_at,
				   c.id, c.name, c.created_at, c.updated_at
			FROM subcategories s
			INNER JOIN categories c ON s.category_id = c.id
			ORDER BY s.created_at DESC
		`
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list subcategories: %w", err)
	}
	defer rows.Close()

	var subcategories subcategory.SubcategorySlice
	for rows.Next() {
		var s subcategory.Subcategory
		s.Category = &category.Category{}
		if err := rows.Scan(
			&s.ID, &s.CategoryID, &s.Name, &s.CreatedAt, &s.UpdatedAt,
			&s.Category.ID, &s.Category.Name, &s.Category.CreatedAt, &s.Category.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subcategory: %w", err)
		}
		subcategories = append(subcategories, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subcategories: %w", err)
	}

	return subcategories, nil
}

func (r *subcategoryRepository) Detail(ctx context.Context, model *subcategory.Subcategory) (*subcategory.Subcategory, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection is not set in context")
	}

	query := `
		SELECT s.id, s.category_id, s.name, s.created_at, s.updated_at,
			   c.id, c.name, c.created_at, c.updated_at
		FROM subcategories s
		INNER JOIN categories c ON s.category_id = c.id
		WHERE s.id = ?
	`

	var s subcategory.Subcategory
	s.Category = &category.Category{}
	err := db.QueryRowContext(ctx, query, model.ID).Scan(
		&s.ID, &s.CategoryID, &s.Name, &s.CreatedAt, &s.UpdatedAt,
		&s.Category.ID, &s.Category.Name, &s.Category.CreatedAt, &s.Category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("subcategory not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get subcategory detail: %w", err)
	}

	return &s, nil
}

func (r *subcategoryRepository) Delete(ctx context.Context, model *subcategory.Subcategory) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection is not set in context")
	}

	query := `DELETE FROM subcategories WHERE id = ?`

	result, err := db.ExecContext(ctx, query, model.ID)
	if err != nil {
		return fmt.Errorf("failed to delete subcategory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subcategory not found")
	}

	return nil
}
