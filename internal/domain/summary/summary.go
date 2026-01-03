package summary

import (
	"context"
	"database/sql"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
)

type ISummaryRepository interface {
	Save(ctx context.Context, model *Summary) error
	List(ctx context.Context, opts ListOptions) (*ListResult, error)
	Detail(ctx context.Context, model *Summary) (*Summary, error)
	Delete(ctx context.Context, model *Summary) error
}

type ListOptions struct {
	Category string
	Limit    int
	Offset   int
}

type ListResult struct {
	Items   SummarySlice `json:"items"`
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
	HasNext bool         `json:"has_next"`
}

type Summary struct {
	domain.WYHBaseModel
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Content     string         `json:"content"`
	Category    sql.NullString `json:"category"`
	UserID      string         `json:"user_id"`
	User        *user.User     `json:"user"`
}

type SummarySlice []*Summary
