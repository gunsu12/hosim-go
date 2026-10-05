package account

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Repository adalah port database persistence untuk bagan akun (Chart of Accounts)
type Repository interface {
	Create(ctx context.Context, acc *ChartOfAccount) error
	CreateWithParentUpdate(ctx context.Context, acc *ChartOfAccount, parentID string) error
	Update(ctx context.Context, acc *ChartOfAccount) error
	Delete(ctx context.Context, id, deletedBy string) error
	DeleteWithParentCheck(ctx context.Context, id string, parentID *string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*ChartOfAccount, error)
	FindByCode(ctx context.Context, code string) (*ChartOfAccount, error)
	FindAll(ctx context.Context, params AccountListParams) ([]ChartOfAccount, int64, error)
	FindAllRaw(ctx context.Context, activeOnly bool) ([]ChartOfAccount, error)
	FindChildrenCount(ctx context.Context, parentID string) (int64, error)
	IsCodeUsedInTariff(ctx context.Context, code string) (bool, error)
	SetPostable(ctx context.Context, id string, isPostable bool) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, acc *ChartOfAccount) error {
	return r.db.WithContext(ctx).Select("*").Create(acc).Error
}

func (r *repository) CreateWithParentUpdate(ctx context.Context, acc *ChartOfAccount, parentID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Simpan akun anak baru
		if err := tx.Select("*").Create(acc).Error; err != nil {
			return err
		}

		// 2. Ubah akun induk menjadi non-postable (is_postable = false)
		if err := tx.Table("chart_of_accounts").
			Where("id = ? AND deleted_at IS NULL", parentID).
			Updates(map[string]any{
				"is_postable": false,
				"updated_at":  time.Now(),
			}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) Update(ctx context.Context, acc *ChartOfAccount) error {
	return r.db.WithContext(ctx).
		Model(acc).
		Select("Name", "Description", "Position", "IsTreasuryAccount", "BankName", "BankAccountNumber", "Currency", "IsActive", "UpdatedBy", "UpdatedAt").
		Updates(acc).Error
}

func (r *repository) Delete(ctx context.Context, id, deletedBy string) error {
	return r.db.WithContext(ctx).
		Table("chart_of_accounts").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"deleted_at": time.Now(),
			"deleted_by": deletedBy,
		}).Error
}

func (r *repository) DeleteWithParentCheck(ctx context.Context, id string, parentID *string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Soft delete akun terkait
		if err := tx.Table("chart_of_accounts").
			Where("id = ? AND deleted_at IS NULL", id).
			Updates(map[string]any{
				"deleted_at": time.Now(),
				"deleted_by": deletedBy,
			}).Error; err != nil {
			return err
		}

		// 2. Jika akun memiliki induk, periksa apakah induk masih memiliki anak aktif lain
		if parentID != nil && *parentID != "" {
			var count int64
			err := tx.Table("chart_of_accounts").
				Where("parent_id = ? AND deleted_at IS NULL", *parentID).
				Count(&count).Error
			if err != nil {
				return err
			}

			// Jika tidak ada lagi anak aktif, kembalikan induk menjadi postable (is_postable = true)
			if count == 0 {
				if err := tx.Table("chart_of_accounts").
					Where("id = ? AND deleted_at IS NULL", *parentID).
					Updates(map[string]any{
						"is_postable": true,
						"updated_at":  time.Now(),
					}).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*ChartOfAccount, error) {
	var acc ChartOfAccount
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&acc).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*ChartOfAccount, error) {
	var acc ChartOfAccount
	err := r.db.WithContext(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&acc).Error
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *repository) FindAll(ctx context.Context, params AccountListParams) ([]ChartOfAccount, int64, error) {
	var accounts []ChartOfAccount
	var total int64

	query := r.db.WithContext(ctx).
		Table("chart_of_accounts").
		Where("deleted_at IS NULL")

	if params.Search != "" {
		s := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("(LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", s, s, s)
	}
	if params.Type != nil && *params.Type != "" {
		query = query.Where("type = ?", *params.Type)
	}
	if params.Position != nil && *params.Position != "" {
		query = query.Where("position = ?", *params.Position)
	}
	if params.ParentID != nil {
		if *params.ParentID == "" {
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", *params.ParentID)
		}
	}
	if params.IsPostable != nil {
		query = query.Where("is_postable = ?", *params.IsPostable)
	}
	if params.IsTreasury != nil {
		query = query.Where("is_treasury_account = ?", *params.IsTreasury)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.Level != nil {
		query = query.Where("account_level = ?", *params.Level)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	offset := (params.Page - 1) * params.Limit

	err := query.Order("code ASC").
		Limit(params.Limit).
		Offset(offset).
		Find(&accounts).Error

	return accounts, total, err
}

func (r *repository) FindAllRaw(ctx context.Context, activeOnly bool) ([]ChartOfAccount, error) {
	var accounts []ChartOfAccount
	query := r.db.WithContext(ctx).
		Table("chart_of_accounts").
		Where("deleted_at IS NULL")

	if activeOnly {
		query = query.Where("is_active = true")
	}

	err := query.Order("account_level ASC, code ASC").Find(&accounts).Error
	return accounts, err
}

func (r *repository) FindChildrenCount(ctx context.Context, parentID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("chart_of_accounts").
		Where("parent_id = ? AND deleted_at IS NULL", parentID).
		Count(&count).Error
	return count, err
}

func (r *repository) IsCodeUsedInTariff(ctx context.Context, code string) (bool, error) {
	// 1. Cek di master komponen tarif
	var compCount int64
	err := r.db.WithContext(ctx).
		Table("tariff_components").
		Where("default_coa_code = ? AND deleted_at IS NULL", code).
		Count(&compCount).Error
	if err != nil {
		return false, err
	}
	if compCount > 0 {
		return true, nil
	}

	// 2. Cek di rincian komponen buku tarif
	var planCompCount int64
	err = r.db.WithContext(ctx).
		Table("tariff_price_plan_item_components").
		Where("coa_code = ?", code).
		Count(&planCompCount).Error
	if err != nil {
		return false, err
	}
	if planCompCount > 0 {
		return true, nil
	}

	return false, nil
}

func (r *repository) SetPostable(ctx context.Context, id string, isPostable bool) error {
	return r.db.WithContext(ctx).
		Table("chart_of_accounts").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{
			"is_postable": isPostable,
			"updated_at":  time.Now(),
		}).Error
}
