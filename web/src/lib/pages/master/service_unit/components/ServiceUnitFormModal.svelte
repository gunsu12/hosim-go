<script lang="ts">
  import { Stethoscope, Check, Loader2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select from '$lib/components/m3/M3Select.svelte';
  import M3ToggleCard from '$lib/components/m3/M3ToggleCard.svelte';
  import { createServiceUnit, updateServiceUnit } from '$lib/api/organization';
  import type {
    ServiceUnitRecord,
    CreateServiceUnitDTO,
    UpdateServiceUnitDTO
  } from '../../../../types/master/service_unit';
  import type { DepartementRecord } from '../../../../types/master/departement';

  interface Props {
    isOpen: boolean;
    mode: 'create' | 'edit';
    unit?: ServiceUnitRecord | null;
    departments?: DepartementRecord[];
    onClose: () => void;
    onSuccess: (unit: ServiceUnitRecord) => void;
  }

  let {
    isOpen = false,
    mode = 'create',
    unit = null,
    departments = [],
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  let formCode = $state('');
  let formName = $state('');
  let formDeptId = $state('');
  let formPhone = $state('');
  let formEmail = $state('');
  let formAddress = $state('');
  let formIhsLocId = $state('');
  let formIsRegistrationTarget = $state(true);
  let formIsActive = $state(true);
  let formError = $state<string | null>(null);
  let isSaving = $state(false);

  $effect(() => {
    if (isOpen) {
      formError = null;
      isSaving = false;
      if (mode === 'edit' && unit) {
        formCode = unit.code;
        formName = unit.name;
        formDeptId = unit.departement_id || '';
        formPhone = unit.phone || '';
        formEmail = unit.email || '';
        formAddress = unit.address || '';
        formIhsLocId = unit.ihs_location_id || '';
        formIsRegistrationTarget = unit.is_registration_target;
        formIsActive = unit.is_active;
      } else {
        formCode = '';
        formName = '';
        formDeptId = '';
        formPhone = '';
        formEmail = '';
        formAddress = '';
        formIhsLocId = '';
        formIsRegistrationTarget = true;
        formIsActive = true;
      }
    }
  });

  async function handleSave() {
    formError = null;
    if (!formCode.trim()) {
      formError = 'Kode unit layanan wajib diisi';
      return;
    }
    if (!formName.trim()) {
      formError = 'Nama unit layanan wajib diisi';
      return;
    }

    isSaving = true;
    try {
      const payload: CreateServiceUnitDTO = {
        code: formCode.trim().toUpperCase(),
        name: formName.trim(),
        departement_id: formDeptId || undefined,
        phone: formPhone.trim() || undefined,
        email: formEmail.trim() || undefined,
        address: formAddress.trim() || undefined,
        ihs_location_id: formIhsLocId.trim() || undefined,
        is_registration_target: formIsRegistrationTarget,
        is_active: formIsActive
      };

      let savedUnit: ServiceUnitRecord;
      if (mode === 'create') {
        savedUnit = await createServiceUnit(payload);
      } else if (unit) {
        const updatePayload: UpdateServiceUnitDTO = { ...payload };
        savedUnit = await updateServiceUnit(unit.id, updatePayload);
      } else {
        return;
      }

      onSuccess(savedUnit);
      onClose();
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan unit layanan';
    } finally {
      isSaving = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'create' ? 'Tambah Unit Layanan / Poli Baru' : 'Edit Unit Layanan'}
  subtitle="Master poliklinik, instalasi, dan unit pelayanan rawat jalan / inap"
  icon={Stethoscope}
  size="2xl"
  {onClose}
  onsubmit={handleSave}
>
  <div class="space-y-4 text-xs">
    <!-- Error Feedback -->
    {#if formError}
      <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
        <span class="font-bold">Gagal:</span>
        <span>{formError}</span>
      </div>
    {/if}

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label for="su-form-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
          Kode Unit <span class="text-rose-500">*</span>
        </label>
        <input
          id="su-form-code"
          type="text"
          bind:value={formCode}
          placeholder="Contoh: POLI-BEDAH"
          disabled={mode === 'edit'}
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] disabled:opacity-60 disabled:cursor-not-allowed focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <M3Select
          id="su-form-dept"
          label="Induk Departemen"
          bind:value={formDeptId}
          placeholder="-- Tanpa Induk Departemen --"
          options={[
            { value: '', label: '-- Tanpa Induk Departemen --' },
            ...departments.map((d) => ({
              value: d.id,
              label: `${d.name} (${d.code})`
            }))
          ]}
        />
      </div>
    </div>

    <div>
      <label for="su-form-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Nama Unit Layanan / Poli <span class="text-rose-500">*</span>
      </label>
      <input
        id="su-form-name"
        type="text"
        bind:value={formName}
        placeholder="Contoh: Poliklinik Bedah Umum"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label for="su-form-phone" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
          Telepon / Ext
        </label>
        <input
          id="su-form-phone"
          type="text"
          bind:value={formPhone}
          placeholder="Contoh: Ext. 101"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <label for="su-form-email" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
          Email Unit
        </label>
        <input
          id="su-form-email"
          type="email"
          bind:value={formEmail}
          placeholder="Contoh: bedah@hosim.local"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>
    </div>

    <div>
      <label for="su-form-ihs" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        SATUSEHAT Location ID
      </label>
      <input
        id="su-form-ihs"
        type="text"
        bind:value={formIhsLocId}
        placeholder="Contoh: LOC-3273-0101"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-1 block">
        ID Lokasi faskes terdaftar di platform SatuSehat Kemenkes
      </span>
    </div>

    <div>
      <label for="su-form-address" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Lokasi / Ruang / Lantai
      </label>
      <input
        id="su-form-address"
        type="text"
        bind:value={formAddress}
        placeholder="Contoh: Lantai 2 Sayap Barat Kamar 204"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="pt-2 space-y-3">
      <M3ToggleCard
        bind:checked={formIsRegistrationTarget}
        title="Tujuan Pendaftaran Pasien"
        subtitle="Dapat dipilih sebagai poli/tujuan saat pasien mendaftar di loket admisi maupun online."
        colorVariant="primary"
      />

      <M3ToggleCard
        bind:checked={formIsActive}
        title="Status Unit Aktif"
        subtitle="Unit layanan aktif dan melayani operasional medis harian."
        colorVariant="emerald"
      />
    </div>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSaving}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSaving}>
      {#if isSaving}
        <Loader2 class="w-4 h-4 mr-1.5 animate-spin" />
        <span>Menyimpan...</span>
      {:else}
        <Check class="w-4 h-4 mr-1.5" />
        <span>{mode === 'create' ? 'Simpan Unit Layanan' : 'Perbarui Unit Layanan'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
