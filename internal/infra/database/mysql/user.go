package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
)

type User struct {
	ID       string
	Version  int
	Name     string
	Email    string
	UserType string
}

func NewUserRepository() user.IUserRepository {
	return &User{}
}

func (u *User) Save(ctx context.Context, model *user.User) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection not found in context")
	}

	query := `INSERT INTO users (id, name, email, user_type, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())`
	_, err := db.ExecContext(ctx, query, model.ID, model.Name, model.Email, model.UserType)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (u *User) List(ctx context.Context) (user.UserSlice, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	query := `SELECT id, name, email, user_type, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get user list: %w", err)
	}
	defer rows.Close()

	var users user.UserSlice
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.UserType, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (u *User) Detail(ctx context.Context, model *user.User) (*user.User, error) {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	query := `SELECT id, name, email, user_type, created_at, updated_at FROM users WHERE id = ? AND deleted_at IS NULL`
	var result user.User
	err := db.QueryRowContext(ctx, query, model.ID).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.UserType,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user detail: %w", err)
	}

	return &result, nil
}

func (u *User) Delete(ctx context.Context, model *user.User) error {
	db := Ctx.GetDB(ctx)
	if db == nil {
		return fmt.Errorf("database connection not found in context")
	}

	query := `UPDATE users SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result, err := db.ExecContext(ctx, query, model.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
