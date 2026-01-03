package category

import (
	"context"

	"github.com/o-ga09/web-ya-hime/internal/domain"
)

type ICategoryRepository interface {
	Save(ctx context.Context, model *Category) error
	List(ctx context.Context) (CategorySlice, error)
	Detail(ctx context.Context, model *Category) (*Category, error)
	Delete(ctx context.Context, model *Category) error
}

type Category struct {
	domain.WYHBaseModel
	Name string `json:"name"`
}

type CategorySlice []*Category
