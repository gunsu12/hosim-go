package priceschema

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createSchemaUC     CreatePriceSchemaUseCase
	updateSchemaUC     UpdatePriceSchemaUseCase
	deactivateSchemaUC DeactivatePriceSchemaUseCase
	getSchemaUC        GetPriceSchemaUseCase
	listSchemasUC      ListPriceSchemasUseCase
	upsertRuleUC       UpsertPriceSchemaRuleUseCase
	deleteRuleUC       DeletePriceSchemaRuleUseCase
	listRulesUC        ListPriceSchemaRulesUseCase
	simulateTariffUC   SimulateTariffUseCase
}

func NewHandler(
	createSchemaUC CreatePriceSchemaUseCase,
	updateSchemaUC UpdatePriceSchemaUseCase,
	deactivateSchemaUC DeactivatePriceSchemaUseCase,
	getSchemaUC GetPriceSchemaUseCase,
	listSchemasUC ListPriceSchemasUseCase,
	upsertRuleUC UpsertPriceSchemaRuleUseCase,
	deleteRuleUC DeletePriceSchemaRuleUseCase,
	listRulesUC ListPriceSchemaRulesUseCase,
	simulateTariffUC SimulateTariffUseCase,
) *Handler {
	return &Handler{
		createSchemaUC:     createSchemaUC,
		updateSchemaUC:     updateSchemaUC,
		deactivateSchemaUC: deactivateSchemaUC,
		getSchemaUC:        getSchemaUC,
		listSchemasUC:      listSchemasUC,
		upsertRuleUC:       upsertRuleUC,
		deleteRuleUC:       deleteRuleUC,
		listRulesUC:        listRulesUC,
		simulateTariffUC:   simulateTariffUC,
	}
}

func getActor(c *gin.Context) OperatorContext {
	principal := middleware.GetPrincipal(c)
	corrID := c.GetHeader("X-Correlation-ID")
	if corrID == "" {
		corrID = c.GetHeader("X-Request-ID")
	}

	userID := principal.UserID
	if userID == "" {
		userID = "SYSTEM"
	}

	username := principal.Username
	if username == "" {
		username = userID
	}

	role := principal.Role
	if role == "" {
		role = "OPERATOR"
	}

	return OperatorContext{
		UserID:        userID,
		Username:      username,
		Role:          role,
		CorrelationID: corrID,
	}
}

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/finance/price-schemas")
	{
		// Simulasi perhitungan tarif
		group.POST("/simulate", middleware.RequireAnyPermission("price_schema:simulate", "price_schema:read", "billing:create", "outpatient:register"), h.SimulateTariff)

		// Pengelolaan Header Price Schema
		group.GET("", middleware.RequirePermission("price_schema:read"), h.ListSchemas)
		group.POST("", middleware.RequirePermission("price_schema:create"), h.CreateSchema)
		group.GET("/:id", middleware.RequirePermission("price_schema:read"), h.GetSchemaByID)
		group.PUT("/:id", middleware.RequirePermission("price_schema:update"), h.UpdateSchema)
		group.POST("/:id/deactivate", middleware.RequirePermission("price_schema:update"), h.DeactivateSchema)

		// Pengelolaan Rule Price Schema
		group.GET("/:id/rules", middleware.RequirePermission("price_schema:read"), h.ListRules)
		group.PUT("/:id/rules", middleware.RequirePermission("price_schema:update"), h.UpsertRule)
		group.DELETE("/:id/rules/:rule_id", middleware.RequirePermission("price_schema:delete"), h.DeleteRule)
	}
}

func (h *Handler) ListSchemas(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("q")

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		b := strings.EqualFold(activeStr, "true") || activeStr == "1"
		isActive = &b
	}

	params := ListSchemaParams{
		Search:   search,
		IsActive: isActive,
		Page:     page,
		Limit:    limit,
	}

	list, meta, err := h.listSchemasUC.Execute(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar price schema", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Daftar price schema berhasil dimuat", list, meta)
}

func isClientValidationError(err error) bool {
	return errors.Is(err, ErrSchemaCodeRequired) ||
		errors.Is(err, ErrSchemaNameRequired) ||
		errors.Is(err, ErrScopeExclusive) ||
		errors.Is(err, ErrFixedPriceRequiresItem) ||
		errors.Is(err, ErrFixedPriceNegative) ||
		errors.Is(err, ErrPercentMinimum) ||
		errors.Is(err, ErrNoneValueMustBeZero) ||
		errors.Is(err, ErrCostShareSumMismatch) ||
		errors.Is(err, ErrCostShareOutOfRange) ||
		errors.Is(err, ErrTaxOutOfRange) ||
		errors.Is(err, ErrInvalidAdjustmentType) ||
		errors.Is(err, ErrInvalidMedicationAdjustType) ||
		errors.Is(err, ErrDepartmentInactiveOrNotFound) ||
		errors.Is(err, ErrServiceUnitInactiveOrNotFound) ||
		errors.Is(err, ErrItemInactiveOrNotFound) ||
		errors.Is(err, ErrDuplicateRuleScope) ||
		errors.Is(err, ErrMultipleDefaultRules)
}

func (h *Handler) CreateSchema(c *gin.Context) {
	var req CreatePriceSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format payload tidak valid", err.Error())
		return
	}

	actor := getActor(c)
	res, err := h.createSchemaUC.Execute(c.Request.Context(), req, actor)
	if err != nil {
		if errors.Is(err, ErrSchemaCodeExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if isClientValidationError(err) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal membuat price schema", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Price schema berhasil dibuat", res)
}

func (h *Handler) GetSchemaByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.getSchemaUC.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil price schema", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Detail price schema berhasil dimuat", res)
}

func (h *Handler) UpdateSchema(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePriceSchemaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format payload tidak valid", err.Error())
		return
	}

	actor := getActor(c)
	res, err := h.updateSchemaUC.Execute(c.Request.Context(), id, req, actor)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrSchemaNameRequired) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal memperbarui price schema", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Price schema berhasil diperbarui", res)
}

func (h *Handler) DeactivateSchema(c *gin.Context) {
	id := c.Param("id")
	actor := getActor(c)

	res, err := h.deactivateSchemaUC.Execute(c.Request.Context(), id, actor)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal menonaktifkan price schema", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Price schema berhasil dinonaktifkan", res)
}

func (h *Handler) ListRules(c *gin.Context) {
	schemaID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	params := ListRuleParams{
		Level:         c.Query("level"),
		DepartementID: c.Query("departement_id"),
		ServiceUnitID: c.Query("service_unit_id"),
		ItemID:        c.Query("item_id"),
		Page:          page,
		Limit:         limit,
	}

	rules, meta, err := h.listRulesUC.Execute(c.Request.Context(), schemaID, params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mengambil daftar rule price schema", err.Error())
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Daftar rule price schema berhasil dimuat", rules, meta)
}

func (h *Handler) UpsertRule(c *gin.Context) {
	schemaID := c.Param("id")
	var req UpsertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format payload tidak valid", err.Error())
		return
	}

	actor := getActor(c)
	res, err := h.upsertRuleUC.Execute(c.Request.Context(), schemaID, req, actor)
	if err != nil {
		if errors.Is(err, ErrSchemaNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if isClientValidationError(err) {
			response.Error(c, http.StatusBadRequest, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal menyimpan rule price schema", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Rule price schema berhasil disimpan", res)
}

func (h *Handler) DeleteRule(c *gin.Context) {
	schemaID := c.Param("id")
	ruleID := c.Param("rule_id")
	actor := getActor(c)

	err := h.deleteRuleUC.Execute(c.Request.Context(), schemaID, ruleID, actor)
	if err != nil {
		if errors.Is(err, ErrRuleNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCannotDeleteDefaultRule) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "Gagal menghapus rule price schema", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Rule price schema berhasil dihapus", nil)
}

func (h *Handler) SimulateTariff(c *gin.Context) {
	var req SimulateTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Format payload simulasi tidak valid", err.Error())
		return
	}

	res, err := h.simulateTariffUC.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrTariffNotConfigured) || errors.Is(err, ErrSchemaNotFound) || errors.Is(err, ErrNoMatchingRule) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Simulasi tarif berhasil dihitung", res)
}
