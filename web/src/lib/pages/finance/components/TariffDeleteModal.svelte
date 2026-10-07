<script lang="ts">
  import { AlertTriangle, X } from '@lucide/svelte';
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

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-sm p-6 flex flex-col gap-4">
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-3">
        <div class="flex items-center gap-2 text-[#c5221f]">
          <AlertTriangle class="w-5 h-5" />
          <h3 class="text-sm font-bold text-[#1f1f1f]">{title}</h3>
        </div>
        <button onclick={onClose} type="button" class="text-[#747775] hover:text-[#1f1f1f]">
          <X class="w-4 h-4" />
        </button>
      </div>

      <div class="text-xs text-[#444746] flex flex-col gap-2">
        <p>Apakah Anda yakin ingin menghapus data berikut?</p>
        <div class="p-3 rounded-xl bg-[#fce8e6]/50 border border-[#fad2cf] flex flex-col gap-1">
          <span class="font-mono font-bold text-[#c5221f] text-xs">[{itemCode}]</span>
          <span class="font-semibold text-[#1f1f1f]">{itemName}</span>
        </div>
        <p class="text-[11px] text-[#747775]">Data yang dihapus akan ditandai soft delete di database dan dinonaktifkan dari referensi pelayanan.</p>
      </div>

      <div class="flex justify-end gap-2 pt-2 border-t border-[#e1e5ea]">
        <button
          type="button"
          onclick={onClose}
          class="px-3.5 py-1.5 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] hover:bg-[#f0f4f9] transition-colors"
        >
          Batal
        </button>
        <M3Button variant="filled" onclick={onConfirm} disabled={isDeleting}>
          <span class="text-white">{isDeleting ? 'Menghapus...' : 'Ya, Hapus'}</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}
