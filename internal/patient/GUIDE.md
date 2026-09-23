# Panduan Step-by-Step: Membuat Master Data Pasien (internal/master/patient)

Halo! Dokumen ini disiapkan sebagai panduan langkah demi langkah saat kamu siap mengetik modul **Master Data Pasien** secara mandiri.

---

## 1. Domain Model: `internal/master/patient/entity.go`

Tentukan field yang dibutuhkan oleh rumah sakit dan standar SatuSehat (Kemenkes):
- `ID` (UUID atau string)
- `MedicalRecordNo` (No. RM unik, misal: `RM-000001`)
- `NIK` (16 digit kependudukan)
- `Name` (Nama lengkap pasien)
- `Gender` (`L` / `P`)
- `BirthDate` (Tanggal lahir)
- `BirthPlace` (Tempat lahir)
- `Address` (Alamat KTP/domisili)
- `Phone` (No. HP/WhatsApp aktif)
- `CreatedAt`, `UpdatedAt`, `DeletedAt` (Soft delete GORM)

---

## 2. Database Layer: `internal/master/patient/repository.go`

Buat interface `Repository` dan implementasinya menggunakan GORM:
```go
type Repository interface {
    Create(ctx context.Context, patient *Patient) error
    FindByID(ctx context.Context, id string) (*Patient, error)
    FindByMedicalRecordNo(ctx context.Context, rmNo string) (*Patient, error)
    FindByNIK(ctx context.Context, nik string) (*Patient, error)
    FindAll(ctx context.Context, limit, offset int, search string) ([]Patient, int64, error)
    Update(ctx context.Context, patient *Patient) error
    Delete(ctx context.Context, id string) error
}
```

---

## 3. Business Logic: `internal/master/patient/service.go`

Tulis logika bisnis rumah sakit:
- Cek apakah NIK sudah pernah terdaftar (mencegah duplikasi data pasien).
- Penomoran otomatis No. Rekam Medis (auto-generate nomor urut).
- Validasi usia atau format data.

---

## 4. HTTP Controller: `internal/master/patient/handler.go`

Tangani request Gin:
- Struct Request (`CreatePatientRequest`, `UpdatePatientRequest`) dengan tag `binding:"required,len=16"` dll.
- Gunakan `response.Success` dan `response.Error` dari `pkg/response`.
- Daftarkan route ke grup router `/api/v1/patients`.
