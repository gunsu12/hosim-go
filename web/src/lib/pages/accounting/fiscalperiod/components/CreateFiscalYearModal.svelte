<script lang="ts">
  import { CalendarPlus, AlertCircle, Loader2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import { createFiscalYear } from '$lib/api/accounting';
  import type { FiscalYearRecord } from '$lib/types/accounting/fiscalPeriod';

  interface Props {
    isOpen: boolean;
    existingYears?: number[];
    onClose: () => void;
    onSuccess: (created: FiscalYearRecord) => void;
  }

  let { isOpen, existingYears = [], onClose, onSuccess }: Props = $props();

  const currentYear = new Date().getFullYear();
  let year = $state<number>(currentYear);
  let notes = $state('');
  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      // Auto-suggest next available year
      if (existingYears.length > 0) {
        const maxYear = Math.max(...existingYears);
        year = maxYear >= currentYear ? maxYear + 1 : currentYear;
      } else {
        year = currentYear;
      }
      notes = `Tahun Buku Akuntansi RS ${year}`;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!year || year < 2000 || year > 2100) {
      errorMessage = 'Tahun harus di antara 2000 dan 2100';
      return;
    }

    if (existingYears.includes(year)) {
      errorMessage = `Tahun fiskal ${year} sudah terdaftar dalam sistem`;
      return;
    }

    isSubmitting = true;
    errorMessage = null;

    try {
      const res = await createFiscalYear({
        year,
        notes: notes.trim() || undefined
      });
      onSuccess(res);
      onClose();
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal membuat tahun fiskal baru';
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Buka Tahun Buku Fiskal Baru"
  subtitle="Otomatis men-generate 12 periode bulanan standar (01 Jan - 31 Des)"
  icon={CalendarPlus}
  size="md"
  onsubmit={handleSubmit}
  {onClose}
>
  <div class="space-y-4">
    <div class="p-3.5 rounded-2xl bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-800 text-blue-800 dark:text-blue-200 text-xs leading-relaxed">
      <p class="font-bold mb-1">Ketentuan Otomatisasi Sistem:</p>
      <p>
        Sistem backend akan secara otomatis membuat <strong>1 tahun buku</strong> dan <strong>12 periode bulanan</strong> (Januari s/d Desember) dengan status awal <strong>TERBUKA (OPEN)</strong> serta penanggalan kabisat yang akurat.
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
      <label for="year-input" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Tahun Kalender Fiskal <span class="text-rose-500">*</span>
      </label>
      <input
        id="year-input"
        type="number"
        min="2000"
        max="2100"
        bind:value={year}
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <p class="mt-1 text-[11px] text-[#444746] dark:text-[#c4c7c5]">
        Rentang periode: 01 Januari {year || 'YYYY'} s/d 31 Desember {year || 'YYYY'}
      </p>
    </div>

    <div>
      <label for="year-notes" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1">
        Keterangan / Catatan
      </label>
      <textarea
        id="year-notes"
        bind:value={notes}
        rows="2"
        placeholder="Misal: Tahun Anggaran & Operasional RS 2027..."
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
      >
        {#if isSubmitting}
          <Loader2 class="w-4 h-4 animate-spin mr-1.5" />
          <span>Membuat...</span>
        {:else}
          <CalendarPlus class="w-4 h-4 mr-1.5" />
          <span>Buat Tahun Fiskal</span>
        {/if}
      </M3Button>
    </div>
  {/snippet}
</M3Modal>
