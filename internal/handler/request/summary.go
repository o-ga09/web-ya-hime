package request

import (
	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	nullvalue "github.com/o-ga09/web-ya-hime/pkg/null_value"
	"github.com/o-ga09/web-ya-hime/pkg/ptr"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

// SaveSummaryRequest は保存リクエストの構造体
type SaveSummaryRequest struct {
	ID            *string `json:"id,omitempty"`
	Title         string  `json:"title" validate:"required,max=255"`
	Description   string  `json:"description" validate:"max=5000"`
	Content       string  `json:"content" validate:"required"`
	CategoryID    *string `json:"category_id" validate:"omitempty"`
	SubcategoryID *string `json:"subcategory_id" validate:"omitempty"`
	UserID        string  `json:"user_id"`
}

// ListSummaryRequest はリスト取得リクエストの構造体
type ListSummaryRequest struct {
	Category      *string `query:"category"`
	CategoryID    *string `query:"category_id"`
	SubcategoryID *string `query:"subcategory_id"`
	Limit         int     `query:"limit" validate:"min=1,max=100"`
	Offset        int     `query:"offset" validate:"min=0"`
}

// DetailSummaryRequest は詳細取得リクエストの構造体
type DetailSummaryRequest struct {
	ID string `path:"id" validate:"required"`
}

// DeleteSummaryRequest は削除リクエストの構造体
type DeleteSummaryRequest struct {
	ID string `path:"id" validate:"required"`
}

func (s *SaveSummaryRequest) ToModel() *summary.Summary {
	id := uuid.GenerateID()
	if s.ID != nil {
		id = ptr.PtrToString(s.ID)
	}

	return &summary.Summary{
		WYHBaseModel: domain.WYHBaseModel{
			ID: id,
		},
		Title:         s.Title,
		Description:   s.Description,
		Content:       s.Content,
		CategoryID:    nullvalue.PointerToSqlString(s.CategoryID),
		SubcategoryID: nullvalue.PointerToSqlString(s.SubcategoryID),
		UserID:        s.UserID,
	}
}
