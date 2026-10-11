<script lang="ts">
  import { AlertTriangle, Trash2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';

  interface Props {
    isOpen: boolean;
    title: string;
    itemName: string;
    itemCode: string;
    isDeleting: boolean;
    onClose: () => void;
    onConfirm: () => Promise<void>;
  }

  let {
    isOpen,
    title,
    itemName,
    itemCode,
    isDeleting,
    onClose,
    onConfirm
  }: Props = $props();
</script>

<M3Modal
  {isOpen}
  {title}
  subtitle="Tindakan ini memerlukan konfirmasi otorisasi"
  icon={AlertTriangle}
  iconClass="bg-rose-500/10 text-rose-600 dark:bg-rose-500/20 dark:text-rose-400"
  size="sm"
  {onClose}
>
  <div class="text-xs text-[#444746] dark:text-[#c4c7c5] flex flex-col gap-3">
    <p>Apakah Anda yakin ingin menghapus data tarif berikut?</p>
    
    <div class="p-3.5 rounded-2xl bg-rose-50/70 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-800/60 flex flex-col gap-1">
      <span class="font-mono font-bold text-rose-700 dark:text-rose-400 text-xs">
        [{itemCode}]
      </span>
      <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] text-sm">
        {itemName}
      </span>
    </div>

    <p class="text-[11px] text-[#747775] dark:text-[#8e918f] leading-relaxed">
      Data yang dihapus akan ditandai soft delete di database dan dinonaktifkan dari referensi pelayanan transaksi aktif.
    </p>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isDeleting}>
      Batal
    </M3Button>
    <M3Button
      variant="filled"
      class="bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white dark:bg-rose-600 dark:hover:bg-rose-700"
      onclick={onConfirm}
      disabled={isDeleting}
    >
      {#if isDeleting}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menghapus...</span>
      {:else}
        <Trash2 class="w-4 h-4 mr-1" />
        <span>Ya, Hapus Data</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
