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
    Clock,
    AlertCircle
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type {
    PricePlanRecord,
    PricePlanStatus
  } from '$lib/types/finance/price_plan';

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
    onArchivePlan
  }: Props = $props();

  let localSearch = $state('');

  $effect(() => {
    localSearch = searchQuery;
  });

  const statusBadgeConfig: Record<PricePlanStatus, { label: string; bg: string; text: string; border: string }> = {
    DRAFT: { label: 'DRAFT', bg: 'bg-amber-50', text: 'text-amber-800', border: 'border-amber-200' },
    SUBMITTED: { label: 'SUBMITTED', bg: 'bg-blue-50', text: 'text-blue-800', border: 'border-blue-200' },
    APPROVED: { label: 'APPROVED', bg: 'bg-purple-50', text: 'text-purple-800', border: 'border-purple-200' },
    ACTIVE: { label: 'ACTIVE', bg: 'bg-emerald-50', text: 'text-emerald-800', border: 'border-emerald-200' },
    ARCHIVED: { label: 'ARCHIVED', bg: 'bg-slate-100', text: 'text-slate-600', border: 'border-slate-300' }
  };
</script>

<div class="flex flex-col gap-4">
  <!-- Controls Bar -->
  <div class="flex flex-wrap items-center justify-between gap-3 bg-[#f8fafd] p-3 rounded-2xl border border-[#e1e5ea]">
    <div class="flex-1 min-w-[240px] relative">
      <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={localSearch}
        onkeydown={(e) => e.key === 'Enter' && onSearchChange(localSearch)}
        placeholder="Cari kode atau nama buku tarif (tekan Enter)..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-white border border-[#e1e5ea] focus:border-[#0b57d0] text-xs text-[#1f1f1f] focus:outline-hidden transition-all shadow-2xs"
      />
    </div>

    <div class="flex items-center gap-2 flex-wrap">
      <!-- Status Filter -->
      <select
        value={filterStatus}
        onchange={(e) => onStatusChange(e.currentTarget.value)}
        class="h-10 px-3 rounded-xl border border-[#e1e5ea] bg-white text-xs text-[#444746] focus:border-[#0b57d0] focus:outline-hidden cursor-pointer shadow-2xs"
      >
        <option value="ALL">Semua Status</option>
        <option value="DRAFT">DRAFT (Penyusunan)</option>
        <option value="SUBMITTED">SUBMITTED (Diajukan)</option>
        <option value="APPROVED">APPROVED (Disetujui)</option>
        <option value="ACTIVE">ACTIVE (Aktif Digunakan)</option>
        <option value="ARCHIVED">ARCHIVED (Diarsipkan)</option>
      </select>

      <M3Button variant="filled" onclick={onAdd}>
        <Plus class="w-4 h-4" />
        <span>Buat Buku Tarif</span>
      </M3Button>
    </div>
  </div>

  <!-- Table Container -->
  <div class="border border-[#e1e5ea] rounded-2xl bg-white overflow-hidden shadow-2xs">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="bg-[#f0f4f9] text-[#444746] font-semibold border-b border-[#e1e5ea]">
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
        <tbody class="divide-y divide-[#e1e5ea]">
          {#if isLoading}
            <tr>
              <td colspan="9" class="py-12 text-center text-[#747775]">
                <div class="flex flex-col items-center justify-center gap-2">
                  <div class="w-6 h-6 border-2 border-[#0b57d0] border-t-transparent rounded-full animate-spin"></div>
                  <span>Memuat daftar buku tarif...</span>
                </div>
              </td>
            </tr>
          {:else if !plans || plans.length === 0}
            <tr>
              <td colspan="9" class="py-12 text-center text-[#747775]">
                <div class="flex flex-col items-center justify-center gap-2">
                  <Layers class="w-8 h-8 text-[#b0b8c4]" />
                  <span class="font-medium text-xs">Belum ada buku tarif yang ditemukan</span>
                  <p class="text-[11px] text-[#747775]">Mulai dengan membuat Buku Tarif standar RS baru.</p>
                </div>
              </td>
            </tr>
          {:else}
            {#each plans || [] as plan, idx (plan.id)}
              {@const badge = statusBadgeConfig[plan.status] || statusBadgeConfig.DRAFT}
              <tr class="hover:bg-[#f8fafd] transition-colors">
                <td class="py-3 px-4 text-center font-mono text-[#747775]">
                  {(currentPage - 1) * 10 + idx + 1}
                </td>
                <td class="py-3 px-4">
                  <div class="flex flex-col">
                    <span class="font-mono font-bold text-xs text-[#0b57d0] bg-[#e8f0fe] px-2 py-0.5 rounded-md inline-block w-fit">
                      {plan.code}
                    </span>
                    {#if plan.is_default}
                      <span class="mt-1 text-[10px] font-semibold text-emerald-700 bg-emerald-50 border border-emerald-200 px-1.5 py-0.2 rounded w-fit">
                        DEFAULT RS
                      </span>
                    {/if}
                  </div>
                </td>
                <td class="py-3 px-4 font-semibold text-[#1f1f1f]">
                  <div>{plan.name}</div>
                  {#if plan.description}
                    <div class="text-[11px] text-[#747775] font-normal truncate max-w-xs">{plan.description}</div>
                  {/if}
                </td>
                <td class="py-3 px-4 text-[#444746]">
                  <div class="flex items-center gap-1.5">
                    <Calendar class="w-3.5 h-3.5 text-[#747775] shrink-0" />
                    <span>{plan.effective_from}</span>
                    <span class="text-[#747775]">s/d</span>
                    <span>{plan.effective_to || 'Permanen'}</span>
                  </div>
                </td>
                <td class="py-3 px-4">
                  {#if plan.customer_name}
                    <div class="flex items-center gap-1 text-[#1f1f1f] font-medium">
                      <Building2 class="w-3.5 h-3.5 text-blue-600 shrink-0" />
                      <span class="truncate max-w-[120px]">{plan.customer_name}</span>
                    </div>
                  {:else}
                    <span class="text-[11px] text-[#747775] italic">Semua Pasien / Umum</span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-center font-mono font-bold text-[#1f1f1f]">
                  +{plan.default_cito_percent}%
                </td>
                <td class="py-3 px-4 text-center">
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-bold border {badge.bg} {badge.text} {badge.border}">
                    {badge.label}
                  </span>
                </td>
                <td class="py-3 px-4 text-center">
                  <button
                    type="button"
                    onclick={() => onManageItems(plan)}
                    class="px-2 py-1 rounded-lg bg-[#e8f0fe] hover:bg-[#d2e3fc] text-[#0b57d0] text-xs font-bold flex items-center gap-1 mx-auto cursor-pointer"
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
                      class="p-1.5 rounded-lg text-[#0b57d0] hover:bg-[#e8f0fe] transition-colors cursor-pointer"
                      title="Kelola Item Tarif & Komponen Biaya"
                    >
                      <List class="w-3.5 h-3.5" />
                    </button>

                    <!-- Edit Metadata (Hanya DRAFT) -->
                    {#if plan.status === 'DRAFT'}
                      <button
                        type="button"
                        onclick={() => onEdit(plan)}
                        class="p-1.5 rounded-lg text-[#0b57d0] hover:bg-[#e8f0fe] transition-colors cursor-pointer"
                        title="Edit Metadata Buku Tarif"
                      >
                        <Edit3 class="w-3.5 h-3.5" />
                      </button>
                    {/if}

                    <!-- Clone Plan -->
                    <button
                      type="button"
                      onclick={() => onClone(plan)}
                      class="p-1.5 rounded-lg text-[#747775] hover:bg-[#f0f4f9] hover:text-[#1f1f1f] transition-colors cursor-pointer"
                      title="Gandakan (Clone) ke Buku Tarif Draf Baru"
                    >
                      <Copy class="w-3.5 h-3.5" />
                    </button>

                    <!-- STATE MACHINE ACTIONS -->
                    {#if plan.status === 'DRAFT'}
                      <button
                        type="button"
                        onclick={() => onSubmitPlan(plan)}
                        class="px-2 py-1 rounded-lg bg-blue-50 text-blue-700 hover:bg-blue-100 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Ajukan Persetujuan (Submit)"
                      >
                        <Send class="w-3 h-3" />
                        <span>Submit</span>
                      </button>
                    {:else if plan.status === 'SUBMITTED'}
                      <button
                        type="button"
                        onclick={() => onApprovePlan(plan)}
                        class="px-2 py-1 rounded-lg bg-purple-50 text-purple-700 hover:bg-purple-100 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Setujui (Approve) & Validasi Balancing"
                      >
                        <CheckCircle class="w-3 h-3" />
                        <span>Approve</span>
                      </button>
                    {:else if plan.status === 'APPROVED'}
                      <button
                        type="button"
                        onclick={() => onActivatePlan(plan)}
                        class="px-2 py-1 rounded-lg bg-emerald-50 text-emerald-700 hover:bg-emerald-100 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
                        title="Aktifkan (Activate) sebagai Acuan Kasir"
                      >
                        <Power class="w-3 h-3" />
                        <span>Aktifkan</span>
                      </button>
                    {:else if plan.status === 'ACTIVE'}
                      <button
                        type="button"
                        onclick={() => onArchivePlan(plan)}
                        class="px-2 py-1 rounded-lg bg-slate-100 text-slate-700 hover:bg-slate-200 font-semibold text-[10px] flex items-center gap-1 cursor-pointer"
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
      <div class="p-3 border-t border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between text-xs text-[#444746]">
        <span>Total <strong>{totalCount}</strong> buku tarif</span>
        <div class="flex items-center gap-1">
          <button
            type="button"
            disabled={currentPage <= 1 || isLoading}
            onclick={() => onPageChange(currentPage - 1)}
            class="p-1.5 rounded-lg border border-[#e1e5ea] bg-white hover:bg-[#f0f4f9] disabled:opacity-40 cursor-pointer"
          >
            <ChevronLeft class="w-3.5 h-3.5" />
          </button>
          <span class="px-2 font-medium">Hal {currentPage} / {totalPages}</span>
          <button
            type="button"
            disabled={currentPage >= totalPages || isLoading}
            onclick={() => onPageChange(currentPage + 1)}
            class="p-1.5 rounded-lg border border-[#e1e5ea] bg-white hover:bg-[#f0f4f9] disabled:opacity-40 cursor-pointer"
          >
            <ChevronRight class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>
