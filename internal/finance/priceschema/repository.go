package priceschema

import (
	"context"
	"errors"
	"strings"

	"hosim-go/internal/database"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository membuat instance repository GORM untuk Price Schema.
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) getDB(ctx context.Context) *gorm.DB {
	return database.DBFromContextOr(ctx, r.db)
}

func (r *repository) CreateSchema(ctx context.Context, schema *PriceSchema) error {
	return r.getDB(ctx).Create(schema).Error
}

func (r *repository) UpdateSchema(ctx context.Context, schema *PriceSchema) error {
	return r.getDB(ctx).
		Model(schema).
		Select("Name", "Description", "BasePricePlanID", "IsActive", "UpdatedBy", "UpdatedAt").
		Updates(schema).Error
}

func (r *repository) FindSchemaByID(ctx context.Context, id string) (*PriceSchema, error) {
	var schema PriceSchema
	err := r.getDB(ctx).
		Model(&PriceSchema{}).
		Select("price_schemas.*, tariff_price_plans.code AS base_price_plan_code, tariff_price_plans.name AS base_price_plan_name").
		Joins("LEFT JOIN tariff_price_plans ON tariff_price_plans.id = price_schemas.base_price_plan_id AND tariff_price_plans.deleted_at IS NULL").
		Preload("Rules", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Where("price_schemas.id = ?", id).
		First(&schema).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchemaNotFound
		}
		return nil, err
	}
	return &schema, nil
}

func (r *repository) FindSchemaByCode(ctx context.Context, code string) (*PriceSchema, error) {
	var schema PriceSchema
	err := r.getDB(ctx).
		Model(&PriceSchema{}).
		Select("price_schemas.*, tariff_price_plans.code AS base_price_plan_code, tariff_price_plans.name AS base_price_plan_name").
		Joins("LEFT JOIN tariff_price_plans ON tariff_price_plans.id = price_schemas.base_price_plan_id AND tariff_price_plans.deleted_at IS NULL").
		Where("price_schemas.code = ?", code).
		First(&schema).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSchemaNotFound
		}
		return nil, err
	}
	return &schema, nil
}

func (r *repository) ListSchemas(ctx context.Context, params ListSchemaParams) ([]PriceSchema, int64, error) {
	db := r.getDB(ctx).Model(&PriceSchema{})

	if params.Search != "" {
		s := "%" + strings.ToLower(params.Search) + "%"
		db = db.Where("LOWER(price_schemas.code) LIKE ? OR LOWER(price_schemas.name) LIKE ?", s, s)
	}
	if params.IsActive != nil {
		db = db.Where("price_schemas.is_active = ?", *params.IsActive)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
	}
	offset := (params.Page - 1) * params.Limit

	var list []PriceSchema
	err := db.
		Select("price_schemas.*, tariff_price_plans.code AS base_price_plan_code, tariff_price_plans.name AS base_price_plan_name").
		Joins("LEFT JOIN tariff_price_plans ON tariff_price_plans.id = price_schemas.base_price_plan_id AND tariff_price_plans.deleted_at IS NULL").
		Order("price_schemas.created_at DESC").
		Offset(offset).
		Limit(params.Limit).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) CreateRule(ctx context.Context, rule *PriceSchemaRule) error {
	return r.getDB(ctx).Create(rule).Error
}

func (r *repository) UpdateRule(ctx context.Context, rule *PriceSchemaRule) error {
	return r.getDB(ctx).
		Model(rule).
		Select("TariffAdjustmentType", "TariffAdjustmentValue", "MedicationAdjustmentType", "MedicationAdjustmentValue",
			"PayByPatientPercent", "PayByCustomerPercent", "TaxPercent", "UpdatedBy", "UpdatedAt").
		Updates(rule).Error
}

func (r *repository) DeleteRule(ctx context.Context, id string) error {
	return r.getDB(ctx).Where("id = ?", id).Delete(&PriceSchemaRule{}).Error
}

func (r *repository) FindRuleByID(ctx context.Context, id string) (*PriceSchemaRule, error) {
	var rule PriceSchemaRule
	err := r.getDB(ctx).Where("id = ?", id).First(&rule).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRuleNotFound
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) FindRuleByScope(ctx context.Context, schemaID string, deptID, unitID, itemID *string) (*PriceSchemaRule, error) {
	db := r.getDB(ctx).Where("price_schema_id = ?", schemaID)

	if deptID != nil && *deptID != "" {
		db = db.Where("departement_id = ?", *deptID)
	} else {
		db = db.Where("departement_id IS NULL")
	}

	if unitID != nil && *unitID != "" {
		db = db.Where("service_unit_id = ?", *unitID)
	} else {
		db = db.Where("service_unit_id IS NULL")
	}

	if itemID != nil && *itemID != "" {
		db = db.Where("item_id = ?", *itemID)
	} else {
		db = db.Where("item_id IS NULL")
	}

	var rule PriceSchemaRule
	if err := db.First(&rule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rule, nil
}

func (r *repository) ListRulesBySchemaID(ctx context.Context, schemaID string, params ListRuleParams) ([]PriceSchemaRule, int64, error) {
	db := r.getDB(ctx).Model(&PriceSchemaRule{}).Where("price_schema_id = ?", schemaID)

	if params.DepartementID != "" {
		db = db.Where("departement_id = ?", params.DepartementID)
	}
	if params.ServiceUnitID != "" {
		db = db.Where("service_unit_id = ?", params.ServiceUnitID)
	}
	if params.ItemID != "" {
		db = db.Where("item_id = ?", params.ItemID)
	}

	if params.Level != "" {
		switch MatchLevel(params.Level) {
		case MatchLevelItemUnit:
			db = db.Where("item_id IS NOT NULL AND service_unit_id IS NOT NULL")
		case MatchLevelItemDept:
			db = db.Where("item_id IS NOT NULL AND departement_id IS NOT NULL")
		case MatchLevelItem:
			db = db.Where("item_id IS NOT NULL AND service_unit_id IS NULL AND departement_id IS NULL")
		case MatchLevelUnit:
			db = db.Where("item_id IS NULL AND service_unit_id IS NOT NULL")
		case MatchLevelDept:
			db = db.Where("item_id IS NULL AND departement_id IS NOT NULL")
		case MatchLevelSchemaDefault:
			db = db.Where("item_id IS NULL AND service_unit_id IS NULL AND departement_id IS NULL")
		}
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	} else if params.Limit > 100 {
		params.Limit = 100
	}
	offset := (params.Page - 1) * params.Limit

	var list []PriceSchemaRule
	err := db.Order("created_at ASC").
		Offset(offset).
		Limit(params.Limit).
		Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) FindBestMatchingRule(ctx context.Context, schemaID, itemID, serviceUnitID, departementID string) (*PriceSchemaRule, error) {
	conditions := []string{"(departement_id IS NULL AND service_unit_id IS NULL AND item_id IS NULL)"}
	args := []interface{}{}

	if itemID != "" && serviceUnitID != "" {
		conditions = append(conditions, "(item_id = ? AND service_unit_id = ?)")
		args = append(args, itemID, serviceUnitID)
	}
	if itemID != "" && departementID != "" {
		conditions = append(conditions, "(item_id = ? AND departement_id = ?)")
		args = append(args, itemID, departementID)
	}
	if itemID != "" {
		conditions = append(conditions, "(item_id = ? AND service_unit_id IS NULL AND departement_id IS NULL)")
		args = append(args, itemID)
	}
	if serviceUnitID != "" {
		conditions = append(conditions, "(service_unit_id = ? AND item_id IS NULL)")
		args = append(args, serviceUnitID)
	}
	if departementID != "" {
		conditions = append(conditions, "(departement_id = ? AND item_id IS NULL)")
		args = append(args, departementID)
	}

	orClause := strings.Join(conditions, " OR ")

	var candidates []PriceSchemaRule
	err := r.getDB(ctx).
		Where("price_schema_id = ?", schemaID).
		Where(orClause, args...).
		Find(&candidates).Error
	if err != nil {
		return nil, err
	}

	if len(candidates) == 0 {
		return nil, ErrNoMatchingRule
	}

	priorityOf := func(rule PriceSchemaRule) int {
		switch rule.Level() {
		case MatchLevelItemUnit:
			return 1
		case MatchLevelItemDept:
			return 2
		case MatchLevelItem:
			return 3
		case MatchLevelUnit:
			return 4
		case MatchLevelDept:
			return 5
		default:
			return 6
		}
	}

	best := candidates[0]
	bestPriority := priorityOf(best)
	for _, cand := range candidates[1:] {
		p := priorityOf(cand)
		if p < bestPriority {
			best = cand
			bestPriority = p
		}
	}

	return &best, nil
}
