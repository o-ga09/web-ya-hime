package request

import (
	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	"github.com/o-ga09/web-ya-hime/pkg/ptr"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

// SaveUserRequest は保存リクエストの構造体
type SaveUserRequest struct {
	ID       *string `json:"id,omitempty"`
	Name     string  `json:"name" validate:"required,max=100"`
	Email    string  `json:"email" validate:"required,max=255"`
	UserType string  `json:"user_type" validate:"required"`
}

// ListUserRequest はリスト取得リクエストの構造体
type ListUserRequest struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

// DetailUserRequest は詳細取得リクエストの構造体
type DetailUserRequest struct {
	ID string `path:"id" validate:"required"`
}

// DeleteUserRequest は削除リクエストの構造体
type DeleteUserRequest struct {
	ID string `path:"id" validate:"required"`
}

func (s *SaveUserRequest) ToModel() *user.User {
	id := uuid.GenerateID()
	if s.ID != nil {
		id = ptr.PtrToString(s.ID)
	}

	return &user.User{
		WYHBaseModel: domain.WYHBaseModel{
			ID: id,
		},
		Name:     s.Name,
		Email:    s.Email,
		UserType: s.UserType,
	}
}
