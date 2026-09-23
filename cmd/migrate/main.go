package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"hosim-go/internal/config"
	"hosim-go/internal/database"
	"hosim-go/migrations"

	"github.com/pressly/goose/v3"
)

const usage = `Perintah Migrasi Database Goose (PostgreSQL)

Penggunaan:
  go run cmd/migrate/main.go <perintah> [argumen]

Daftar Perintah:
  up                   Jalankan semua migrasi pending
  up-by-one            Jalankan 1 migrasi berikutnya
  up-to VERSION        Migrasi naik ke versi tertentu
  down                 Rollback 1 migrasi terakhir
  down-to VERSION      Rollback turun ke versi tertentu
  status               Tampilkan status seluruh migrasi
  version              Tampilkan nomor versi skema database saat ini
  reset                Rollback seluruh migrasi dari awal
  create NAME [sql]    Buat template file migrasi SQL baru di direktori migrations/

Contoh:
  go run cmd/migrate/main.go status
  go run cmd/migrate/main.go up
  go run cmd/migrate/main.go down
  go run cmd/migrate/main.go create create_billing_table sql
`

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]

	// Jika perintahnya adalah membuat file migrasi baru (create), tulis langsung ke filesystem lokal
	if command == "create" {
		if len(args) < 2 {
			log.Fatal("[ERROR] Harap tentukan nama file migrasi. Contoh: go run cmd/migrate/main.go create nama_migrasi sql")
		}
		migrationType := "sql"
		if len(args) >= 3 {
			migrationType = args[2]
		}
		if err := goose.Create(nil, "migrations", args[1], migrationType); err != nil {
			log.Fatalf("[ERROR] Gagal membuat file migrasi: %v\n", err)
		}
		log.Printf("[SUCCESS] File migrasi berhasil dibuat di direktori migrations/\n")
		return
	}

	// 1. Baca Konfigurasi Lingkungan
	cfg := config.LoadConfig()

	// 2. Hubungkan ke Database PostgreSQL
	gormDB, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Gagal membuka koneksi database: %v\n", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Fatalf("[FATAL] Gagal mengakses generic sql.DB: %v\n", err)
	}
	defer sqlDB.Close()

	// 3. Set Dialect PostgreSQL & Gunakan Embedded Migration Files
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("[FATAL] Gagal mengatur dialect goose: %v\n", err)
	}

	goose.SetBaseFS(migrations.FS)

	// 4. Eksekusi Perintah Goose
	ctx := context.Background()
	cmdArgs := args[1:]

	log.Printf("[GOOSE] Menjalankan perintah '%s' pada database %s...\n", command, cfg.DBName)
	if err := goose.RunContext(ctx, command, sqlDB, ".", cmdArgs...); err != nil {
		log.Fatalf("[ERROR] Eksekusi migrasi '%s' gagal: %v\n", command, err)
	}

	log.Println("[GOOSE] Operasi migrasi selesai dengan sukses.")
}
