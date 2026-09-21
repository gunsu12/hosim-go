package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config menampung semua konfigurasi environment aplikasi
type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DBDriver   string // "sqlite" atau "postgres"
	DBName     string // Nama file untuk SQLite, atau nama database untuk Postgres
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
}

// LoadConfig membaca environment variables dari .env (jika ada) dan OS env
func LoadConfig() *Config {
	// Cari .env di direktori kerja saat ini atau naik ke direktori induk (penting saat unit test berjalan di subfolder)
	dir, err := os.Getwd()
	if err == nil {
		for {
			envPath := filepath.Join(dir, ".env")
			if _, err := os.Stat(envPath); err == nil {
				_ = godotenv.Load(envPath)
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	return &Config{
		AppName:    getEnv("APP_NAME", "HOSIM-GO"),
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBName:     getEnv("DB_NAME", "hosim_go"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// IsProduction mengecek apakah aplikasi berjalan di environment produksi
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
