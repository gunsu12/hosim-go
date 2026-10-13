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
| **Batas Ukuran File** | **Maksimal ~500 baris kode per file** | Menumpuk semua tab/modal dalam 1 file monolitik > 500 baris |
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

- **Batas Ukuran File**: Satu file komponen Svelte atau TypeScript **WAJIB ditargetkan tidak melebihi ~500 baris kode**.
- **Strategi Pemecahan**:
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
