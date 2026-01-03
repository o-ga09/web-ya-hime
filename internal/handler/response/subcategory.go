package response

import "time"

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
