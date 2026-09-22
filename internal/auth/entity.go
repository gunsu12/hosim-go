package auth

import (
	"time"
	"uuid"

	"gorm.io/gorm"
)

// Permission merepresentasikan izin aksi atau akses menu spesifik dalam sistem
type Permission struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	Code        string    `gorm:"uniqueIndex:idx_permissions_code;size:100;not null" json:"code"` // e.g. "department:read", "outpatient:register"
	Name        string    `gorm:"size:150;not null" json:"name"`                                 // e.g. "Lihat Master Departemen"
	Module      string    `gorm:"size:100;index;not null" json:"module"`                          // e.g. "Master Departemen", "Rawat Jalan"
	Description string    `gorm:"size:255" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.NewV7().String()
	}
	return nil
}

// Role merepresentasikan peran pengguna dalam sistem (e.g. ADMIN, DOCTOR, NURSE, PARAMEDIC, STAFF)
type Role struct {
	ID          string       `gorm:"primaryKey;size:36" json:"id"`
	Code        string       `gorm:"uniqueIndex:idx_roles_code;size:50;not null" json:"code"` // e.g. "ADMIN", "DOCTOR"
	Name        string       `gorm:"size:150;not null" json:"name"`                           // e.g. "Administrator", "Dokter"
	Description string       `gorm:"size:255" json:"description,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewV7().String()
	}
	return nil
}

// User merepresentasikan entitas akun pengguna dalam sistem
type User struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	Username  string         `gorm:"uniqueIndex:idx_users_username;size:100;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex:idx_users_email;size:150;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Name      string         `gorm:"size:150;not null" json:"name"`
	RoleID    *string        `gorm:"size:36;index" json:"role_id,omitempty"`
	Role      *Role          `gorm:"foreignKey:RoleID;references:ID" json:"role,omitempty"`
	IsActive  bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP;OnUpdate:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy string         `gorm:"default:'SYSTEM'" json:"-"`
	UpdatedBy string         `gorm:"default:'SYSTEM'" json:"-"`
	DeletedBy string         `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.NewV7().String()
	}
	return nil
}

// GetRoleCode mengembalikan kode role pengguna (contoh: "ADMIN", "DOCTOR")
func (u *User) GetRoleCode() string {
	if u.Role != nil {
		return u.Role.Code
	}
	return ""
}

// GetPermissionCodes mengembalikan daftar kode izin yang dimiliki pengguna
func (u *User) GetPermissionCodes() []string {
	if u.Role == nil || len(u.Role.Permissions) == 0 {
		return []string{}
	}
	codes := make([]string, len(u.Role.Permissions))
	for i, p := range u.Role.Permissions {
		codes[i] = p.Code
	}
	return codes
}

// RefreshToken merepresentasikan token refresh sesi user yang tersimpan di database
type RefreshToken struct {
	ID        string     `gorm:"primaryKey;size:36" json:"id"`
	UserID    string     `gorm:"size:36;index;not null" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Token     string     `gorm:"uniqueIndex:idx_refresh_tokens_token;size:255;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	RevokedAt *time.Time `gorm:"index" json:"revoked_at,omitempty"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == "" {
		rt.ID = uuid.NewV7().String()
	}
	return nil
}

// IsValid mengecek apakah refresh token belum kedaluwarsa dan belum di-revoke
func (rt *RefreshToken) IsValid() bool {
	if rt.RevokedAt != nil {
		return false
	}
	return time.Now().Before(rt.ExpiresAt)
}
