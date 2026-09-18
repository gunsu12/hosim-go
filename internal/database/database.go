package database

import (
	"fmt"
	"log"
	"time"

	"hosim-go/internal/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase membuat koneksi database sesuai driver yang ditentukan di konfigurasi
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// Tentukan log level GORM: Info saat dev (agar query SQL terlihat), Silent/Warn saat prod
	gormLogLevel := logger.Info
	if cfg.IsProduction() {
		gormLogLevel = logger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	}

	switch cfg.DBDriver {
	case "sqlite":
		// Driver Pure Go SQLite (tanpa CGO di Windows)
		// Mengaktifkan foreign keys, WAL journal mode (untuk concurrency), dan busy timeout
		dsn := fmt.Sprintf("%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", cfg.DBName)
		dialector = sqlite.Open(dsn)
		log.Printf("[DATABASE] Menggunakan SQLite: %s\n", cfg.DBName)

	case "postgres":
		// Driver PostgreSQL untuk staging / production
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		dialector = postgres.Open(dsn)
		log.Printf("[DATABASE] Menggunakan PostgreSQL di %s:%s/%s\n", cfg.DBHost, cfg.DBPort, cfg.DBName)

	default:
		return nil, fmt.Errorf("driver database '%s' tidak didukung (gunakan 'sqlite' atau 'postgres')", cfg.DBDriver)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	// Konfigurasi Connection Pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengakses generic database pool: %w", err)
	}

	if cfg.DBDriver == "sqlite" {
		// SQLite paling aman dengan 1 max open conn untuk menghindari 'database is locked' saat write
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
	} else {
		// PostgreSQL untuk skala production rumah sakit
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(1 * time.Hour)
	}

	return db, nil
}
