<p align="center">
  <img src="assets/banner.png" alt="HOSIM-GO Banner" width="75%" />
</p>

<h2 align="center">HOSIM-GO (Hospital Information Management System in Go)</h2>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Svelte-5_Runes-FF3E00?style=flat-square&logo=svelte&logoColor=white" alt="Svelte 5" />
  <img src="https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=flat-square&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Tailwind_CSS-v4-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white" alt="Tailwind CSS" />
  <img src="https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat-square&logo=typescript&logoColor=white" alt="TypeScript" />
  <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License MIT" />
</p>

HOSIM-GO adalah sistem informasi manajemen rumah sakit (SIMRS) dan rekam medis elektronik (EHR / RME) modern berbasis **Go (Gin)** dan **Svelte 5 (Runes)**. Didesain untuk performa tinggi, keandalan operasional, dan kepatuhan standar interoperabilitas kesehatan (Kemenkes SATUSEHAT & BPJS Kesehatan).

> [!NOTE]
> ### ⚠️ Pembelajaran & Kolaborasi
> Proyek ini dikembangkan secara mandiri untuk kebutuhan pembelajaran dan eksplorasi teknologi SIMRS modern.
> Diskusi, masukan, dan kontribusi sangat terbuka. Hubungi: **gunawansuarna@gmail.com**.

---

## 🛠️ Tech Stack

| Layer | Label & Komponen Utama |
| :--- | :--- |
| **Backend** | ![Go](https://img.shields.io/badge/Go_1.22+-00ADD8?style=flat-square&logo=go&logoColor=white) ![Gin](https://img.shields.io/badge/Gin_Web_Framework-008ECF?style=flat-square&logo=gin&logoColor=white) ![GORM](https://img.shields.io/badge/GORM-7952B3?style=flat-square) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white) ![Goose](https://img.shields.io/badge/Goose_Migrations-black?style=flat-square) ![JWT](https://img.shields.io/badge/JWT_Auth-000000?style=flat-square&logo=jsonwebtokens&logoColor=white) |
| **Frontend** | ![Svelte 5](https://img.shields.io/badge/Svelte_5_(Runes)-FF3E00?style=flat-square&logo=svelte&logoColor=white) ![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=flat-square&logo=typescript&logoColor=white) ![Tailwind CSS v4](https://img.shields.io/badge/Tailwind_CSS_v4-06B6D4?style=flat-square&logo=tailwindcss&logoColor=white) ![Vite](https://img.shields.io/badge/Vite-646CFF?style=flat-square&logo=vite&logoColor=white) ![Lucide Icons](https://img.shields.io/badge/Lucide_Icons-F56565?style=flat-square) |
| **Interoperabilitas** | ![SATUSEHAT](https://img.shields.io/badge/Kemenkes-SATUSEHAT_(HL7_FHIR)-009B4D?style=flat-square) ![BPJS](https://img.shields.io/badge/BPJS_Kesehatan-VClaim-2E7D32?style=flat-square) |

---

## ✨ Fitur Utama

- **Rekam Medis Elektronik (RME / EHR)**:
  - Interactive Body Diagram (Status Lokalis) berbasis SVG dengan pemetaan cedera/nyeri NRS.
  - Catatan Medis Terstruktur (SOAP & ICD-10/ICD-9).
  - Odontogram interaktif (32 gigi FDI standard).
  - Order Penunjang (Laboratorium, Radiologi, E-Resep).
- **Alur Pelayanan & Care Settings**:
  - Gawat Darurat (IGD), Rawat Jalan (Poliklinik), dan Rawat Inap (Bangsal/Kamar).
- **Master Data Rumah Sakit**:
  - Pasien, Tenaga Medis (Dokter/Perawat), Departemen, Unit Layanan, Ruangan/Bed, Penjamin/Asuransi, Faskes Rujukan, dan Tarif.
- **Keamanan & RBAC**:
  - Otentikasi berbasis JWT, role-based access control, dan audit log aktivitas.

---

## 🚀 Memulai (Quick Start)

### Prasyarat
- **Go** v1.22+
- **Node.js** v18+ atau v20+ & **npm**
- **PostgreSQL** aktif

### 1. Backend (Go API)

1. Salin dan sesuaikan konfigurasi database di `.env`:
   ```bash
   cp .env.example .env
   ```
2. Jalankan migrasi dan seeder awal:
   ```bash
   go run cmd/migrate/main.go up
   go run cmd/seed/main.go
   ```
3. Jalankan server backend:
   ```bash
   go run cmd/api/main.go
   ```
   API berjalan di `http://localhost:8080`. Health check: `http://localhost:8080/health`.

### 2. Frontend (Svelte 5 App)

1. Masuk ke direktori web dan instal dependensi:
   ```bash
   cd web
   npm install
   ```
2. Jalankan development server:
   ```bash
   npm run dev
   ```
   Buka browser di `http://127.0.0.1:5173`.
   *(Preset login cepat tersedia di halaman login: Dokter / Admin)*.

---

## 📖 Panduan Teknis & Arsitektur

Detail arsitektur Clean Architecture, aturan modular monolith, konvensi penulisan kode, boundary transaksi, dan panduan engineering selengkapnya dapat dibaca di:
👉 **[AGENTS.md](AGENTS.md)**

---

## 📜 Lisensi

Didistribusikan di bawah lisensi [MIT](LICENSE).
