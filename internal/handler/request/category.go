package request

type CategorySaveRequest struct {
	ID   string `json:"id" path:"id"`
	Name string `json:"name" validate:"required"`
}

type CategoryDetailRequest struct {
	ID string `json:"id" path:"id" validate:"required"`
}

type CategoryDeleteRequest struct {
	ID string `json:"id" path:"id" validate:"required"`
}
