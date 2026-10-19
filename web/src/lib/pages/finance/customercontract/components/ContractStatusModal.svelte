<script lang="ts">
  import {
    CheckCircle2,
    PauseCircle,
    PlayCircle,
    XCircle,
    AlertCircle,
    ShieldAlert,
    FileSpreadsheet,
  } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { CustomerContractRecord } from '$lib/types/finance/customer_contract';

  export type TransitionAction = 'ACTIVATE' | 'SUSPEND' | 'REACTIVATE' | 'TERMINATE';

  interface Props {
    isOpen: boolean;
    contract: CustomerContractRecord | null;
    action: TransitionAction;
    isSubmitting?: boolean;
    errorMessage?: string | null;
    onClose: () => void;
    onConfirm: (reason?: string) => Promise<void> | void;
  }

  let {
    isOpen,
    contract,
    action = 'ACTIVATE',
    isSubmitting = false,
    errorMessage = null,
    onClose,
    onConfirm,
  }: Props = $props();

  let reason = $state('');

  $effect(() => {
    if (isOpen) {
      reason = '';
    }
  });

  const modalConfig = $derived.by(() => {
    switch (action) {
      case 'ACTIVATE':
        return {
          title: 'Aktifkan Kontrak PKS?',
          subtitle: 'Kontrak aktif akan dapat dipilih dalam proses registrasi pasien',
          icon: CheckCircle2,
          iconClass: 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400',
          btnClass: 'bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white',
          btnLabel: 'Aktifkan Kontrak',
          requiresReason: false,
        };
      case 'SUSPEND':
        return {
          title: 'Tangguhkan Kontrak PKS (Suspend)?',
          subtitle: 'Kontrak yang ditangguhkan tidak dapat dipilih untuk registrasi baru',
          icon: PauseCircle,
          iconClass: 'bg-amber-50 text-amber-600 dark:bg-amber-950/60 dark:text-amber-400',
          btnClass: 'bg-amber-600 hover:bg-amber-700 active:bg-amber-800 text-white',
          btnLabel: 'Tangguhkan Kontrak',
          requiresReason: true,
        };
      case 'REACTIVATE':
        return {
          title: 'Aktifkan Kembali Kontrak (Reactivate)?',
          subtitle: 'Mengembalikan status kontrak ke aktif operasional',
          icon: PlayCircle,
          iconClass: 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/60 dark:text-emerald-400',
          btnClass: 'bg-emerald-600 hover:bg-emerald-700 active:bg-emerald-800 text-white',
          btnLabel: 'Aktifkan Kembali',
          requiresReason: false,
        };
      case 'TERMINATE':
        return {
          title: 'Putus / Hentikan Kontrak (Terminate)?',
          subtitle: 'PKS dihentikan permanen dan tidak dapat diaktifkan kembali',
          icon: XCircle,
          iconClass: 'bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400',
          btnClass: 'bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white',
          btnLabel: 'Putus Kontrak Permanen',
          requiresReason: true,
        };
    }
  });

  // Invariant validation: Activate requires at least 1 active schema and effective_date
  let activeSchemasCount = $derived(
    (contract?.price_schemas || []).filter((s) => s.is_active).length,
  );
  let hasSchemasForActivation = $derived(
    action !== 'ACTIVATE' && action !== 'REACTIVATE' ? true : activeSchemasCount > 0,
  );
  let hasEffectiveDate = $derived(Boolean(contract?.effective_date));

  let canSubmit = $derived.by(() => {
    if (modalConfig.requiresReason && !reason.trim()) return false;
    if ((action === 'ACTIVATE' || action === 'REACTIVATE') && (!hasSchemasForActivation || !hasEffectiveDate)) {
      return false;
    }
    return true;
  });

  function handleSubmit() {
    if (!canSubmit) return;
    onConfirm(modalConfig.requiresReason ? reason.trim() : undefined);
  }
</script>

<M3Modal
  {isOpen}
  title={modalConfig.title}
  subtitle={modalConfig.subtitle}
  icon={modalConfig.icon}
  iconClass={modalConfig.iconClass}
  size="md"
  {onClose}
>
  {#if contract}
    <div class="space-y-4">
      <!-- Error Message in Modal -->
      {#if errorMessage}
        <div
          class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-300 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs flex items-start gap-2.5 shadow-2xs"
        >
          <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
          <span class="font-medium leading-relaxed">{errorMessage}</span>
        </div>
      {/if}

      <!-- Contract Snapshot Card -->
      <div
        class="p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] text-xs space-y-1.5"
      >
        <div class="flex justify-between items-center">
          <span class="text-[#747775] dark:text-[#8e918f]">Nomor Kontrak:</span>
          <span class="font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa]">
            {contract.contract_number}
          </span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-[#747775] dark:text-[#8e918f]">Penjamin:</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] truncate max-w-[200px]">
            {contract.customer_name || contract.customer?.name || '-'}
          </span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-[#747775] dark:text-[#8e918f]">Skema Tarif Terhubung:</span>
          <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
            {activeSchemasCount} skema aktif
          </span>
        </div>
      </div>

      <!-- Activation Invariant Warning -->
      {#if (action === 'ACTIVATE' || action === 'REACTIVATE') && !hasSchemasForActivation}
        <div
          class="p-3 rounded-xl bg-amber-50 dark:bg-amber-950/40 border border-amber-300 dark:border-amber-800 text-xs text-amber-900 dark:text-amber-200 flex items-start gap-2"
        >
          <ShieldAlert class="w-4 h-4 text-amber-600 shrink-0 mt-0.5" />
          <div>
            <p class="font-bold">Invariant Modul Terlanggar</p>
            <p class="text-[11px] mt-0.5 leading-relaxed">
              Kontrak belum memiliki <strong>Price Schema</strong> aktif. Hubungkan minimal 1 skema tarif
              terlebih dahulu sebelum mengaktifkan kontrak.
            </p>
          </div>
        </div>
      {/if}

      <!-- Mandatory Reason Field -->
      {#if modalConfig.requiresReason}
        <div>
          <label
            for="status-reason-input"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1"
          >
            Alasan Perubahan Status <span class="text-rose-500">*</span>
          </label>
          <textarea
            id="status-reason-input"
            bind:value={reason}
            rows="3"
            placeholder="Jelaskan alasan penetapan status ini untuk kebutuhan audit finansial..."
            class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
          ></textarea>
          <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
            Alasan ini akan dicatat ke dalam audit trail append-only HOSIM-GO.
          </p>
        </div>
      {/if}

      <!-- Warning for Terminate -->
      {#if action === 'TERMINATE'}
        <div
          class="p-3 rounded-xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/60 text-xs text-rose-700 dark:text-rose-300"
        >
          <strong>Catatan:</strong> Status TERMINATED bersifat permanen. Jika kerjasama diperbarui di kemudian hari, Anda wajib membuat perpanjangan kontrak (Renew).
        </div>
      {/if}
    </div>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      Batal
    </M3Button>
    <button
      type="button"
      onclick={handleSubmit}
      disabled={isSubmitting || !canSubmit}
      class="inline-flex items-center justify-center gap-2 h-10 px-5 text-sm font-medium rounded-full transition-all duration-150 cursor-pointer {modalConfig.btnClass} shadow-xs disabled:opacity-40 disabled:cursor-not-allowed"
    >
      {#if isSubmitting}
        <div
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-1"
        ></div>
        <span>Memproses...</span>
      {:else}
        <span>{modalConfig.btnLabel}</span>
      {/if}
    </button>
  {/snippet}
</M3Modal>
