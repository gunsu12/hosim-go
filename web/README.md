# HOSIM Web — Frontend Rekam Medis & SIMRS (Svelte 5 + TypeScript)

Aplikasi web front-end untuk sistem informasi manajemen rumah sakit (SIMRS) dan rekam medis elektronik (EHR / RME) **HOSIM-GO**, dibangun menggunakan **Svelte 5 Runes**, **TypeScript (Strict)**, **Vite**, dan **Tailwind CSS v4**.

---

## 🎨 Prinsip Desain & Ergonomi Medis

Antarmuka dirancang dengan standar desain **Google Workspace / Material Design 3** yang disesuaikan khusus untuk lingkungan klinis rumah sakit:

- **Estetika Bersih (Clean Canvas)**: Menggunakan palet Google Blue (`#0b57d0`), latar lembut (`#f8fafd`), dan White Island Card (`rounded-3xl border border-[#e1e5ea]`).
- **Layar Lebih Lega (Spacious Real-Estate)**: Sidebar menu menggunakan pola **Compact Navigation Rail (72px)** secara default agar area grafik tubuh, formulir SOAP, dan tabel rekam medis mendapatkan bidang kerja maksimal.
- **Dual Sub-Menu Interaction**:
  - Saat compact (72px): Sub-menu muncul sebagai **Floating Flyout Popover** saat di-hover/klik.
  - Saat expanded (hamburger toggle): Sub-menu tampil sebagai **Accordion Collapsible Tree**.
- **Tanpa AI Slop**: Tidak ada efek neon glow berlebihan atau visual mengganggu; mengutamakan densitas informasi dan kecepatan respon.

---

## 📁 Struktur Direktori Frontend (`web/src/`)

Arsitektur kode menerapkan pola **Domain-Driven Feature Folders** yang terisolasi dan konsisten dengan modul backend Go (`internal/master/<domain>`):

```text
web/src/
├── lib/
│   ├── components/                 # Komponen antarmuka (Reusable UI & App Shell)
│   │   ├── m3/                     # Atomic Material 3 Primitives
│   │   │   ├── M3Button.svelte     # Tombol M3 (filled, tonal, outlined, drive-new) dengan Snippet
│   │   │   ├── M3Card.svelte       # Kartu kontainer M3
│   │   │   ├── M3Chip.svelte       # Chip status & filter
│   │   │   └── M3SegmentedButton.svelte # Toggle segmen view interaktif
│   │   ├── BodyDiagram.svelte      # Diagram anatomi tubuh interaktif (Anterior & Posterior)
│   │   ├── LoginPage.svelte        # Halaman autentikasi Material You
│   │   ├── NavigationRail.svelte   # Compact side rail navigasi + dual sub-menu
│   │   ├── PatientBanner.svelte    # Banner ringkasan identitas & tanda vital pasien
│   │   └── TopBar.svelte           # Header Google Drive style + live health ping
│   │
│   ├── pages/                      # Modul Halaman Berbasis Domain
│   │   ├── clinical/               # Modul Rekam Medis & Pelayanan Klinis
│   │   │   ├── PhysicalExamPage.svelte # Status Lokalis & Body Diagram + JSON Inspector
│   │   │   ├── SoapPage.svelte         # Formulir Anamnesis & Resume SOAP (ICD-10)
│   │   │   ├── OdontogramPage.svelte   # Pemetaan Gigi 32 Gigi FDI World Dental Federation
│   │   │   ├── LabPage.svelte          # Permintaan Cito Laboratorium Patologi
│   │   │   ├── RadiologyPage.svelte    # Pencitraan Radiologi & Diagnostik
│   │   │   ├── PharmacyPage.svelte     # E-Prescription Farmasi Cito
│   │   │   └── index.ts                # Barrel export terpusat modul klinis
│   │   │
│   │   └── master/                 # Modul Master Data Rumah Sakit (1-to-1 dengan internal/master/)
│   │       ├── patient/            # Domain Pasien (IndexPage.svelte)
│   │       ├── practitioner/       # Domain Dokter & Nakes (IndexPage.svelte)
│   │       ├── departement/        # Domain Instalasi / Departemen (IndexPage.svelte)
│   │       ├── service_unit/       # Domain Poliklinik & Unit Layanan (IndexPage.svelte)
│   │       ├── room/               # Domain Kamar Bedah, Rawat Inap, & Bed (IndexPage.svelte)
│   │       ├── customer/           # Domain Debitur/Customer BPJS & Asuransi (IndexPage.svelte)
│   │       ├── referal/            # Domain Faskes Rujukan Faskes 1/RS (IndexPage.svelte)
│   │       ├── tariff_class/       # Domain Kelas Tarif Layanan & KRIS (IndexPage.svelte)
│   │       └── index.ts            # Barrel export terpusat modul master data
│   │
│   ├── stores/                     # State Management Universal Svelte 5
│   │   └── auth.svelte.ts          # Reactive Auth Store via $state & reactive getters
│   │
│   ├── types/                      # Shared TypeScript Interfaces per Domain
│   │   ├── common.ts               # BaseEntity, ApiResponse<T>, PaginatedResult<T>
│   │   ├── auth.ts                 # UserProfile, AuthSession, LoginPayload
│   │   ├── clinical.ts             # BodyFinding, ClinicalNotes, FindingCategory
│   │   ├── master/                 # Domain Master Types (patient, room, customer, referal, dll.)
│   │   └── index.ts                # Root types barrel export
│   │
│   ├── api.ts                      # HTTP Client (Login, Backend Health Check)
│   └── router.ts                   # Hash-based SPA Router dengan sinkronisasi URL Bar
│
├── App.svelte                      # Main App Shell / Layout Container
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
   - `$state`: Menyimpan reaktivitas data lokal dan form secara granular.
   - `$derived`: Kalkulasi otomatis untuk pemfilteran data pencarian tabel dan filter titik temuan tubuh.
   - `$props` & `$bindable`: Pengikatan data dua arah antar komponen induk-anak secara eksplisit tanpa `export let`.
   - `$effect`: Efek samping deklaratif untuk sinkronisasi URL hash dan auto-expand accordion.

2. **Universal Reactivity (`.svelte.ts`)**:
   - Toko state global [`auth.svelte.ts`](file:///c:/laragon/www/hosim-go/web/src/lib/stores/auth.svelte.ts) menggunakan fitur `.svelte.ts` murni dengan `$state` dan getter function. Ini meniadakan kebutuhan `writable()` atau langganan `$` manual warisan Svelte 4.

3. **Snippets Modern (`Snippet` & `{@render}`)**:
   - Menggantikan elemen `<slot />` dengan potongan template yang ter-type-check secara ketat.

4. **Deklaratif Window Listener**:
   - Menggunakan `<svelte:window onhashchange={...} />` di [`App.svelte`](file:///c:/laragon/www/hosim-go/web/src/App.svelte) untuk otomatis menangani pembersihan memory leak saat unmount.

5. **Path Aliasing `$lib`**:
   - Seluruh impor komponen, type, dan modul internal menggunakan alias `$lib` (contoh: `import M3Button from '$lib/components/m3/M3Button.svelte'`).

---

## 🧭 Sistem Routing SPA (Hash-Based Persistence)

Menggunakan modul [`router.ts`](file:///c:/laragon/www/hosim-go/web/src/lib/router.ts) untuk mengelola rute aplikasi:

| Nav ID | URL Hash Browser | Deskripsi Modul |
| :--- | :--- | :--- |
| `physical` | `#/clinical/physical` | Pemeriksaan Fisik & Body Diagram Anatomi |
| `anamnesis` | `#/clinical/anamnesis` | Catatan Anamnesis & Resume SOAP |
| `odontogram` | `#/clinical/odontogram` | Pemeriksaan Odontogram Gigi 32 Gigi |
| `lab` | `#/diagnostics/lab` | Permintaan Cito Laboratorium |
| `radiology` | `#/diagnostics/radiology` | Pencitraan Radiologi & Diagnostik |
| `pharmacy` | `#/diagnostics/pharmacy` | E-Prescription Farmasi |
| `master-patient` | `#/master/patient` | Master Database Pasien |
| `master-practitioner` | `#/master/practitioner` | Master Tenaga Medis & Dokter |
| `master-departement` | `#/master/departement` | Master Instalasi & Departemen |
| `master-service-unit` | `#/master/service-unit` | Master Poliklinik & Unit Layanan |
| `master-room` | `#/master/room` | Master Ruangan & Tempat Tidur |
| `master-customer` | `#/master/customer` | Master Debitur & Penjamin (Customer) |
| `master-referal` | `#/master/referal` | Master Faskes Asal/Tujuan Rujukan |
| `master-tariff-class` | `#/master/tariff-class` | Master Kelas Tarif & Standar KRIS |

> [!TIP]
> **Keunggulan Hash-Based Router**: Saat browser di-refresh (`F5`), dibagikan lewat tautan langsung, atau menggunakan tombol Back/Forward browser, halaman **tidak akan reset ke index** dan tidak memerlukan konfigurasi URL rewrite khusus pada web server (Apache/Nginx/Laragon).

---

## 🚀 Skrip Pengembangan & Validasi

Jalankan perintah berikut di dalam folder `web/`:

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

## 🔗 Integrasi Backend (Vite Proxy)

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

Ketika Go backend aktif di port `8080`, frontend akan terhubung secara live. Apabila backend sedang offline, sistem memiliki mode **Fallback Demo** otomatis menggunakan kredensial bawaan:
- **Dokter**: `dokter` / `dokter123`
- **Admin**: `admin` / `admin123`
