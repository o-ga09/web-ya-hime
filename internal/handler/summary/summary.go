package summary

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	"github.com/o-ga09/web-ya-hime/internal/handler/request"
	"github.com/o-ga09/web-ya-hime/internal/handler/response"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
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

	// リクエストボディを読み取り
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// JSONをリクエスト構造体にデコード
	var req request.SaveSummaryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// バリデーション
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ドメインモデルに変換
	model := &summary.Summary{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		UserID:      req.UserID,
	}

	// リポジトリに保存
	if err := s.repo.Save(ctx, model); err != nil {
		http.Error(w, "Failed to save summary", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusCreated)
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
		http.Error(w, "Failed to get summary list", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, response.ListSummary{
		Summary: response.ToListSummary(summaries),
		Total:   len(summaries),
	})
}

func (s *summaryHandler) Detail(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// クエリパラメータからIDを取得
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// リクエスト構造体を作成してバリデーション
	req := request.DetailSummaryRequest{
		ID: id,
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
	httputil.Response(&w, http.StatusOK, response.DetailSummary{
		Summary: response.ToSummaryResponse(detail),
	})
}

func (s *summaryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// リクエストボディを読み取り
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// JSONをリクエスト構造体にデコード
	var req request.DeleteSummaryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// バリデーション
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
