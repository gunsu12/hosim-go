# HOSIM-GO (Hospital Information Management System in Go)

HOSIM-GO adalah sistem informasi manajemen rumah sakit (SIMRS) open-source yang dibangun dengan bahasa **Go** menggunakan arsitektur **Clean Architecture & Modular Monolith**.

Didesain untuk keandalan tinggi (24/7), integritas data transaksi medis yang ketat, dan siap diintegrasikan dengan standar nasional (BPJS V-Claim & Kemenkes SatuSehat HL7 FHIR).

---

## 🏗️ Arsitektur & Teknologi

- **Backend**: Go (Gin Framework)
- **Database (Development)**: SQLite (Pure Go via `github.com/glebarez/sqlite`, zero CGO)
- **Database (Production Ready)**: PostgreSQL via GORM
- **Frontend (Rencana)**: Svelte (Monorepo di folder `web/`, siap di-embed via `go:embed`)
- **Pola Arsitektur**: Modular Monolith & Clean Architecture (Domain -> Repository -> Service -> Handler)

---

## 📁 Struktur Direktori

```text
hosim-go/
├── cmd/
│   └── api/
│       └── main.go              # Entry point server Gin & Graceful Shutdown
├── internal/
│   ├── config/                  # Pengaturan environment (.env loader)
│   ├── database/                # Koneksi DB abstrak (SQLite / PostgreSQL)
│   └── master/                  # [Modul] Master Data (Pasien, Dokter, Poli, Kamar)
├── pkg/
│   └── response/                # Helper standard envelope response JSON
├── .env.example                 # Template konfigurasi environment
├── .gitignore                   # Ignore SQLite db, bin, dan temporary files
├── go.mod                       # Modul Go & dependencies
└── README.md
```

---

## 🚀 Cara Menjalankan

### 1. Prasyarat
- Go 1.22 atau lebih baru.

### 2. Konfigurasi Lingkungan
Secara default, file `.env` sudah disiapkan untuk menggunakan SQLite:
```env
APP_NAME=HOSIM-GO
APP_ENV=development
APP_PORT=8080
DB_DRIVER=sqlite
DB_NAME=hosim.db
```

### 3. Menjalankan Aplikasi
```bash
go run cmd/api/main.go
```
Buka browser atau API client (Postman/curl):
- `http://localhost:8080/health`
- `http://localhost:8080/api/v1/health`

---

## 📜 Lisensi
MIT License - Proyek Open Source untuk kemajuan digitalisasi layanan kesehatan.
