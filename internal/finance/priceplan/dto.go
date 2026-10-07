package priceplan

import (
	"time"
)

// ==========================================
// Request DTOs
// ==========================================

// CreatePricePlanRequest merepresentasikan payload pembuatan buku tarif baru
type CreatePricePlanRequest struct {
	Code               string   `json:"code" binding:"required"`
	Name               string   `json:"name" binding:"required"`
	Description        *string  `json:"description"`
	EffectiveFrom      string   `json:"effective_from" binding:"required"` // Format: YYYY-MM-DD
	EffectiveTo        *string  `json:"effective_to"`                     // Format: YYYY-MM-DD (opsional)
	IsDefault          bool     `json:"is_default"`
	CustomerID         *string  `json:"customer_id"`
	DefaultCitoPercent *float64 `json:"default_cito_percent"`
}

// UpdatePricePlanRequest merepresentasikan payload update metadata buku tarif (hanya DRAFT)
type UpdatePricePlanRequest struct {
	Name               string   `json:"name" binding:"required"`
	Description        *string  `json:"description"`
	EffectiveFrom      string   `json:"effective_from" binding:"required"` // Format: YYYY-MM-DD
	EffectiveTo        *string  `json:"effective_to"`                     // Format: YYYY-MM-DD (opsional)
	IsDefault          bool     `json:"is_default"`
	CustomerID         *string  `json:"customer_id"`
	DefaultCitoPercent *float64 `json:"default_cito_percent"`
}

// ClonePricePlanRequest merepresentasikan payload kloning buku tarif
type ClonePricePlanRequest struct {
	NewCode            string   `json:"new_code" binding:"required"`
	NewName            string   `json:"new_name" binding:"required"`
	Description        *string  `json:"description"`
	EffectiveFrom      string   `json:"effective_from" binding:"required"`
	EffectiveTo        *string  `json:"effective_to"`
	IsDefault          bool     `json:"is_default"`
	CustomerID         *string  `json:"customer_id"`
	DefaultCitoPercent *float64 `json:"default_cito_percent"`
}

// ItemComponentRequest mendefinisikan rincian pecahan komponen biaya pada suatu tarif item
type ItemComponentRequest struct {
	ComponentID string   `json:"component_id" binding:"required"`
	BaseAmount  float64  `json:"base_amount" binding:"gte=0"`
	CitoAmount  *float64 `json:"cito_amount"`
	COACode     *string  `json:"coa_code"`
}

// AddItemRequest mendefinisikan payload penambahan item tarif ke buku tarif
type AddItemRequest struct {
	ItemID         string                 `json:"item_id" binding:"required"`
	TariffClassID  string                 `json:"tariff_class_id" binding:"required"`
	TotalBasePrice float64                `json:"total_base_price" binding:"gte=0"`
	TotalCitoPrice *float64               `json:"total_cito_price"`
	IsActive       *bool                  `json:"is_active"`
	Components     []ItemComponentRequest `json:"components" binding:"required,dive"`
}

// UpdateItemRequest mendefinisikan payload pembaruan tarif dan komponen item
type UpdateItemRequest struct {
	TotalBasePrice float64                `json:"total_base_price" binding:"gte=0"`
	TotalCitoPrice *float64               `json:"total_cito_price"`
	IsActive       *bool                  `json:"is_active"`
	Components     []ItemComponentRequest `json:"components" binding:"required,dive"`
}

// BatchUpsertItemsRequest mendefinisikan payload impor masal item tarif
type BatchUpsertItemsRequest struct {
	Items []AddItemRequest `json:"items" binding:"required,dive"`
}

// LookupTariffRequest mendefinisikan payload pencarian tarif transaksi kasir / billing
type LookupTariffRequest struct {
	CustomerID      *string `json:"customer_id"`
	TariffClassID   string  `json:"tariff_class_id" binding:"required"`
	ItemID          string  `json:"item_id" binding:"required"`
	TransactionDate string  `json:"transaction_date" binding:"required"` // Format: RFC3339 atau YYYY-MM-DD
	IsCito          bool    `json:"is_cito"`
}

// ==========================================
// Filter / List Query Params
// ==========================================

// PlanListParams mendefinisikan parameter pencarian list buku tarif
type PlanListParams struct {
	Page       int
	Limit      int
	Search     string
	Status     *PricePlanStatus
	CustomerID *string
	IsDefault  *bool
}

// ItemListParams mendefinisikan parameter pencarian item dalam buku tarif
type ItemListParams struct {
	Page          int
	Limit         int
	Search        string
	ItemID        *string
	TariffClassID *string
	IsActive      *bool
}

// ==========================================
// Response DTOs
// ==========================================

// PricePlanResponse adalah representasi output buku tarif
type PricePlanResponse struct {
	ID                 string          `json:"id"`
	Code               string          `json:"code"`
	Name               string          `json:"name"`
	Description        *string         `json:"description,omitempty"`
	EffectiveFrom      string          `json:"effective_from"`
	EffectiveTo        *string         `json:"effective_to,omitempty"`
	Status             PricePlanStatus `json:"status"`
	IsDefault          bool            `json:"is_default"`
	CustomerID         *string         `json:"customer_id,omitempty"`
	CustomerName       *string         `json:"customer_name,omitempty"`
	DefaultCitoPercent float64         `json:"default_cito_percent"`
	ApprovedAt         *time.Time      `json:"approved_at,omitempty"`
	ApprovedBy         *string         `json:"approved_by,omitempty"`
	TotalItems         int64           `json:"total_items"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

// PricePlanItemResponse adalah representasi output baris item tarif
type PricePlanItemResponse struct {
	ID              string                           `json:"id"`
	PricePlanID     string                           `json:"price_plan_id"`
	ItemID          string                           `json:"item_id"`
	ItemCode        string                           `json:"item_code,omitempty"`
	ItemName        string                           `json:"item_name,omitempty"`
	TariffClassID   string                           `json:"tariff_class_id"`
	TariffClassCode string                           `json:"tariff_class_code,omitempty"`
	TariffClassName string                           `json:"tariff_class_name,omitempty"`
	TotalBasePrice  float64                          `json:"total_base_price"`
	TotalCitoPrice  *float64                         `json:"total_cito_price,omitempty"`
	IsActive        bool                             `json:"is_active"`
	Components      []PricePlanItemComponentResponse `json:"components,omitempty"`
}

// PricePlanItemComponentResponse adalah representasi output komponen biaya
type PricePlanItemComponentResponse struct {
	ID            string   `json:"id"`
	PlanItemID    string   `json:"plan_item_id"`
	ComponentID   string   `json:"component_id"`
	ComponentCode string   `json:"component_code,omitempty"`
	ComponentName string   `json:"component_name,omitempty"`
	ComponentType string   `json:"component_type,omitempty"`
	BaseAmount    float64  `json:"base_amount"`
	CitoAmount    *float64 `json:"cito_amount,omitempty"`
	COACode       *string  `json:"coa_code,omitempty"`
}

// ResolvedComponentResponse adalah representasi rincian komponen hasil lookup
type ResolvedComponentResponse struct {
	ComponentID   string  `json:"component_id"`
	ComponentName string  `json:"component_name"`
	ComponentType string  `json:"component_type"`
	Amount        float64 `json:"amount"`
	COACode       string  `json:"coa_code"`
}

// LookupTariffResponse adalah representasi output resolusi lookup tarif (PRD § 6.4)
type LookupTariffResponse struct {
	PricePlanID     string                      `json:"price_plan_id"`
	PricePlanCode   string                      `json:"price_plan_code"`
	IsCustomPlan    bool                        `json:"is_custom_plan"`
	ItemID          string                      `json:"item_id"`
	ItemCode        string                      `json:"item_code,omitempty"`
	ItemName        string                      `json:"item_name,omitempty"`
	TariffClassID   string                      `json:"tariff_class_id"`
	TariffClassCode string                      `json:"tariff_class_code,omitempty"`
	TariffClassName string                      `json:"tariff_class_name,omitempty"`
	IsCito          bool                        `json:"is_cito"`
	TotalPrice      float64                     `json:"total_price"`
	Components      []ResolvedComponentResponse `json:"components"`
}
