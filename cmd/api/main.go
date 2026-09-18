package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hosim-go/internal/config"
	"hosim-go/internal/database"
	"hosim-go/internal/master/patient"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Konfigurasi
	cfg := config.LoadConfig()

	// 2. Set mode Gin (Release untuk production, Debug untuk development)
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 3. Inisialisasi Koneksi Database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Inisialisasi database gagal: %v\n", err)
	}

	// Auto-Migrate tabel-tabel Master Pasien
	err = db.AutoMigrate(
		&patient.Patient{},
		&patient.PatientEmergencyContact{},
		&patient.PatientRelation{},
		&patient.PatientAllergy{},
		&patient.PatientAddress{},
		&patient.PatientDrugHistory{},
		&patient.PatientChronicalDisease{},
	)
	if err != nil {
		log.Fatalf("[FATAL] AutoMigrate database gagal: %v\n", err)
	}
	log.Println("[DATABASE] Skema tabel Master Pasien berhasil dimigrasi.")

	// 4. Inisialisasi Modul Master Pasien (Dependency Injection)
	patientRepo := patient.NewRepository(db)
	patientService := patient.NewService(patientRepo)
	patientHandler := patient.NewHandler(patientService)

	// 5. Inisialisasi Router Gin
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Middleware CORS sederhana agar frontend (Svelte) dapat mengakses API
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 6. Health Check Endpoint
	healthHandler := func(c *gin.Context) {
		sqlDB, err := db.DB()
		dbStatus := "connected"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "disconnected"
		}

		response.Success(c, http.StatusOK, "HOSIM-GO Backend API is running", gin.H{
			"app_name":  cfg.AppName,
			"env":       cfg.AppEnv,
			"database":  dbStatus,
			"driver":    cfg.DBDriver,
			"timestamp": time.Now().Format(time.RFC3339),
		})
	}

	r.GET("/", healthHandler)
	r.GET("/health", healthHandler)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthHandler)
		// Daftarkan route modul Master Pasien
		patientHandler.RegisterRoutes(v1)
	}

	// 7. Jalankan Server dengan Graceful Shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Jalankan server di background goroutine
	go func() {
		log.Printf("[SERVER] %s berjalan di port http://localhost:%s\n", cfg.AppName, cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server error: %v\n", err)
		}
	}()

	// 8. Menunggu Sinyal Shutdown (SIGINT atau SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[SHUTDOWN] Mematikan server secara graceful...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v\n", err)
	}

	// Tutup koneksi database
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
		log.Println("[DATABASE] Koneksi database ditutup.")
	}

	log.Println("[SHUTDOWN] Server HOSIM-GO berhenti dengan aman.")
}
