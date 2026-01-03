package response

import (
	"time"

	"github.com/o-ga09/web-ya-hime/internal/domain/category"
)

type CategoryResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

func ToCategoryResponse(catRes *category.Category) *CategoryResponse {
	return &CategoryResponse{
		ID:        catRes.ID,
		Name:      catRes.Name,
		CreatedAt: catRes.CreatedAt,
		UpdatedAt: catRes.UpdatedAt,
	}
}

func ToCategoriesListResponse(cats []*category.Category) []*CategoryResponse {
	res := make([]*CategoryResponse, len(cats))
	for i, c := range cats {
		res[i] = ToCategoryResponse(c)
	}
	return res
}
