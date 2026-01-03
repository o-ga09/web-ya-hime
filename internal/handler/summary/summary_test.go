package summary

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
	"github.com/o-ga09/web-ya-hime/internal/domain/summary"
	"github.com/o-ga09/web-ya-hime/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSummaryRepository はsummary.ISummaryRepositoryのモック
type MockSummaryRepository struct {
	mock.Mock
}

func (m *MockSummaryRepository) Save(ctx context.Context, model *summary.Summary) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockSummaryRepository) List(ctx context.Context, opts summary.ListOptions) (*summary.ListResult, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*summary.ListResult), args.Error(1)
}

func (m *MockSummaryRepository) Detail(ctx context.Context, model *summary.Summary) (*summary.Summary, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*summary.Summary), args.Error(1)
}

func (m *MockSummaryRepository) Delete(ctx context.Context, model *summary.Summary) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

// MockSubcategoryRepository はsubcategory.ISubcategoryRepositoryのモック
type MockSubcategoryRepository struct {
	mock.Mock
}

func (m *MockSubcategoryRepository) Save(ctx context.Context, model *subcategory.Subcategory) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func (m *MockSubcategoryRepository) List(ctx context.Context, categoryID string) (subcategory.SubcategorySlice, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(subcategory.SubcategorySlice), args.Error(1)
}

func (m *MockSubcategoryRepository) Detail(ctx context.Context, model *subcategory.Subcategory) (*subcategory.Subcategory, error) {
	args := m.Called(ctx, model)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*subcategory.Subcategory), args.Error(1)
}

func (m *MockSubcategoryRepository) Delete(ctx context.Context, model *subcategory.Subcategory) error {
	args := m.Called(ctx, model)
	return args.Error(0)
}

func TestSummaryHandler_Save(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           map[string]interface{}
		mockSetup      func(*MockSummaryRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:   "成功ケース: サマリーが正常に保存される",
			method: http.MethodPost,
			body: map[string]interface{}{
				"title":       "Test Title",
				"description": "Test Description",
				"content":     "Test Content",
				"category":    "雑談",
				"user_id":     "user-123",
			},
			mockSetup: func(m *MockSummaryRepository) {
				m.On("Save", mock.Anything, mock.AnythingOfType("*summary.Summary")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]string
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.NotEmpty(t, res["summary_id"])
			},
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodGet,
			body:           map[string]interface{}{},
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "失敗ケース: リクエストボディが不正",
			method: http.MethodPost,
			body: map[string]interface{}{
				"invalid": "data",
			},
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: バリデーションエラー（titleが必須）",
			method: http.MethodPost,
			body: map[string]interface{}{
				"content": "Test Content",
				"user_id": "user-123",
			},
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "失敗ケース: リポジトリでエラー",
			method: http.MethodPost,
			body: map[string]interface{}{
				"title":       "Test Title",
				"description": "Test Description",
				"content":     "Test Content",
				"category":    "雑談",
				"user_id":     "user-123",
			},
			mockSetup: func(m *MockSummaryRepository) {
				m.On("Save", mock.Anything, mock.AnythingOfType("*summary.Summary")).Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSummaryRepository)
			mockSubcatRepo := new(MockSubcategoryRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo, mockSubcatRepo)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/summaries", bytes.NewBuffer(bodyBytes))
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

func TestSummaryHandler_List(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		method         string
		mockSetup      func(*MockSummaryRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:   "成功ケース: サマリー一覧を取得",
			method: http.MethodGet,
			mockSetup: func(m *MockSummaryRepository) {
				summaries := summary.SummarySlice{
					&summary.Summary{
						WYHBaseModel: domain.WYHBaseModel{
							ID:        "summary-1",
							CreatedAt: now,
							UpdatedAt: now,
						},
						Title:       "Title 1",
						Description: "Description 1",
						Content:     "Content 1",
						Category:    sql.NullString{String: "雑談", Valid: true},
						UserID:      "user-1",
						User: &user.User{
							WYHBaseModel: domain.WYHBaseModel{
								ID:        "user-1",
								CreatedAt: now,
								UpdatedAt: now,
							},
							Name:     "User Name",
							Email:    "user@example.com",
							UserType: "admin",
						},
					},
				}
				result := &summary.ListResult{
					Items:   summaries,
					Total:   1,
					Limit:   20,
					Offset:  0,
					HasNext: false,
				}
				m.On("List", mock.Anything, mock.Anything).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]interface{}
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, float64(1), res["total"])
				summaries := res["summaries"].([]interface{})
				assert.Len(t, summaries, 1)
			},
		},
		{
			name:   "成功ケース: 空のリストを返す",
			method: http.MethodGet,
			mockSetup: func(m *MockSummaryRepository) {
				result := &summary.ListResult{
					Items:   summary.SummarySlice{},
					Total:   0,
					Limit:   20,
					Offset:  0,
					HasNext: false,
				}
				m.On("List", mock.Anything, mock.Anything).Return(result, nil)
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
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "失敗ケース: リポジトリでエラー",
			method: http.MethodGet,
			mockSetup: func(m *MockSummaryRepository) {
				m.On("List", mock.Anything, mock.Anything).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSummaryRepository)
			mockSubcatRepo := new(MockSubcategoryRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo, mockSubcatRepo)

			req := httptest.NewRequest(tt.method, "/summaries", nil)
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

func TestSummaryHandler_Detail(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		method         string
		summaryID      string
		mockSetup      func(*MockSummaryRepository)
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:      "成功ケース: サマリー詳細を取得",
			method:    http.MethodGet,
			summaryID: "summary-1",
			mockSetup: func(m *MockSummaryRepository) {
				summaryData := &summary.Summary{
					WYHBaseModel: domain.WYHBaseModel{
						ID:        "summary-1",
						CreatedAt: now,
						UpdatedAt: now,
					},
					Title:       "Title 1",
					Description: "Description 1",
					Content:     "Content 1",
					UserID:      "user-1",
					User: &user.User{
						WYHBaseModel: domain.WYHBaseModel{
							ID:        "user-1",
							CreatedAt: now,
							UpdatedAt: now,
						},
						Name:     "User Name",
						Email:    "user@example.com",
						UserType: "admin",
					},
				}
				m.On("Detail", mock.Anything, mock.MatchedBy(func(s *summary.Summary) bool {
					return s.ID == "summary-1"
				})).Return(summaryData, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var res map[string]interface{}
				err := json.Unmarshal([]byte(body), &res)
				assert.NoError(t, err)
				assert.Equal(t, "summary-1", res["id"])
				assert.Equal(t, "Content 1", res["content"])
			},
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodPost,
			summaryID:      "summary-1",
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "失敗ケース: IDが空",
			method:         http.MethodGet,
			summaryID:      "",
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "失敗ケース: リポジトリでエラー",
			method:    http.MethodGet,
			summaryID: "summary-1",
			mockSetup: func(m *MockSummaryRepository) {
				m.On("Detail", mock.Anything, mock.MatchedBy(func(s *summary.Summary) bool {
					return s.ID == "summary-1"
				})).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSummaryRepository)
			mockSubcatRepo := new(MockSubcategoryRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo, mockSubcatRepo)

			// パスパラメータをシミュレートするために、Go 1.22の新しいルーティングを使用
			url := "/summaries/{id}"
			req := httptest.NewRequest(tt.method, url, nil)
			req.SetPathValue("id", tt.summaryID)

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

func TestSummaryHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		summaryID      string
		mockSetup      func(*MockSummaryRepository)
		expectedStatus int
	}{
		{
			name:      "成功ケース: サマリーが削除される",
			method:    http.MethodDelete,
			summaryID: "summary-1",
			mockSetup: func(m *MockSummaryRepository) {
				m.On("Delete", mock.Anything, mock.MatchedBy(func(s *summary.Summary) bool {
					return s.ID == "summary-1"
				})).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "失敗ケース: メソッドが不正",
			method:         http.MethodGet,
			summaryID:      "summary-1",
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "失敗ケース: IDが空",
			method:         http.MethodDelete,
			summaryID:      "",
			mockSetup:      func(m *MockSummaryRepository) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "失敗ケース: リポジトリでエラー",
			method:    http.MethodDelete,
			summaryID: "summary-1",
			mockSetup: func(m *MockSummaryRepository) {
				m.On("Delete", mock.Anything, mock.MatchedBy(func(s *summary.Summary) bool {
					return s.ID == "summary-1"
				})).Return(errors.New("delete error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSummaryRepository)
			mockSubcatRepo := new(MockSubcategoryRepository)
			tt.mockSetup(mockRepo)

			handler := New(mockRepo, mockSubcatRepo)

			// パスパラメータをシミュレートするために、Go 1.22の新しいルーティングを使用
			url := "/summaries/{id}"
			req := httptest.NewRequest(tt.method, url, nil)
			req.SetPathValue("id", tt.summaryID)

			w := httptest.NewRecorder()

			handler.Delete(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
