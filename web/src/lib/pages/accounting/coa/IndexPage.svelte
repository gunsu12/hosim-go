<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Search,
    Plus,
    RefreshCw,
    FolderTree,
    ListFilter,
    ShieldAlert,
    Edit2,
    Trash2,
    ChevronRight,
    ChevronDown,
    Landmark,
    Check,
    X,
    AlertCircle,
    Info,
    Layers,
    PlusCircle
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
  } from '../../../api';
  import type {
    AccountRecord,
    AccountTreeNode,
    AccountType,
    AccountPosition,
    CreateAccountDTO,
    UpdateAccountDTO
  } from '../../../types/accounting/account';

  // ==========================================
  // STATE MANAGEMENT
  // ==========================================

  let viewMode = $state<string>('tree'); // 'tree' | 'list'
  let searchQuery = $state('');
  let selectedTypeFilter = $state<AccountType | 'ALL'>('ALL');
  let activeOnlyFilter = $state(false);
  let postableFilter = $state<'ALL' | 'POSTABLE' | 'HEADER'>('ALL');

  let treeData = $state<AccountTreeNode[]>([]);
  let flatAccounts = $state<AccountRecord[]>([]);
  let totalCount = $state(0);
  let currentPage = $state(1);
  let pageSize = $state(25);

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

  // Form Fields
  let formCode = $state('');
  let formName = $state('');
  let formParentId = $state<string>('');
  let formType = $state<AccountType>('ASSET');
  let formPosition = $state<AccountPosition>('DEBIT');
  let formDescription = $state('');
  let formIsPostable = $state(true);
  let formIsTreasury = $state(false);
  let formBankName = $state('');
  let formBankAccountNumber = $state('');
  let formCurrency = $state('IDR');
  let formIsActive = $state(true);
  let formError = $state<string | null>(null);

  // ==========================================
  // OPTIONS & CONFIG
  // ==========================================

  const accountTypeConfig: Record<AccountType, { label: string; prefix: string; color: string; badgeBg: string; badgeText: string; defaultPos: AccountPosition }> = {
    ASSET: {
      label: 'Aset / Aktiva',
      prefix: '1',
      color: 'blue',
      badgeBg: 'bg-blue-100 border-blue-200',
      badgeText: 'text-blue-800',
      defaultPos: 'DEBIT'
    },
    LIABILITY: {
      label: 'Kewajiban / Liabilitas',
      prefix: '2',
      color: 'amber',
      badgeBg: 'bg-amber-100 border-amber-200',
      badgeText: 'text-amber-800',
      defaultPos: 'CREDIT'
    },
    EQUITY: {
      label: 'Ekuitas / Modal',
      prefix: '3',
      color: 'purple',
      badgeBg: 'bg-purple-100 border-purple-200',
      badgeText: 'text-purple-800',
      defaultPos: 'CREDIT'
    },
    REVENUE: {
      label: 'Pendapatan / Revenue',
      prefix: '4',
      color: 'emerald',
      badgeBg: 'bg-emerald-100 border-emerald-200',
      badgeText: 'text-emerald-800',
      defaultPos: 'CREDIT'
    },
    EXPENSE: {
      label: 'Beban / Biaya',
      prefix: '5',
      color: 'rose',
      badgeBg: 'bg-rose-100 border-rose-200',
      badgeText: 'text-rose-800',
      defaultPos: 'DEBIT'
    }
  };

  // ==========================================
  // COMPUTED / DERIVED
  // ==========================================

  // KPI Metrics
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

  // Filtered Tree Data
  const filteredTree = $derived.by(() => {
    if (selectedTypeFilter === 'ALL' && !searchQuery.trim() && !activeOnlyFilter && postableFilter === 'ALL') {
      return treeData;
    }

    function matchNode(node: AccountTreeNode): boolean {
      const matchSearch = !searchQuery.trim() ||
        node.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
        node.name.toLowerCase().includes(searchQuery.toLowerCase());
      const matchType = selectedTypeFilter === 'ALL' || node.type === selectedTypeFilter;
      const matchActive = !activeOnlyFilter || node.is_active;
      const matchPostable = postableFilter === 'ALL' ||
        (postableFilter === 'POSTABLE' && node.is_postable) ||
        (postableFilter === 'HEADER' && !node.is_postable);

      return matchSearch && matchType && matchActive && matchPostable;
    }

    function filterTreeRecursive(nodes: AccountTreeNode[]): AccountTreeNode[] {
      const result: AccountTreeNode[] = [];
      for (const node of nodes) {
        const filteredChildren = node.children && node.children.length > 0
          ? filterTreeRecursive(node.children)
          : [];
        const isSelfMatch = matchNode(node);

        if (isSelfMatch || filteredChildren.length > 0) {
          result.push({
            ...node,
            children: filteredChildren
          });
        }
      }
      return result;
    }

    return filterTreeRecursive(treeData);
  });

  // Filtered Flat Accounts for Table View
  const filteredFlatAccounts = $derived.by(() => {
    return flatAccounts.filter(acc => {
      const matchSearch = !searchQuery.trim() ||
        acc.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
        acc.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (acc.bank_name && acc.bank_name.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (acc.bank_account_number && acc.bank_account_number.toLowerCase().includes(searchQuery.toLowerCase()));
      const matchType = selectedTypeFilter === 'ALL' || acc.type === selectedTypeFilter;
      const matchActive = !activeOnlyFilter || acc.is_active;
      const matchPostable = postableFilter === 'ALL' ||
        (postableFilter === 'POSTABLE' && acc.is_postable) ||
        (postableFilter === 'HEADER' && !acc.is_postable);

      return matchSearch && matchType && matchActive && matchPostable;
    });
  });

  // ==========================================
  // DATA FETCHING
  // ==========================================

  async function loadData() {
    isLoading = true;
    errorMessage = null;
    try {
      // 1. Ambil data struktur pohon hierarkis
      const treeRes = await getAccountTree(false);
      treeData = treeRes || [];

      // 2. Ambil seluruh akun flat untuk KPI, dropdown parent, dan tabel list view
      const flatRes = await getAccounts({ limit: 1000 });
      flatAccounts = flatRes.items || [];
      totalCount = flatRes.total || flatAccounts.length;

      // Inisialisasi node level 1 dan 2 agar terbuka otomatis di awal
      if (expandedNodeIds.size === 0) {
        expandLevels(treeData, 2);
      }
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memuat Bagan Akun (COA)';
    } finally {
      isLoading = false;
    }
  }

  function expandLevels(nodes: AccountTreeNode[], maxLevel: number) {
    const nextSet = new Set(expandedNodeIds);
    function traverse(list: AccountTreeNode[]) {
      for (const n of list) {
        if (n.account_level <= maxLevel) {
          nextSet.add(n.id);
        }
        if (n.children && n.children.length > 0) {
          traverse(n.children);
        }
      }
    }
    traverse(nodes);
    expandedNodeIds = nextSet;
  }

  function expandAll() {
    const nextSet = new Set<string>();
    function traverse(list: AccountTreeNode[]) {
      for (const n of list) {
        nextSet.add(n.id);
        if (n.children && n.children.length > 0) {
          traverse(n.children);
        }
      }
    }
    traverse(treeData);
    expandedNodeIds = nextSet;
  }

  function collapseAll() {
    expandedNodeIds = new Set();
  }

  function toggleNode(nodeId: string) {
    const nextSet = new Set(expandedNodeIds);
    if (nextSet.has(nodeId)) {
      nextSet.delete(nodeId);
    } else {
      nextSet.add(nodeId);
    }
    expandedNodeIds = nextSet;
  }

  // ==========================================
  // MODAL & FORM HANDLERS
  // ==========================================

  function openCreateModal(parent?: AccountRecord) {
    modalMode = 'create';
    selectedAccount = null;
    parentForNewChild = parent || null;

    formParentId = parent ? parent.id : '';
    formType = parent ? parent.type : 'ASSET';
    formPosition = parent ? parent.position : accountTypeConfig[formType].defaultPos;
    formCode = parent ? `${parent.code}` : '';
    formName = '';
    formDescription = '';
    formIsPostable = true;
    formIsTreasury = false;
    formBankName = '';
    formBankAccountNumber = '';
    formCurrency = 'IDR';
    formIsActive = true;
    formError = null;

    showModal = true;
  }

  function openEditModal(account: AccountRecord) {
    modalMode = 'edit';
    selectedAccount = account;
    parentForNewChild = null;

    formCode = account.code;
    formName = account.name;
    formParentId = account.parent_id || '';
    formType = account.type;
    formPosition = account.position;
    formDescription = account.description || '';
    formIsPostable = account.is_postable;
    formIsTreasury = account.is_treasury_account;
    formBankName = account.bank_name || '';
    formBankAccountNumber = account.bank_account_number || '';
    formCurrency = account.currency || 'IDR';
    formIsActive = account.is_active;
    formError = null;

    showModal = true;
  }

  // Otomatis sinkronisasi posisi saldo normal saat tipe akun diubah (Create Mode)
  function handleTypeChange(newType: AccountType) {
    formType = newType;
    if (accountTypeConfig[newType]) {
      formPosition = accountTypeConfig[newType].defaultPos;
    }
  }

  // Otomatis sinkronisasi tipe akun dari parent jika parent dipilih
  function handleParentSelect(parentId: string) {
    formParentId = parentId;
    if (parentId) {
      const parentAcc = flatAccounts.find(a => a.id === parentId);
      if (parentAcc) {
        formType = parentAcc.type;
        formPosition = parentAcc.position;
      }
    }
  }

  async function handleSaveAccount() {
    formError = null;

    if (!formCode.trim()) {
      formError = 'Kode akun wajib diisi';
      return;
    }
    if (!formName.trim()) {
      formError = 'Nama akun wajib diisi';
      return;
    }
    if (formIsTreasury && !formBankName.trim()) {
      formError = 'Nama Bank / Kas wajib diisi untuk Akun Kas & Bank';
      return;
    }

    isSaving = true;
    try {
      if (modalMode === 'create') {
        const payload: CreateAccountDTO = {
          code: formCode.trim(),
          name: formName.trim(),
          parent_id: formParentId.trim() ? formParentId : null,
          type: formType,
          position: formPosition,
          description: formDescription.trim() ? formDescription.trim() : null,
          is_postable: formIsPostable,
          is_treasury_account: formIsTreasury,
          bank_name: formIsTreasury && formBankName.trim() ? formBankName.trim() : null,
          bank_account_number: formIsTreasury && formBankAccountNumber.trim() ? formBankAccountNumber.trim() : null,
          currency: formCurrency.trim() || 'IDR',
          is_active: formIsActive
        };

        const created = await createAccount(payload);
        successMessage = `Akun ${created.code} - ${created.name} berhasil dibuat!`;
      } else if (modalMode === 'edit' && selectedAccount) {
        const payload: UpdateAccountDTO = {
          name: formName.trim(),
          position: formPosition,
          description: formDescription.trim() ? formDescription.trim() : null,
          is_treasury_account: formIsTreasury,
          bank_name: formIsTreasury && formBankName.trim() ? formBankName.trim() : null,
          bank_account_number: formIsTreasury && formBankAccountNumber.trim() ? formBankAccountNumber.trim() : null,
          currency: formCurrency.trim() || 'IDR',
          is_active: formIsActive
        };

        const updated = await updateAccount(selectedAccount.id, payload);
        successMessage = `Akun ${updated.code} - ${updated.name} berhasil diperbarui!`;
      }

      showModal = false;
      await loadData();
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan data akun';
    } finally {
      isSaving = false;
    }
  }

  // ==========================================
  // DELETE HANDLERS
  // ==========================================

  function confirmDelete(account: AccountRecord) {
    accountToDelete = account;
    showDeleteConfirm = true;
  }

  async function handleDeleteAccount() {
    if (!accountToDelete) return;
    isDeleting = true;
    errorMessage = null;
    try {
      await deleteAccount(accountToDelete.id);
      successMessage = `Akun ${accountToDelete.code} - ${accountToDelete.name} berhasil dihapus.`;
      showDeleteConfirm = false;
      accountToDelete = null;
      await loadData();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal menghapus akun. Pastikan akun tidak memiliki sub-akun anak atau transaksi aktif.';
      showDeleteConfirm = false;
    } finally {
      isDeleting = false;
    }
  }

  onMount(() => {
    loadData();
  });
</script>

<div class="flex flex-col gap-5">
  <!-- 1. Header Banner & Action Bar -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <div class="flex items-center gap-2.5">
        <div class="p-2 rounded-xl bg-[#0b57d0]/10 text-[#0b57d0]">
          <FolderTree class="w-5 h-5" />
        </div>
        <div>
          <div class="flex items-center gap-2">
            <h2 class="text-xl font-bold text-[#1f1f1f] tracking-tight">Bagan Akun Standar (Chart of Accounts)</h2>
            <span class="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
              Akuntansi RS
            </span>
          </div>
          <p class="text-xs text-[#444746] mt-0.5">
            Daftar hierarki akun buku besar standar rumah sakit untuk pencatatan transaksi kasir, billing, farmasi, dan jurnal keuangan.
          </p>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2.5">
      <button
        type="button"
        onclick={loadData}
        class="flex items-center gap-1.5 px-3.5 py-2 text-xs font-medium text-[#444746] bg-white border border-[#e1e5ea] rounded-xl hover:bg-[#f0f4f9] transition-all cursor-pointer shadow-2xs"
        title="Muat ulang data dari server"
      >
        <RefreshCw class="w-3.5 h-3.5 {isLoading ? 'animate-spin text-[#0b57d0]' : ''}" />
        <span>Segarkan</span>
      </button>

      {#if auth.hasPermission('accounting:create')}
        <M3Button variant="filled" onclick={() => openCreateModal()}>
          <Plus class="w-4 h-4" />
          <span>Tambah Akun Baru</span>
        </M3Button>
      {:else}
        <div class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-[#f0f4f9] text-[#747775] text-xs font-medium border border-[#e1e5ea]">
          <ShieldAlert class="w-3.5 h-3.5 text-amber-600 shrink-0" />
          <span>Hanya Lihat (Read-Only)</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- 2. KPI Summary Cards -->
  <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
    <div class="p-3.5 rounded-2xl bg-white border border-[#e1e5ea] shadow-2xs flex flex-col justify-between">
      <div class="flex items-center justify-between text-[#747775]">
        <span class="text-xs font-medium">Total Akun</span>
        <Layers class="w-4 h-4 text-[#0b57d0]" />
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold text-[#1f1f1f]">{stats.total}</span>
        <span class="text-[11px] text-[#747775]">({stats.headers} Induk, {stats.postable} Pos)</span>
      </div>
    </div>

    <div class="p-3.5 rounded-2xl bg-white border border-[#e1e5ea] shadow-2xs flex flex-col justify-between">
      <div class="flex items-center justify-between text-[#747775]">
        <span class="text-xs font-medium">1. Aset (Aktiva)</span>
        <span class="w-2.5 h-2.5 rounded-full bg-blue-500"></span>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold text-blue-700">{stats.byType.ASSET}</span>
        <span class="text-[11px] text-blue-600 font-mono">1xxx (Debit)</span>
      </div>
    </div>

    <div class="p-3.5 rounded-2xl bg-white border border-[#e1e5ea] shadow-2xs flex flex-col justify-between">
      <div class="flex items-center justify-between text-[#747775]">
        <span class="text-xs font-medium">2 & 3. Kewajiban & Modal</span>
        <span class="w-2.5 h-2.5 rounded-full bg-amber-500"></span>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold text-amber-700">{stats.byType.LIABILITY + stats.byType.EQUITY}</span>
        <span class="text-[11px] text-amber-600 font-mono">2xxx / 3xxx (Kredit)</span>
      </div>
    </div>

    <div class="p-3.5 rounded-2xl bg-white border border-[#e1e5ea] shadow-2xs flex flex-col justify-between">
      <div class="flex items-center justify-between text-[#747775]">
        <span class="text-xs font-medium">4 & 5. Pendapatan & Beban</span>
        <span class="w-2.5 h-2.5 rounded-full bg-emerald-500"></span>
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold text-emerald-700">{stats.byType.REVENUE + stats.byType.EXPENSE}</span>
        <span class="text-[11px] text-emerald-600 font-mono">4xxx / 5xxx</span>
      </div>
    </div>

    <div class="p-3.5 rounded-2xl bg-white border border-[#e1e5ea] shadow-2xs flex flex-col justify-between col-span-2 sm:col-span-1">
      <div class="flex items-center justify-between text-[#747775]">
        <span class="text-xs font-medium">Kas & Bank (Treasury)</span>
        <Landmark class="w-4 h-4 text-teal-600" />
      </div>
      <div class="mt-2 flex items-baseline gap-2">
        <span class="text-2xl font-bold text-teal-700">{stats.treasury}</span>
        <span class="text-[11px] text-teal-600">Rekening RS</span>
      </div>
    </div>
  </div>

  <!-- Notification Toasts -->
  {#if successMessage}
    <div class="flex items-center justify-between p-3.5 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-800 text-xs shadow-2xs animate-fadeIn">
      <div class="flex items-center gap-2.5">
        <Check class="w-4 h-4 text-emerald-600 shrink-0" />
        <span class="font-medium">{successMessage}</span>
      </div>
      <button type="button" onclick={() => (successMessage = null)} class="text-emerald-600 hover:text-emerald-900 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div class="flex items-center justify-between p-3.5 rounded-2xl bg-rose-50 border border-rose-200 text-rose-800 text-xs shadow-2xs animate-fadeIn">
      <div class="flex items-center gap-2.5">
        <AlertCircle class="w-4 h-4 text-rose-600 shrink-0" />
        <span class="font-medium">{errorMessage}</span>
      </div>
      <button type="button" onclick={() => (errorMessage = null)} class="text-rose-600 hover:text-rose-900 cursor-pointer">
        <X class="w-4 h-4" />
      </button>
    </div>
  {/if}

  <!-- 3. Filter Controls & View Switcher -->
  <div class="flex flex-col gap-3 bg-[#f8fafd] p-4 rounded-3xl border border-[#e1e5ea]">
    <!-- Row 1: Search, Mode Switcher, & Tree Buttons -->
    <div class="flex flex-wrap items-center justify-between gap-3">
      <!-- Search Input -->
      <div class="flex-1 min-w-[240px] relative">
        <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775]" />
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Cari kode akun, nama akun, nama bank, atau no. rekening..."
          class="w-full pl-10 pr-4 py-2.5 text-xs bg-white border border-[#e1e5ea] rounded-2xl focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0] transition-all"
        />
        {#if searchQuery}
          <button
            type="button"
            onclick={() => searchQuery = ''}
            class="absolute right-3 top-1/2 -translate-y-1/2 text-[#747775] hover:text-[#1f1f1f] cursor-pointer"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        {/if}
      </div>

      <!-- View Switcher (Tree vs Table) -->
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
            onclick={expandAll}
            class="px-3 py-1.5 text-[11px] font-medium text-[#444746] bg-white border border-[#e1e5ea] rounded-xl hover:bg-[#f0f4f9] transition-all cursor-pointer"
            title="Buka semua cabang pohon"
          >
            Buka Semua
          </button>
          <button
            type="button"
            onclick={collapseAll}
            class="px-3 py-1.5 text-[11px] font-medium text-[#444746] bg-white border border-[#e1e5ea] rounded-xl hover:bg-[#f0f4f9] transition-all cursor-pointer"
            title="Tutup semua cabang pohon"
          >
            Tutup Semua
          </button>
        {/if}
      </div>
    </div>

    <!-- Row 2: Account Type Filter Tabs & Secondary Filters -->
    <div class="flex flex-wrap items-center justify-between gap-3 pt-2 border-t border-[#e1e5ea]/60">
      <!-- Tabs Kategori Akun -->
      <div class="flex flex-wrap items-center gap-1.5">
        <button
          type="button"
          onclick={() => selectedTypeFilter = 'ALL'}
          class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === 'ALL' ? 'bg-[#0b57d0] text-white shadow-xs font-semibold' : 'bg-white text-[#444746] border border-[#e1e5ea] hover:bg-[#f0f4f9]'}"
        >
          Semua Akun
        </button>
        <button
          type="button"
          onclick={() => selectedTypeFilter = 'ASSET'}
          class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === 'ASSET' ? 'bg-blue-600 text-white shadow-xs font-semibold' : 'bg-white text-blue-800 border border-blue-200 hover:bg-blue-50'}"
        >
          1. Aset
        </button>
        <button
          type="button"
          onclick={() => selectedTypeFilter = 'LIABILITY'}
          class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === 'LIABILITY' ? 'bg-amber-600 text-white shadow-xs font-semibold' : 'bg-white text-amber-800 border border-amber-200 hover:bg-amber-50'}"
        >
          2. Kewajiban
        </button>
        <button
          type="button"
          onclick={() => selectedTypeFilter = 'EQUITY'}
          class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === 'EQUITY' ? 'bg-purple-600 text-white shadow-xs font-semibold' : 'bg-white text-purple-800 border border-purple-200 hover:bg-purple-50'}"
        >
          3. Ekuitas
        </button>
        <button
          type="button"
          onclick={() => selectedTypeFilter = 'REVENUE'}
          class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === 'REVENUE' ? 'bg-emerald-600 text-white shadow-xs font-semibold' : 'bg-white text-emerald-800 border border-emerald-200 hover:bg-emerald-50'}"
        >
          4. Pendapatan
        </button>
        <button
          type="button"
          onclick={() => selectedTypeFilter = 'EXPENSE'}
          class="px-3 py-1 text-xs rounded-full font-medium transition-all cursor-pointer {selectedTypeFilter === 'EXPENSE' ? 'bg-rose-600 text-white shadow-xs font-semibold' : 'bg-white text-rose-800 border border-rose-200 hover:bg-rose-50'}"
        >
          5. Beban
        </button>
      </div>

      <!-- Secondary Filters -->
      <div class="flex items-center gap-3">
        <!-- Filter Postable -->
        <select
          bind:value={postableFilter}
          class="px-2.5 py-1 text-xs bg-white border border-[#e1e5ea] rounded-xl text-[#444746] focus:outline-none focus:border-[#0b57d0]"
        >
          <option value="ALL">Semua Sifat Akun</option>
          <option value="POSTABLE">Hanya Akun Transaksi (Postable)</option>
          <option value="HEADER">Hanya Akun Induk (Header)</option>
        </select>

        <!-- Toggle Hanya Aktif -->
        <label class="flex items-center gap-1.5 text-xs text-[#444746] cursor-pointer select-none">
          <input
            type="checkbox"
            bind:checked={activeOnlyFilter}
            class="rounded border-[#c4c7c5] text-[#0b57d0] focus:ring-[#0b57d0]"
          />
          <span>Hanya Aktif</span>
        </label>
      </div>
    </div>
  </div>

  <!-- 4. MAIN CONTENT AREA -->
  {#if isLoading && treeData.length === 0}
    <div class="p-12 text-center text-[#747775]">
      <RefreshCw class="w-8 h-8 text-[#0b57d0] animate-spin mx-auto mb-3" />
      <p class="font-medium text-sm">Memuat Bagan Akun RS...</p>
      <p class="text-xs text-[#8e9196] mt-0.5">Mengambil struktur hierarki dan konfigurasi dari server</p>
    </div>

  {:else if viewMode === 'tree'}
    <!-- ========================================== -->
    <!-- VIEW A: HIERARCHICAL TREE VIEW -->
    <!-- ========================================== -->
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xs overflow-hidden">
      <!-- Tree Table Header -->
      <div class="grid grid-cols-12 gap-3 px-4 py-3 bg-[#f8fafd] border-b border-[#e1e5ea] text-xs font-semibold text-[#444746]">
        <div class="col-span-6 sm:col-span-5 flex items-center gap-2">
          <span>Struktur Akun & Nama Buku Besar</span>
        </div>
        <div class="col-span-2 text-center">
          <span>Tipe & Golongan</span>
        </div>
        <div class="col-span-1 text-center">
          <span>Posisi</span>
        </div>
        <div class="col-span-1 text-center hidden sm:block">
          <span>Level</span>
        </div>
        <div class="col-span-2 text-center">
          <span>Sifat & Bank</span>
        </div>
        <div class="col-span-1 text-right">
          <span>Aksi</span>
        </div>
      </div>

      <!-- Recursive Tree Snippet Definition -->
      {#snippet renderTreeNode(node: AccountTreeNode, depth: number)}
        {@const hasChildren = node.children && node.children.length > 0}
        {@const isExpanded = expandedNodeIds.has(node.id)}
        {@const typeMeta = accountTypeConfig[node.type] || accountTypeConfig.ASSET}
        {@const isHeader = !node.is_postable}

        <div
          class="grid grid-cols-12 gap-3 px-4 py-2.5 items-center border-b border-[#e1e5ea]/50 hover:bg-[#f8fafd] transition-colors {isHeader ? 'bg-[#fafcff]/50 font-medium' : ''}"
        >
          <!-- Column 1: Tree Indent, Caret, Code & Name -->
          <div class="col-span-6 sm:col-span-5 flex items-center min-w-0" style="padding-left: {depth * 20}px;">
            <!-- Expand / Collapse Toggle -->
            {#if hasChildren}
              <button
                type="button"
                onclick={() => toggleNode(node.id)}
                class="p-1 -ml-1 mr-1 text-[#747775] hover:text-[#0b57d0] hover:bg-[#e8f0fe] rounded-md transition-colors cursor-pointer shrink-0"
                title={isExpanded ? 'Tutup cabang' : 'Buka cabang'}
              >
                {#if isExpanded}
                  <ChevronDown class="w-4 h-4 text-[#0b57d0]" />
                {:else}
                  <ChevronRight class="w-4 h-4" />
                {/if}
              </button>
            {:else}
              <span class="w-6 shrink-0 inline-block text-center text-[#c4c7c5] text-xs">•</span>
            {/if}

            <!-- Account Code & Name -->
            <div class="flex items-baseline gap-2 min-w-0 truncate">
              <span class="font-mono text-xs font-bold {isHeader ? 'text-[#0b57d0]' : 'text-[#444746]'} shrink-0">
                {node.code}
              </span>
              <span class="text-xs truncate {isHeader ? 'font-bold text-[#1f1f1f]' : 'text-[#1f1f1f]'}" title={node.name}>
                {node.name}
              </span>
            </div>
          </div>

          <!-- Column 2: Account Type Badge -->
          <div class="col-span-2 text-center">
            <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border {typeMeta.badgeBg} {typeMeta.badgeText}">
              {node.type}
            </span>
          </div>

          <!-- Column 3: Normal Balance Position -->
          <div class="col-span-1 text-center">
            {#if node.position === 'DEBIT'}
              <span class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-teal-50 text-teal-700 border border-teal-200">
                D
              </span>
            {:else}
              <span class="px-1.5 py-0.5 rounded text-[10px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-200">
                K
              </span>
            {/if}
          </div>

          <!-- Column 4: Level -->
          <div class="col-span-1 text-center hidden sm:block">
            <span class="text-[11px] font-mono text-[#747775]">L{node.account_level}</span>
          </div>

          <!-- Column 5: Postable / Treasury / Bank Info -->
          <div class="col-span-2 text-center flex flex-wrap items-center justify-center gap-1">
            {#if node.is_postable}
              <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                Postable
              </span>
            {:else}
              <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600">
                Induk
              </span>
            {/if}

            {#if node.is_treasury_account}
              <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-teal-100 text-teal-800 flex items-center gap-1" title="{node.bank_name || 'Bank'}: {node.bank_account_number || ''}">
                <Landmark class="w-2.5 h-2.5" />
                <span>Kas/Bank</span>
              </span>
            {/if}

            {#if !node.is_active}
              <span class="px-1.5 py-0.5 rounded text-[9px] font-bold bg-rose-100 text-rose-700">
                Non-Aktif
              </span>
            {/if}
          </div>

          <!-- Column 6: Action Buttons -->
          <div class="col-span-1 text-right">
            <div class="flex items-center justify-end gap-1">
              {#if auth.hasPermission('accounting:create')}
                <button
                  type="button"
                  onclick={() => openCreateModal(node)}
                  class="p-1 text-emerald-600 hover:bg-emerald-50 rounded-lg transition-colors cursor-pointer"
                  title="Tambah sub-akun anak di bawah {node.code}"
                >
                  <PlusCircle class="w-3.5 h-3.5" />
                </button>
              {/if}

              {#if auth.hasPermission('accounting:update')}
                <button
                  type="button"
                  onclick={() => openEditModal(node)}
                  class="p-1 text-[#0b57d0] hover:bg-[#e8f0fe] rounded-lg transition-colors cursor-pointer"
                  title="Edit data akun"
                >
                  <Edit2 class="w-3.5 h-3.5" />
                </button>
              {/if}

              {#if auth.hasPermission('accounting:delete')}
                <button
                  type="button"
                  onclick={() => confirmDelete(node)}
                  class="p-1 text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
                  title="Hapus akun"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              {/if}
            </div>
          </div>
        </div>

        <!-- Render Children Recursively if Expanded -->
        {#if hasChildren && isExpanded}
          {#each node.children as child (child.id)}
            {@render renderTreeNode(child, depth + 1)}
          {/each}
        {/if}
      {/snippet}

      <!-- Tree Nodes Body -->
      {#if filteredTree.length === 0}
        <div class="p-12 text-center text-[#747775]">
          <FolderTree class="w-10 h-10 text-[#a8abb0] mx-auto mb-2 opacity-50" />
          <p class="font-medium text-sm">Tidak ada akun yang sesuai dengan filter</p>
          <p class="text-xs text-[#8e9196] mt-0.5">Coba atur ulang pencarian atau pilih kategori akun lainnya.</p>
        </div>
      {:else}
        <div class="divide-y divide-[#e1e5ea]/30">
          {#each filteredTree as rootNode (rootNode.id)}
            {@render renderTreeNode(rootNode, 0)}
          {/each}
        </div>
      {/if}
    </div>

  {:else}
    <!-- ========================================== -->
    <!-- VIEW B: FLAT TABLE VIEW -->
    <!-- ========================================== -->
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xs overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse text-xs">
          <thead>
            <tr class="bg-[#f8fafd] border-b border-[#e1e5ea] text-[#444746] select-none">
              <th class="py-3 px-4 font-semibold">Kode Akun</th>
              <th class="py-3 px-4 font-semibold">Nama Akun</th>
              <th class="py-3 px-4 font-semibold">Tipe Akun</th>
              <th class="py-3 px-4 font-semibold text-center">Posisi</th>
              <th class="py-3 px-4 font-semibold text-center">Level</th>
              <th class="py-3 px-4 font-semibold">Akun Induk</th>
              <th class="py-3 px-4 font-semibold text-center">Sifat</th>
              <th class="py-3 px-4 font-semibold">Rekening Kas/Bank</th>
              <th class="py-3 px-4 font-semibold text-center">Status</th>
              <th class="py-3 px-4 font-semibold text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
            {#if filteredFlatAccounts.length === 0}
              <tr>
                <td colspan="10" class="py-12 text-center text-[#747775]">
                  <FolderTree class="w-8 h-8 text-[#a8abb0] mx-auto mb-2 opacity-50" />
                  <p class="font-medium">Tidak ada data akun ditemukan</p>
                </td>
              </tr>
            {:else}
              {#each filteredFlatAccounts as acc (acc.id)}
                {@const typeMeta = accountTypeConfig[acc.type] || accountTypeConfig.ASSET}
                <tr class="hover:bg-[#f8fafd] transition-colors {acc.is_postable ? '' : 'bg-[#fafcff]/60 font-medium'}">
                  <td class="py-3 px-4 font-mono font-bold {acc.is_postable ? 'text-[#444746]' : 'text-[#0b57d0]'}">
                    {acc.code}
                  </td>
                  <td class="py-3 px-4">
                    <div class="font-semibold text-[#1f1f1f]">{acc.name}</div>
                    {#if acc.description}
                      <div class="text-[11px] text-[#747775] truncate max-w-xs">{acc.description}</div>
                    {/if}
                  </td>
                  <td class="py-3 px-4">
                    <span class="inline-flex items-center px-2 py-0.5 rounded-md text-[10px] font-semibold border {typeMeta.badgeBg} {typeMeta.badgeText}">
                      {acc.type}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-center">
                    {#if acc.position === 'DEBIT'}
                      <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-teal-50 text-teal-700 border border-teal-200">
                        DEBIT
                      </span>
                    {:else}
                      <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-indigo-50 text-indigo-700 border border-indigo-200">
                        KREDIT
                      </span>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-center font-mono text-[11px] text-[#747775]">
                    L{acc.account_level}
                  </td>
                  <td class="py-3 px-4 font-mono text-xs text-[#747775]">
                    {acc.parent_code || '-'}
                  </td>
                  <td class="py-3 px-4 text-center">
                    {#if acc.is_postable}
                      <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                        Postable
                      </span>
                    {:else}
                      <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600">
                        Induk
                      </span>
                    {/if}
                  </td>
                  <td class="py-3 px-4">
                    {#if acc.is_treasury_account}
                      <div class="flex items-center gap-1.5 text-xs text-teal-800 font-medium">
                        <Landmark class="w-3.5 h-3.5 text-teal-600 shrink-0" />
                        <div>
                          <div>{acc.bank_name || 'Kas RS'}</div>
                          {#if acc.bank_account_number}
                            <div class="text-[10px] font-mono text-[#747775]">{acc.bank_account_number}</div>
                          {/if}
                        </div>
                      </div>
                    {:else}
                      <span class="text-[#a8abb0]">-</span>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-center">
                    {#if acc.is_active}
                      <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                        AKTIF
                      </span>
                    {:else}
                      <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-gray-100 text-gray-600">
                        NON-AKTIF
                      </span>
                    {/if}
                  </td>
                  <td class="py-3 px-4 text-right">
                    <div class="flex items-center justify-end gap-1.5">
                      {#if auth.hasPermission('accounting:create')}
                        <button
                          type="button"
                          onclick={() => openCreateModal(acc)}
                          class="p-1.5 text-emerald-600 hover:bg-emerald-50 rounded-lg transition-colors cursor-pointer"
                          title="Tambah sub-akun"
                        >
                          <PlusCircle class="w-3.5 h-3.5" />
                        </button>
                      {/if}
                      {#if auth.hasPermission('accounting:update')}
                        <button
                          type="button"
                          onclick={() => openEditModal(acc)}
                          class="p-1.5 text-[#0b57d0] hover:bg-[#e8f0fe] rounded-lg transition-colors cursor-pointer"
                          title="Edit Akun"
                        >
                          <Edit2 class="w-3.5 h-3.5" />
                        </button>
                      {/if}
                      {#if auth.hasPermission('accounting:delete')}
                        <button
                          type="button"
                          onclick={() => confirmDelete(acc)}
                          class="p-1.5 text-rose-600 hover:bg-rose-50 rounded-lg transition-colors cursor-pointer"
                          title="Hapus Akun"
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
      <div class="p-3 bg-[#f8fafd] border-t border-[#e1e5ea] text-xs text-[#747775] flex items-center justify-between">
        <span>Menampilkan {filteredFlatAccounts.length} dari total {flatAccounts.length} akun</span>
        <span class="font-mono text-[11px]">HOSIM-GO Chart of Accounts Engine</span>
      </div>
    </div>
  {/if}
</div>

<!-- ========================================== -->
<!-- MODAL FORM INPUT & EDIT AKUN -->
<!-- ========================================== -->
{#if showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-fadeIn">
    <div class="w-full max-w-2xl bg-white rounded-3xl shadow-xl border border-[#e1e5ea] overflow-hidden flex flex-col max-h-[92vh]">
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-[#e1e5ea] bg-[#f8fafd]">
        <div class="flex items-center gap-2.5">
          <div class="p-2 rounded-xl {modalMode === 'create' ? 'bg-blue-50 text-blue-700' : 'bg-amber-50 text-amber-700'}">
            {#if modalMode === 'create'}
              <Plus class="w-5 h-5" />
            {:else}
              <Edit2 class="w-5 h-5" />
            {/if}
          </div>
          <div>
            <h3 class="text-base font-bold text-[#1f1f1f]">
              {modalMode === 'create' ? 'Tambah Akun Buku Besar Baru' : `Edit Akun: ${formCode}`}
            </h3>
            <p class="text-xs text-[#444746]">
              {modalMode === 'create' ? 'Konfigurasi akun baru pada bagan akun standar RS' : 'Perbarui nama, posisi saldo, atau informasi perbankan'}
            </p>
          </div>
        </div>
        <button
          type="button"
          onclick={() => showModal = false}
          class="p-1.5 text-[#747775] hover:text-[#1f1f1f] hover:bg-[#e1e5ea]/50 rounded-xl transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Modal Body (Scrollable) -->
      <div class="p-6 overflow-y-auto flex flex-col gap-4 text-xs">
        {#if formError}
          <div class="p-3 rounded-xl bg-rose-50 border border-rose-200 text-rose-800 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0 text-rose-600" />
            <span>{formError}</span>
          </div>
        {/if}

        <!-- Mode Tambah Sub-Akun Banner Info -->
        {#if modalMode === 'create' && parentForNewChild}
          <div class="p-3 rounded-2xl bg-blue-50 border border-blue-200 text-blue-900 flex items-center justify-between">
            <div class="flex items-center gap-2">
              <Info class="w-4 h-4 text-blue-700 shrink-0" />
              <span>Membuat sub-akun di bawah induk: <strong>{parentForNewChild.code} - {parentForNewChild.name}</strong></span>
            </div>
            <span class="text-[10px] font-mono px-2 py-0.5 rounded bg-blue-100 text-blue-800 font-bold">
              Level {parentForNewChild.account_level + 1}
            </span>
          </div>
        {/if}

        <!-- Section 1: Kode & Nama Akun -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label for="coa-form-code" class="block font-semibold text-[#1f1f1f] mb-1.5">
              Kode Akun <span class="text-rose-500">*</span>
            </label>
            <input
              id="coa-form-code"
              type="text"
              bind:value={formCode}
              disabled={modalMode === 'edit'}
              placeholder="Contoh: 11101 atau 111.01"
              class="w-full px-3 py-2 text-xs bg-white border border-[#e1e5ea] rounded-xl font-mono focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0] disabled:bg-gray-100 disabled:text-gray-500"
            />
            {#if modalMode === 'edit'}
              <span class="text-[10px] text-[#747775] mt-1 block">
                Kode akun bersifat tetap (immutable) untuk menjaga konsistensi buku besar.
              </span>
            {/if}
          </div>

          <div>
            <label for="coa-form-name" class="block font-semibold text-[#1f1f1f] mb-1.5">
              Nama Akun <span class="text-rose-500">*</span>
            </label>
            <input
              id="coa-form-name"
              type="text"
              bind:value={formName}
              placeholder="Contoh: Kas Kasir Rawat Jalan"
              class="w-full px-3 py-2 text-xs bg-white border border-[#e1e5ea] rounded-xl focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0]"
            />
          </div>
        </div>

        <!-- Section 2: Parent Akun & Tipe Akun -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <!-- Dropdown Parent Akun -->
          <div>
            <label for="coa-form-parent" class="block font-semibold text-[#1f1f1f] mb-1.5">
              Akun Induk (Parent)
            </label>
            <select
              id="coa-form-parent"
              value={formParentId}
              onchange={(e) => handleParentSelect((e.target as HTMLSelectElement).value)}
              disabled={modalMode === 'edit'}
              class="w-full px-3 py-2 text-xs bg-white border border-[#e1e5ea] rounded-xl focus:outline-none focus:border-[#0b57d0] disabled:bg-gray-100 disabled:text-gray-500"
            >
              <option value="">(Tanpa Induk - Akun Kepala Level 1)</option>
              {#each flatAccounts as pAcc}
                {#if modalMode === 'create' || pAcc.id !== selectedAccount?.id}
                  <option value={pAcc.id}>
                    [{pAcc.code}] {pAcc.name} (L{pAcc.account_level} • {pAcc.type})
                  </option>
                {/if}
              {/each}
            </select>
            {#if modalMode === 'edit'}
              <span class="text-[10px] text-[#747775] mt-1 block">
                Struktur hierarki induk tidak dapat dipindahkan setelah akun memiliki buku besar.
              </span>
            {/if}
          </div>

          <!-- Tipe Akun -->
          <div>
            <label for="coa-form-type" class="block font-semibold text-[#1f1f1f] mb-1.5">
              Klasifikasi Tipe Akun <span class="text-rose-500">*</span>
            </label>
            <select
              id="coa-form-type"
              value={formType}
              onchange={(e) => handleTypeChange((e.target as HTMLSelectElement).value as AccountType)}
              disabled={modalMode === 'edit' || (modalMode === 'create' && !!formParentId)}
              class="w-full px-3 py-2 text-xs bg-white border border-[#e1e5ea] rounded-xl focus:outline-none focus:border-[#0b57d0] disabled:bg-gray-100 disabled:text-gray-500"
            >
              <option value="ASSET">ASSET (1xxx - Aset / Aktiva)</option>
              <option value="LIABILITY">LIABILITY (2xxx - Kewajiban / Hutang)</option>
              <option value="EQUITY">EQUITY (3xxx - Ekuitas / Modal)</option>
              <option value="REVENUE">REVENUE (4xxx - Pendapatan)</option>
              <option value="EXPENSE">EXPENSE (5xxx - Beban Operasional)</option>
            </select>
          </div>
        </div>

        <!-- Section 3: Posisi Saldo Normal & Sifat Postable -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 rounded-2xl bg-[#f8fafd] border border-[#e1e5ea]">
          <div>
            <span class="block font-semibold text-[#1f1f1f] mb-1.5">
              Posisi Saldo Normal <span class="text-rose-500">*</span>
            </span>
            <div class="flex items-center gap-3">
              <label class="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="position"
                  value="DEBIT"
                  checked={formPosition === 'DEBIT'}
                  onchange={() => formPosition = 'DEBIT'}
                  class="text-[#0b57d0] focus:ring-[#0b57d0]"
                />
                <span class="font-medium">DEBIT</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="position"
                  value="CREDIT"
                  checked={formPosition === 'CREDIT'}
                  onchange={() => formPosition = 'CREDIT'}
                  class="text-[#0b57d0] focus:ring-[#0b57d0]"
                />
                <span class="font-medium">KREDIT</span>
              </label>
            </div>
            <span class="text-[10px] text-[#747775] mt-1 block">
              Standar: Aset & Beban bertambah di Debit; Kewajiban, Ekuitas, & Pendapatan bertambah di Kredit.
            </span>
          </div>

          <div>
            <span class="block font-semibold text-[#1f1f1f] mb-1.5">
              Sifat Akun Transaksi
            </span>
            <label class="flex items-start gap-2.5 cursor-pointer mt-1">
              <input
                type="checkbox"
                bind:checked={formIsPostable}
                disabled={modalMode === 'edit'}
                class="mt-0.5 rounded border-[#c4c7c5] text-[#0b57d0] focus:ring-[#0b57d0] disabled:opacity-50"
              />
              <div>
                <span class="font-medium text-[#1f1f1f]">Akun Postable (Dapat Dijurnal)</span>
                <span class="text-[10px] text-[#747775] block">
                  Hilangkan centang jika akun ini hanya sebagai Akun Induk / Header ringkasan.
                </span>
              </div>
            </label>
          </div>
        </div>

        <!-- Section 4: Kas & Bank / Treasury Flag -->
        <div class="p-4 rounded-2xl border border-teal-200 bg-teal-50/40 flex flex-col gap-3">
          <label class="flex items-center gap-2.5 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={formIsTreasury}
              class="rounded border-teal-300 text-teal-600 focus:ring-teal-600"
            />
            <div class="flex items-center gap-1.5">
              <Landmark class="w-4 h-4 text-teal-700" />
              <span class="font-semibold text-teal-900">Merupakan Akun Kas / Bank (Treasury Account)</span>
            </div>
          </label>
          <p class="text-[11px] text-teal-800">
            Centang opsi ini jika akun ini mewakili kas tunai, rekening giro operasional, rekening penampung BPJS, atau kas kecil RS.
          </p>

          {#if formIsTreasury}
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2 border-t border-teal-200/60 animate-fadeIn">
              <div>
                <label for="coa-form-bank-name" class="block font-medium text-teal-950 mb-1">
                  Nama Bank / Entitas Kas <span class="text-rose-500">*</span>
                </label>
                <input
                  id="coa-form-bank-name"
                  type="text"
                  bind:value={formBankName}
                  placeholder="Contoh: Bank Mandiri, BCA, atau Kas Induk RS"
                  class="w-full px-3 py-2 text-xs bg-white border border-teal-300 rounded-xl focus:outline-none focus:border-teal-600"
                />
              </div>

              <div>
                <label for="coa-form-bank-account" class="block font-medium text-teal-950 mb-1">
                  Nomor Rekening Bank
                </label>
                <input
                  id="coa-form-bank-account"
                  type="text"
                  bind:value={formBankAccountNumber}
                  placeholder="Contoh: 137-00-1234567-8"
                  class="w-full px-3 py-2 text-xs bg-white border border-teal-300 rounded-xl font-mono focus:outline-none focus:border-teal-600"
                />
              </div>
            </div>
          {/if}
        </div>

        <!-- Section 5: Keterangan & Status -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="sm:col-span-2">
            <label for="coa-form-desc" class="block font-semibold text-[#1f1f1f] mb-1.5">
              Keterangan / Deskripsi Akun
            </label>
            <textarea
              id="coa-form-desc"
              bind:value={formDescription}
              rows="2"
              placeholder="Penjelasan fungsi akun buku besar ini..."
              class="w-full px-3 py-2 text-xs bg-white border border-[#e1e5ea] rounded-xl focus:outline-none focus:border-[#0b57d0]"
            ></textarea>
          </div>

          <div class="flex flex-col justify-between">
            <div>
              <label for="coa-form-currency" class="block font-semibold text-[#1f1f1f] mb-1.5">
                Mata Uang
              </label>
              <input
                id="coa-form-currency"
                type="text"
                bind:value={formCurrency}
                class="w-full px-3 py-2 text-xs bg-white border border-[#e1e5ea] rounded-xl font-mono text-center font-bold"
              />
            </div>

            <label class="flex items-center gap-2 cursor-pointer mt-3">
              <input
                type="checkbox"
                bind:checked={formIsActive}
                class="rounded border-[#c4c7c5] text-[#0b57d0] focus:ring-[#0b57d0]"
              />
              <span class="font-medium text-[#1f1f1f]">Akun Aktif Digunakan</span>
            </label>
          </div>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="flex items-center justify-end gap-2.5 px-6 py-4 border-t border-[#e1e5ea] bg-[#f8fafd]">
        <button
          type="button"
          onclick={() => showModal = false}
          disabled={isSaving}
          class="px-4 py-2 text-xs font-medium text-[#444746] hover:bg-[#e1e5ea]/60 rounded-xl transition-all cursor-pointer"
        >
          Batal
        </button>
        <M3Button variant="filled" onclick={handleSaveAccount} disabled={isSaving}>
          {#if isSaving}
            <RefreshCw class="w-4 h-4 animate-spin" />
            <span>Menyimpan...</span>
          {:else}
            <Check class="w-4 h-4" />
            <span>{modalMode === 'create' ? 'Buat Akun' : 'Simpan Perubahan'}</span>
          {/if}
        </M3Button>
      </div>
    </div>
  </div>
{/if}

<!-- ========================================== -->
<!-- MODAL CONFIRM DELETE -->
<!-- ========================================== -->
{#if showDeleteConfirm && accountToDelete}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-xs p-4 animate-fadeIn">
    <div class="w-full max-w-md bg-white rounded-3xl shadow-xl border border-[#e1e5ea] overflow-hidden p-6 flex flex-col gap-4">
      <div class="flex items-center gap-3">
        <div class="p-3 rounded-2xl bg-rose-100 text-rose-700 shrink-0">
          <Trash2 class="w-6 h-6" />
        </div>
        <div>
          <h3 class="text-base font-bold text-[#1f1f1f]">Hapus Akun Buku Besar?</h3>
          <p class="text-xs text-[#747775]">Tindakan ini tidak dapat dibatalkan jika akun memiliki dependensi.</p>
        </div>
      </div>

      <div class="p-4 rounded-2xl bg-[#f8fafd] border border-[#e1e5ea] text-xs flex flex-col gap-1.5">
        <div class="flex justify-between">
          <span class="text-[#747775]">Kode Akun:</span>
          <span class="font-mono font-bold text-[#0b57d0]">{accountToDelete.code}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-[#747775]">Nama Akun:</span>
          <span class="font-semibold text-[#1f1f1f]">{accountToDelete.name}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-[#747775]">Klasifikasi:</span>
          <span class="font-medium">{accountToDelete.type} (Level {accountToDelete.account_level})</span>
        </div>
      </div>

      <p class="text-xs text-rose-700 leading-relaxed">
        <strong>Peringatan:</strong> Akun hanya dapat dihapus jika belum memiliki sub-akun anak dan belum pernah digunakan dalam jurnal transaksi buku besar keuangan.
      </p>

      <div class="flex items-center justify-end gap-2.5 pt-2">
        <button
          type="button"
          onclick={() => showDeleteConfirm = false}
          disabled={isDeleting}
          class="px-4 py-2 text-xs font-medium text-[#444746] hover:bg-[#e1e5ea]/60 rounded-xl transition-all cursor-pointer"
        >
          Batal
        </button>
        <button
          type="button"
          onclick={handleDeleteAccount}
          disabled={isDeleting}
          class="flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-rose-600 hover:bg-rose-700 rounded-xl transition-all cursor-pointer shadow-xs"
        >
          {#if isDeleting}
            <RefreshCw class="w-3.5 h-3.5 animate-spin" />
            <span>Menghapus...</span>
          {:else}
            <Trash2 class="w-3.5 h-3.5" />
            <span>Hapus Akun</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
