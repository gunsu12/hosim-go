# HOSIM-GO (Hospital Information Management System in Go)

HOSIM-GO adalah sistem informasi manajemen rumah sakit (SIMRS) dan rekam medis elektronik (EHR / RME) modern yang dibangun menggunakan arsitektur **Clean Architecture & Modular Monolith** di backend (**Go Gin**) dan antarmuka web modern berbasis **Svelte 5 Runes & TypeScript**.

Didesain untuk performa tinggi, keandalan 24/7, integritas data transaksi medis yang ketat, serta kemudahan interoperabilitas (Kemenkes SatuSehat HL7 FHIR & BPJS Kesehatan VClaim).

> [!NOTE]
> ### ⚠️ ALERT & DISCLAIMER
> Project ini adalah project yang rencananya akan saya kembangkan untuk kebutuhan pembelajaran pribadi saya. Segala kekurangan dan ketidaksempurnaan project ini mohon dimaklumi.
> 
> Kalau ada yang berminat untuk berdiskusi atau menjadi kontributor, silakan hubungi email: **gunawansuarna@gmail.com**

---

## 🏗️ Arsitektur & Teknologi

### Backend (Go Gin)
- **Bahasa & Framework**: Go 1.22+ dengan [Gin Web Framework](https://github.com/gin-gonic/gin).
- **Database**: PostgreSQL via [GORM](https://gorm.io/) & [Goose v3](https://github.com/pressly/goose) untuk skema migrasi.
- **Autentikasi**: JWT (JSON Web Token) dengan stateless authorization & role-based permissions.
- **Pola Arsitektur**: Clean Architecture & Domain-Driven Modular Monolith (`Entity` -> `Repository` -> `Service` -> `Handler`).

### Frontend (`web/`)
- **Framework**: [Svelte 5](https://svelte.dev/) dengan paradigma **Runes** (`$state`, `$derived`, `$props`, `$bindable`, `$effect`).
- **Bahasa & Build Tool**: TypeScript (Strict Mode) + Vite 8.
- **Styling**: Tailwind CSS v4 dengan estetika desain **Google Workspace / Material Design 3** (*clean, compact navigation rail, white content island, zero AI slop*).
- **Routing**: Lightweight Hash-based SPA Router (`#/clinical/physical`, `#/master/patient`) yang persisten terhadap page refresh (`F5`).
- **Icons**: Official `@lucide/svelte` v1.47+.

---

## 🌟 Fitur Utama yang Sudah Tersedia

### 1. Rekam Medis Elektronik (RME) & Pelayanan Klinis
- **Interactive Body Diagram (Status Lokalis)**: Pemetaan visual titik keluhan/cedera pasien pada diagram anatomi tubuh SVG (Tampak Depan & Belakang) dengan koordinat persentase akurat, visualisasi skala nyeri NRS (0–10), kategori luka robek, fraktur, hematoma, dan inspektur payload JSON API cito.
- **Catatan Anamnesis SOAP**: Pencatatan terstruktur *Subjective*, *Objective*, *Assessment* (ICD-10), dan *Plan* (Disposisi rawat inap / cito).
- **Odontogram Gigi 32 Gigi**: Pemetaan interaktif formula gigi FDI World Dental Federation (Rahang Atas Maxilla 18–28 & Rahang Bawah Mandibula 48–38).
- **Order Penunjang Medis**: Permintaan Laboratorium Cito, Pencitraan Radiologi, dan E-Prescription Farmasi.

### 2. Master Data Rumah Sakit (Tersinkronisasi per Domain Backend)
Setiap domain master data diisolasi ke dalam folder terpisah di backend (`internal/master/<domain>`) dan frontend (`web/src/lib/pages/master/<domain>/`):
1. **Pasien (`patient`)**: Pengelolaan identitas NIK/KTP, No. Rekam Medis (MRN), kontak, dan jenis penjamin.
2. **Tenaga Medis (`practitioner`)**: Database dokter spesialis, perawat, nomor SIP/STR, dan SATUSEHAT Practitioner ID.
3. **Departemen / Instalasi (`departement`)**: Master instalasi utama (Rawat Jalan, Rawat Inap, Bedah Sentral IBS, IGD, Lab).
4. **Unit Layanan (`service_unit`)**: Daftar poliklinik dan unit tujuan pendaftaran.
5. **Ruangan & Bed (`room`)**: Manajemen kapasitas tempat tidur, ruang operasi, bangsal rawat inap, dan ICU.
6. **Debitur / Penjamin (`customer`)**: Konfigurasi BPJS Kesehatan (PBI & Non-PBI), asuransi komersial, dan pasien umum.
7. **Faskes Rujukan (`referal`)**: Registrasi faskes pengirim rujukan (Puskesmas/Faskes 1, Klinik Pratama, RS Tipe B/C).
8. **Kelas Tarif (`tariff_class`)**: Klasifikasi akomodasi kamar (VVIP, VIP, KRIS Kelas 1-3) dan koefisien tarif layanan.

### 3. Otentikasi & Keamanan
- Halaman login Material You dengan preset profil cepat:
  - **Dokter**: `dokter` / `dokter123`
  - **Admin**: `admin` / `admin123`
- Session persistence via reactive store (`auth.svelte.ts`) dan indikator live ping status koneksi backend.

---

## 📁 Struktur Direktori Proyek

```text
hosim-go/
├── cmd/
│   ├── api/                     # Entry point Go Gin API & Graceful Shutdown
│   ├── migrate/                 # Tool migrasi database skema via Goose
│   └── seed/                    # Tool database seeder (RBAC, SatuSehat, Dokter)
├── internal/
│   ├── organization/            # Domain Organisasi & Fasilitas:
│   │   ├── department/          # Instalasi & Departemen
│   │   ├── serviceunit/         # Poliklinik & Unit Layanan
│   │   ├── room/                # Ruangan Perawatan & Tindakan
│   │   ├── referal/             # Faskes Mitra Rujukan
│   │   ├── bed/                 # Bed Management
│   │   ├── storage/             # Gudang Logistik/Farmasi Pusat
│   │   └── depo/                # Depo Farmasi Ruangan
│   ├── patient/                 # Domain Pasien & Rekam Medis
│   ├── practitioner/            # Domain Dokter, Nakes, & Spesialisasi
│   ├── catalog/                 # Katalog Medis (item, icd10, icd9, vitalsign)
│   ├── finance/                 # Keuangan & Penjamin (customer, tarifclass, tariff)
│   ├── clinical/                # Pelayanan Klinis EHR (encounter, observation, condition, procedure, dll)
│   ├── emergency/               # Alur Pelayanan Gawat Darurat (IGD)
│   ├── outpatient/              # Alur Pelayanan Rawat Jalan (Poliklinik)
│   ├── inpatient/               # Alur Pelayanan Rawat Inap (Bangsal)
│   ├── pharmacy/                # Farmasi (prescription, dispensing)
│   ├── inventory/               # Manajemen Stok & Logistik
│   ├── billing/                 # Kasir, Tagihan, & Invoice
│   ├── integration/             # Interoperabilitas (satusehat, bpjs, lis, pacs)
│   ├── audit/                   # Audit Trail & Jejak Aktivitas
│   ├── auth/                    # Modul Autentikasi JWT & RBAC
│   ├── config/                  # Pengaturan environment (.env loader)
│   ├── database/                # Koneksi database GORM & Seeder
│   └── middleware/              # Middleware JWT Auth, CORS, Logger
├── pkg/
│   ├── enums/                   # Enum tipe data medis & faskes
│   ├── jwt/                     # Helper generator & verifikasi token
│   └── response/                # Helper standard JSON response envelope
├── web/                         # Frontend Web App (Svelte 5 + TypeScript + Vite)
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/      # UI Components (TopBar, NavigationRail, BodyDiagram, M3)
│   │   │   ├── pages/
│   │   │   │   ├── clinical/    # Modul Halaman Klinis (PhysicalExam, SOAP, Odontogram, dll.)
│   │   │   │   └── master/      # Modul Halaman Master per domain backend
│   │   │   ├── stores/          # Universal Reactivity Stores (auth.svelte.ts)
│   │   │   ├── types/           # Shared domain TypeScript interfaces
│   │   │   ├── api.ts           # Axios / Fetch client ke Go Gin backend
│   │   │   └── router.ts        # Hash-based persistent SPA router
│   │   ├── App.svelte           # Main App Shell & Layout
│   │   └── main.ts              # Entry point Svelte 5
│   ├── package.json
│   ├── tsconfig.json            # Konfigurasi TypeScript ($lib alias)
│   └── vite.config.js           # Vite config & API reverse-proxy ke localhost:8080
├── .env.example                 # Template konfigurasi environment backend
├── go.mod                       # Modul Go & dependencies
└── README.md
```

---

## 🚀 Cara Menjalankan

### 1. Prasyarat
- **Go**: 1.22 atau lebih baru.
- **Node.js**: v18+ atau v20+ (disertai `npm`).

---

### 2. Menjalankan Backend (Go Gin)

1. Pastikan file `.env` sudah siap (menggunakan PostgreSQL):
   ```env
   APP_NAME=HOSIM-GO
   APP_ENV=development
   APP_PORT=8080
   DB_DRIVER=postgres
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=hosim_go
   DB_SSLMODE=disable
   JWT_SECRET=super-secret-jwt-key-hosim-2026
   ```

2. Jalankan migrasi skema database (Goose):
   ```bash
   go run cmd/migrate/main.go up
   ```

3. Jalankan seeder data awal (RBAC, superadmin, data referensi SatuSehat):
   ```bash
   go run cmd/seed/main.go
   ```

4. Jalankan server API backend:
   ```bash
   go run cmd/api/main.go
   ```
   Server backend akan aktif di: `http://localhost:8080`.
   Endpoint cek status: `http://localhost:8080/health` atau `http://localhost:8080/api/v1/health`.

---

### 3. Menjalankan Frontend (Svelte 5 Web App)

1. Buka terminal baru dan masuk ke direktori `web`:
   ```bash
   cd web
   ```

2. Instal dependensi:
   ```bash
   npm install
   ```

3. Jalankan Vite Development Server:
   ```bash
   npm run dev
   ```
   Aplikasi frontend akan aktif di: `http://127.0.0.1:5173`.
   *(Vite secara otomatis mem-proxy request `/api/*` ke Go Gin backend di port `8080`)*.

4. Buka browser di `http://127.0.0.1:5173` dan login menggunakan tombol **Dokter** atau **Admin**.

---

### 4. Validasi & Typecheck

Untuk memastikan kualitas kode frontend dan tidak ada kesalahan tipe data:
```bash
cd web
npx svelte-check --tsconfig ./tsconfig.json
npm run build
```

---

## 📜 Lisensi
MIT License - Proyek Open Source untuk kemajuan digitalisasi layanan kesehatan dan pembelajaran mandiri.
