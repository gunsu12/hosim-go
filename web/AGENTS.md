# AGENTS.md — HOSIM-GO Web Frontend Engineering Rules

Dokumen ini berlaku untuk seluruh folder `web/` di repositori HOSIM-GO. Semua coding agent **WAJIB** membaca dan mematuhi ringkasan aturan ini sebelum membuat, mengubah, atau me-refactor kode frontend.

Kata **WAJIB**, **DILARANG**, dan **BOLEH** bersifat normatif.

---

## 1. Quick Reference: Aturan Emas Frontend

| Kebutuhan UI / Logic | WAJIB Gunakan | DILARANG Keras |
|---|---|---|
| **Dialog / Form Modal** | [`M3Modal.svelte`](./src/lib/components/m3/M3Modal.svelte) | Backdrop manual `fixed inset-0 bg-black/40`, frame card custom |
| **Dropdown / Pilihan Opsi** | [`M3Select.svelte`](./src/lib/components/m3/M3Select.svelte) | Tag native HTML `<select>` & `<option>` (merusak radius & dark mode) |
| **Toggle / Status Boolean** | [`M3ToggleCard.svelte`](./src/lib/components/m3/M3ToggleCard.svelte) | Tag native `<input type="checkbox">` |
| **Dokter / Paramedis** | [`PractitionerLookupInput.svelte`](./src/lib/components/PractitionerLookupInput.svelte) | Native `<select>` atau load array statis tanpa server search |
| **Bagan Akun (COA)** | [`CoaLookupInput.svelte`](./src/lib/components/CoaLookupInput.svelte) | Input teks biasa atau dropdown statis |
| **Katalog Tarif & Tindakan** | [`TariffLookupInput.svelte`](./src/lib/components/TariffLookupInput.svelte) | Input teks biasa atau dropdown statis |
| **Tombol Aksi** | [`M3Button.svelte`](./src/lib/components/m3/M3Button.svelte) | Tag `<button>` tanpa styling standar M3 |
| **Reaktivitas Komponen** | Svelte 5 Runes (`$state`, `$derived`, `$props`, `$effect`) | Sintaks Svelte 3/4 (`export let`, `$:` reactive declaration, `<slot />`) |
| **Event Dispatcher** | Callback props (`onclick`, `onClose`, `onSuccess`) | `createEventDispatcher` |
| **Batas Ukuran File** | **Target ~500 baris per file** (fleksibel bila orkestrasi belum dapat dirampingkan) | Menumpuk semua logika tanpa dekomposisi |
| **Aksi Baris Tabel** | 1 Aksi Utama + `MoreVertical` (`fixed` & smart flipping) | 4–5 tombol berderet atau dropdown `absolute` terpotong di `overflow` |
| **Filter & Pencarian** | Single-row filter bar + status pills horizontal | Filter bertingkat 2–3 baris yang memakan ruang vertikal layar |
| **Feedback Operasi** | Floating toast notification (`fixed bottom-6 right-6`) | Banner alert di atas tabel yang memicu layout shift |
| **Tema & Tampilan** | Full Dark Mode adaptif (`dark:...`) di setiap elemen | Warna hardcoded tanpa padanan kelas dark mode |

---

## 2. Standar Dialog & Modal (Google M3)

Pola desain pada [`CustomerModal.svelte`](./src/lib/pages/finance/customer/components/CustomerModal.svelte) adalah **benchmark resmi (Gold Standard)** untuk form dialog.

### 2.1. Aturan Penggunaan M3Modal
- **WAJIB** memakai `M3Modal.svelte` sebagai container dialog.
- **DILARANG** menduplikasi markup backdrop (`fixed inset-0 bg-black/40...`), tombol close `X`, atau listener tombol ESC secara manual.
- **Wajib properti**: `isOpen`, `title`, `subtitle`, `icon`, `onClose`, dan `onsubmit` (bila bertipe form).
- **Pilihan ukuran (`size`)**:
  - `'sm'` / `'md'`: Konfirmasi aksi singkat / delete confirmation dialog.
  - `'lg'` / `'xl'`: Form ringkas (3–5 field).
  - `'2xl'` (**Default**): Form master data & entitas konfigurasi standar.
  - `'3xl'` / `'4xl'`: Form multi-tab kompleks atau embedded data table.
- **Tombol Aksi Footer**: Letakkan di dalam `{#snippet footer()}` pada `M3Modal`:
  - Batal: `<M3Button variant="text" onclick={onClose} disabled={isSubmitting}>Batal</M3Button>`
  - Submit: `<M3Button variant="filled" type="submit" disabled={isSubmitting}>` dengan ikon Lucide atau spinner saat submit.

### 2.2. Modal Konfirmasi Destruktif (Delete Modal)
- Benchmark resmi: [`CoaDeleteModal.svelte`](./src/lib/pages/accounting/components/CoaDeleteModal.svelte).
- Gunakan `size="md"`, `icon={Trash2}`, dan `iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"`.
- **Proteksi Double Submit**: Tombol hapus dan batal wajib `disabled={isDeleting}` dengan indikator spinner.
- **Inline Error Feedback**: Jika server menolak penghapusan (e.g. HTTP 422 karena masih memiliki dependensi), dialog **DILARANG ditutup mendadak**, tampilkan error secara inline di dalam modal.
- **Proteksi Entitas Pokok**: Jika data adalah entitas bawaan sistem (e.g. Akun Level 1, Role ADMIN), tombol hapus **WAJIB dinonaktifkan** (`disabled={isDeleting || isProtected}`) dan disertai kartu peringatan edukatif (`ShieldAlert`).

---

## 3. Standar Form, Input & Selection

### 3.1. Dropdown & Selection
- **Opsi Biasa**: **WAJIB** menggunakan `M3Select.svelte`. DILARANG menggunakan `<select>` dan `<option>` native karena elemen option OS tidak mendukung `border-radius` dan merusak estetika rounded-2xl M3.
- **Toggle Boolean**: **WAJIB** menggunakan `M3ToggleCard.svelte`. Gunakan `colorVariant="emerald"` untuk status aktif operasional, dan `colorVariant="primary"` untuk aturan fungsional.
- **Master Data Skala Besar**: Entitas yang berpotensi berjumlah ratusan/ribuan wajib memakai lookup component:
  - Dokter / Nakes: `PractitionerLookupInput.svelte`
  - Akun Buku Besar: `CoaLookupInput.svelte`
  - Tarif Layanan: `TariffLookupInput.svelte`
  *DILARANG* memuat ratusan data ke dalam memori client atau array statis; wajib menggunakan server-side search (debounce 300ms) dan on-demand pagination.

### 3.2. Form Layout & Styling
- **Grid Layout**: Gunakan `grid grid-cols-1 sm:grid-cols-2 gap-4` untuk form 2 kolom.
- **Label**: `block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1`. Beri tanda wajib: `<span class="text-rose-500">*</span>`.
- **Text Input & Textarea**:
  `w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all`
- **Kode / Referensi / ID**: Wajib tambahkan `font-mono uppercase font-bold`.
- **Error Banner**:
  ```svelte
  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <span class="font-bold">Gagal:</span>
      <span>{errorMessage}</span>
    </div>
  {/if}
  ```

---

## 4. Batasan Ukuran File & Dekomposisi Komponen

- **Batas Ukuran File**: Satu file komponen Svelte atau TypeScript **ditargetkan tidak melebihi ~500 baris kode**.
- Aturan ini adalah **pedoman kualitas dan keterbacaan, bukan batasan kaku** yang memaksa pemecahan buatan. Jika sebuah halaman indeks mengorkestrasikan banyak dialog dan belum ada abstraksi yang lebih sederhana, membiarkannya tetap utuh lebih diutamakan daripada memecahnya menjadi file-file kecil yang justru menyulitkan pelacakan alur.
- **Strategi Pemecahan Alami**:
  1. **Tab Views**: Pisahkan konten per tab menjadi komponen mandiri (e.g. `UsersTab.svelte`, `RolesTab.svelte`, `PermissionsTab.svelte`).
  2. **Sub-Modals**: Ekstrak form add/edit dan dialog delete ke file terpisah (e.g. `UserFormModal.svelte`, `UserDeleteModal.svelte`).
  3. **Complex Rows/Cards**: Ekstrak card iterasi tabel/matriks yang rumit ke komponen card mandiri.
  4. **Konstanta & Struktur**: Pindahkan array menu/definisi statis ke file `.ts` pendukung.

---

## 5. Standar Svelte 5 & TypeScript

- **Gunakan Svelte 5 Runes Eksklusif**:
  - State: `let val = $state(initial);`
  - Derived: `let computed = $derived(expr);` atau `$derived.by(() => { ... })`
  - Props: `let { propA, propB = defaultVal }: Props = $props();`
  - Side effect: `$effect(() => { ... });`
- **DILARANG** sintaks lama: `export let foo;`, `$: bar = ...`, `<slot />`, `createEventDispatcher`.
- **TypeScript Hygiene**:
  - Selalu definisikan interface `Props` eksplisit untuk setiap komponen.
  - Impor tipe dari modul pemilik resminya di `$lib/types/...`. DILARANG mendeklarasikan ulang tipe domain di modul lain.

---

## 6. Standar Tabel, Action Menu & Tampilan "To The Point"

### 6.1. Action Menu pada Tabel (Anti-Clipping & Clean Table)
- **Prinsip Tampilan Bersih**: Kolom AKSI pada tabel maksimal menampilkan **1 tombol aksi cepat** (misalnya *View Detail*) dan **1 tombol opsi** (`MoreVertical`). Hindari menaruh 4–5 tombol berjejer di dalam baris karena membuat tabel bising secara visual (*cluttered*).
- **Anti-Clipping Container**: Dropdown menu opsi pada tabel **WAJIB** diposisikan dengan `position: fixed` (atau dirender di luar pembungkus tabel `overflow-x-auto`) menggunakan koordinat pemicu dari `button.getBoundingClientRect()`.
- **Smart Flipping**: Menu **WAJIB** mendeteksi ruang vertikal layar: jika posisi tombol berada di baris bawah dan ruang ke bawah tidak mencukupi, menu otomatis membuka ke atas (`bottom`) agar tidak terpotong tepi layar.
- **Auto-Dismiss**: Menu otomatis menutup saat pengguna mengklik di luar area, menekan tombol `Escape`, atau melakukan *scroll* pada jendela maupun container tabel (`onscroll={closeMenu}`).

### 6.2. Header & Filter Bar "To The Point" (Single-Row Bar + Status Pills)
- **Status Filter Pills**: Gunakan baris tombol pil horizontal ringkas dengan counter jumlah data (`Semua (X)`, `Aktif (Y)`, `Draft (Z)`).
- **Single-Row Filter Bar**: Satukan input pencarian multi-field (nomor, nama, relasi), selector filter utama (`M3Select`), dan info counter data dalam **satu kontainer baris tunggal (`h-10`)**.
- **DILARANG** membuat form filter bertingkat 2–3 baris yang memakan ruang vertikal layar sebelum tabel terlihat. Pengguna harus langsung disajikan data tabel tanpa terganggu oleh filter yang memakan tempat.

### 6.3. Prinsip "One Dialog per Intent" untuk Workflow Transaksional
- Pada entitas yang memiliki siklus hidup (*lifecycle* seperti DRAFT, ACTIVE, SUSPENDED, TERMINATED), **WAJIB** memisahkan dialog berdasarkan tujuan bisnis (*intent*):
  - Form dialog khusus: hanya untuk pembuatan baru dan revisi draft awal (`FooFormModal.svelte`).
  - Dialog transisi status: hanya meminta alasan/justifikasi aktivasi, penangguhan, atau pemutusan (`FooStatusModal.svelte`).
  - Dialog amandemen / perpanjangan: hanya menampilkan field yang diizinkan untuk diubah (`FooAmendModal.svelte`, `FooRenewModal.svelte`).
  - Dialog pemetaan / aturan khusus: dialog terpisah untuk konfigurasi relasi sekunder (`FooPriceSchemasModal.svelte`).
- **Simplicity Advantage**: Form menjadi pendek, fokus, dan tidak membingungkan pengguna dengan puluhan field nonaktif yang tidak relevan dengan tindakan saat itu.

### 6.4. Feedback Operasional yang Ringkas (Floating Toast vs Layout Shift)
- **Respon Operasi CRUD**: Gunakan floating toast notification di sudut kanan bawah (`fixed bottom-6 right-6 z-50`) dengan waktu *auto-dismiss* 3–4 detik untuk konfirmasi keberhasilan aksi (simpan, ubah status, hapus).
- **DILARANG** menampilkan banner alert sukses/gagal di atas tabel yang menyebabkan pergeseran layout halaman (*layout shift*).
- **Inline Error**: Banner error inline (`bg-rose-50 ...`) hanya digunakan di dalam dialog modal saat proses validasi atau panggilan API gagal, sehingga pengguna dapat langsung memperbaiki input tanpa kehilangan konteks.
