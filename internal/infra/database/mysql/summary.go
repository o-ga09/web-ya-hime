package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/category"
	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
)

type summaryRepository struct{}

func NewSummaryRepository() summary.ISummaryRepository {
	return &summaryRepository{}
}

func (s *summaryRepository) Save(ctx context.Context, model *summary.Summary) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection not found in context")
	}

	query := `INSERT INTO summaries (id, title, description, content, category, category_id, subcategory_id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	_, err := db.ExecContext(ctx, query, model.ID, model.Title, model.Description, model.Content, model.Category, model.CategoryID, model.SubcategoryID, model.UserID)
	if err != nil {
		return fmt.Errorf("failed to save summary: %w", err)
	}

	return nil
}

func (s *summaryRepository) List(ctx context.Context, opts summary.ListOptions) (*summary.ListResult, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	// デフォルト値の設定
	if opts.Limit <= 0 {
		opts.Limit = 20
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// WHERE句の構築
	whereClause := "WHERE s.deleted_at IS NULL"
	args := []interface{}{}
	if opts.Category != "" {
		whereClause += " AND s.category = ?"
		args = append(args, opts.Category)
	}
	if opts.CategoryID != "" {
		whereClause += " AND s.category_id = ?"
		args = append(args, opts.CategoryID)
	}
	if opts.SubcategoryID != "" {
		whereClause += " AND s.subcategory_id = ?"
		args = append(args, opts.SubcategoryID)
	}

	// 総件数を取得
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM summaries s %s`, whereClause)
	var total int
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// データ取得
	query := fmt.Sprintf(`
		SELECT 
			s.id, s.title, s.description, s.content, s.category, s.category_id, s.subcategory_id, s.user_id, s.created_at, s.updated_at,
			u.id, u.name, u.email, u.user_type, u.created_at, u.updated_at,
			c.id, c.name, c.created_at, c.updated_at,
			sc.id, sc.category_id, sc.name, sc.created_at, sc.updated_at
		FROM summaries s
		LEFT JOIN users u ON s.user_id = u.id AND u.deleted_at IS NULL
		LEFT JOIN categories c ON s.category_id = c.id
		LEFT JOIN subcategories sc ON s.subcategory_id = sc.id
		%s
		ORDER BY s.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	queryArgs := append(args, opts.Limit+1, opts.Offset) // +1で次のページの有無を判定
	rows, err := db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary list: %w", err)
	}
	defer rows.Close()

	var summaries summary.SummarySlice
	for rows.Next() {
		var s summary.Summary
		var u user.User
		var catID, catName sql.NullString
		var catCreatedAt, catUpdatedAt sql.NullTime
		var subID, subCategoryID, subName sql.NullString
		var subCreatedAt, subUpdatedAt sql.NullTime

		if err := rows.Scan(
			&s.ID, &s.Title, &s.Description, &s.Content, &s.Category, &s.CategoryID, &s.SubcategoryID, &s.UserID, &s.CreatedAt, &s.UpdatedAt,
			&u.ID, &u.Name, &u.Email, &u.UserType, &u.CreatedAt, &u.UpdatedAt,
			&catID, &catName, &catCreatedAt, &catUpdatedAt,
			&subID, &subCategoryID, &subName, &subCreatedAt, &subUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan summary: %w", err)
		}

		s.User = &user.User{
			WYHBaseModel: domain.WYHBaseModel{
				ID:        u.ID,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			},
			Name:     u.Name,
			Email:    u.Email,
			UserType: u.UserType,
		}

		if catID.Valid {
			s.CategoryObj = &category.Category{
				WYHBaseModel: domain.WYHBaseModel{
					ID:        catID.String,
					CreatedAt: catCreatedAt.Time,
					UpdatedAt: catUpdatedAt.Time,
				},
				Name: catName.String,
			}
		}

		if subID.Valid {
			s.Subcategory = &subcategory.Subcategory{
				WYHBaseModel: domain.WYHBaseModel{
					ID:        subID.String,
					CreatedAt: subCreatedAt.Time,
					UpdatedAt: subUpdatedAt.Time,
				},
				CategoryID: subCategoryID.String,
				Name:       subName.String,
			}
		}

		summaries = append(summaries, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// 次のページがあるか判定
	hasNext := len(summaries) > opts.Limit
	if hasNext {
		summaries = summaries[:opts.Limit]
	}

	return &summary.ListResult{
		Items:   summaries,
		Total:   total,
		Limit:   opts.Limit,
		Offset:  opts.Offset,
		HasNext: hasNext,
	}, nil
}

func (s *summaryRepository) Detail(ctx context.Context, model *summary.Summary) (*summary.Summary, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	query := `
		SELECT 
			s.id, s.title, s.description, s.content, s.category, s.category_id, s.subcategory_id, s.user_id, s.created_at, s.updated_at,
			u.id, u.name, u.email, u.user_type, u.created_at, u.updated_at,
			c.id, c.name, c.created_at, c.updated_at,
			sc.id, sc.category_id, sc.name, sc.created_at, sc.updated_at
		FROM summaries s
		LEFT JOIN users u ON s.user_id = u.id AND u.deleted_at IS NULL
		LEFT JOIN categories c ON s.category_id = c.id
		LEFT JOIN subcategories sc ON s.subcategory_id = sc.id
		WHERE s.id = ? AND s.deleted_at IS NULL
	`
	var result summary.Summary
	var userID, userName, userEmail, userType sql.NullString
	var userCreatedAt, userUpdatedAt sql.NullTime
	var catID, catName sql.NullString
	var catCreatedAt, catUpdatedAt sql.NullTime
	var subID, subCategoryID, subName sql.NullString
	var subCreatedAt, subUpdatedAt sql.NullTime

	err := db.QueryRowContext(ctx, query, model.ID).Scan(
		&result.ID,
		&result.Title,
		&result.Description,
		&result.Content,
		&result.Category,
		&result.CategoryID,
		&result.SubcategoryID,
		&result.UserID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&userID,
		&userName,
		&userEmail,
		&userType,
		&userCreatedAt,
		&userUpdatedAt,
		&catID,
		&catName,
		&catCreatedAt,
		&catUpdatedAt,
		&subID,
		&subCategoryID,
		&subName,
		&subCreatedAt,
		&subUpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("summary not found")
		}
		return nil, fmt.Errorf("failed to get summary detail: %w", err)
	}

	if userID.Valid {
		result.User = &user.User{
			WYHBaseModel: domain.WYHBaseModel{
				ID:        userID.String,
				CreatedAt: userCreatedAt.Time,
				UpdatedAt: userUpdatedAt.Time,
			},
			Name:     userName.String,
			Email:    userEmail.String,
			UserType: userType.String,
		}
	}

	if catID.Valid {
		result.CategoryObj = &category.Category{
			WYHBaseModel: domain.WYHBaseModel{
				ID:        catID.String,
				CreatedAt: catCreatedAt.Time,
				UpdatedAt: catUpdatedAt.Time,
			},
			Name: catName.String,
		}
	}

	if subID.Valid {
		result.Subcategory = &subcategory.Subcategory{
			WYHBaseModel: domain.WYHBaseModel{
				ID:        subID.String,
				CreatedAt: subCreatedAt.Time,
				UpdatedAt: subUpdatedAt.Time,
			},
			CategoryID: subCategoryID.String,
			Name:       subName.String,
		}
	}

	return &result, nil
}

func (s *summaryRepository) Delete(ctx context.Context, model *summary.Summary) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection not found in context")
	}

	query := `UPDATE summaries SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result, err := db.ExecContext(ctx, query, model.ID)
	if err != nil {
		return fmt.Errorf("failed to delete summary: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("summary not found")
	}

	return nil
}
