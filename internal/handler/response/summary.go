package response

import SummaryDomain "github.com/o-ga09/web-ya-hime/internal/domain/summary"

// ListResponse はリスト取得のレスポンス構造体
type ListSummary struct {
	Summaries []*DetailSummary `json:"summaries"`
	Total     int              `json:"total"`
}

// DetailSummary はサマリーの詳細構造体
type DetailSummary struct {
	ID        string `json:"id"`
	User      *user  `json:"user"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
		ID:        s.ID,
		Content:   s.Content,
		CreatedAt: s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if s.User != nil {
		res.User = ToUserResponse(s.User)
	}

	return res
}
