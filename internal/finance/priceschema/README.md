# Price Schema Domain (Skema Harga & Penyesuaian Tarif PKS)

Dokumentasi ini ditujukan bagi kontributor (developer dan AI agent) agar memahami arsitektur, batas modul, invariant bisnis, dan acuan implementasi sub-modul `priceschema` di HOSIM-GO sesuai [`PRD.md`](./PRD.md).

---

## 1. 🎯 Gambaran & Konteks Bisnis

Sub-modul `internal/finance/priceschema` bertindak sebagai **pusat aturan penyesuaian tarif (*Pricing Adjustment Rules Engine*)** untuk pasien rekanan/penjamin berdasarkan **PKS (Perjanjian Kerja Sama)** tanpa menduplikasi buku tarif dasar rumah sakit.

- **Buku Tarif Dasar RS (`internal/finance/priceplan`)**: Menyimpan harga dasar per item tindakan medis × kelas layanan (`tariff_classes`) beserta rincian komponen biaya baku.
- **Skema Harga (`internal/finance/priceschema`)**: Menyimpan aturan penyesuaian (markup, diskon, fixed price, PPN, dan pembagian tanggungan pasien vs penjamin/cost sharing) yang diterapkan di atas harga dasar.

### Hierarki Pencocokan Rule (Urutan Override)
Pencocokan rule dilakukan secara spesifik dari level paling presisi hingga default:
1. `ITEM_UNIT`: Item tertentu pada Service Unit tertentu
2. `ITEM_DEPT`: Item tertentu pada Departemen tertentu
3. `ITEM`: Item tertentu di seluruh RS
4. `UNIT`: Seluruh item pada Service Unit tertentu
5. `DEPT`: Seluruh item pada Departemen tertentu
6. `SCHEMA_DEFAULT`: Aturan standar skema jika tidak ada rule di atas yang cocok
7. *(Fallback)*: Jika pasien tidak memiliki skema harga (pasien umum/pribadi), gunakan harga dasar Price Plan tanpa penyesuaian dengan tanggungan 100% pasien.

---

## 2. 🧮 Single Source of Truth (SSoT) Kalkulasi Tarif

Sistem HOSIM-GO menetapkan bahwa **seluruh perhitungan matematis tarif tindakan/layanan medis** dipusatkan secara eksklusif pada file:

> 📍 [`internal/finance/priceschema/calculation.go`](./calculation.go)

File ini berisi *pure domain functions* (deterministik, tanpa I/O, teruji table-driven unit test):
- **`ApplyAdjustment(basePrice, adjType, adjValue)`**: Menghitung harga setelah markup/diskon persentase, nominal tetap, atau fixed price PKS.
- **`CalculateTax(finalPrice, taxPercent)`**: Menghitung nilai PPN dari harga setelah penyesuaian.
- **`CalculateCostSharing(gross, payByPatientPercent)`**: Membagi total tagihan bruto ke porsi tanggungan pasien dan penjamin.
- **`DistributeToComponents(finalPriceTanpaCito, totalBasePrice, isCito, baseComponents)`**: Mendistribusikan selisih/delta tarif secara proporsional ke rincian komponen biaya (Jasa Sarana RS, Jasa Medis Dokter, Anestesi, Paramedis) dengan mengalokasikan sisa pembulatan (*rounding difference*) ke komponen Jasa Sarana RS.

---

## 3. 🛡️ Invariant & Aturan Bisnis (Non-Negotiable)

1. **Prinsip Override Utuh (Bukan Aditif)**:
   - Rule paling spesifik yang cocok dipakai **sepenuhnya** untuk seluruh nilai (penyesuaian tarif, PPN, dan cost sharing). Tidak ada penjumlahan atau pewarisan antar-level hierarki.
2. **Scope Exclusivity**:
   - Kolom `departement_id` dan `service_unit_id` pada rule **tidak boleh diisi bersamaan**. Departemen diturunkan secara dinamis dari unit organisasi saat lookup untuk mencegah inkonsistensi bila unit dipindah departemen.
3. **Diskon PKS vs CITO**:
   - Diskon/markup PKS **tidak memotong biaya tambahan CITO**. Biaya CITO dihitung murni dari selisih CITO pada Price Plan (`total_cito_price - total_base_price`) dan ditambahkan secara utuh setelah penyesuaian tarif reguler selesai.
4. **Audit Trail Mutasi (Jalur B)**:
   - Setiap operasi create, update, deactivate schema, atau upsert/delete rule **wajib mencatat audit append-only** (*before/after snapshot*) dalam transaksi database yang sama via `TransactionManager`.
5. **Immutable Snapshot di Modul Billing**:
   - Perubahan skema/rule tidak pernah menghitung ulang tagihan lampau. Modul penagihan (*Billing*) wajib menyimpan snapshot hasil perhitungan saat charge dibuat.

---

## 4. 📂 Peta File Sumber Daya (Code Tour)

| File | Peran & Deskripsi |
|---|---|
| [`PRD.md`](./PRD.md) | Dokumen spesifikasi kebutuhan produk & acuan arsitektur lengkap. |
| [`calculation.go`](./calculation.go) | **Core Math Engine / SSoT**: Fungsi matematis murni untuk adjustment, PPN, cost-sharing, dan distribusi komponen. |
| [`calculation_test.go`](./calculation_test.go) | Table-driven unit test komprehensif untuk seluruh skenario kalkulasi & pembulatan. |
| [`entity.go`](./entity.go) | Entitas `PriceSchema` & `PriceSchemaRule`, validasi invariant, dan enum level match. |
| [`ports.go`](./ports.go) | Kontrak domain & interface dependency (`RuleMatcher`, `Repository`, readers lintas modul). |
| [`matcher.go`](./matcher.go) | Implementasi `RuleMatcher` yang menjalankan resolusi override bertingkat. |
| [`repository.go`](./repository.go) | Implementasi GORM persistence untuk query dan mutasi data rule/schema. |
| [`usecases.go`](./usecases.go) | Orchestrator Jalur B: Create, Update, Deactivate Schema, Upsert/Delete Rule, dan `SimulateTariffUseCase`. |
| [`adapters.go`](./adapters.go) | Adapter integrasi ke modul luar (`department`, `serviceunit`, `item`, dan `priceplan`). |
| [`dto.go`](./dto.go) | Request & Response DTO untuk HTTP transport. |
| [`handler.go`](./handler.go) | Gin HTTP Handler dengan otorisasi RBAC permission `finance.price_schema.*`. |
| [`usecase_test.go`](./usecase_test.go) | Integration & unit test untuk use case mutasi, simulasi harga, dan verifikasi audit. |

---

## 5. ⚠️ Roadmap & Technical Debt Tracker

Bagian ini mencatat checklist teknis yang harus diperhatikan saat mengintegrasikan modul operasional lain:

### [ ] 1. Integrasi Modul Billing: *Patient Visit Charge Tariff*
- Saat modul `internal/billing` dibangun untuk mencatat tagihan tindakan pasien (rawat jalan, rawat inap, IGD):
  1. Ambil tarif dasar RS dari `priceplan.LookupTariffBase`.
  2. Ambil aturan penyesuaian PKS dari `priceschema.RuleMatcher`.
  3. Hitung nominal final, PPN, cost sharing, dan pecahan komponen melalui fungsi-fungsi di [`calculation.go`](./calculation.go).
  4. Simpan **Immutable Snapshot** pada tabel charge line billing (`priced_at`, `price_plan_id`, `price_schema_id`, `final_price`, komponen).
- **Aturan**: Modul pelayanan (*outpatient*, *inpatient*, *emergency*) maupun *billing* dilarang membuat rumus penentuan tarif ad-hoc sendiri.

### [ ] 2. Deprecasi Langkah 1 di `priceplan/lookup_tariff.go`
- Di [`internal/finance/priceplan/lookup_tariff.go`](../priceplan/lookup_tariff.go#L41-L56), masih terdapat query legacy `FindActiveCustomPlan` via `customer_id`.
- Sesuai keputusan arsitektur [PRD.md §9](./PRD.md#L358-L367), hapus Langkah 1 tersebut saat modul billing mulai menggunakan `Price Schema`, sehingga `priceplan` murni menjadi rujukan buku tarif dasar RS.

### [ ] 3. Single Source of Functionality untuk Obat & BMHP (Fase 2)
- Kolom `medication_adjustment_type` dan `medication_adjustment_value` sudah tersedia di tabel `price_schema_rules`.
- Untuk kalkulasi harga jual barang farmasi (obat dan BMHP):
  - Wajib dibuatkan **satu engine kalkulasi terpusat** (mengikuti pola murni `calculation.go`) yang menangani: harga dasar/HNA/HPP, margin/markup apotek, diskon PKS, tuslah/embalase resep, PPN, dan cost sharing.
  - Modul `billing` (kasir) dan `pharmacy` (dispensing farmasi rawat jalan & ranap) wajib mengonsumsi engine kalkulasi yang sama guna mencegah disparitas harga tagihan obat.
