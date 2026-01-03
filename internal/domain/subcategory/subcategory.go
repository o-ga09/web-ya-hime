package subcategory

import (
	"context"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/category"
)

type ISubcategoryRepository interface {
	Save(ctx context.Context, model *Subcategory) error
	List(ctx context.Context, categoryID string) (SubcategorySlice, error)
	Detail(ctx context.Context, model *Subcategory) (*Subcategory, error)
	Delete(ctx context.Context, model *Subcategory) error
}

type Subcategory struct {
	domain.WYHBaseModel
	CategoryID string             `json:"category_id"`
	Name       string             `json:"name"`
	Category   *category.Category `json:"category,omitempty"`
}

type SubcategorySlice []*Subcategory
