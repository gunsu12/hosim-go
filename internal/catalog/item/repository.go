package item

import (
	"context"
	"errors"
	"strings"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrItemNotFound            = errors.New("item tidak ditemukan")
	ErrCodeAlreadyExists       = errors.New("kode item sudah digunakan")
	ErrCategoryNotFound        = errors.New("kategori item tidak ditemukan")
	ErrProductLineNotFound     = errors.New("lini produk item tidak ditemukan")
	ErrUnitNotFound            = errors.New("satuan item tidak ditemukan")
	ErrDuplicateUnitName            = errors.New("nama satuan sudah terdaftar untuk item ini")
	ErrBaseUnitAlreadyExists        = errors.New("item sudah memiliki base unit")
	ErrBaseUnitRequired             = errors.New("base unit wajib diisi")
	ErrBaseUnitCannotBeDeleted      = errors.New("base unit tidak boleh dihapus")
	ErrSubtypeMismatch              = errors.New("subtipe item tidak sesuai dengan item_type")
	ErrCategoryCodeAlreadyExists    = errors.New("kode kategori item sudah digunakan")
	ErrProductLineCodeAlreadyExists = errors.New("kode lini produk item sudah digunakan")
)

// ListParams mendefinisikan parameter pencarian katalog universal
type ListParams struct {
	Page          int
	Limit         int
	Search        string
	ItemType      *enums.ItemType
	CategoryID    *string
	ProductLineID *string
	IsActive      *bool
}

// CategoryListParams mendefinisikan parameter pencarian master kategori item
type CategoryListParams struct {
	Page     int
	Limit    int
	Search   string
	ItemType *enums.ItemType
	IsActive *bool
}

// ProductLineListParams mendefinisikan parameter pencarian master lini produk item
type ProductLineListParams struct {
	Page     int
	Limit    int
	Search   string
	IsActive *bool
}

// Repository adalah port persistence untuk entitas item dan seluruh subtipe CTI
type Repository interface {
	// Item Aggregate
	CreateItem(ctx context.Context, item *Item) error
	UpdateItem(ctx context.Context, item *Item) error
	FindItemByID(ctx context.Context, id string) (*Item, error)
	FindItemByCode(ctx context.Context, code string) (*Item, error)
	FindAllItems(ctx context.Context, params ListParams) ([]Item, int64, error)
	DeleteItem(ctx context.Context, id string, deletedBy string) error

	// Subtypes
	CreateMedication(ctx context.Context, med *ItemMedication) error
	UpdateMedication(ctx context.Context, med *ItemMedication) error
	FindMedicationByItemID(ctx context.Context, itemID string) (*ItemMedication, error)

	CreateGeneral(ctx context.Context, gen *ItemGeneral) error
	UpdateGeneral(ctx context.Context, gen *ItemGeneral) error
	FindGeneralByItemID(ctx context.Context, itemID string) (*ItemGeneral, error)

	CreateAsset(ctx context.Context, asset *ItemAsset) error
	UpdateAsset(ctx context.Context, asset *ItemAsset) error
	FindAssetByItemID(ctx context.Context, itemID string) (*ItemAsset, error)

	CreateTariff(ctx context.Context, tariff *ItemTariff) error
	UpdateTariff(ctx context.Context, tariff *ItemTariff) error
	FindTariffByItemID(ctx context.Context, itemID string) (*ItemTariff, error)

	// Units (Multi-UOM)
	CreateUnit(ctx context.Context, unit *ItemUnit) error
	UpdateUnit(ctx context.Context, unit *ItemUnit) error
	DeleteUnit(ctx context.Context, id string) error
	FindUnitByID(ctx context.Context, id string) (*ItemUnit, error)
	FindUnitsByItemID(ctx context.Context, itemID string) ([]ItemUnit, error)
	// Category Operations
	CreateCategory(ctx context.Context, cat *ItemCategory) error
	UpdateCategory(ctx context.Context, cat *ItemCategory) error
	DeleteCategory(ctx context.Context, id string, deletedBy string) error
	FindCategoryByID(ctx context.Context, id string) (*ItemCategory, error)
	FindCategoryByCode(ctx context.Context, code string) (*ItemCategory, error)
	FindAllCategories(ctx context.Context, params CategoryListParams) ([]ItemCategory, int64, error)
	CategoryExists(ctx context.Context, categoryID string) (bool, error)

	// Product Line Operations
	CreateProductLine(ctx context.Context, pl *ItemProductLine) error
	UpdateProductLine(ctx context.Context, pl *ItemProductLine) error
	DeleteProductLine(ctx context.Context, id string, deletedBy string) error
	FindProductLineByID(ctx context.Context, id string) (*ItemProductLine, error)
	FindProductLineByCode(ctx context.Context, code string) (*ItemProductLine, error)
	FindAllProductLines(ctx context.Context, params ProductLineListParams) ([]ItemProductLine, int64, error)
	ProductLineExists(ctx context.Context, productLineID string) (bool, error)

	// Transaction support
	WithTransaction(ctx context.Context, fn func(txRepo Repository) error) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository membuat instance baru Repository berbasis GORM
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) WithTransaction(ctx context.Context, fn func(txRepo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&gormRepository{db: tx})
	})
}

// ==========================================
// Item Base Operations
// ==========================================

func (r *gormRepository) CreateItem(ctx context.Context, item *Item) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Create(item).Error
}

func (r *gormRepository) UpdateItem(ctx context.Context, item *Item) error {
	return r.db.WithContext(ctx).Omit(clause.Associations).Save(item).Error
}

func (r *gormRepository) FindItemByID(ctx context.Context, id string) (*Item, error) {
	var item Item
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("ProductLine").
		Preload("Medication").
		Preload("General").
		Preload("Asset").
		Preload("Tariff").
		Preload("Units").
		Where("id = ?", id).
		First(&item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *gormRepository) FindItemByCode(ctx context.Context, code string) (*Item, error) {
	var item Item
	err := r.db.WithContext(ctx).Where("code = ?", strings.TrimSpace(code)).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *gormRepository) FindAllItems(ctx context.Context, params ListParams) ([]Item, int64, error) {
	var items []Item
	var total int64

	query := r.db.WithContext(ctx).Model(&Item{}).
		Preload("Category").
		Preload("ProductLine").
		Preload("Units")

	if params.Search != "" {
		s := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(generic_name) LIKE ?", s, s, s)
	}

	if params.ItemType != nil {
		query = query.Where("item_type = ?", *params.ItemType)
	}

	if params.CategoryID != nil && *params.CategoryID != "" {
		query = query.Where("category_id = ?", *params.CategoryID)
	}

	if params.ProductLineID != nil && *params.ProductLineID != "" {
		query = query.Where("product_line_id = ?", *params.ProductLineID)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *gormRepository) DeleteItem(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Model(&Item{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": deletedBy,
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

// ==========================================
// Subtype Operations
// ==========================================

func (r *gormRepository) CreateMedication(ctx context.Context, med *ItemMedication) error {
	return r.db.WithContext(ctx).Create(med).Error
}

func (r *gormRepository) UpdateMedication(ctx context.Context, med *ItemMedication) error {
	return r.db.WithContext(ctx).Save(med).Error
}

func (r *gormRepository) FindMedicationByItemID(ctx context.Context, itemID string) (*ItemMedication, error) {
	var med ItemMedication
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&med).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &med, nil
}

func (r *gormRepository) CreateGeneral(ctx context.Context, gen *ItemGeneral) error {
	return r.db.WithContext(ctx).Create(gen).Error
}

func (r *gormRepository) UpdateGeneral(ctx context.Context, gen *ItemGeneral) error {
	return r.db.WithContext(ctx).Save(gen).Error
}

func (r *gormRepository) FindGeneralByItemID(ctx context.Context, itemID string) (*ItemGeneral, error) {
	var gen ItemGeneral
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&gen).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &gen, nil
}

func (r *gormRepository) CreateAsset(ctx context.Context, asset *ItemAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *gormRepository) UpdateAsset(ctx context.Context, asset *ItemAsset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *gormRepository) FindAssetByItemID(ctx context.Context, itemID string) (*ItemAsset, error) {
	var asset ItemAsset
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &asset, nil
}

func (r *gormRepository) CreateTariff(ctx context.Context, tariff *ItemTariff) error {
	return r.db.WithContext(ctx).Create(tariff).Error
}

func (r *gormRepository) UpdateTariff(ctx context.Context, tariff *ItemTariff) error {
	return r.db.WithContext(ctx).Save(tariff).Error
}

func (r *gormRepository) FindTariffByItemID(ctx context.Context, itemID string) (*ItemTariff, error) {
	var tariff ItemTariff
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).First(&tariff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &tariff, nil
}

// ==========================================
// Unit Operations
// ==========================================

func (r *gormRepository) CreateUnit(ctx context.Context, unit *ItemUnit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *gormRepository) UpdateUnit(ctx context.Context, unit *ItemUnit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

func (r *gormRepository) DeleteUnit(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&ItemUnit{}).Error
}

func (r *gormRepository) FindUnitByID(ctx context.Context, id string) (*ItemUnit, error) {
	var unit ItemUnit
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnitNotFound
		}
		return nil, err
	}
	return &unit, nil
}

func (r *gormRepository) FindUnitsByItemID(ctx context.Context, itemID string) ([]ItemUnit, error) {
	var units []ItemUnit
	err := r.db.WithContext(ctx).Where("item_id = ?", itemID).Order("conversion_factor ASC").Find(&units).Error
	return units, err
}

func (r *gormRepository) FindBaseUnitByItemID(ctx context.Context, itemID string) (*ItemUnit, error) {
	var unit ItemUnit
	err := r.db.WithContext(ctx).Where("item_id = ? AND is_base_unit = true", itemID).First(&unit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &unit, nil
}

// ==========================================
// Category Operations
// ==========================================

func (r *gormRepository) CreateCategory(ctx context.Context, cat *ItemCategory) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

func (r *gormRepository) UpdateCategory(ctx context.Context, cat *ItemCategory) error {
	return r.db.WithContext(ctx).Save(cat).Error
}

func (r *gormRepository) DeleteCategory(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Model(&ItemCategory{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": deletedBy,
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *gormRepository) FindCategoryByID(ctx context.Context, id string) (*ItemCategory, error) {
	var cat ItemCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *gormRepository) FindCategoryByCode(ctx context.Context, code string) (*ItemCategory, error) {
	var cat ItemCategory
	err := r.db.WithContext(ctx).Where("code = ?", strings.TrimSpace(code)).First(&cat).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &cat, nil
}

func (r *gormRepository) FindAllCategories(ctx context.Context, params CategoryListParams) ([]ItemCategory, int64, error) {
	var categories []ItemCategory
	var total int64

	query := r.db.WithContext(ctx).Model(&ItemCategory{})

	if params.Search != "" {
		s := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", s, s)
	}

	if params.ItemType != nil {
		query = query.Where("item_type = ?", *params.ItemType)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *gormRepository) CategoryExists(ctx context.Context, categoryID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ItemCategory{}).Where("id = ? AND is_active = true", categoryID).Count(&count).Error
	return count > 0, err
}

// ==========================================
// Product Line Operations
// ==========================================

func (r *gormRepository) CreateProductLine(ctx context.Context, pl *ItemProductLine) error {
	return r.db.WithContext(ctx).Create(pl).Error
}

func (r *gormRepository) UpdateProductLine(ctx context.Context, pl *ItemProductLine) error {
	return r.db.WithContext(ctx).Save(pl).Error
}

func (r *gormRepository) DeleteProductLine(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Model(&ItemProductLine{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_by": deletedBy,
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *gormRepository) FindProductLineByID(ctx context.Context, id string) (*ItemProductLine, error) {
	var pl ItemProductLine
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&pl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductLineNotFound
		}
		return nil, err
	}
	return &pl, nil
}

func (r *gormRepository) FindProductLineByCode(ctx context.Context, code string) (*ItemProductLine, error) {
	var pl ItemProductLine
	err := r.db.WithContext(ctx).Where("code = ?", strings.TrimSpace(code)).First(&pl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductLineNotFound
		}
		return nil, err
	}
	return &pl, nil
}

func (r *gormRepository) FindAllProductLines(ctx context.Context, params ProductLineListParams) ([]ItemProductLine, int64, error) {
	var productLines []ItemProductLine
	var total int64

	query := r.db.WithContext(ctx).Model(&ItemProductLine{})

	if params.Search != "" {
		s := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ?", s, s)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	if err := query.Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&productLines).Error; err != nil {
		return nil, 0, err
	}

	return productLines, total, nil
}

func (r *gormRepository) ProductLineExists(ctx context.Context, productLineID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&ItemProductLine{}).Where("id = ? AND is_active = true", productLineID).Count(&count).Error
	return count > 0, err
}
