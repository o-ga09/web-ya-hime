package summary

import (
	"context"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
)

type ISummaryRepository interface {
	Save(ctx context.Context, model *Summary) error
	List(ctx context.Context) (SummarySlice, error)
	Detail(ctx context.Context, model *Summary) (*Summary, error)
	Delete(ctx context.Context, model *Summary) error
}

type Summary struct {
	domain.WYHBaseModel
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Content     string     `json:"content"`
	UserID      string     `json:"user_id"`
	User        *user.User `json:"user"`
}

type SummarySlice []*Summary
