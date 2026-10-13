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
    RefreshCw,
    Layers,
    Stethoscope
  } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import M3Select from '../../../components/m3/M3Select.svelte';
  import { auth } from '../../../stores/auth.svelte';
  import {
    getServiceUnits,
    getDepartments
  } from '$lib/api/organization';
  import type { ServiceUnitRecord } from '../../../types/master/service_unit';
  import type { DepartementRecord } from '../../../types/master/departement';

  import ServiceUnitFormModal from './components/ServiceUnitFormModal.svelte';
  import ServiceUnitDeleteModal from './components/ServiceUnitDeleteModal.svelte';

  let searchQuery = $state('');
  let serviceUnits = $state<ServiceUnitRecord[]>([]);
  let departments = $state<DepartementRecord[]>([]);
  let totalCount = $state(0);
  let isLoading = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Filter Departemen dropdown di atas tabel
  let selectedDeptFilter = $state('');

  // Modal State
  let showModal = $state(false);
  let modalMode = $state<'create' | 'edit'>('create');
  let selectedUnit = $state<ServiceUnitRecord | null>(null);

  // Delete Dialog State
  let showDeleteConfirm = $state(false);
  let unitToDelete = $state<ServiceUnitRecord | null>(null);

  async function loadInitial() {
    await Promise.all([loadServiceUnits(), loadDepartmentsList()]);
  }

  async function loadDepartmentsList() {
    try {
      const res = await getDepartments({ limit: 100 });
      departments = res.data || [];
    } catch (e) {
      console.error('Gagal memuat daftar departemen', e);
    }
  }

  async function loadServiceUnits() {
    isLoading = true;
    errorMessage = null;
    try {
      const res = await getServiceUnits({
        search: searchQuery.trim(),
        departement_id: selectedDeptFilter || undefined
      });
      serviceUnits = res.data || [];
      totalCount = res.count ?? serviceUnits.length;
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memuat data unit layanan';
    } finally {
      isLoading = false;
    }
  }

  let searchTimeout: any;
  function handleSearchInput() {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      loadServiceUnits();
    }, 300);
  }

  function handleFilterDeptChange() {
    loadServiceUnits();
  }

  function openCreateModal() {
    modalMode = 'create';
    selectedUnit = null;
    showModal = true;
  }

  function openEditModal(u: ServiceUnitRecord) {
    modalMode = 'edit';
    selectedUnit = u;
    showModal = true;
  }

  function confirmDelete(u: ServiceUnitRecord) {
    unitToDelete = u;
    showDeleteConfirm = true;
  }

  function showSuccess(msg: string) {
    successMessage = msg;
    setTimeout(() => {
      if (successMessage === msg) successMessage = null;
    }, 4000);
  }

  onMount(() => {
    loadInitial();
  });
</script>

<div class="flex flex-col gap-4">
  <!-- Header Banner -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">Master Unit Layanan (Service Unit)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa]">
          API: /api/v1/service-units
        </span>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">Daftar poliklinik, instalasi penunjang, dan unit tujuan pendaftaran pasien.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="tonal" onclick={loadServiceUnits} disabled={isLoading}>
        <RefreshCw class="w-3.5 h-3.5 mr-1 {isLoading ? 'animate-spin' : ''}" />
        <span>Refresh</span>
      </M3Button>

      {#if auth.hasPermission('service_unit:create')}
        <M3Button variant="filled" onclick={openCreateModal}>
          <Plus class="w-4 h-4 mr-1" />
          <span>Tambah Unit Layanan</span>
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

  <!-- Filter & Search Bar -->
  <div class="p-3 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] flex flex-wrap items-center justify-between gap-3 shadow-xs">
    <div class="flex flex-wrap items-center gap-2 flex-1 max-w-xl">
      <div class="relative flex-1 min-w-[200px]">
        <Search class="w-4 h-4 text-[#747775] dark:text-[#8e918f] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
        <input
          type="text"
          bind:value={searchQuery}
          oninput={handleSearchInput}
          placeholder="Cari nama poliklinik atau kode unit..."
          class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <!-- Departemen Filter -->
      <div class="w-52">
        <M3Select
          bind:value={selectedDeptFilter}
          onchange={handleFilterDeptChange}
          placeholder="Semua Instalasi"
          options={[
            { value: '', label: 'Semua Instalasi' },
            ...departments.map((d) => ({
              value: d.id,
              label: d.name,
              badge: d.code
            }))
          ]}
        />
      </div>
    </div>

    <span class="text-xs text-[#444746] dark:text-[#c4c7c5]">
      Total: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3] font-mono">{totalCount}</strong> unit
    </span>
  </div>

  <!-- Table Container -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden relative min-h-[220px]">
    {#if isLoading && serviceUnits.length === 0}
      <div class="absolute inset-0 flex flex-col items-center justify-center bg-white/70 dark:bg-[#1e1f20]/70 backdrop-blur-xs z-10">
        <Loader2 class="w-7 h-7 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
        <span class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-2 font-medium">Memuat data unit layanan dari backend...</span>
      </div>
    {/if}

    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Unit Layanan / Poli</th>
            <th class="py-3 px-4 font-semibold">Induk Departemen</th>
            <th class="py-3 px-4 font-semibold">Telepon / Ext</th>
            <th class="py-3 px-4 font-semibold">SATUSEHAT Location ID</th>
            <th class="py-3 px-4 font-semibold text-center">Tujuan Pendaftaran</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#if serviceUnits.length === 0 && !isLoading}
            <tr>
              <td colspan="8" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                <Stethoscope class="w-8 h-8 text-[#a8abb0] dark:text-[#5e5e5e] mx-auto mb-2 opacity-50" />
                <p class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">Belum ada unit layanan</p>
                <p class="text-[11px] text-[#8e9196] dark:text-[#8e918f] mt-0.5">Silakan klik "Tambah Unit Layanan" untuk mulai mendaftarkan poli / unit.</p>
              </td>
            </tr>
          {:else}
            {#each serviceUnits as u (u.id)}
              <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/50 transition-colors">
                <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa]">{u.code}</td>
                <td class="py-3.5 px-4 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
                  <div>{u.name}</div>
                  {#if u.address}
                    <div class="text-[11px] font-normal text-[#747775] dark:text-[#8e918f]">{u.address}</div>
                  {/if}
                </td>
                <td class="py-3.5 px-4">
                  {#if u.departement}
                    <span class="px-2 py-0.5 rounded-md font-medium text-[11px] bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa] border border-[#d2e3fc] dark:border-[#2b3a55]">
                      {u.departement.name}
                    </span>
                  {:else}
                    <span class="text-[#8e9196] dark:text-[#747775]">-</span>
                  {/if}
                </td>
                <td class="py-3.5 px-4 font-mono text-[#444746] dark:text-[#c4c7c5]">{u.phone || '-'}</td>
                <td class="py-3.5 px-4 font-mono text-[#0b57d0] dark:text-[#a8c7fa]">{u.ihs_location_id || '-'}</td>
                <td class="py-3.5 px-4 text-center">
                  {#if u.is_registration_target}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-100 dark:bg-blue-950/60 text-blue-800 dark:text-blue-300">
                      Bisa Didaftar
                    </span>
                  {:else}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300">
                      Internal
                    </span>
                  {/if}
                </td>
                <td class="py-3.5 px-4 text-center">
                  {#if u.is_active}
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
                    {#if auth.hasPermission('service_unit:update')}
                      <button
                        type="button"
                        onclick={() => openEditModal(u)}
                        class="p-1.5 text-[#0b57d0] dark:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] rounded-lg transition-colors cursor-pointer"
                        title="Edit Unit Layanan"
                      >
                        <Edit2 class="w-3.5 h-3.5" />
                      </button>
                    {/if}
                    {#if auth.hasPermission('service_unit:delete')}
                      <button
                        type="button"
                        onclick={() => confirmDelete(u)}
                        class="p-1.5 text-rose-600 dark:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/50 rounded-lg transition-colors cursor-pointer"
                        title="Hapus Unit Layanan"
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
<ServiceUnitFormModal
  isOpen={showModal}
  mode={modalMode}
  unit={selectedUnit}
  {departments}
  onClose={() => (showModal = false)}
  onSuccess={() => {
    showSuccess(modalMode === 'create' ? 'Unit layanan berhasil ditambahkan' : 'Unit layanan berhasil diperbarui');
    loadServiceUnits();
  }}
/>

<ServiceUnitDeleteModal
  isOpen={showDeleteConfirm}
  unit={unitToDelete}
  onClose={() => { showDeleteConfirm = false; unitToDelete = null; }}
  onSuccess={() => {
    showSuccess('Unit layanan berhasil dihapus');
    loadServiceUnits();
  }}
/>
