<script lang="ts">
  import { Trash2 } from '@lucide/svelte';
  import type { ItemSummaryRecord } from '../../../../types/master/item';

  interface Props {
    isOpen: boolean;
    item: ItemSummaryRecord | null;
    isLoading?: boolean;
    onClose: () => void;
    onConfirm: (item: ItemSummaryRecord) => void;
  }

  let {
    isOpen,
    item,
    isLoading = false,
    onClose,
    onConfirm
  }: Props = $props();
</script>

{#if isOpen && item}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-md p-6 flex flex-col gap-4">
      <div class="w-12 h-12 rounded-2xl bg-rose-50 text-rose-600 flex items-center justify-center mx-auto">
        <Trash2 class="w-6 h-6" />
      </div>
      <div class="text-center">
        <h3 class="text-sm font-bold text-[#1f1f1f]">Hapus Item dari Katalog?</h3>
        <p class="text-xs text-[#747775] mt-1">
          Anda akan menghapus item <strong class="text-[#1f1f1f]">{item.name}</strong> ({item.code}). Riwayat kartu stok dan penagihan masa lalu tetap terlindungi secara aman (soft-delete).
        </p>
      </div>

      <div class="flex items-center justify-center gap-2 pt-2">
        <button
          onclick={onClose}
          type="button"
          disabled={isLoading}
          class="px-4 py-2 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] hover:bg-[#f0f4f9] cursor-pointer"
        >
          Batal
        </button>
        <button
          onclick={() => onConfirm(item)}
          type="button"
          disabled={isLoading}
          class="px-5 py-2 rounded-xl bg-[#c5221f] text-white text-xs font-semibold hover:bg-[#a51d24] transition-colors cursor-pointer shadow-xs disabled:opacity-50"
        >
          {isLoading ? 'Menghapus...' : 'Ya, Hapus Item'}
        </button>
      </div>
    </div>
  </div>
{/if}
