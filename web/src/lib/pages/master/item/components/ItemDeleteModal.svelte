<script lang="ts">
  import { Trash2 } from "@lucide/svelte";
  import M3Modal from "../../../../components/m3/M3Modal.svelte";
  import M3Button from "../../../../components/m3/M3Button.svelte";
  import type { ItemSummaryRecord } from "../../../../types/master/item";

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
    onConfirm,
  }: Props = $props();
</script>

<M3Modal {isOpen} title="Hapus Item dari Katalog?" size="md" {onClose}>
  {#if item}
    <div class="flex flex-col gap-4 text-center">
      <div
        class="w-12 h-12 rounded-2xl bg-rose-50 dark:bg-rose-950/40 text-rose-600 dark:text-rose-400 flex items-center justify-center mx-auto"
      >
        <Trash2 class="w-6 h-6" />
      </div>
      <div>
        <p class="text-xs text-[#444746] dark:text-[#c4c7c5] leading-relaxed">
          Anda akan menghapus item <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]">{item.name}</strong> ({item.code}).
          Riwayat kartu stok dan penagihan masa lalu tetap terlindungi secara aman (soft-delete).
        </p>
      </div>

      <div class="flex items-center justify-center gap-2 pt-2">
        <M3Button variant="text" onclick={onClose} disabled={isLoading}>
          Batal
        </M3Button>
        <button
          onclick={() => onConfirm(item)}
          type="button"
          disabled={isLoading}
          class="px-5 py-2 rounded-xl bg-[#c5221f] text-white text-xs font-semibold hover:bg-[#a51d24] transition-colors cursor-pointer shadow-xs disabled:opacity-50"
        >
          {isLoading ? "Menghapus..." : "Ya, Hapus Item"}
        </button>
      </div>
    </div>
  {/if}
</M3Modal>
