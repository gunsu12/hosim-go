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
  import { getDepartments } from '$lib/api/organization';
  import type {
    DepartementRecord,
    DepartementType
  } from '../../../types/master/departement';

  import DepartmentFormModal from './components/DepartmentFormModal.svelte';
  import DepartmentDeleteModal from './components/DepartmentDeleteModal.svelte';

  let searchQuery = $state('');
  let departments = $state<DepartementRecord[]>([]);
  let totalCount = $state(0);
  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Modal State
  let showModal = $state(false);
  let modalMode = $state<'create' | 'edit'>('create');
  let selectedDepartment = $state<DepartementRecord | null>(null);

  // Delete Dialog State
  let showDeleteConfirm = $state(false);
  let departmentToDelete = $state<DepartementRecord | null>(null);

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
    selectedDepartment = null;
    showModal = true;
  }

  function openEditModal(d: DepartementRecord) {
    modalMode = 'edit';
    selectedDepartment = d;
    showModal = true;
  }

  function confirmDelete(d: DepartementRecord) {
    departmentToDelete = d;
    showDeleteConfirm = true;
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
        <h2 class="text-lg font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">Master Departemen / Instalasi</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa]">
          API: /api/v1/departments
        </span>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">Manajemen instalasi utama rumah sakit (Rawat Jalan, Rawat Inap, Kamar Operasi, IGD, Lab).</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="tonal" onclick={loadData} disabled={isLoading}>
        <RefreshCw class="w-3.5 h-3.5 mr-1 {isLoading ? 'animate-spin' : ''}" />
        <span>Refresh</span>
      </M3Button>

      {#if auth.hasPermission('department:create')}
        <M3Button variant="filled" onclick={openCreateModal}>
          <Plus class="w-4 h-4 mr-1" />
          <span>Tambah Departemen</span>
        </M3Button>
      {:else}
        <div class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] text-[#747775] dark:text-[#8e918f] text-xs font-medium border border-[#e1e5ea] dark:border-[#353738]">
          <ShieldAlert class="w-3.5 h-3.5 text-amber-600 dark:text-amber-400 shrink-0" />
          <span>Read-Only</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Notification Toasts -->
  {#if successMessage}
    <div class="flex items-center justify-between p-3.5 rounded-2xl bg-[#e6f4ea] dark:bg-emerald-950/40 border border-[#ceead6] dark:border-emerald-800 text-[#137333] dark:text-emerald-300 text-xs shadow-xs animate-fadeIn">
      <div class="flex items-center gap-2">
        <Check class="w-4 h-4 text-[#137333] dark:text-emerald-400 shrink-0" />
        <span>{successMessage}</span>
      </div>
      <button type="button" onclick={() => (successMessage = null)} class="text-[#137333] dark:text-emerald-400 hover:opacity-75 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div class="flex items-center justify-between p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 text-xs shadow-xs animate-fadeIn">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0" />
        <span>{errorMessage}</span>
      </div>
      <button type="button" onclick={() => (errorMessage = null)} class="text-rose-600 dark:text-rose-400 hover:opacity-75 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  <!-- Search and Stats Bar -->
  <div class="p-3 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] dark:text-[#8e918f] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        oninput={handleSearchInput}
        placeholder="Cari nama instalasi atau kode departemen..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
    <span class="text-xs text-[#444746] dark:text-[#c4c7c5]">
      Total: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3] font-mono">{totalCount}</strong> instalasi
    </span>
  </div>

  <!-- Table Container -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden relative min-h-[220px]">
    {#if isLoading && departments.length === 0}
      <div class="absolute inset-0 flex flex-col items-center justify-center bg-white/70 dark:bg-[#1e1f20]/70 backdrop-blur-xs z-10">
        <Loader2 class="w-7 h-7 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
        <span class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-2 font-medium">Memuat data departemen dari backend...</span>
      </div>
    {/if}

    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Instalasi / Departemen</th>
            <th class="py-3 px-4 font-semibold">Kategori Layanan</th>
            <th class="py-3 px-4 font-semibold">Telepon Internal</th>
            <th class="py-3 px-4 font-semibold">SATUSEHAT Org ID</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#if departments.length === 0 && !isLoading}
            <tr>
              <td colspan="7" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                <Building2 class="w-8 h-8 text-[#a8abb0] dark:text-[#5e5e5e] mx-auto mb-2 opacity-50" />
                <p class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">Belum ada data departemen</p>
                <p class="text-[11px] text-[#8e9196] dark:text-[#8e918f] mt-0.5">Silakan klik tombol "Tambah Departemen" untuk menambahkan baru.</p>
              </td>
            </tr>
          {:else}
            {#each departments as d (d.id)}
              <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/50 transition-colors">
                <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa]">{d.code}</td>
                <td class="py-3.5 px-4">
                  <div class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{d.name}</div>
                  {#if d.email}
                    <div class="text-[11px] text-[#747775] dark:text-[#8e918f]">{d.email}</div>
                  {/if}
                </td>
                <td class="py-3.5 px-4">
                  <span class="px-2 py-0.5 rounded-md font-medium text-[11px] bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] border border-[#e1e5ea] dark:border-[#353738]">
                    {getDeptTypeLabel(d.departement_type)}
                  </span>
                </td>
                <td class="py-3.5 px-4 font-mono text-[#444746] dark:text-[#c4c7c5]">{d.phone || '-'}</td>
                <td class="py-3.5 px-4 font-mono text-[#0b57d0] dark:text-[#a8c7fa]">{d.ihs_organization_id || '-'}</td>
                <td class="py-3.5 px-4 text-center">
                  {#if d.is_active}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300">
                      AKTIF
                    </span>
                  {:else}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300">
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
                        class="p-1.5 text-[#0b57d0] dark:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] rounded-lg transition-colors cursor-pointer"
                        title="Edit Departemen"
                      >
                        <Edit2 class="w-3.5 h-3.5" />
                      </button>
                    {/if}
                    {#if auth.hasPermission('department:delete')}
                      <button
                        type="button"
                        onclick={() => confirmDelete(d)}
                        class="p-1.5 text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/50 rounded-lg transition-colors cursor-pointer"
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

<!-- Modals -->
<DepartmentFormModal
  isOpen={showModal}
  mode={modalMode}
  department={selectedDepartment}
  onClose={() => (showModal = false)}
  onSuccess={() => {
    showSuccess(modalMode === 'create' ? 'Departemen berhasil ditambahkan' : 'Departemen berhasil diperbarui');
    loadData();
  }}
/>

<DepartmentDeleteModal
  isOpen={showDeleteConfirm}
  department={departmentToDelete}
  onClose={() => { showDeleteConfirm = false; departmentToDelete = null; }}
  onSuccess={() => {
    showSuccess('Departemen berhasil dihapus');
    loadData();
  }}
/>
