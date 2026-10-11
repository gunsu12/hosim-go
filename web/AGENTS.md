# AGENTS.md — HOSIM-GO Web Frontend Engineering Rules

Dokumen ini berlaku untuk seluruh folder `web/` di repositori HOSIM-GO. Dokumen ini melengkapi `AGENTS.md` di root repositori dengan aturan spesifik untuk arsitektur frontend, UI/UX, Svelte 5, TypeScript, dan Google Material Design 3 (M3).

Semua coding agent **WAJIB** membaca dan mematuhi aturan ini sebelum membuat, mengubah, atau me-refactor komponen dan halaman di `web/`.

Kata **WAJIB**, **DILARANG**, dan **BOLEH** bersifat normatif.

---

## 1. Filosofi & Design System (Google Material Design 3)

Antarmuka HOSIM-GO mengadopsi **Google Workspace / Drive Authentic Material Design 3 (M3)**:
- **Clean, Spacious, Soft-curved**: Menggunakan border-radius M3 yang konsisten (`rounded-3xl` untuk dialog/modal permukaan, `rounded-2xl` untuk cards & avatar squircle, `rounded-xl` untuk form inputs, `rounded-full` untuk pills & buttons).
- **Surface Elevation & Tokens**: Mengikuti token permukaan M3 di `src/app.css` (`--md-surface`, `--md-surface-container`, `--md-primary`, dll.).
- **Dukungan Dark Mode Penuh**: Semua komponen, form field, teks, border, dan modal **WAJIB** memiliki class Dark Mode (`dark:...`) yang adaptif. DILARANG menggunakan warna hardcoded tanpa padanan dark mode.

---

## 2. Standar Form & Modal (M3 Dialog Standards)

Pola desain pada `CustomerModal.svelte` (`$lib/pages/finance/customer/components/CustomerModal.svelte`) adalah **benchmark resmi (Gold Standard)** untuk semua form dan dialog modal di HOSIM-GO.

### 2.1. Aturan Penggunaan Komponen Modal

- **WAJIB** menggunakan [`M3Modal.svelte`](file:///c:/laragon/www/hosim-go/web/src/lib/components/m3/M3Modal.svelte) sebagai container dialog (baik form mutasi data maupun detail viewer).
- **DILARANG** menduplikasi markup backdrop (`fixed inset-0 bg-black/40...`), frame dialog card, close button `X`, dan ESC event listener secara manual di halaman/komponen baru.
- **WAJIB** menyertakan:
  - `title`: Judul modal yang ringkas dan jelas.
  - `subtitle`: Keterangan konteks atau kode referensi.
  - `icon`: Icon komponen Lucide yang representatif (squircle avatar 40x40px).
  - `onClose`: Callback fungsi untuk menutup modal.
  - `onsubmit`: Handler submit form (bila modal bertipe form).
- **Pilihan Ukuran Modal (`size`)**:
  - `'sm'` / `'md'`: Konfirmasi aksi singkat atau input 1–2 field sederhana.
  - `'lg'` / `'xl'`: Form ringkas (3–5 field).
  - `'2xl'` (**Default**): Form master data standar, konfigurasi entitas, atau form transaksi standar.
  - `'3xl'` / `'4xl'`: Form kompleks dengan multi-tab atau embedded data table / order items.

### 2.2. Aturan Form Inputs & Grid Layout

- **Field Layout**:
  - Gunakan grid responsif `grid grid-cols-1 sm:grid-cols-2 gap-4` untuk penataan field 2 kolom.
  - Field tunggal lebar penuh (nama instansi, alamat, deskripsi) ditaruh di luar grid atau span 2 kolom.
- **Label**:
  - Class: `block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1`.
  - Tanda wajib diisi: `<span class="text-rose-500">*</span>`.
- **Input, Select & Textarea**:
  - Text Input / Select:  
    `w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all`
  - Kode / Nomor Referensi / NIK / ID:  
    Wajib tambahkan `font-mono uppercase font-bold`.
  - Textarea:  
    `w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all`
  - Helper / Deskripsi di bawah field:  
    `<p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">`
- **Section Grouping (Pembatas Kelompok)**:
  - Pisahkan kelompok logis form dengan divider:  
    `<div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">`
  - Heading section:  
    `<h4 class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider mb-3">`

### 2.3. Aturan Konfigurasi Opsi & Switch Boolean

- **WAJIB** menggunakan [`M3ToggleCard.svelte`](file:///c:/laragon/www/hosim-go/web/src/lib/components/m3/M3ToggleCard.svelte) untuk konfigurasi aturan sistem, status aktif/nonaktif, atau opsi boolean lainnya.
- **DILARANG** membuat checkbox input raw HTML atau styling toggle manual berulang kali.
- Gunakan `colorVariant`:
  - `'primary'`: Konfigurasi aturan fungsional (e.g. Wajib Kartu Fisik, Wajib Verifikasi).
  - `'emerald'`: Status operasional aktif/nonaktif.
  - `'amber'`: Opsi bersyarat / peringatan.
  - `'rose'`: Flag bahaya / high alert / isolasi.

### 2.4. Aturan Tombol Aksi (Footer Actions)

- Letakkan tombol aksi di dalam `{#snippet footer()}` pada `M3Modal`.
- **Tombol Batal**: `<M3Button variant="text" onclick={onClose} disabled={isSubmitting}>Batal</M3Button>`.
- **Tombol Submit**: `<M3Button variant="filled" type="submit" disabled={isSubmitting}>` dengan spinner indicator saat `isSubmitting === true` dan icon Lucide (`CheckCircle2`, `Save`, dsb.) saat normal.

### 2.5. Error Handling & Validation Feedback

- Validasi format dasar dilakukan di sisi client sebelum pemanggilan API.
- Tampilkan error banner di bagian atas isi form:
  ```svelte
  {#if errorMessage}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <span class="font-bold">Gagal:</span>
      <span>{errorMessage}</span>
    </div>
  {/if}
  ```
- Reset `errorMessage = null` setiap kali modal dibuka (`isOpen`).

### 2.4. Standar Dialog Konfirmasi Destruktif (Delete Confirmation Modal)

Untuk dialog konfirmasi penghapusan data (seperti `CoaDeleteModal.svelte`):
- **Ukuran & Ikon**: Gunakan `size="md"` dengan `icon={Trash2}` dan `iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"`.
- **Proteksi Double-Submission**: Tombol hapus dan batal **WAJIB** memiliki `disabled={isDeleting}` serta menampilkan ikon loading spinner (`RefreshCw class="animate-spin"`) saat penghapusan sedang berlangsung.
- **Inline Error Feedback**: Jika server menolak penghapusan (misal: HTTP 422 karena masih memiliki turunan anak atau keterikatan tarif/transaksi), modal **DILARANG ditutup mendadak**. Dialog harus tetap terbuka dan menampilkan `errorMessage` secara inline di dalam modal agar pengguna langsung mengetahui alasan penolakan tanpa harus mencari banner error di posisi scroll lain.
- **Proteksi Entitas Pokok Sistem**: Jika data yang hendak dihapus merupakan entitas klasifikasi sistem pokok yang dilindungi (misal: Akun Buku Besar Level 1 atau Role sistem bawaan), tombol hapus **WAJIB dinonaktifkan** (`disabled={isDeleting || isLevelOne}`) dan disertai kartu peringatan edukatif (`ShieldAlert`).

---

## 3. Standar Svelte 5 & TypeScript

Proyek frontend HOSIM-GO menggunakan **Svelte 5 Runes**:

### 3.1. Gunakan Svelte 5 Runes Eksklusif
- **State**: Gunakan `let foo = $state(val);`.
- **Derived**: Gunakan `let bar = $derived(expr);` atau `$derived.by(() => { ... })`.
- **Props**: Gunakan `let { propA, propB = defaultVal }: Props = $props();`.
- **Side Effect**: Gunakan `$effect(() => { ... });`.
- **DILARANG** menggunakan sintaks Svelte 3/4 lama:
  - ❌ `export let foo;`
  - ❌ `$: bar = foo * 2;`
  - ❌ `<slot />` atau `$$slots` (Gunakan Svelte 5 Snippets: `Snippet`, `{#snippet ...}`, `{@render ...}`).
  - ❌ `createEventDispatcher` (Gunakan callback props: `onclick`, `onClose`, `onSuccess`).

### 3.2. TypeScript & Data Contracts
- Buat interface eksplisit untuk `Props` setiap komponen.
- Model dan DTO **WAJIB** diimpor dari `$lib/types/...` (misalnya `$lib/types/finance`, `$lib/types/master/...`).
- Hindari penggunaan tipe `any` kecuali untuk payload eksternal yang belum terdefinisi.

---

## 4. Template Baku Komponen Modal Baru

Gunakan template di bawah ini saat membuat modal baru:

```svelte
<script lang="ts">
  import { ShieldCheck, CheckCircle2 } from "@lucide/svelte";
  import M3Modal from "$lib/components/m3/M3Modal.svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3ToggleCard from "$lib/components/m3/M3ToggleCard.svelte";
  import type { ExampleRecord, CreateExampleDTO, UpdateExampleDTO } from "$lib/types/...";

  interface Props {
    isOpen: boolean;
    mode: "CREATE" | "EDIT";
    item?: ExampleRecord | null;
    onClose: () => void;
    onSuccess: (saved: ExampleRecord) => void;
  }

  let {
    isOpen = false,
    mode = "CREATE",
    item = null,
    onClose = () => {},
    onSuccess = () => {},
  }: Props = $props();

  let code = $state("");
  let name = $state("");
  let isActive = $state(true);
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      if (mode === "EDIT" && item) {
        code = item.code || "";
        name = item.name || "";
        isActive = item.is_active;
      } else {
        code = "";
        name = "";
        isActive = true;
      }
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    errorMessage = null;

    if (!code.trim()) {
      errorMessage = "Kode wajib diisi";
      return;
    }
    if (!name.trim()) {
      errorMessage = "Nama wajib diisi";
      return;
    }

    isSubmitting = true;
    try {
      if (mode === "CREATE") {
        // Panggil create API...
      } else {
        // Panggil update API...
      }
      onSuccess(...);
    } catch (err: any) {
      errorMessage = err?.message || "Terjadi kesalahan saat menyimpan data";
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === "CREATE" ? "Tambah Data Baru" : "Ubah Data"}
  subtitle={mode === "CREATE" ? "Konfigurasi master data baru" : `Kode: ${item?.code || "-"}`}
  icon={ShieldCheck}
  size="2xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if errorMessage}
    <div
      class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5"
    >
      <span class="font-bold">Gagal:</span>
      <span>{errorMessage}</span>
    </div>
  {/if}

  <!-- Section: Identitas Utama (2 Kolom) -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="modal-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Kode <span class="text-rose-500">*</span>
      </label>
      <input
        id="modal-code"
        type="text"
        bind:value={code}
        placeholder="Contoh: KODE-01"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
        Harus unik di dalam sistem.
      </p>
    </div>

    <div>
      <label for="modal-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Nama <span class="text-rose-500">*</span>
      </label>
      <input
        id="modal-name"
        type="text"
        bind:value={name}
        placeholder="Contoh: Nama Entitas"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <!-- Section: Kebijakan & Status -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738] space-y-3">
    <h4 class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider mb-2">
      Status Operasional
    </h4>

    <M3ToggleCard
      title="Status Aktif"
      description="Nonaktifkan bila data tidak boleh dipilih pada transaksi pelayanan baru."
      bind:checked={isActive}
      colorVariant="emerald"
    />
  </div>

  <!-- Sticky Footer Action Buttons -->
  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSubmitting}>
      {#if isSubmitting}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>{mode === "CREATE" ? "Simpan Data" : "Perbarui Data"}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
```

---

## 4. Batasan Ukuran File & Dekomposisi Komponen (Maksimal 500 Baris)

Sesuai standar pemeliharaan kode (maintainability) HOSIM-GO:

- **Target Ukuran File**: Satu file komponen Svelte atau TypeScript **WAJIB ditargetkan tidak melebihi ~500 baris kode**.
- **Strategi Dekomposisi**: Jika sebuah file mulai mendekati atau melampaui batas ini, agent **WAJIB** memecahnya ke arsitektur komponen modular:
  1. **Tab & View Pages**: Pisahkan konten per tab menjadi komponen tersendiri (misal: `PricePlanTab.svelte`, `TariffClassTab.svelte`, `TariffComponentTab.svelte` pada `IndexPage.svelte`).
  2. **Sub-Modal & Nested Forms**: Ekstrak sub-modal atau form anak ke file komponen terpisah (misal: `PricePlanItemFormModal.svelte` diekstrak dari `PricePlanItemsModal.svelte`).
  3. **Complex Repeater Rows/Cards**: Ekstrak kartu iterasi tabel/matriks yang rumit menjadi card component tersendiri (misal: `TariffMatrixClassCard.svelte`).
  4. **Data Struktur & Menu Statis**: Pindahkan array navigasi, definisi modul, atau konstanta jenis akun ke file TypeScript murni (misal: `navigationStructure.ts`, `launcherModules.ts`, `constants.ts`).
  5. **Helper & Filter Logic**: Pindahkan algoritma rekursif atau penyaringan tree ke utilitas pendukung (misal: `coaUtils.ts`).

---

## 5. TypeScript Types, Domain Barrel, & Export Hygiene

- **Ambiguitas Re-export Barrel**: Pada barrel pusat (`$lib/types/index.ts`), **DILARANG** mengekspor ganda interface dengan nama identik dari domain berbeda tanpa alias yang jelas.
- **Kepemilikan Domain Tipe**: Tipe domain yang sudah memiliki modul pemilik resmi (misalnya `TariffClassRecord` yang merupakan domain `finance`) **DILARANG** dideklarasikan atau diekspor ulang dari modul lain (seperti `master`). Seluruh komponen dan consumer wajib mengimpor dari modul pemilik definitifnya.


