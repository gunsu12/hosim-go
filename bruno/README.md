# Bruno Collection - HOSIM-GO API

Koleksi API request untuk aplikasi [Bruno](https://www.usebruno.com/).

## Cara Menggunakan

1. Buka aplikasi **Bruno**.
2. Pilih **Open Collection**, lalu arahkan ke folder `bruno/` di root repository ini.
3. Pilih Environment **Local** di pojok kanan atas Bruno.
4. Jalankan request **Auth > Login** terlebih dahulu:
   - Request ini akan otomatis menyimpan `access_token` ke environment variable `token`.
5. Anda dapat langsung menjalankan request di folder:
   - **`Finance`**:
     - **`Customers`**: `List Customer Types`, `List Customers`, `Create Customer` (otomatis simpan `{{customerId}}`), `Get by ID`, `Update`, `Delete`.
     - **`Tariff Classes`**: `List`, `Create` (otomatis simpan `{{tariffClassId}}`), `Get by ID`, `Update`, `Delete`.
     - **`Tariff Components`**: `List`, `Create` (otomatis simpan `{{tariffComponentId}}`), `Get by ID`, `Update`, `Delete`.
   - **`Organization`**:
     - **`Departments`**: `List`, `Create` (otomatis simpan `{{departmentId}}`), `Get by ID`, `Update`, `Delete`.
     - **`Service Units`**: `List`, `Create` (otomatis simpan `{{serviceUnitId}}`), `Get by ID`, `Update`, `Delete`.
     - **`Rooms`**: `List`, `Create` (otomatis simpan `{{roomId}}`), `Get by ID`, `Update`, `Delete`.
     - **`Referals`**: `List`, `Create` (otomatis simpan `{{referalId}}`), `Get by ID`, `Update`, `Delete`.
   - **`Catalog`**:
     - **`ItemCategory`**:
       - `List Categories`: List master kategori item dengan filter pencarian & tipe item.
       - `Create Category`: Membuat master kategori baru (mapping COA Pendapatan) $\rightarrow$ otomatis simpan `{{categoryId}}`.
       - `Get Category by ID`: Detail kategori.
       - `Update Category`: Pembaruan data kategori.
       - `Delete Category`: Soft delete kategori.
     - **`ItemProductLine`**:
       - `List Product Lines`: List master lini produk persediaan.
       - `Create Product Line`: Membuat master lini produk baru (mapping COA Persediaan, HPP, Aset) $\rightarrow$ otomatis simpan `{{productLineId}}`.
       - `Get Product Line by ID`: Detail lini produk.
       - `Update Product Line`: Pembaruan data lini produk.
       - `Delete Product Line`: Soft delete lini produk.
     - **`ItemMedication`**:
       - `List Medications`: Filter list khusus obat (`item_type=MEDICATION`).
       - `Create Medication`: Registrasi obat (dosis, sediaan, KFA SATUSEHAT, BPOM NIE, LASA, HAM) $\rightarrow$ otomatis simpan `{{medicationId}}`.
       - `Get Medication by ID`: Detail obat lengkap beserta subtipe dan satuan.
       - `Update Medication`: Pembaruan data obat.
       - `Delete Medication`: Soft delete data obat.
       - `List Medication Units`: Daftar satuan konversi obat.
       - `Add Unit to Medication`: Menambah konversi satuan obat $\rightarrow$ otomatis simpan `{{unitId}}`.
     - **`ItemGeneral`**:
       - `List Generals`: Filter list BMHP & barang umum (`item_type=GENERAL`).
       - `Create General Item`: Registrasi BMHP/umum (CSSD, sterilitas, disposable) $\rightarrow$ otomatis simpan `{{generalId}}`.
       - `Get General by ID`: Detail barang umum.
       - `Update General Item`: Pembaruan data barang umum.
       - `Delete General Item`: Soft delete data barang umum.
       - `List General Units`: Daftar satuan konversi barang umum.
       - `Add Unit to General`: Menambah konversi satuan general $\rightarrow$ otomatis simpan `{{unitId}}`.
     - **`ItemAsset`**:
       - `List Assets`: Filter list alat modal/aset RS (`item_type=ASSET`).
       - `Create Asset Item`: Registrasi alat modal (merk, model, masa manfaat, metode depresiasi, interval kalibrasi) $\rightarrow$ otomatis simpan `{{assetId}}`.
       - `Get Asset by ID`: Detail aset modal.
       - `Update Asset Item`: Pembaruan data aset modal.
       - `Delete Asset Item`: Soft delete data aset.
       - `List Asset Units`: Daftar satuan aset.
       - `Add Unit to Asset`: Menambah satuan aset $\rightarrow$ otomatis simpan `{{unitId}}`.
     - **`ItemTariff`**:
       - `List Tariffs`: Filter list tarif layanan pasien (`item_type=TARIFF`).
       - `Create Tariff Item`: Registrasi tarif non-inventory (tindakan, akomodasi, dsb) $\rightarrow$ otomatis simpan `{{tariffId}}`.
       - `Get Tariff by ID`: Detail tarif layanan.
       - `Update Tariff Item`: Pembaruan data tarif layanan.
       - `Delete Tariff Item`: Soft delete tarif layanan.
       - `List Tariff Units`: Daftar satuan penagihan tarif.
       - `Add Unit to Tariff`: Menambah satuan penagihan alternatif $\rightarrow$ otomatis simpan `{{unitId}}`.
     - **`ItemUnit`**:
       - `List Units by Item ID`: Memuat seluruh satuan untuk item apapun (`{{itemId}}`).
       - `Add Unit to Item`: Mendaftarkan satuan alternatif (UOM) ke item $\rightarrow$ otomatis simpan `{{unitId}}`.
       - `Update Unit`: Mengubah nama satuan atau rasio konversi.
       - `Delete Unit`: Menghapus satuan alternatif (terproteksi larangan menghapus base unit).
