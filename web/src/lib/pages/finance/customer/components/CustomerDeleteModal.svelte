<script lang="ts">
  import { Trash2, AlertCircle, Loader2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
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

<M3Modal
  {isOpen}
  title="Hapus Penjamin?"
  subtitle="Konfirmasi penghapusan data penjamin / debitur"
  icon={Trash2}
  iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"
  size="md"
  {onClose}
>
  {#if customer}
    <div class="space-y-4 text-xs">
      {#if errorMessage}
        <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-300 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs flex items-start gap-2.5 shadow-2xs">
          <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
          <span class="font-medium leading-relaxed">{errorMessage}</span>
        </div>
      {/if}

      <p class="text-[#444746] dark:text-[#c4c7c5] leading-relaxed">
        Apakah Anda yakin ingin menghapus data debitur/penjamin berikut?
      </p>

      <div class="p-3.5 rounded-2xl bg-[#f0f4f9] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] space-y-1.5">
        <div class="flex items-center gap-2">
          <span class="font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa]">{customer.code}</span>
          <span class="text-[#747775] dark:text-[#8e918f]">•</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{customer.name}</span>
        </div>
        {#if customer.customer_type}
          <p class="text-[11px] text-[#444746] dark:text-[#c4c7c5]">
            Kategori: <strong>{customer.customer_type.name}</strong>
          </p>
        {/if}
        {#if customer.contact_person}
          <p class="text-[11px] text-[#444746] dark:text-[#c4c7c5]">
            PIC: {customer.contact_person}
          </p>
        {/if}
      </div>

      <div class="p-3 rounded-xl bg-rose-50/80 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/60">
        <p class="text-xs text-rose-700 dark:text-rose-300 leading-relaxed">
          <strong>Peringatan:</strong> Tindakan ini akan menonaktifkan/menghapus penjamin. Pastikan tidak ada tagihan piutang berjalan yang menggunakan penjamin ini.
        </p>
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isDeleting}>
      Batal
    </M3Button>
    <button
      type="button"
      onclick={handleDelete}
      disabled={isDeleting}
      class="flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-rose-600 hover:bg-rose-700 dark:bg-rose-600 dark:hover:bg-rose-500 rounded-xl transition-all cursor-pointer shadow-xs disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {#if isDeleting}
        <Loader2 class="w-3.5 h-3.5 animate-spin" />
        <span>Menghapus...</span>
      {:else}
        <Trash2 class="w-3.5 h-3.5" />
        <span>Hapus Penjamin</span>
      {/if}
    </button>
  {/snippet}
</M3Modal>
