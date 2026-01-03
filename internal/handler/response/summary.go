package response

import (
	SummaryDomain "github.com/o-ga09/web-ya-hime/internal/domain/summary"
	"github.com/o-ga09/web-ya-hime/pkg/date"
)

// ListResponse はリスト取得のレスポンス構造体
type ListSummary struct {
	Summaries []*DetailSummary `json:"summaries"`
	Total     int              `json:"total"`
	Limit     int              `json:"limit"`
	Offset    int              `json:"offset"`
	HasNext   bool             `json:"has_next"`
}

// DetailSummary はサマリーの詳細構造体
type DetailSummary struct {
	ID          string               `json:"id"`
	User        *user                `json:"user"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Content     string               `json:"content"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	Category    *CategoryResponse    `json:"category,omitempty"`
	SubCategory *SubcategoryResponse `json:"subcategory,omitempty"`
}

func ToListSummary(summaries []*SummaryDomain.Summary) []*DetailSummary {
	res := make([]*DetailSummary, len(summaries))
	for i, s := range summaries {
		res[i] = ToSummaryResponse(s)
	}
	return res
}

func ToSummaryResponse(s *SummaryDomain.Summary) *DetailSummary {
	res := &DetailSummary{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		Content:     s.Content,
		CreatedAt:   date.FormatDefault(s.CreatedAt),
		UpdatedAt:   date.FormatDefault(s.UpdatedAt),
	}

	if s.User != nil {
		res.User = ToUserResponse(s.User)
	}
	if s.Category != nil {
		res.Category = ToCategoryResponse(s.Category)
	}
	if s.Subcategory != nil {
		res.SubCategory = ToSubCategoryResponse(s.Subcategory)
	}

	return res
}
