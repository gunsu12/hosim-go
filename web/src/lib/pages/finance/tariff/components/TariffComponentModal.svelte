<script lang="ts">
  import { Layers, Building2, UserCheck, HeartHandshake, CheckCircle2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3ToggleCard from '$lib/components/m3/M3ToggleCard.svelte';
  import CoaLookupInput from '$lib/components/CoaLookupInput.svelte';
  import type {
    TariffComponentRecord,
    CreateTariffComponentDTO,
    UpdateTariffComponentDTO
  } from '$lib/types/finance/tariff';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
    tariffComponent: TariffComponentRecord | null;
    isSaving: boolean;
    onClose: () => void;
    onSave: (payload: CreateTariffComponentDTO | UpdateTariffComponentDTO) => Promise<void>;
  }

  let {
    isOpen,
    mode,
    tariffComponent,
    isSaving,
    onClose,
    onSave
  }: Props = $props();

  let code = $state('');
  let name = $state('');
  let componentType = $state('JASA_MEDIS');
  let isHospitalRevenue = $state(false);
  let isOperatorRevenue = $state(true);
  let isParamedicRevenue = $state(false);
  let defaultCoaCode = $state('');
  let description = $state('');
  let isActive = $state(true);
  let formError = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      if (mode === 'EDIT' && tariffComponent) {
        code = tariffComponent.code;
        name = tariffComponent.name;
        componentType = tariffComponent.component_type || 'LAINNYA';
        isHospitalRevenue = tariffComponent.is_hospital_revenue;
        isOperatorRevenue = tariffComponent.is_operator_revenue;
        isParamedicRevenue = tariffComponent.is_paramedic_revenue;
        defaultCoaCode = tariffComponent.default_coa_code || '';
        description = tariffComponent.description || '';
        isActive = tariffComponent.is_active;
      } else {
        code = '';
        name = '';
        componentType = 'JASA_MEDIS';
        isHospitalRevenue = false;
        isOperatorRevenue = true;
        isParamedicRevenue = false;
        defaultCoaCode = '';
        description = '';
        isActive = true;
      }
      formError = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!code.trim()) {
      formError = 'Kode komponen tarif wajib diisi';
      return;
    }
    if (!name.trim()) {
      formError = 'Nama komponen tarif wajib diisi';
      return;
    }

    formError = null;
    await onSave({
      code: code.trim().toUpperCase(),
      name: name.trim(),
      component_type: componentType,
      is_hospital_revenue: isHospitalRevenue,
      is_operator_revenue: isOperatorRevenue,
      is_paramedic_revenue: isParamedicRevenue,
      default_coa_code: defaultCoaCode.trim() || undefined,
      description: description.trim() || undefined,
      is_active: isActive
    });
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'CREATE' ? 'Tambah Komponen Tarif Baru' : 'Edit Komponen Tarif'}
  subtitle="Elemen pemecahan tagihan biaya (jasa dokter, sarana RS, BHP)"
  icon={Layers}
  iconClass="bg-emerald-500/10 text-emerald-600 dark:bg-emerald-500/20 dark:text-emerald-400"
  size="xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if formError}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <span class="font-bold">Gagal:</span>
      <span>{formError}</span>
    </div>
  {/if}

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="tc-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Kode Komponen <span class="text-rose-500">*</span>
      </label>
      <input
        id="tc-code"
        type="text"
        bind:value={code}
        placeholder="Misal: JM-001, JS-RS, BHP-LAB"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div>
      <label for="tc-type" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Tipe Komponen <span class="text-rose-500">*</span>
      </label>
      <select
        id="tc-type"
        bind:value={componentType}
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      >
        <option value="JASA_MEDIS">Jasa Medis (Dokter / Nakes)</option>
        <option value="JASA_RS">Jasa Sarana Rumah Sakit</option>
        <option value="SEWA_ALAT">Sewa Alat & Aset Medis</option>
        <option value="BAHAN_ALKES">Bahan Habis Pakai & Alkes</option>
        <option value="ADMINISTRASI">Administrasi & Rekam Medis</option>
        <option value="LAINNYA">Lain-lain</option>
      </select>
    </div>
  </div>

  <div>
    <label for="tc-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Nama Komponen Tarif <span class="text-rose-500">*</span>
    </label>
    <input
      id="tc-name"
      type="text"
      bind:value={name}
      placeholder="Contoh: Jasa Medis Dokter Spesialis Penanggung Jawab"
      required
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    />
  </div>

  <!-- Revenue Distribution Checklist -->
  <div class="p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] flex flex-col gap-2">
    <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs block">
      Alokasi Distribusi Penerima Pendapatan:
    </span>
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
      <label class="flex items-center gap-2.5 p-2.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] cursor-pointer hover:border-blue-400 dark:hover:border-blue-500 transition-colors">
        <input type="checkbox" bind:checked={isHospitalRevenue} class="rounded text-[#0b57d0]" />
        <div class="flex items-center gap-1.5 text-xs text-[#1f1f1f] dark:text-[#e3e3e3] font-semibold">
          <Building2 class="w-4 h-4 text-blue-600 dark:text-blue-400" />
          <span>Rumah Sakit</span>
        </div>
      </label>

      <label class="flex items-center gap-2.5 p-2.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] cursor-pointer hover:border-emerald-400 dark:hover:border-emerald-500 transition-colors">
        <input type="checkbox" bind:checked={isOperatorRevenue} class="rounded text-[#0b57d0]" />
        <div class="flex items-center gap-1.5 text-xs text-[#1f1f1f] dark:text-[#e3e3e3] font-semibold">
          <UserCheck class="w-4 h-4 text-emerald-600 dark:text-emerald-400" />
          <span>Dokter</span>
        </div>
      </label>

      <label class="flex items-center gap-2.5 p-2.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] cursor-pointer hover:border-amber-400 dark:hover:border-amber-500 transition-colors">
        <input type="checkbox" bind:checked={isParamedicRevenue} class="rounded text-[#0b57d0]" />
        <div class="flex items-center gap-1.5 text-xs text-[#1f1f1f] dark:text-[#e3e3e3] font-semibold">
          <HeartHandshake class="w-4 h-4 text-amber-600 dark:text-amber-400" />
          <span>Paramedis</span>
        </div>
      </label>
    </div>
  </div>

  <!-- Dynamic COA Account Lookup -->
  <CoaLookupInput
    id="tc-coa"
    label="Default Pemetaan Akun Bagan Akun (COA)"
    bind:value={defaultCoaCode}
    placeholder="Cari atau pilih akun pendapatan COA (4xx)..."
    preferredType="REVENUE"
    helperText="Akun buku besar yang akan dikreditkan saat kasir memproses tagihan komponen ini."
  />

  <div>
    <label for="tc-desc" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Deskripsi Tambahan
    </label>
    <textarea
      id="tc-desc"
      bind:value={description}
      rows="2"
      placeholder="Catatan aturan honorarium, sarana kamar, atau ketentuan khusus..."
      class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
    ></textarea>
  </div>

  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
    <M3ToggleCard
      title="Komponen Tarif Aktif Digunakan"
      description="Nonaktifkan jika komponen pemecahan tarif ini sudah tidak lagi dihitung dalam penagihan."
      bind:checked={isActive}
      colorVariant="emerald"
    />
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSaving}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSaving}>
      {#if isSaving}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>{mode === 'CREATE' ? 'Simpan Komponen' : 'Perbarui Komponen'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
