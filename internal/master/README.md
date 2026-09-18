# Modul Master Data (HOSIM-GO)

Folder ini diperuntukkan untuk seluruh sub-modul **Master Data** SIMRS.

## Sub-Modul yang Direncanakan:
- `patient/`: Master Data Pasien & Rekam Medis (Akan kamu buat manual sambil belajar).
- `practitioner/`: Dokter, Perawat, Petugas Medis, Spesialisasi.
- `organization/`: Poliklinik, Instalasi/Departemen, Ruangan/Bangsal, Tempat Tidur (Bed Management).
- `clinical/`: Master ICD-10 (Diagnosa) dan ICD-9-CM (Prosedur/Tindakan Medis).

## Blueprint Struktur Sub-Modul `patient/`:
Saat kamu siap mengetik kode Master Pasien, kamu dapat membuat file-file berikut di dalam `internal/master/patient/`:
1. `entity.go`: Struct model database Pasien (GORM tags, JSON tags).
2. `repository.go`: Interface & implementasi query database (Create, FindByID, FindAll, Update, Delete).
3. `service.go`: Business logic (validasi NIK, generate nomor rekam medis, dll).
4. `handler.go`: HTTP controller Gin (request binding, validation error formatting, response).
