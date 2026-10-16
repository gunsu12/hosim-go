<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    Search,
    Plus,
    ShieldCheck,
    Building2,
    RefreshCw,
    AlertCircle,
    CheckCircle2,
    Trash2,
    Edit3,
    Eye,
    Filter,
    CreditCard,
    PhoneCall,
    Check,
    X,
  } from "@lucide/svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3Select from "$lib/components/m3/M3Select.svelte";
  import { getCustomers, getCustomerTypes } from "$lib/api/finance";
  import type { CustomerRecord, CustomerTypeRecord } from "$lib/types/finance";

  import CustomerModal from "./components/CustomerModal.svelte";
  import CustomerDeleteModal from "./components/CustomerDeleteModal.svelte";
  import CustomerDetailModal from "./components/CustomerDetailModal.svelte";

  // State Data
  let customers = $state<CustomerRecord[]>([]);
  let customerTypes = $state<CustomerTypeRecord[]>([]);
  let totalCount = $state<number>(0);
  let isLoading = $state<boolean>(true);
  let errorMessage = $state<string | null>(null);

  // Filter & Search
  let searchQuery = $state<string>("");
  let selectedTypeId = $state<string>("");
  let selectedStatus = $state<string>("all"); // 'all' | 'active' | 'inactive'
  let currentPage = $state<number>(1);
  let pageSize = $state<number>(10);

  // Modal States
  let isModalOpen = $state<boolean>(false);
  let modalMode = $state<"CREATE" | "EDIT">("CREATE");
  let selectedCustomer = $state<CustomerRecord | null>(null);

  let isDeleteModalOpen = $state<boolean>(false);
  let customerToDelete = $state<CustomerRecord | null>(null);

  let isDetailModalOpen = $state<boolean>(false);
  let customerToView = $state<CustomerRecord | null>(null);

  // Alert Feedback
  let toastMessage = $state<string | null>(null);
  let toastType = $state<"success" | "error">("success");

  function showToast(msg: string, type: "success" | "error" = "success") {
    toastMessage = msg;
    toastType = type;
    setTimeout(() => {
      toastMessage = null;
    }, 4000);
  }

  // Load Data
  async function loadData() {
    isLoading = true;
    errorMessage = null;

    try {
      const [custRes, typeRes] = await Promise.all([
        getCustomers({
          search: searchQuery.trim() || undefined,
          customer_type_id: selectedTypeId || undefined,
          is_active:
            selectedStatus === "all" ? undefined : selectedStatus === "active",
          page: currentPage,
          limit: pageSize,
        }),
        customerTypes.length === 0
          ? getCustomerTypes()
          : Promise.resolve(customerTypes),
      ]);

      customers = custRes.data || [];
      totalCount = custRes.count ?? customers.length;
      if (typeRes && typeRes.length > 0) {
        customerTypes = typeRes;
      }
    } catch (err: any) {
      console.error("Failed to load customers:", err);
      errorMessage =
        err?.message || "Gagal terhubung ke server untuk memuat data penjamin.";
    } finally {
      isLoading = false;
    }
  }

  onMount(() => {
    loadData();
  });

  onDestroy(() => {
    if (searchTimer) clearTimeout(searchTimer);
  });

  // Debounce search input
  let searchTimer: any;
  function handleSearchInput() {
    clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      currentPage = 1;
      loadData();
    }, 350);
  }

  function handleFilterChange() {
    currentPage = 1;
    loadData();
  }

  // Pagination calculation
  let totalPages = $derived(Math.max(1, Math.ceil(totalCount / pageSize)));

  function openCreateModal() {
    selectedCustomer = null;
    modalMode = "CREATE";
    isModalOpen = true;
  }

  function openEditModal(c: CustomerRecord) {
    selectedCustomer = c;
    modalMode = "EDIT";
    isModalOpen = true;
  }

  function openDetailModal(c: CustomerRecord) {
    customerToView = c;
    isDetailModalOpen = true;
  }

  function openDeleteModal(c: CustomerRecord) {
    customerToDelete = c;
    isDeleteModalOpen = true;
  }

  function handleModalSuccess(saved: CustomerRecord) {
    isModalOpen = false;
    showToast(
      modalMode === "CREATE"
        ? `Debitur "${saved.name}" berhasil ditambahkan.`
        : `Debitur "${saved.name}" berhasil diperbarui.`,
    );
    loadData();
  }

  function handleDeleteSuccess() {
    isDeleteModalOpen = false;
    showToast("Debitur/penjamin berhasil dihapus.");
    loadData();
  }
</script>

<div class="flex flex-col gap-5">
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
        class="ml-2 opacity-70 hover:opacity-100"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    </div>
  {/if}

  <!-- Header Section -->
  <div class="flex flex-wrap items-center justify-between gap-4">
    <div>
      <div class="flex items-center gap-2.5">
        <h2
          class="text-xl font-bold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight"
        >
          Master Penjamin (Customer)
        </h2>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">
        Konfigurasi penjamin pembiayaan pasien (BPJS Kesehatan, Asuransi Swasta,
        Korporat/Perusahaan, & Pasien Umum).
      </p>
    </div>
  </div>

  <!-- Search & Filter Bar -->
  <div
    class="flex flex-wrap items-center justify-between gap-3 bg-[#f8fafd] dark:bg-[#1e1f20] p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#444746]"
  >
    <div class="flex flex-wrap items-center gap-3 flex-1 min-w-[280px]">
      <!-- Search Input -->
      <div class="relative flex-1 min-w-[200px] max-w-md">
        <Search
          class="w-4 h-4 text-[#747775] dark:text-[#c4c7c5] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
        />
        <input
          type="text"
          bind:value={searchQuery}
          oninput={handleSearchInput}
          placeholder="Cari kode, nama penjamin, atau PIC..."
          class="w-full h-10 pl-10 pr-4 rounded-xl bg-white dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#444746] focus:border-[#0b57d0] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all shadow-2xs"
        />
      </div>

      <!-- Kategori Filter -->
      <div class="w-48">
        <M3Select
          bind:value={selectedTypeId}
          placeholder="Semua Kategori"
          onchange={handleFilterChange}
          options={[
            { value: "", label: "Semua Kategori" },
            ...customerTypes.map((ct) => ({
              value: ct.id,
              label: ct.name,
              badge: ct.code,
            })),
          ]}
        />
      </div>

      <!-- Status Filter -->
      <div class="w-36">
        <M3Select
          bind:value={selectedStatus}
          placeholder="Semua Status"
          onchange={handleFilterChange}
          options={[
            { value: "all", label: "Semua Status" },
            { value: "active", label: "Hanya Aktif" },
            { value: "inactive", label: "Hanya Non-Aktif" },
          ]}
        />
      </div>

      <div class="flex items-center gap-2.5">
        <M3Button variant="filled" onclick={openCreateModal}>
          <Plus class="w-4 h-4" />
          <span>Tambah Penjamin</span>
        </M3Button>
      </div>
    </div>
  </div>

  <!-- Error Alert Banner -->
  {#if errorMessage}
    <div
      class="p-4 rounded-2xl bg-rose-50 border border-rose-200 text-rose-700 flex items-center justify-between gap-3 text-xs"
    >
      <div class="flex items-center gap-2">
        <AlertCircle class="w-4 h-4 shrink-0" />
        <span>{errorMessage}</span>
      </div>
      <M3Button variant="tonal" onclick={loadData}>Coba Lagi</M3Button>
    </div>
  {/if}

  <!-- Data Table Section -->
  <div
    class="bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden"
  >
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr
            class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]"
          >
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Debitur / Penjamin</th>
            <th class="py-3 px-4 font-semibold">Kategori Penjamin</th>
            <th class="py-3 px-4 font-semibold">PIC / Kontak Hotline</th>
            <th class="py-3 px-4 font-semibold text-center">Wajib Kartu</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>

        <tbody
          class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]"
        >
          {#if isLoading}
            <tr>
              <td colspan="7" class="py-12 text-center text-[#747775]">
                <div class="inline-flex items-center gap-2">
                  <RefreshCw class="w-4 h-4 animate-spin text-[#0b57d0]" />
                  <span>Memuat data penjamin dari database...</span>
                </div>
              </td>
            </tr>
          {:else if customers.length === 0}
            <tr>
              <td colspan="7" class="py-12 text-center text-[#747775]">
                <Building2 class="w-10 h-10 mx-auto text-[#c4c7c5] mb-2" />
                <p class="font-semibold text-[#1f1f1f]">
                  Tidak ada data penjamin ditemukan
                </p>
                <p class="text-[11px] mt-0.5">
                  {searchQuery
                    ? `Tidak ada hasil pencarian "${searchQuery}". Coba kata kunci lain.`
                    : "Belum ada debitur/penjamin yang dikonfigurasi."}
                </p>
                <div class="mt-3">
                  <M3Button variant="tonal" onclick={openCreateModal}>
                    <Plus class="w-3.5 h-3.5 mr-1" />
                    <span>Tambah Penjamin Pertama</span>
                  </M3Button>
                </div>
              </td>
            </tr>
          {:else}
            {#each customers as c (c.id)}
              <tr
                class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/50 transition-colors"
              >
                <!-- Kode -->
                <td
                  class="py-3.5 px-4 font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa]"
                >
                  {c.code}
                </td>

                <!-- Nama & Catatan -->
                <td class="py-3.5 px-4">
                  <div class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
                    {c.name}
                  </div>
                  {#if c.address}
                    <div
                      class="text-[11px] text-[#747775] dark:text-[#8e918f] truncate max-w-xs"
                    >
                      {c.address}
                    </div>
                  {/if}
                </td>

                <!-- Kategori Tipe -->
                <td class="py-3.5 px-4">
                  <span
                    class="px-2 py-0.5 rounded-lg text-[11px] font-semibold bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] border border-[#e1e5ea] dark:border-[#353738]"
                  >
                    {c.customer_type?.name || "Umum"}
                  </span>
                </td>

                <!-- Kontak PIC -->
                <td class="py-3.5 px-4 text-[#444746] dark:text-[#c4c7c5]">
                  {#if c.contact_person}
                    <div class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">
                      {c.contact_person}
                    </div>
                  {/if}
                  {#if c.phone}
                    <div
                      class="font-mono text-[11px] text-[#0b57d0] dark:text-[#a8c7fa]"
                    >
                      {c.phone}
                    </div>
                  {:else if !c.contact_person}
                    <span class="text-[#747775] dark:text-[#8e918f]">-</span>
                  {/if}
                </td>

                <!-- Wajib Kartu -->
                <td class="py-3.5 px-4 text-center">
                  {#if c.require_card}
                    <span
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300"
                    >
                      Wajib Kartu/SEP
                    </span>
                  {:else}
                    <span
                      class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300"
                    >
                      Tanpa Kartu
                    </span>
                  {/if}
                </td>

                <!-- Status Aktif -->
                <td class="py-3.5 px-4 text-center">
                  {#if c.is_active}
                    <span
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300"
                    >
                      AKTIF
                    </span>
                  {:else}
                    <span
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-gray-200 dark:bg-gray-800 text-gray-700 dark:text-gray-300"
                    >
                      NON-AKTIF
                    </span>
                  {/if}
                </td>

                <!-- Aksi -->
                <td class="py-3.5 px-4 text-right">
                  <div class="inline-flex items-center justify-end gap-1">
                    <button
                      type="button"
                      onclick={() => openDetailModal(c)}
                      class="p-1.5 rounded-lg text-[#747775] dark:text-[#8e918f] hover:text-[#0b57d0] dark:hover:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] transition-colors cursor-pointer"
                      title="Lihat Detail Penjamin"
                    >
                      <Eye class="w-4 h-4" />
                    </button>

                    <button
                      type="button"
                      onclick={() => openEditModal(c)}
                      class="p-1.5 rounded-lg text-[#747775] dark:text-[#8e918f] hover:text-[#0b57d0] dark:hover:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] transition-colors cursor-pointer"
                      title="Edit Data Penjamin"
                    >
                      <Edit3 class="w-4 h-4" />
                    </button>

                    <button
                      type="button"
                      onclick={() => openDeleteModal(c)}
                      class="p-1.5 rounded-lg text-[#747775] dark:text-[#8e918f] hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-950/50 transition-colors cursor-pointer"
                      title="Hapus Penjamin"
                    >
                      <Trash2 class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>

    <!-- Pagination Footer -->
    {#if totalCount > 0}
      <div
        class="px-5 py-3 border-t border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a] flex flex-wrap items-center justify-between gap-3 text-xs"
      >
        <div class="text-[#444746] dark:text-[#c4c7c5]">
          Menampilkan baris <strong>{(currentPage - 1) * pageSize + 1}</strong>
          – <strong>{Math.min(currentPage * pageSize, totalCount)}</strong> dari
          <strong>{totalCount}</strong> penjamin
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            disabled={currentPage <= 1 || isLoading}
            onclick={() => {
              currentPage--;
              loadData();
            }}
            class="px-3 py-1.5 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-[#1f1f1f] dark:text-[#e3e3e3] font-medium hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
          >
            Sebelumnya
          </button>

          <span
            class="font-mono px-2 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]"
          >
            {currentPage} / {totalPages}
          </span>

          <button
            type="button"
            disabled={currentPage >= totalPages || isLoading}
            onclick={() => {
              currentPage++;
              loadData();
            }}
            class="px-3 py-1.5 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-[#1f1f1f] dark:text-[#e3e3e3] font-medium hover:bg-[#f0f4f9] dark:hover:bg-[#28292a] disabled:opacity-40 disabled:cursor-not-allowed transition-colors cursor-pointer"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>

<!-- Modal Components -->
<CustomerModal
  isOpen={isModalOpen}
  mode={modalMode}
  customer={selectedCustomer}
  {customerTypes}
  onClose={() => (isModalOpen = false)}
  onSuccess={handleModalSuccess}
/>

<CustomerDeleteModal
  isOpen={isDeleteModalOpen}
  customer={customerToDelete}
  onClose={() => (isDeleteModalOpen = false)}
  onSuccess={handleDeleteSuccess}
/>

<CustomerDetailModal
  isOpen={isDetailModalOpen}
  customer={customerToView}
  onClose={() => (isDetailModalOpen = false)}
  onEdit={(c) => openEditModal(c)}
/>
