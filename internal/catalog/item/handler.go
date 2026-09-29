package item

import (
	"errors"
	"net/http"

	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	createMedicationUC CreateMedicationUseCase
	createGeneralUC    CreateGeneralUseCase
	createAssetUC      CreateAssetUseCase
	createTariffUC     CreateTariffUseCase
	updateItemUC       UpdateItemUseCase
	manageUnitUC       ManageUnitUseCase
	getItemUC          GetItemUseCase
}

func NewHandler(
	createMedicationUC CreateMedicationUseCase,
	createGeneralUC CreateGeneralUseCase,
	createAssetUC CreateAssetUseCase,
	createTariffUC CreateTariffUseCase,
	updateItemUC UpdateItemUseCase,
	manageUnitUC ManageUnitUseCase,
	getItemUC GetItemUseCase,
) *Handler {
	return &Handler{
		createMedicationUC: createMedicationUC,
		createGeneralUC:    createGeneralUC,
		createAssetUC:      createAssetUC,
		createTariffUC:     createTariffUC,
		updateItemUC:       updateItemUC,
		manageUnitUC:       manageUnitUC,
		getItemUC:          getItemUC,
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
	items := router.Group("/items")
	{
		items.GET("", middleware.RequirePermission("item:read"), h.List)
		items.GET("/:id", middleware.RequirePermission("item:read"), h.GetByID)
		items.DELETE("/:id", middleware.RequirePermission("item:delete"), h.DeleteItem)

		items.POST("/medications", middleware.RequirePermission("item:create"), h.CreateMedication)
		items.PUT("/medications/:id", middleware.RequirePermission("item:update"), h.UpdateMedication)

		items.POST("/generals", middleware.RequirePermission("item:create"), h.CreateGeneral)
		items.PUT("/generals/:id", middleware.RequirePermission("item:update"), h.UpdateGeneral)

		items.POST("/assets", middleware.RequirePermission("item:create"), h.CreateAsset)
		items.PUT("/assets/:id", middleware.RequirePermission("item:update"), h.UpdateAsset)

		items.POST("/tariffs", middleware.RequirePermission("item:create"), h.CreateTariff)
		items.PUT("/tariffs/:id", middleware.RequirePermission("item:update"), h.UpdateTariff)

		items.GET("/:id/units", middleware.RequirePermission("item:read"), h.ListUnits)
		items.POST("/:id/units", middleware.RequirePermission("item:update"), h.AddUnit)
		items.PUT("/:id/units/:unit_id", middleware.RequirePermission("item:update"), h.UpdateUnit)
		items.DELETE("/:id/units/:unit_id", middleware.RequirePermission("item:update"), h.DeleteUnit)
	}
}

// CreateMedication menangani registrasi master obat
func (h *Handler) CreateMedication(c *gin.Context) {
	var req CreateMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload obat gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.createMedicationUC.Execute(c.Request.Context(), req, operator)
	if err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCategoryNotFound) || errors.Is(err, ErrProductLineNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Item obat berhasil didaftarkan", result)
}

// CreateGeneral menangani registrasi master barang umum / BMHP
func (h *Handler) CreateGeneral(c *gin.Context) {
	var req CreateGeneralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload barang umum/BMHP gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.createGeneralUC.Execute(c.Request.Context(), req, operator)
	if err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCategoryNotFound) || errors.Is(err, ErrProductLineNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Item umum/BMHP berhasil didaftarkan", result)
}

// CreateAsset menangani registrasi master alat kesehatan / barang modal
func (h *Handler) CreateAsset(c *gin.Context) {
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload aset modal gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.createAssetUC.Execute(c.Request.Context(), req, operator)
	if err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCategoryNotFound) || errors.Is(err, ErrProductLineNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Item aset berhasil didaftarkan", result)
}

// CreateTariff menangani registrasi master layanan / tarif jasa
func (h *Handler) CreateTariff(c *gin.Context) {
	var req CreateTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload tarif gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.createTariffUC.Execute(c.Request.Context(), req, operator)
	if err != nil {
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Item tarif layanan berhasil didaftarkan", result)
}

// AddUnit menangani penambahan satuan alternatif (UOM)
func (h *Handler) AddUnit(c *gin.Context) {
	itemID := c.Param("id")
	var req AddUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi satuan gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	unit, err := h.manageUnitUC.AddUnit(c.Request.Context(), itemID, req, operator)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrDuplicateUnitName) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Satuan baru berhasil ditambahkan", unit)
}

// DeleteUnit menangani penghapusan satuan alternatif
func (h *Handler) DeleteUnit(c *gin.Context) {
	itemID := c.Param("id")
	unitID := c.Param("unit_id")

	err := h.manageUnitUC.DeleteUnit(c.Request.Context(), itemID, unitID)
	if err != nil {
		if errors.Is(err, ErrUnitNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrBaseUnitCannotBeDeleted) {
			response.Error(c, http.StatusForbidden, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Satuan berhasil dihapus", nil)
}

// GetByID mengambil detail lengkap item beserta subtipe dan satuan
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	item, err := h.getItemUC.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item tidak ditemukan", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Detail item berhasil dimuat", item)
}

// List melakukan pencarian katalog universal dengan filter dan pagination
func (h *Handler) List(c *gin.Context) {
	var req QueryItemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter query tidak valid: "+err.Error(), nil)
		return
	}

	items, meta, err := h.getItemUC.List(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Katalog item berhasil dimuat", items, meta)
}

// UpdateMedication menangani pembaruan data master obat
func (h *Handler) UpdateMedication(c *gin.Context) {
	id := c.Param("id")
	var req UpdateMedicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload obat gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.updateItemUC.UpdateMedication(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item obat tidak ditemukan", nil)
			return
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Data obat berhasil diperbarui", result)
}

// UpdateGeneral menangani pembaruan data barang umum / BMHP
func (h *Handler) UpdateGeneral(c *gin.Context) {
	id := c.Param("id")
	var req UpdateGeneralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload barang umum/BMHP gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.updateItemUC.UpdateGeneral(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item umum/BMHP tidak ditemukan", nil)
			return
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Data umum/BMHP berhasil diperbarui", result)
}

// UpdateAsset menangani pembaruan data alat modal / aset
func (h *Handler) UpdateAsset(c *gin.Context) {
	id := c.Param("id")
	var req UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload aset gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.updateItemUC.UpdateAsset(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item aset tidak ditemukan", nil)
			return
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Data aset berhasil diperbarui", result)
}

// UpdateTariff menangani pembaruan data tarif layanan
func (h *Handler) UpdateTariff(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTariffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload tarif gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	result, err := h.updateItemUC.UpdateTariff(c.Request.Context(), id, req, operator)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item tarif tidak ditemukan", nil)
			return
		}
		if errors.Is(err, ErrCodeAlreadyExists) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Data tarif layanan berhasil diperbarui", result)
}

// DeleteItem menangani soft-delete item
func (h *Handler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	operator := getOperator(c)

	err := h.updateItemUC.DeleteItem(c.Request.Context(), id, operator)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item tidak ditemukan", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Item berhasil dihapus", nil)
}

// ListUnits mengambil daftar satuan untuk suatu item
func (h *Handler) ListUnits(c *gin.Context) {
	itemID := c.Param("id")
	units, err := h.manageUnitUC.ListUnits(c.Request.Context(), itemID)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			response.Error(c, http.StatusNotFound, "Item tidak ditemukan", nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Daftar satuan berhasil dimuat", units)
}

// UpdateUnit menangani pembaruan data satuan alternatif
func (h *Handler) UpdateUnit(c *gin.Context) {
	itemID := c.Param("id")
	unitID := c.Param("unit_id")
	var req UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validasi payload satuan gagal: "+err.Error(), nil)
		return
	}

	operator := getOperator(c)
	unit, err := h.manageUnitUC.UpdateUnit(c.Request.Context(), itemID, unitID, req, operator)
	if err != nil {
		if errors.Is(err, ErrUnitNotFound) {
			response.Error(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		if errors.Is(err, ErrDuplicateUnitName) {
			response.Error(c, http.StatusConflict, err.Error(), nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "Satuan berhasil diperbarui", unit)
}
