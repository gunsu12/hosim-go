package item

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"
	"uuid"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/enums"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// ==========================================
// Category DTOs
// ==========================================

type CreateCategoryRequest struct {
	Code            string         `json:"code" binding:"required"`
	Name            string         `json:"name" binding:"required"`
	ItemType        enums.ItemType `json:"item_type" binding:"required"`
	IncomeCOACode   *string        `json:"income_coa_code"`
	DiscountCOACode *string        `json:"discount_coa_code"`
	SalesTaxCOACode *string        `json:"sales_tax_coa_code"`
	IsActive        *bool          `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Code            string         `json:"code" binding:"required"`
	Name            string         `json:"name" binding:"required"`
	ItemType        enums.ItemType `json:"item_type" binding:"required"`
	IncomeCOACode   *string        `json:"income_coa_code"`
	DiscountCOACode *string        `json:"discount_coa_code"`
	SalesTaxCOACode *string        `json:"sales_tax_coa_code"`
	IsActive        *bool          `json:"is_active"`
}

type CategoryResponse struct {
	ID              string         `json:"id"`
	Code            string         `json:"code"`
	Name            string         `json:"name"`
	ItemType        enums.ItemType `json:"item_type"`
	IncomeCOACode   *string        `json:"income_coa_code,omitempty"`
	DiscountCOACode *string        `json:"discount_coa_code,omitempty"`
	SalesTaxCOACode *string        `json:"sales_tax_coa_code,omitempty"`
	IsActive        bool           `json:"is_active"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func ToCategoryResponse(c *ItemCategory) CategoryResponse {
	return CategoryResponse{
		ID:              c.ID,
		Code:            c.Code,
		Name:            c.Name,
		ItemType:        c.ItemType,
		IncomeCOACode:   c.IncomeCOACode,
		DiscountCOACode: c.DiscountCOACode,
		SalesTaxCOACode: c.SalesTaxCOACode,
		IsActive:        c.IsActive,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

// ==========================================
// Category Service (Jalur A)
// ==========================================

type CategoryService interface {
	Create(ctx context.Context, req CreateCategoryRequest, actor string) (*CategoryResponse, error)
	Update(ctx context.Context, id string, req UpdateCategoryRequest, actor string) (*CategoryResponse, error)
	Delete(ctx context.Context, id string, actor string) error
	GetByID(ctx context.Context, id string) (*CategoryResponse, error)
	List(ctx context.Context, params CategoryListParams) ([]CategoryResponse, response.PaginationMeta, error)
}

type categoryService struct {
	repo Repository
}

func NewCategoryService(repo Repository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, req CreateCategoryRequest, actor string) (*CategoryResponse, error) {
	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)

	if code == "" || name == "" {
		return nil, errors.New("kode dan nama kategori wajib diisi")
	}
	if !req.ItemType.IsValid() {
		return nil, errors.New("item_type kategori tidak valid")
	}

	existing, err := s.repo.FindCategoryByCode(ctx, code)
	if err == nil && existing != nil {
		return nil, ErrCategoryCodeAlreadyExists
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	now := time.Now()
	cat := ItemCategory{
		ID:              uuid.NewV7().String(),
		Code:            code,
		Name:            name,
		ItemType:        req.ItemType,
		IncomeCOACode:   req.IncomeCOACode,
		DiscountCOACode: req.DiscountCOACode,
		SalesTaxCOACode: req.SalesTaxCOACode,
		IsActive:        isActive,
		CreatedAt:       now,
		CreatedBy:       actor,
		UpdatedAt:       now,
		UpdatedBy:       actor,
	}

	if err := s.repo.CreateCategory(ctx, &cat); err != nil {
		return nil, err
	}

	res := ToCategoryResponse(&cat)
	return &res, nil
}

func (s *categoryService) Update(ctx context.Context, id string, req UpdateCategoryRequest, actor string) (*CategoryResponse, error) {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	if code == "" || name == "" {
		return nil, errors.New("kode dan nama kategori wajib diisi")
	}
	if !req.ItemType.IsValid() {
		return nil, errors.New("item_type kategori tidak valid")
	}

	if code != cat.Code {
		existing, err := s.repo.FindCategoryByCode(ctx, code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrCategoryCodeAlreadyExists
		}
	}

	if req.IsActive != nil {
		cat.IsActive = *req.IsActive
	}

	cat.Code = code
	cat.Name = name
	cat.ItemType = req.ItemType
	cat.IncomeCOACode = req.IncomeCOACode
	cat.DiscountCOACode = req.DiscountCOACode
	cat.SalesTaxCOACode = req.SalesTaxCOACode
	cat.UpdatedAt = time.Now()
	cat.UpdatedBy = actor

	if err := s.repo.UpdateCategory(ctx, cat); err != nil {
		return nil, err
	}

	res := ToCategoryResponse(cat)
	return &res, nil
}

func (s *categoryService) Delete(ctx context.Context, id string, actor string) error {
	_, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteCategory(ctx, id, actor)
}

func (s *categoryService) GetByID(ctx context.Context, id string) (*CategoryResponse, error) {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := ToCategoryResponse(cat)
	return &res, nil
}

func (s *categoryService) List(ctx context.Context, params CategoryListParams) ([]CategoryResponse, response.PaginationMeta, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	categories, total, err := s.repo.FindAllCategories(ctx, params)
	if err != nil {
		return nil, response.PaginationMeta{}, err
	}

	res := make([]CategoryResponse, len(categories))
	for i := range categories {
		res[i] = ToCategoryResponse(&categories[i])
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}

	meta := response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return res, meta, nil
}

// ==========================================
// Category HTTP Handler
// ==========================================

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) RegisterRoutes(router *gin.RouterGroup) {
	categories := router.Group("/item-categories")
	{
		categories.GET("", middleware.RequirePermission("item:read"), h.List)
		categories.GET("/:id", middleware.RequirePermission("item:read"), h.GetByID)
		categories.POST("", middleware.RequirePermission("item:create"), h.Create)
		categories.PUT("/:id", middleware.RequirePermission("item:update"), h.Update)
		categories.DELETE("/:id", middleware.RequirePermission("item:delete"), h.Delete)
	}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload kategori gagal: "+err.Error(), nil)
		return
	}

	actor := getOperator(c)
	result, err := h.service.Create(c.Request.Context(), req, actor)
	if err != nil {
		if errors.Is(err, ErrCategoryCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Kategori item berhasil dibuat", result)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload kategori gagal: "+err.Error(), nil)
		return
	}

	actor := getOperator(c)
	result, err := h.service.Update(c.Request.Context(), id, req, actor)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCategoryCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Kategori item berhasil diperbarui", result)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	actor := getOperator(c)

	err := h.service.Delete(c.Request.Context(), id, actor)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Kategori item berhasil dihapus", nil)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Detail kategori berhasil dimuat", result)
}

type QueryCategoriesRequest struct {
	Page     int             `form:"page,default=1"`
	Limit    int             `form:"limit,default=10"`
	Search   string          `form:"search"`
	ItemType *enums.ItemType `form:"item_type"`
	IsActive *bool           `form:"is_active"`
}

func (h *CategoryHandler) List(c *gin.Context) {
	var req QueryCategoriesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter query tidak valid: "+err.Error(), nil)
		return
	}

	params := CategoryListParams{
		Page:     req.Page,
		Limit:    req.Limit,
		Search:   req.Search,
		ItemType: req.ItemType,
		IsActive: req.IsActive,
	}

	categories, meta, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Daftar kategori berhasil dimuat", categories, meta)
}
