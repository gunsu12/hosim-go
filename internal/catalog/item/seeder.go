package item

import (
	"log"
	"time"
	"uuid"

	"hosim-go/pkg/enums"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func strPtr(s string) *string {
	return &s
}

// DefaultItemCategories memuat data referensi kategori item standar rumah sakit
var DefaultItemCategories = []ItemCategory{
	// MEDICATION
	{
		Code:            "CAT-MED-PATEN",
		Name:            "Obat Paten & Resep",
		ItemType:        enums.ItemTypeMedication,
		IncomeCOACode:   strPtr("410.02.001"),
		DiscountCOACode: strPtr("410.02.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-MED-GENERIK",
		Name:            "Obat Generik",
		ItemType:        enums.ItemTypeMedication,
		IncomeCOACode:   strPtr("410.02.003"),
		DiscountCOACode: strPtr("410.02.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-MED-OTC",
		Name:            "Obat Bebas & OTC (Over The Counter)",
		ItemType:        enums.ItemTypeMedication,
		IncomeCOACode:   strPtr("410.02.004"),
		DiscountCOACode: strPtr("410.02.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-MED-KRONIS",
		Name:            "Obat Kronis & PRB (Program Rujuk Balik)",
		ItemType:        enums.ItemTypeMedication,
		IncomeCOACode:   strPtr("410.02.005"),
		DiscountCOACode: strPtr("410.02.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},

	// GENERAL (BMHP, Linen, Logistik)
	{
		Code:            "CAT-GEN-BMHP-MEDIS",
		Name:            "Bahan Medis Habis Pakai (BMHP) Bedah & Rawat",
		ItemType:        enums.ItemTypeGeneral,
		IncomeCOACode:   strPtr("410.03.001"),
		DiscountCOACode: strPtr("410.03.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-GEN-LINEN-CSSD",
		Name:            "Linen Steril & Instrumen CSSD",
		ItemType:        enums.ItemTypeGeneral,
		IncomeCOACode:   strPtr("410.03.003"),
		DiscountCOACode: strPtr("410.03.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-GEN-LOGISTIK",
		Name:            "Logistik Umum, ATK & Kebersihan",
		ItemType:        enums.ItemTypeGeneral,
		IncomeCOACode:   strPtr("410.03.004"),
		DiscountCOACode: strPtr("410.03.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},

	// ASSET (Alat Modal / Alkes)
	{
		Code:            "CAT-AST-MEDIS",
		Name:            "Alat Kesehatan & Modal Medis",
		ItemType:        enums.ItemTypeAsset,
		IncomeCOACode:   strPtr("410.04.001"),
		DiscountCOACode: strPtr("410.04.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-AST-NONMEDIS",
		Name:            "Aset Tetap & Peralatan Non-Medis",
		ItemType:        enums.ItemTypeAsset,
		IncomeCOACode:   strPtr("410.04.003"),
		DiscountCOACode: strPtr("410.04.002"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},

	// TARIFF (Layanan & Tindakan Jasa Pasien)
	{
		Code:            "CAT-TAR-ADM",
		Name:            "Administrasi Pendaftaran & Rekam Medis",
		ItemType:        enums.ItemTypeTariff,
		IncomeCOACode:   strPtr("410.01.001"),
		DiscountCOACode: strPtr("410.01.005"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-TAR-POLI",
		Name:            "Konsultasi & Pemeriksaan Poliklinik Rawat Jalan",
		ItemType:        enums.ItemTypeTariff,
		IncomeCOACode:   strPtr("410.01.002"),
		DiscountCOACode: strPtr("410.01.005"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-TAR-RANAP",
		Name:            "Akomodasi & Sewa Kamar Rawat Inap",
		ItemType:        enums.ItemTypeTariff,
		IncomeCOACode:   strPtr("410.01.003"),
		DiscountCOACode: strPtr("410.01.005"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-TAR-TINDAKAN",
		Name:            "Jasa Tindakan Medis & Keperawatan (OK/IGD)",
		ItemType:        enums.ItemTypeTariff,
		IncomeCOACode:   strPtr("410.01.004"),
		DiscountCOACode: strPtr("410.01.005"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
	{
		Code:            "CAT-TAR-PENUNJANG",
		Name:            "Pemeriksaan Penunjang (Laboratorium & Radiologi)",
		ItemType:        enums.ItemTypeTariff,
		IncomeCOACode:   strPtr("410.01.006"),
		DiscountCOACode: strPtr("410.01.005"),
		SalesTaxCOACode: strPtr("215.01.001"),
		IsActive:        true,
	},
}

// DefaultItemProductLines memuat data referensi lini produk persediaan standar rumah sakit
var DefaultItemProductLines = []ItemProductLine{
	{
		Code:                    "PL-OBAT-FARMASI",
		Name:                    "Persediaan Obat Paten & Generik Farmasi",
		InventoryCOACode:        strPtr("114.01.001"),
		CogsCOACode:             strPtr("510.01.001"),
		PurchaseDiscountCOACode: strPtr("510.01.002"),
		PurchaseTaxCOACode:      strPtr("115.01.001"),
		IsActive:                true,
	},
	{
		Code:                    "PL-BMHP-MEDIS",
		Name:                    "Persediaan Bahan Medis Habis Pakai (BMHP)",
		InventoryCOACode:        strPtr("114.02.001"),
		CogsCOACode:             strPtr("510.02.001"),
		PurchaseDiscountCOACode: strPtr("510.02.002"),
		PurchaseTaxCOACode:      strPtr("115.01.001"),
		IsActive:                true,
	},
	{
		Code:                    "PL-LOGISTIK-UMUM",
		Name:                    "Persediaan Logistik Umum & ATK",
		InventoryCOACode:        strPtr("114.03.001"),
		CogsCOACode:             strPtr("510.03.001"),
		PurchaseDiscountCOACode: strPtr("510.03.002"),
		PurchaseTaxCOACode:      strPtr("115.01.001"),
		IsActive:                true,
	},
	{
		Code:                    "PL-GIZI-DAPUR",
		Name:                    "Persediaan Bahan Makanan & Gizi",
		InventoryCOACode:        strPtr("114.04.001"),
		CogsCOACode:             strPtr("510.04.001"),
		PurchaseDiscountCOACode: strPtr("510.04.002"),
		PurchaseTaxCOACode:      strPtr("115.01.001"),
		IsActive:                true,
	},
	{
		Code:                            "PL-ASET-MEDIS",
		Name:                            "Aset Tetap Peralatan Medis & Alkes Modal",
		InventoryCOACode:                strPtr("120.01.001"),
		CogsCOACode:                     strPtr("520.01.001"),
		AssetCOACode:                    strPtr("120.01.001"),
		AssetAccumulationCOACode:        strPtr("121.01.001"),
		AssetDepreciationExpenseCOACode: strPtr("530.01.001"),
		IsActive:                        true,
	},
	{
		Code:                            "PL-ASET-NONMEDIS",
		Name:                            "Aset Tetap Peralatan Non-Medis & Kantor",
		InventoryCOACode:                strPtr("120.02.001"),
		CogsCOACode:                     strPtr("520.02.001"),
		AssetCOACode:                    strPtr("120.02.001"),
		AssetAccumulationCOACode:        strPtr("121.02.001"),
		AssetDepreciationExpenseCOACode: strPtr("530.02.001"),
		IsActive:                        true,
	},
}

// SeedItemCategories mengisi data master kategori item secara idempoten
func SeedItemCategories(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding master Kategori Item (Item Categories)...")
	now := time.Now()

	for _, cat := range DefaultItemCategories {
		var count int64
		db.Model(&ItemCategory{}).Where("code = ?", cat.Code).Count(&count)
		if count == 0 {
			cat.ID = uuid.NewV7().String()
			cat.CreatedBy = "system-seeder"
			cat.UpdatedBy = "system-seeder"
			cat.CreatedAt = now
			cat.UpdatedAt = now
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&cat).Error; err != nil {
				return err
			}
		}
	}

	log.Println("[SEEDER] Seeding master Kategori Item berhasil.")
	return nil
}

// SeedItemProductLines mengisi data master lini produk persediaan secara idempoten
func SeedItemProductLines(db *gorm.DB) error {
	log.Println("[SEEDER] Memeriksa & melakukan seeding master Lini Produk Item (Item Product Lines)...")
	now := time.Now()

	for _, pl := range DefaultItemProductLines {
		var count int64
		db.Model(&ItemProductLine{}).Where("code = ?", pl.Code).Count(&count)
		if count == 0 {
			pl.ID = uuid.NewV7().String()
			pl.CreatedBy = "system-seeder"
			pl.UpdatedBy = "system-seeder"
			pl.CreatedAt = now
			pl.UpdatedAt = now
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&pl).Error; err != nil {
				return err
			}
		}
	}

	log.Println("[SEEDER] Seeding master Lini Produk Item berhasil.")
	return nil
}

// SeedCatalogMasterData menjalankan seluruh seeder untuk master kategori dan lini produk
func SeedCatalogMasterData(db *gorm.DB) error {
	if err := SeedItemCategories(db); err != nil {
		return err
	}
	if err := SeedItemProductLines(db); err != nil {
		return err
	}
	return nil
}
