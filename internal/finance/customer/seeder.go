package customer

import (
	"log"
	"time"
	"uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DefaultCustomerTypes memuat referensi kategori/tipe penjamin rumah sakit standar
var DefaultCustomerTypes = []CustomerType{
	{
		Code: "PERSONAL",
		Name: "Pribadi / Perorangan / Umum",
	},
	{
		Code: "BPJS",
		Name: "Pemerintah / BPJS Kesehatan",
	},
	{
		Code: "ASURANSI",
		Name: "Asuransi Swasta Komersial",
	},
	{
		Code: "PERUSAHAAN",
		Name: "Perusahaan / Korporasi Rekanan",
	},
}

// SeedCustomerTypes mengisi master kategori/tipe penjamin secara idempoten
func SeedCustomerTypes(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding master Kategori Penjamin (Customer Types)...")
	now := time.Now()

	for _, ct := range DefaultCustomerTypes {
		var count int64
		db.Model(&CustomerType{}).Where("code = ?", ct.Code).Count(&count)
		if count == 0 {
			ct.ID = uuid.NewV7().String()
			ct.CreatedBy = "system-seeder"
			ct.UpdatedBy = "system-seeder"
			ct.CreatedAt = now
			ct.UpdatedAt = now
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ct).Error; err != nil {
				return err
			}
		}
	}

	log.Println("[SEEDER] Seeding master Kategori Penjamin berhasil.")
	return nil
}

// SeedDefaultImmutableCustomers menginisialisasi tipe penjamin dan customer immutable sistem untuk tipe penjamin PERSONAL:
// 1. PERSONAL (Pasien Personal / Umum - Pribadi / Mandiri)
// 2. KITAS (Pasien WNA Pemegang KITAS / KITAP)
// 3. PASSPORT (Pasien WNA Paspor / Asing - Turis / Non-KITAS)
func SeedDefaultImmutableCustomers(db *gorm.DB) error {
	// 1. Pastikan Customer Types sudah ada
	if err := SeedCustomerTypes(db); err != nil {
		return err
	}

	log.Println("[SEEDER] Memeriksa & melakukan seeding penjamin sistem / immutable (Personal, KITAS, Pasport/Asing)...")

	// 2. Cari tipe customer PERSONAL
	var personalType CustomerType
	if err := db.Where("code = ?", "PERSONAL").First(&personalType).Error; err != nil {
		return err
	}

	falseVal := false
	trueVal := true
	descPersonal := "Penjamin untuk pasien umum / personal / bayar mandiri (out-of-pocket)."
	descKitas := "Penjamin untuk Warga Negara Asing (WNA) pemegang kartu izin tinggal terbatas/tetap (KITAS/KITAP)."
	descPassport := "Penjamin untuk Warga Negara Asing (WNA) wisatawan / pemegang paspor non-KITAS."

	immutableCustomers := []Customer{
		{
			Code:           "PERSONAL",
			Name:           "Pasien Personal / Umum (Pribadi)",
			CustomerTypeID: &personalType.ID,
			RequireCard:    &falseVal,
			Description:    &descPersonal,
			IsActive:       true,
			IsImmutable:    true,
		},
		{
			Code:           "KITAS",
			Name:           "Pasien WNA Pemegang KITAS / KITAP",
			CustomerTypeID: &personalType.ID,
			RequireCard:    &trueVal,
			Description:    &descKitas,
			IsActive:       true,
			IsImmutable:    true,
		},
		{
			Code:           "PASSPORT",
			Name:           "Pasien WNA Paspor / Asing (Turis)",
			CustomerTypeID: &personalType.ID,
			RequireCard:    &trueVal,
			Description:    &descPassport,
			IsActive:       true,
			IsImmutable:    true,
		},
	}

	now := time.Now()
	for _, cust := range immutableCustomers {
		var existing Customer
		err := db.Where("code = ?", cust.Code).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			cust.ID = uuid.NewV7().String()
			cust.CreatedBy = "system-seeder"
			cust.UpdatedBy = "system-seeder"
			cust.CreatedAt = now
			cust.UpdatedAt = now
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&cust).Error; err != nil {
				return err
			}
			log.Printf("[SEEDER] Penjamin immutable '%s' (%s) berhasil dibuat.\n", cust.Code, cust.Name)
		} else if err == nil {
			// Jika data sudah ada, pastikan is_immutable bernilai true dan customer_type_id terhubung
			updates := map[string]interface{}{
				"is_immutable":     true,
				"customer_type_id": personalType.ID,
			}
			if err := db.Model(&existing).Updates(updates).Error; err != nil {
				return err
			}
			log.Printf("[SEEDER] Penjamin '%s' diperbarui menjadi immutable.\n", cust.Code)
		} else {
			return err
		}
	}

	log.Println("[SEEDER] Seeding penjamin immutable berhasil.")
	return nil
}
