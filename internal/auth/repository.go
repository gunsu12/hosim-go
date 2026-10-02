package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsernameOrEmail(ctx context.Context, identifier string) (*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	CountUsers(ctx context.Context) (int64, error)
	ListUsers(ctx context.Context, search string, roleID string, isActive *bool, page, limit int) ([]User, int64, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string, deletedBy string) error

	SaveRefreshToken(ctx context.Context, token *RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	RevokeAllUserTokens(ctx context.Context, userID string) error

	// Role operations
	FindRoleByCode(ctx context.Context, code string) (*Role, error)
	FindRoleByID(ctx context.Context, id string) (*Role, error)
	FindAllRoles(ctx context.Context) ([]Role, error)
	CreateRole(ctx context.Context, role *Role, permIDs []string) (*Role, error)
	UpdateRole(ctx context.Context, role *Role, permIDs []string) (*Role, error)
	DeleteRole(ctx context.Context, id string) error

	// Permission operations
	FindAllPermissions(ctx context.Context) ([]Permission, error)
	FindPermissionByID(ctx context.Context, id string) (*Permission, error)
	FindPermissionByCode(ctx context.Context, code string) (*Permission, error)
	CreatePermission(ctx context.Context, perm *Permission) (*Permission, error)
	UpdatePermission(ctx context.Context, perm *Permission) error
	DeletePermission(ctx context.Context, id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Preload("Role.Permissions").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Preload("Role.Permissions").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByUsernameOrEmail(ctx context.Context, identifier string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Preload("Role.Permissions").Where("username = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindUserByID(ctx context.Context, id string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Preload("Role.Permissions").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) ListUsers(ctx context.Context, search string, roleID string, isActive *bool, page, limit int) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})
	if search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ? OR LOWER(name) LIKE ?", s, s, s)
	}
	if roleID != "" {
		query = query.Where("role_id = ?", roleID)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	if err := query.Preload("Role.Permissions").Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *repository) DeleteUser(ctx context.Context, id string, deletedBy string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_by": deletedBy,
		"is_active":  false,
		"deleted_at": &now,
	}).Error
}

func (r *repository) SaveRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *repository) FindRefreshToken(ctx context.Context, token string) (*RefreshToken, error) {
	var rt RefreshToken
	if err := r.db.WithContext(ctx).Preload("User.Role.Permissions").Where("token = ?", token).First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, token string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("token = ? AND revoked_at IS NULL", token).
		Update("revoked_at", &now).Error
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}

func (r *repository) FindRoleByCode(ctx context.Context, code string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) FindRoleByID(ctx context.Context, id string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) FindAllRoles(ctx context.Context) ([]Role, error) {
	var roles []Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Order("code ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *repository) CreateRole(ctx context.Context, role *Role, permIDs []string) (*Role, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}
		if len(permIDs) > 0 {
			var perms []Permission
			if err := tx.Where("id IN ?", permIDs).Find(&perms).Error; err != nil {
				return err
			}
			if err := tx.Model(role).Association("Permissions").Replace(perms); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.FindRoleByID(ctx, role.ID)
}

func (r *repository) UpdateRole(ctx context.Context, role *Role, permIDs []string) (*Role, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(role).Updates(map[string]interface{}{
			"name":        role.Name,
			"description": role.Description,
		}).Error; err != nil {
			return err
		}
		if permIDs != nil {
			var perms []Permission
			if len(permIDs) > 0 {
				if err := tx.Where("id IN ?", permIDs).Find(&perms).Error; err != nil {
					return err
				}
			}
			if err := tx.Model(role).Association("Permissions").Replace(perms); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.FindRoleByID(ctx, role.ID)
}

func (r *repository) DeleteRole(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&User{}).Where("role_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("role tidak dapat dihapus karena masih digunakan oleh pengguna aktif")
		}
		var role Role
		if err := tx.First(&role, "id = ?", id).Error; err != nil {
			return err
		}
		_ = tx.Model(&role).Association("Permissions").Clear()
		return tx.Delete(&role).Error
	})
}

func (r *repository) FindAllPermissions(ctx context.Context) ([]Permission, error) {
	var perms []Permission
	if err := r.db.WithContext(ctx).Order("module ASC, code ASC").Find(&perms).Error; err != nil {
		return nil, err
	}
	return perms, nil
}

func (r *repository) FindPermissionByID(ctx context.Context, id string) (*Permission, error) {
	var perm Permission
	if err := r.db.WithContext(ctx).First(&perm, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *repository) FindPermissionByCode(ctx context.Context, code string) (*Permission, error) {
	var perm Permission
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (r *repository) CreatePermission(ctx context.Context, perm *Permission) (*Permission, error) {
	if err := r.db.WithContext(ctx).Create(perm).Error; err != nil {
		return nil, err
	}
	return perm, nil
}

func (r *repository) UpdatePermission(ctx context.Context, perm *Permission) error {
	return r.db.WithContext(ctx).Model(perm).Updates(map[string]interface{}{
		"name":        perm.Name,
		"module":      perm.Module,
		"description": perm.Description,
	}).Error
}

func (r *repository) DeletePermission(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var perm Permission
		if err := tx.First(&perm, "id = ?", id).Error; err != nil {
			return err
		}
		_ = tx.Table("role_permissions").Where("permission_id = ?", id).Delete(nil)
		return tx.Delete(&perm).Error
	})
}
