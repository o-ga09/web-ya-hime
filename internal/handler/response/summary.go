package response

import SummaryDomain "github.com/o-ga09/web-ya-hime/internal/domain/summary"

// ListResponse はリスト取得のレスポンス構造体
type ListSummary struct {
	Summary []*summary `json:"summaries"`
	Total   int        `json:"total"`
}

// DetailResponse は詳細取得のレスポンス構造体
type DetailSummary struct {
	Summary *summary `json:"data"`
}

type summary struct {
	ID        string `json:"id"`
	User      *user  `json:"user"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToListSummary(summaries []*SummaryDomain.Summary) []*summary {
	res := make([]*summary, len(summaries))
	for i, s := range summaries {
		res[i] = ToSummaryResponse(s)
	}
	return res
}

func ToSummaryResponse(s *SummaryDomain.Summary) *summary {
	res := &summary{
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
