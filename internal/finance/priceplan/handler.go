package priceplan

import (
	"errors"
	"net/http"
	"strconv"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createPlanUC  CreatePricePlanUseCase
	updatePlanUC  UpdatePricePlanUseCase
	clonePlanUC   ClonePricePlanUseCase
	submitPlanUC  SubmitPricePlanUseCase
	approvePlanUC ApprovePricePlanUseCase
	activatePlanUC ActivatePricePlanUseCase
	archivePlanUC ArchivePricePlanUseCase
	getPlanUC     GetPricePlanUseCase
	manageItemsUC ManageItemsUseCase
	lookupTariffUC LookupTariffUseCase
}

func NewHandler(
	createPlanUC CreatePricePlanUseCase,
	updatePlanUC UpdatePricePlanUseCase,
	clonePlanUC ClonePricePlanUseCase,
	submitPlanUC SubmitPricePlanUseCase,
	approvePlanUC ApprovePricePlanUseCase,
	activatePlanUC ActivatePricePlanUseCase,
	archivePlanUC ArchivePricePlanUseCase,
	getPlanUC GetPricePlanUseCase,
	manageItemsUC ManageItemsUseCase,
	lookupTariffUC LookupTariffUseCase,
) *Handler {
	return &Handler{
		createPlanUC:   createPlanUC,
		updatePlanUC:   updatePlanUC,
		clonePlanUC:    clonePlanUC,
		submitPlanUC:   submitPlanUC,
		approvePlanUC:  approvePlanUC,
		activatePlanUC: activatePlanUC,
		archivePlanUC:  archivePlanUC,
		getPlanUC:      getPlanUC,
		manageItemsUC:  manageItemsUC,
		lookupTariffUC: lookupTariffUC,
	}
}

func getOperator(c *gin.Context) string {
	if op := c.GetHeader("X-User-ID"); op != "" {
		return op
	}
	if username := middleware.GetUsername(c); username != "" {
		return username
	}
	if uid := middleware.GetUserID(c); uid != "" {
		return uid
	}
	return "SYSTEM"
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	plans := router.Group("/price-plans")
	{
		// PRD § 6.4 Core Engine Lookup Tarif
		plans.POST("/lookup", middleware.RequireAnyPermission("price_plan:lookup", "price_plan:read", "billing:create", "outpatient:register"), h.LookupTariff)

		// PRD § 6.2 Pengelolaan Buku Tarif / Price Plan
		plans.GET("", middleware.RequirePermission("price_plan:read"), h.ListPlans)
		plans.POST("", middleware.RequirePermission("price_plan:create"), h.CreatePlan)
		plans.GET("/:id", middleware.RequirePermission("price_plan:read"), h.GetPlanByID)
		plans.PUT("/:id", middleware.RequirePermission("price_plan:update"), h.UpdatePlan)
		plans.POST("/:id/clone", middleware.RequirePermission("price_plan:create"), h.ClonePlan)
		plans.POST("/:id/submit", middleware.RequirePermission("price_plan:submit"), h.SubmitPlan)
		plans.POST("/:id/approve", middleware.RequirePermission("price_plan:approve"), h.ApprovePlan)
		plans.POST("/:id/activate", middleware.RequirePermission("price_plan:activate"), h.ActivatePlan)
		plans.POST("/:id/archive", middleware.RequirePermission("price_plan:archive"), h.ArchivePlan)

		// PRD § 6.3 Pengelolaan Item Tarif & Komponen di Dalam Buku Tarif
		plans.GET("/:id/items", middleware.RequirePermission("price_plan:read"), h.ListItems)
		plans.POST("/:id/items", middleware.RequirePermission("price_plan:update"), h.AddItem)
		plans.GET("/:id/items/:item_id", middleware.RequirePermission("price_plan:read"), h.GetItem)
		plans.PUT("/:id/items/:item_id", middleware.RequirePermission("price_plan:update"), h.UpdateItem)
		plans.DELETE("/:id/items/:item_id", middleware.RequirePermission("price_plan:update"), h.DeleteItem)
		plans.POST("/:id/items/batch-upsert", middleware.RequirePermission("price_plan:update"), h.BatchUpsertItems)
	}
}

// -------------------------------------------------------------
// Core Engine Lookup Handler (PRD § 6.4)
// -------------------------------------------------------------

func (h *Handler) LookupTariff(c *gin.Context) {
	var req LookupTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format parameter lookup tidak sesuai", err.Error())
		return
	}

	result, err := h.lookupTariffUC.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrTariffNotConfigured) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Tarif berhasil ditemukan", result)
}

// -------------------------------------------------------------
// Price Plan Header Handlers (PRD § 6.2)
// -------------------------------------------------------------

func (h *Handler) ListPlans(c *gin.Context) {
	params := PlanListParams{
		Page:   1,
		Limit:  10,
		Search: c.Query("search"),
	}

	if p := c.Query("page"); p != "" {
		if page, err := strconv.Atoi(p); err == nil && page > 0 {
			params.Page = page
		}
	}
	if l := c.Query("limit"); l != "" {
		if limit, err := strconv.Atoi(l); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100
			}
			params.Limit = limit
		}
	}
	if st := c.Query("status"); st != "" {
		status := PricePlanStatus(st)
		params.Status = &status
	}
	if cust := c.Query("customer_id"); cust != "" {
		params.CustomerID = &cust
	}
	if def := c.Query("is_default"); def != "" {
		b := def == "true" || def == "1"
		params.IsDefault = &b
	}

	plans, total, err := h.getPlanUC.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)
	response.SuccessWithMeta(c, http.StatusOK, "Daftar buku tarif berhasil diambil", plans, response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  int(totalPages),
	})
}

func (h *Handler) GetPlanByID(c *gin.Context) {
	id := c.Param("id")
	plan, err := h.getPlanUC.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Detail buku tarif berhasil diambil", plan)
}

func (h *Handler) CreatePlan(c *gin.Context) {
	var req CreatePricePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.createPlanUC.Execute(c.Request.Context(), req, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrEffectiveDateInvalid) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Buku tarif baru berhasil dibuat", result)
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePricePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.updatePlanUC.Execute(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPricePlanImmutable) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Metadata buku tarif berhasil diperbarui", result)
}

func (h *Handler) ClonePlan(c *gin.Context) {
	id := c.Param("id")
	var req ClonePricePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.clonePlanUC.Execute(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPricePlanCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Buku tarif berhasil dikloning ke draf baru", result)
}

func (h *Handler) SubmitPlan(c *gin.Context) {
	id := c.Param("id")
	operator := getOperator(c)

	result, err := h.submitPlanUC.Execute(c.Request.Context(), id, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidStateTransition) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Buku tarif berhasil diajukan untuk persetujuan", result)
}

func (h *Handler) ApprovePlan(c *gin.Context) {
	id := c.Param("id")
	operator := getOperator(c)

	result, err := h.approvePlanUC.Execute(c.Request.Context(), id, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrTariffComponentsMismatch) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidStateTransition) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Buku tarif berhasil disetujui dan dikunci secara permanen", result)
}

func (h *Handler) ActivatePlan(c *gin.Context) {
	id := c.Param("id")
	operator := getOperator(c)

	result, err := h.activatePlanUC.Execute(c.Request.Context(), id, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrDefaultPricePlanConflict) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidStateTransition) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Buku tarif berhasil diaktifkan menjadi acuan operasional", result)
}

func (h *Handler) ArchivePlan(c *gin.Context) {
	id := c.Param("id")
	operator := getOperator(c)

	result, err := h.archivePlanUC.Execute(c.Request.Context(), id, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrInvalidStateTransition) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Buku tarif berhasil diarsipkan", result)
}

// -------------------------------------------------------------
// Price Plan Items Handlers (PRD § 6.3)
// -------------------------------------------------------------

func (h *Handler) ListItems(c *gin.Context) {
	planID := c.Param("id")
	params := ItemListParams{
		Page:   1,
		Limit:  10,
		Search: c.Query("search"),
	}

	if p := c.Query("page"); p != "" {
		if page, err := strconv.Atoi(p); err == nil && page > 0 {
			params.Page = page
		}
	}
	if l := c.Query("limit"); l != "" {
		if limit, err := strconv.Atoi(l); err == nil && limit > 0 {
			if limit > 100 {
				limit = 100
			}
			params.Limit = limit
		}
	}
	if cid := c.Query("tariff_class_id"); cid != "" {
		params.TariffClassID = &cid
	}
	if act := c.Query("is_active"); act != "" {
		b := act == "true" || act == "1"
		params.IsActive = &b
	}

	items, total, err := h.manageItemsUC.ListItems(c.Request.Context(), planID, params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	totalPages := (total + int64(params.Limit) - 1) / int64(params.Limit)
	response.SuccessWithMeta(c, http.StatusOK, "Daftar item tarif berhasil diambil", items, response.PaginationMeta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalItems:  total,
		TotalPages:  int(totalPages),
	})
}

func (h *Handler) GetItem(c *gin.Context) {
	planID := c.Param("id")
	itemID := c.Param("item_id")

	item, err := h.manageItemsUC.GetItem(c.Request.Context(), planID, itemID)
	if err != nil {
		if errors.Is(err, ErrPricePlanItemNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Detail item tarif berhasil diambil", item)
}

func (h *Handler) AddItem(c *gin.Context) {
	planID := c.Param("id")
	var req AddItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.manageItemsUC.AddItem(c.Request.Context(), planID, req, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPricePlanImmutable) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrTariffComponentsMismatch) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrItemAlreadyInPlan) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Item tarif berhasil ditambahkan ke buku tarif", result)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	planID := c.Param("id")
	itemID := c.Param("item_id")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	result, err := h.manageItemsUC.UpdateItem(c.Request.Context(), planID, itemID, req, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) || errors.Is(err, ErrPricePlanItemNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPricePlanImmutable) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrTariffComponentsMismatch) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Item tarif berhasil diperbarui", result)
}

func (h *Handler) DeleteItem(c *gin.Context) {
	planID := c.Param("id")
	itemID := c.Param("item_id")
	operator := getOperator(c)

	err := h.manageItemsUC.DeleteItem(c.Request.Context(), planID, itemID, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) || errors.Is(err, ErrPricePlanItemNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPricePlanImmutable) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Item tarif berhasil dihapus dari buku tarif", nil)
}

func (h *Handler) BatchUpsertItems(c *gin.Context) {
	planID := c.Param("id")
	var req BatchUpsertItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi gagal: format data tidak sesuai", err.Error())
		return
	}

	operator := getOperator(c)
	err := h.manageItemsUC.BatchUpsert(c.Request.Context(), planID, req, operator)
	if err != nil {
		if errors.Is(err, ErrPricePlanNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrPricePlanImmutable) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrTariffComponentsMismatch) {
			response.Error(c, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Batch upsert item tarif berhasil diproses", nil)
}
