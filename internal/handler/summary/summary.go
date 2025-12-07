package summary

import (
	"net/http"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	"github.com/o-ga09/web-ya-hime/internal/handler/request"
	"github.com/o-ga09/web-ya-hime/internal/handler/response"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
)

type ISummaryHandler interface {
	Save(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Detail(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type summaryHandler struct {
	repo summary.ISummaryRepository
}

func New(repo summary.ISummaryRepository) ISummaryHandler {
	return &summaryHandler{
		repo: repo,
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

	// リポジトリからリストを取得
	summaries, err := s.repo.List(ctx)
	if err != nil {
		logger.Error(ctx, "error", err)
		http.Error(w, "Failed to get summary list", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, response.ToListSummary(summaries))
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

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, response.ToSummaryResponse(detail))
}

func (s *summaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.DeleteSummaryRequest
	// リクエスト構造体を作成してバリデーション
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ドメインモデルに変換
	model := &summary.Summary{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	// リポジトリから削除
	if err := s.repo.Delete(ctx, model); err != nil {
		http.Error(w, "Failed to delete summary", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusNoContent)
}
