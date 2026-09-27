# Customer ( Penjamin Bayar )

Dokumentasi ini ditujukan bagi kontributor (developer dan AI agent) agar memahami konteks bisnis, integritas data, dan aturan operasional entitas `Customer` (Penjamin Bayar) di HOSIM-GO.

## 1. Gambaran dan Definisi

Dalam konteks manajemen rumah sakit, `Customer` mengacu pada **entitas penjamin pembayaran** yang mewakili entitas legal atau korporat yang bertanggung jawab atas pembiayaan layanan pasien. Entitas ini berfungsi sebagai "tuan rumah" (*parent entity*) yang menaungi pasien, klaim asuransi, serta kontrak layanan medis yang disepakati.

> Tipe Customer
>> 1. Personal ( pasien individual / umum )
>> 2. Corporate ( perusahaan / instansi )
>> 3. BPJS ( Badan Penyelenggara Jaminan Sosial )
>> 4. Insurance ( Asuransi ) 
>> 5. Internal ( untuk karyawan rumah sakit biasanya ada benefit tersendiri )

## 2. 🛡️ Invariant & Aturan Bisnis (Non-Negotiable)
1. **Imutability** untuk customer dengan tipe personal itu sudah pasti harus ada. tidak boleh di update 
2. **Untuk customer tipe BPJS** itu tidak boleh di hapus, hanya bisa di nonaktifkan
3. **Untuk customer tipe Corporate dan Insurance** itu boleh di hapus atau di nonaktifkan, karena biasanya data ini akan di update seiring waktu

## 3. Relasi dan Kardinalitas Data
```
erDiagram
    customer_types ||--o{ customers : "has many"

    customers {
        string id PK "size:36"
        string code UK "size:255"
        string name "size:255"
        string address "nullable, size:255"
        string phone "nullable, size:255"
        string email "nullable, size:255"
        string website "nullable, size:255"
        string contact_person "nullable, size:255"
        bool require_card "nullable, default:true"
        string description "nullable, size:255"
        string customer_type_id FK "nullable, size:36"
        bool is_active "default:true"
        bool is_immutable "default:false"
        timestamp created_at "default:CURRENT_TIMESTAMP"
        timestamp updated_at "default:CURRENT_TIMESTAMP"
        timestamp deleted_at "nullable"
        string created_by "default:SYSTEM"
        string updated_by "default:SYSTEM"
        string deleted_by "nullable"
    }

    customer_types {
        string id PK "size:36"
        string code UK "size:255"
        string name "size:255"
        timestamp created_at "default:CURRENT_TIMESTAMP"
        timestamp updated_at "default:CURRENT_TIMESTAMP"
        timestamp deleted_at "nullable"
        string created_by "default:SYSTEM"
        string updated_by "default:SYSTEM"
        string deleted_by "nullable"
    }
```

## 4. 📂 Peta File Sumber Daya (Code Tour)

- [`entity.go`](file:///c:/laragon/www/hosim-go/internal/finance/customer/entity.go): Definisi struct `Customer`, GORM tag, relasi, dan hook ID UUIDv7.
- [`repository.go`](file:///c:/laragon/www/hosim-go/internal/finance/customer/repository.go): Implementasi repository interface `IRepository`.
- [`service.go`](file:///c:/laragon/www/hosim-go/internal/finance/customer/service.go): Implementasi service interface `IService`.
- [`handler.go`](file:///c:/laragon/www/hosim-go/internal/finance/customer/handler.go): Implementasi handler interface `IHandler`.
- [`doc.go`](file:///c:/laragon/www/hosim-go/internal/finance/customer/doc.go): Dokumentasi package Go.
