# HOSIM Web — Frontend Rekam Medis & SIMRS (Svelte 5 + TypeScript)

Aplikasi web front-end untuk sistem informasi manajemen rumah sakit (SIMRS) dan rekam medis elektronik (EHR / RME) **HOSIM-GO**, dibangun menggunakan **Svelte 5 Runes**, **TypeScript (Strict)**, **Vite**, dan **Tailwind CSS v4**.

---

## 🎨 Prinsip Desain & Ergonomi Medis

Antarmuka dirancang dengan standar desain **Google Workspace / Material Design 3** yang disesuaikan khusus untuk kenyamanan dan efisiensi tenaga medis:

- **Estetika Bersih (Clean Canvas)**: Menggunakan palet Google Blue (`#0b57d0`), latar lembut (`#f8fafd`), dan White Island Card (`rounded-3xl border border-[#e1e5ea]`).
- **Desktop Launcher (Root Hub `#/desktop`)**:
  - Halaman beranda terpusat yang menampilkan katalog modul SIMRS secara modular.
  - Pengelompokan kategori:
    1. *Pelayanan Pasien & Asuhan Medis* (Klinis EHR, Rawat Jalan, IGD, Rawat Inap).
    2. *Penunjang Medis & Logistik* (Farmasi, Inventori & Logistik RS).
    3. *Administrasi, Keuangan & Tata Kelola* (Master Data RS, Akuntansi RS, Billing Kasir, Akses & Pengguna RBAC, Audit Trail).
  - Dilengkapi pencarian instan (live search), badge status modul (*Aktif & Terintegrasi* vs *Tahap Integrasi*), statistik ringkas, serta filter hak akses (*permission-aware*).
- **Contextual Focused Navigation Rail**:
  - Sidebar bernavigasi kontekstual yang otomatis menyesuaikan menu dengan modul aktif (Klinis, Master Data, Akuntansi, atau Auth).
  - Mode **Compact Rail (72px)** secara default agar area grafik tubuh, formulir SOAP, dan tabel rekam medis mendapatkan bidang kerja maksimal.
  - Mode **Expanded (Hamburger Toggle)** untuk melihat struktur menu secara mendalam.
  - Tombol cepat *Kembali ke Desktop Launcher*.
  - **Dual Sub-Menu Interaction**:
    - Saat compact (72px): Sub-menu muncul sebagai **Floating Flyout Popover** saat di-hover/klik.
    - Saat expanded: Sub-menu tampil sebagai **Accordion Collapsible Tree**.
- **TopBar & Header Terintegrasi**:
  - *Module Context Switcher & Breadcrumbs*: Menampilkan ikon modul aktif dan tombol kembali ke launcher.
  - *Live Backend Health Ping*: Monitor status koneksi ke Go Gin backend (`/api/v1/health`) secara otomatis tiap 15 detik dengan indikator WiFi visual.
  - *User Profile & Session Dropdown*: Kartu profil pengguna, role aktif, dan aksi logout dengan deteksi klik luar (*click-outside listener*) untuk penutupan otomatis.
- **Tanpa AI Slop**: Tidak ada efek neon glow berlebihan atau visual mengganggu; mengutamakan densitas informasi, kejelasan tipografi klinis, dan kecepatan respon.

---

## 📁 Struktur Direktori Frontend (`web/src/`)

Arsitektur kode menerapkan pola **Domain-Driven Feature Folders** dan **Composables API** yang terisolasi serta konsisten dengan modul backend Go (`internal/<domain>`):

```text
web/src/
├── lib/
│   ├── api/                        # Modular HTTP Composables & API Clients
│   │   ├── client.ts               # Base fetcher, Token storage, Mutex Refresh Token, Health check
│   │   ├── accounting/             # API client akuntansi (coa.ts)
│   │   ├── auth/                   # API client autentikasi (session, users, roles, permissions)
│   │   ├── organization/           # API client organisasi (departments, service_units)
│   │   └── index.ts                # Central API barrel export
│   │
│   ├── components/                 # Komponen antarmuka (Reusable UI & App Shell)
│   │   ├── m3/                     # Atomic Material 3 Primitives
│   │   │   ├── M3Button.svelte     # Tombol M3 (filled, tonal, outlined, drive-new) dengan Snippet
│   │   │   ├── M3Card.svelte       # Kartu kontainer M3
│   │   │   ├── M3Chip.svelte       # Chip status & filter
│   │   │   └── M3SegmentedButton.svelte # Toggle segmen view interaktif
│   │   ├── BodyDiagram.svelte      # Diagram anatomi tubuh interaktif (Anterior & Posterior)
│   │   ├── LoginPage.svelte        # Halaman autentikasi Material You
│   │   ├── NavigationRail.svelte   # Contextual side rail navigasi + dual sub-menu
│   │   ├── PatientBanner.svelte    # Banner ringkasan identitas & tanda vital pasien
│   │   └── TopBar.svelte           # Header Google Drive style, health ping & profile dropdown
│   │
│   ├── pages/                      # Modul Halaman Berbasis Domain
│   │   ├── launcher/               # Root Hub Navigasi SIMRS
│   │   │   ├── DesktopLauncher.svelte   # Grid modul interaktif, pencarian, dan peluncur kerja
│   │   │   └── ModulePlaceholder.svelte # Pratinjau modul roadmap & integrasi endpoint
│   │   │
│   │   ├── clinical/               # Modul Rekam Medis & Pelayanan Klinis
│   │   │   ├── PhysicalExamPage.svelte  # Status Lokalis & Body Diagram + JSON Inspector
│   │   │   ├── SoapPage.svelte          # Formulir Anamnesis & Resume SOAP (ICD-10)
│   │   │   └── index.ts                 # Barrel export modul klinis
│   │   │
│   │   ├── master/                 # Modul Master Data Rumah Sakit (1-to-1 dengan internal/master/)
│   │   │   ├── patient/            # Domain Pasien (IndexPage.svelte)
│   │   │   ├── practitioner/       # Domain Dokter & Nakes (IndexPage.svelte)
│   │   │   ├── departement/        # Domain Instalasi / Departemen (IndexPage.svelte)
│   │   │   ├── service_unit/       # Domain Poliklinik & Unit Layanan (IndexPage.svelte)
│   │   │   ├── room/               # Domain Ruangan, Kamar, & Tempat Tidur (IndexPage.svelte)
│   │   │   ├── customer/           # Domain Debitur/Customer BPJS & Asuransi (IndexPage.svelte)
│   │   │   ├── referal/            # Domain Faskes Asal & Rujukan (IndexPage.svelte)
│   │   │   ├── tariff_class/       # Domain Kelas Tarif Layanan & KRIS (IndexPage.svelte)
│   │   │   └── index.ts            # Barrel export modul master data
│   │   │
│   │   ├── accounting/             # Modul Akuntansi & Keuangan RS
│   │   │   ├── coa/                # Bagan Akun Standar / Chart of Accounts (IndexPage.svelte)
│   │   │   └── index.ts            # Barrel export modul accounting
│   │   │
│   │   └── auth/                   # Modul Manajemen Akses & Keamanan
│   │       ├── IndexPage.svelte    # Tab pengelolaan Pengguna, Peran (Roles), & Izin (Permissions)
│   │       └── index.ts            # Barrel export modul auth
│   │
│   ├── stores/                     # State Management Universal Svelte 5
│   │   └── auth.svelte.ts          # Reactive Auth Store via $state, token decode & permission helpers
│   │
│   ├── types/                      # Shared TypeScript Interfaces per Domain
│   │   ├── common.ts               # BaseEntity, ApiResponse<T>, PaginatedResult<T>
│   │   ├── auth.ts                 # UserProfile, AuthSession, LoginPayload, Role, Permission
│   │   ├── clinical.ts             # BodyFinding, ClinicalNotes, FindingCategory
│   │   ├── accounting/             # Account types (coa, klasifikasi akun)
│   │   ├── master/                 # Domain Master Types (patient, room, customer, referal, dll.)
│   │   └── index.ts                # Root types barrel export
│   │
│   └── router.ts                   # Hash-based SPA Router & Module Context Resolver
│
├── App.svelte                      # Main App Shell / Layout Container & Route Switcher
├── app.css                         # Tailwind CSS v4 design tokens & base resets
├── main.ts                         # Application entrypoint
├── package.json                    # Dependencies & scripts
├── tsconfig.json                   # TypeScript configuration ($lib path alias)
└── vite.config.js                  # Vite configuration & proxy API ke localhost:8080
```

---

## ⚡ Implementasi Svelte 5 Best Practices

Aplikasi ini dibangun menggunakan arsitektur resmi **Svelte 5 Runes**:

1. **Reactivity Berbasis Runes**:
   - `$state`: Menyimpan reaktivitas data lokal, filter pencarian tabel, dan form input secara granular.
   - `$derived`: Kalkulasi otomatis untuk pemfilteran modul launcher, pemfilteran tabel master data, titik temuan tubuh, dan penentuan modul aktif dari rute hash.
   - `$props` & `$bindable`: Pengikatan data dua arah antar komponen induk-anak secara eksplisit tanpa sintaks `export let`.
   - `$effect`: Efek samping deklaratif untuk sinkronisasi URL hash, pembaruan rute, dan state transisi.

2. **Universal Reactivity (`.svelte.ts`)**:
   - Toko state global [`auth.svelte.ts`](file:///c:/laragon/www/hosim-go/web/src/lib/stores/auth.svelte.ts) menggunakan fitur `.svelte.ts` murni dengan `$state` dan method getter/helper (`isAuthenticated`, `hasPermission`, `hasAnyRole`). Menghilangkan kebutuhan `writable()` atau subscribe `$` manual warisan Svelte 4.

3. **Mutex & Deduplicated JWT Refresh**:
   - Modul [`client.ts`](file:///c:/laragon/www/hosim-go/web/src/lib/api/client.ts) dilengkapi penanganan pembaharuan access token otomatis saat menerima respons HTTP `401 Unauthorized`.
   - Menggunakan promise singleton (`refreshPromise`) untuk mencegah balapan konkurensi (race condition) ketika banyak request API dikirimkan bersamaan.

4. **Snippets Modern (`Snippet` & `{@render}`)**:
   - Menggantikan elemen `<slot />` dengan potongan template yang ter-type-check secara ketat pada komponen dasar Material 3.

5. **Deklaratif Window Event Listener**:
   - Menggunakan `<svelte:window onhashchange={...} onclick={...} />` untuk sinkronisasi tombol navigasi browser serta auto-close popover/dropdown tanpa memory leak.

6. **Path Aliasing `$lib`**:
   - Seluruh impor komponen, type, dan modul internal menggunakan alias `$lib` (contoh: `import { auth } from '$lib/stores/auth.svelte'`).

---

## 🧭 Sistem Routing SPA (Hash-Based Persistence)

Routing dikelola oleh modul [`router.ts`](file:///c:/laragon/www/hosim-go/web/src/lib/router.ts) dengan dukungan URL hash sinkron dan deteksi modul otomatis:

| Nav ID | URL Hash Browser | Deskripsi Modul & Halaman |
| :--- | :--- | :--- |
| `desktop` | `#/desktop` | **Desktop Launcher (Beranda Utama Hub SIMRS)** |
| `physical` | `#/clinical/physical` | Rekam Medis: Status Lokalis & Body Diagram Anatomi |
| `anamnesis` | `#/clinical/anamnesis` | Rekam Medis: Catatan Anamnesis & Resume SOAP |
| `master-patient` | `#/master/patient` | Master: Database Pasien & Rekam Medis Demografi |
| `master-practitioner` | `#/master/practitioner` | Master: Tenaga Medis, Dokter, & Nakes |
| `master-departement` | `#/master/departement` | Master: Instalasi & Departemen Rumah Sakit |
| `master-service-unit` | `#/master/service-unit` | Master: Poliklinik & Unit Layanan Pelayanan |
| `master-room` | `#/master/room` | Master: Ruangan, Kamar Inap, & Tempat Tidur (Bed) |
| `master-customer` | `#/master/customer` | Master: Debitur, Penjamin BPJS & Asuransi Swasta |
| `master-referal` | `#/master/referal` | Master: Faskes Asal & Tujuan Rujukan Pasien |
| `master-tariff-class` | `#/master/tariff-class` | Master: Kelas Tarif Layanan & Standar KRIS |
| `accounting-coa` | `#/accounting/coa` | Akuntansi: Bagan Akun Standar (Chart of Accounts) |
| `accounting-journals` | `#/accounting/journals` | Akuntansi: Jurnal Umum & Transaksi Memorial *(Preview)* |
| `accounting-ledger` | `#/accounting/ledger` | Akuntansi: Buku Besar & Kas/Bank Treasury *(Preview)* |
| `accounting-reports` | `#/accounting/reports` | Akuntansi: Laporan Keuangan & Neraca Saldo *(Preview)* |
| `auth-users` | `#/auth/users` | Akses: Manajemen Pengguna & Akun Nakes |
| `auth-roles` | `#/auth/roles` | Akses: Manajemen Peran (Roles) & Penugasan Izin |
| `auth-permissions` | `#/auth/permissions` | Akses: Katalog Hak Akses (Permissions Matrix) |
| `outpatient-workspace` | `#/outpatient/workspace` | Pelayanan: Rawat Jalan Poliklinik *(Preview Roadmap)* |
| `emergency-workspace` | `#/emergency/workspace` | Pelayanan: Triase & IGD Cito *(Preview Roadmap)* |
| `inpatient-workspace` | `#/inpatient/workspace` | Pelayanan: Rawat Inap & Admisi Bangsal *(Preview Roadmap)* |
| `pharmacy-workspace` | `#/pharmacy/workspace` | Penunjang: E-Resep & Farmasi Cito *(Preview Roadmap)* |
| `inventory-workspace` | `#/inventory/workspace` | Logistik: Gudang Farmasi & Depo Logistik *(Preview Roadmap)* |
| `billing-workspace` | `#/billing/workspace` | Keuangan: Kasir & Billing Pasien *(Preview Roadmap)* |
| `audit-workspace` | `#/audit/workspace` | Tata Kelola: Append-Only Audit Trail *(Preview Roadmap)* |
| `history` | `#/activity/history` | Riwayat Kunjungan RME Pasien |
| `schedule` | `#/activity/schedule` | Jadwal Kontrol Poliklinik Pasien |

> [!TIP]
> **Keunggulan Hash-Based Router**: Saat browser di-refresh (`F5`), dibuka lewat tautan langsung, atau menggunakan tombol Back/Forward browser, halaman **tidak akan reset ke index** dan tidak memerlukan konfigurasi URL rewrite khusus pada web server (Apache/Nginx/Laragon).

---

## 🚀 Skrip Pengembangan & Validasi

Jalankan perintah berikut di dalam direktori `web/`:

```bash
# Menjalankan development server (Vite HMR)
npm run dev

# Memeriksa diagnostik TypeScript & Svelte (Type-Check)
npx svelte-check --tsconfig ./tsconfig.json

# Membangun bundle produksi
npm run build

# Menjalankan preview dari build produksi
npm run preview
```

---

## 🔗 Integrasi Backend (Vite Proxy & Autentikasi)

File `vite.config.js` telah dikonfigurasi untuk mem-proxy seluruh panggilan rute `/api` ke Go Gin Backend:

```javascript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}
```

Ketika backend Go aktif di port `8080`, frontend akan berkomunikasi secara live dengan API backend (termasuk validasi JWT dan RBAC). Apabila backend offline, sistem menyediakan mode **Fallback Demo** otomatis menggunakan kredensial bawaan:
- **Dokter**: `dokter` / `dokter123`
- **Admin**: `admin` / `admin123`

