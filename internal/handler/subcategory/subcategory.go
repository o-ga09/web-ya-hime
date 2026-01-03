package subcategory

import (
	"net/http"

	"github.com/o-ga09/web-ya-hime/internal/domain"
	"github.com/o-ga09/web-ya-hime/internal/domain/subcategory"
	"github.com/o-ga09/web-ya-hime/internal/handler/request"
	"github.com/o-ga09/web-ya-hime/internal/handler/response"
	"github.com/o-ga09/web-ya-hime/pkg/httputil"
	"github.com/o-ga09/web-ya-hime/pkg/logger"
	"github.com/o-ga09/web-ya-hime/pkg/uuid"
)

type ISubcategoryHandler interface {
	Save(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	Detail(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type subcategoryHandler struct {
	repo subcategory.ISubcategoryRepository
}

func New(repo subcategory.ISubcategoryRepository) ISubcategoryHandler {
	return &subcategoryHandler{
		repo: repo,
	}
}

func (h *subcategoryHandler) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.SubcategorySaveRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// IDが指定されていない場合は生成
	id := req.ID
	if id == "" {
		id = uuid.GenerateID()
	}

	model := &subcategory.Subcategory{
		WYHBaseModel: domain.WYHBaseModel{
			ID: id,
		},
		CategoryID: req.CategoryID,
		Name:       req.Name,
	}

	if err := h.repo.Save(ctx, model); err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to save subcategory", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusOK, map[string]string{
		"subcategory_id": model.ID,
	})
}

func (h *subcategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.SubcategoryListRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subcategories, err := h.repo.List(ctx, req.CategoryID)
	if err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to get subcategory list", http.StatusInternalServerError)
		return
	}

	var subcategoryResponses []response.SubcategoryResponse
	for _, subcat := range subcategories {
		subcatRes := response.SubcategoryResponse{
			ID:         subcat.ID,
			CategoryID: subcat.CategoryID,
			Name:       subcat.Name,
			CreatedAt:  subcat.CreatedAt,
			UpdatedAt:  subcat.UpdatedAt,
		}
		if subcat.Category != nil {
			subcatRes.Category = &response.CategoryResponse{
				ID:        subcat.Category.ID,
				Name:      subcat.Category.Name,
				CreatedAt: subcat.Category.CreatedAt,
				UpdatedAt: subcat.Category.UpdatedAt,
			}
		}
		subcategoryResponses = append(subcategoryResponses, subcatRes)
	}

	res := response.SubcategoryListResponse{
		Subcategories: subcategoryResponses,
	}
	httputil.Response(&w, http.StatusOK, res)
}

func (h *subcategoryHandler) Detail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.SubcategoryDetailRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := &subcategory.Subcategory{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	subcat, err := h.repo.Detail(ctx, model)
	if err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to get subcategory detail", http.StatusNotFound)
		return
	}

	res := response.SubcategoryResponse{
		ID:         subcat.ID,
		CategoryID: subcat.CategoryID,
		Name:       subcat.Name,
		CreatedAt:  subcat.CreatedAt,
		UpdatedAt:  subcat.UpdatedAt,
	}
	if subcat.Category != nil {
		res.Category = &response.CategoryResponse{
			ID:        subcat.Category.ID,
			Name:      subcat.Category.Name,
			CreatedAt: subcat.Category.CreatedAt,
			UpdatedAt: subcat.Category.UpdatedAt,
		}
	}
	httputil.Response(&w, http.StatusOK, res)
}

func (h *subcategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()

	var req request.SubcategoryDeleteRequest
	if err := request.Bind(r, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := request.Validate(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := &subcategory.Subcategory{
		WYHBaseModel: domain.WYHBaseModel{
			ID: req.ID,
		},
	}

	if err := h.repo.Delete(ctx, model); err != nil {
		logger.Error(ctx, err.Error())
		http.Error(w, "Failed to delete subcategory", http.StatusInternalServerError)
		return
	}

	httputil.Response(&w, http.StatusOK, map[string]string{
		"message": "Subcategory deleted successfully",
	})
}
