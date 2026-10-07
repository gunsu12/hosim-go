<script lang="ts">
  import {
    Pill,
    Package,
    Stethoscope,
    Receipt,
    Boxes,
    CheckCircle2,
    XCircle,
    ChevronRight,
    Eye,
    Edit3,
    Trash2,
    RefreshCw,
    Plus
  } from '@lucide/svelte';
  import M3Button from '../../../../components/m3/M3Button.svelte';
  import type { ItemSummaryRecord, ItemType } from '../../../../types/master/item';

  interface Props {
    items: ItemSummaryRecord[];
    isLoading?: boolean;
    currentPage: number;
    totalPages: number;
    totalItems: number;
    onPageChange: (page: number) => void;
    onViewDetail: (item: ItemSummaryRecord) => void;
    onEdit: (item: ItemSummaryRecord) => void;
    onDelete: (item: ItemSummaryRecord) => void;
    onCreate: () => void;
  }

  let {
    items,
    isLoading = false,
    currentPage,
    totalPages,
    totalItems,
    onPageChange,
    onViewDetail,
    onEdit,
    onDelete,
    onCreate
  }: Props = $props();

  function getTypeBadge(type: ItemType) {
    switch (type) {
      case 'MEDICATION':
        return { label: 'Obat', bg: 'bg-emerald-50 text-emerald-700 border-emerald-200', icon: Pill };
      case 'GENERAL':
        return { label: 'BMHP / Umum', bg: 'bg-amber-50 text-amber-700 border-amber-200', icon: Package };
      case 'ASSET':
        return { label: 'Aset / Alkes', bg: 'bg-indigo-50 text-indigo-700 border-indigo-200', icon: Stethoscope };
      case 'TARIFF':
        return { label: 'Tarif Jasa', bg: 'bg-blue-50 text-blue-700 border-blue-200', icon: Receipt };
    }
  }
</script>

<div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden">
  {#if isLoading}
    <div class="p-12 flex flex-col items-center justify-center gap-3 text-center">
      <RefreshCw class="w-8 h-8 text-[#0b57d0] animate-spin" />
      <p class="text-xs text-[#747775]">Memuat katalog item rumah sakit...</p>
    </div>
  {:else if items.length === 0}
    <div class="p-12 flex flex-col items-center justify-center gap-3 text-center">
      <div class="w-12 h-12 rounded-full bg-[#f0f4f9] text-[#747775] flex items-center justify-center">
        <Boxes class="w-6 h-6" />
      </div>
      <div>
        <h3 class="text-sm font-semibold text-[#1f1f1f]">Katalog Item Masih Kosong</h3>
        <p class="text-xs text-[#747775] mt-1 max-w-sm">
          Belum ada data item yang sesuai dengan kriteria filter. Tambahkan master obat, BMHP, aset, atau tarif baru.
        </p>
      </div>
      <M3Button variant="tonal" onclick={onCreate}>
        <Plus class="w-4 h-4" />
        <span>Tambah Item Sekarang</span>
      </M3Button>
    </div>
  {:else}
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746] font-semibold">
            <th class="py-3 px-4">Kode Item / SKU</th>
            <th class="py-3 px-4">Nama Produk / Layanan</th>
            <th class="py-3 px-4">Tipe & Kategori</th>
            <th class="py-3 px-4">Satuan Dasar (Base UOM)</th>
            <th class="py-3 px-4">Lini Produk / COA</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#each items as item (item.id)}
            {@const badge = getTypeBadge(item.item_type)}
            <tr class="hover:bg-[#f8fafd] transition-colors">
              <!-- Kode Item -->
              <td class="py-3.5 px-4 font-mono font-bold text-[#0b57d0] whitespace-nowrap">
                <button
                  type="button"
                  onclick={() => onViewDetail(item)}
                  class="hover:underline flex items-center gap-1 cursor-pointer"
                >
                  <span>{item.code}</span>
                  <ChevronRight class="w-3 h-3 opacity-60" />
                </button>
              </td>

              <!-- Nama Item -->
              <td class="py-3.5 px-4 font-semibold text-[#1f1f1f] max-w-[260px]">
                <div class="truncate text-xs font-semibold">{item.name}</div>
                {#if item.generic_name}
                  <div class="text-[11px] text-[#747775] font-normal truncate">
                    Generik: {item.generic_name}
                  </div>
                {/if}
              </td>

              <!-- Tipe & Kategori -->
              <td class="py-3.5 px-4 whitespace-nowrap">
                <div class="flex items-center gap-1.5">
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium border {badge.bg}">
                    <badge.icon class="w-3 h-3" />
                    <span>{badge.label}</span>
                  </span>
                </div>
                {#if item.category}
                  <div class="text-[11px] text-[#444746] mt-0.5 truncate max-w-[180px]">
                    {item.category.name}
                  </div>
                {/if}
              </td>

              <!-- Satuan Dasar -->
              <td class="py-3.5 px-4 whitespace-nowrap font-medium text-[#444746]">
                <span class="px-2 py-0.5 rounded bg-[#f0f4f9] text-[#1f1f1f] font-mono text-[11px]">
                  {item.base_unit?.unit_name || 'Unit'}
                </span>
              </td>

              <!-- Lini Produk -->
              <td class="py-3.5 px-4 text-[#444746] whitespace-nowrap">
                {#if item.product_line}
                  <div class="font-medium text-xs">{item.product_line.name}</div>
                  {#if item.product_line.inventory_coa_code}
                    <div class="text-[10px] text-[#747775] font-mono">
                      COA: {item.product_line.inventory_coa_code}
                    </div>
                  {/if}
                {:else if item.item_type === 'TARIFF'}
                  <span class="text-[11px] text-[#747775] italic">Non-Persediaan (Jasa)</span>
                {:else}
                  <span class="text-[11px] text-[#747775]">-</span>
                {/if}
              </td>

              <!-- Status -->
              <td class="py-3.5 px-4 whitespace-nowrap">
                {#if item.is_active}
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-[#e6f4ea] text-[#137333]">
                    <CheckCircle2 class="w-3 h-3" />
                    <span>Aktif</span>
                  </span>
                {:else}
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-[#fce8e6] text-[#c5221f]">
                    <XCircle class="w-3 h-3" />
                    <span>Nonaktif</span>
                  </span>
                {/if}
              </td>

              <!-- Aksi -->
              <td class="py-3.5 px-4 text-right whitespace-nowrap">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    onclick={() => onViewDetail(item)}
                    class="p-1.5 rounded-lg text-[#0b57d0] hover:bg-[#0b57d0]/10 transition-colors cursor-pointer"
                    title="Lihat Rincian Lengkap"
                  >
                    <Eye class="w-4 h-4" />
                  </button>

                  <button
                    type="button"
                    onclick={() => onEdit(item)}
                    class="p-1.5 rounded-lg text-[#444746] hover:bg-[#f0f4f9] transition-colors cursor-pointer"
                    title="Edit Item"
                  >
                    <Edit3 class="w-4 h-4" />
                  </button>

                  <button
                    type="button"
                    onclick={() => onDelete(item)}
                    class="p-1.5 rounded-lg text-[#c5221f] hover:bg-[#c5221f]/10 transition-colors cursor-pointer"
                    title="Hapus Item"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Pagination Controls -->
    {#if totalPages > 1}
      <div class="flex items-center justify-between p-4 border-t border-[#e1e5ea] bg-[#f8fafd]">
        <div class="text-xs text-[#747775]">
          Menampilkan {items.length} dari total {totalItems} item
        </div>
        <div class="flex items-center gap-1.5">
          <button
            type="button"
            onclick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1}
            class="px-3 py-1.5 rounded-lg border border-[#e1e5ea] bg-white text-xs text-[#1f1f1f] disabled:opacity-40 hover:bg-[#f0f4f9] cursor-pointer"
          >
            Sebelumnya
          </button>
          <span class="px-2 text-xs font-mono font-semibold">{currentPage} / {totalPages}</span>
          <button
            type="button"
            onclick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            class="px-3 py-1.5 rounded-lg border border-[#e1e5ea] bg-white text-xs text-[#1f1f1f] disabled:opacity-40 hover:bg-[#f0f4f9] cursor-pointer"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>
