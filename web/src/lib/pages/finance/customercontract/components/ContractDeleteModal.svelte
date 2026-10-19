<script lang="ts">
  import { Trash2, AlertCircle, ShieldAlert, FileText } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { CustomerContractRecord } from '$lib/types/finance/customer_contract';

  interface Props {
    isOpen: boolean;
    contract: CustomerContractRecord | null;
    isDeleting?: boolean;
    errorMessage?: string | null;
    onClose: () => void;
    onConfirm: () => Promise<void> | void;
  }

  let {
    isOpen,
    contract,
    isDeleting = false,
    errorMessage = null,
    onClose,
    onConfirm,
  }: Props = $props();

  let isDraft = $derived(contract?.status === 'DRAFT');
</script>

<M3Modal
  {isOpen}
  title="Hapus Draft Kontrak PKS?"
  subtitle="Tindakan ini akan menghapus data draft kontrak secara permanen"
  icon={Trash2}
  iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"
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

      <!-- Contract Summary Card -->
      <div
        class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] text-xs flex flex-col gap-2"
      >
        <div class="flex justify-between items-center">
          <span class="text-[#747775] dark:text-[#8e918f]">Nomor Kontrak:</span>
          <span class="font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa]">
            {contract.contract_number}
          </span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-[#747775] dark:text-[#8e918f]">Nama Kontrak:</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{contract.name}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-[#747775] dark:text-[#8e918f]">Status Saat Ini:</span>
          <span
            class="px-2 py-0.5 rounded-full font-bold text-[10px] {contract.status === 'DRAFT'
              ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300'
              : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300'}"
          >
            {contract.status}
          </span>
        </div>
      </div>

      <!-- Invariant Guard: Hanya DRAFT yang boleh dihapus -->
      {#if !isDraft}
        <div
          class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-300 dark:border-amber-800/60 flex items-start gap-2.5 text-xs text-amber-900 dark:text-amber-200"
        >
          <ShieldAlert class="w-4 h-4 text-amber-600 dark:text-amber-400 shrink-0 mt-0.5" />
          <div>
            <p class="font-bold">Proteksi Integritas Audit Finansial</p>
            <p class="text-[11px] text-amber-800/90 dark:text-amber-300/80 mt-0.5">
              Kontrak berstatus <strong>{contract.status}</strong> tidak boleh dihapus dari sistem
              karena telah memiliki rekam jejak audit dan acuan tarif pelayanan. Gunakan menu
              <strong>Terminate (Pemutusan Kontrak)</strong> jika kerjasama telah berakhir.
            </p>
          </div>
        </div>
      {:else}
        <!-- Warning for Draft Deletion -->
        <div
          class="p-3.5 rounded-2xl bg-rose-50/80 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/60 text-xs text-rose-700 dark:text-rose-300 flex items-start gap-2"
        >
          <FileText class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
          <p class="leading-relaxed">
            Draft kontrak dan seluruh relasi skema tarif yang terlampir akan dihapus dari sistem.
            Tindakan ini tidak dapat dibatalkan.
          </p>
        </div>
      {/if}
    </div>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isDeleting}>
      Batal
    </M3Button>
    <button
      type="button"
      onclick={onConfirm}
      disabled={isDeleting || !isDraft}
      class="inline-flex items-center justify-center gap-2 h-10 px-5 text-sm font-medium rounded-full transition-all duration-150 cursor-pointer bg-rose-600 hover:bg-rose-700 active:bg-rose-800 text-white shadow-xs disabled:opacity-40 disabled:cursor-not-allowed"
    >
      {#if isDeleting}
        <div
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-1"
        ></div>
        <span>Menghapus...</span>
      {:else}
        <Trash2 class="w-4 h-4 mr-1" />
        <span>Hapus Draft</span>
      {/if}
    </button>
  {/snippet}
</M3Modal>
