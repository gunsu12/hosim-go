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
  import { auth } from '../../../stores/auth.svelte';
  import {
    getServiceUnits,
    createServiceUnit,
    updateServiceUnit,
    deleteServiceUnit,
    getDepartments
  } from '../../../api';
  import type {
    ServiceUnitRecord,
    CreateServiceUnitDTO,
    UpdateServiceUnitDTO
  } from '../../../types/master/service_unit';
  import type { DepartementRecord } from '../../../types/master/departement';

  let searchQuery = $state('');
  let serviceUnits = $state<ServiceUnitRecord[]>([]);
  let departments = $state<DepartementRecord[]>([]);
  let totalCount = $state(0);
  let isLoading = $state(false);
  let isSaving = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Filter Departemen dropdown di atas tabel
  let selectedDeptFilter = $state('');

  // Modal State
  let showModal = $state(false);
  let modalMode = $state<'create' | 'edit'>('create');
  let selectedId = $state<string | null>(null);

  // Delete Dialog State
  let showDeleteConfirm = $state(false);
  let unitToDelete = $state<ServiceUnitRecord | null>(null);
  let isDeleting = $state(false);

  // Form Fields
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
    selectedId = null;
    formCode = '';
    formName = '';
    formDeptId = departments.length > 0 ? departments[0].id : '';
    formPhone = '';
    formEmail = '';
    formAddress = '';
    formIhsLocId = '';
    formIsRegistrationTarget = true;
    formIsActive = true;
    formError = null;
    showModal = true;
  }

  function openEditModal(u: ServiceUnitRecord) {
    modalMode = 'edit';
    selectedId = u.id;
    formCode = u.code;
    formName = u.name;
    formDeptId = u.departement_id || '';
    formPhone = u.phone || '';
    formEmail = u.email || '';
    formAddress = u.address || '';
    formIhsLocId = u.ihs_location_id || '';
    formIsRegistrationTarget = u.is_registration_target;
    formIsActive = u.is_active;
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
        departement_id: formDeptId ? formDeptId : undefined,
        phone: formPhone.trim() || undefined,
        email: formEmail.trim() || undefined,
        address: formAddress.trim() || undefined,
        ihs_location_id: formIhsLocId.trim() || undefined,
        is_registration_target: formIsRegistrationTarget,
        is_active: formIsActive
      };

      if (modalMode === 'create') {
        await createServiceUnit(payload);
        showSuccess('Unit layanan berhasil ditambahkan');
      } else if (selectedId) {
        await updateServiceUnit(selectedId, payload as UpdateServiceUnitDTO);
        showSuccess('Unit layanan berhasil diperbarui');
      }

      closeModal();
      await loadServiceUnits();
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan unit layanan';
    } finally {
      isSaving = false;
    }
  }

  function confirmDelete(u: ServiceUnitRecord) {
    unitToDelete = u;
    showDeleteConfirm = true;
  }

  async function handleDelete() {
    if (!unitToDelete) return;
    isDeleting = true;
    try {
      await deleteServiceUnit(unitToDelete.id);
      showDeleteConfirm = false;
      unitToDelete = null;
      showSuccess('Unit layanan berhasil dihapus');
      await loadServiceUnits();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal menghapus unit layanan';
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
    loadInitial();
  });
</script>

<div class="flex flex-col gap-4">
  <!-- Header Banner -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Unit Layanan (Service Unit)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          API: /api/v1/service-units
        </span>
      </div>
      <p class="text-xs text-[#444746]">Daftar poliklinik, instalasi penunjang, dan unit tujuan pendaftaran pasien.</p>
    </div>

    <div class="flex items-center gap-2">
      <button
        type="button"
        onclick={loadServiceUnits}
        class="flex items-center gap-1 px-3 py-2 text-xs font-medium text-[#444746] bg-white border border-[#e1e5ea] rounded-xl hover:bg-[#f0f4f9] transition-all cursor-pointer"
        title="Muat ulang data"
      >
        <RefreshCw class="w-3.5 h-3.5 {isLoading ? 'animate-spin text-[#0b57d0]' : ''}" />
        <span>Refresh</span>
      </button>

      {#if auth.hasPermission('service_unit:create')}
        <M3Button variant="filled" onclick={openCreateModal}>
          <Plus class="w-4 h-4" />
          <span>Tambah Unit Layanan</span>
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

  <!-- Filter & Search Bar -->
  <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex flex-wrap items-center justify-between gap-3 shadow-xs">
    <div class="flex items-center gap-2 flex-1 max-w-xl">
      <div class="relative flex-1">
        <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
        <input
          type="text"
          bind:value={searchQuery}
          oninput={handleSearchInput}
          placeholder="Cari nama poliklinik atau kode unit..."
          class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all"
        />
      </div>

      <!-- Departemen Filter -->
      <select
        bind:value={selectedDeptFilter}
        onchange={handleFilterDeptChange}
        class="h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all cursor-pointer max-w-[200px]"
      >
        <option value="">Semua Instalasi</option>
        {#each departments as dept}
          <option value={dept.id}>{dept.name}</option>
        {/each}
      </select>
    </div>

    <span class="text-xs text-[#444746]">
      Total: <strong>{totalCount}</strong> unit
    </span>
  </div>

  <!-- Table Container -->
  <div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden relative min-h-[220px]">
    {#if isLoading && serviceUnits.length === 0}
      <div class="absolute inset-0 flex flex-col items-center justify-center bg-white/70 backdrop-blur-xs z-10">
        <Loader2 class="w-7 h-7 animate-spin text-[#0b57d0]" />
        <span class="text-xs text-[#444746] mt-2 font-medium">Memuat data unit layanan dari backend...</span>
      </div>
    {/if}

    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746]">
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
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#if serviceUnits.length === 0 && !isLoading}
            <tr>
              <td colspan="8" class="py-12 text-center text-[#747775]">
                <Stethoscope class="w-8 h-8 text-[#a8abb0] mx-auto mb-2 opacity-50" />
                <p class="font-medium">Belum ada unit layanan</p>
                <p class="text-[11px] text-[#8e9196] mt-0.5">Silakan klik "Tambah Unit Layanan" untuk mulai mendaftarkan poli / unit.</p>
              </td>
            </tr>
          {:else}
            {#each serviceUnits as u}
              <tr class="hover:bg-[#f8fafd] transition-colors">
                <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0]">{u.code}</td>
                <td class="py-3.5 px-4 font-semibold text-[#1f1f1f]">
                  <div>{u.name}</div>
                  {#if u.address}
                    <div class="text-[11px] font-normal text-[#747775]">{u.address}</div>
                  {/if}
                </td>
                <td class="py-3.5 px-4">
                  {#if u.departement}
                    <span class="px-2 py-0.5 rounded-md font-medium text-[11px] bg-[#e8f0fe] text-[#0b57d0]">
                      {u.departement.name}
                    </span>
                  {:else}
                    <span class="text-[#8e9196]">-</span>
                  {/if}
                </td>
                <td class="py-3.5 px-4 font-mono text-[#444746]">{u.phone || '-'}</td>
                <td class="py-3.5 px-4 font-mono text-[#0b57d0]">{u.ihs_location_id || '-'}</td>
                <td class="py-3.5 px-4 text-center">
                  {#if u.is_registration_target}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-100 text-blue-800">
                      Bisa Didaftar
                    </span>
                  {:else}
                    <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600">
                      Internal
                    </span>
                  {/if}
                </td>
                <td class="py-3.5 px-4 text-center">
                  {#if u.is_active}
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
                    {#if auth.hasPermission('service_unit:update')}
                      <button
                        type="button"
                        onclick={() => openEditModal(u)}
                        class="p-1.5 text-[#0b57d0] hover:bg-[#e8f0fe] rounded-lg transition-colors cursor-pointer"
                        title="Edit Unit Layanan"
                      >
                        <Edit2 class="w-3.5 h-3.5" />
                      </button>
                    {/if}
                    {#if auth.hasPermission('service_unit:delete')}
                      <button
                        type="button"
                        onclick={() => confirmDelete(u)}
                        class="p-1.5 text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
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

<!-- Modal Form Tambah / Edit Service Unit -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-fadeIn">
    <div class="w-full max-w-lg bg-white rounded-3xl shadow-xl border border-[#e1e5ea] overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-[#e1e5ea] bg-[#f8fafd]">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-xl bg-[#e8f0fe] text-[#0b57d0]">
            <Stethoscope class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-semibold text-[#1f1f1f]">
              {modalMode === 'create' ? 'Tambah Unit Layanan / Poli Baru' : 'Edit Unit Layanan'}
            </h3>
            <p class="text-[11px] text-[#747775]">Masukkan detail informasi unit layanan atau poliklinik</p>
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
            <label for="su-form-code" class="block text-[11px] font-semibold text-[#444746] mb-1">Kode Unit *</label>
            <input
              id="su-form-code"
              type="text"
              bind:value={formCode}
              placeholder="Contoh: POLI-BEDAH"
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs font-mono uppercase focus:outline-none transition-all"
            />
          </div>

          <div>
            <label for="su-form-dept" class="block text-[11px] font-semibold text-[#444746] mb-1">Induk Departemen</label>
            <select
              id="su-form-dept"
              bind:value={formDeptId}
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all cursor-pointer"
            >
              <option value="">-- Tanpa Induk Departemen --</option>
              {#each departments as dept}
                <option value={dept.id}>{dept.name} ({dept.code})</option>
              {/each}
            </select>
          </div>
        </div>

        <div>
          <label for="su-form-name" class="block text-[11px] font-semibold text-[#444746] mb-1">Nama Unit Layanan / Poli *</label>
          <input
            id="su-form-name"
            type="text"
            bind:value={formName}
            placeholder="Contoh: Poliklinik Bedah Umum"
            class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
          />
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label for="su-form-phone" class="block text-[11px] font-semibold text-[#444746] mb-1">Telepon / Ext</label>
            <input
              id="su-form-phone"
              type="text"
              bind:value={formPhone}
              placeholder="Contoh: Ext. 101"
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
            />
          </div>

          <div>
            <label for="su-form-email" class="block text-[11px] font-semibold text-[#444746] mb-1">Email Unit</label>
            <input
              id="su-form-email"
              type="email"
              bind:value={formEmail}
              placeholder="Contoh: bedah@hosim.local"
              class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
            />
          </div>
        </div>

        <div>
          <label for="su-form-ihs" class="block text-[11px] font-semibold text-[#444746] mb-1">SATUSEHAT Location ID</label>
          <input
            id="su-form-ihs"
            type="text"
            bind:value={formIhsLocId}
            placeholder="Contoh: LOC-3273-0101"
            class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs font-mono focus:outline-none transition-all"
          />
          <span class="text-[10px] text-[#747775]">ID Lokasi faskes terdaftar di platform SatuSehat Kemenkes</span>
        </div>

        <div>
          <label for="su-form-address" class="block text-[11px] font-semibold text-[#444746] mb-1">Lokasi / Ruang / Lantai</label>
          <input
            id="su-form-address"
            type="text"
            bind:value={formAddress}
            placeholder="Contoh: Lantai 2 Sayap Barat Kamar 204"
            class="w-full h-10 px-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs focus:outline-none transition-all"
          />
        </div>

        <div class="flex flex-col gap-2 pt-1 border-t border-[#e1e5ea]">
          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              id="formIsRegTarget"
              bind:checked={formIsRegistrationTarget}
              class="w-4 h-4 rounded text-[#0b57d0] focus:ring-0 cursor-pointer"
            />
            <label for="formIsRegTarget" class="text-xs text-[#1f1f1f] select-none cursor-pointer">
              Tujuan Pendaftaran Pasien (Bisa dipilih saat pasien mendaftar di loket / online)
            </label>
          </div>

          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              id="formIsActiveUnit"
              bind:checked={formIsActive}
              class="w-4 h-4 rounded text-[#0b57d0] focus:ring-0 cursor-pointer"
            />
            <label for="formIsActiveUnit" class="text-xs text-[#1f1f1f] select-none cursor-pointer">
              Unit Layanan Aktif dan Beroperasi
            </label>
          </div>
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
            <span>{modalMode === 'create' ? 'Simpan Unit Layanan' : 'Perbarui'}</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}

<!-- Dialog Konfirmasi Hapus -->
{#if showDeleteConfirm && unitToDelete}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-fadeIn">
    <div class="w-full max-w-sm bg-white rounded-3xl shadow-xl border border-[#e1e5ea] p-6 flex flex-col gap-4 text-xs">
      <div class="flex items-center gap-3 text-rose-600">
        <div class="p-2.5 rounded-2xl bg-rose-50">
          <Trash2 class="w-5 h-5" />
        </div>
        <div>
          <h3 class="text-sm font-semibold text-[#1f1f1f]">Hapus Unit Layanan?</h3>
          <p class="text-[11px] text-[#747775]">Tindakan ini akan menghapus data unit layanan</p>
        </div>
      </div>

      <p class="text-[#444746] leading-relaxed">
        Apakah Anda yakin ingin menghapus unit <strong>{unitToDelete.name}</strong> ({unitToDelete.code})?
      </p>

      <div class="flex items-center justify-end gap-2 pt-2">
        <button
          type="button"
          onclick={() => { showDeleteConfirm = false; unitToDelete = null; }}
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
