<script lang="ts">
  import { Trash2, RefreshCw, ShieldAlert, FolderTree, AlertCircle } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { AccountRecord } from '$lib/types/accounting/account';

  interface Props {
    isOpen: boolean;
    account: AccountRecord | null;
    isDeleting?: boolean;
    errorMessage?: string | null;
    onClose: () => void;
    onConfirm: () => Promise<void> | void;
  }

  let {
    isOpen,
    account,
    isDeleting = false,
    errorMessage = null,
    onClose,
    onConfirm
  }: Props = $props();

  let isLevelOne = $derived(account?.account_level === 1);
</script>

<M3Modal
  {isOpen}
  title="Hapus Akun Buku Besar?"
  subtitle="Tindakan ini tidak dapat dibatalkan jika akun memiliki dependensi"
  icon={Trash2}
  iconClass="bg-rose-50 text-rose-600 dark:bg-rose-950/60 dark:text-rose-400"
  size="md"
  {onClose}
>
  {#if account}
    <div class="space-y-4">
      <!-- Error Message in Modal -->
      {#if errorMessage}
        <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-300 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs flex items-start gap-2.5 shadow-2xs">
          <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
          <span class="font-medium leading-relaxed">{errorMessage}</span>
        </div>
      {/if}

      <!-- Account Summary Card -->
      <div class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-surface-container border border-outline-variant/60 text-xs flex flex-col gap-2">
        <div class="flex justify-between items-center">
          <span class="text-on-surface-variant">Kode Akun:</span>
          <span class="font-mono font-bold text-primary">{account.code}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-on-surface-variant">Nama Akun:</span>
          <span class="font-semibold text-on-surface">{account.name}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-on-surface-variant">Klasifikasi:</span>
          <span class="font-medium text-on-surface">{account.type} (Level {account.account_level})</span>
        </div>
      </div>

      <!-- Level 1 Guard Alert -->
      {#if isLevelOne}
        <div class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-300 dark:border-amber-800/60 flex items-start gap-2.5 text-xs text-amber-900 dark:text-amber-200">
          <ShieldAlert class="w-4 h-4 text-amber-600 dark:text-amber-400 shrink-0 mt-0.5" />
          <div>
            <p class="font-bold">Akun Klasifikasi Utama (Level 1)</p>
            <p class="text-[11px] text-amber-800/90 dark:text-amber-300/80 mt-0.5">
              Akun level 1 ({account.code} - {account.name}) merupakan klasifikasi pokok sistem akuntansi rumah sakit dan dilindungi secara permanen dari penghapusan.
            </p>
          </div>
        </div>
      {:else if !account.is_postable}
        <div class="p-3 rounded-xl bg-blue-50/70 dark:bg-blue-950/40 border border-blue-200 dark:border-blue-900/50 text-xs text-blue-800 dark:text-blue-300 flex items-start gap-2">
          <FolderTree class="w-4 h-4 text-blue-600 dark:text-blue-400 shrink-0 mt-0.5" />
          <span>Akun ini berstatus <strong>Header/Induk</strong>. Penghapusan akan otomatis ditolak jika akun masih memiliki sub-akun anak aktif di bawahnya.</span>
        </div>
      {/if}

      <!-- Standard Warning -->
      <div class="p-3 rounded-xl bg-rose-50/80 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-900/60">
        <p class="text-xs text-rose-700 dark:text-rose-300 leading-relaxed">
          <strong>Peringatan:</strong> Akun hanya dapat dihapus jika belum memiliki sub-akun anak dan belum pernah digunakan dalam jurnal transaksi buku besar keuangan maupun konfigurasi tarif/katalog.
        </p>
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <M3Button
      variant="text"
      onclick={onClose}
      disabled={isDeleting}
    >
      Batal
    </M3Button>
    <button
      type="button"
      onclick={onConfirm}
      disabled={isDeleting || isLevelOne}
      title={isLevelOne ? 'Akun level 1 tidak dapat dihapus' : 'Hapus Akun'}
      class="flex items-center gap-1.5 px-4 py-2 text-xs font-semibold text-white bg-rose-600 hover:bg-rose-700 dark:bg-rose-600 dark:hover:bg-rose-500 rounded-xl transition-all cursor-pointer shadow-xs disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {#if isDeleting}
        <RefreshCw class="w-3.5 h-3.5 animate-spin" />
        <span>Menghapus...</span>
      {:else}
        <Trash2 class="w-3.5 h-3.5" />
        <span>Hapus Akun</span>
      {/if}
    </button>
  {/snippet}
</M3Modal>
