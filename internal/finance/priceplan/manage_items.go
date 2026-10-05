package priceplan

import (
	"context"
	"errors"
	"math"
	"strings"

	"gorm.io/gorm"
)

// ManageItemsUseCase mengorkestrasi pengelolaan item tarif dan pecahan komponen (hanya saat status DRAFT)
type ManageItemsUseCase interface {
	AddItem(ctx context.Context, planID string, req AddItemRequest, operator string) (*PricePlanItemResponse, error)
	UpdateItem(ctx context.Context, planID string, itemID string, req UpdateItemRequest, operator string) (*PricePlanItemResponse, error)
	DeleteItem(ctx context.Context, planID string, itemID string, operator string) error
	GetItem(ctx context.Context, planID string, itemID string) (*PricePlanItemResponse, error)
	ListItems(ctx context.Context, planID string, params ItemListParams) ([]PricePlanItemResponse, int64, error)
	BatchUpsert(ctx context.Context, planID string, req BatchUpsertItemsRequest, operator string) error
}

type manageItemsUseCase struct {
	repo Repository
}

func NewManageItemsUseCase(repo Repository) ManageItemsUseCase {
	return &manageItemsUseCase{repo: repo}
}

func (uc *manageItemsUseCase) AddItem(ctx context.Context, planID string, req AddItemRequest, operator string) (*PricePlanItemResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if !plan.CanModify() {
		return nil, ErrPricePlanImmutable
	}

	req.ItemID = strings.TrimSpace(req.ItemID)
	req.TariffClassID = strings.TrimSpace(req.TariffClassID)
	if req.ItemID == "" || req.TariffClassID == "" {
		return nil, errors.New("item_id dan tariff_class_id wajib diisi")
	}

	if len(req.Components) == 0 {
		return nil, errors.New("minimal satu rincian komponen biaya wajib disertakan")
	}

	// Hitung dan validasi balancing komponen
	var sum float64
	seenComp := make(map[string]bool, len(req.Components))
	components := make([]TariffPricePlanItemComponent, len(req.Components))
	for i, c := range req.Components {
		cID := strings.TrimSpace(c.ComponentID)
		if cID == "" {
			return nil, errors.New("component_id pada rincian komponen wajib diisi")
		}
		if seenComp[cID] {
			return nil, ErrDuplicateComponent
		}
		seenComp[cID] = true
		sum += c.BaseAmount
		components[i] = TariffPricePlanItemComponent{
			ComponentID: cID,
			BaseAmount:  c.BaseAmount,
			CitoAmount:  c.CitoAmount,
			COACode:     c.COACode,
			CreatedBy:   operator,
			UpdatedBy:   operator,
		}
	}

	if math.Abs(sum-req.TotalBasePrice) > 0.01 {
		return nil, ErrTariffComponentsMismatch
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item := &TariffPricePlanItem{
		PricePlanID:    planID,
		ItemID:         req.ItemID,
		TariffClassID:  req.TariffClassID,
		TotalBasePrice: req.TotalBasePrice,
		TotalCitoPrice: req.TotalCitoPrice,
		IsActive:       isActive,
		CreatedBy:      operator,
		UpdatedBy:      operator,
		Components:     components,
	}

	if err := uc.repo.AddItem(ctx, item); err != nil {
		return nil, err
	}

	// Reload item dengan informasi nama relasi
	savedItem, err := uc.repo.FindItemByID(ctx, planID, item.ID)
	if err != nil {
		return nil, err
	}

	return toItemResponse(savedItem), nil
}

func (uc *manageItemsUseCase) UpdateItem(ctx context.Context, planID string, itemID string, req UpdateItemRequest, operator string) (*PricePlanItemResponse, error) {
	plan, err := uc.repo.FindPlanByID(ctx, planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanNotFound
		}
		return nil, err
	}

	if !plan.CanModify() {
		return nil, ErrPricePlanImmutable
	}

	existingItem, err := uc.repo.FindItemByID(ctx, planID, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanItemNotFound
		}
		return nil, err
	}

	if len(req.Components) == 0 {
		return nil, errors.New("minimal satu rincian komponen biaya wajib disertakan")
	}

	var sum float64
	seenComp := make(map[string]bool, len(req.Components))
	components := make([]TariffPricePlanItemComponent, len(req.Components))
	for i, c := range req.Components {
		cID := strings.TrimSpace(c.ComponentID)
		if cID == "" {
			return nil, errors.New("component_id pada rincian komponen wajib diisi")
		}
		if seenComp[cID] {
			return nil, ErrDuplicateComponent
		}
		seenComp[cID] = true
		sum += c.BaseAmount
		components[i] = TariffPricePlanItemComponent{
			ComponentID: cID,
			BaseAmount:  c.BaseAmount,
			CitoAmount:  c.CitoAmount,
			COACode:     c.COACode,
			CreatedBy:   operator,
			UpdatedBy:   operator,
		}
	}

	if math.Abs(sum-req.TotalBasePrice) > 0.01 {
		return nil, ErrTariffComponentsMismatch
	}

	existingItem.TotalBasePrice = req.TotalBasePrice
	existingItem.TotalCitoPrice = req.TotalCitoPrice
	if req.IsActive != nil {
		existingItem.IsActive = *req.IsActive
	}
	existingItem.UpdatedBy = operator
	existingItem.Components = components

	if err := uc.repo.UpdateItem(ctx, existingItem); err != nil {
		return nil, err
	}

	savedItem, err := uc.repo.FindItemByID(ctx, planID, itemID)
	if err != nil {
		return nil, err
	}

	return toItemResponse(savedItem), nil
}

func (uc *manageItemsUseCase) DeleteItem(ctx context.Context, planID string, itemID string, operator string) error {
	plan, err := uc.repo.FindPlanByID(ctx, planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPricePlanNotFound
		}
		return err
	}

	if !plan.CanModify() {
		return ErrPricePlanImmutable
	}

	_, err = uc.repo.FindItemByID(ctx, planID, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPricePlanItemNotFound
		}
		return err
	}

	return uc.repo.DeleteItem(ctx, planID, itemID, operator)
}

func (uc *manageItemsUseCase) GetItem(ctx context.Context, planID string, itemID string) (*PricePlanItemResponse, error) {
	item, err := uc.repo.FindItemByID(ctx, planID, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPricePlanItemNotFound
		}
		return nil, err
	}

	return toItemResponse(item), nil
}

func (uc *manageItemsUseCase) ListItems(ctx context.Context, planID string, params ItemListParams) ([]PricePlanItemResponse, int64, error) {
	items, total, err := uc.repo.FindItemsByPlanID(ctx, planID, params)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]PricePlanItemResponse, len(items))
	for i := range items {
		responses[i] = *toItemResponse(&items[i])
	}

	return responses, total, nil
}

func (uc *manageItemsUseCase) BatchUpsert(ctx context.Context, planID string, req BatchUpsertItemsRequest, operator string) error {
	plan, err := uc.repo.FindPlanByID(ctx, planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPricePlanNotFound
		}
		return err
	}

	if !plan.CanModify() {
		return ErrPricePlanImmutable
	}

	items := make([]TariffPricePlanItem, len(req.Items))
	for idx, itReq := range req.Items {
		itReq.ItemID = strings.TrimSpace(itReq.ItemID)
		itReq.TariffClassID = strings.TrimSpace(itReq.TariffClassID)
		if itReq.ItemID == "" || itReq.TariffClassID == "" {
			return errors.New("item_id dan tariff_class_id wajib diisi pada setiap item")
		}

		var sum float64
		seenComp := make(map[string]bool, len(itReq.Components))
		components := make([]TariffPricePlanItemComponent, len(itReq.Components))
		for cIdx, c := range itReq.Components {
			cID := strings.TrimSpace(c.ComponentID)
			if cID == "" {
				return errors.New("component_id wajib diisi pada setiap rincian komponen")
			}
			if seenComp[cID] {
				return ErrDuplicateComponent
			}
			seenComp[cID] = true
			sum += c.BaseAmount
			components[cIdx] = TariffPricePlanItemComponent{
				ComponentID: cID,
				BaseAmount:  c.BaseAmount,
				CitoAmount:  c.CitoAmount,
				COACode:     c.COACode,
			}
		}

		if math.Abs(sum-itReq.TotalBasePrice) > 0.01 {
			return ErrTariffComponentsMismatch
		}

		isActive := true
		if itReq.IsActive != nil {
			isActive = *itReq.IsActive
		}

		items[idx] = TariffPricePlanItem{
			PricePlanID:    planID,
			ItemID:         itReq.ItemID,
			TariffClassID:  itReq.TariffClassID,
			TotalBasePrice: itReq.TotalBasePrice,
			TotalCitoPrice: itReq.TotalCitoPrice,
			IsActive:       isActive,
			Components:     components,
		}
	}

	return uc.repo.BatchUpsertItems(ctx, planID, items, operator)
}

func toItemResponse(i *TariffPricePlanItem) *PricePlanItemResponse {
	comps := make([]PricePlanItemComponentResponse, len(i.Components))
	for idx, c := range i.Components {
		comps[idx] = PricePlanItemComponentResponse{
			ID:            c.ID,
			PlanItemID:    c.PlanItemID,
			ComponentID:   c.ComponentID,
			ComponentCode: c.ComponentCode,
			ComponentName: c.ComponentName,
			ComponentType: c.ComponentType,
			BaseAmount:    c.BaseAmount,
			CitoAmount:    c.CitoAmount,
			COACode:       c.COACode,
		}
	}

	return &PricePlanItemResponse{
		ID:              i.ID,
		PricePlanID:     i.PricePlanID,
		ItemID:          i.ItemID,
		ItemCode:        i.ItemCode,
		ItemName:        i.ItemName,
		TariffClassID:   i.TariffClassID,
		TariffClassCode: i.TariffClassCode,
		TariffClassName: i.TariffClassName,
		TotalBasePrice:  i.TotalBasePrice,
		TotalCitoPrice:  i.TotalCitoPrice,
		IsActive:        i.IsActive,
		Components:      comps,
	}
}
