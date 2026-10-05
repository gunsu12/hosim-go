package priceplan

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository adalah port persistence untuk buku tarif dan lookup engine
type Repository interface {
	// Plan Operations
	CreatePlan(ctx context.Context, plan *TariffPricePlan) error
	UpdatePlan(ctx context.Context, plan *TariffPricePlan) error
	FindPlanByID(ctx context.Context, id string) (*TariffPricePlan, error)
	FindPlanByCode(ctx context.Context, code string) (*TariffPricePlan, error)
	FindAllPlans(ctx context.Context, params PlanListParams) ([]TariffPricePlan, int64, error)
	CheckDefaultConflict(ctx context.Context, planID string, effectiveFrom time.Time, effectiveTo *time.Time) (bool, error)

	// Item Operations
	AddItem(ctx context.Context, item *TariffPricePlanItem) error
	UpdateItem(ctx context.Context, item *TariffPricePlanItem) error
	DeleteItem(ctx context.Context, planID, itemID string, deletedBy string) error
	FindItemByID(ctx context.Context, planID, itemID string) (*TariffPricePlanItem, error)
	FindItemsByPlanID(ctx context.Context, planID string, params ItemListParams) ([]TariffPricePlanItem, int64, error)
	FindAllItemsByPlanID(ctx context.Context, planID string) ([]TariffPricePlanItem, error)
	BatchUpsertItems(ctx context.Context, planID string, items []TariffPricePlanItem, operator string) error

	// Lookup Engine Operations (PRD § 4.4 & § 6.4)
	FindActiveCustomPlan(ctx context.Context, customerID string, txDate time.Time) (*TariffPricePlan, error)
	FindActiveDefaultPlan(ctx context.Context, txDate time.Time) (*TariffPricePlan, error)
	FindPlanItemForLookup(ctx context.Context, planID, itemID, tariffClassID string) (*TariffPricePlanItem, error)

	// Clone Operation
	ClonePlan(ctx context.Context, sourcePlanID string, newPlan *TariffPricePlan) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreatePlan(ctx context.Context, plan *TariffPricePlan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *repository) UpdatePlan(ctx context.Context, plan *TariffPricePlan) error {
	return r.db.WithContext(ctx).
		Model(plan).
		Select("Name", "Description", "EffectiveFrom", "EffectiveTo", "Status", "IsDefault", "CustomerID", "DefaultCitoPercent", "ApprovedAt", "ApprovedBy", "UpdatedBy", "UpdatedAt").
		Updates(plan).Error
}

func (r *repository) FindPlanByID(ctx context.Context, id string) (*TariffPricePlan, error) {
	var plan TariffPricePlan
	query := r.db.WithContext(ctx).
		Table("tariff_price_plans").
		Select("tariff_price_plans.*, customers.name AS customer_name, (SELECT COUNT(1) FROM tariff_price_plan_items WHERE tariff_price_plan_items.price_plan_id = tariff_price_plans.id AND tariff_price_plan_items.deleted_at IS NULL) AS total_items").
		Joins("LEFT JOIN customers ON customers.id = tariff_price_plans.customer_id").
		Where("tariff_price_plans.id = ? AND tariff_price_plans.deleted_at IS NULL", id)

	if err := query.First(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *repository) FindPlanByCode(ctx context.Context, code string) (*TariffPricePlan, error) {
	var plan TariffPricePlan
	err := r.db.WithContext(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *repository) FindAllPlans(ctx context.Context, params PlanListParams) ([]TariffPricePlan, int64, error) {
	var plans []TariffPricePlan
	var total int64

	query := r.db.WithContext(ctx).
		Table("tariff_price_plans").
		Select("tariff_price_plans.*, customers.name AS customer_name, (SELECT COUNT(1) FROM tariff_price_plan_items WHERE tariff_price_plan_items.price_plan_id = tariff_price_plans.id AND tariff_price_plan_items.deleted_at IS NULL) AS total_items").
		Joins("LEFT JOIN customers ON customers.id = tariff_price_plans.customer_id").
		Where("tariff_price_plans.deleted_at IS NULL")

	if params.Search != "" {
		s := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("(LOWER(tariff_price_plans.code) LIKE ? OR LOWER(tariff_price_plans.name) LIKE ? OR LOWER(tariff_price_plans.description) LIKE ?)", s, s, s)
	}
	if params.Status != nil && *params.Status != "" {
		query = query.Where("tariff_price_plans.status = ?", *params.Status)
	}
	if params.CustomerID != nil && *params.CustomerID != "" {
		query = query.Where("tariff_price_plans.customer_id = ?", *params.CustomerID)
	}
	if params.IsDefault != nil {
		query = query.Where("tariff_price_plans.is_default = ?", *params.IsDefault)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.Limit

	err := query.Order("tariff_price_plans.effective_from DESC, tariff_price_plans.created_at DESC").
		Limit(params.Limit).
		Offset(offset).
		Find(&plans).Error

	return plans, total, err
}

func (r *repository) CheckDefaultConflict(ctx context.Context, planID string, effectiveFrom time.Time, effectiveTo *time.Time) (bool, error) {
	// Memeriksa apakah ada buku tarif default lain yang berstatus ACTIVE dengan overlap tanggal
	query := r.db.WithContext(ctx).
		Table("tariff_price_plans").
		Where("is_default = true AND status = ? AND deleted_at IS NULL", PricePlanStatusActive)

	if planID != "" {
		query = query.Where("id != ?", planID)
	}

	// Overlap check logic:
	// A.from <= B.to AND (A.to IS NULL OR A.to >= B.from)
	if effectiveTo == nil {
		query = query.Where("(effective_to IS NULL OR effective_to >= ?)", effectiveFrom)
	} else {
		query = query.Where("effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)", *effectiveTo, effectiveFrom)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) AddItem(ctx context.Context, item *TariffPricePlanItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Validasi duplikasi (item_id, tariff_class_id) dalam price_plan
		var count int64
		err := tx.Table("tariff_price_plan_items").
			Where("price_plan_id = ? AND item_id = ? AND tariff_class_id = ? AND deleted_at IS NULL",
				item.PricePlanID, item.ItemID, item.TariffClassID).
			Count(&count).Error
		if err != nil {
			return err
		}
		if count > 0 {
			return ErrItemAlreadyInPlan
		}

		if err := tx.Omit(clause.Associations).Create(item).Error; err != nil {
			return err
		}

		for idx := range item.Components {
			item.Components[idx].PlanItemID = item.ID
			if err := tx.Create(&item.Components[idx]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) UpdateItem(ctx context.Context, item *TariffPricePlanItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Save(item).Error; err != nil {
			return err
		}

		// Gantikan seluruh rincian komponen
		if err := tx.Where("plan_item_id = ?", item.ID).Delete(&TariffPricePlanItemComponent{}).Error; err != nil {
			return err
		}

		for idx := range item.Components {
			item.Components[idx].ID = "" // generate uuid baru via BeforeCreate
			item.Components[idx].PlanItemID = item.ID
			if err := tx.Create(&item.Components[idx]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) DeleteItem(ctx context.Context, planID, itemID string, deletedBy string) error {
	return r.db.WithContext(ctx).
		Table("tariff_price_plan_items").
		Where("id = ? AND price_plan_id = ? AND deleted_at IS NULL", itemID, planID).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		}).Error
}

func (r *repository) FindItemByID(ctx context.Context, planID, itemID string) (*TariffPricePlanItem, error) {
	var item TariffPricePlanItem
	err := r.db.WithContext(ctx).
		Table("tariff_price_plan_items").
		Select("tariff_price_plan_items.*, items.name AS item_name, items.code AS item_code, tariff_classes.name AS tariff_class_name, tariff_classes.code AS tariff_class_code").
		Joins("JOIN items ON items.id = tariff_price_plan_items.item_id").
		Joins("JOIN tariff_classes ON tariff_classes.id = tariff_price_plan_items.tariff_class_id").
		Where("tariff_price_plan_items.id = ? AND tariff_price_plan_items.price_plan_id = ? AND tariff_price_plan_items.deleted_at IS NULL", itemID, planID).
		First(&item).Error
	if err != nil {
		return nil, err
	}

	// Load komponen beserta informasi nama master komponen
	var components []TariffPricePlanItemComponent
	err = r.db.WithContext(ctx).
		Table("tariff_price_plan_item_components").
		Select("tariff_price_plan_item_components.*, tariff_components.code AS component_code, tariff_components.name AS component_name, tariff_components.component_type AS component_type, tariff_components.default_coa_code AS default_coa_code").
		Joins("JOIN tariff_components ON tariff_components.id = tariff_price_plan_item_components.component_id").
		Where("tariff_price_plan_item_components.plan_item_id = ?", item.ID).
		Find(&components).Error
	if err != nil {
		return nil, err
	}

	item.Components = components
	return &item, nil
}

func (r *repository) FindItemsByPlanID(ctx context.Context, planID string, params ItemListParams) ([]TariffPricePlanItem, int64, error) {
	var items []TariffPricePlanItem
	var total int64

	query := r.db.WithContext(ctx).
		Table("tariff_price_plan_items").
		Select("tariff_price_plan_items.*, items.name AS item_name, items.code AS item_code, tariff_classes.name AS tariff_class_name, tariff_classes.code AS tariff_class_code").
		Joins("JOIN items ON items.id = tariff_price_plan_items.item_id").
		Joins("JOIN tariff_classes ON tariff_classes.id = tariff_price_plan_items.tariff_class_id").
		Where("tariff_price_plan_items.price_plan_id = ? AND tariff_price_plan_items.deleted_at IS NULL", planID)

	if params.Search != "" {
		s := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("(LOWER(items.name) LIKE ? OR LOWER(items.code) LIKE ?)", s, s)
	}
	if params.TariffClassID != nil && *params.TariffClassID != "" {
		query = query.Where("tariff_price_plan_items.tariff_class_id = ?", *params.TariffClassID)
	}
	if params.IsActive != nil {
		query = query.Where("tariff_price_plan_items.is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.Limit

	err := query.Order("items.name ASC, tariff_classes.code ASC").
		Limit(params.Limit).
		Offset(offset).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// Preload components untuk items yang diambil
	if len(items) > 0 {
		itemIDs := make([]string, len(items))
		for i, it := range items {
			itemIDs[i] = it.ID
		}

		var components []TariffPricePlanItemComponent
		err = r.db.WithContext(ctx).
			Table("tariff_price_plan_item_components").
			Select("tariff_price_plan_item_components.*, tariff_components.code AS component_code, tariff_components.name AS component_name, tariff_components.component_type AS component_type, tariff_components.default_coa_code AS default_coa_code").
			Joins("JOIN tariff_components ON tariff_components.id = tariff_price_plan_item_components.component_id").
			Where("tariff_price_plan_item_components.plan_item_id IN ?", itemIDs).
			Find(&components).Error
		if err != nil {
			return nil, 0, err
		}

		compMap := make(map[string][]TariffPricePlanItemComponent)
		for _, c := range components {
			compMap[c.PlanItemID] = append(compMap[c.PlanItemID], c)
		}

		for i := range items {
			items[i].Components = compMap[items[i].ID]
		}
	}

	return items, total, nil
}

func (r *repository) FindAllItemsByPlanID(ctx context.Context, planID string) ([]TariffPricePlanItem, error) {
	var items []TariffPricePlanItem
	err := r.db.WithContext(ctx).
		Table("tariff_price_plan_items").
		Where("price_plan_id = ? AND deleted_at IS NULL", planID).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return items, nil
	}

	itemIDs := make([]string, len(items))
	for i, it := range items {
		itemIDs[i] = it.ID
	}

	var components []TariffPricePlanItemComponent
	err = r.db.WithContext(ctx).
		Table("tariff_price_plan_item_components").
		Where("plan_item_id IN ?", itemIDs).
		Find(&components).Error
	if err != nil {
		return nil, err
	}

	compMap := make(map[string][]TariffPricePlanItemComponent)
	for _, c := range components {
		compMap[c.PlanItemID] = append(compMap[c.PlanItemID], c)
	}

	for i := range items {
		items[i].Components = compMap[items[i].ID]
	}

	return items, nil
}

func (r *repository) BatchUpsertItems(ctx context.Context, planID string, items []TariffPricePlanItem, operator string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			it.PricePlanID = planID
			it.CreatedBy = operator
			it.UpdatedBy = operator

			// Upsert header item
			var existing TariffPricePlanItem
			err := tx.Where("price_plan_id = ? AND item_id = ? AND tariff_class_id = ? AND deleted_at IS NULL",
				planID, it.ItemID, it.TariffClassID).First(&existing).Error

			if err == nil {
				// Update
				existing.TotalBasePrice = it.TotalBasePrice
				existing.TotalCitoPrice = it.TotalCitoPrice
				existing.IsActive = it.IsActive
				existing.UpdatedBy = operator
				if err := tx.Omit(clause.Associations).Save(&existing).Error; err != nil {
					return err
				}
				it.ID = existing.ID
				// Bersihkan komponen lama
				if err := tx.Where("plan_item_id = ?", it.ID).Delete(&TariffPricePlanItemComponent{}).Error; err != nil {
					return err
				}
			} else if errors.Is(err, gorm.ErrRecordNotFound) {
				// Insert baru
				if err := tx.Omit(clause.Associations).Create(&it).Error; err != nil {
					return err
				}
			} else {
				return err
			}

			// Simpan rincian komponen
			for _, comp := range it.Components {
				comp.ID = ""
				comp.PlanItemID = it.ID
				comp.CreatedBy = operator
				comp.UpdatedBy = operator
				if err := tx.Create(&comp).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// -------------------------------------------------------------
// Lookup Engine Repository Methods (PRD § 4.4 & § 6.4)
// -------------------------------------------------------------

// FindActiveCustomPlan mencari buku tarif aktif khusus penjamin (PKS) yang berlaku pada txDate
func (r *repository) FindActiveCustomPlan(ctx context.Context, customerID string, txDate time.Time) (*TariffPricePlan, error) {
	var plan TariffPricePlan
	dateStr := txDate.Format("2006-01-02")

	err := r.db.WithContext(ctx).
		Where("customer_id = ? AND status = ? AND deleted_at IS NULL", customerID, PricePlanStatusActive).
		Where("effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)", dateStr, dateStr).
		Order("effective_from DESC, created_at DESC").
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindActiveDefaultPlan mencari buku tarif standar RS (is_default = true) yang berlaku pada txDate
func (r *repository) FindActiveDefaultPlan(ctx context.Context, txDate time.Time) (*TariffPricePlan, error) {
	var plan TariffPricePlan
	dateStr := txDate.Format("2006-01-02")

	err := r.db.WithContext(ctx).
		Where("is_default = true AND status = ? AND deleted_at IS NULL", PricePlanStatusActive).
		Where("effective_from <= ? AND (effective_to IS NULL OR effective_to >= ?)", dateStr, dateStr).
		Order("effective_from DESC, created_at DESC").
		First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindPlanItemForLookup mengambil baris item beserta seluruh komponen biaya dan metadata item/komponen untuk lookup
func (r *repository) FindPlanItemForLookup(ctx context.Context, planID, itemID, tariffClassID string) (*TariffPricePlanItem, error) {
	var item TariffPricePlanItem
	err := r.db.WithContext(ctx).
		Table("tariff_price_plan_items").
		Select("tariff_price_plan_items.*, items.name AS item_name, items.code AS item_code, tariff_classes.name AS tariff_class_name, tariff_classes.code AS tariff_class_code").
		Joins("JOIN items ON items.id = tariff_price_plan_items.item_id").
		Joins("JOIN tariff_classes ON tariff_classes.id = tariff_price_plan_items.tariff_class_id").
		Where("tariff_price_plan_items.price_plan_id = ? AND tariff_price_plan_items.item_id = ? AND tariff_price_plan_items.tariff_class_id = ? AND tariff_price_plan_items.is_active = true AND tariff_price_plan_items.deleted_at IS NULL",
			planID, itemID, tariffClassID).
		First(&item).Error
	if err != nil {
		return nil, err
	}

	var components []TariffPricePlanItemComponent
	err = r.db.WithContext(ctx).
		Table("tariff_price_plan_item_components").
		Select("tariff_price_plan_item_components.*, tariff_components.code AS component_code, tariff_components.name AS component_name, tariff_components.component_type AS component_type, tariff_components.default_coa_code AS default_coa_code").
		Joins("JOIN tariff_components ON tariff_components.id = tariff_price_plan_item_components.component_id").
		Where("tariff_price_plan_item_components.plan_item_id = ?", item.ID).
		Find(&components).Error
	if err != nil {
		return nil, err
	}

	item.Components = components
	return &item, nil
}

// ClonePlan menyalin seluruh isi buku tarif acuan ke buku tarif baru berstatus DRAFT
func (r *repository) ClonePlan(ctx context.Context, sourcePlanID string, newPlan *TariffPricePlan) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Simpan header buku tarif baru
		if err := tx.Create(newPlan).Error; err != nil {
			return err
		}

		// 2. Ambil seluruh item aktif dari buku tarif sumber
		var sourceItems []TariffPricePlanItem
		err := tx.Table("tariff_price_plan_items").
			Where("price_plan_id = ? AND deleted_at IS NULL", sourcePlanID).
			Find(&sourceItems).Error
		if err != nil {
			return err
		}

		if len(sourceItems) == 0 {
			return nil
		}

		sourceItemIDs := make([]string, len(sourceItems))
		for i, it := range sourceItems {
			sourceItemIDs[i] = it.ID
		}

		var sourceComponents []TariffPricePlanItemComponent
		err = tx.Table("tariff_price_plan_item_components").
			Where("plan_item_id IN ?", sourceItemIDs).
			Find(&sourceComponents).Error
		if err != nil {
			return err
		}

		compMap := make(map[string][]TariffPricePlanItemComponent)
		for _, c := range sourceComponents {
			compMap[c.PlanItemID] = append(compMap[c.PlanItemID], c)
		}

		// 3. Gandakan baris item dan komponen ke newPlan
		for _, sItem := range sourceItems {
			newItem := TariffPricePlanItem{
				PricePlanID:    newPlan.ID,
				ItemID:         sItem.ItemID,
				TariffClassID:  sItem.TariffClassID,
				TotalBasePrice: sItem.TotalBasePrice,
				TotalCitoPrice: sItem.TotalCitoPrice,
				IsActive:       sItem.IsActive,
				CreatedBy:      newPlan.CreatedBy,
				UpdatedBy:      newPlan.CreatedBy,
			}
			if err := tx.Omit(clause.Associations).Create(&newItem).Error; err != nil {
				return err
			}

			sComps := compMap[sItem.ID]
			newComps := make([]TariffPricePlanItemComponent, len(sComps))
			for cIdx, sc := range sComps {
				newComps[cIdx] = TariffPricePlanItemComponent{
					PlanItemID:  newItem.ID,
					ComponentID: sc.ComponentID,
					BaseAmount:  sc.BaseAmount,
					CitoAmount:  sc.CitoAmount,
					COACode:     sc.COACode,
					CreatedBy:   newPlan.CreatedBy,
					UpdatedBy:   newPlan.CreatedBy,
				}
			}
			if len(newComps) > 0 {
				if err := tx.CreateInBatches(newComps, 100).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
