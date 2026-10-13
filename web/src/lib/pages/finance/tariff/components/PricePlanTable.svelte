<script lang="ts">
  import {
    Search,
    Plus,
    Edit3,
    Copy,
    Send,
    CheckCircle,
    Power,
    Archive,
    List,
    Layers,
    ChevronLeft,
    ChevronRight,
    Building2,
    Calendar,
  } from "@lucide/svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3Select from "$lib/components/m3/M3Select.svelte";
  import type {
    PricePlanRecord,
    PricePlanStatus,
  } from "$lib/types/finance/price_plan";

  interface Props {
    plans: PricePlanRecord[];
    isLoading: boolean;
    searchQuery: string;
    filterStatus: string;
    currentPage: number;
    totalPages: number;
    totalCount: number;
    onSearchChange: (q: string) => void;
    onStatusChange: (s: string) => void;
    onPageChange: (p: number) => void;
    onAdd: () => void;
    onEdit: (plan: PricePlanRecord) => void;
    onClone: (plan: PricePlanRecord) => void;
    onManageItems: (plan: PricePlanRecord) => void;
    onSubmitPlan: (plan: PricePlanRecord) => void;
    onApprovePlan: (plan: PricePlanRecord) => void;
    onActivatePlan: (plan: PricePlanRecord) => void;
    onArchivePlan: (plan: PricePlanRecord) => void;
  }

  let {
    plans,
    isLoading,
    searchQuery,
    filterStatus,
    currentPage,
    totalPages,
    totalCount,
    onSearchChange,
    onStatusChange,
    onPageChange,
    onAdd,
    onEdit,
    onClone,
    onManageItems,
    onSubmitPlan,
    onApprovePlan,
    onActivatePlan,
    onArchivePlan,
  }: Props = $props();

  let localSearch = $state("");

  $effect(() => {
    localSearch = searchQuery;
  });

  const statusFilterOptions = [
    { value: "ALL", label: "Semua Status" },
    { value: "DRAFT", label: "DRAFT (Penyusunan)" },
    { value: "SUBMITTED", label: "SUBMITTED (Diajukan)" },
    { value: "APPROVED", label: "APPROVED (Disetujui)" },
    { value: "ACTIVE", label: "ACTIVE (Aktif Digunakan)" },
    { value: "ARCHIVED", label: "ARCHIVED (Diarsipkan)" },
  ];

  const statusBadgeConfig: Record<
    PricePlanStatus,
    { label: string; bg: string; text: string; border: string }
  > = {
    DRAFT: {
      label: "DRAFT",
      bg: "bg-amber-50 dark:bg-amber-950/40",
      text: "text-amber-800 dark:text-amber-300",
      border: "border-amber-200 dark:border-amber-800",
    },
    SUBMITTED: {
      label: "SUBMITTED",
      bg: "bg-blue-50 dark:bg-blue-950/40",
      text: "text-blue-800 dark:text-blue-300",
      border: "border-blue-200 dark:border-blue-800",
    },
    APPROVED: {
      label: "APPROVED",
      bg: "bg-purple-50 dark:bg-purple-950/40",
      text: "text-purple-800 dark:text-purple-300",
      border: "border-purple-200 dark:border-purple-800",
    },
    ACTIVE: {
      label: "ACTIVE",
      bg: "bg-emerald-50 dark:bg-emerald-950/40",
      text: "text-emerald-800 dark:text-emerald-300",
      border: "border-emerald-200 dark:border-emerald-800",
    },
    ARCHIVED: {
      label: "ARCHIVED",
      bg: "bg-slate-100 dark:bg-slate-800",
      text: "text-slate-600 dark:text-slate-400",
      border: "border-slate-300 dark:border-slate-700",
    },
  };
</script>

<div class="flex flex-col gap-4">
  <!-- Controls Bar -->
  <div
    class="flex flex-wrap items-center justify-between gap-3 bg-[#f8fafd] dark:bg-[#1e1f20] p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#444746]"
  >
    <div class="flex-1 min-w-[240px] relative">
      <Search
        class="w-4 h-4 text-[#747775] dark:text-[#c4c7c5] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none"
      />
      <input
        type="text"
        bind:value={localSearch}
        onkeydown={(e) => e.key === "Enter" && onSearchChange(localSearch)}
        placeholder="Cari kode atau nama buku tarif (tekan Enter)..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-white dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#444746] focus:border-[#0b57d0] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all shadow-2xs"
      />
    </div>

    <div class="flex items-center gap-2 flex-wrap">
      <!-- Status Filter -->
      <div class="w-48">
        <M3Select
          value={filterStatus}
          options={statusFilterOptions}
          onchange={(val) => onStatusChange(val)}
          placeholder="Pilih Status"
        />
      </div>

      <M3Button variant="filled" onclick={onAdd}>
        <Plus class="w-4 h-4" />
        <span>Buat Buku Tarif</span>
      </M3Button>
    </div>
  </div>

  <!-- Table Container -->
  <div
    class="border border-[#e1e5ea] dark:border-[#444746] rounded-2xl bg-white dark:bg-[#1e1f20] overflow-hidden shadow-2xs"
  >
    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr
            class="bg-[#f0f4f9] dark:bg-[#282a2c] text-[#444746] dark:text-[#c4c7c5] font-semibold border-b border-[#e1e5ea] dark:border-[#444746]"
          >
            <th class="py-3 px-4 w-12 text-center">No</th>
            <th class="py-3 px-4 w-36">Kode SK / Dokumen</th>
            <th class="py-3 px-4">Nama Buku Tarif</th>
            <th class="py-3 px-4">Masa Berlaku</th>
            <th class="py-3 px-4">Penjamin / Rekanan</th>
            <th class="py-3 px-4 text-center">CITO</th>
            <th class="py-3 px-4 text-center">Status</th>
            <th class="py-3 px-4 text-center">Item</th>
            <th class="py-3 px-4 text-right">Aksi & Transisi Status</th>
          </tr>
        </thead>
        <tbody
          class="divide-y divide-[#e1e5ea] dark:divide-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3]"
        >
          {#if isLoading}
            <tr>
              <td
                colspan="9"
                class="py-12 text-center text-[#747775] dark:text-[#c4c7c5]"
              >
                <div class="flex flex-col items-center justify-center gap-2">
                  <div
                    class="w-6 h-6 border-2 border-[#0b57d0] dark:border-blue-400 border-t-transparent rounded-full animate-spin"
                  ></div>
                  <span>Memuat daftar buku tarif...</span>
                </div>
              </td>
            </tr>
          {:else if !plans || plans.length === 0}
            <tr>
              <td
                colspan="9"
                class="py-12 text-center text-[#747775] dark:text-[#c4c7c5]"
              >
                <div class="flex flex-col items-center justify-center gap-2">
                  <Layers class="w-8 h-8 text-[#b0b8c4] dark:text-[#555]" />
                  <span class="font-medium text-xs"
                    >Belum ada buku tarif yang ditemukan</span
                  >
                  <p class="text-[11px] text-[#747775] dark:text-[#c4c7c5]">
                    Mulai dengan membuat Buku Tarif standar RS baru.
                  </p>
                </div>
              </td>
            </tr>
          {:else}
            {#each plans || [] as plan, idx (plan.id)}
              {@const badge =
                statusBadgeConfig[plan.status] || statusBadgeConfig.DRAFT}
              <tr
                class="hover:bg-[#f8fafd] dark:hover:bg-[#282a2c] transition-colors"
              >
                <td
                  class="py-3 px-4 text-center font-mono text-[#747775] dark:text-[#c4c7c5]"
                >
                  {(currentPage - 1) * 10 + idx + 1}
                </td>
                <td class="py-3 px-4">
                  <div class="flex flex-col">
                    <span
                      class="font-mono font-bold text-xs text-[#0b57d0] dark:text-blue-400 bg-[#e8f0fe] dark:bg-[#0b57d0]/20 px-2 py-0.5 rounded-md inline-block w-fit text-nowrap"
                    >
                      {plan.code}
                    </span>
                    {#if plan.is_default}
                      <span
                        class="mt-1 text-[10px] font-semibold text-emerald-700 dark:text-emerald-400 bg-emerald-50 dark:bg-emerald-950/40 border border-emerald-200 dark:border-emerald-800 px-1.5 py-0.2 rounded w-fit"
                      >
                        DEFAULT RS
                      </span>
                    {/if}
                  </div>
                </td>
                <td
                  class="py-3 px-4 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]"
                >
                  <div>{plan.name}</div>
                  {#if plan.description}
                    <div
                      class="text-[11px] text-[#747775] dark:text-[#c4c7c5] font-normal truncate max-w-xs"
                    >
                      {plan.description}
                    </div>
                  {/if}
                </td>
                <td class="py-3 px-4 text-[#444746] dark:text-[#c4c7c5]">
                  <div class="flex items-center gap-1.5 text-nowrap">
                    <Calendar
                      class="w-3.5 h-3.5 text-[#747775] dark:text-[#c4c7c5] shrink-0"
                    />
                    <span>{plan.effective_from}</span>
                    <span class="text-[#747775] dark:text-[#c4c7c5]">s/d</span>
                    <span>{plan.effective_to || "Permanen"}</span>
                  </div>
                </td>
                <td class="py-3 px-4">
                  {#if plan.customer_name}
                    <div
                      class="flex items-center gap-1 text-[#1f1f1f] dark:text-[#e3e3e3] font-medium"
                    >
                      <Building2
                        class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400 shrink-0"
                      />
                      <span class="truncate max-w-[120px]"
                        >{plan.customer_name}</span
                      >
                    </div>
                  {:else}
                    <span
                      class="text-[11px] text-[#747775] dark:text-[#c4c7c5] italic"
                      >Semua Pasien / Umum</span
                    >
                  {/if}
                </td>
                <td
                  class="py-3 px-4 text-center font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3]"
                >
                  +{plan.default_cito_percent}%
                </td>
                <td class="py-3 px-4 text-center">
                  <span
                    class="px-2 py-0.5 rounded-full text-[10px] font-bold border {badge.bg} {badge.text} {badge.border}"
                  >
                    {badge.label}
                  </span>
                </td>
                <td class="py-3 px-4 text-center">
                  <button
                    type="button"
                    onclick={() => onManageItems(plan)}
                    class="px-2 py-1 rounded-lg bg-[#e8f0fe] dark:bg-[#0b57d0]/20 hover:bg-[#d2e3fc] dark:hover:bg-[#0b57d0]/30 text-[#0b57d0] dark:text-blue-400 text-xs font-bold flex items-center gap-1 mx-auto cursor-pointer"
                    title="Buka rincian item tarif"
                  >
                    <List class="w-3 h-3" />
                    <span>{plan.total_items} Item</span>
                  </button>
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-1 flex-wrap">
                    <!-- Kelola Item -->
                    <button
                      type="button"
                      onclick={() => onManageItems(plan)}
                      class="p-1.5 rounded-lg text-[#0b57d0] dark:text-blue-400 hover:bg-[#e8f0fe] dark:hover:bg-[#0b57d0]/20 transition-colors cursor-pointer"
                      title="Kelola Item Tarif & Komponen Biaya"
                    >
                      <List class="w-3.5 h-3.5" />
                    </button>

                    <!-- Edit Metadata (Hanya DRAFT) -->
                    {#if plan.status === "DRAFT"}
                      <button
                        type="button"
                        onclick={() => onEdit(plan)}
                        class="p-1.5 rounded-lg text-[#0b57d0] dark:text-blue-400 hover:bg-[#e8f0fe] dark:hover:bg-[#0b57d0]/20 transition-colors cursor-pointer"
                        title="Edit Metadata Buku Tarif"
                      >
                        <Edit3 class="w-3.5 h-3.5" />
                      </button>
                    {/if}

                    <!-- Clone Plan -->
                    <button
                      type="button"
                      onclick={() => onClone(plan)}
                      class="p-1.5 rounded-lg text-[#747775] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3] transition-colors cursor-pointer"
                      title="Gandakan (Clone) ke Buku Tarif Draf Baru"
                    >
                      <Copy class="w-3.5 h-3.5" />
                    </button>

                    <!-- STATE MACHINE ACTIONS -->
                    {#if plan.status === "DRAFT"}
                      <button
                        type="button"
                        onclick={() => onSubmitPlan(plan)}
                        class="px-2 py-1 rounded-lg bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-300 hover:bg-blue-100 dark:hover:bg-blue-900/60 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Ajukan Persetujuan (Submit)"
                      >
                        <Send class="w-3 h-3" />
                        <span>Submit</span>
                      </button>
                    {:else if plan.status === "SUBMITTED"}
                      <button
                        type="button"
                        onclick={() => onApprovePlan(plan)}
                        class="px-2 py-1 rounded-lg bg-purple-50 dark:bg-purple-950/40 text-purple-700 dark:text-purple-300 hover:bg-purple-100 dark:hover:bg-purple-900/60 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Setujui (Approve) & Validasi Balancing"
                      >
                        <CheckCircle class="w-3 h-3" />
                        <span>Approve</span>
                      </button>
                    {:else if plan.status === "APPROVED"}
                      <button
                        type="button"
                        onclick={() => onActivatePlan(plan)}
                        class="px-2 py-1 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 hover:bg-emerald-100 dark:hover:bg-emerald-900/60 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Aktifkan (Activate) sebagai Acuan Kasir"
                      >
                        <Power class="w-3 h-3" />
                        <span>Aktifkan</span>
                      </button>
                    {:else if plan.status === "ACTIVE"}
                      <button
                        type="button"
                        onclick={() => onArchivePlan(plan)}
                        class="px-2 py-1 rounded-lg bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Arsipkan (Archive)"
                      >
                        <Archive class="w-3 h-3" />
                        <span>Arsipkan</span>
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

    <!-- Pagination Footer -->
    {#if totalPages > 1}
      <div
        class="p-3 border-t border-[#e1e5ea] dark:border-[#444746] bg-[#f8fafd] dark:bg-[#282a2c] flex items-center justify-between text-xs text-[#444746] dark:text-[#c4c7c5]"
      >
        <span
          >Total <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]"
            >{totalCount}</strong
          > buku tarif</span
        >
        <div class="flex items-center gap-1">
          <button
            type="button"
            disabled={currentPage <= 1 || isLoading}
            onclick={() => onPageChange(currentPage - 1)}
            class="p-1.5 rounded-lg border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#1e1f20] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c] disabled:opacity-40 cursor-pointer"
          >
            <ChevronLeft class="w-3.5 h-3.5" />
          </button>
          <span class="px-2 font-medium">Hal {currentPage} / {totalPages}</span>
          <button
            type="button"
            disabled={currentPage >= totalPages || isLoading}
            onclick={() => onPageChange(currentPage + 1)}
            class="p-1.5 rounded-lg border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#1e1f20] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c] disabled:opacity-40 cursor-pointer"
          >
            <ChevronRight class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>
