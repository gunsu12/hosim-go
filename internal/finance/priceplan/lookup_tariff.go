package priceplan

import (
	"context"
	"errors"
	"strings"
	"time"
)

// LookupTariffUseCase mengorkestrasi pencarian tarif dengan algoritma fallback bertingkat (PRD § 4.4 & § 6.4)
type LookupTariffUseCase interface {
	Execute(ctx context.Context, req LookupTariffRequest) (*LookupTariffResponse, error)
}

type lookupTariffUseCase struct {
	repo Repository
}

func NewLookupTariffUseCase(repo Repository) LookupTariffUseCase {
	return &lookupTariffUseCase{repo: repo}
}

func (uc *lookupTariffUseCase) Execute(ctx context.Context, req LookupTariffRequest) (*LookupTariffResponse, error) {
	req.TariffClassID = strings.TrimSpace(req.TariffClassID)
	req.ItemID = strings.TrimSpace(req.ItemID)

	if req.TariffClassID == "" || req.ItemID == "" {
		return nil, ErrTariffNotConfigured
	}

	// Parsing tanggal transaksi (mendukung RFC3339 atau format YYYY-MM-DD)
	txDate, err := parseDate(req.TransactionDate)
	if err != nil {
		return nil, errors.New("format tanggal transaksi tidak valid, gunakan RFC3339 atau YYYY-MM-DD")
	}

	var matchedPlan *TariffPricePlan
	var matchedItem *TariffPricePlanItem
	isCustomPlan := false

	// =========================================================================
	// LANGKAH 1: Cek Buku Tarif Khusus Penjamin / PKS (Custom Plan)
	// =========================================================================
	if req.CustomerID != nil && strings.TrimSpace(*req.CustomerID) != "" {
		custID := strings.TrimSpace(*req.CustomerID)
		customPlan, err := uc.repo.FindActiveCustomPlan(ctx, custID, txDate)
		if err == nil && customPlan != nil {
			// Cari apakah tindakan + kelas tarif terdaftar pada custom plan ini
			item, err := uc.repo.FindPlanItemForLookup(ctx, customPlan.ID, req.ItemID, req.TariffClassID)
			if err == nil && item != nil && item.IsActive {
				matchedPlan = customPlan
				matchedItem = item
				isCustomPlan = true
			}
		}
	}

	// =========================================================================
	// LANGKAH 2: Fallback ke Buku Tarif Standar RS (General/Standard Plan)
	// =========================================================================
	if matchedPlan == nil {
		defaultPlan, err := uc.repo.FindActiveDefaultPlan(ctx, txDate)
		if err == nil && defaultPlan != nil {
			item, err := uc.repo.FindPlanItemForLookup(ctx, defaultPlan.ID, req.ItemID, req.TariffClassID)
			if err == nil && item != nil && item.IsActive {
				matchedPlan = defaultPlan
				matchedItem = item
				isCustomPlan = false
			}
		}
	}

	// =========================================================================
	// LANGKAH 3: Penanganan Not Found Error
	// =========================================================================
	if matchedPlan == nil || matchedItem == nil {
		return nil, ErrTariffNotConfigured
	}

	// Kalkulasi harga final dan pecahan komponen (menangani kondisi reguler maupun CITO)
	finalTotalPrice, resolvedComponents := matchedItem.CalculatePrice(req.IsCito, matchedPlan.DefaultCitoPercent)

	respComponents := make([]ResolvedComponentResponse, len(resolvedComponents))
	for i, c := range resolvedComponents {
		respComponents[i] = ResolvedComponentResponse{
			ComponentID:   c.ComponentID,
			ComponentName: c.ComponentName,
			ComponentType: c.ComponentType,
			Amount:        c.Amount,
			COACode:       c.COACode,
		}
	}

	return &LookupTariffResponse{
		PricePlanID:     matchedPlan.ID,
		PricePlanCode:   matchedPlan.Code,
		IsCustomPlan:    isCustomPlan,
		ItemID:          matchedItem.ItemID,
		ItemCode:        matchedItem.ItemCode,
		ItemName:        matchedItem.ItemName,
		TariffClassID:   matchedItem.TariffClassID,
		TariffClassCode: matchedItem.TariffClassCode,
		TariffClassName: matchedItem.TariffClassName,
		IsCito:          req.IsCito,
		TotalPrice:      finalTotalPrice,
		Components:      respComponents,
	}, nil
}

// parseDate membantu parsing string tanggal ke time.Time
func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	// Coba format RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// Coba format YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	// Coba format datetime tanpa timezone (contoh: 2006-01-02 15:04:05)
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, nil
	}
	return time.Time{}, errors.New("invalid date format")
}
