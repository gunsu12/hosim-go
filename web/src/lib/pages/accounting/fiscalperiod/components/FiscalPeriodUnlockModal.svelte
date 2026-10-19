<script lang="ts">
  import { Unlock, AlertCircle, Loader2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { unlockFiscalPeriod } from '$lib/api/accounting';
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
      const res = await unlockFiscalPeriod(period.id, {
        notes: notes.trim() || undefined
      });
      onSuccess(res);
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal membuka kunci periode fiskal';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Buka Kunci Periode Fiskal"
  subtitle={period ? `${period.name} (${period.period_key})` : ''}
  icon={Unlock}
  iconClass="bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400"
  size="md"
  onsubmit={handleSubmit}
  {onClose}
>
  <div class="space-y-4">
    <div class="p-3.5 rounded-2xl bg-emerald-50 dark:bg-emerald-950/30 border border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-200 text-xs leading-relaxed">
      <p class="font-bold mb-1">Membuka Kunci (OPEN):</p>
      <p>
        Periode ini akan kembali berstatus TERBUKA (OPEN). Seluruh staf dengan hak akses pencatatan transaksi jurnal dapat kembali memposting transaksi ke periode ini.
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
      <label for="unlock-notes" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Catatan Pembukaan Kunci
      </label>
      <textarea
        id="unlock-notes"
        bind:value={notes}
        rows="3"
        placeholder="Misal: Dibuka kembali atas persetujuan manajer keuangan untuk koreksi saldo..."
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
        class="bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white"
      >
        {#if isSubmitting}
          <Loader2 class="w-4 h-4 animate-spin mr-1.5" />
          <span>Membuka...</span>
        {:else}
          <Unlock class="w-4 h-4 mr-1.5" />
          <span>Buka Kunci Periode</span>
        {/if}
      </M3Button>
    </div>
  {/snippet}
</M3Modal>
