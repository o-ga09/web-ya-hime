package response

import UserDomain "github.com/o-ga09/web-ya-hime/internal/domain/user"

type ListUser struct {
	User  []*user `json:"users"`
	Total int     `json:"total"`
}

type DetailUser struct {
	User *user `json:"user"`
}

type user struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	UserType  string `json:"user_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToListUser(users []*UserDomain.User) []*user {
	res := make([]*user, len(users))
	for i, u := range users {
		res[i] = ToUserResponse(u)
	}
	return res
}

func ToUserResponse(u *UserDomain.User) *user {
	return &user{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		UserType: u.UserType,
	}
}
