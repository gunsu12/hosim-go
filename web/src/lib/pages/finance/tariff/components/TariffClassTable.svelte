<script lang="ts">
  import {
    Search,
    Plus,
    Edit3,
    Trash2,
    CheckCircle2,
    XCircle,
    ChevronLeft,
    ChevronRight,
    Layers,
    Tag
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { TariffClassRecord } from '$lib/types/finance/tariff';

  interface Props {
    classes: TariffClassRecord[];
    isLoading: boolean;
    searchQuery: string;
    filterStatus: 'ALL' | 'ACTIVE' | 'INACTIVE';
    currentPage: number;
    totalPages: number;
    totalCount: number;
    onSearchChange: (q: string) => void;
    onStatusChange: (s: 'ALL' | 'ACTIVE' | 'INACTIVE') => void;
    onPageChange: (p: number) => void;
    onAdd: () => void;
    onEdit: (tc: TariffClassRecord) => void;
    onDelete: (tc: TariffClassRecord) => void;
  }

  let {
    classes,
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
    onDelete
  }: Props = $props();

  let localSearch = $state('');

  $effect(() => {
    localSearch = searchQuery;
  });
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
        placeholder="Cari kode kelas atau nama kelas tarif (tekan Enter)..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-white border border-[#e1e5ea] focus:border-[#0b57d0] text-xs text-[#1f1f1f] focus:outline-hidden transition-all shadow-2xs"
      />
    </div>

    <div class="flex items-center gap-2">
      <!-- Status Filter -->
      <select
        value={filterStatus}
        onchange={(e) => onStatusChange(e.currentTarget.value as any)}
        class="h-10 px-3 rounded-xl border border-[#e1e5ea] bg-white text-xs text-[#444746] focus:border-[#0b57d0] focus:outline-hidden cursor-pointer shadow-2xs"
      >
        <option value="ALL">Semua Status</option>
        <option value="ACTIVE">Hanya Aktif</option>
        <option value="INACTIVE">Hanya Nonaktif</option>
      </select>

      <M3Button variant="filled" onclick={onAdd}>
        <Plus class="w-4 h-4" />
        <span>Tambah Kelas Tarif</span>
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
            <th class="py-3 px-4 w-32">Kode Kelas</th>
            <th class="py-3 px-4">Nama Kelas Perawatan / Tarif</th>
            <th class="py-3 px-4">Keterangan / Fasilitas</th>
            <th class="py-3 px-4 w-28 text-center">Status</th>
            <th class="py-3 px-4 w-24 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea]">
          {#if isLoading}
            <tr>
              <td colspan="6" class="py-12 text-center text-[#747775]">
                <div class="flex flex-col items-center justify-center gap-2">
                  <div class="w-6 h-6 border-2 border-[#0b57d0] border-t-transparent rounded-full animate-spin"></div>
                  <span>Memuat data kelas tarif...</span>
                </div>
              </td>
            </tr>
          {:else if classes.length === 0}
            <tr>
              <td colspan="6" class="py-12 text-center text-[#747775]">
                <div class="flex flex-col items-center justify-center gap-2">
                  <Layers class="w-8 h-8 text-[#b0b8c4]" />
                  <span class="font-medium text-xs">Belum ada kelas tarif yang terdaftar</span>
                  <p class="text-[11px] text-[#747775]">Silakan tambahkan kelas tarif baru seperti VVIP, VIP, Kelas 1, dsb.</p>
                </div>
              </td>
            </tr>
          {:else}
            {#each classes as item, idx (item.id)}
              <tr class="hover:bg-[#f8fafd] transition-colors">
                <td class="py-3 px-4 text-center font-mono text-[#747775]">
                  {(currentPage - 1) * 10 + idx + 1}
                </td>
                <td class="py-3 px-4">
                  <span class="font-mono font-bold text-xs text-[#0b57d0] bg-[#e8f0fe] px-2 py-0.5 rounded-md">
                    {item.code}
                  </span>
                </td>
                <td class="py-3 px-4 font-semibold text-[#1f1f1f]">
                  <div class="flex items-center gap-2">
                    <Tag class="w-3.5 h-3.5 text-[#0b57d0]" />
                    <span>{item.name}</span>
                  </div>
                </td>
                <td class="py-3 px-4 text-[#444746] max-w-xs truncate">
                  {item.description || '-'}
                </td>
                <td class="py-3 px-4 text-center">
                  {#if item.is_active}
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#e6f4ea] text-[#137333] border border-[#ceead6]">
                      <CheckCircle2 class="w-3 h-3" />
                      Aktif
                    </span>
                  {:else}
                    <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[10px] font-semibold bg-[#fce8e6] text-[#c5221f] border border-[#fad2cf]">
                      <XCircle class="w-3 h-3" />
                      Nonaktif
                    </span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <button
                      type="button"
                      onclick={() => onEdit(item)}
                      class="p-1.5 rounded-lg text-[#0b57d0] hover:bg-[#e8f0fe] transition-colors cursor-pointer"
                      title="Edit Kelas Tarif"
                    >
                      <Edit3 class="w-3.5 h-3.5" />
                    </button>
                    <button
                      type="button"
                      onclick={() => onDelete(item)}
                      class="p-1.5 rounded-lg text-[#c5221f] hover:bg-[#fce8e6] transition-colors cursor-pointer"
                      title="Hapus Kelas Tarif"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
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
    {#if totalPages > 1}
      <div class="p-3 border-t border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between text-xs text-[#444746]">
        <span>Total <strong>{totalCount}</strong> kelas tarif</span>
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
