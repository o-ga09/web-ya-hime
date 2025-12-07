package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	Ctx "github.com/o-ga09/web-ya-hime/pkg/context"
	"github.com/stretchr/testify/assert"
)

func TestSummaryRepository_Save(t *testing.T) {
	repo := NewSummaryRepository()

	tests := []struct {
		name    string
		summary *summary.Summary
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功ケース: サマリーが正常に保存される",
			summary: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "test-id-1",
				},
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
				UserID:      "user-id-1",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO summaries").
					WithArgs("test-id-1", "Test Title", "Test Description", "Test Content", "user-id-1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "失敗ケース: データベース接続がcontextに存在しない",
			summary: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "test-id-2",
				},
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
				UserID:      "user-id-1",
			},
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: データベースエラー",
			summary: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "test-id-3",
				},
				Title:       "Test Title",
				Description: "Test Description",
				Content:     "Test Content",
				UserID:      "user-id-1",
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO summaries").
					WithArgs("test-id-3", "Test Title", "Test Description", "Test Content", "user-id-1").
					WillReturnError(fmt.Errorf("db error"))
			},
			wantErr: true,
			errMsg:  "failed to save summary",
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

			err = repo.Save(ctx, tt.summary)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSummaryRepository_List(t *testing.T) {
	repo := NewSummaryRepository()
	now := time.Now()

	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    int
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功ケース: サマリー一覧を取得",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "content", "user_id", "created_at", "updated_at",
					"id", "name", "email", "user_type", "created_at", "updated_at",
				}).
					AddRow("summary-1", "Title 1", "Description 1", "Content 1", "user-1", now, now,
						"user-1", "User Name 1", "user1@example.com", "admin", now, now).
					AddRow("summary-2", "Title 2", "Description 2", "Content 2", "user-2", now, now,
						"user-2", "User Name 2", "user2@example.com", "user", now, now)

				mock.ExpectQuery("SELECT (.+) FROM summaries").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "成功ケース: 空の結果",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "content", "user_id", "created_at", "updated_at",
					"id", "name", "email", "user_type", "created_at", "updated_at",
				})
				mock.ExpectQuery("SELECT (.+) FROM summaries").
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
				mock.ExpectQuery("SELECT (.+) FROM summaries").
					WillReturnError(fmt.Errorf("query error"))
			},
			wantErr: true,
			errMsg:  "failed to get summary list",
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

func TestSummaryRepository_Detail(t *testing.T) {
	repo := NewSummaryRepository()
	now := time.Now()

	tests := []struct {
		name    string
		input   *summary.Summary
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
		check   func(t *testing.T, result *summary.Summary)
	}{
		{
			name: "成功ケース: サマリー詳細を取得",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "summary-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "title", "description", "content", "user_id", "created_at", "updated_at",
					"id", "name", "email", "user_type", "created_at", "updated_at",
				}).
					AddRow("summary-1", "Title 1", "Description 1", "Content 1", "user-1", now, now,
						"user-1", "User Name 1", "user1@example.com", "admin", now, now)

				mock.ExpectQuery("SELECT (.+) FROM summaries (.+) WHERE s.id = ?").
					WithArgs("summary-1").
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, result *summary.Summary) {
				assert.NotNil(t, result)
				assert.Equal(t, "summary-1", result.ID)
				assert.Equal(t, "Title 1", result.Title)
				assert.Equal(t, "Description 1", result.Description)
				assert.Equal(t, "Content 1", result.Content)
				assert.NotNil(t, result.User)
				assert.Equal(t, "user-1", result.User.ID)
				assert.Equal(t, "User Name 1", result.User.Name)
			},
		},
		{
			name: "失敗ケース: サマリーが見つからない",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "non-existent",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM summaries (.+) WHERE s.id = ?").
					WithArgs("non-existent").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  "summary not found",
		},
		{
			name: "失敗ケース: データベース接続がcontextに存在しない",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "summary-1",
				},
			},
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: クエリエラー",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "summary-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM summaries (.+) WHERE s.id = ?").
					WithArgs("summary-1").
					WillReturnError(fmt.Errorf("query error"))
			},
			wantErr: true,
			errMsg:  "failed to get summary detail",
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

func TestSummaryRepository_Delete(t *testing.T) {
	repo := NewSummaryRepository()

	tests := []struct {
		name    string
		input   *summary.Summary
		mockFn  func(mock sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功ケース: サマリーが削除される",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "summary-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM summaries WHERE id = ?").
					WithArgs("summary-1").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "失敗ケース: サマリーが見つからない",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "non-existent",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM summaries WHERE id = ?").
					WithArgs("non-existent").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			errMsg:  "summary not found",
		},
		{
			name: "失敗ケース: データベース接続がcontextに存在しない",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "summary-1",
				},
			},
			mockFn:  func(mock sqlmock.Sqlmock) {},
			wantErr: true,
			errMsg:  "database connection not found in context",
		},
		{
			name: "失敗ケース: 削除エラー",
			input: &summary.Summary{
				WYHBaseModel: domain.WYHBaseModel{
					ID: "summary-1",
				},
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM summaries WHERE id = ?").
					WithArgs("summary-1").
					WillReturnError(fmt.Errorf("delete error"))
			},
			wantErr: true,
			errMsg:  "failed to delete summary",
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
