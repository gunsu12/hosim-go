<script lang="ts">
  import { ShieldAlert, AlertCircle, Loader2, CheckCircle2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { closeFiscalYear } from '$lib/api/accounting';
  import type { FiscalYearRecord } from '$lib/types/accounting/fiscalPeriod';

  interface Props {
    isOpen: boolean;
    fiscalYear: FiscalYearRecord | null;
    onClose: () => void;
    onSuccess: (closed: FiscalYearRecord) => void;
  }

  let { isOpen, fiscalYear, onClose, onSuccess }: Props = $props();

  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  // Cek apakah seluruh 12 periode sudah berstatus CLOSED
  const openOrLockedPeriods = $derived.by(() => {
    if (!fiscalYear || !fiscalYear.periods) return [];
    return fiscalYear.periods.filter((p) => p.status !== 'CLOSED');
  });

  const canClose = $derived(
    Boolean(
      fiscalYear &&
      !fiscalYear.is_closed &&
      fiscalYear.periods &&
      fiscalYear.periods.length === 12 &&
      openOrLockedPeriods.length === 0
    )
  );

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!fiscalYear || !canClose) return;

    isSubmitting = true;
    errorMessage = null;

    try {
      const res = await closeFiscalYear(fiscalYear.id);
      onSuccess(res);
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal menutup tahun buku fiskal';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Tutup Buku Akhir Tahun Fiskal (Year-End Closing)"
  subtitle={fiscalYear ? `Tahun Buku ${fiscalYear.year}` : ''}
  icon={ShieldAlert}
  iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"
  size="md"
  onsubmit={handleSubmit}
  {onClose}
>
  <div class="space-y-4">
    {#if !canClose}
      <div class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-800 text-amber-800 dark:text-amber-200 text-xs leading-relaxed">
        <div class="flex items-center gap-2 font-bold mb-1 text-amber-700 dark:text-amber-300">
          <AlertCircle class="w-4 h-4 shrink-0" />
          <span>Prasyarat Tutup Buku Belum Terpenuhi:</span>
        </div>
        <p>
          Seluruh <strong>12 periode bulanan</strong> pada tahun {fiscalYear?.year} wajib berstatus <strong>CLOSED</strong> terlebih dahulu sebelum tutup buku tahunan dapat disahkan.
        </p>
        {#if openOrLockedPeriods.length > 0}
          <div class="mt-2.5 pt-2 border-t border-amber-200/60 dark:border-amber-800/60">
            <span class="font-semibold block mb-1">Periode yang belum ditutup ({openOrLockedPeriods.length}):</span>
            <div class="flex flex-wrap gap-1.5 max-h-24 overflow-y-auto">
              {#each openOrLockedPeriods as op}
                <span class="px-2 py-0.5 rounded-md bg-amber-100 dark:bg-amber-900/60 text-amber-900 dark:text-amber-100 font-mono text-[10px] font-bold">
                  {op.period_key} ({op.status})
                </span>
              {/each}
            </div>
          </div>
        {/if}
      </div>
    {:else}
      <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs leading-relaxed">
        <div class="flex items-center gap-2 font-bold mb-1 text-rose-700 dark:text-rose-300">
          <ShieldAlert class="w-4 h-4 shrink-0" />
          <span>PERINGATAN TUTUP BUKU TAHUNAN:</span>
        </div>
        <p>
          Anda akan menutup buku tahun fiskal <strong>{fiscalYear?.year}</strong> secara permanen. Tindakan ini menandai berakhirnya siklus akuntansi tahun berjalan dan mengunci seluruh mutasi buku besar secara permanen.
        </p>
      </div>
    {/if}

    {#if errorMessage}
      <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
        <AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
        <div>
          <span class="font-bold">Gagal: </span>
          <span>{errorMessage}</span>
        </div>
      </div>
    {/if}

    <div class="p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] text-xs text-[#444746] dark:text-[#c4c7c5] space-y-1">
      <div class="flex justify-between">
        <span>Tahun Fiskal:</span>
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">{fiscalYear?.year}</span>
      </div>
      <div class="flex justify-between">
        <span>Rentang Kalender:</span>
        <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{fiscalYear?.start_date} s/d {fiscalYear?.end_date}</span>
      </div>
      <div class="flex justify-between">
        <span>Status 12 Bulan:</span>
        <span class="font-bold {canClose ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'}">
          {canClose ? '12/12 Periode CLOSED (Siap Tutup)' : `${openOrLockedPeriods.length} Periode Belum CLOSED`}
        </span>
      </div>
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
        disabled={isSubmitting || !canClose}
        class="bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white"
      >
        {#if isSubmitting}
          <Loader2 class="w-4 h-4 animate-spin mr-1.5" />
          <span>Memproses Tutup Tahun...</span>
        {:else}
          <CheckCircle2 class="w-4 h-4 mr-1.5" />
          <span>Tutup Tahun Fiskal</span>
        {/if}
      </M3Button>
    </div>
  {/snippet}
</M3Modal>
