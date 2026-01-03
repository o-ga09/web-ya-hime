package response

import (
	"time"

	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
)

type SubcategoryResponse struct {
	ID         string            `json:"id"`
	CategoryID string            `json:"category_id"`
	Name       string            `json:"name"`
	Category   *CategoryResponse `json:"category,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

type SubcategoryListResponse struct {
	Subcategories []SubcategoryResponse `json:"subcategories"`
}

func ToSubCategoryResponse(subcatRes *subcategory.Subcategory) *SubcategoryResponse {
	res := &SubcategoryResponse{
		ID:         subcatRes.ID,
		CategoryID: subcatRes.CategoryID,
		Name:       subcatRes.Name,
		CreatedAt:  subcatRes.CreatedAt,
		UpdatedAt:  subcatRes.UpdatedAt,
	}
	if subcatRes.Category != nil {
		res.Category = ToCategoryResponse(subcatRes.Category)
	}
	return res
}

func ToSubcategoriesListResponse(subcats []*subcategory.Subcategory) []*SubcategoryResponse {
	res := make([]*SubcategoryResponse, len(subcats))
	for i, s := range subcats {
		res[i] = ToSubCategoryResponse(s)
	}
	return res
}
