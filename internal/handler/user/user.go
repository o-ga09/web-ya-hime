package user

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	"github.com/o-ga09/web-ya-hime/internal/handler/request"
	"github.com/o-ga09/web-ya-hime/internal/handler/response"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
)

type IUserHandler interface {
	Save(http.ResponseWriter, *http.Request)
	List(http.ResponseWriter, *http.Request)
	Detail(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}
type userHandler struct {
	repo user.IUserRepository
}

func New(repo user.IUserRepository) IUserHandler {
	return &userHandler{
		repo: repo,
	}
}

func (u *userHandler) Save(w http.ResponseWriter, r *http.Request) {
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
	var req request.SaveUserRequest
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
	model := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		UserType: req.UserType,
	}

	// リポジトリに保存
	if err := u.repo.Save(ctx, model); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusCreated)
}

func (u *userHandler) List(w http.ResponseWriter, r *http.Request) {
	// メソッドチェック
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	// リポジトリからリストを取得
	users, err := u.repo.List(ctx)
	if err != nil {
		http.Error(w, "Failed to get user list", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, response.ListUser{
		User:  response.ToListUser(users),
		Total: len(users),
	})
}

func (u *userHandler) Detail(w http.ResponseWriter, r *http.Request) {
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
	req := request.DetailUserRequest{
		ID: id,
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ドメインモデルを作成
	model := &user.User{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	// リポジトリから詳細を取得
	detail, err := u.repo.Detail(ctx, model)
	if err != nil {
		http.Error(w, "Failed to get user detail", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, response.DetailUser{
		User: response.ToUserResponse(detail),
	})
}

func (u *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	var req request.DeleteUserRequest
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
	model := &user.User{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	// リポジトリから削除
	if err := u.repo.Delete(ctx, model); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusNoContent)
}
