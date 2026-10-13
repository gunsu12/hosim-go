<script lang="ts">
  import { Building2, Check, Loader2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select from '$lib/components/m3/M3Select.svelte';
  import M3ToggleCard from '$lib/components/m3/M3ToggleCard.svelte';
  import { createDepartment, updateDepartment } from '$lib/api/organization';
  import type {
    DepartementRecord,
    DepartementType,
    CreateDepartementDTO,
    UpdateDepartementDTO
  } from '../../../../types/master/departement';

  interface Props {
    isOpen: boolean;
    mode: 'create' | 'edit';
    department?: DepartementRecord | null;
    onClose: () => void;
    onSuccess: (dept: DepartementRecord) => void;
  }

  let {
    isOpen = false,
    mode = 'create',
    department = null,
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  let formCode = $state('');
  let formName = $state('');
  let formType = $state<DepartementType>('outpatient');
  let formPhone = $state('');
  let formEmail = $state('');
  let formAddress = $state('');
  let formIhsOrgId = $state('');
  let formIsActive = $state(true);
  let formError = $state<string | null>(null);
  let isSaving = $state(false);

  const deptTypeOptions: { value: DepartementType; label: string }[] = [
    { value: 'outpatient', label: 'Rawat Jalan (Outpatient)' },
    { value: 'inpatient', label: 'Rawat Inap (Inpatient)' },
    { value: 'emergency', label: 'Gawat Darurat (Emergency / IGD)' },
    { value: 'diagnostic', label: 'Penunjang Diagnostik / Lab' },
    { value: 'medical_checkup', label: 'Medical Check-Up (MCU)' },
    { value: 'other', label: 'Lainnya (Other)' }
  ];

  $effect(() => {
    if (isOpen) {
      formError = null;
      isSaving = false;
      if (mode === 'edit' && department) {
        formCode = department.code;
        formName = department.name;
        formType = department.departement_type || 'outpatient';
        formPhone = department.phone || '';
        formEmail = department.email || '';
        formAddress = department.address || '';
        formIhsOrgId = department.ihs_organization_id || '';
        formIsActive = department.is_active;
      } else {
        formCode = '';
        formName = '';
        formType = 'outpatient';
        formPhone = '';
        formEmail = '';
        formAddress = '';
        formIhsOrgId = '';
        formIsActive = true;
      }
    }
  });

  async function handleSave() {
    formError = null;
    if (!formCode.trim()) {
      formError = 'Kode departemen wajib diisi';
      return;
    }
    if (!formName.trim()) {
      formError = 'Nama departemen wajib diisi';
      return;
    }

    isSaving = true;
    try {
      const payload: CreateDepartementDTO = {
        code: formCode.trim().toUpperCase(),
        name: formName.trim(),
        departement_type: formType,
        phone: formPhone.trim() || undefined,
        email: formEmail.trim() || undefined,
        address: formAddress.trim() || undefined,
        ihs_organization_id: formIhsOrgId.trim() || undefined,
        is_active: formIsActive
      };

      let savedDept: DepartementRecord;
      if (mode === 'create') {
        savedDept = await createDepartment(payload);
      } else if (department) {
        const updatePayload: UpdateDepartementDTO = { ...payload };
        savedDept = await updateDepartment(department.id, updatePayload);
      } else {
        return;
      }

      onSuccess(savedDept);
      onClose();
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan departemen';
    } finally {
      isSaving = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'create' ? 'Tambah Instalasi / Departemen Baru' : 'Edit Instalasi / Departemen'}
  subtitle="Master struktur organisasi instalasi layanan rumah sakit"
  icon={Building2}
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
        <label for="dept-form-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
          Kode Instalasi <span class="text-rose-500">*</span>
        </label>
        <input
          id="dept-form-code"
          type="text"
          bind:value={formCode}
          placeholder="Contoh: INST-RJ"
          disabled={mode === 'edit'}
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] disabled:opacity-60 disabled:cursor-not-allowed focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <M3Select
          id="dept-form-type"
          label="Kategori Layanan"
          required={true}
          bind:value={formType}
          options={deptTypeOptions}
          placeholder="Pilih Kategori Layanan"
        />
      </div>
    </div>

    <div>
      <label for="dept-form-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Nama Instalasi / Departemen <span class="text-rose-500">*</span>
      </label>
      <input
        id="dept-form-name"
        type="text"
        bind:value={formName}
        placeholder="Contoh: Instalasi Rawat Jalan (Poliklinik)"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label for="dept-form-phone" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
          Telepon / Ext
        </label>
        <input
          id="dept-form-phone"
          type="text"
          bind:value={formPhone}
          placeholder="Contoh: 022-710001 / Ext 101"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <label for="dept-form-email" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
          Email Internal
        </label>
        <input
          id="dept-form-email"
          type="email"
          bind:value={formEmail}
          placeholder="Contoh: rawatjalan@hosim.local"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>
    </div>

    <div>
      <label for="dept-form-ihs" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        SATUSEHAT Organization ID (IHS Org)
      </label>
      <input
        id="dept-form-ihs"
        type="text"
        bind:value={formIhsOrgId}
        placeholder="Contoh: ORG-10009823"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-1 block">
        Digunakan untuk integrasi interoperabilitas platform Kemenkes SATUSEHAT
      </span>
    </div>

    <div>
      <label for="dept-form-address" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Lokasi / Gedung
      </label>
      <input
        id="dept-form-address"
        type="text"
        bind:value={formAddress}
        placeholder="Contoh: Gedung A Lantai 1 & 2"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="pt-2">
      <M3ToggleCard
        bind:checked={formIsActive}
        title="Status Departemen Aktif"
        subtitle="Departemen beroperasi dan dapat menaungi unit layanan / poliklinik."
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
        <span>{mode === 'create' ? 'Simpan Departemen' : 'Perbarui Departemen'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
