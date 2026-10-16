<script lang="ts">
  import { ShieldAlert, X, Power, AlertCircle } from '@lucide/svelte';
  import type { PriceSchemaRecord } from '$lib/types/finance/price_schema';

  interface Props {
    isOpen: boolean;
    schema: PriceSchemaRecord | null;
    isProcessing: boolean;
    onClose: () => void;
    onConfirm: () => Promise<void>;
  }

  let { isOpen, schema, isProcessing, onClose, onConfirm }: Props = $props();
</script>

{#if isOpen && schema}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-xs animate-in fade-in duration-150">
    <div class="bg-white dark:bg-surface-container rounded-3xl border border-outline-variant/60 shadow-xl max-w-md w-full overflow-hidden flex flex-col">
      <!-- Modal Header -->
      <div class="p-6 pb-4 flex items-start gap-4">
        <div class="w-10 h-10 rounded-2xl {schema.is_active ? 'bg-amber-50 text-amber-600 dark:bg-amber-950/50 dark:text-amber-400' : 'bg-emerald-50 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400'} flex items-center justify-center shrink-0">
          <ShieldAlert class="w-5 h-5" />
        </div>
        <div class="flex-1 min-w-0">
          <h3 class="text-sm font-bold text-on-surface">
            {schema.is_active ? 'Nonaktifkan Skema Tarif?' : 'Aktifkan Kembali Skema Tarif?'}
          </h3>
          <p class="text-xs text-on-surface-variant mt-1 leading-relaxed">
            {#if schema.is_active}
              Skema <span class="font-bold text-on-surface font-mono">{schema.code}</span> ({schema.name}) tidak akan dapat dipilih untuk pendaftaran/registrasi pasien baru. Tagihan dan rekam transaksi masa lalu tetap aman.
            {:else}
              Skema <span class="font-bold text-on-surface font-mono">{schema.code}</span> ({schema.name}) akan diaktifkan kembali dan dapat dipilih pada pendaftaran pasien.
            {/if}
          </p>
        </div>
        <button
          type="button"
          onclick={onClose}
          class="p-1 rounded-full text-on-surface-variant hover:text-on-surface cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <div class="px-6 py-2">
        <div class="p-3 rounded-2xl bg-surface-container-high/60 border border-outline-variant/60 text-[11px] text-on-surface-variant flex items-center gap-2">
          <AlertCircle class="w-4 h-4 text-primary shrink-0" />
          <span>Sesuai aturan finansial HOSIM-GO, skema harga tidak di-hard delete untuk menjaga integritas data billing historis.</span>
        </div>
      </div>

      <!-- Actions -->
      <div class="p-6 pt-4 flex items-center justify-end gap-2">
        <button
          type="button"
          onclick={onClose}
          class="px-4 py-2 text-xs font-semibold rounded-xl text-on-surface-variant hover:bg-surface-container transition-colors cursor-pointer"
          disabled={isProcessing}
        >
          Batal
        </button>
        <button
          type="button"
          onclick={onConfirm}
          class="px-4 py-2 text-xs font-semibold rounded-xl {schema.is_active ? 'bg-amber-600 hover:bg-amber-700 text-white' : 'bg-emerald-600 hover:bg-emerald-700 text-white'} transition-colors flex items-center gap-1.5 cursor-pointer shadow-xs"
          disabled={isProcessing}
        >
          {#if isProcessing}
            <span class="loading loading-spinner loading-xs"></span>
            <span>Memproses...</span>
          {:else}
            <Power class="w-3.5 h-3.5" />
            <span>{schema.is_active ? 'Nonaktifkan' : 'Aktifkan Kembali'}</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
