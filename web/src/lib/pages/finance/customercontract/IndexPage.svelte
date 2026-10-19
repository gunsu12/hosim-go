<script lang="ts">
  import { onMount } from "svelte";
  import {
    Plus,
    RefreshCw,
    Search,
    CheckCircle2,
    AlertCircle,
    X,
  } from "@lucide/svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3Select, {
    type SelectOption,
  } from "$lib/components/m3/M3Select.svelte";
  import type {
    CustomerContractRecord,
    CreateContractDTO,
    UpdateContractDraftDTO,
    UpdateContractInfoDTO,
    AmendContractTermsDTO,
    RenewContractDTO,
  } from "$lib/types/finance/customer_contract";
  import type { CustomerRecord } from "$lib/types/finance/customer";
  import {
    getCustomerContracts,
    createCustomerContract,
    updateCustomerContract,
    amendCustomerContract,
    activateCustomerContract,
    reactivateCustomerContract,
    suspendCustomerContract,
    terminateCustomerContract,
    renewCustomerContract,
    deleteCustomerContract,
    setContractPriceSchemas,
  } from "$lib/api/finance/customer_contract";
  import { getCustomers } from "$lib/api/finance/customer";

  import ContractTable from "./components/ContractTable.svelte";
  import ContractFormModal from "./components/ContractFormModal.svelte";
  import ContractDetailModal from "./components/ContractDetailModal.svelte";
  import ContractAmendModal from "./components/ContractAmendModal.svelte";
  import ContractStatusModal, {
    type TransitionAction,
  } from "./components/ContractStatusModal.svelte";
  import ContractRenewModal from "./components/ContractRenewModal.svelte";
  import ContractDeleteModal from "./components/ContractDeleteModal.svelte";
  import ContractPriceSchemasModal from "./components/ContractPriceSchemasModal.svelte";

  // Data states
  let contracts = $state<CustomerContractRecord[]>([]);
  let customers = $state<CustomerRecord[]>([]);
  let isLoading = $state(true);
  let isSaving = $state(false);

  // Search & Filters
  let searchQuery = $state("");
  let selectedCustomerId = $state("");
  let activeStatusTab = $state<
    "ALL" | "ACTIVE" | "DRAFT" | "SUSPENDED" | "TERMINATED" | "EXPIRING"
  >("ALL");

  // Modals state
  let isFormModalOpen = $state(false);
  let formModalMode = $state<"CREATE" | "EDIT">("CREATE");
  let selectedContractForForm = $state<CustomerContractRecord | null>(null);

  let isDetailModalOpen = $state(false);
  let selectedContractForDetail = $state<CustomerContractRecord | null>(null);

  let isAmendModalOpen = $state(false);
  let selectedContractForAmend = $state<CustomerContractRecord | null>(null);

  let isStatusModalOpen = $state(false);
  let statusTransitionAction = $state<TransitionAction>("ACTIVATE");
  let selectedContractForStatus = $state<CustomerContractRecord | null>(null);

  let isRenewModalOpen = $state(false);
  let selectedContractForRenew = $state<CustomerContractRecord | null>(null);

  let isDeleteModalOpen = $state(false);
  let selectedContractForDelete = $state<CustomerContractRecord | null>(null);

  let isPriceSchemasModalOpen = $state(false);
  let selectedContractForSchemas = $state<CustomerContractRecord | null>(null);

  // Toast feedback
  let toastMessage = $state<string | null>(null);
  let toastType = $state<"success" | "error">("success");

  function showToast(msg: string, isError = false) {
    toastMessage = msg;
    toastType = isError ? "error" : "success";
    setTimeout(() => {
      toastMessage = null;
    }, 4000);
  }

  const customerOptions = $derived<SelectOption[]>([
    { value: "", label: "Semua Penjamin" },
    ...customers.map((c) => ({
      value: c.id,
      label: `${c.code} • ${c.name}`,
    })),
  ]);

  // Statistics counts for the pill filter
  const stats = $derived.by(() => {
    const total = contracts.length;
    const active = contracts.filter((c) => c.status === "ACTIVE").length;
    const draft = contracts.filter((c) => c.status === "DRAFT").length;
    const suspended = contracts.filter((c) => c.status === "SUSPENDED").length;
    const terminated = contracts.filter(
      (c) => c.status === "TERMINATED",
    ).length;
    const expiring = contracts.filter((c) => {
      if (c.status !== "ACTIVE" || !c.expiration_date) return false;
      const exp = new Date(c.expiration_date);
      const now = new Date();
      now.setHours(0, 0, 0, 0);
      const diff = Math.ceil(
        (exp.getTime() - now.getTime()) / (1000 * 60 * 60 * 24),
      );
      return diff >= 0 && diff <= 30;
    }).length;

    return { total, active, draft, suspended, terminated, expiring };
  });

  // Client filtered contracts based on search, customer, and status tab
  const filteredContracts = $derived.by(() => {
    let result = contracts;

    if (activeStatusTab === "ACTIVE") {
      result = result.filter((c) => c.status === "ACTIVE");
    } else if (activeStatusTab === "DRAFT") {
      result = result.filter((c) => c.status === "DRAFT");
    } else if (activeStatusTab === "SUSPENDED") {
      result = result.filter((c) => c.status === "SUSPENDED");
    } else if (activeStatusTab === "TERMINATED") {
      result = result.filter((c) => c.status === "TERMINATED");
    } else if (activeStatusTab === "EXPIRING") {
      const now = new Date();
      now.setHours(0, 0, 0, 0);
      result = result.filter((c) => {
        if (c.status !== "ACTIVE" || !c.expiration_date) return false;
        const exp = new Date(c.expiration_date);
        const diff = Math.ceil(
          (exp.getTime() - now.getTime()) / (1000 * 60 * 60 * 24),
        );
        return diff >= 0 && diff <= 30;
      });
    }

    if (selectedCustomerId) {
      result = result.filter((c) => c.customer_id === selectedCustomerId);
    }

    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      result = result.filter(
        (c) =>
          c.contract_number.toLowerCase().includes(q) ||
          c.name.toLowerCase().includes(q) ||
          (c.customer_name && c.customer_name.toLowerCase().includes(q)) ||
          (c.customer_code && c.customer_code.toLowerCase().includes(q)),
      );
    }

    return result;
  });

  async function loadData() {
    isLoading = true;
    try {
      const [contractRes, custRes] = await Promise.all([
        getCustomerContracts({ limit: 100 }),
        customers.length === 0
          ? getCustomers({ limit: 100 })
          : Promise.resolve({ data: customers }),
      ]);

      contracts = contractRes.items || [];
      if (custRes.data && custRes.data.length > 0) {
        customers = custRes.data;
      }
    } catch (err: any) {
      console.error("Gagal memuat kontrak:", err);
      showToast(err?.message || "Gagal memuat data kontrak", true);
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  // Action handlers
  function openCreateModal() {
    selectedContractForForm = null;
    formModalMode = "CREATE";
    isFormModalOpen = true;
  }

  function handleEditContract(contract: CustomerContractRecord) {
    selectedContractForForm = contract;
    formModalMode = "EDIT";
    isFormModalOpen = true;
  }

  function handleViewDetail(contract: CustomerContractRecord) {
    selectedContractForDetail = contract;
    isDetailModalOpen = true;
  }

  function handleAmendContract(contract: CustomerContractRecord) {
    selectedContractForAmend = contract;
    isAmendModalOpen = true;
  }

  function handleStatusTransition(
    contract: CustomerContractRecord,
    action: TransitionAction,
  ) {
    selectedContractForStatus = contract;
    statusTransitionAction = action;
    isStatusModalOpen = true;
  }

  function handleRenewContract(contract: CustomerContractRecord) {
    selectedContractForRenew = contract;
    isRenewModalOpen = true;
  }

  function handleDeleteContract(contract: CustomerContractRecord) {
    selectedContractForDelete = contract;
    isDeleteModalOpen = true;
  }

  function handleManageSchemas(contract: CustomerContractRecord) {
    selectedContractForSchemas = contract;
    isPriceSchemasModalOpen = true;
  }

  // Submit operations
  async function handleFormSubmit(
    dto: CreateContractDTO | UpdateContractDraftDTO | UpdateContractInfoDTO,
  ) {
    isSaving = true;
    try {
      if (formModalMode === "CREATE") {
        const created = await createCustomerContract(dto as CreateContractDTO);
        isFormModalOpen = false;
        showToast(`Draft kontrak '${created.contract_number}' berhasil dibuat`);
      } else if (selectedContractForForm) {
        const updated = await updateCustomerContract(
          selectedContractForForm.id,
          dto as UpdateContractDraftDTO,
        );
        isFormModalOpen = false;
        showToast(
          `Perubahan kontrak '${updated.contract_number}' berhasil disimpan`,
        );
      }
      await loadData();
    } catch (e: any) {
      showToast(e?.message || "Gagal menyimpan kontrak", true);
    } finally {
      isSaving = false;
    }
  }

  async function handleAmendSubmit(dto: AmendContractTermsDTO) {
    if (!selectedContractForAmend) return;
    isSaving = true;
    try {
      const res = await amendCustomerContract(selectedContractForAmend.id, dto);
      isAmendModalOpen = false;
      showToast(`Amandemen PKS '${res.contract_number}' berhasil disimpan`);
      await loadData();
    } catch (e: any) {
      showToast(e?.message || "Gagal menyimpan amandemen", true);
    } finally {
      isSaving = false;
    }
  }

  async function handleStatusConfirm(reason?: string) {
    if (!selectedContractForStatus) return;
    isSaving = true;
    try {
      const targetId = selectedContractForStatus.id;
      let res: CustomerContractRecord;
      switch (statusTransitionAction) {
        case "ACTIVATE":
          res = await activateCustomerContract(targetId);
          showToast(`Kontrak PKS '${res.contract_number}' aktif`);
          break;
        case "SUSPEND":
          res = await suspendCustomerContract(targetId, reason || "");
          showToast(`Kontrak PKS '${res.contract_number}' ditangguhkan`);
          break;
        case "REACTIVATE":
          res = await reactivateCustomerContract(targetId);
          showToast(`Kontrak PKS '${res.contract_number}' aktif kembali`);
          break;
        case "TERMINATE":
          res = await terminateCustomerContract(targetId, reason || "");
          showToast(`Kontrak PKS '${res.contract_number}' diputus`);
          break;
      }
      isStatusModalOpen = false;
      await loadData();
    } catch (e: any) {
      showToast(e?.message || "Gagal mengubah status", true);
    } finally {
      isSaving = false;
    }
  }

  async function handleRenewSubmit(dto: RenewContractDTO) {
    if (!selectedContractForRenew) return;
    isSaving = true;
    try {
      const renewed = await renewCustomerContract(
        selectedContractForRenew.id,
        dto,
      );
      isRenewModalOpen = false;
      showToast(
        `Perpanjangan PKS '${renewed.contract_number}' berhasil dibuat`,
      );
      await loadData();
    } catch (e: any) {
      showToast(e?.message || "Gagal memperpanjang kontrak", true);
    } finally {
      isSaving = false;
    }
  }

  async function handleDeleteConfirm() {
    if (!selectedContractForDelete) return;
    isSaving = true;
    try {
      await deleteCustomerContract(selectedContractForDelete.id);
      isDeleteModalOpen = false;
      showToast(
        `Draft kontrak '${selectedContractForDelete.contract_number}' berhasil dihapus`,
      );
      await loadData();
    } catch (e: any) {
      showToast(e?.message || "Gagal menghapus draft", true);
    } finally {
      isSaving = false;
    }
  }

  async function handlePriceSchemasSubmit(
    schemas: Array<{
      price_schema_id: string;
      priority: number;
      is_active: boolean;
      price_schema_name?: string;
      price_schema_code?: string;
    }>,
  ) {
    if (!selectedContractForSchemas) return;
    isSaving = true;
    try {
      const res = await setContractPriceSchemas(
        selectedContractForSchemas.id,
        schemas,
      );
      isPriceSchemasModalOpen = false;
      showToast(`Skema tarif '${res.contract_number}' berhasil diperbarui`);
      await loadData();
      if (
        selectedContractForDetail &&
        selectedContractForDetail.id === res.id
      ) {
        selectedContractForDetail = res;
      }
    } catch (e: any) {
      showToast(e?.message || "Gagal memperbarui skema", true);
    } finally {
      isSaving = false;
    }
  }
</script>

<div class="flex flex-col gap-4 max-w-7xl mx-auto w-full pb-8">
  <!-- Toast Feedback -->
  {#if toastMessage}
    <div
      class={`fixed bottom-6 right-6 z-50 px-4 py-3 rounded-2xl shadow-xl flex items-center gap-3 text-xs font-semibold animate-in slide-in-from-bottom-5 duration-200 ${
        toastType === "success"
          ? "bg-[#0f5223] text-white"
          : "bg-[#ba1a1a] text-white"
      }`}
    >
      {#if toastType === "success"}
        <CheckCircle2 class="w-4 h-4 shrink-0" />
      {:else}
        <AlertCircle class="w-4 h-4 shrink-0" />
      {/if}
      <span>{toastMessage}</span>
      <button
        onclick={() => (toastMessage = null)}
        type="button"
        class="ml-2 opacity-70 hover:opacity-100 cursor-pointer"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  <!-- Header Section -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <h2
        class="text-xl font-bold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight"
      >
        Kontrak Kerjasama Penjamin (PKS)
      </h2>
      <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">
        Perjanjian kerja sama penjamin, batas kredit klaim, dan pemetaan skema
        tarif pelayanan.
      </p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button
        variant="tonal"
        onclick={loadData}
        disabled={isLoading}
        class="h-10"
      >
        <RefreshCw class="w-4 h-4 mr-1.5 {isLoading ? 'animate-spin' : ''}" />
        <span>Refresh</span>
      </M3Button>
      <M3Button variant="filled" onclick={openCreateModal} class="h-10">
        <Plus class="w-4 h-4 mr-1.5" />
        <span>Tambah Kontrak</span>
      </M3Button>
    </div>
  </div>

  <!-- Status Filter Pills (Simple & Clean) -->
  <div class="flex items-center gap-1.5 overflow-x-auto pb-1 text-xs">
    <button
      type="button"
      onclick={() => (activeStatusTab = "ALL")}
      class="px-3 py-1.5 rounded-xl font-medium transition-all shrink-0 cursor-pointer flex items-center gap-1.5 {activeStatusTab ===
      'ALL'
        ? 'bg-[#0b57d0] text-white font-semibold shadow-2xs'
        : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea]'}"
    >
      <span>Semua</span>
      <span
        class="px-1.5 py-0.2 rounded-full text-[10px] {activeStatusTab === 'ALL'
          ? 'bg-white/20 text-white'
          : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'}"
      >
        {stats.total}
      </span>
    </button>

    <button
      type="button"
      onclick={() => (activeStatusTab = "ACTIVE")}
      class="px-3 py-1.5 rounded-xl font-medium transition-all shrink-0 cursor-pointer flex items-center gap-1.5 {activeStatusTab ===
      'ACTIVE'
        ? 'bg-emerald-600 text-white font-semibold shadow-2xs'
        : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea]'}"
    >
      <span>Aktif</span>
      <span
        class="px-1.5 py-0.2 rounded-full text-[10px] {activeStatusTab ===
        'ACTIVE'
          ? 'bg-white/20 text-white'
          : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'}"
      >
        {stats.active}
      </span>
    </button>

    <button
      type="button"
      onclick={() => (activeStatusTab = "DRAFT")}
      class="px-3 py-1.5 rounded-xl font-medium transition-all shrink-0 cursor-pointer flex items-center gap-1.5 {activeStatusTab ===
      'DRAFT'
        ? 'bg-amber-600 text-white font-semibold shadow-2xs'
        : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea]'}"
    >
      <span>Draft</span>
      <span
        class="px-1.5 py-0.2 rounded-full text-[10px] {activeStatusTab ===
        'DRAFT'
          ? 'bg-white/20 text-white'
          : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'}"
      >
        {stats.draft}
      </span>
    </button>

    <button
      type="button"
      onclick={() => (activeStatusTab = "SUSPENDED")}
      class="px-3 py-1.5 rounded-xl font-medium transition-all shrink-0 cursor-pointer flex items-center gap-1.5 {activeStatusTab ===
      'SUSPENDED'
        ? 'bg-orange-600 text-white font-semibold shadow-2xs'
        : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea]'}"
    >
      <span>Ditangguhkan</span>
      <span
        class="px-1.5 py-0.2 rounded-full text-[10px] {activeStatusTab ===
        'SUSPENDED'
          ? 'bg-white/20 text-white'
          : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'}"
      >
        {stats.suspended}
      </span>
    </button>

    {#if stats.expiring > 0}
      <button
        type="button"
        onclick={() => (activeStatusTab = "EXPIRING")}
        class="px-3 py-1.5 rounded-xl font-medium transition-all shrink-0 cursor-pointer flex items-center gap-1.5 {activeStatusTab ===
        'EXPIRING'
          ? 'bg-purple-600 text-white font-semibold shadow-2xs'
          : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea]'}"
      >
        <span>Sisa &le; 30 Hari</span>
        <span
          class="px-1.5 py-0.2 rounded-full text-[10px] {activeStatusTab ===
          'EXPIRING'
            ? 'bg-white/20 text-white'
            : 'bg-purple-100 text-purple-700 dark:bg-purple-950/60 dark:text-purple-300 font-bold'}"
        >
          {stats.expiring}
        </span>
      </button>
    {/if}

    {#if stats.terminated > 0}
      <button
        type="button"
        onclick={() => (activeStatusTab = "TERMINATED")}
        class="px-3 py-1.5 rounded-xl font-medium transition-all shrink-0 cursor-pointer flex items-center gap-1.5 {activeStatusTab ===
        'TERMINATED'
          ? 'bg-rose-600 text-white font-semibold shadow-2xs'
          : 'bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e1e5ea]'}"
      >
        <span>Diputus</span>
        <span
          class="px-1.5 py-0.2 rounded-full text-[10px] {activeStatusTab ===
          'TERMINATED'
            ? 'bg-white/20 text-white'
            : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'}"
        >
          {stats.terminated}
        </span>
      </button>
    {/if}
  </div>

  <!-- Single-Row Search & Filter Bar (Minimalist & Focused) -->
  <div
    class="flex flex-wrap items-center justify-between gap-3 bg-[#f8fafd] dark:bg-[#1e1f20] p-2.5 rounded-2xl border border-[#e1e5ea] dark:border-[#353738]"
  >
    <div class="flex flex-wrap items-center gap-2.5 flex-1 min-w-[260px]">
      <!-- Search Input -->
      <div class="relative flex-1 min-w-[200px] max-w-md">
        <Search
          class="w-4 h-4 text-[#747775] dark:text-[#c4c7c5] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
        />
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Cari nomor kontrak, nama, atau penjamin..."
          class="w-full h-10 pl-10 pr-8 rounded-xl bg-white dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all shadow-2xs"
        />
        {#if searchQuery}
          <button
            type="button"
            onclick={() => (searchQuery = "")}
            class="absolute right-3 top-1/2 -translate-y-1/2 text-[#747775] hover:text-[#1f1f1f] cursor-pointer"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        {/if}
      </div>

      <!-- Filter Penjamin -->
      <div class="w-56">
        <M3Select
          options={customerOptions}
          bind:value={selectedCustomerId}
          searchable
          placeholder="Semua Penjamin"
        />
      </div>
    </div>

    <div
      class="text-[11px] text-[#747775] dark:text-[#8e918f] font-medium pr-2"
    >
      Menampilkan <strong>{filteredContracts.length}</strong> dari {contracts.length}
      kontrak
    </div>
  </div>

  <!-- Clean Data Table -->
  <ContractTable
    contracts={filteredContracts}
    {isLoading}
    onViewDetail={handleViewDetail}
    onEdit={handleEditContract}
    onAmend={handleAmendContract}
    onManageSchemas={handleManageSchemas}
    onStatusTransition={handleStatusTransition}
    onRenew={handleRenewContract}
    onDelete={handleDeleteContract}
  />
</div>

<!-- MODALS -->
<ContractFormModal
  isOpen={isFormModalOpen}
  mode={formModalMode}
  contract={selectedContractForForm}
  isSubmitting={isSaving}
  onClose={() => (isFormModalOpen = false)}
  onSubmit={handleFormSubmit}
/>

<ContractDetailModal
  isOpen={isDetailModalOpen}
  contract={selectedContractForDetail}
  onClose={() => (isDetailModalOpen = false)}
  onEdit={handleEditContract}
  onAmend={handleAmendContract}
  onManageSchemas={handleManageSchemas}
  onStatusTransition={handleStatusTransition}
  onRenew={handleRenewContract}
  onDelete={handleDeleteContract}
/>

<ContractAmendModal
  isOpen={isAmendModalOpen}
  contract={selectedContractForAmend}
  isSubmitting={isSaving}
  onClose={() => (isAmendModalOpen = false)}
  onSubmit={handleAmendSubmit}
/>

<ContractStatusModal
  isOpen={isStatusModalOpen}
  contract={selectedContractForStatus}
  action={statusTransitionAction}
  isSubmitting={isSaving}
  onClose={() => (isStatusModalOpen = false)}
  onConfirm={handleStatusConfirm}
/>

<ContractRenewModal
  isOpen={isRenewModalOpen}
  contract={selectedContractForRenew}
  isSubmitting={isSaving}
  onClose={() => (isRenewModalOpen = false)}
  onSubmit={handleRenewSubmit}
/>

<ContractDeleteModal
  isOpen={isDeleteModalOpen}
  contract={selectedContractForDelete}
  isDeleting={isSaving}
  onClose={() => (isDeleteModalOpen = false)}
  onConfirm={handleDeleteConfirm}
/>

<ContractPriceSchemasModal
  isOpen={isPriceSchemasModalOpen}
  contract={selectedContractForSchemas}
  isSubmitting={isSaving}
  onClose={() => (isPriceSchemasModalOpen = false)}
  onSubmit={handlePriceSchemasSubmit}
/>
