<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    Plus,
    Building2,
    ShieldAlert,
    Edit2,
    Trash2,
    X,
    Check,
    Loader2,
    AlertCircle,
    RefreshCw
  } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { auth } from '../../../stores/auth.svelte';
  import {
    getDepartments,
    createDepartment,
    updateDepartment,
    deleteDepartment
  } from '../../../api';
  import type {
    DepartementRecord,
    DepartementType,
    CreateDepartementDTO,
    UpdateDepartementDTO
  } from '../../../types/master/departement';

  let searchQuery = $state('');
  let departments = $state<DepartementRecord[]>([]);
  let totalCount = $state(0);
  let isLoading = $state(false);
  let isSaving = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Modal State
  let showModal = $state(false);
  let modalMode = $state<'create' | 'edit'>('create');
  let selectedId = $state<string | null>(null);

  // Delete Dialog State
  let showDeleteConfirm = $state(false);
  let departmentToDelete = $state<DepartementRecord | null>(null);
  let isDeleting = $state(false);

  // Form Fields
  let formCode = $state('');
  let formName = $state('');
  let formType = $state<DepartementType>('outpatient');
  let formPhone = $state('');
  let formEmail = $state('');
  let formAddress = $state('');
  let formIhsOrgId = $state('');
  let formIsActive = $state(true);
  let formError = $state<string | null>(null);

  const deptTypeOptions: { value: DepartementType; label: string }[] = [
    { value: 'outpatient', label: 'Rawat Jalan (Outpatient)' },
    { value: 'inpatient', label: 'Rawat Inap (Inpatient)' },
    { value: 'emergency', label: 'Gawat Darurat (Emergency / IGD)' },
    { value: 'diagnostic', label: 'Penunjang Diagnostik / Lab' },
    { value: 'medical_checkup', label: 'Medical Check-Up (MCU)' },
    { value: 'other', label: 'Lainnya (Other)' }
  ];

  function getDeptTypeLabel(type: string): string {
    const found = deptTypeOptions.find(o => o.value === type);
    return found ? found.label.split(' (')[0] : type;
  }

  async function loadData() {
    isLoading = true;
    errorMessage = null;
    try {
      const res = await getDepartments({ search: searchQuery.trim() });
      departments = res.data || [];
      totalCount = res.count ?? departments.length;
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memuat data departemen';
    } finally {
      isLoading = false;
    }
  }

  let searchTimeout: any;
  function handleSearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      loadData();
    }, 300);
  }

  function openCreateModal() {
    modalMode = 'create';
    selectedId = null;
    formCode = '';
    formName = '';
    formType = 'outpatient';
    formPhone = '';
    formEmail = '';
    formAddress = '';
    formIhsOrgId = '';
    formIsActive = true;
    formError = null;
    showModal = true;
  }

  function openEditModal(d: DepartementRecord) {
    modalMode = 'edit';
    selectedId = d.id;
    formCode = d.code;
    formName = d.name;
    formType = d.departement_type || 'outpatient';
    formPhone = d.phone || '';
    formEmail = d.email || '';
    formAddress = d.address || '';
    formIhsOrgId = d.ihs_organization_id || '';
    formIsActive = d.is_active;
    formError = null;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    formError = null;
  }

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

      if (modalMode === 'create') {
        await createDepartment(payload);
        showSuccess('Departemen berhasil ditambahkan');
      } else if (selectedId) {
        await updateDepartment(selectedId, payload as UpdateDepartementDTO);
        showSuccess('Departemen berhasil diperbarui');
      }

      closeModal();
      await loadData();
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan departemen';
    } finally {
      isSaving = false;
    }
  }

  function confirmDelete(d: DepartementRecord) {
    departmentToDelete = d;
    showDeleteConfirm = true;
  }

  async function handleDelete() {
    if (!departmentToDelete) return;
    isDeleting = true;
    try {
      await deleteDepartment(departmentToDelete.id);
      showDeleteConfirm = false;
      departmentToDelete = null;
      showSuccess('Departemen berhasil dihapus');
      await loadData();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal menghapus departemen';
    } finally {
      isDeleting = false;
    }
  }

  function showSuccess(msg: string) {
    successMessage = msg;
    setTimeout(() => {
      if (successMessage === msg) successMessage = null;
    }, 4000);
  }

  onMount(() => {
    loadData();
  });
</script>

<div class="flex flex-col gap-4">
  <!-- Header Banner -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Departemen / Instalasi</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          API: /api/v1/departments
        </span>
      </div>
      <p class="text-xs text-[#444746]">Manajemen instalasi utama rumah sakit (Rawat Jalan, Rawat Inap, Kamar Operasi, IGD, Lab).</p>
    </div>

    <div class="flex items-center gap-2">
      <button
        type="button"
        onclick={loadData}
        class="flex items-center gap-1 px-3 py-2 text-xs font-medium text-[#444746] bg-white border border-[#e1e5ea] rounded-xl hover:bg-[#f0f4f9] transition-all cursor-pointer"
        title="Muat ulang data"
      >
        <RefreshCw class="w-3.5 h-3.5 {isLoading ? 'animate-spin text-[#0b57d0]' : ''}" />
        <span>Refresh</span>
      </button>

      {#if auth.hasPermission('department:create')}
        <M3Button variant="filled" onclick={openCreateModal}>
          <Plus class="w-4 h-4" />
          <span>Tambah Departemen</span>
        </M3Button>
      {:else}
        <div class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-[#f0f4f9] text-[#747775] text-xs font-medium border border-[#e1e5ea]">
          <ShieldAlert class="w-3.5 h-3.5 text-amber-600 shrink-0" />
          <span>Read-Only</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Notification Toasts -->
  {#if successMessage}
    <div class="flex items-center justify-between p-3 rounded-xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs shadow-xs animate-fadeIn">
      <div class="flex items-center gap-2">
        <Check class="w-4 h-4 text-emerald-600 shrink-0" />
        <span>{successMessage}</span>
      </div>
      <button type="button" onclick={() => (successMessage = null)} class="text-emerald-600 hover:text-emerald-900 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div class="flex items-center justify-between p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-800 text-xs shadow-xs animate-fadeIn">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
        <span>{errorMessage}</span>
      </div>
      <button type="button" onclick={() => (errorMessage = null)} class="text-rose-600 hover:text-rose-900 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  <!-- Search and Stats Bar -->
  <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        oninput={handleSearchInput}
        placeholder="Cari nama instalasi atau kode departemen..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all"
      />
    </div>
    <span class="text-xs text-[#444746]">
      Total: <strong>{totalCount}</strong> instalasi
    </span>
  </div>

  <!-- Table Container -->
  <div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden relative min-h-[220px]">
    {#if isLoading && departments.length === 0}
      <div class="absolute inset-0 flex flex-col items-center justify-center bg-white/70 backdrop-blur-xs z-10">
        <Loader2 class="w-7 h-7 animate-spin text-[#0b57d0]" />
        <span class="text-xs text-[#444746] mt-2 font-medium">Memuat data departemen dari backend...</span>
      </div>
    {/if}

    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746]">
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Instalasi / Departemen</th>
            <th class="py-3 px-4 font-semibold">Kategori Layanan</th>
            <th class="py-3 px-4 font-semibold">Telepon Internal</th>
            <th class="py-3 px-4 font-semibold">SATUSEHAT Org ID</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#if departments.length === 0 && !isLoading}
            <tr>
              <td colspan="7" class="py-12 text-center text-[#747775]">
                <Building2 class="w-8 h-8 text-[#a8abb0] mx-auto mb-2 opacity-50" />
                <p class="font-medium">Belum ada data departemen</p>
                <p class="text-[11px] text-[#8e9196] mt-0.5">Silakan klik tombol "Tambah Departemen" untuk menambahkan baru.</p>
              </td>
            </tr>
          {:else}
            {#each departments as d}
              <tr class="hover:bg-[#f8fafd] transition-colors">
                <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0]">{d.code}</td>
                <td class="py-3.5 px-4">
                  <div class="font-semibold text-[#1f1f1f]">{d.name}</div>
                  {#if d.email}
                    <div class="text-[11px] text-[#747775]">{d.email}</div>
                  {/if}
                </td>
                <td class="py-3.5 px-4">
                  <span class="px-2 py-0.5 rounded-md font-medium text-[11px] bg-[#f0f4f9] text-[#444746] border border-[#e1e5ea]">
                    {getDeptTypeLabel(d.departement_type)}
                  </span>
                </td>
                <td class="py-3.5 px-4 font-mono text-[#444746]">{d.phone || '-'}</td>
                <td class="py-3.5 px-4 font-mono text-[#0b57d0]">{d.ihs_organization_id || '-'}</td>
                <td class="py-3.5 px-4 text-center">
                  {#if d.is_active}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                      AKTIF
                    </span>
                  {:else}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-gray-100 text-gray-600">
                      NON-AKTIF
                    </span>
                  {/if}
                </td>
                <td class="py-3.5 px-4 text-right">
                  <div class="flex items-center justify-end gap-1.5">
                    {#if auth.hasPermission('department:update')}
                      <button
                        type="button"
                        onclick={() => openEditModal(d)}
                        class="p-1.5 text-[#0b57d0] hover:bg-[#e8f0fe] rounded-lg transition-colors cursor-pointer"
                        title="Edit Departemen"
                      >
                        <Edit2 class="w-3.5 h-3.5" />
                      </button>
                    {/if}
                    {#if auth.hasPermission('department:delete')}
                      <button
                        type="button"
                        onclick={() => confirmDelete(d)}
                        class="p-1.5 text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
                        title="Hapus Departemen"
                      >
                        <Trash2 class="w-3.5 h-3.5" />
                      </button>
                    {/if}
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Modal Form Tambah / Edit Departemen -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-fadeIn">
    <div class="w-full max-w-lg bg-white rounded-3xl shadow-xl border border-[#e1e5ea] overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-[#e1e5ea] bg-[#f8fafd]">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-xl bg-[#e8f0fe] text-[#0b57d0]">
            <Building2 class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-semibold text-[#1f1f1f]">
              {modalMode === 'create' ? 'Tambah Departemen / Instalasi Baru' : 'Edit Departemen'}
            </h3>
            <p class="text-[11px] text-[#747775]">Masukkan detail informasi instalasi rumah sakit</p>
          </div>
        </div>
        <button
          type="button"
          onclick={closeModal}
          class="p-1 text-[#747775] hover:text-[#1f1f1f] hover:bg-[#e1e5ea] rounded-full transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto flex flex-col gap-4 text-xs">
        {#if formError}
          <div class="flex items-center gap-2 p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-800">
            <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
            <span>{formError}</span>
          </div>
        {/if}

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="dept-form-code" class="block text-[11px] font-semibold text-[#444746] mb-1">Kode Instalasi *</label>
            <input
              id="dept-form-code"
              type="text"
              bind:value={formCode}
              placeholder="Contoh: INST-RJ"
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs font-mono uppercase focus:outline-none transition-all"
            />
          </div>

          <div>
            <label for="dept-form-type" class="block text-[11px] font-semibold text-[#444746] mb-1">Kategori Layanan *</label>
            <select
              id="dept-form-type"
              bind:value={formType}
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all cursor-pointer"
            >
              {#each deptTypeOptions as opt}
                <option value={opt.value}>{opt.label}</option>
              {/each}
            </select>
          </div>
        </div>

        <div>
          <label for="dept-form-name" class="block text-[11px] font-semibold text-[#444746] mb-1">Nama Instalasi / Departemen *</label>
          <input
            id="dept-form-name"
            type="text"
            bind:value={formName}
            placeholder="Contoh: Instalasi Rawat Jalan (Poliklinik)"
            class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="dept-form-phone" class="block text-[11px] font-semibold text-[#444746] mb-1">Telepon / Ext</label>
            <input
              id="dept-form-phone"
              type="text"
              bind:value={formPhone}
              placeholder="Contoh: 022-710001 / Ext 101"
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
            />
          </div>

          <div>
            <label for="dept-form-email" class="block text-[11px] font-semibold text-[#444746] mb-1">Email Internal</label>
            <input
              id="dept-form-email"
              type="email"
              bind:value={formEmail}
              placeholder="Contoh: rawatjalan@hosim.local"
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
            />
          </div>
        </div>

        <div>
          <label for="dept-form-ihs" class="block text-[11px] font-semibold text-[#444746] mb-1">SATUSEHAT Organization ID (IHS Org)</label>
          <input
            id="dept-form-ihs"
            type="text"
            bind:value={formIhsOrgId}
            placeholder="Contoh: ORG-10009823"
            class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs font-mono focus:outline-none transition-all"
          />
          <span class="text-[10px] text-[#747775]">Digunakan untuk integrasi interoperabilitas platform Kemenkes SATUSEHAT</span>
        </div>

        <div>
          <label for="dept-form-address" class="block text-[11px] font-semibold text-[#444746] mb-1">Lokasi / Gedung</label>
          <input
            id="dept-form-address"
            type="text"
            bind:value={formAddress}
            placeholder="Contoh: Gedung A Lantai 1 & 2"
            class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
          />
        </div>

        <div class="flex items-center gap-2 pt-1">
          <input
            type="checkbox"
            id="formIsActiveDept"
            bind:checked={formIsActive}
            class="w-4 h-4 rounded text-[#0b57d0] focus:ring-0 cursor-pointer"
          />
          <label for="formIsActiveDept" class="text-xs text-[#1f1f1f] select-none cursor-pointer">
            Departemen Aktif dan Beroperasi
          </label>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="flex items-center justify-end gap-2 px-6 py-4 border-t border-[#e1e5ea] bg-[#f8fafd]">
        <button
          type="button"
          onclick={closeModal}
          class="px-4 py-2 text-xs font-medium text-[#444746] hover:bg-[#e1e5ea] rounded-xl transition-colors cursor-pointer"
        >
          Batal
        </button>
        <button
          type="button"
          onclick={handleSave}
          disabled={isSaving}
          class="flex items-center gap-1.5 px-4 py-2 text-xs font-medium text-white bg-[#0b57d0] hover:bg-[#0842a0] disabled:opacity-50 rounded-xl shadow-xs transition-colors cursor-pointer"
        >
          {#if isSaving}
            <Loader2 class="w-4 h-4 animate-spin" />
            <span>Menyimpan...</span>
          {:else}
            <Check class="w-4 h-4" />
            <span>{modalMode === 'create' ? 'Simpan Departemen' : 'Perbarui'}</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Dialog Konfirmasi Hapus -->
{#if showDeleteConfirm && departmentToDelete}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-fadeIn">
    <div class="w-full max-w-sm bg-white rounded-3xl shadow-xl border border-[#e1e5ea] p-6 flex flex-col gap-4 text-xs">
      <div class="flex items-center gap-3 text-rose-600">
        <div class="p-2.5 rounded-2xl bg-rose-50">
          <Trash2 class="w-5 h-5" />
        </div>
        <div>
          <h3 class="text-sm font-semibold text-[#1f1f1f]">Hapus Departemen?</h3>
          <p class="text-[11px] text-[#747775]">Tindakan ini akan menonaktifkan departemen</p>
        </div>
      </div>

      <p class="text-[#444746] leading-relaxed">
        Apakah Anda yakin ingin menghapus instalasi <strong>{departmentToDelete.name}</strong> ({departmentToDelete.code})?
      </p>

      <div class="flex items-center justify-end gap-2 pt-2">
        <button
          type="button"
          onclick={() => { showDeleteConfirm = false; departmentToDelete = null; }}
          class="px-3.5 py-2 text-xs font-medium text-[#444746] hover:bg-[#f0f4f9] rounded-xl transition-colors cursor-pointer"
        >
          Batal
        </button>
        <button
          type="button"
          onclick={handleDelete}
          disabled={isDeleting}
          class="flex items-center gap-1.5 px-4 py-2 text-xs font-medium text-white bg-rose-600 hover:bg-rose-700 disabled:opacity-50 rounded-xl shadow-xs transition-colors cursor-pointer"
        >
          {#if isDeleting}
            <Loader2 class="w-3.5 h-3.5 animate-spin" />
            <span>Menghapus...</span>
          {:else}
            <span>Ya, Hapus</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
