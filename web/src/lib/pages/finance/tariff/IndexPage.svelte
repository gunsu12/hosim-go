<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Receipt,
    Layers,
    Tag,
    RefreshCw,
    CheckCircle2,
    AlertCircle,
    X,
    BookOpen,
    Calculator,
    Sliders,
    Building2,
    UserCheck,
    CheckCheck,
    Archive
  } from '@lucide/svelte';
  import {
    getTariffClasses,
    createTariffClass,
    updateTariffClass,
    deleteTariffClass,
    getTariffComponents,
    createTariffComponent,
    updateTariffComponent,
    deleteTariffComponent,
    getPricePlans,
    createPricePlan,
    updatePricePlan,
    submitPricePlan,
    approvePricePlan,
    activatePricePlan,
    archivePricePlan,
    clonePricePlan
  } from '$lib/api/finance';
  import type {
    TariffClassRecord,
    TariffComponentRecord,
    CreateTariffClassDTO,
    UpdateTariffClassDTO,
    CreateTariffComponentDTO,
    UpdateTariffComponentDTO,
    PricePlanRecord,
    CreatePricePlanDTO,
    UpdatePricePlanDTO,
    ClonePricePlanDTO
  } from '$lib/types/finance';

  // Subcomponents
  import PricePlanTable from './components/PricePlanTable.svelte';
  import PricePlanModal from './components/PricePlanModal.svelte';
  import PricePlanCloneModal from './components/PricePlanCloneModal.svelte';
  import PricePlanItemsModal from './components/PricePlanItemsModal.svelte';
  import TariffLookupSimulator from './components/TariffLookupSimulator.svelte';
  import TariffMatrixEditor from './components/TariffMatrixEditor.svelte';

  import TariffClassTable from './components/TariffClassTable.svelte';
  import TariffComponentTable from './components/TariffComponentTable.svelte';
  import TariffClassModal from './components/TariffClassModal.svelte';
  import TariffComponentModal from './components/TariffComponentModal.svelte';
  import TariffDeleteModal from './components/TariffDeleteModal.svelte';

  interface Props {
    initialTab?: 'PLAN' | 'MATRIX' | 'SIMULATOR' | 'CLASS' | 'COMPONENT';
  }

  let { initialTab = 'PLAN' }: Props = $props();

  // Active Tab State
  let activeTab = $state<'PLAN' | 'MATRIX' | 'SIMULATOR' | 'CLASS' | 'COMPONENT'>('PLAN');

  $effect(() => {
    if (initialTab) {
      activeTab = initialTab;
    }
  });

  // Toast / Feedback State
  let successMessage = $state<string | null>(null);
  let errorMessage = $state<string | null>(null);

  function showToast(msg: string, isError = false) {
    if (isError) {
      errorMessage = msg;
      setTimeout(() => (errorMessage = null), 6000);
    } else {
      successMessage = msg;
      setTimeout(() => (successMessage = null), 4000);
    }
  }

  // ==========================================
  // STATE: PRICE PLANS
  // ==========================================
  let plans = $state<PricePlanRecord[]>([]);
  let isPlansLoading = $state(false);
  let planSearchQuery = $state('');
  let planFilterStatus = $state<string>('ALL');
  let planCurrentPage = $state(1);
  let planTotalPages = $state(1);
  let planTotalCount = $state(0);

  // Plan Modal State
  let showPlanModal = $state(false);
  let planModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let selectedPlan = $state<PricePlanRecord | null>(null);
  let isPlanSaving = $state(false);

  // Clone Modal State
  let showCloneModal = $state(false);
  let sourceClonePlan = $state<PricePlanRecord | null>(null);
  let isPlanCloning = $state(false);

  // Items Management Modal State
  let showItemsModal = $state(false);
  let itemsTargetPlan = $state<PricePlanRecord | null>(null);

  // ==========================================
  // STATE: TARIFF CLASSES
  // ==========================================
  let classes = $state<TariffClassRecord[]>([]);
  let isClassesLoading = $state(false);
  let classSearchQuery = $state('');
  let classFilterStatus = $state<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL');
  let classCurrentPage = $state(1);
  let classTotalPages = $state(1);
  let classTotalCount = $state(0);

  // Class Modal State
  let showClassModal = $state(false);
  let classModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let selectedClass = $state<TariffClassRecord | null>(null);
  let isClassSaving = $state(false);

  // ==========================================
  // STATE: TARIFF COMPONENTS
  // ==========================================
  let components = $state<TariffComponentRecord[]>([]);
  let isComponentsLoading = $state(false);
  let compSearchQuery = $state('');
  let compFilterType = $state('ALL');
  let compFilterStatus = $state<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL');
  let compCurrentPage = $state(1);
  let compTotalPages = $state(1);
  let compTotalCount = $state(0);

  // Component Modal State
  let showCompModal = $state(false);
  let compModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let selectedComp = $state<TariffComponentRecord | null>(null);
  let isCompSaving = $state(false);

  // Delete Modal State (for Class & Component)
  let showDeleteModal = $state(false);
  let deleteItemType = $state<'CLASS' | 'COMPONENT'>('CLASS');
  let deleteItemId = $state('');
  let deleteItemCode = $state('');
  let deleteItemName = $state('');
  let isDeleting = $state(false);

  // ==========================================
  // STATS (DERIVED)
  // ==========================================
  let statsActivePlans = $derived((plans || []).filter((p) => p.status === 'ACTIVE').length);
  let statsDraftPlans = $derived((plans || []).filter((p) => p.status === 'DRAFT').length);
  let statsActiveClasses = $derived((classes || []).filter((c) => c.is_active).length);
  let statsHospitalComps = $derived((components || []).filter((c) => c.is_hospital_revenue).length);
  let statsOperatorComps = $derived((components || []).filter((c) => c.is_operator_revenue).length);

  // ==========================================
  // DATA FETCHING
  // ==========================================
  async function loadPlans() {
    isPlansLoading = true;
    try {
      const res = await getPricePlans({
        page: planCurrentPage,
        limit: 10,
        search: planSearchQuery || undefined,
        status: planFilterStatus !== 'ALL' ? (planFilterStatus as any) : undefined
      });
      plans = Array.isArray(res.data) ? res.data : [];
      planTotalCount = res.meta?.total_items ?? plans.length;
      planTotalPages = res.meta?.total_pages ?? 1;
    } catch (err: any) {
      plans = [];
      showToast(err.message || 'Gagal memuat daftar buku tarif', true);
    } finally {
      isPlansLoading = false;
    }
  }

  async function loadClasses() {
    isClassesLoading = true;
    try {
      const res = await getTariffClasses({
        page: classCurrentPage,
        limit: 10,
        search: classSearchQuery || undefined,
        is_active: classFilterStatus === 'ALL' ? undefined : classFilterStatus === 'ACTIVE'
      });
      classes = res.data;
      classTotalCount = res.meta.total;
      classTotalPages = res.meta.total_pages;
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat daftar kelas tarif', true);
    } finally {
      isClassesLoading = false;
    }
  }

  async function loadComponents() {
    isComponentsLoading = true;
    try {
      const res = await getTariffComponents({
        page: compCurrentPage,
        limit: 10,
        search: compSearchQuery || undefined,
        component_type: compFilterType !== 'ALL' ? compFilterType : undefined,
        is_active: compFilterStatus === 'ALL' ? undefined : compFilterStatus === 'ACTIVE'
      });
      components = res.data;
      compTotalCount = res.meta.total;
      compTotalPages = res.meta.total_pages;
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat daftar komponen tarif', true);
    } finally {
      isComponentsLoading = false;
    }
  }

  function handleRefresh() {
    if (activeTab === 'PLAN') {
      loadPlans();
    } else if (activeTab === 'CLASS') {
      loadClasses();
    } else if (activeTab === 'COMPONENT') {
      loadComponents();
    }
  }

  onMount(() => {
    loadPlans();
    loadClasses();
    loadComponents();
  });

  // Re-fetch when switching tabs
  $effect(() => {
    if (activeTab === 'PLAN') {
      loadPlans();
    } else if (activeTab === 'CLASS') {
      loadClasses();
    } else if (activeTab === 'COMPONENT') {
      loadComponents();
    }
  });

  // ==========================================
  // PRICE PLAN ACTIONS
  // ==========================================
  function openCreatePlan() {
    planModalMode = 'CREATE';
    selectedPlan = null;
    showPlanModal = true;
  }

  function openEditPlan(plan: PricePlanRecord) {
    planModalMode = 'EDIT';
    selectedPlan = plan;
    showPlanModal = true;
  }

  function openClonePlan(plan: PricePlanRecord) {
    sourceClonePlan = plan;
    showCloneModal = true;
  }

  function openManageItems(plan: PricePlanRecord) {
    itemsTargetPlan = plan;
    showItemsModal = true;
  }

  async function handleSavePlan(payload: CreatePricePlanDTO | UpdatePricePlanDTO) {
    isPlanSaving = true;
    try {
      if (planModalMode === 'CREATE') {
        await createPricePlan(payload as CreatePricePlanDTO);
        showToast(`Buku tarif ${payload.name} berhasil dibuat`);
      } else if (selectedPlan) {
        await updatePricePlan(selectedPlan.id, payload as UpdatePricePlanDTO);
        showToast(`Buku tarif ${payload.name} berhasil diperbarui`);
      }
      showPlanModal = false;
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan buku tarif', true);
    } finally {
      isPlanSaving = false;
    }
  }

  async function handleClonePlan(payload: ClonePricePlanDTO) {
    if (!sourceClonePlan) return;
    isPlanCloning = true;
    try {
      await clonePricePlan(sourceClonePlan.id, payload);
      showToast(`Buku tarif ${payload.name} berhasil diduplikasi ke draf baru`);
      showCloneModal = false;
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal menggandakan buku tarif', true);
    } finally {
      isPlanCloning = false;
    }
  }

  async function handleSubmitPlan(plan: PricePlanRecord) {
    if (!confirm(`Ajukan buku tarif "${plan.name}" untuk direview/approval?`)) return;
    try {
      await submitPricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} berhasil diajukan (SUBMITTED)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal mengajukan buku tarif', true);
    }
  }

  async function handleApprovePlan(plan: PricePlanRecord) {
    if (!confirm(`Setujui buku tarif "${plan.name}"? Pastikan total pecahan komponen seimbang.`)) return;
    try {
      await approvePricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} berhasil disetujui (APPROVED)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyetujui buku tarif', true);
    }
  }

  async function handleActivatePlan(plan: PricePlanRecord) {
    if (!confirm(`Aktifkan buku tarif "${plan.name}" sebagai acuan tarif transaksi pelayanan?`)) return;
    try {
      await activatePricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} sekarang AKTIF (ACTIVE)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal mengaktifkan buku tarif', true);
    }
  }

  async function handleArchivePlan(plan: PricePlanRecord) {
    if (!confirm(`Arsipkan buku tarif "${plan.name}"? Buku tarif ini tidak akan digunakan lagi untuk transaksi baru.`)) return;
    try {
      await archivePricePlan(plan.id);
      showToast(`Buku tarif ${plan.name} telah diarsipkan (ARCHIVED)`);
      loadPlans();
    } catch (err: any) {
      showToast(err.message || 'Gagal mengarsipkan buku tarif', true);
    }
  }

  // ==========================================
  // TARIFF CLASS ACTIONS
  // ==========================================
  function openCreateClass() {
    classModalMode = 'CREATE';
    selectedClass = null;
    showClassModal = true;
  }

  function openEditClass(tc: TariffClassRecord) {
    classModalMode = 'EDIT';
    selectedClass = tc;
    showClassModal = true;
  }

  function openDeleteClass(tc: TariffClassRecord) {
    deleteItemType = 'CLASS';
    deleteItemId = tc.id;
    deleteItemCode = tc.code;
    deleteItemName = tc.name;
    showDeleteModal = true;
  }

  async function handleSaveClass(payload: CreateTariffClassDTO | UpdateTariffClassDTO) {
    isClassSaving = true;
    try {
      if (classModalMode === 'CREATE') {
        await createTariffClass(payload as CreateTariffClassDTO);
        showToast(`Kelas tarif ${payload.name} berhasil ditambahkan`);
      } else if (selectedClass) {
        await updateTariffClass(selectedClass.id, payload as UpdateTariffClassDTO);
        showToast(`Kelas tarif ${payload.name} berhasil diperbarui`);
      }
      showClassModal = false;
      loadClasses();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan kelas tarif', true);
    } finally {
      isClassSaving = false;
    }
  }

  // ==========================================
  // TARIFF COMPONENT ACTIONS
  // ==========================================
  function openCreateComponent() {
    compModalMode = 'CREATE';
    selectedComp = null;
    showCompModal = true;
  }

  function openEditComponent(tc: TariffComponentRecord) {
    compModalMode = 'EDIT';
    selectedComp = tc;
    showCompModal = true;
  }

  function openDeleteComponent(tc: TariffComponentRecord) {
    deleteItemType = 'COMPONENT';
    deleteItemId = tc.id;
    deleteItemCode = tc.code;
    deleteItemName = tc.name;
    showDeleteModal = true;
  }

  async function handleSaveComponent(payload: CreateTariffComponentDTO | UpdateTariffComponentDTO) {
    isCompSaving = true;
    try {
      if (compModalMode === 'CREATE') {
        await createTariffComponent(payload as CreateTariffComponentDTO);
        showToast(`Komponen tarif ${payload.name} berhasil ditambahkan`);
      } else if (selectedComp) {
        await updateTariffComponent(selectedComp.id, payload as UpdateTariffComponentDTO);
        showToast(`Komponen tarif ${payload.name} berhasil diperbarui`);
      }
      showCompModal = false;
      loadComponents();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan komponen tarif', true);
    } finally {
      isCompSaving = false;
    }
  }

  // ==========================================
  // DELETE CONFIRMATION HANDLER
  // ==========================================
  async function handleConfirmDelete() {
    if (!deleteItemId) return;
    isDeleting = true;
    try {
      if (deleteItemType === 'CLASS') {
        await deleteTariffClass(deleteItemId);
        showToast(`Kelas tarif ${deleteItemName} berhasil dihapus`);
        loadClasses();
      } else {
        await deleteTariffComponent(deleteItemId);
        showToast(`Komponen tarif ${deleteItemName} berhasil dihapus`);
        loadComponents();
      }
      showDeleteModal = false;
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus data', true);
    } finally {
      isDeleting = false;
    }
  }
</script>

<div class="flex flex-col gap-5">
  <!-- Toast Feedback Alerts -->
  {#if successMessage}
    <div class="p-3.5 rounded-2xl bg-[#e6f4ea] border border-[#ceead6] text-[#137333] flex items-center justify-between shadow-xs animate-in fade-in duration-200">
      <div class="flex items-center gap-2 text-xs font-semibold">
        <CheckCircle2 class="w-4 h-4 text-[#137333]" />
        <span>{successMessage}</span>
      </div>
      <button onclick={() => (successMessage = null)} type="button" class="text-[#137333] hover:opacity-75">
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div class="p-3.5 rounded-2xl bg-[#fce8e6] border border-[#fad2cf] text-[#c5221f] flex items-center justify-between shadow-xs animate-in fade-in duration-200">
      <div class="flex items-center gap-2 text-xs font-semibold">
        <AlertCircle class="w-4 h-4 text-[#c5221f]" />
        <span>{errorMessage}</span>
      </div>
      <button onclick={() => (errorMessage = null)} type="button" class="text-[#c5221f] hover:opacity-75">
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  <!-- Header Domain -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <div class="flex items-center gap-2.5">
        <div class="w-9 h-9 rounded-2xl bg-emerald-600/10 text-emerald-700 flex items-center justify-center font-bold">
          <Receipt class="w-5 h-5" />
        </div>
        <div>
          <h2 class="text-lg font-bold text-[#1f1f1f] tracking-tight">
            Tarif & Keuangan Rumah Sakit (Finance)
          </h2>
          <div class="flex items-center gap-2 mt-0.5">
            <span class="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
              domain: internal/finance
            </span>
            <span class="text-xs text-[#747775]">• Buku Tarif (Price Plan), Kelas, Pecahan Komponen & Simulator</span>
          </div>
        </div>
      </div>
      <p class="text-xs text-[#444746] mt-1.5 max-w-3xl leading-relaxed">
        Pusat pengelolaan buku tarif pelayanan rumah sakit (standar RS & rekanan PKS), validasi keseimbangan pecahan komponen tagihan, siklus persetujuan bertahap, serta kalkulator lookup pencarian tarif real-time.
      </p>
    </div>

    <div class="flex items-center gap-2">
      <button
        type="button"
        onclick={handleRefresh}
        class="h-10 px-3.5 rounded-xl border border-[#e1e5ea] bg-white text-[#444746] hover:bg-[#f8fafd] text-xs font-medium flex items-center gap-1.5 transition-colors disabled:opacity-50 cursor-pointer shadow-xs"
        title="Muat ulang data"
      >
        <RefreshCw class="w-3.5 h-3.5 {(isPlansLoading || isClassesLoading || isComponentsLoading) ? 'animate-spin text-[#0b57d0]' : ''}" />
        <span>Refresh</span>
      </button>
    </div>
  </div>

  <!-- KPI / Stats Summary Cards -->
  <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
    <!-- Stat 1: Buku Tarif (Price Plans) -->
    <button
      type="button"
      onclick={() => (activeTab = 'PLAN')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab === 'PLAN' ? 'bg-indigo-50/80 border-indigo-300 ring-2 ring-indigo-400/20' : 'bg-white border-[#e1e5ea] hover:border-indigo-200'}"
    >
      <div class="flex items-center justify-between text-xs font-medium text-indigo-800 mb-1">
        <span class="flex items-center gap-1.5">
          <BookOpen class="w-4 h-4 text-indigo-600" />
          Buku Tarif
        </span>
        <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-indigo-100 text-indigo-800">PLAN</span>
      </div>
      <div class="text-lg font-bold text-[#1f1f1f]">{planTotalCount}</div>
      <div class="text-[11px] text-[#747775] mt-0.5">{statsActivePlans} Aktif • {statsDraftPlans} Draf</div>
    </button>

    <!-- Stat 2: Simulator & Lookup -->
    <button
      type="button"
      onclick={() => (activeTab = 'SIMULATOR')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab === 'SIMULATOR' ? 'bg-purple-50/80 border-purple-300 ring-2 ring-purple-400/20' : 'bg-white border-[#e1e5ea] hover:border-purple-200'}"
    >
      <div class="flex items-center justify-between text-xs font-medium text-purple-800 mb-1">
        <span class="flex items-center gap-1.5">
          <Calculator class="w-4 h-4 text-purple-600" />
          Lookup Engine
        </span>
        <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-purple-100 text-purple-800">SIM</span>
      </div>
      <div class="text-lg font-bold text-[#1f1f1f]">3-Step</div>
      <div class="text-[11px] text-[#747775] mt-0.5">Hierarki Rekanan & Standar RS</div>
    </button>

    <!-- Stat 3: Kelas Tarif -->
    <button
      type="button"
      onclick={() => (activeTab = 'CLASS')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab === 'CLASS' ? 'bg-blue-50/80 border-blue-300 ring-2 ring-blue-400/20' : 'bg-white border-[#e1e5ea] hover:border-blue-200'}"
    >
      <div class="flex items-center justify-between text-xs font-medium text-blue-800 mb-1">
        <span class="flex items-center gap-1.5">
          <Tag class="w-4 h-4 text-blue-600" />
          Kelas Tarif
        </span>
        <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-blue-100 text-blue-800">CLASS</span>
      </div>
      <div class="text-lg font-bold text-[#1f1f1f]">{classTotalCount}</div>
      <div class="text-[11px] text-[#747775] mt-0.5">{statsActiveClasses} Kelas Aktif Digunakan</div>
    </button>

    <!-- Stat 4: Total Komponen -->
    <button
      type="button"
      onclick={() => (activeTab = 'COMPONENT')}
      class="p-3.5 rounded-2xl border text-left transition-all cursor-pointer {activeTab === 'COMPONENT' ? 'bg-emerald-50/80 border-emerald-300 ring-2 ring-emerald-400/20' : 'bg-white border-[#e1e5ea] hover:border-emerald-200'}"
    >
      <div class="flex items-center justify-between text-xs font-medium text-emerald-800 mb-1">
        <span class="flex items-center gap-1.5">
          <Layers class="w-4 h-4 text-emerald-600" />
          Komponen Tarif
        </span>
        <span class="text-[10px] font-mono px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-800">COMP</span>
      </div>
      <div class="text-lg font-bold text-[#1f1f1f]">{compTotalCount}</div>
      <div class="text-[11px] text-[#747775] mt-0.5">{statsHospitalComps} RS • {statsOperatorComps} Dokter</div>
    </button>
  </div>

  <!-- Navigation Tabs -->
  <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-1 gap-2 overflow-x-auto">
    <div class="flex items-center gap-1.5">
      <!-- Tab 1: Price Plan -->
      <button
        type="button"
        onclick={() => (activeTab = 'PLAN')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab === 'PLAN' ? 'bg-indigo-600 text-white shadow-xs' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
      >
        <BookOpen class="w-3.5 h-3.5" />
        <span>Buku Tarif ({planTotalCount})</span>
      </button>

      <!-- Tab 2: Editor Matriks Multi-Kelas -->
      <button
        type="button"
        onclick={() => (activeTab = 'MATRIX')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab === 'MATRIX' ? 'bg-[#0b57d0] text-white shadow-xs' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
      >
        <Sliders class="w-3.5 h-3.5" />
        <span>Editor Matriks Tarif (Multi-Kelas)</span>
      </button>

      <!-- Tab 3: Lookup Simulator -->
      <button
        type="button"
        onclick={() => (activeTab = 'SIMULATOR')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab === 'SIMULATOR' ? 'bg-purple-600 text-white shadow-xs' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
      >
        <Calculator class="w-3.5 h-3.5" />
        <span>Simulator & Lookup Tarif</span>
      </button>

      <!-- Tab 4: Kelas Tarif -->
      <button
        type="button"
        onclick={() => (activeTab = 'CLASS')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab === 'CLASS' ? 'bg-blue-600 text-white shadow-xs' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
      >
        <Tag class="w-3.5 h-3.5" />
        <span>Kelas Perawatan ({classTotalCount})</span>
      </button>

      <!-- Tab 5: Komponen Tarif -->
      <button
        type="button"
        onclick={() => (activeTab = 'COMPONENT')}
        class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab === 'COMPONENT' ? 'bg-emerald-600 text-white shadow-xs' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
      >
        <Layers class="w-3.5 h-3.5" />
        <span>Komponen Tarif ({compTotalCount})</span>
      </button>
    </div>
  </div>

  <!-- TAB CONTENT RENDER -->
  {#if activeTab === 'PLAN'}
    <PricePlanTable
      plans={plans}
      isLoading={isPlansLoading}
      searchQuery={planSearchQuery}
      filterStatus={planFilterStatus}
      currentPage={planCurrentPage}
      totalPages={planTotalPages}
      totalCount={planTotalCount}
      onSearchChange={(q) => {
        planSearchQuery = q;
        planCurrentPage = 1;
        loadPlans();
      }}
      onStatusChange={(s) => {
        planFilterStatus = s;
        planCurrentPage = 1;
        loadPlans();
      }}
      onPageChange={(p) => {
        planCurrentPage = p;
        loadPlans();
      }}
      onAdd={openCreatePlan}
      onEdit={openEditPlan}
      onClone={openClonePlan}
      onManageItems={openManageItems}
      onSubmitPlan={handleSubmitPlan}
      onApprovePlan={handleApprovePlan}
      onActivatePlan={handleActivatePlan}
      onArchivePlan={handleArchivePlan}
    />
  {:else if activeTab === 'MATRIX'}
    <TariffMatrixEditor showToast={showToast} />
  {:else if activeTab === 'SIMULATOR'}
    <TariffLookupSimulator />
  {:else if activeTab === 'CLASS'}
    <TariffClassTable
      classes={classes}
      isLoading={isClassesLoading}
      searchQuery={classSearchQuery}
      filterStatus={classFilterStatus}
      currentPage={classCurrentPage}
      totalPages={classTotalPages}
      totalCount={classTotalCount}
      onSearchChange={(q) => {
        classSearchQuery = q;
        classCurrentPage = 1;
        loadClasses();
      }}
      onStatusChange={(s) => {
        classFilterStatus = s;
        classCurrentPage = 1;
        loadClasses();
      }}
      onPageChange={(p) => {
        classCurrentPage = p;
        loadClasses();
      }}
      onAdd={openCreateClass}
      onEdit={openEditClass}
      onDelete={openDeleteClass}
    />
  {:else}
    <TariffComponentTable
      components={components}
      isLoading={isComponentsLoading}
      searchQuery={compSearchQuery}
      filterType={compFilterType}
      filterStatus={compFilterStatus}
      currentPage={compCurrentPage}
      totalPages={compTotalPages}
      totalCount={compTotalCount}
      onSearchChange={(q) => {
        compSearchQuery = q;
        compCurrentPage = 1;
        loadComponents();
      }}
      onTypeChange={(t) => {
        compFilterType = t;
        compCurrentPage = 1;
        loadComponents();
      }}
      onStatusChange={(s) => {
        compFilterStatus = s;
        compCurrentPage = 1;
        loadComponents();
      }}
      onPageChange={(p) => {
        compCurrentPage = p;
        loadComponents();
      }}
      onAdd={openCreateComponent}
      onEdit={openEditComponent}
      onDelete={openDeleteComponent}
    />
  {/if}
</div>

<!-- Modal Create / Edit Buku Tarif -->
<PricePlanModal
  isOpen={showPlanModal}
  mode={planModalMode}
  plan={selectedPlan}
  isSaving={isPlanSaving}
  onClose={() => (showPlanModal = false)}
  onSave={handleSavePlan}
/>

<!-- Modal Kloning Buku Tarif -->
<PricePlanCloneModal
  isOpen={showCloneModal}
  sourcePlan={sourceClonePlan}
  isCloning={isPlanCloning}
  onClose={() => (showCloneModal = false)}
  onClone={handleClonePlan}
/>

<!-- Drawer / Modal Manajemen Item & Komponen Buku Tarif -->
<PricePlanItemsModal
  isOpen={showItemsModal}
  plan={itemsTargetPlan}
  onClose={() => {
    showItemsModal = false;
    loadPlans();
  }}
  showToast={showToast}
/>

<!-- Modal Kelas Tarif -->
<TariffClassModal
  isOpen={showClassModal}
  mode={classModalMode}
  tariffClass={selectedClass}
  isSaving={isClassSaving}
  onClose={() => (showClassModal = false)}
  onSave={handleSaveClass}
/>

<!-- Modal Komponen Tarif -->
<TariffComponentModal
  isOpen={showCompModal}
  mode={compModalMode}
  tariffComponent={selectedComp}
  isSaving={isCompSaving}
  onClose={() => (showCompModal = false)}
  onSave={handleSaveComponent}
/>

<!-- Modal Hapus Data (Kelas / Komponen) -->
<TariffDeleteModal
  isOpen={showDeleteModal}
  title={deleteItemType === 'CLASS' ? 'Hapus Kelas Tarif' : 'Hapus Komponen Tarif'}
  itemCode={deleteItemCode}
  itemName={deleteItemName}
  isDeleting={isDeleting}
  onClose={() => (showDeleteModal = false)}
  onConfirm={handleConfirmDelete}
/>
