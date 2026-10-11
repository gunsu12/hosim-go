<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    Plus,
    RefreshCw,
    FolderTree,
    ListFilter,
    ShieldAlert,
    Check,
    X,
    AlertCircle
  } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import M3SegmentedButton from '../../../components/m3/M3SegmentedButton.svelte';
  import { auth } from '../../../stores/auth.svelte';
  import {
    getAccountTree,
    getAccounts,
    createAccount,
    updateAccount,
    deleteAccount
  } from '$lib/api/accounting';
  import type {
    AccountRecord,
    AccountTreeNode,
    AccountType,
    CreateAccountDTO,
    UpdateAccountDTO
  } from '../../../types/accounting/account';

  import CoaStatsCards from './components/CoaStatsCards.svelte';
  import CoaTreeView from './components/CoaTreeView.svelte';
  import CoaListView from './components/CoaListView.svelte';
  import CoaFormModal from './components/CoaFormModal.svelte';
  import CoaDeleteModal from './components/CoaDeleteModal.svelte';
  import {
    filterTreeRecursive,
    filterFlatAccounts,
    collectExpandedNodeIds
  } from './coaUtils';

  // State Management
  let viewMode = $state<string>('tree'); // 'tree' | 'list'
  let searchQuery = $state('');
  let selectedTypeFilter = $state<AccountType | 'ALL'>('ALL');
  let activeOnlyFilter = $state(false);
  let postableFilter = $state<'ALL' | 'POSTABLE' | 'HEADER'>('ALL');

  let treeData = $state<AccountTreeNode[]>([]);
  let flatAccounts = $state<AccountRecord[]>([]);
  let isLoading = $state(false);
  let isSaving = $state(false);
  let isDeleting = $state(false);
  let errorMessage = $state<string | null>(null);
  let successMessage = $state<string | null>(null);

  // Expanded Tree Node IDs
  let expandedNodeIds = $state<Set<string>>(new Set());

  // Modal State
  let showModal = $state(false);
  let modalMode = $state<'create' | 'edit'>('create');
  let selectedAccount = $state<AccountRecord | null>(null);
  let parentForNewChild = $state<AccountRecord | null>(null);

  // Delete Confirmation Modal State
  let showDeleteConfirm = $state(false);
  let accountToDelete = $state<AccountRecord | null>(null);
  let deleteErrorMessage = $state<string | null>(null);

  // Derived
  const stats = $derived.by(() => {
    let total = flatAccounts.length;
    let headers = 0;
    let postable = 0;
    let treasury = 0;
    let byType: Record<AccountType, number> = {
      ASSET: 0,
      LIABILITY: 0,
      EQUITY: 0,
      REVENUE: 0,
      EXPENSE: 0
    };

    for (const acc of flatAccounts) {
      if (acc.is_postable) postable++;
      else headers++;
      if (acc.is_treasury_account) treasury++;
      if (acc.type in byType) byType[acc.type]++;
    }

    return { total, headers, postable, treasury, byType };
  });

  const filteredTree = $derived(
    filterTreeRecursive(treeData, searchQuery, selectedTypeFilter, activeOnlyFilter, postableFilter)
  );

  const filteredFlatAccounts = $derived(
    filterFlatAccounts(flatAccounts, searchQuery, selectedTypeFilter, activeOnlyFilter, postableFilter)
  );

  // Data Fetching
  async function loadData() {
    isLoading = true;
    errorMessage = null;
    try {
      const [treeRes, flatRes] = await Promise.all([
        getAccountTree(false),
        getAccounts({ limit: 1000 })
      ]);
      treeData = treeRes || [];
      flatAccounts = flatRes.items || [];

      if (expandedNodeIds.size === 0) {
        expandedNodeIds = collectExpandedNodeIds(treeData, 2);
      }
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memuat Bagan Akun (COA)';
    } finally {
      isLoading = false;
    }
  }

  function toggleNode(nodeId: string) {
    const nextSet = new Set(expandedNodeIds);
    if (nextSet.has(nodeId)) nextSet.delete(nodeId);
    else nextSet.add(nodeId);
    expandedNodeIds = nextSet;
  }

  // Modal Handlers
  function openCreateModal(parent?: AccountRecord) {
    modalMode = 'create';
    selectedAccount = null;
    parentForNewChild = parent || null;
    showModal = true;
  }

  function openEditModal(account: AccountRecord) {
    modalMode = 'edit';
    selectedAccount = account;
    parentForNewChild = null;
    showModal = true;
  }

  async function handleSaveAccount(payload: CreateAccountDTO | UpdateAccountDTO) {
    isSaving = true;
    try {
      if (modalMode === 'create') {
        const created = await createAccount(payload as CreateAccountDTO);
        successMessage = `Akun ${created.code} - ${created.name} berhasil dibuat!`;
      } else if (modalMode === 'edit' && selectedAccount) {
        const updated = await updateAccount(selectedAccount.id, payload as UpdateAccountDTO);
        successMessage = `Akun ${updated.code} - ${updated.name} berhasil diperbarui!`;
      }
      showModal = false;
      await loadData();
    } catch (err: any) {
      throw err;
    } finally {
      isSaving = false;
    }
  }

  function confirmDelete(account: AccountRecord) {
    accountToDelete = account;
    deleteErrorMessage = null;
    showDeleteConfirm = true;
  }

  async function handleDeleteAccount() {
    if (!accountToDelete) return;
    isDeleting = true;
    deleteErrorMessage = null;
    try {
      await deleteAccount(accountToDelete.id);
      successMessage = `Akun ${accountToDelete.code} - ${accountToDelete.name} berhasil dihapus.`;
      showDeleteConfirm = false;
      accountToDelete = null;
      await loadData();
    } catch (err: any) {
      deleteErrorMessage = err.message || 'Gagal menghapus akun. Pastikan akun tidak memiliki sub-akun anak atau keterikatan tarif.';
    } finally {
      isDeleting = false;
    }
  }

  onMount(() => {
    loadData();
  });
</script>

<div class="flex flex-col gap-5">
  <!-- Header Banner & Action Bar -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div class="flex items-center gap-2.5">
      <div class="p-2 rounded-xl bg-primary/10 text-primary">
        <FolderTree class="w-5 h-5" />
      </div>
      <div>
        <div class="flex items-center gap-2">
          <h2 class="text-xl font-bold text-on-surface tracking-tight">Bagan Akun Standar (COA)</h2>
          <span class="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-primary-container text-on-primary-container">
            Akuntansi RS
          </span>
        </div>
        <p class="text-xs text-on-surface-variant mt-0.5">
          Hierarki akun buku besar standar untuk transaksi kasir, billing, farmasi, dan jurnal keuangan.
        </p>
      </div>
    </div>

    <div class="flex items-center gap-2.5">
      <button
        type="button"
        onclick={loadData}
        class="flex items-center gap-1.5 px-3.5 py-2 text-xs font-medium text-on-surface-variant bg-white dark:bg-surface-container border border-outline-variant/60 rounded-xl hover:bg-surface-container-high transition-all cursor-pointer shadow-2xs"
        title="Muat ulang data"
      >
        <RefreshCw class="w-3.5 h-3.5 {isLoading ? 'animate-spin text-primary' : ''}" />
        <span>Segarkan</span>
      </button>

      {#if auth.hasPermission('accounting:create')}
        <M3Button variant="filled" onclick={() => openCreateModal()}>
          <Plus class="w-4 h-4" />
          <span>Tambah Akun Baru</span>
        </M3Button>
      {:else}
        <div class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-surface-container text-on-surface-variant text-xs font-medium border border-outline-variant/60">
          <ShieldAlert class="w-3.5 h-3.5 text-amber-600 shrink-0" />
          <span>Hanya Lihat</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- KPI Summary Cards -->
  <CoaStatsCards {stats} />

  <!-- Notifications -->
  {#if successMessage}
    <div class="flex items-center justify-between p-3.5 rounded-2xl bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800/60 text-emerald-800 dark:text-emerald-300 text-xs shadow-2xs">
      <div class="flex items-center gap-2.5">
        <Check class="w-4 h-4 text-emerald-600 dark:text-emerald-400 shrink-0" />
        <span class="font-medium">{successMessage}</span>
      </div>
      <button type="button" onclick={() => (successMessage = null)} class="text-emerald-600 dark:text-emerald-400 hover:text-emerald-900 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div class="flex items-center justify-between p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800/60 text-rose-800 dark:text-rose-300 text-xs shadow-2xs">
      <div class="flex items-center gap-2.5">
        <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0" />
        <span class="font-medium">{errorMessage}</span>
      </div>
      <button type="button" onclick={() => (errorMessage = null)} class="text-rose-600 dark:text-rose-400 hover:text-rose-900 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  <!-- Filters & View Switcher -->
  <div class="flex flex-col gap-3 bg-[#f8fafd] dark:bg-surface-container p-4 rounded-3xl border border-outline-variant/60">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <!-- Search Input -->
      <div class="flex-1 min-w-[240px] relative">
        <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-on-surface-variant" />
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Cari kode akun, nama akun, nama bank, atau no. rekening..."
          class="w-full pl-10 pr-4 py-2.5 text-xs bg-white dark:bg-surface-container-high border border-outline-variant/60 rounded-2xl focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all text-on-surface"
        />
        {#if searchQuery}
          <button type="button" onclick={() => (searchQuery = '')} class="absolute right-3 top-1/2 -translate-y-1/2 text-on-surface-variant hover:text-on-surface cursor-pointer">
            <X class="w-3.5 h-3.5" />
          </button>
        {/if}
      </div>

      <!-- View Switcher -->
      <div class="flex items-center gap-2">
        <M3SegmentedButton
          bind:selected={viewMode}
          options={[
            { value: 'tree', label: 'Hierarki Pohon', icon: FolderTree },
            { value: 'list', label: 'Daftar Tabel', icon: ListFilter }
          ]}
        />
        {#if viewMode === 'tree'}
          <button
            type="button"
            onclick={() => (expandedNodeIds = collectExpandedNodeIds(treeData))}
            class="px-3 py-1.5 text-[11px] font-medium text-on-surface-variant bg-white dark:bg-surface-container-high border border-outline-variant/60 rounded-xl hover:bg-surface-container-highest transition-all cursor-pointer"
          >
            Buka Semua
          </button>
          <button
            type="button"
            onclick={() => (expandedNodeIds = new Set())}
            class="px-3 py-1.5 text-[11px] font-medium text-on-surface-variant bg-white dark:bg-surface-container-high border border-outline-variant/60 rounded-xl hover:bg-surface-container-highest transition-all cursor-pointer"
          >
            Tutup Semua
          </button>
        {/if}
      </div>
    </div>

    <!-- Category Tabs & Secondary Filters -->
    <div class="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-outline-variant/40">
      <div class="flex flex-wrap items-center gap-1.5">
        {#each [
          { val: 'ALL', label: 'Semua Akun', bg: 'bg-primary text-on-primary' },
          { val: 'ASSET', label: '1. Aset', bg: 'bg-blue-600 text-white' },
          { val: 'LIABILITY', label: '2. Kewajiban', bg: 'bg-amber-600 text-white' },
          { val: 'EQUITY', label: '3. Ekuitas', bg: 'bg-purple-600 text-white' },
          { val: 'REVENUE', label: '4. Pendapatan', bg: 'bg-emerald-600 text-white' },
          { val: 'EXPENSE', label: '5. Beban', bg: 'bg-rose-600 text-white' }
        ] as tab}
          <button
            type="button"
            onclick={() => (selectedTypeFilter = tab.val as any)}
            class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === tab.val
              ? `${tab.bg} shadow-xs font-semibold`
              : 'bg-white dark:bg-surface-container-high text-on-surface-variant border border-outline-variant/60 hover:bg-surface-container-highest'}"
          >
            {tab.label}
          </button>
        {/each}
      </div>

      <div class="flex items-center gap-3">
        <select
          bind:value={postableFilter}
          class="px-2.5 py-1 text-xs bg-white dark:bg-surface-container-high border border-outline-variant/60 rounded-xl text-on-surface-variant focus:outline-none focus:border-primary"
        >
          <option value="ALL">Semua Sifat Akun</option>
          <option value="POSTABLE">Hanya Transaksi (Postable)</option>
          <option value="HEADER">Hanya Induk (Header)</option>
        </select>

        <label class="flex items-center gap-1.5 text-xs text-on-surface-variant cursor-pointer select-none">
          <input type="checkbox" bind:checked={activeOnlyFilter} class="rounded border-outline-variant text-primary focus:ring-primary" />
          <span>Hanya Aktif</span>
        </label>
      </div>
    </div>
  </div>

  <!-- Content Area -->
  {#if isLoading && treeData.length === 0}
    <div class="p-12 text-center text-on-surface-variant">
      <RefreshCw class="w-8 h-8 text-primary animate-spin mx-auto mb-3" />
      <p class="font-medium text-sm">Memuat Bagan Akun RS...</p>
    </div>
  {:else if viewMode === 'tree'}
    <CoaTreeView
      treeNodes={filteredTree}
      {expandedNodeIds}
      canCreate={auth.hasPermission('accounting:create')}
      canUpdate={auth.hasPermission('accounting:update')}
      canDelete={auth.hasPermission('accounting:delete')}
      onToggleNode={toggleNode}
      onCreateChild={openCreateModal}
      onEdit={openEditModal}
      onDelete={confirmDelete}
    />
  {:else}
    <CoaListView
      accounts={filteredFlatAccounts}
      totalAccounts={flatAccounts.length}
      canCreate={auth.hasPermission('accounting:create')}
      canUpdate={auth.hasPermission('accounting:update')}
      canDelete={auth.hasPermission('accounting:delete')}
      onCreateChild={openCreateModal}
      onEdit={openEditModal}
      onDelete={confirmDelete}
    />
  {/if}
</div>

<!-- Modals -->
<CoaFormModal
  isOpen={showModal}
  mode={modalMode}
  account={selectedAccount}
  {parentForNewChild}
  {flatAccounts}
  {isSaving}
  onClose={() => (showModal = false)}
  onSave={handleSaveAccount}
/>

<CoaDeleteModal
  isOpen={showDeleteConfirm}
  account={accountToDelete}
  {isDeleting}
  errorMessage={deleteErrorMessage}
  onClose={() => {
    showDeleteConfirm = false;
    accountToDelete = null;
    deleteErrorMessage = null;
  }}
  onConfirm={handleDeleteAccount}
/>
