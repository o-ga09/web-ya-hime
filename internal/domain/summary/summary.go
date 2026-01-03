package summary

import (
	"context"
	"database/sql"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/category"
	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
)

type ISummaryRepository interface {
	Save(ctx context.Context, model *Summary) error
	List(ctx context.Context, opts ListOptions) (*ListResult, error)
	Detail(ctx context.Context, model *Summary) (*Summary, error)
	Delete(ctx context.Context, model *Summary) error
}

type ListOptions struct {
	Category      string
	CategoryID    string
	SubcategoryID string
	Limit         int
	Offset        int
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
	Title         string                   `json:"title"`
	Description   string                   `json:"description"`
	Content       string                   `json:"content"`
	CategoryID    sql.NullString           `json:"category_id"`
	SubcategoryID sql.NullString           `json:"subcategory_id"`
	UserID        string                   `json:"user_id"`
	User          *user.User               `json:"user,omitempty"`
	Category      *category.Category       `json:"category,omitempty"`
	Subcategory   *subcategory.Subcategory `json:"subcategory,omitempty"`
}

type SummarySlice []*Summary
