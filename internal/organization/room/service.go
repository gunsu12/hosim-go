package room

import (
	"context"
	"errors"
	"strings"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
)

var (
	ErrRoomNotFound          = errors.New("ruangan tidak ditemukan")
	ErrRoomCodeAlreadyExists = errors.New("ruangan dengan kode tersebut sudah terdaftar")
	ErrInvalidRoomType       = errors.New("tipe ruangan tidak valid (pilihan: OUTPATIENT, INPATIENT, EMERGENCY, OPERATING)")
	ErrNegativeCapacity      = errors.New("kapasitas ruangan tidak boleh negatif")
	ErrCodeRequired          = errors.New("kode ruangan wajib diisi")
	ErrNameRequired          = errors.New("nama ruangan wajib diisi")
	ErrServiceUnitIDRequired = errors.New("unit layanan (service_unit_id) wajib diisi")
)

type CreateRoomRequest struct {
	Code          string         `json:"code" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	Capacity      int            `json:"capacity"`
	Location      string         `json:"location"`
	RoomType      enums.RoomType `json:"room_type" binding:"required"`
	ServiceUnitID string         `json:"service_unit_id" binding:"required"`
}

type UpdateRoomRequest struct {
	Code          string         `json:"code" binding:"required"`
	Name          string         `json:"name" binding:"required"`
	Capacity      int            `json:"capacity"`
	Location      string         `json:"location"`
	RoomType      enums.RoomType `json:"room_type" binding:"required"`
	ServiceUnitID string         `json:"service_unit_id" binding:"required"`
}

func isValidRoomType(t enums.RoomType) bool {
	switch t {
	case enums.RoomTypeOutpatient,
		enums.RoomTypeInpatient,
		enums.RoomTypeEmergency,
		enums.RoomTypeOperating:
		return true
	default:
		return false
	}
}

type Service interface {
	CreateRoom(ctx context.Context, req CreateRoomRequest, operatorID string) (*Room, error)
	UpdateRoom(ctx context.Context, id string, req UpdateRoomRequest, operatorID string) (*Room, error)
	DeleteRoom(ctx context.Context, id string, operatorID string) error
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	ListRooms(ctx context.Context, params ListParams) ([]Room, int64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateRoom(ctx context.Context, req CreateRoomRequest, operatorID string) (*Room, error) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Location = strings.TrimSpace(req.Location)
	req.ServiceUnitID = strings.TrimSpace(req.ServiceUnitID)

	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if req.ServiceUnitID == "" {
		return nil, ErrServiceUnitIDRequired
	}
	if req.Capacity < 0 {
		return nil, ErrNegativeCapacity
	}
	if !isValidRoomType(req.RoomType) {
		return nil, ErrInvalidRoomType
	}

	existing, err := s.repo.FindByCode(ctx, req.Code)
	if err == nil && existing != nil {
		return nil, ErrRoomCodeAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	room := &Room{
		Code:          req.Code,
		Name:          req.Name,
		Capacity:      req.Capacity,
		Location:      req.Location,
		RoomType:      req.RoomType,
		ServiceUnitID: req.ServiceUnitID,
		CreatedBy:     operatorID,
		UpdatedBy:     operatorID,
	}

	return s.repo.Create(ctx, room)
}

func (s *service) UpdateRoom(ctx context.Context, id string, req UpdateRoomRequest, operatorID string) (*Room, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Location = strings.TrimSpace(req.Location)
	req.ServiceUnitID = strings.TrimSpace(req.ServiceUnitID)

	if req.Code == "" {
		return nil, ErrCodeRequired
	}
	if req.Name == "" {
		return nil, ErrNameRequired
	}
	if req.ServiceUnitID == "" {
		return nil, ErrServiceUnitIDRequired
	}
	if req.Capacity < 0 {
		return nil, ErrNegativeCapacity
	}
	if !isValidRoomType(req.RoomType) {
		return nil, ErrInvalidRoomType
	}

	if req.Code != room.Code {
		existing, err := s.repo.FindByCode(ctx, req.Code)
		if err == nil && existing != nil && existing.ID != id {
			return nil, ErrRoomCodeAlreadyExists
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	room.Code = req.Code
	room.Name = req.Name
	room.Capacity = req.Capacity
	room.Location = req.Location
	room.RoomType = req.RoomType
	room.ServiceUnitID = req.ServiceUnitID
	room.UpdatedBy = operatorID

	return s.repo.Update(ctx, room)
}

func (s *service) DeleteRoom(ctx context.Context, id string, operatorID string) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoomNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id, operatorID)
}

func (s *service) GetRoomByID(ctx context.Context, id string) (*Room, error) {
	room, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRoomNotFound
		}
		return nil, err
	}
	return room, nil
}

func (s *service) ListRooms(ctx context.Context, params ListParams) ([]Room, int64, error) {
	return s.repo.FindAll(ctx, params)
}
