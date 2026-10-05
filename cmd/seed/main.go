package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"hosim-go/internal/config"
	"hosim-go/internal/database"
	"hosim-go/internal/database/seeder"
)

const usage = `Perintah Seeder Database (PostgreSQL)

Penggunaan:
  go run cmd/seed/main.go [opsi]

Opsi Flag:
  -type string   Kategori seeder yang ingin dijalankan (default: "all")
                 Pilihan: all, rbac, satusehat, practitioner, catalog, accounting, coa-head

Contoh:
  go run cmd/seed/main.go
  go run cmd/seed/main.go -type=accounting
  go run cmd/seed/main.go -type=coa-head
  go run cmd/seed/main.go -type=catalog
  go run cmd/seed/main.go -type=rbac
`

func main() {
	seedType := flag.String("type", "all", "Kategori seeder: all | rbac | satusehat | practitioner | catalog | accounting | coa-head")
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	// 1. Baca Konfigurasi
	cfg := config.LoadConfig()

	// 2. Hubungkan ke Database PostgreSQL
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Gagal membuka koneksi database: %v\n", err)
	}

	// 3. Eksekusi Seeder Sesuai Tipe
	log.Printf("[SEEDER] Menjalankan kategori seeder '%s' pada database %s...\n", *seedType, cfg.DBName)

	var seedErr error
	switch *seedType {
	case "all":
		seedErr = seeder.RunAll(db)
	case "rbac":
		seedErr = seeder.RunRBAC(db)
	case "satusehat":
		seedErr = seeder.RunSatuSehatReferences(db)
	case "practitioner":
		seedErr = seeder.RunPractitioners(db)
	case "catalog":
		seedErr = seeder.RunCatalog(db)
	case "accounting", "coa":
		seedErr = seeder.RunAccountingCOA(db)
	case "accounting-head", "coa-head":
		seedErr = seeder.RunAccountingHeadAccounts(db)
	default:
		log.Printf("[ERROR] Tipe seeder '%s' tidak valid. Gunakan: all, rbac, satusehat, practitioner, catalog, accounting, coa-head\n", *seedType)
		flag.Usage()
		os.Exit(1)
	}

	if seedErr != nil {
		log.Fatalf("[FATAL] Eksekusi seeder gagal: %v\n", seedErr)
	}

	log.Println("[SUCCESS] Seeding data database selesai.")
}
