package user

import (
	"context"

	"github.com/o-ga09/web-ya-hime/internal/domain"
)

type IUserRepository interface {
	Save(ctx context.Context, model *User) error
	List(ctx context.Context) (UserSlice, error)
	Detail(ctx context.Context, model *User) (*User, error)
	Delete(ctx context.Context, model *User) error
}

type User struct {
	domain.WYHBaseModel
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}
type UserSlice []*User
