package seeder

import (
	"fmt"
	"log"

	"hosim-go/internal/auth"
	"hosim-go/internal/master/practitioner"

	"gorm.io/gorm"
)

// RunAll menjalankan seluruh seeder esensial (RBAC, Referensi SatuSehat, Tenaga Medis Default)
func RunAll(db *gorm.DB) error {
	log.Println("[SEEDER] Memulai proses seeding database...")

	if err := RunRBAC(db); err != nil {
		return fmt.Errorf("seeder RBAC gagal: %w", err)
	}

	if err := RunSatuSehatReferences(db); err != nil {
		return fmt.Errorf("seeder referensi SatuSehat gagal: %w", err)
	}

	if err := RunPractitioners(db); err != nil {
		return fmt.Errorf("seeder practitioner gagal: %w", err)
	}

	log.Println("[SEEDER] Seluruh data seeder berhasil diinisialisasi.")
	return nil
}

// RunRBAC menginisialisasi permissions, roles, mapping role-permission, dan superadmin
func RunRBAC(db *gorm.DB) error {
	log.Println("[SEEDER] Menjalankan seeder Roles, Permissions & Default Admin...")
	return auth.SeedRBAC(db)
}

// RunSatuSehatReferences menginisialisasi data master Profesi dan Spesialisasi standar SatuSehat
func RunSatuSehatReferences(db *gorm.DB) error {
	log.Println("[SEEDER] Menjalankan seeder Profesi & Spesialisasi (SatuSehat)...")
	return practitioner.SeedProfessionsAndSpecialties(db)
}

// RunPractitioners menginisialisasi contoh data dokter spesialis default
func RunPractitioners(db *gorm.DB) error {
	log.Println("[SEEDER] Menjalankan seeder Tenaga Medis (Dokter)...")
	return practitioner.SeedDefaultPractitioner(db)
}
