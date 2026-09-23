package database

import (
	"fmt"
	"log"
	"time"

	"hosim-go/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase membuat koneksi database PostgreSQL sesuai konfigurasi
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	// Tentukan log level GORM: Info saat dev (agar query SQL terlihat), Silent/Warn saat prod
	gormLogLevel := logger.Info
	if cfg.IsProduction() {
		gormLogLevel = logger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	}

	// Driver PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	dialector := postgres.Open(dsn)
	log.Printf("[DATABASE] Menghubungkan ke PostgreSQL di %s:%s/%s\n", cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database PostgreSQL: %w", err)
	}

	// Konfigurasi Connection Pool untuk production rumah sakit
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengakses generic database pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}
