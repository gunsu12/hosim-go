<script lang="ts">
  import { X, Layers, Building2, UserCheck, HeartHandshake } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
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
      code: code.trim(),
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

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-lg p-6 flex flex-col gap-4 max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-emerald-50 text-emerald-700 flex items-center justify-center font-bold">
            <Layers class="w-4 h-4" />
          </div>
          <div>
            <h3 class="text-sm font-bold text-[#1f1f1f]">
              {mode === 'CREATE' ? 'Tambah Komponen Tarif Baru' : 'Edit Komponen Tarif'}
            </h3>
            <p class="text-[11px] text-[#747775]">Elemen pemecahan tagihan biaya (jasa dokter, sarana RS, BHP)</p>
          </div>
        </div>
        <button
          onclick={onClose}
          type="button"
          class="text-[#747775] hover:text-[#1f1f1f] p-1 rounded-full hover:bg-[#f0f4f9] transition-colors"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body Form -->
      <form onsubmit={handleSubmit} class="flex flex-col gap-3 text-xs">
        {#if formError}
          <div class="p-2.5 rounded-xl bg-[#fce8e6] text-[#c5221f] text-xs font-semibold">
            {formError}
          </div>
        {/if}

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div>
            <label for="tc-code" class="block font-semibold text-[#444746] mb-1">Kode Komponen *</label>
            <input
              id="tc-code"
              type="text"
              bind:value={code}
              placeholder="Misal: JM^001, JS-RS, BHP-LAB"
              required
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0] focus:outline-hidden"
            />
          </div>

          <div>
            <label for="tc-type" class="block font-semibold text-[#444746] mb-1">Tipe Komponen *</label>
            <select
              id="tc-type"
              bind:value={componentType}
              class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
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
          <label for="tc-name" class="block font-semibold text-[#444746] mb-1">Nama Komponen Tarif *</label>
          <input
            id="tc-name"
            type="text"
            bind:value={name}
            placeholder="Contoh: Jasa Medis Dokter Spesialis Penanggung Jawab"
            required
            class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden"
          />
        </div>

        <!-- Revenue Distribution Checklist -->
        <div class="p-3.5 rounded-2xl bg-[#f8fafd] border border-[#e1e5ea] flex flex-col gap-2">
          <span class="font-semibold text-[#1f1f1f] text-[11px] block">
            Alokasi Distribusi Penerima Pendapatan:
          </span>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <label class="flex items-center gap-2 p-2 rounded-xl bg-white border border-[#e1e5ea] cursor-pointer hover:border-blue-300">
              <input type="checkbox" bind:checked={isHospitalRevenue} class="rounded text-[#0b57d0]" />
              <div class="flex items-center gap-1.5 text-xs text-[#1f1f1f] font-medium">
                <Building2 class="w-3.5 h-3.5 text-blue-600" />
                <span>Rumah Sakit</span>
              </div>
            </label>

            <label class="flex items-center gap-2 p-2 rounded-xl bg-white border border-[#e1e5ea] cursor-pointer hover:border-emerald-300">
              <input type="checkbox" bind:checked={isOperatorRevenue} class="rounded text-[#0b57d0]" />
              <div class="flex items-center gap-1.5 text-xs text-[#1f1f1f] font-medium">
                <UserCheck class="w-3.5 h-3.5 text-emerald-600" />
                <span>Dokter</span>
              </div>
            </label>

            <label class="flex items-center gap-2 p-2 rounded-xl bg-white border border-[#e1e5ea] cursor-pointer hover:border-amber-300">
              <input type="checkbox" bind:checked={isParamedicRevenue} class="rounded text-[#0b57d0]" />
              <div class="flex items-center gap-1.5 text-xs text-[#1f1f1f] font-medium">
                <HeartHandshake class="w-3.5 h-3.5 text-amber-600" />
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
          <label for="tc-desc" class="block font-semibold text-[#444746] mb-1">Deskripsi Tambahan</label>
          <textarea
            id="tc-desc"
            bind:value={description}
            rows="2"
            placeholder="Catatan aturan honorarium, sarana kamar, atau ketentuan khusus..."
            class="w-full p-2.5 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0] focus:outline-hidden resize-none"
          ></textarea>
        </div>

        <div class="pt-1">
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" bind:checked={isActive} class="rounded text-[#0b57d0]" />
            <span class="font-medium text-[#1f1f1f]">Komponen Tarif Aktif Digunakan</span>
          </label>
        </div>

        <!-- Footer Actions -->
        <div class="flex justify-end gap-2 pt-3 border-t border-[#e1e5ea] mt-1">
          <button
            type="button"
            onclick={onClose}
            class="px-3.5 py-1.5 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] hover:bg-[#f0f4f9] transition-colors"
          >
            Batal
          </button>
          <M3Button variant="filled" type="submit" disabled={isSaving}>
            <span>{mode === 'CREATE' ? 'Simpan Komponen' : 'Perbarui Komponen'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
