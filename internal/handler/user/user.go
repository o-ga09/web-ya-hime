package user

import (
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

	var req request.SaveUserRequest
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
	if err := u.repo.Save(ctx, model); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	// レスポンスを返す
	httputil.Response(&w, http.StatusOK, map[string]string{
		"user_id": model.ID,
	})
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

	// リクエスト構造体を作成してバリデーション
	var req request.DetailUserRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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

	var req request.DeleteUserRequest
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
