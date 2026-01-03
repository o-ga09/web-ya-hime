package category

import (
	"net/http"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/category"
	"github.com/o-ga09/web-ya-hime/internal/handler/request"
	"github.com/o-ga09/web-ya-hime/internal/handler/response"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

type ICategoryHandler interface {
	Save(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Detail(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	repo category.ICategoryRepository
}

func New(repo category.ICategoryRepository) ICategoryHandler {
	return &categoryHandler{
		repo: repo,
	}
}

func (h *categoryHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.CategorySaveRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPut && req.ID == "" {
		http.Error(w, "Category ID is required for update", http.StatusBadRequest)
		return
	}

	// IDが指定されていない場合は生成
	id := req.ID
	if id == "" {
		id = uuid.GenerateID()
	}

	model := &category.Category{
		WYHBaseModel: domain.WYHBaseModel{
			ID: id,
		},
		Name: req.Name,
	}

	if err := h.repo.Save(ctx, model); err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to save category", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusOK, map[string]string{
		"category_id": model.ID,
	})
}

func (h *categoryHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	categories, err := h.repo.List(ctx)
	if err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to get category list", http.StatusInternalServerError)
		return
	}

	var categoryResponses []response.CategoryResponse
	for _, cat := range categories {
		categoryResponses = append(categoryResponses, response.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		})
	}

	res := response.CategoryListResponse{
		Categories: categoryResponses,
	}
	httputil.Response(&w, http.StatusOK, res)
}

func (h *categoryHandler) Detail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.CategoryDetailRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := &category.Category{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	cat, err := h.repo.Detail(ctx, model)
	if err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to get category detail", http.StatusNotFound)
		return
	}

	res := response.CategoryResponse{
		ID:        cat.ID,
		Name:      cat.Name,
		CreatedAt: cat.CreatedAt,
		UpdatedAt: cat.UpdatedAt,
	}
	httputil.Response(&w, http.StatusOK, res)
}

func (h *categoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.CategoryDeleteRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := &category.Category{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	if err := h.repo.Delete(ctx, model); err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusOK, map[string]string{
		"message": "Category deleted successfully",
	})
}
