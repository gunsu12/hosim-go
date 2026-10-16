<script lang="ts">
  import { onMount } from "svelte";
  import {
    ArrowLeft,
    SlidersHorizontal,
    Search,
    Plus,
    RefreshCw,
    FolderTree,
    ChevronDown,
    ChevronRight,
    Calculator,
    CheckCircle2,
    AlertCircle,
    X,
    Filter,
    Edit2,
  } from "@lucide/svelte";
  import type {
    PriceSchemaRecord,
    PriceSchemaRuleRecord,
    PriceSchemaTreeNode,
    UpsertRuleDTO,
  } from "$lib/types/finance/price_schema";
  import type { DepartementRecord } from "$lib/types/master/departement";
  import type { ServiceUnitRecord } from "$lib/types/master/service_unit";
  import type { ItemSummaryRecord } from "$lib/types/master/item";
  import type { TariffClassRecord } from "$lib/types/finance/tariff";

  import {
    getPriceSchemaRules,
    upsertPriceSchemaRule,
    deletePriceSchemaRule,
  } from "$lib/api/finance/price_schema";
  import { getDepartments } from "$lib/api/organization/departments";
  import { getServiceUnits } from "$lib/api/organization/service_units";
  import { getItems } from "$lib/api/catalog/item";
  import { getTariffClasses } from "$lib/api/finance/tariff";

  import PriceSchemaDefaultCard from "./components/PriceSchemaDefaultCard.svelte";
  import PriceSchemaTreeView from "./components/PriceSchemaTreeView.svelte";
  import PriceSchemaRuleDrawer from "./components/PriceSchemaRuleDrawer.svelte";
  import PriceSchemaSimulator from "./components/PriceSchemaSimulator.svelte";

  import {
    buildPriceSchemaTree,
    filterTreeRecursive,
    collectExpandedNodeIds,
  } from "./treeUtils";

  interface Props {
    schema: PriceSchemaRecord;
    onBack: () => void;
    onEditHeader: (schema: PriceSchemaRecord) => void;
  }

  let { schema, onBack, onEditHeader }: Props = $props();

  // Active View Tab: 'TREE' or 'SIMULATOR'
  let activeTab = $state<"TREE" | "SIMULATOR">("TREE");

  // Search and Filter State
  let searchQuery = $state("");
  let filterMode = $state<"ALL" | "OVERRIDE" | "INHERITED">("ALL");

  // Data State
  let rules = $state<PriceSchemaRuleRecord[]>([]);
  let departments = $state<DepartementRecord[]>([]);
  let serviceUnits = $state<ServiceUnitRecord[]>([]);
  let items = $state<ItemSummaryRecord[]>([]);
  let tariffClasses = $state<TariffClassRecord[]>([]);

  let isLoading = $state(false);
  let isSavingRule = $state(false);
  let feedbackMessage = $state<{ text: string; isError?: boolean } | null>(
    null,
  );

  // Tree Nodes & Expanded State
  let expandedNodeIds = $state<Set<string>>(new Set());

  // Rule Drawer State
  let isDrawerOpen = $state(false);
  let targetNodeForDrawer = $state<PriceSchemaTreeNode | null>(null);

  function showToast(text: string, isError = false) {
    feedbackMessage = { text, isError };
    setTimeout(() => {
      feedbackMessage = null;
    }, 4500);
  }

  // Find SCHEMA_DEFAULT rule
  const defaultRule = $derived(
    rules.find((r) => r.level === "SCHEMA_DEFAULT") || null,
  );

  // Build items map
  const itemsMap = $derived.by(() => {
    const map = new Map<string, ItemSummaryRecord>();
    for (const it of items) {
      map.set(it.id, it);
    }
    return map;
  });

  // Build raw tree from data
  const rawTree = $derived(
    buildPriceSchemaTree(
      defaultRule,
      rules,
      departments,
      serviceUnits,
      itemsMap,
    ),
  );

  // Filtered tree according to search and mode
  const filteredTree = $derived(
    filterTreeRecursive(rawTree, searchQuery, filterMode),
  );

  async function loadAllData() {
    isLoading = true;
    try {
      const [rulesRes, deptsRes, unitsRes, itemsRes, classesRes] =
        await Promise.all([
          getPriceSchemaRules(schema.id, { limit: 1000 }),
          getDepartments({ limit: 100 }),
          getServiceUnits({ limit: 500 }),
          getItems({ limit: 1000 }),
          getTariffClasses({ limit: 100 }),
        ]);

      rules = rulesRes.items || [];
      departments = deptsRes.data || [];
      serviceUnits = unitsRes.data || [];
      items = itemsRes.data || [];
      tariffClasses = classesRes.data || [];

      // Auto-expand first 2 levels on first load
      if (expandedNodeIds.size === 0) {
        const initialTree = buildPriceSchemaTree(
          rules.find((r) => r.level === "SCHEMA_DEFAULT") || null,
          rules,
          departments,
          serviceUnits,
          itemsMap,
        );
        expandedNodeIds = collectExpandedNodeIds(initialTree, 2);
      }
    } catch (err: any) {
      showToast(err.message || "Gagal memuat aturan skema tarif", true);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadAllData();
  });

  function toggleNode(nodeId: string) {
    const nextSet = new Set(expandedNodeIds);
    if (nextSet.has(nodeId)) {
      nextSet.delete(nodeId);
    } else {
      nextSet.add(nodeId);
    }
    expandedNodeIds = nextSet;
  }

  function expandAll() {
    expandedNodeIds = collectExpandedNodeIds(rawTree, 10);
  }

  function collapseAll() {
    expandedNodeIds = new Set();
  }

  // Drawer Open Handlers
  function handleOpenDefaultRule() {
    targetNodeForDrawer = {
      id: defaultRule?.id || "synthetic_default",
      nodeType: "DEFAULT",
      scopeType: "SCHEMA_DEFAULT",
      name: "Aturan Dasar Skema (SCHEMA_DEFAULT)",
      isInherited: false,
      ruleId: defaultRule?.id,
      tariffAdjustmentType: defaultRule?.tariff_adjustment_type || "NONE",
      tariffAdjustmentValue: Number(defaultRule?.tariff_adjustment_value) || 0,
      medicationAdjustmentType:
        defaultRule?.medication_adjustment_type || "NONE",
      medicationAdjustmentValue:
        Number(defaultRule?.medication_adjustment_value) || 0,
      payByPatientPercent: Number(defaultRule?.pay_by_patient_percent) || 0,
      payByCustomerPercent: Number(defaultRule?.pay_by_customer_percent) || 100,
      taxPercent: Number(defaultRule?.tax_percent) || 11,
    };
    isDrawerOpen = true;
  }

  function handleCreateOverride(node: PriceSchemaTreeNode) {
    targetNodeForDrawer = node;
    isDrawerOpen = true;
  }

  function handleEditRule(node: PriceSchemaTreeNode) {
    targetNodeForDrawer = node;
    isDrawerOpen = true;
  }

  function handleAddSubItem(node: PriceSchemaTreeNode) {
    targetNodeForDrawer = node;
    isDrawerOpen = true;
  }

  async function handleResetRule(node: PriceSchemaTreeNode) {
    if (!node.ruleId) return;
    if (
      !confirm(
        `Hapus override untuk ${node.name}? Layanan ini akan kembali mewarisi aturan induk di atasnya.`,
      )
    ) {
      return;
    }

    try {
      await deletePriceSchemaRule(schema.id, node.ruleId);
      showToast(`Override ${node.name} berhasil dihapus.`);
      await loadAllData();
    } catch (err: any) {
      showToast(err.message || "Gagal menghapus override", true);
    }
  }

  async function handleSaveDrawerRule(payload: UpsertRuleDTO) {
    isSavingRule = true;
    try {
      await upsertPriceSchemaRule(schema.id, payload);
      isDrawerOpen = false;
      targetNodeForDrawer = null;
      showToast("Aturan penyesuaian tarif berhasil disimpan.");
      await loadAllData();
    } catch (err: any) {
      showToast(err.message || "Gagal menyimpan aturan tarif", true);
    } finally {
      isSavingRule = false;
    }
  }
</script>

<div class="flex flex-col gap-5">
  <!-- Toast Feedback Alerts -->
  {#if feedbackMessage}
    <div
      class="p-3.5 rounded-2xl {feedbackMessage.isError
        ? 'bg-rose-50 text-rose-700 border-rose-200 dark:bg-rose-950/40 dark:text-rose-300 dark:border-rose-800/60'
        : 'bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-800/60'} border flex items-center justify-between shadow-xs animate-in fade-in duration-200"
    >
      <div class="flex items-center gap-2 text-xs font-semibold">
        {#if feedbackMessage.isError}
          <AlertCircle class="w-4 h-4 shrink-0" />
        {:else}
          <CheckCircle2 class="w-4 h-4 shrink-0" />
        {/if}
        <span>{feedbackMessage.text}</span>
      </div>
      <button
        onclick={() => (feedbackMessage = null)}
        type="button"
        class="hover:opacity-75 cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  <!-- Breadcrumb & Top Bar -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div class="flex flex-col gap-1">
      <!-- Breadcrumb -->
      <nav
        class="flex items-center gap-2 text-xs text-on-surface-variant font-medium"
      >
        <button
          type="button"
          onclick={onBack}
          class="hover:text-primary transition-colors cursor-pointer flex items-center gap-1"
        >
          <ArrowLeft class="w-3.5 h-3.5" />
          <span>Skema Tarif Rekanan</span>
        </button>
        <span>/</span>
        <span class="font-bold text-on-surface font-mono">{schema.code}</span>
      </nav>

      <!-- Title & Schema Info -->
      <div class="flex items-center gap-2.5 mt-1 flex-wrap">
        <h2 class="text-lg font-bold text-on-surface tracking-tight">
          {schema.name}
        </h2>
        {#if schema.is_active}
          <span
            class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300"
          >
            Aktif
          </span>
        {:else}
          <span
            class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-bold bg-zinc-200 text-zinc-700 dark:bg-zinc-800 dark:text-zinc-400"
          >
            Nonaktif
          </span>
        {/if}

        {#if schema.base_price_plan_code}
          <span
            class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-lg text-xs font-semibold bg-purple-50 text-purple-700 dark:bg-purple-900/50 dark:text-purple-300 border border-purple-200 dark:border-purple-800"
          >
            <span class="text-[10px] text-purple-500 font-medium"
              >Buku Tarif:</span
            >
            <span class="font-mono">[{schema.base_price_plan_code}]</span>
            <span>{schema.base_price_plan_name || ""}</span>
          </span>
        {:else}
          <span
            class="inline-flex items-center gap-1 px-2 py-0.5 rounded-lg text-[11px] font-medium bg-surface-container text-on-surface-variant border border-outline-variant/60"
          >
            <span>Buku Tarif: Standar RS Terkini (Otomatis)</span>
          </span>
        {/if}
      </div>

      {#if schema.description}
        <p class="text-xs text-on-surface-variant max-w-3xl leading-relaxed">
          {schema.description}
        </p>
      {/if}
    </div>

    <!-- Right Header Actions -->
    <div class="flex items-center gap-2">
      <button
        type="button"
        onclick={loadAllData}
        class="p-2 rounded-xl text-on-surface-variant hover:text-on-surface hover:bg-surface-container border border-outline-variant transition-colors cursor-pointer"
        title="Muat ulang aturan"
        disabled={isLoading}
      >
        <RefreshCw class="w-3.5 h-3.5 {isLoading ? 'animate-spin' : ''}" />
      </button>
    </div>
  </div>

  <!-- Navigation Switcher: Pohon Aturan vs Simulator -->
  <div class="flex items-center gap-2 border-b border-outline-variant/60 pb-1">
    <button
      type="button"
      onclick={() => (activeTab = "TREE")}
      class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
      'TREE'
        ? 'bg-primary text-on-primary shadow-xs'
        : 'text-on-surface-variant hover:bg-surface-container'}"
    >
      <FolderTree class="w-3.5 h-3.5" />
      <span>Kebijakan Skema</span>
    </button>

    <button
      type="button"
      onclick={() => (activeTab = "SIMULATOR")}
      class="px-4 py-2 text-xs font-semibold rounded-xl transition-all flex items-center gap-2 cursor-pointer {activeTab ===
      'SIMULATOR'
        ? 'bg-purple-600 text-white shadow-xs'
        : 'text-on-surface-variant hover:bg-surface-container'}"
    >
      <Calculator class="w-3.5 h-3.5" />
      <span>Simulator & Uji Lookup Tarif</span>
    </button>
  </div>

  {#if activeTab === "TREE"}
    <!-- 1. Default Fallback Anchor Card -->
    <PriceSchemaDefaultCard
      {defaultRule}
      onEditDefault={handleOpenDefaultRule}
    />

    <!-- 2. Tree Toolbar -->
    <div class="flex flex-wrap items-center justify-between gap-3 text-xs">
      <!-- Search & Filters -->
      <div class="flex flex-wrap items-center gap-2 flex-1 max-w-lg">
        <div class="relative flex-1 min-w-[200px]">
          <Search
            class="w-3.5 h-3.5 absolute left-3 top-2.5 text-on-surface-variant"
          />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Cari departemen, unit poli, atau tindakan..."
            class="w-full pl-8 pr-3 py-2 rounded-xl border border-outline-variant bg-surface text-xs focus:outline-hidden focus:border-primary"
          />
        </div>

        <select
          bind:value={filterMode}
          class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs"
        >
          <option value="ALL">Semua Layanan</option>
          <option value="OVERRIDE">Hanya Override (Aturan Khusus)</option>
          <option value="INHERITED">Hanya Inherited (Mewarisi)</option>
        </select>
      </div>

      <!-- Expand / Collapse Controls -->
      <div class="flex items-center gap-1.5">
        <button
          type="button"
          onclick={expandAll}
          class="px-3 py-1.5 rounded-xl border border-outline-variant/60 bg-surface hover:bg-surface-container text-on-surface text-xs font-semibold transition-colors cursor-pointer"
        >
          Buka Semua
        </button>
        <button
          type="button"
          onclick={collapseAll}
          class="px-3 py-1.5 rounded-xl border border-outline-variant/60 bg-surface hover:bg-surface-container text-on-surface text-xs font-semibold transition-colors cursor-pointer"
        >
          Tutup Semua
        </button>
      </div>
    </div>

    <!-- 3. Hierarchical Tree Grid -->
    {#if isLoading}
      <div
        class="p-16 text-center text-on-surface-variant bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60"
      >
        <div class="inline-flex items-center gap-2">
          <span class="loading loading-spinner loading-sm text-primary"></span>
          <span class="text-xs">Menyusun struktur hierarkis skema tarif...</span
          >
        </div>
      </div>
    {:else}
      <PriceSchemaTreeView
        treeNodes={filteredTree}
        {expandedNodeIds}
        onToggleNode={toggleNode}
        onCreateOverride={handleCreateOverride}
        onEditRule={handleEditRule}
        onResetRule={handleResetRule}
        onAddSubItem={handleAddSubItem}
      />
    {/if}
  {:else}
    <!-- Tab 2: Simulator -->
    <PriceSchemaSimulator
      priceSchemaId={schema.id}
      basePricePlanId={schema.base_price_plan_id}
      {departments}
      {serviceUnits}
      {tariffClasses}
      {items}
    />
  {/if}

  <!-- Offcanvas Rule Drawer for Edit / Override -->
  <PriceSchemaRuleDrawer
    isOpen={isDrawerOpen}
    targetNode={targetNodeForDrawer}
    isSaving={isSavingRule}
    availableItems={items}
    onClose={() => {
      isDrawerOpen = false;
      targetNodeForDrawer = null;
    }}
    onSave={handleSaveDrawerRule}
  />
</div>
