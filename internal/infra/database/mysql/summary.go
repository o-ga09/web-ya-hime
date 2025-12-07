package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/o-ga09/web-ya-hime/internal/domain"
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

	query := `INSERT INTO summaries (id, title, description, content, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())`
	_, err := db.ExecContext(ctx, query, model.ID, model.Title, model.Description, model.Content, model.UserID)
	if err != nil {
		return fmt.Errorf("failed to save summary: %w", err)
	}

	return nil
}

func (s *summaryRepository) List(ctx context.Context) (summary.SummarySlice, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	query := `
		SELECT 
			s.id, s.title, s.description, s.content, s.user_id, s.created_at, s.updated_at,
			u.id, u.name, u.email, u.user_type, u.created_at, u.updated_at
		FROM summaries s
		LEFT JOIN users u ON s.user_id = u.id
		ORDER BY s.created_at DESC
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get summary list: %w", err)
	}
	defer rows.Close()

	var summaries summary.SummarySlice
	for rows.Next() {
		var s summary.Summary
		var userID, userName, userEmail, userType sql.NullString
		var userCreatedAt, userUpdatedAt sql.NullTime

		if err := rows.Scan(
			&s.ID, &s.Title, &s.Description, &s.Content, &s.UserID, &s.CreatedAt, &s.UpdatedAt,
			&userID, &userName, &userEmail, &userType, &userCreatedAt, &userUpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan summary: %w", err)
		}

		if userID.Valid {
			s.User = &user.User{
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

		summaries = append(summaries, &s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return summaries, nil
}

func (s *summaryRepository) Detail(ctx context.Context, model *summary.Summary) (*summary.Summary, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	query := `
		SELECT 
			s.id, s.title, s.description, s.content, s.user_id, s.created_at, s.updated_at,
			u.id, u.name, u.email, u.user_type, u.created_at, u.updated_at
		FROM summaries s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.id = ?
	`
	var result summary.Summary
	var userID, userName, userEmail, userType sql.NullString
	var userCreatedAt, userUpdatedAt sql.NullTime

	err := db.QueryRowContext(ctx, query, model.ID).Scan(
		&result.ID,
		&result.Title,
		&result.Description,
		&result.Content,
		&result.UserID,
		&result.CreatedAt,
		&result.UpdatedAt,
		&userID,
		&userName,
		&userEmail,
		&userType,
		&userCreatedAt,
		&userUpdatedAt,
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

	return &result, nil
}

func (s *summaryRepository) Delete(ctx context.Context, model *summary.Summary) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection not found in context")
	}

	query := `DELETE FROM summaries WHERE id = ?`
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
