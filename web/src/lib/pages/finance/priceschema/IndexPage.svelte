<script lang="ts">
  import { onMount } from "svelte";
  import {
    Plus,
    RefreshCw,
    Search,
    FileSpreadsheet,
    CheckCircle2,
    AlertCircle,
    X,
    Receipt,
    Layers,
    ShieldAlert,
  } from "@lucide/svelte";
  import type {
    PriceSchemaRecord,
    CreatePriceSchemaDTO,
    UpdatePriceSchemaDTO,
  } from "$lib/types/finance/price_schema";
  import {
    getPriceSchemas,
    createPriceSchema,
    updatePriceSchema,
    deactivatePriceSchema,
  } from "$lib/api/finance/price_schema";

  import PriceSchemaTable from "./components/PriceSchemaTable.svelte";
  import PriceSchemaFormModal from "./components/PriceSchemaFormModal.svelte";
  import PriceSchemaDeleteModal from "./components/PriceSchemaDeleteModal.svelte";
  import SchemaDetailPage from "./SchemaDetailPage.svelte";

  // State Management
  let schemas = $state<PriceSchemaRecord[]>([]);
  let isLoading = $state(false);
  let isSaving = $state(false);
  let isProcessingStatus = $state(false);

  // Search & Filter
  let searchQuery = $state("");
  let statusFilter = $state<"ALL" | "ACTIVE" | "INACTIVE">("ALL");

  // Selected Schema for Detail View (if set, renders SchemaDetailPage)
  let selectedSchemaForDetail = $state<PriceSchemaRecord | null>(null);

  // Modal State
  let isFormModalOpen = $state(false);
  let formModalMode = $state<"create" | "edit">("create");
  let selectedSchemaForForm = $state<PriceSchemaRecord | null>(null);

  // Deactivate Modal State
  let isDeactivateModalOpen = $state(false);
  let selectedSchemaForDeactivate = $state<PriceSchemaRecord | null>(null);

  // Toast Feedback State
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

  // Filtered Schemas
  const filteredSchemas = $derived.by(() => {
    let result = schemas;

    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      result = result.filter(
        (s) =>
          s.name.toLowerCase().includes(q) ||
          s.code.toLowerCase().includes(q) ||
          (s.description && s.description.toLowerCase().includes(q)),
      );
    }

    if (statusFilter === "ACTIVE") {
      result = result.filter((s) => s.is_active);
    } else if (statusFilter === "INACTIVE") {
      result = result.filter((s) => !s.is_active);
    }

    return result;
  });

  // KPI Statistics
  const stats = $derived.by(() => {
    const total = schemas.length;
    const active = schemas.filter((s) => s.is_active).length;
    const inactive = total - active;
    return { total, active, inactive };
  });

  async function loadSchemas() {
    isLoading = true;
    errorMessage = null;
    try {
      const res = await getPriceSchemas({ limit: 100 });
      schemas = res.items || [];
    } catch (err: any) {
      showToast(err.message || "Gagal memuat daftar skema tarif", true);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadSchemas();
  });

  // Action Handlers
  function handleOpenCreateModal() {
    formModalMode = "create";
    selectedSchemaForForm = null;
    isFormModalOpen = true;
  }

  function handleOpenEditModal(schema: PriceSchemaRecord) {
    formModalMode = "edit";
    selectedSchemaForForm = schema;
    isFormModalOpen = true;
  }

  function handleOpenToggleStatus(schema: PriceSchemaRecord) {
    selectedSchemaForDeactivate = schema;
    isDeactivateModalOpen = true;
  }

  function handleManageSchema(schema: PriceSchemaRecord) {
    selectedSchemaForDetail = schema;
  }

  async function handleSaveSchema(
    payload: CreatePriceSchemaDTO | UpdatePriceSchemaDTO,
  ) {
    isSaving = true;
    try {
      if (formModalMode === "create") {
        const created = await createPriceSchema(
          payload as CreatePriceSchemaDTO,
        );
        showToast(`Skema ${created.code} berhasil dibuat.`);
      } else if (selectedSchemaForForm) {
        const updated = await updatePriceSchema(
          selectedSchemaForForm.id,
          payload as UpdatePriceSchemaDTO,
        );
        showToast(`Skema ${updated.code} berhasil diperbarui.`);
        if (selectedSchemaForDetail?.id === updated.id) {
          selectedSchemaForDetail = updated;
        }
      }
      isFormModalOpen = false;
      await loadSchemas();
    } catch (err: any) {
      showToast(err.message || "Gagal menyimpan skema tarif", true);
    } finally {
      isSaving = false;
    }
  }

  async function handleConfirmToggleStatus() {
    if (!selectedSchemaForDeactivate) return;
    isProcessingStatus = true;
    try {
      if (selectedSchemaForDeactivate.is_active) {
        await deactivatePriceSchema(selectedSchemaForDeactivate.id);
        showToast(
          `Skema ${selectedSchemaForDeactivate.code} berhasil dinonaktifkan.`,
        );
      } else {
        await updatePriceSchema(selectedSchemaForDeactivate.id, {
          name: selectedSchemaForDeactivate.name,
          is_active: true,
        });
        showToast(
          `Skema ${selectedSchemaForDeactivate.code} berhasil diaktifkan kembali.`,
        );
      }
      isDeactivateModalOpen = false;
      await loadSchemas();
    } catch (err: any) {
      showToast(err.message || "Gagal memperbarui status skema", true);
    } finally {
      isProcessingStatus = false;
    }
  }
</script>

{#if selectedSchemaForDetail}
  <!-- Sub-View: Detail & Aturan Hierarkis Skema -->
  <SchemaDetailPage
    schema={selectedSchemaForDetail}
    onBack={() => {
      selectedSchemaForDetail = null;
      loadSchemas();
    }}
    onEditHeader={(s) => handleOpenEditModal(s)}
  />
{:else}
  <!-- Master View: Daftar Skema Tarif Rekanan -->
  <div class="flex flex-col gap-5">
    <!-- Toast Feedback Alerts -->
    {#if successMessage}
      <div
        class="p-3.5 rounded-2xl bg-[#e6f4ea] dark:bg-emerald-950/40 border border-[#ceead6] dark:border-emerald-800/60 text-[#137333] dark:text-emerald-300 flex items-center justify-between shadow-xs animate-in fade-in duration-200"
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
        class="p-3.5 rounded-2xl bg-[#fce8e6] dark:bg-rose-950/40 border border-[#fad2cf] dark:border-rose-800/60 text-[#c5221f] dark:text-rose-300 flex items-center justify-between shadow-xs animate-in fade-in duration-200"
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
          <div>
            <h2 class="text-lg font-bold text-on-surface tracking-tight">
              Skema Tarif Rekanan (Price Schema)
            </h2>
            <p
              class="text-xs text-on-surface-variant mt-0.5 max-w-2xl leading-relaxed"
            >
              Kelola aturan penyesuaian tarif (diskon/markup/fixed price),
              pembagian beban tanggungan (*cost sharing*), dan pajak PKS rumah
              sakit tanpa menduplikasi buku tarif dasar.
            </p>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          onclick={loadSchemas}
          class="p-2 rounded-xl text-on-surface-variant hover:text-on-surface hover:bg-surface-container border border-outline-variant transition-colors cursor-pointer"
          title="Muat ulang data"
          disabled={isLoading}
        >
          <RefreshCw class="w-4 h-4 {isLoading ? 'animate-spin' : ''}" />
        </button>

        <button
          type="button"
          onclick={handleOpenCreateModal}
          class="px-4 py-2 text-xs font-bold rounded-xl bg-primary text-on-primary hover:bg-primary/90 transition-colors flex items-center gap-1.5 cursor-pointer shadow-xs"
        >
          <Plus class="w-4 h-4" />
          <span>Tambah Skema</span>
        </button>
      </div>
    </div>
    <!-- Filter & Search Toolbar -->
    <div class="flex flex-wrap items-center justify-between gap-3 text-xs">
      <div class="flex flex-wrap items-center gap-2 flex-1 max-w-md">
        <div class="relative flex-1 min-w-[200px]">
          <Search
            class="w-3.5 h-3.5 absolute left-3 top-2.5 text-on-surface-variant"
          />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Cari kode atau nama skema tarif..."
            class="w-full pl-8 pr-3 py-2 rounded-xl border border-outline-variant bg-surface text-xs focus:outline-hidden focus:border-primary"
          />
        </div>

        <select
          bind:value={statusFilter}
          class="px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface text-xs"
        >
          <option value="ALL">Semua Status</option>
          <option value="ACTIVE">Hanya Aktif</option>
          <option value="INACTIVE">Hanya Nonaktif</option>
        </select>
      </div>
    </div>

    <!-- Table of Price Schemas -->
    <PriceSchemaTable
      schemas={filteredSchemas}
      {isLoading}
      onManage={handleManageSchema}
      onEdit={handleOpenEditModal}
      onToggleStatus={handleOpenToggleStatus}
    />

    <!-- Modal Form (Create / Edit Header) -->
    <PriceSchemaFormModal
      isOpen={isFormModalOpen}
      mode={formModalMode}
      schema={selectedSchemaForForm}
      {isSaving}
      onClose={() => (isFormModalOpen = false)}
      onSave={handleSaveSchema}
    />

    <!-- Modal Delete / Deactivate -->
    <PriceSchemaDeleteModal
      isOpen={isDeactivateModalOpen}
      schema={selectedSchemaForDeactivate}
      isProcessing={isProcessingStatus}
      onClose={() => (isDeactivateModalOpen = false)}
      onConfirm={handleConfirmToggleStatus}
    />
  </div>
{/if}
