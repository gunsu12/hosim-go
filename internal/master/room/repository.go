package room

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ListParams struct {
	Page          int
	Limit         int
	Search        string
	ServiceUnitID string
	RoomType      string
}

type Repository interface {
	Create(ctx context.Context, r *Room) (*Room, error)
	Update(ctx context.Context, r *Room) (*Room, error)
	Delete(ctx context.Context, id string, deletedBy string) error
	FindByID(ctx context.Context, id string) (*Room, error)
	FindByCode(ctx context.Context, code string) (*Room, error)
	FindAll(ctx context.Context, params ListParams) ([]Room, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, room *Room) (*Room, error) {
	result := r.db.WithContext(ctx).Create(room)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, room.ID)
}

func (r *repository) Update(ctx context.Context, room *Room) (*Room, error) {
	result := r.db.WithContext(ctx).Save(room)
	if result.Error != nil {
		return nil, result.Error
	}
	return r.FindByID(ctx, room.ID)
}

func (r *repository) Delete(ctx context.Context, id string, deletedBy string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Room{}).Where("id = ?", id).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		return tx.Delete(&Room{}, "id = ?", id).Error
	})
}

func (r *repository) FindByID(ctx context.Context, id string) (*Room, error) {
	var rm Room
	err := r.db.WithContext(ctx).Preload("ServiceUnit").First(&rm, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *repository) FindByCode(ctx context.Context, code string) (*Room, error) {
	var rm Room
	err := r.db.WithContext(ctx).Preload("ServiceUnit").First(&rm, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *repository) FindAll(ctx context.Context, params ListParams) ([]Room, int64, error) {
	var rooms []Room
	var count int64

	query := r.db.WithContext(ctx).Model(&Room{})

	if params.ServiceUnitID != "" {
		query = query.Where("service_unit_id = ?", params.ServiceUnitID)
	}

	if params.RoomType != "" {
		query = query.Where("room_type = ?", params.RoomType)
	}

	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(strings.TrimSpace(params.Search)) + "%"
		query = query.Where(
			"LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(location) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	offset := (params.Page - 1) * params.Limit

	err := query.Preload("ServiceUnit").Order("created_at DESC").Limit(params.Limit).Offset(offset).Find(&rooms).Error
	if err != nil {
		return nil, 0, err
	}

	return rooms, count, nil
}
