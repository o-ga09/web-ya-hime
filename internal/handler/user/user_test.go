package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository はuser.IUserRepositoryのモック
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, model *user.User) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context) (user.UserSlice, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(user.UserSlice), args.Error(1)
}

func (m *MockUserRepository) Detail(ctx context.Context, model *user.User) (*user.User, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, model *user.User) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func TestUserHandler_Save(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           map[string]interface{}
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:   "成功ケース: ユーザーが正常に保存される",
			method: http.MethodPost,
			body: map[string]interface{}{
				"name":      "Test User",
				"email":     "test@example.com",
				"user_type": "admin",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]string
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res["user_id"])
			},
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodGet,
			body:           map[string]interface{}{},
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "失敗ケース: リクエストボディが不正",
			method: http.MethodPost,
			body: map[string]interface{}{
				"invalid": "data",
			},
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: バリデーションエラー（nameが必須）",
			method: http.MethodPost,
			body: map[string]interface{}{
				"email":     "test@example.com",
				"user_type": "admin",
			},
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: バリデーションエラー（emailが必須）",
			method: http.MethodPost,
			body: map[string]interface{}{
				"name":      "Test User",
				"user_type": "admin",
			},
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: リポジトリでエラー",
			method: http.MethodPost,
			body: map[string]interface{}{
				"name":      "Test User",
				"email":     "test@example.com",
				"user_type": "admin",
			},
			mockSetup: func(m *MockUserRepository) {
				m.On("Save", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/users", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Save(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		method         string
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:   "成功ケース: ユーザー一覧を取得",
			method: http.MethodGet,
			mockSetup: func(m *MockUserRepository) {
				users := user.UserSlice{
					&user.User{
						WYHBaseModel: domain.WYHBaseModel{
							ID:        "user-1",
							CreatedAt: now,
							UpdatedAt: now,
						},
						Name:     "User Name 1",
						Email:    "user1@example.com",
						UserType: "admin",
					},
					&user.User{
						WYHBaseModel: domain.WYHBaseModel{
							ID:        "user-2",
							CreatedAt: now,
							UpdatedAt: now,
						},
						Name:     "User Name 2",
						Email:    "user2@example.com",
						UserType: "user",
					},
				}
				m.On("List", mock.Anything).Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]interface{}
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, float64(2), res["total"])
				users := res["users"].([]interface{})
				assert.Len(t, users, 2)
			},
		},
		{
			name:   "成功ケース: 空のリストを返す",
			method: http.MethodGet,
			mockSetup: func(m *MockUserRepository) {
				m.On("List", mock.Anything).Return(user.UserSlice{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]interface{}
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, float64(0), res["total"])
			},
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodPost,
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "失敗ケース: リポジトリでエラー",
			method: http.MethodGet,
			mockSetup: func(m *MockUserRepository) {
				m.On("List", mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo)

			req := httptest.NewRequest(tt.method, "/users", nil)
			w := httptest.NewRecorder()

			handler.List(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Detail(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		method         string
		userID         string
		mockSetup      func(*MockUserRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:   "成功ケース: ユーザー詳細を取得",
			method: http.MethodGet,
			userID: "user-1",
			mockSetup: func(m *MockUserRepository) {
				userData := &user.User{
					WYHBaseModel: domain.WYHBaseModel{
						ID:        "user-1",
						CreatedAt: now,
						UpdatedAt: now,
					},
					Name:     "User Name 1",
					Email:    "user1@example.com",
					UserType: "admin",
				}
				m.On("Detail", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
					return u.ID == "user-1"
				})).Return(userData, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]interface{}
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				userObj := res["user"].(map[string]interface{})
				assert.Equal(t, "user-1", userObj["id"])
				assert.Equal(t, "User Name 1", userObj["name"])
			},
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodPost,
			userID:         "user-1",
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "失敗ケース: IDが空",
			method:         http.MethodGet,
			userID:         "",
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: リポジトリでエラー",
			method: http.MethodGet,
			userID: "user-1",
			mockSetup: func(m *MockUserRepository) {
				m.On("Detail", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
					return u.ID == "user-1"
				})).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo)

			url := "/users/{id}"
			req := httptest.NewRequest(tt.method, url, nil)
			req.SetPathValue("id", tt.userID)

			w := httptest.NewRecorder()

			handler.Detail(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.String())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		mockSetup      func(*MockUserRepository)
		expectedStatus int
	}{
		{
			name:   "成功ケース: ユーザーが削除される",
			method: http.MethodDelete,
			userID: "user-1",
			mockSetup: func(m *MockUserRepository) {
				m.On("Delete", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
					return u.ID == "user-1"
				})).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodGet,
			userID:         "user-1",
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "失敗ケース: IDが空",
			method:         http.MethodDelete,
			userID:         "",
			mockSetup:      func(m *MockUserRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: リポジトリでエラー",
			method: http.MethodDelete,
			userID: "user-1",
			mockSetup: func(m *MockUserRepository) {
				m.On("Delete", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
					return u.ID == "user-1"
				})).Return(errors.New("delete error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo)

			url := "/users/{id}"
			req := httptest.NewRequest(tt.method, url, nil)
			req.SetPathValue("id", tt.userID)

			w := httptest.NewRecorder()

			handler.Delete(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
