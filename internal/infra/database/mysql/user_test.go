package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Save(t *testing.T) {
	repo := NewUserRepository()

	tests := []struct {
		name    string
		user    *user.User
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功ケース: ユーザーが正常に保存される",
			user: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "test-user-1",
				},
				Name:     "Test User",
				Email:    "test@example.com",
				UserType: "admin",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test-user-1", "Test User", "test@example.com", "admin").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "失敗ケース: データベース接続がcontextに存在しない",
			user: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "test-user-2",
				},
				Name:     "Test User",
				Email:    "test@example.com",
				UserType: "user",
			},
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: データベースエラー",
			user: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "test-user-3",
				},
				Name:     "Test User",
				Email:    "test@example.com",
				UserType: "admin",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("test-user-3", "Test User", "test@example.com", "admin").
					WillReturnError(fmt.Errorf("db error"))
			},
			wantErr: true,
			errMsg:  "failed to save user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockFn(mock)

			ctx := context.Background()
			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				ctx = Ctx.SetDB(ctx, db)
			}

			err = repo.Save(ctx, tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestUserRepository_List(t *testing.T) {
	repo := NewUserRepository()
	now := time.Now()

	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功ケース: ユーザー一覧を取得",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "user_type", "created_at", "updated_at",
				}).
					AddRow("user-1", "User Name 1", "user1@example.com", "admin", now, now).
					AddRow("user-2", "User Name 2", "user2@example.com", "user", now, now)

				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "成功ケース: 空の結果",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "user_type", "created_at", "updated_at",
				})
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name:    "失敗ケース: データベース接続がcontextに存在しない",
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: クエリエラー",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnError(fmt.Errorf("query error"))
			},
			wantErr: true,
			errMsg:  "failed to get user list",
		},
		{
			name: "失敗ケース: スキャンエラー",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "user_type", "created_at", "updated_at",
				}).
					AddRow("user-1", "User Name 1", "user1@example.com", "admin", "invalid-time", now)

				mock.ExpectQuery("SELECT (.+) FROM users").
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "failed to scan user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockFn(mock)

			ctx := context.Background()
			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				ctx = Ctx.SetDB(ctx, db)
			}

			result, err := repo.List(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.want)
			}

			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestUserRepository_Detail(t *testing.T) {
	repo := NewUserRepository()
	now := time.Now()

	tests := []struct {
		name    string
		input   *user.User
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
		check   func(t *testing.T, result *user.User)
	}{
		{
			name: "成功ケース: ユーザー詳細を取得",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "name", "email", "user_type", "created_at", "updated_at",
				}).
					AddRow("user-1", "User Name 1", "user1@example.com", "admin", now, now)

				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("user-1").
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, result *user.User) {
				assert.NotNil(t, result)
				assert.Equal(t, "user-1", result.ID)
				assert.Equal(t, "User Name 1", result.Name)
				assert.Equal(t, "user1@example.com", result.Email)
				assert.Equal(t, "admin", result.UserType)
			},
		},
		{
			name: "失敗ケース: ユーザーが見つからない",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "non-existent",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  "user not found",
		},
		{
			name: "失敗ケース: データベース接続がcontextに存在しない",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: クエリエラー",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id = ?").
					WithArgs("user-1").
					WillReturnError(fmt.Errorf("query error"))
			},
			wantErr: true,
			errMsg:  "failed to get user detail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockFn(mock)

			ctx := context.Background()
			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				ctx = Ctx.SetDB(ctx, db)
			}

			result, err := repo.Detail(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.check != nil {
					tt.check(t, result)
				}
			}

			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := NewUserRepository()

	tests := []struct {
		name    string
		input   *user.User
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功ケース: ユーザーが削除される",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = ?").
					WithArgs("user-1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "失敗ケース: ユーザーが見つからない",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "non-existent",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = ?").
					WithArgs("non-existent").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			errMsg:  "user not found",
		},
		{
			name: "失敗ケース: データベース接続がcontextに存在しない",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: 削除エラー",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = ?").
					WithArgs("user-1").
					WillReturnError(fmt.Errorf("delete error"))
			},
			wantErr: true,
			errMsg:  "failed to delete user",
		},
		{
			name: "失敗ケース: RowsAffected取得エラー",
			input: &user.User{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "user-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users WHERE id = ?").
					WithArgs("user-1").
					WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("rows affected error")))
			},
			wantErr: true,
			errMsg:  "failed to get rows affected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.mockFn(mock)

			ctx := context.Background()
			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				ctx = Ctx.SetDB(ctx, db)
			}

			err = repo.Delete(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			if tt.name != "失敗ケース: データベース接続がcontextに存在しない" {
				assert.NoError(t, mock.ExpectationsWereMet())
			}
		})
	}
}
