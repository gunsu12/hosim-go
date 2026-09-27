# Product Requirement Document (PRD): Catalog Item Domain

## 1. Ringkasan Eksekutif & Latar Belakang

Domain `internal/catalog/item` bertindak sebagai master katalog barang dan tarif sentral di HOSIM-GO. Katalog ini melayani seluruh kebutuhan operasional rumah sakit, mulai dari pelayanan klinis (Farmasi, Rawat Jalan, Rawat Inap, IGD), logistik & persediaan (Gudang Medis, Gudang Non-Medis, Depo), pengadaan (Procurement), hingga penagihan dan akuntansi (Billing & Finance).

Sebelumnya, pencatatan item rumah sakit berisiko tercampur atau tidak memiliki atribut yang presisi sesuai karakteristik fungsionalnya. Oleh karena itu, katalog item dirancang menggunakan pola **Shared Base Table (Class Table Inheritance)** yang memisahkan identitas umum barang/layanan dari atribut spesifik masing-masing subtipe:
1. **`item_medications`**: Obat-obatan, infus, dan vaksin dengan atribut dosis, bentuk sediaan, regulasi farmasi, dan integrasi SATUSEHAT KFA.
2. **`item_generals`**: Bahan Medis Habis Pakai (BMHP), instrumen medis/linen steril (siklus CSSD), dan barang umum/logistik non-medis (ATK, kebersihan, gizi).
3. **`item_assets`**: Barang modal dan alat kesehatan yang memiliki umur ekonomis, penyusutan (depresiasi), dan jadwal pemeliharaan/kalibrasi.
4. **`item_tariffs`**: Master katalog layanan dan jasa pasien (administrasi pendaftaran, akomodasi kamar/ranap, tindakan medis/keperawatan, dan pemeriksaan penunjang). Seluruh kebijakan penetapan harga, pemetaan kelas rawat (`tariff_classes`), versioning buku tarif, surcharge CITO, dan pecahan komponen biaya (jasa sarana, jasa medis dokter) didelegasikan sepenuhnya ke **Modul Finance (Buku Tarif / Tariff Price Plan)**.

Selain pemisahan subtipe fisik dan jasa, katalog ini mengintegrasikan **Akuntansi Persediaan dan Pendapatan** secara modular:
- **`item_categories`**: Memetakan kelompok barang/jasa ke **COA Pendapatan (Revenue Account)** pada modul billing.
- **`item_product_lines`**: Memetakan kelompok persediaan ke **COA Persediaan (Inventory Asset Account)** dan **COA Beban/HPP (COGS Account)** pada modul akuntansi/logistik. Untuk item bertipe `TARIFF`, relasi ini bernilai `NULL` karena bersifat non-persediaan (tidak memiliki kartu stok fisik).

---

## 2. Arsitektur Data & Model Relasional

### 2.1 Pola Inheritance: Shared Base Table

Setiap barang dan layanan tercatat pada tabel induk universal `items`. Modul `inventory` (kartu stok, batch/lot, mutasi), `billing` (tagihan, invoice kasir), dan `procurement` (purchase order) cukup mereferensikan satu foreign key universal: `item_id`.

```mermaid
erDiagram
    ITEM_CATEGORIES ||--o{ ITEMS : "mapping COA pendapatan"
    ITEM_PRODUCT_LINES ||--o{ ITEMS : "mapping COA persediaan & HPP (nullable untuk TARIFF)"
    ITEMS ||--|| ITEM_MEDICATIONS : "1:1 subtype obat"
    ITEMS ||--|| ITEM_GENERALS : "1:1 subtype umum / BMHP"
    ITEMS ||--|| ITEM_ASSETS : "1:1 subtype aset modal"
    ITEMS ||--|| ITEM_TARIFFS : "1:1 subtype tarif / jasa"
    ITEMS ||--o{ ITEM_UNITS : "1:N multi-satuan (UOM)"
    ITEMS ||--o{ TARIFF_PRICE_PLAN_ITEMS : "diberi tarif di modul finance (Buku Tarif)"

    ITEMS {
        uuid id PK
        string code UK "SKU / Barcode unik"
        string name "Nama resmi produk / layanan"
        string generic_name "Nama generik / zat aktif / nama alternatif"
        enum item_type "MEDICATION | GENERAL | ASSET | TARIFF"
        uuid category_id FK "Relasi ke item_categories (COA Pendapatan)"
        uuid product_line_id FK "Relasi ke item_product_lines (nullable untuk TARIFF)"
        boolean is_active "Status aktif katalog"
    }

    ITEM_MEDICATIONS {
        uuid item_id PK_FK
        string kfa_code "Kode KFA Kemenkes SATUSEHAT"
        string bpom_nie "Nomor Izin Edar BPOM"
        string dosage_form "Bentuk sediaan (Tablet, Sirup, Injeksi)"
        string strength_amount "Kekuatan dosis (500, 10, dll)"
        string strength_unit "Satuan kekuatan (mg, mcg, ml, %, dll)"
        string default_route "Rute pemberian (Oral, IV, IM, Topikal)"
        string medication_type "Bebas, Keras, Narkotika, Psikotropika, Prekursor"
        boolean is_high_alert "High Alert Medication (HAM)"
        boolean is_lasa "Look Alike Sound Alike (LASA / NORUM)"
        boolean is_fornas "Formularium Nasional"
        boolean is_antibiotic "Flag Antibiotik"
        string storage_temperature "Suhu simpan (Ruang, Kulkas 2-8C, Beku)"
    }

    ITEM_GENERALS {
        uuid item_id PK_FK
        string general_type "BMHP_MEDIS | INSTRUMEN_MEDIS | ATK | LINEN | KEBERSIHAN | DAPUR | LAINNYA"
        boolean is_sterile "Flag sterilitas produk akhir"
        boolean is_disposable "Single-use (true) vs Reusable (false)"
        boolean is_cssd_item "Apakah item ini dikelola dalam siklus sterilisasi CSSD"
        string sterilization_method "Metode sterilisasi: STEAM_AUTOCLAVE | EO_GAS | PLASMA | DRY_HEAT"
    }

    ITEM_ASSETS {
        uuid item_id PK_FK
        string brand "Merk / Brand manufaktur"
        string model_name "Tipe / Model mesin/alat"
        boolean is_medical_equipment "Alat medis vs Non-medis"
        int expected_life_years "Estimasi umur ekonomis (tahun)"
        string depreciation_method "Metode penyusutan (Straight Line, dsb)"
        int maintenance_interval_days "Interval servis / kalibrasi berkala (hari)"
    }

    ITEM_TARIFFS {
        uuid item_id PK_FK
        string name_alias "bisa dimanfaatkan untuk English name of the tariff/service atau alias lainnya"
        string charge_type "ADMINISTRASI | AKOMODASI | TINDAKAN | PENUNJANG | LAINNYA"
        string notes "Catatan operasional / deskripsi penagihan"
    }

    ITEM_UNITS {
        uuid id PK
        uuid item_id FK
        string unit_name "Nama satuan (Box, Strip, Tablet, Kali, Tindakan, Hari)"
        numeric conversion_factor "Rasio pengali ke base unit"
        boolean is_base_unit "Satuan terkecil (faktor = 1)"
        boolean is_purchase_unit "Satuan default saat PO"
        boolean is_dispense_unit "Satuan default saat resep/dispensing/billing"
    }

    ITEM_CATEGORIES {
        uuid id PK
        string code UK "Kode Kategori"
        string name "Nama Kategori"
        string item_type "Filter subtipe item (MEDICATION, GENERAL, ASSET, TARIFF)"
        string income_coa_code "Kode COA Akun Pendapatan"
        string discount_coa_code "Kode COA Akun Diskon Pendapatan"
        string sales_tax_coa_code "Kode COA Akun Pajak Penjualan"
        boolean is_active
    }

    ITEM_PRODUCT_LINES {
        uuid id PK
        string code UK "Kode Lini Produk"
        string name "Nama Lini Produk"
        string inventory_coa_code "Kode COA Akun Persediaan (Neraca/Aset)"
        string cogs_coa_code "Kode COA Akun Beban Pokok / HPP"
        string purchase_discount_coa_code "Kode COA Akun Diskon Pembelian"
        string purchase_tax_coa_code "Kode COA Akun Pajak Pembelian"
        string asset_coa_code "Kode COA Akun Aset Tetap (untuk item tipe ASSET)"
        string asset_accumulation_coa_code "Kode COA Akun Akumulasi Penyusutan (untuk item tipe ASSET)"
        string asset_depreciation_expense_coa_code "Kode COA Akun Beban Penyusutan (untuk item tipe ASSET)"
        boolean is_active
    }
```

---

## 3. Integrasi Akuntansi Finansial & Penjurnalan Ledger

Pemisahan antara `item_categories`, `item_product_lines`, dan `item_tariff_components` memberikan arsitektur akuntansi yang tangguh dan presisi:

| Entitas | Peran Akuntansi | Sisi Laporan | Trigger Transaksi | Berlaku Untuk |
|---|---|---|---|---|
| **`item_product_lines`** | **COA Persediaan** (`inventory_coa_code`) & **COA Beban/HPP** (`cogs_coa_code`) | **Neraca (Aset Lancar)** & **Laba Rugi (Beban Pokok)** | Saat penerimaan barang PO gudang (Debet: Persediaan), dan saat pemakaian/dispensing barang (Kredit: Persediaan, Debet: HPP). | `MEDICATION`, `GENERAL`, `ASSET` |
| **`item_tariff_components` (Modul Finance)** | **COA Komponen Spesifik** (`coa_code`) | **Neraca (Liabilitas) / Laba Rugi** | Dikelola pada Buku Tarif Modul Finance (`tariff_price_plan_item_components`). Memecah pendapatan sarana RS dan utang jasa medis dokter saat penagihan di kasir. | `TARIFF` |

### 3.1 Contoh Kasus Akuntansi Barang Fisik vs Tarif Jasa

1. **Obat Amoxicillin 500mg** (`MEDICATION`):
   - `product_line`: "Persediaan Obat Paten & Generik" -> Akun Persediaan: `114.01.001`, Akun HPP: `510.01.001`.
   - `category`: "Pendapatan Farmasi Rawat Jalan" -> Akun Pendapatan: `410.02.001`.
   - *Jurnal Penyerahan*: Debet Piutang/Kas, Kredit Pendapatan `410.02.001`; serta Debet HPP `510.01.001`, Kredit Persediaan `114.01.001`.

2. **Kasa Steril 10x10** (`GENERAL` - BMHP):
   - `product_line`: "Persediaan BMHP Bedah" -> Akun Persediaan: `114.02.001`, Akun HPP: `510.02.001`.
   - `category`: "Pendapatan Tindakan Medis / Bahan Medis" -> Akun Pendapatan: `410.03.001`.

3. **Konsultasi Dokter Spesialis (Poli Rawat Jalan)** (`TARIFF`):
   - `product_line_id`: `NULL` (non-inventory, tidak mempengaruhi saldo persediaan gudang).
   - `category`: "Pendapatan Layanan Rawat Jalan" -> Akun Utama: `410.01.001`.
   - `tariff_rates` (Kelas Poliklinik): Total Rp 200.000.
     - Komponen 1 (*Jasa Sarana RS*): Rp 50.000 -> COA `410.01.001` (Pendapatan Jasa Sarana RS).
     - Komponen 2 (*Jasa Medis Dokter*): Rp 150.000 -> COA `214.01.001` (Utang Jasa Medis Dokter / Fee Sharing).
   - *Jurnal Ledger Billing*:
     - **Debet**: Kas / Piutang Pasien `112.01.001` sebesar Rp 200.000
     - **Kredit**: Pendapatan Jasa Sarana `410.01.001` sebesar Rp 50.000
     - **Kredit**: Utang Jasa Medis Dokter `214.01.001` sebesar Rp 150.000

4. **Sewa Kamar Rawat Inap Kelas 1** (`TARIFF`):
   - `charge_type`: `AKOMODASI`
   - `product_line_id`: `NULL`
   - `tariff_rates` (Kelas 1): Total Rp 400.000 / Hari.
     - Komponen 1 (*Sewa Kamar & Fasilitas*): Rp 300.000 -> COA `410.02.001` (Pendapatan Akomodasi Ranap).
     - Komponen 2 (*Asuhan Keperawatan Standar*): Rp 100.000 -> COA `410.02.002` (Pendapatan Keperawatan).

---

## 4. Multi-Satuan Ukur (Multi-UOM) & Konversi

Rumah sakit membutuhkan konversi multi-satuan fleksibel untuk alur suplai barang logistik maupun standar penagihan tarif layanan:
- **`is_base_unit` (Satuan Terkecil)**: Nilai acuan stok inventory di database atau acuan kuantitas dasar layanan. Faktor konversi selalu `1.0`. Contoh barang: *Tablet*, *Kapsul*, *Pcs*, *ml*. Contoh tarif: *Kali*, *Tindakan*, *Hari*, *Sesi*, *Paket*.
- **`is_purchase_unit` (Satuan Pengadaan/Beli)**: Default saat membuat Purchase Requisition (PR) dan Purchase Order (PO). Contoh: *Box* (faktor = 100 Tablet), *Karton*, *Rim*. (Hanya untuk item barang fisik).
- **`is_dispense_unit` (Satuan Pelayanan/Resep/Billing)**: Default saat peresepan dokter, penyerahan obat, maupun penagihan charge tarif pasien. Contoh obat: *Strip* (faktor = 10). Contoh tarif: *Kali*, *Hari*.

### 4.1 Standarisasi UOM untuk Item Tipe TARIFF

Setiap item bertipe `TARIFF` **wajib mendaftarkan minimal 1 base unit** pada `item_units` dengan konfigurasi:
- `conversion_factor = 1.0`
- `is_base_unit = true`
- `is_dispense_unit = true`
- `is_purchase_unit = false`

Pendekatan ini memastikan modul `billing` memperlakukan seluruh item secara **polimorfik** tanpa memerlukan conditional logic khusus barang vs jasa saat membaca unit penagihan (`unit_id`).

### 4.2 Contoh Data Konversi Riil

| Kategori Item | Satuan | `conversion_factor` | `is_base_unit` | `is_purchase_unit` | `is_dispense_unit` | Keterangan Operasional |
|---|---|---|---|---|---|---|
| **Amoxicillin 500mg** | Tablet | 1 | `true` | `false` | `true` | Satuan stok dasar & eceran |
| | Strip | 10 | `false` | `false` | `true` | Satuan resep dokter lazim |
| | Box | 100 | `false` | `true` | `false` | Satuan beli kemasan distributor |
| **Paracetamol Sirup** | Botol | 1 | `true` | `false` | `true` | Satuan konsumsi eceran |
| | Karton | 24 | `false` | `true` | `false` | Satuan beli grosir |
| **Sarung Tangan Bedah** | Pasang | 1 | `true` | `false` | `true` | Satuan pemakaian tindakan |
| | Box | 50 | `false` | `false` | `false` | Distribusi depo ruang OK |
| **Konsul Dokter Spesialis** | Kali | 1 | `true` | `false` | `true` | Satuan billing per sesi kunjungan |
| **Sewa Kamar VIP** | Hari | 1 | `true` | `false` | `true` | Satuan billing akomodasi per hari |
| **Tindakan Jahit Luka** | Tindakan | 1 | `true` | `false` | `true` | Satuan billing tindakan IGD/Bedah |

---

## 5. Regulasi & Karakteristik Khusus Subtipe Item

### 5.1 Kepatuhan Regulasi Farmasi & SATUSEHAT (`item_medications`)

1. **Integrasi KFA Kemenkes (SATUSEHAT)**: Kolom `kfa_code` menyimpan kode identifikasi resmi dari Kamus Farmasi dan Alat Kesehatan Kemenkes RI untuk interoperabilitas RME.
2. **Izin Edar BPOM**: Kolom `bpom_nie` untuk mencatat validitas registrasi izin edar obat di Indonesia.
3. **Keselamatan Pasien (Patient Safety)**:
   - `is_high_alert`: Obat yang memerlukan kewaspadaan tinggi.
   - `is_lasa`: Obat dengan nama, rupa, atau ucapan mirip (*Look-Alike Sound-Alike* / NORUM).
   - `is_antibiotic`: Kontrol penggunaan antibiotik dan program resistensi antimikroba (PPRA).
4. **Regulasi Narkotika & Psikotropika**: Klasifikasi ketat untuk pelaporan berkala SIPNAP.

### 5.2 Standarisasi Sterilisasi & Siklus CSSD (`item_generals`)

1. **Klasifikasi Barang CSSD vs Non-CSSD**:
   - `is_cssd_item = true` menandai bahwa siklus pakai item menjadi wewenang instalasi CSSD.
   - `is_disposable = false` (reusable) dan `is_sterile = true` mengidentifikasi instrumen atau linen yang wajib diproses ulang.
2. **Spesifikasi Mesin & Metode Sterilisasi (`sterilization_method`)**: `STEAM_AUTOCLAVE`, `PLASMA`, `EO_GAS`, `DRY_HEAT`.

### 5.3 Spesifikasi Service Charge & Komponen Tarif (`item_tariffs`)

Item tarif berfungsi sebagai representasi master layanan dan beban biaya yang ditagihkan kepada pasien (*patient service charges*):
1. **Klasifikasi Beban (`charge_type`)**:
   - `ADMINISTRASI`: Beban administrasi loket pendaftaran, pembuatan kartu berobat, atau pengelolaan berkas rekam medis.
   - `AKOMODASI`: Beban sewa tempat tidur dan kamar rawat inap (VIP, Kelas 1, 2, 3, ICU, Isolasi).
   - `TINDAKAN`: Jasa tindakan medis, konsultasi dokter, pemeriksaan fisik, atau tindakan keperawatan.
   - `PENUNJANG`: Biaya pemeriksaan laboratorium, radiologi, ambulans, atau layanan penunjang diagnostik.
   - `LAINNYA`: Biaya jasa pendukung lainnya.
2. **Pendelegasian Kebijakan Harga ke Modul Finance (Buku Tarif)**:
   - Modul `catalog` tidak mengunci nominal harga ataupun aturan kelas rawat inap secara kaku.
   - Penentuan harga per kelas rawat (`tariff_classes`), rincian komponen biaya (`tariff_components`), status approval SK Direksi, tarif CITO, serta diferensiasi buku tarif umum vs penjamin (PKS Asuransi/Perusahaan) dikelola secara berversi (*versioned*) di **Modul Finance (Buku Tarif / Tariff Price Plan)**.

---

## 6. Desain Endpoint API

### 6.1 Agregat & Universal Search (Untuk Modul Inventory, Billing, Procurement)
- `GET /api/v1/catalog/items`: Pencarian umum seluruh tipe item dengan filter `item_type` (`MEDICATION`, `GENERAL`, `ASSET`, `TARIFF`), `category_id`, `product_line_id`, `search` (kode/nama), pagination.
- `GET /api/v1/catalog/items/:id`: Detail lengkap item beserta ekstensi subtipe (`medication`, `general`, `asset`, `tariff`), daftar `units`, dan untuk tarif menyertakan relasi `rates` beserta `components`.

### 6.2 Pengelolaan Khusus Subtipe (Input Form Terpisah & Validasi Ketat)
- **Obat**:
  - `POST /api/v1/catalog/items/medications`: Mendaftarkan item baru bertipe medication beserta field dosis & regulasi.
  - `PUT /api/v1/catalog/items/medications/:id`: Memperbarui data umum dan atribut obat.
- **Barang Umum / BMHP**:
  - `POST /api/v1/catalog/items/generals`: Mendaftarkan item umum / BMHP.
  - `PUT /api/v1/catalog/items/generals/:id`: Memperbarui data umum dan atribut general.
- **Aset / Alkes Modal**:
  - `POST /api/v1/catalog/items/assets`: Mendaftarkan master tipe aset (spesifikasi, masa manfaat, metode depresiasi).
  - `PUT /api/v1/catalog/items/assets/:id`: Memperbarui master tipe aset.
- **Tarif Layanan / Jasa**:
  - `POST /api/v1/catalog/items/tariffs`: Mendaftarkan master layanan/tindakan baru (data umum item, base unit penagihan, subtipe `charge_type`, serta catatan operasional).
  - `PUT /api/v1/catalog/items/tariffs/:id`: Memperbarui data umum item dan atribut `item_tariffs` (`charge_type`, `notes`).
  - *(Catatan: Penetapan nominal harga, kelas, dan komponen biaya dikelola melalui API Modul Finance `/api/v1/finance/price-plans`)*.

### 6.4 Pengelolaan Satuan & Akuntansi
- `POST /api/v1/catalog/items/:id/units`: Menambah satuan alternatif & rasio konversi.
- `PUT /api/v1/catalog/items/:id/units/:unit_id`: Mengubah konfigurasi satuan atau rasio konversi.
- `DELETE /api/v1/catalog/items/:id/units/:unit_id`: Menghapus satuan alternatif (selama bukan base unit).
- `CRUD /api/v1/catalog/item-categories`: Pengelolaan master kategori & mapping COA Pendapatan.
- `CRUD /api/v1/catalog/item-product-lines`: Pengelolaan master lini produk & mapping COA Persediaan/HPP.

---

## 7. Rencana Kerja Implementasi (Roadmap)

1. **Database Migration (Goose)**:
   - Buat migrasi `migrations/00007_create_master_item_tables.sql` mencakup tabel:
     - `item_categories`
     - `item_product_lines`
     - `items` (`product_line_id` nullable, enum `item_type` mencakup `TARIFF`)
     - `item_medications`
     - `item_generals`
     - `item_assets`
     - `item_tariffs`
     - `item_units`
2. **Domain Entities & Enums**:
   - `pkg/enums/item.go`: Enum tipe item (`MEDICATION`, `GENERAL`, `ASSET`, `TARIFF`), general type, medication type, depreciation method, charge type (`ADMINISTRASI`, `AKOMODASI`, `TINDAKAN`, `PENUNJANG`, `LAINNYA`).
   - `internal/catalog/item/entity.go`: Struct domain GORM dengan audit trail dan UUID v7 untuk tabel-tabel master item di atas.
3. **Repository Layer**:
   - `internal/catalog/item/repository.go`: Query teroptimasi, eager-loading subtype (`Medication`, `General`, `Asset`, `Tariff`), filter pagination.
4. **Service & Validation Layer (Jalur A Clean Architecture)**:
   - `internal/catalog/item/service.go`: Validasi integritas base unit, kode unik, dan transaksi atomik pembuatan base + subtype.
5. **Transport Layer (Gin HTTP Handler & DTO)**:
   - `internal/catalog/item/dto.go`: Request dan response DTO terpisah dari model entity.
   - `internal/catalog/item/handler.go`: Handler endpoint universal, subtipe obat, general, aset, dan layanan tarif.
6. **Testing & Database Seeder**:
   - `internal/catalog/item/service_test.go`: Unit test bisnis logic dan konversi satuan.
   - `internal/catalog/item/seeder.go`: Sampel awal kategori, lini produk, obat KFA, BMHP, aset, serta sampel master tarif layanan pendaftaran, akomodasi, dan tindakan rawat jalan.

