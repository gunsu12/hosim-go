<script lang="ts">
  import { AlertTriangle, X, Trash2 } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { deleteCustomer } from '$lib/api/finance';
  import type { CustomerRecord } from '$lib/types/finance';

  interface Props {
    isOpen: boolean;
    customer?: CustomerRecord | null;
    onClose: () => void;
    onSuccess: () => void;
  }

  let {
    isOpen = false,
    customer = null,
    onClose = () => {},
    onSuccess = () => {},
  }: Props = $props();

  let isDeleting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      isDeleting = false;
    }
  });

  async function handleDelete() {
    if (!customer) return;

    isDeleting = true;
    errorMessage = null;

    try {
      await deleteCustomer(customer.id);
      onSuccess();
    } catch (err: any) {
      console.error('Failed to delete customer:', err);
      errorMessage = err?.message || 'Gagal menghapus penjamin. Pastikan tidak ada transaksi aktif yang mengikat.';
    } finally {
      isDeleting = false;
    }
  }
</script>

{#if isOpen && customer}
  <div
    class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b border-[#e1e5ea] flex items-center justify-between bg-[#f8fafd]">
        <div class="flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-2xl bg-rose-100 text-rose-700 flex items-center justify-center">
            <AlertTriangle class="w-5 h-5" />
          </div>
          <h3 class="text-sm font-bold text-[#1f1f1f]">
            Konfirmasi Hapus Penjamin
          </h3>
        </div>
        <button
          onclick={onClose}
          type="button"
          class="text-[#747775] hover:text-[#1f1f1f] p-1.5 rounded-full hover:bg-[#e1e5ea] transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 text-xs space-y-4 text-[#1f1f1f]">
        {#if errorMessage}
          <div class="p-3 rounded-2xl bg-rose-50 border border-rose-200 text-rose-700">
            <strong>Gagal:</strong> {errorMessage}
          </div>
        {/if}

        <p class="text-[#444746] leading-relaxed">
          Apakah Anda yakin ingin menghapus data debitur/penjamin berikut?
        </p>

        <div class="p-3.5 rounded-2xl bg-[#f0f4f9] border border-[#e1e5ea] space-y-1">
          <div class="flex items-center gap-2">
            <span class="font-mono font-bold text-[#0b57d0]">{customer.code}</span>
            <span class="text-[#747775]">•</span>
            <span class="font-semibold text-[#1f1f1f]">{customer.name}</span>
          </div>
          {#if customer.customer_type}
            <p class="text-[11px] text-[#444746]">
              Kategori: <strong>{customer.customer_type.name}</strong>
            </p>
          {/if}
          {#if customer.contact_person}
            <p class="text-[11px] text-[#444746]">
              PIC: {customer.contact_person}
            </p>
          {/if}
        </div>

        <p class="text-[11px] text-[#747775]">
          Tindakan ini akan melakukan soft-delete pada penjamin di database.
        </p>
      </div>

      <!-- Modal Footer -->
      <div class="px-6 py-3.5 border-t border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-end gap-2.5">
        <M3Button variant="text" onclick={onClose} disabled={isDeleting}>
          Batal
        </M3Button>

        <button
          type="button"
          onclick={handleDelete}
          disabled={isDeleting}
          class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-rose-600 hover:bg-rose-700 text-white font-semibold text-xs transition-colors disabled:opacity-50 cursor-pointer shadow-xs"
        >
          {#if isDeleting}
            <span class="inline-block animate-spin">⏳</span>
            <span>Menghapus...</span>
          {:else}
            <Trash2 class="w-4 h-4" />
            <span>Hapus Penjamin</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
