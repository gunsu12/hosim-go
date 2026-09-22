package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DefaultPermissions memuat seluruh permission bawaan sistem
var DefaultPermissions = []Permission{
	// Departemen
	{Code: "department:read", Name: "Lihat Departemen", Module: "Master Departemen"},
	{Code: "department:create", Name: "Tambah Departemen", Module: "Master Departemen"},
	{Code: "department:update", Name: "Ubah Departemen", Module: "Master Departemen"},
	{Code: "department:delete", Name: "Hapus Departemen", Module: "Master Departemen"},

	// Ruangan
	{Code: "room:read", Name: "Lihat Ruangan", Module: "Master Ruangan"},
	{Code: "room:create", Name: "Tambah Ruangan", Module: "Master Ruangan"},
	{Code: "room:update", Name: "Ubah Ruangan", Module: "Master Ruangan"},
	{Code: "room:delete", Name: "Hapus Ruangan", Module: "Master Ruangan"},

	// Unit Layanan
	{Code: "service_unit:read", Name: "Lihat Unit Layanan", Module: "Master Unit Layanan"},
	{Code: "service_unit:create", Name: "Tambah Unit Layanan", Module: "Master Unit Layanan"},
	{Code: "service_unit:update", Name: "Ubah Unit Layanan", Module: "Master Unit Layanan"},
	{Code: "service_unit:delete", Name: "Hapus Unit Layanan", Module: "Master Unit Layanan"},

	// Penjamin
	{Code: "payer:read", Name: "Lihat Penjamin", Module: "Master Penjamin"},
	{Code: "payer:create", Name: "Tambah Penjamin", Module: "Master Penjamin"},
	{Code: "payer:update", Name: "Ubah Penjamin", Module: "Master Penjamin"},
	{Code: "payer:delete", Name: "Hapus Penjamin", Module: "Master Penjamin"},

	// Rujukan
	{Code: "referal:read", Name: "Lihat Rujukan", Module: "Master Rujukan"},
	{Code: "referal:create", Name: "Tambah Rujukan", Module: "Master Rujukan"},
	{Code: "referal:update", Name: "Ubah Rujukan", Module: "Master Rujukan"},
	{Code: "referal:delete", Name: "Hapus Rujukan", Module: "Master Rujukan"},

	// Kelas Tarif
	{Code: "tariff_class:read", Name: "Lihat Kelas Tarif", Module: "Master Kelas Tarif"},
	{Code: "tariff_class:create", Name: "Tambah Kelas Tarif", Module: "Master Kelas Tarif"},
	{Code: "tariff_class:update", Name: "Ubah Kelas Tarif", Module: "Master Kelas Tarif"},
	{Code: "tariff_class:delete", Name: "Hapus Kelas Tarif", Module: "Master Kelas Tarif"},

	// Pasien
	{Code: "patient:read", Name: "Lihat Pasien", Module: "Master Pasien"},
	{Code: "patient:create", Name: "Tambah Pasien", Module: "Master Pasien"},
	{Code: "patient:update", Name: "Ubah Pasien", Module: "Master Pasien"},
	{Code: "patient:delete", Name: "Hapus Pasien", Module: "Master Pasien"},

	// Tenaga Medis
	{Code: "practitioner:read", Name: "Lihat Tenaga Medis", Module: "Master Tenaga Medis"},
	{Code: "practitioner:create", Name: "Tambah Tenaga Medis", Module: "Master Tenaga Medis"},
	{Code: "practitioner:update", Name: "Ubah Tenaga Medis", Module: "Master Tenaga Medis"},
	{Code: "practitioner:delete", Name: "Hapus Tenaga Medis", Module: "Master Tenaga Medis"},

	// Modul Rawat Jalan (Outpatient & Appointment)
	{Code: "outpatient:view", Name: "Buka Menu Rawat Jalan", Module: "Rawat Jalan"},
	{Code: "outpatient:register", Name: "Pendaftaran Rawat Jalan", Module: "Rawat Jalan"},
	{Code: "appointment:view", Name: "Lihat Jadwal Janji Temu", Module: "Rawat Jalan"},
	{Code: "appointment:create", Name: "Buat Janji Temu", Module: "Rawat Jalan"},
}

// SeedRBAC menginisialisasi tabel permissions, roles, mapping role_permissions, dan default admin
func SeedRBAC(db *gorm.DB) error {
	// 1. Seed Permissions (Upsert berdasarkan Code)
	for _, p := range DefaultPermissions {
		perm := p
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "module"}),
		}).Create(&perm).Error; err != nil {
			return err
		}
	}

	// Ambil semua permissions dari database
	var allPerms []Permission
	if err := db.Find(&allPerms).Error; err != nil {
		return err
	}

	permMap := make(map[string]Permission, len(allPerms))
	for _, p := range allPerms {
		permMap[p.Code] = p
	}

	// 2. Definisi Roles Bawaan
	defaultRoles := []struct {
		Role        Role
		PermCodes   []string
		AllPermFlag bool
	}{
		{
			Role: Role{
				Code:        "ADMIN",
				Name:        "Administrator Sistem",
				Description: "Akses penuh seluruh modul dan konfigurasi",
			},
			AllPermFlag: true,
		},
		{
			Role: Role{
				Code:        "DOCTOR",
				Name:        "Dokter / Practitioner",
				Description: "Pelayanan medis, rekam medis, dan pendaftaran",
			},
			PermCodes: []string{
				"department:read", "room:read", "service_unit:read", "payer:read", "referal:read", "tariff_class:read",
				"practitioner:read", "patient:read", "patient:create", "patient:update",
				"outpatient:view", "appointment:view",
			},
		},
		{
			Role: Role{
				Code:        "PARAMEDIC",
				Name:        "Perawat / Paramedis",
				Description: "Petugas keperawatan dan pendaftaran rawat jalan",
			},
			PermCodes: []string{
				"department:read", "room:read", "service_unit:read", "payer:read", "referal:read", "tariff_class:read",
				"patient:read", "patient:create", "patient:update",
				"outpatient:view", "outpatient:register", "appointment:view",
			},
		},
		{
			Role: Role{
				Code:        "STAFF",
				Name:        "Staf Administrasi",
				Description: "Melihat data master umum",
			},
			PermCodes: []string{
				"department:read", "room:read", "service_unit:read", "payer:read", "referal:read", "tariff_class:read",
				"practitioner:read", "patient:read",
			},
		},
	}

	for _, dr := range defaultRoles {
		role := dr.Role
		// Cari apakah role sudah ada
		var existingRole Role
		if err := db.Where("code = ?", role.Code).First(&existingRole).Error; err != nil {
			if err := db.Create(&role).Error; err != nil {
				return err
			}
			existingRole = role
		}

		// Tentukan permissions untuk role ini
		var assignedPerms []Permission
		if dr.AllPermFlag {
			assignedPerms = allPerms
		} else {
			for _, code := range dr.PermCodes {
				if p, ok := permMap[code]; ok {
					assignedPerms = append(assignedPerms, p)
				}
			}
		}

		// Asosiasikan permissions ke role
		if err := db.Model(&existingRole).Association("Permissions").Replace(assignedPerms); err != nil {
			return err
		}
	}

	// 3. Seed Default Admin User
	var adminRole Role
	if err := db.Where("code = ?", "ADMIN").First(&adminRole).Error; err != nil {
		return err
	}

	var adminCount int64
	db.Model(&User{}).Where("username = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		admin := User{
			Username:  "admin",
			Email:     "admin@hosim.local",
			Password:  string(hashedPassword),
			Name:      "Administrator Sistem",
			RoleID:    &adminRole.ID,
			IsActive:  true,
			CreatedBy: "SEEDER",
			UpdatedBy: "SEEDER",
		}

		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		log.Println("[SEEDER] Default akun admin berhasil dibuat (username: admin, password: admin123)")
	}

	log.Println("[SEEDER] Seluruh Role & Permission RBAC berhasil diinisialisasi.")
	return nil
}
