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


