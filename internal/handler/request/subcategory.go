package request

type SubcategorySaveRequest struct {
	ID         string `json:"id" path:"id"`
	CategoryID string `json:"category_id" validate:"required"`
	Name       string `json:"name" validate:"required"`
}

type SubcategoryListRequest struct {
	CategoryID string `json:"category_id" query:"category_id"`
}

type SubcategoryDetailRequest struct {
	ID string `json:"id" path:"id" validate:"required"`
}

type SubcategoryDeleteRequest struct {
	ID string `json:"id" path:"id" validate:"required"`
}
