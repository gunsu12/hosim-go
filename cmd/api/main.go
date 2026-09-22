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

	"hosim-go/internal/auth"
	"hosim-go/internal/config"
	"hosim-go/internal/database"
	"hosim-go/internal/master/departement"
	"hosim-go/internal/master/patient"
	"hosim-go/internal/master/payer"
	"hosim-go/internal/master/practitioner"
	"hosim-go/internal/master/referal"
	"hosim-go/internal/master/room"
	serviceunit "hosim-go/internal/master/service_unit"
	tariffclass "hosim-go/internal/master/tariff_class"
	"hosim-go/internal/middleware"
	"hosim-go/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Konfigurasi
	cfg := config.LoadConfig()

	// Peringatan keamanan jika berjalan di production dengan default secret
	if cfg.IsProduction() && cfg.JWTSecret == "hosim-secret-key-change-this-in-production-12345678" {
		log.Println("[WARN] PERINGATAN KEAMANAN: Aplikasi berjalan pada mode production tetapi JWT_SECRET masih menggunakan nilai default!")
	}

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

	// Auto-Migrate tabel-tabel Master
	err = db.AutoMigrate(
		// Master Penjamin / Payer
		&payer.PayerType{},
		&payer.Payer{},

		// Master Departemen / Instalasi, Unit Layanan & Ruangan
		&departement.Departement{},
		&serviceunit.ServiceUnit{},
		&room.Room{},

		// Master Rujukan & Kelas Tarif
		&referal.Referal{},
		&tariffclass.TariffClass{},

		// Master Pasien
		&patient.Patient{},
		&patient.PatientEmergencyContact{},
		&patient.PatientRelation{},
		&patient.PatientAllergy{},
		&patient.PatientAddress{},
		&patient.PatientDrugHistory{},
		&patient.PatientChronicalDisease{},

		// Master Nakes / Practitioner
		&practitioner.Profession{},
		&practitioner.Specialty{},
		&practitioner.Practitioner{},

		// Autentikasi, Role & Hak Akses (RBAC)
		&auth.Permission{},
		&auth.Role{},
		&auth.User{},
		&auth.RefreshToken{},
	)
	if err != nil {
		log.Fatalf("[FATAL] AutoMigrate database gagal: %v\n", err)
	}
	log.Println("[DATABASE] Skema tabel database berhasil dimigrasi.")

	// Jalankan Seeder data referensi (Profesi & Spesialisasi SatuSehat)
	if err := practitioner.SeedProfessionsAndSpecialties(db); err != nil {
		log.Printf("[WARN] Seeder profesi dan spesialisasi gagal: %v\n", err)
	}
	if err := practitioner.SeedDefaultPractitioner(db); err != nil {
		log.Printf("[WARN] Seeder data dokter gagal: %v\n", err)
	}

	// Jalankan Seeder Roles & Permissions RBAC (termasuk default akun admin)
	if err := auth.SeedRBAC(db); err != nil {
		log.Printf("[WARN] Seeder RBAC gagal: %v\n", err)
	}

	// 4. Inisialisasi Modul Autentikasi & Master Data (Dependency Injection)
	// Modul Autentikasi
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg)
	authHandler := auth.NewHandler(authService, cfg.JWTSecret)

	// Master Departemen
	departementRepo := departement.NewRepository(db)
	departementService := departement.NewService(departementRepo)
	departementHandler := departement.NewHandler(departementService)

	// Master Penjamin / Payer
	payerRepo := payer.NewRepository(db)
	payerService := payer.NewService(payerRepo)
	payerHandler := payer.NewHandler(payerService)

	// Master Rujukan / Referal
	referalRepo := referal.NewRepository(db)
	referalService := referal.NewService(referalRepo)
	referalHandler := referal.NewHandler(referalService)

	// Master Unit Layanan / Service Unit
	serviceUnitRepo := serviceunit.NewRepository(db)
	serviceUnitService := serviceunit.NewService(serviceUnitRepo)
	serviceUnitHandler := serviceunit.NewHandler(serviceUnitService)

	// Master Ruangan / Room
	roomRepo := room.NewRepository(db)
	roomService := room.NewService(roomRepo)
	roomHandler := room.NewHandler(roomService)

	// Master Kelas Tarif / Tariff Class
	tariffClassRepo := tariffclass.NewRepository(db)
	tariffClassService := tariffclass.NewService(tariffClassRepo)
	tariffClassHandler := tariffclass.NewHandler(tariffClassService)

	// Master Pasien
	patientRepo := patient.NewRepository(db)
	patientService := patient.NewService(patientRepo)
	patientHandler := patient.NewHandler(patientService)

	// Master Tenaga Medis (Practitioner)
	practitionerRepo := practitioner.NewRepository(db)
	practitionerService := practitioner.NewService(practitionerRepo)
	practitionerHandler := practitioner.NewHandler(practitionerService)

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

		// Daftarkan route Autentikasi (/api/v1/auth/*)
		authHandler.RegisterRoutes(v1)

		// Rute Master Data Terproteksi (Memerlukan JWT Access Token)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			departementHandler.RegisterRoutes(protected)
			payerHandler.RegisterRoutes(protected)
			referalHandler.RegisterRoutes(protected)
			serviceUnitHandler.RegisterRoutes(protected)
			roomHandler.RegisterRoutes(protected)
			tariffClassHandler.RegisterRoutes(protected)
			patientHandler.RegisterRoutes(protected)
			practitionerHandler.RegisterRoutes(protected)
		}
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
