<script lang="ts">
  import { Lock, AlertCircle, Loader2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { lockFiscalPeriod } from '$lib/api/accounting';
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
      const res = await lockFiscalPeriod(period.id, {
        notes: notes.trim() || undefined
      });
      onSuccess(res);
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal mengunci periode fiskal';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Kunci Sementara Periode Fiskal"
  subtitle={period ? `${period.name} (${period.period_key})` : ''}
  icon={Lock}
  iconClass="bg-amber-50 text-amber-600 dark:bg-amber-950/60 dark:text-amber-400"
  size="md"
  onsubmit={handleSubmit}
  {onClose}
>
  <div class="space-y-4">
    <div class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/30 border border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-200 text-xs leading-relaxed">
      <p class="font-bold mb-1">Perhatian Penguncian Sementara (LOCKED):</p>
      <p>
        Mengunci periode akan mencegah posting transaksi baru ke buku besar untuk tanggal transaksi di bulan ini. Periode yang dikunci masih dapat dibuka kembali (Unlock) sebelum tutup buku permanen.
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
      <label for="lock-notes" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Alasan / Catatan Penguncian
      </label>
      <textarea
        id="lock-notes"
        bind:value={notes}
        rows="3"
        placeholder="Misal: Penguncian akhir bulan untuk rekonsiliasi kas/bank & penyesuaian..."
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
        class="bg-amber-600 hover:bg-amber-700 active:bg-amber-800 text-white"
      >
        {#if isSubmitting}
          <Loader2 class="w-4 h-4 animate-spin mr-1.5" />
          <span>Mengunci...</span>
        {:else}
          <Lock class="w-4 h-4 mr-1.5" />
          <span>Kunci Periode</span>
        {/if}
      </M3Button>
    </div>
  {/snippet}
</M3Modal>
