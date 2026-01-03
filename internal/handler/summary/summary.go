package summary

import (
	"net/http"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	"github.com/o-ga09/web-ya-hime/internal/handler/request"
	"github.com/o-ga09/web-ya-hime/internal/handler/response"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
	"github.com/o-ga09/web-ya-hime/pkg/ptr"
)

type ISummaryHandler interface {
	Save(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Detail(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type summaryHandler struct {
	repo       summary.ISummaryRepository
	subcatRepo subcategory.ISubcategoryRepository
}

func New(repo summary.ISummaryRepository, subcatRepo subcategory.ISubcategoryRepository) ISummaryHandler {
	return &summaryHandler{
		repo:       repo,
		subcatRepo: subcatRepo,
	}
}

func (s *summaryHandler) Save(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.SaveSummaryRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// カテゴリとサブカテゴリの組み合わせチェック
	if req.CategoryID != nil && req.SubcategoryID != nil {
		subcatModel := &subcategory.Subcategory{
			WYHBaseModel: domain.WYHBaseModel{
				ID: *req.SubcategoryID,
			},
		}
		subcat, err := s.subcatRepo.Detail(ctx, subcatModel)
		if err != nil {
			http.Error(w, "Invalid subcategory", http.StatusBadRequest)
			return
		}
		if subcat.CategoryID != *req.CategoryID {
			http.Error(w, "Subcategory does not belong to the specified category", http.StatusBadRequest)
			return
		}
	}

	// ドメインモデルに変換
	model := req.ToModel()

	// リポジトリに保存
	if err := s.repo.Save(ctx, model); err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to save summary", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, map[string]string{
		"summary_id": model.ID,
	})
}

func (s *summaryHandler) List(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.ListSummaryRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// デフォルト値の設定
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// リポジトリからリストを取得
	category := ptr.PtrToString(req.Category)
	categoryID := ptr.PtrToString(req.CategoryID)
	subcategoryID := ptr.PtrToString(req.SubcategoryID)

	// サブカテゴリのみ指定された場合、サブカテゴリが属するカテゴリで絞り込む
	if subcategoryID != "" && categoryID == "" {
		subcatModel := &subcategory.Subcategory{
			WYHBaseModel: domain.WYHBaseModel{
				ID: subcategoryID,
			},
		}
		subcat, err := s.subcatRepo.Detail(ctx, subcatModel)
		if err != nil {
			http.Error(w, "Invalid subcategory", http.StatusBadRequest)
			return
		}
		categoryID = subcat.CategoryID
	}

	// カテゴリとサブカテゴリの組み合わせチェック
	if categoryID != "" && subcategoryID != "" {
		subcatModel := &subcategory.Subcategory{
			WYHBaseModel: domain.WYHBaseModel{
				ID: subcategoryID,
			},
		}
		subcat, err := s.subcatRepo.Detail(ctx, subcatModel)
		if err != nil {
			http.Error(w, "Invalid subcategory", http.StatusBadRequest)
			return
		}
		if subcat.CategoryID != categoryID {
			http.Error(w, "Subcategory does not belong to the specified category", http.StatusBadRequest)
			return
		}
	}

	opts := summary.ListOptions{
		Category:      category,
		CategoryID:    categoryID,
		SubcategoryID: subcategoryID,
		Limit:         req.Limit,
		Offset:        req.Offset,
	}

	result, err := s.repo.List(ctx, opts)
	if err != nil {
		logger.Error(ctx, "error", err)
		http.Error(w, "Failed to get summary list", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	res := &response.ListSummary{
		Summaries: response.ToListSummary(result.Items),
		Total:     result.Total,
		Limit:     result.Limit,
		Offset:    result.Offset,
		HasNext:   result.HasNext,
	}
	httputil.Response(&w, http.StatusOK, res)
}

func (s *summaryHandler) Detail(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.DetailSummaryRequest
	// リクエスト構造体を作成してバリデーション
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ドメインモデルを作成
	model := &summary.Summary{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	// リポジトリから詳細を取得
	detail, err := s.repo.Detail(ctx, model)
	if err != nil {
		http.Error(w, "Failed to get summary detail", http.StatusInternalServerError)
		return
	}

	res := response.ToSummaryResponse(detail)
	httputil.Response(&w, http.StatusOK, res)
}

func (s *summaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.DeleteSummaryRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := &summary.Summary{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	if err := s.repo.Delete(ctx, model); err != nil {
		http.Error(w, "Failed to delete summary", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusNoContent, map[string]string{
		"message": "Summary deleted successfully",
	})
}
