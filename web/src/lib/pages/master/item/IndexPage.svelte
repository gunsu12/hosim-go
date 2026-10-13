<script lang="ts">
  import { onMount } from "svelte";
  import {
    Search,
    Plus,
    RefreshCw,
    Pill,
    Package,
    Stethoscope,
    Receipt,
    Settings,
    Boxes,
    CheckCircle2,
    AlertCircle,
    X,
  } from "@lucide/svelte";
  import M3Button from "../../../components/m3/M3Button.svelte";
  import M3Select from "../../../components/m3/M3Select.svelte";
  import {
    getItems,
    deleteItem,
    getItemCategories,
    getItemProductLines,
  } from "../../../api/catalog";
  import type {
    ItemSummaryRecord,
    ItemCategoryRecord,
    ItemProductLineRecord,
    ItemType,
  } from "../../../types/master/item";

  // Sub-components
  import ItemTable from "./components/ItemTable.svelte";
  import ItemDetailModal from "./components/ItemDetailModal.svelte";
  import ItemFormModal from "./components/ItemFormModal.svelte";
  import ItemDeleteModal from "./components/ItemDeleteModal.svelte";
  import ItemConfigTab from "./components/ItemConfigTab.svelte";

  // ==========================================
  // STATE MANAGEMENT
  // ==========================================

  let activeTab = $state<ItemType | "ALL" | "CONFIG">("ALL");
  let searchQuery = $state("");
  let filterCategoryId = $state<string>("");
  let filterProductLineId = $state<string>("");
  let filterStatus = $state<"ALL" | "ACTIVE" | "INACTIVE">("ALL");

  let items = $state<ItemSummaryRecord[]>([]);
  let categories = $state<ItemCategoryRecord[]>([]);
  let productLines = $state<ItemProductLineRecord[]>([]);

  // Pagination
  let currentPage = $state(1);
  let pageSize = $state(15);
  let totalItems = $state(0);
  let totalPages = $state(1);

  // Loading & Feedback
  let isLoading = $state(false);
  let isActionLoading = $state(false);
  let successMessage = $state<string | null>(null);
  let errorMessage = $state<string | null>(null);

  // Modal States
  let showDetailModal = $state(false);
  let selectedDetailItem = $state<ItemSummaryRecord | null>(null);

  let showFormModal = $state(false);
  let formMode = $state<"CREATE" | "EDIT">("CREATE");
  let formTargetType = $state<ItemType>("MEDICATION");
  let selectedEditItem = $state<ItemSummaryRecord | null>(null);

  let showDeleteModal = $state(false);
  let itemToDelete = $state<ItemSummaryRecord | null>(null);

  // Filter options
  let categoryFilterOptions = $derived([
    { value: "", label: "Semua Kategori" },
    ...categories.map((c) => ({
      value: c.id,
      label: `[${c.item_type}] ${c.name}`,
    })),
  ]);

  let productLineFilterOptions = $derived([
    { value: "", label: "Semua Lini Produk" },
    ...productLines.map((pl) => ({
      value: pl.id,
      label: pl.name,
    })),
  ]);

  const statusFilterOptions = [
    { value: "ALL", label: "Semua Status" },
    { value: "ACTIVE", label: "Aktif Saja" },
    { value: "INACTIVE", label: "Nonaktif" },
  ];

  // ==========================================
  // DERIVED & HELPERS
  // ==========================================

  function showToast(msg: string, isError = false) {
    if (isError) {
      errorMessage = msg;
      setTimeout(() => (errorMessage = null), 5000);
    } else {
      successMessage = msg;
      setTimeout(() => (successMessage = null), 4000);
    }
  }

  // ==========================================
  // DATA FETCHING
  // ==========================================

  async function loadItems() {
    isLoading = true;
    try {
      const typeParam =
        activeTab !== "ALL" && activeTab !== "CONFIG" ? activeTab : undefined;
      const res = await getItems({
        page: currentPage,
        limit: pageSize,
        search: searchQuery || undefined,
        item_type: typeParam,
        category_id: filterCategoryId || undefined,
        product_line_id: filterProductLineId || undefined,
        is_active:
          filterStatus === "ALL" ? undefined : filterStatus === "ACTIVE",
      });

      items = res.data;
      totalItems = res.meta.total_items;
      totalPages = res.meta.total_pages;
    } catch (err: any) {
      showToast(err.message || "Gagal memuat katalog item", true);
    } finally {
      isLoading = false;
    }
  }

  async function loadReferences() {
    try {
      const [catRes, plRes] = await Promise.all([
        getItemCategories({ limit: 100 }),
        getItemProductLines({ limit: 100 }),
      ]);
      categories = catRes.data;
      productLines = plRes.data;
    } catch (err: any) {
      console.error("Gagal memuat data referensi kategori & lini produk", err);
    }
  }

  onMount(() => {
    loadReferences();
    loadItems();
  });

  // Re-fetch when tab changes
  $effect(() => {
    if (activeTab) {
      currentPage = 1;
      if (activeTab !== "CONFIG") {
        loadItems();
      }
    }
  });

  // ==========================================
  // MODAL TRIGGER HANDLERS
  // ==========================================

  function handleOpenDetail(item: ItemSummaryRecord) {
    selectedDetailItem = item;
    showDetailModal = true;
  }

  function handleOpenCreate(type?: ItemType) {
    formMode = "CREATE";
    selectedEditItem = null;
    formTargetType =
      type ||
      (activeTab !== "ALL" && activeTab !== "CONFIG"
        ? activeTab
        : "MEDICATION");
    showFormModal = true;
  }

  function handleOpenEdit(item: ItemSummaryRecord) {
    formMode = "EDIT";
    selectedEditItem = item;
    formTargetType = item.item_type;
    showFormModal = true;
  }

  function handleOpenDelete(item: ItemSummaryRecord) {
    itemToDelete = item;
    showDeleteModal = true;
  }

  async function handleConfirmDelete(item: ItemSummaryRecord) {
    isActionLoading = true;
    try {
      await deleteItem(item.id);
      showToast(`Item ${item.name} (${item.code}) berhasil dihapus`);
      showDeleteModal = false;
      itemToDelete = null;
      loadItems();
    } catch (err: any) {
      showToast(err.message || "Gagal menghapus item katalog", true);
    } finally {
      isActionLoading = false;
    }
  }
</script>

<div class="flex flex-col gap-5">
  <!-- Toast Feedback Alerts -->
  {#if successMessage}
    <div
      class="p-3.5 rounded-2xl bg-[#e6f4ea] dark:bg-emerald-950/40 border border-[#ceead6] dark:border-emerald-800 text-[#137333] dark:text-emerald-300 flex items-center justify-between shadow-xs animate-in fade-in duration-200"
    >
      <div class="flex items-center gap-2 text-xs font-semibold">
        <CheckCircle2 class="w-4 h-4 text-[#137333] dark:text-emerald-400" />
        <span>{successMessage}</span>
      </div>
      <button
        onclick={() => (successMessage = null)}
        type="button"
        class="text-[#137333] dark:text-emerald-400 hover:opacity-75 cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  {#if errorMessage}
    <div
      class="p-3.5 rounded-2xl bg-[#fce8e6] dark:bg-rose-950/40 border border-[#fad2cf] dark:border-rose-800 text-[#c5221f] dark:text-rose-300 flex items-center justify-between shadow-xs animate-in fade-in duration-200"
    >
      <div class="flex items-center gap-2 text-xs font-semibold">
        <AlertCircle class="w-4 h-4 text-[#c5221f] dark:text-rose-400" />
        <span>{errorMessage}</span>
      </div>
      <button
        onclick={() => (errorMessage = null)}
        type="button"
        class="text-[#c5221f] dark:text-rose-400 hover:opacity-75 cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  <!-- Header Domain -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <div class="flex items-center gap-2.5">
        <div
          class="w-9 h-9 rounded-2xl bg-[#0b57d0]/10 dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-blue-400 flex items-center justify-center font-bold"
        >
          <Boxes class="w-5 h-5" />
        </div>
        <div>
          <h2 class="text-lg font-bold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">
            Master Katalog Item & Tarif
          </h2>
        </div>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-1.5 max-w-3xl">
        Katalog universal sentral untuk seluruh obat farmasi, BMHP habis pakai,
        aset modal alkes, serta tarif jasa medis RS dengan integrasi
        multi-satuan (UOM) dan pemetaan bagan akun (COA).
      </p>
    </div>

    <!-- Top Action Buttons -->
    <div class="flex items-center gap-2 flex-wrap">
      <button
        onclick={loadItems}
        disabled={isLoading}
        type="button"
        class="h-10 px-3.5 rounded-xl border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#1e1f20] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f8fafd] dark:hover:bg-[#282a2c] text-xs font-medium flex items-center gap-1.5 transition-colors disabled:opacity-50 cursor-pointer shadow-xs"
        title="Muat ulang data"
      >
        <RefreshCw
          class="w-3.5 h-3.5 {isLoading ? 'animate-spin text-[#0b57d0] dark:text-blue-400' : ''}"
        />
        <span>Refresh</span>
      </button>

      <M3Button variant="tonal" onclick={() => handleOpenCreate("MEDICATION")}>
        <Plus class="w-4 h-4" />
        <span>Obat</span>
      </M3Button>

      <M3Button variant="filled" onclick={() => handleOpenCreate()}>
        <Plus class="w-4 h-4" />
        <span>Tambah Item</span>
      </M3Button>
    </div>
  </div>

  <!-- Navigation Category Tabs -->
  <div
    class="flex items-center justify-between border-b border-[#e1e5ea] dark:border-[#444746] pb-1 gap-2 overflow-x-auto"
  >
    <div class="flex items-center gap-1.5">
      <button
        type="button"
        onclick={() => (activeTab = "ALL")}
        class="px-4 py-2 text-xs font-medium rounded-xl transition-all cursor-pointer {activeTab ===
        "ALL"
          ? 'bg-[#0b57d0] text-white shadow-xs'
          : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
      >
        Semua Item ({totalItems})
      </button>

      <button
        type="button"
        onclick={() => (activeTab = "MEDICATION")}
        class="px-4 py-2 text-xs font-medium rounded-xl transition-all flex items-center gap-1.5 cursor-pointer {activeTab ===
        "MEDICATION"
          ? 'bg-emerald-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
      >
        <Pill class="w-3.5 h-3.5" />
        <span>Obat-Obatan</span>
      </button>

      <button
        type="button"
        onclick={() => (activeTab = "GENERAL")}
        class="px-4 py-2 text-xs font-medium rounded-xl transition-all flex items-center gap-1.5 cursor-pointer {activeTab ===
        "GENERAL"
          ? 'bg-amber-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
      >
        <Package class="w-3.5 h-3.5" />
        <span>BMHP & Umum</span>
      </button>

      <button
        type="button"
        onclick={() => (activeTab = "ASSET")}
        class="px-4 py-2 text-xs font-medium rounded-xl transition-all flex items-center gap-1.5 cursor-pointer {activeTab ===
        "ASSET"
          ? 'bg-indigo-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
      >
        <Stethoscope class="w-3.5 h-3.5" />
        <span>Aset & Alkes</span>
      </button>

      <button
        type="button"
        onclick={() => (activeTab = "TARIFF")}
        class="px-4 py-2 text-xs font-medium rounded-xl transition-all flex items-center gap-1.5 cursor-pointer {activeTab ===
        "TARIFF"
          ? 'bg-blue-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
      >
        <Receipt class="w-3.5 h-3.5" />
        <span>Tarif Layanan</span>
      </button>

      <button
        type="button"
        onclick={() => (activeTab = "CONFIG")}
        class="px-4 py-2 text-xs font-medium rounded-xl transition-all flex items-center gap-1.5 cursor-pointer {activeTab ===
        "CONFIG"
          ? 'bg-purple-600 text-white shadow-xs'
          : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
      >
        <Settings class="w-3.5 h-3.5" />
        <span>Kategori & Lini Produk</span>
      </button>
    </div>

    {#if activeTab !== "CONFIG"}
      <div class="text-xs text-[#747775] dark:text-[#c4c7c5] shrink-0 font-mono">
        Halaman {currentPage} dari {totalPages}
      </div>
    {/if}
  </div>

  <!-- TAB VIEW 1: ITEM CATALOG TABLE -->
  {#if activeTab !== "CONFIG"}
    <!-- Search and Filter Bar -->
    <div
      class="flex flex-wrap items-center justify-between gap-3 bg-[#f8fafd] dark:bg-[#1e1f20] p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#444746]"
    >
      <div class="flex-1 min-w-[240px] relative">
        <Search
          class="w-4 h-4 text-[#747775] dark:text-[#c4c7c5] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
        />
        <input
          type="text"
          bind:value={searchQuery}
          onkeydown={(e) => e.key === "Enter" && loadItems()}
          placeholder="Cari kode SKU/KFA, nama item, atau nama generik (tekan Enter)..."
          class="w-full h-10 pl-10 pr-4 rounded-xl bg-white dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#444746] focus:border-[#0b57d0] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none transition-all shadow-2xs"
        />
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <!-- Filter Kategori -->
        <div class="w-48">
          <M3Select
            bind:value={filterCategoryId}
            options={categoryFilterOptions}
            onchange={loadItems}
            placeholder="Semua Kategori"
          />
        </div>

        <!-- Filter Lini Produk -->
        <div class="w-48">
          <M3Select
            bind:value={filterProductLineId}
            options={productLineFilterOptions}
            onchange={loadItems}
            placeholder="Semua Lini Produk"
          />
        </div>

        <!-- Filter Status -->
        <div class="w-36">
          <M3Select
            bind:value={filterStatus}
            options={statusFilterOptions}
            onchange={loadItems}
            placeholder="Semua Status"
          />
        </div>
      </div>
    </div>

    <!-- Data Table Subcomponent -->
    <ItemTable
      {items}
      {isLoading}
      {currentPage}
      {totalPages}
      {totalItems}
      onPageChange={(p) => {
        currentPage = p;
        loadItems();
      }}
      onViewDetail={handleOpenDetail}
      onEdit={handleOpenEdit}
      onDelete={handleOpenDelete}
      onCreate={() => handleOpenCreate()}
    />

    <!-- TAB VIEW 2: KELOLA KATEGORI & LINI PRODUK (CONFIG) -->
  {:else}
    <ItemConfigTab
      {categories}
      {productLines}
      onRefresh={loadReferences}
      {showToast}
    />
  {/if}
</div>

<!-- ======================================================== -->
<!-- MODAL KOMPONEN                                           -->
<!-- ======================================================== -->

<!-- 1. Modal Rincian Item (Detail Modal) -->
<ItemDetailModal
  isOpen={showDetailModal}
  itemSummary={selectedDetailItem}
  onClose={() => (showDetailModal = false)}
  onEdit={(item) => {
    showDetailModal = false;
    handleOpenEdit(item);
  }}
  {showToast}
/>

<!-- 2. Modal Form Tambah / Edit Item -->
<ItemFormModal
  isOpen={showFormModal}
  mode={formMode}
  initialItem={selectedEditItem}
  initialType={formTargetType}
  {categories}
  {productLines}
  onClose={() => (showFormModal = false)}
  onSaved={() => {
    showFormModal = false;
    loadItems();
  }}
  {showToast}
/>

<!-- 3. Modal Konfirmasi Hapus -->
<ItemDeleteModal
  isOpen={showDeleteModal}
  item={itemToDelete}
  isLoading={isActionLoading}
  onClose={() => (showDeleteModal = false)}
  onConfirm={handleConfirmDelete}
/>
