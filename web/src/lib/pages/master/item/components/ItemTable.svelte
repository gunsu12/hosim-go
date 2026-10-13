<script lang="ts">
  import {
    Pill,
    Package,
    Stethoscope,
    Receipt,
    Boxes,
    CheckCircle2,
    XCircle,
    Eye,
    Edit3,
    Trash2,
    RefreshCw,
    Plus,
  } from "@lucide/svelte";
  import M3Button from "../../../../components/m3/M3Button.svelte";
  import type {
    ItemSummaryRecord,
    ItemType,
  } from "../../../../types/master/item";

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
    onCreate,
  }: Props = $props();

  function getTypeBadge(type: ItemType) {
    switch (type) {
      case "MEDICATION":
        return {
          label: "Obat",
          bg: "bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border-emerald-200 dark:border-emerald-800",
          icon: Pill,
        };
      case "GENERAL":
        return {
          label: "BMHP / Umum",
          bg: "bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-400 border-amber-200 dark:border-amber-800",
          icon: Package,
        };
      case "ASSET":
        return {
          label: "Aset / Alkes",
          bg: "bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-400 border-indigo-200 dark:border-indigo-800",
          icon: Stethoscope,
        };
      case "TARIFF":
        return {
          label: "Tarif Jasa",
          bg: "bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-400 border-blue-200 dark:border-blue-800",
          icon: Receipt,
        };
    }
  }
</script>

<div
  class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#444746] shadow-xs overflow-hidden"
>
  {#if isLoading}
    <div
      class="p-12 flex flex-col items-center justify-center gap-3 text-center"
    >
      <RefreshCw class="w-8 h-8 text-[#0b57d0] animate-spin" />
      <p class="text-xs text-[#747775] dark:text-[#c4c7c5]">
        Memuat katalog item rumah sakit...
      </p>
    </div>
  {:else if items.length === 0}
    <div
      class="p-12 flex flex-col items-center justify-center gap-3 text-center"
    >
      <div
        class="w-12 h-12 rounded-full bg-[#f0f4f9] dark:bg-[#282a2c] text-[#747775] dark:text-[#c4c7c5] flex items-center justify-center"
      >
        <Boxes class="w-6 h-6" />
      </div>
      <div>
        <h3 class="text-sm font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
          Katalog Item Masih Kosong
        </h3>
        <p class="text-xs text-[#747775] dark:text-[#c4c7c5] mt-1 max-w-sm">
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
          <tr
            class="bg-[#f0f4f9] dark:bg-[#282a2c] border-b border-[#e1e5ea] dark:border-[#444746] text-[#444746] dark:text-[#c4c7c5] font-semibold"
          >
            <th class="py-3 px-4">Kode Item / SKU</th>
            <th class="py-3 px-4">Nama Produk / Layanan</th>
            <th class="py-3 px-4">Tipe & Kategori</th>
            <th class="py-3 px-4">Satuan Dasar (Base UOM)</th>
            <th class="py-3 px-4">Lini Produk / COA</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#each items as item (item.id)}
            {@const badge = getTypeBadge(item.item_type)}
            <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#2b2b2b] transition-colors">
              <!-- Kode Item -->
              <td class="py-3 px-4 font-mono font-bold text-[#0b57d0] dark:text-blue-400">
                {item.code}
              </td>

              <!-- Nama Produk -->
              <td class="py-3 px-4">
                <div class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{item.name}</div>
                {#if item.generic_name}
                  <div class="text-[11px] text-[#747775] dark:text-[#c4c7c5] mt-0.5">
                    {item.generic_name}
                  </div>
                {/if}
              </td>

              <!-- Tipe & Kategori -->
              <td class="py-3 px-4">
                <div class="flex items-center gap-1.5 flex-wrap">
                  <span
                    class="px-2 py-0.5 rounded-full font-semibold border flex items-center gap-1 text-[11px] {badge.bg}"
                  >
                    <badge.icon class="w-3 h-3" />
                    <span>{badge.label}</span>
                  </span>
                  {#if item.category}
                    <span class="text-xs text-[#444746] dark:text-[#c4c7c5]">{item.category.name}</span>
                  {/if}
                </div>
              </td>

              <!-- Satuan Dasar -->
              <td class="py-3 px-4 font-mono">
                <span
                  class="px-2 py-0.5 rounded-lg bg-[#f0f4f9] dark:bg-[#282a2c] text-[#1f1f1f] dark:text-[#e3e3e3] border border-[#e1e5ea] dark:border-[#444746]"
                >
                  {item.base_unit?.unit_name || item.base_unit_name || '-'}
                </span>
              </td>

              <!-- Lini Produk / COA -->
              <td class="py-3 px-4">
                {#if item.product_line}
                  <div class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">
                    {item.product_line.name}
                  </div>
                  <div class="text-[10px] text-[#747775] dark:text-[#c4c7c5] font-mono mt-0.5">
                    Persediaan: {item.product_line.inventory_coa_code || "-"}
                  </div>
                {:else}
                  <span class="text-[#747775] dark:text-[#c4c7c5] italic">Non-Persediaan</span>
                {/if}
              </td>

              <!-- Status -->
              <td class="py-3 px-4">
                {#if item.is_active}
                  <span
                    class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800"
                  >
                    <CheckCircle2 class="w-3 h-3" />
                    <span>Aktif</span>
                  </span>
                {:else}
                  <span
                    class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-semibold bg-rose-50 dark:bg-rose-950/40 text-rose-700 dark:text-rose-400 border border-rose-200 dark:border-rose-800"
                  >
                    <XCircle class="w-3 h-3" />
                    <span>Nonaktif</span>
                  </span>
                {/if}
              </td>

              <!-- Aksi -->
              <td class="py-3 px-4 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    type="button"
                    onclick={() => onViewDetail(item)}
                    class="p-1.5 rounded-lg text-[#0b57d0] dark:text-blue-400 hover:bg-[#0b57d0]/10 transition-colors cursor-pointer"
                    title="Lihat Rincian Lengkap"
                  >
                    <Eye class="w-4 h-4" />
                  </button>

                  <button
                    type="button"
                    onclick={() => onEdit(item)}
                    class="p-1.5 rounded-lg text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c] transition-colors cursor-pointer"
                    title="Edit Item"
                  >
                    <Edit3 class="w-4 h-4" />
                  </button>

                  <button
                    type="button"
                    onclick={() => onDelete(item)}
                    class="p-1.5 rounded-lg text-[#c5221f] dark:text-rose-400 hover:bg-[#c5221f]/10 transition-colors cursor-pointer"
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
      <div
        class="flex items-center justify-between p-4 border-t border-[#e1e5ea] dark:border-[#444746] bg-[#f8fafd] dark:bg-[#282a2c]"
      >
        <div class="text-xs text-[#747775] dark:text-[#c4c7c5]">
          Menampilkan {items.length} dari total {totalItems} item
        </div>
        <div class="flex items-center gap-1.5">
          <button
            type="button"
            onclick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1}
            class="px-3 py-1.5 rounded-lg border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] disabled:opacity-40 hover:bg-[#f0f4f9] dark:hover:bg-[#383b3d] cursor-pointer"
          >
            Sebelumnya
          </button>
          <span class="px-2 text-xs font-mono font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
            {currentPage} / {totalPages}
          </span>
          <button
            type="button"
            onclick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            class="px-3 py-1.5 rounded-lg border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] disabled:opacity-40 hover:bg-[#f0f4f9] dark:hover:bg-[#383b3d] cursor-pointer"
          >
            Selanjutnya
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>
