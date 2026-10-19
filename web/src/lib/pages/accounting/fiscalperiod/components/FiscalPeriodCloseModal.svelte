<script lang="ts">
  import { ShieldAlert, AlertCircle, Loader2, CheckCircle2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { closeFiscalPeriod } from '$lib/api/accounting';
  import type { FiscalPeriodRecord } from '$lib/types/accounting/fiscalPeriod';

  interface Props {
    isOpen: boolean;
    period: FiscalPeriodRecord | null;
    onClose: () => void;
    onSuccess: (updated: FiscalPeriodRecord) => void;
  }

  let { isOpen, period, onClose, onSuccess }: Props = $props();

  let notes = $state('');
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      notes = period?.notes || '';
      errorMessage = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!period) return;

    isSubmitting = true;
    errorMessage = null;

    try {
      const res = await closeFiscalPeriod(period.id, {
        notes: notes.trim() || undefined
      });
      onSuccess(res);
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal menutup buku periode fiskal';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Tutup Buku Bulanan Permanen (CLOSED)"
  subtitle={period ? `${period.name} (${period.period_key})` : ''}
  icon={ShieldAlert}
  iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"
  size="md"
  onsubmit={handleSubmit}
  {onClose}
>
  <div class="space-y-4">
    <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs leading-relaxed">
      <div class="flex items-center gap-2 font-bold mb-1.5 text-rose-700 dark:text-rose-300">
        <ShieldAlert class="w-4 h-4 shrink-0" />
        <span>PERINGATAN INTEGRITAS AKUNTANSI:</span>
      </div>
      <p>
        Menutup buku periode bulanan adalah <strong>tindakan permanen dan ireversibel</strong>. Begitu periode berstatus <strong>CLOSED</strong>, transaksi jurnal di bulan ini tidak akan pernah dapat ditambah, diedit, atau dihapus selamanya demi kepatuhan audit PSAK RS.
      </p>
    </div>

    {#if errorMessage}
      <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
        <AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
        <div>
          <span class="font-bold">Gagal: </span>
          <span>{errorMessage}</span>
        </div>
      </div>
    {/if}

    <div>
      <label for="close-notes" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Catatan Pengesahan Tutup Buku
      </label>
      <textarea
        id="close-notes"
        bind:value={notes}
        rows="3"
        placeholder="Misal: Tutup buku bulanan disahkan setelah rekonsiliasi kas-bank selesai 100%..."
        class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-medium text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
      ></textarea>
    </div>
  </div>

  {#snippet footer()}
    <div class="flex items-center justify-end gap-2.5 w-full">
      <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
        Batal
      </M3Button>
      <M3Button
        variant="filled"
        type="submit"
        disabled={isSubmitting}
        class="bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white"
      >
        {#if isSubmitting}
          <Loader2 class="w-4 h-4 animate-spin mr-1.5" />
          <span>Memproses Tutup Buku...</span>
        {:else}
          <CheckCircle2 class="w-4 h-4 mr-1.5" />
          <span>Tutup Buku Permanen</span>
        {/if}
      </M3Button>
    </div>
  {/snippet}
</M3Modal>
