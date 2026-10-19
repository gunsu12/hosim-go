<script lang="ts">
  import {
    FileText,
    Building2,
    Calendar,
    Layers,
    RotateCw,
    Edit3,
    FileSignature,
    CheckCircle2,
    PauseCircle,
    XCircle,
  } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { CustomerContractRecord } from '$lib/types/finance/customer_contract';

  interface Props {
    isOpen: boolean;
    contract?: CustomerContractRecord | null;
    onClose: () => void;
    onEdit?: (contract: CustomerContractRecord) => void;
    onAmend?: (contract: CustomerContractRecord) => void;
    onManageSchemas?: (contract: CustomerContractRecord) => void;
    onStatusTransition?: (
      contract: CustomerContractRecord,
      action: 'ACTIVATE' | 'SUSPEND' | 'REACTIVATE' | 'TERMINATE',
    ) => void;
    onRenew?: (contract: CustomerContractRecord) => void;
  }

  let {
    isOpen = false,
    contract = null,
    onClose,
    onEdit = () => {},
    onAmend = () => {},
    onManageSchemas = () => {},
    onStatusTransition = () => {},
    onRenew = () => {},
  }: Props = $props();

  const isExpired = $derived.by(() => {
    if (!contract?.expiration_date) return false;
    const today = new Date().toISOString().split('T')[0];
    return contract.expiration_date < today;
  });
</script>

<M3Modal
  {isOpen}
  title="Detail Kontrak PKS"
  subtitle={contract ? `${contract.contract_number} • ${contract.name}` : ''}
  icon={FileText}
  iconClass="bg-[#0b57d0]/10 text-[#0b57d0] dark:bg-[#a8c7fa]/15 dark:text-[#a8c7fa]"
  size="xl"
  {onClose}
>
  {#if contract}
    <div class="space-y-3.5 text-xs">
      <!-- Title & Status Banner -->
      <div
        class="p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-between gap-3"
      >
        <div class="min-w-0">
          <span class="font-mono text-sm font-bold text-[#0b57d0] dark:text-[#a8c7fa] block truncate">
            {contract.contract_number}
          </span>
          <h4 class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] truncate mt-0.5">
            {contract.name}
          </h4>
        </div>
        <div class="flex items-center gap-1.5 shrink-0">
          <span
            class="px-2.5 py-1 rounded-full font-bold text-[10px] uppercase {contract.status === 'ACTIVE'
              ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300'
              : contract.status === 'DRAFT'
                ? 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300'
                : contract.status === 'SUSPENDED'
                  ? 'bg-orange-100 text-orange-800 dark:bg-orange-950/60 dark:text-orange-300'
                  : 'bg-rose-100 text-rose-800 dark:bg-rose-950/60 dark:text-rose-300'}"
          >
            {contract.status}
          </span>
          {#if isExpired}
            <span
              class="px-2 py-0.5 rounded-full font-bold text-[10px] bg-rose-100 text-rose-700 dark:bg-rose-950/60 dark:text-rose-300"
            >
              Kadaluarsa
            </span>
          {/if}
        </div>
      </div>

      <!-- Info Grid (Clean Key-Value) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <!-- Penjamin -->
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Penjamin / Debitur</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5 truncate">
            {contract.customer_name || contract.customer?.name || '-'}
            {#if contract.customer_code}
              <span class="font-mono text-[11px] text-[#747775]">({contract.customer_code})</span>
            {/if}
          </span>
        </div>

        <!-- Masa Berlaku -->
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Masa Berlaku</span>
          <span class="font-mono font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5">
            {contract.effective_date} s/d {contract.expiration_date || 'Permanen'}
          </span>
        </div>

        <!-- Termin Pembayaran -->
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Termin Pembayaran</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5">
            Net {contract.payment_term_days} Hari
          </span>
          {#if contract.payment_term_notes}
            <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block mt-0.5">
              {contract.payment_term_notes}
            </span>
          {/if}
        </div>

        <!-- Batas Kredit -->
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Batas Kredit / Plafon</span>
          {#if contract.credit_limit !== null && contract.credit_limit !== undefined}
            <span class="font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa] block mt-0.5">
              Rp {contract.credit_limit.toLocaleString('id-ID')}
            </span>
          {:else}
            <span class="font-medium text-emerald-600 dark:text-emerald-400 block mt-0.5">
              Bebas Plafon (Unlimited)
            </span>
          {/if}
        </div>

        <!-- Syarat Jaminan -->
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Syarat Jaminan Kunjungan</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5">
            {contract.guarantee_requirement}
          </span>
          {#if contract.guarantee_notes}
            <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block mt-0.5">
              {contract.guarantee_notes}
            </span>
          {/if}
        </div>

        <!-- Perpanjangan Asal (jika ada) -->
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Asal Kontrak</span>
          {#if contract.renewed_from_contract_number}
            <span class="font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa] block mt-0.5">
              Perpanjangan dari: {contract.renewed_from_contract_number}
            </span>
          {:else}
            <span class="text-[#747775] dark:text-[#8e918f] block mt-0.5">Kontrak Baru (Induk)</span>
          {/if}
        </div>
      </div>

      <!-- Skema Tarif Terhubung -->
      <div
        class="p-3.5 rounded-2xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20] space-y-2"
      >
        <div class="flex items-center justify-between">
          <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs">
            Skema Tarif Terhubung ({contract.price_schemas?.length || 0})
          </span>
          <button
            type="button"
            onclick={() => { onClose(); onManageSchemas(contract!); }}
            class="text-[#0b57d0] dark:text-[#a8c7fa] hover:underline font-semibold text-[11px] cursor-pointer"
          >
            Atur Skema
          </button>
        </div>

        {#if !contract.price_schemas || contract.price_schemas.length === 0}
          <p class="text-[11px] text-[#747775] dark:text-[#8e918f]">Belum ada skema tarif terhubung.</p>
        {:else}
          <div class="space-y-1.5">
            {#each contract.price_schemas as s}
              <div
                class="flex items-center justify-between p-2 rounded-xl bg-[#f8fafd] dark:bg-[#28292a] text-xs"
              >
                <div class="flex items-center gap-2 min-w-0">
                  <span
                    class="w-5 h-5 rounded-full flex items-center justify-center font-bold text-[10px] {s.priority === 1
                      ? 'bg-amber-100 text-amber-800'
                      : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300'}"
                  >
                    {s.priority}
                  </span>
                  <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
                    {s.price_schema_name || s.price_schema_code || 'Skema Tarif'}
                  </span>
                  {#if s.priority === 1}
                    <span class="text-[9px] font-bold text-amber-700 dark:text-amber-400">
                      [Default]
                    </span>
                  {/if}
                </div>
                <span
                  class="text-[10px] font-bold {s.is_active
                    ? 'text-emerald-700 dark:text-emerald-400'
                    : 'text-gray-400'}"
                >
                  {s.is_active ? 'Aktif' : 'Non-aktif'}
                </span>
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Catatan Khusus (jika ada) -->
      {#if contract.notes}
        <div
          class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]"
        >
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block font-semibold mb-0.5">Catatan:</span>
          <p class="text-[#1f1f1f] dark:text-[#e3e3e3] whitespace-pre-wrap">{contract.notes}</p>
        </div>
      {/if}

      <!-- Alasan Transisi (jika ada) -->
      {#if contract.status_reason}
        <div
          class="p-3 rounded-xl bg-amber-50 dark:bg-amber-950/40 border border-amber-200 dark:border-amber-900/50 text-amber-900 dark:text-amber-200"
        >
          <span class="font-bold block text-[11px]">Alasan Perubahan Status:</span>
          <p class="mt-0.5">{contract.status_reason}</p>
        </div>
      {/if}

      <!-- Audit Metadata Footer -->
      <div
        class="pt-2 text-[11px] text-[#747775] dark:text-[#8e918f] flex flex-wrap justify-between gap-2 border-t border-[#e1e5ea] dark:border-[#353738]"
      >
        <span>Dibuat: {contract.created_by}</span>
        <span>ID: <strong class="font-mono">{contract.id}</strong></span>
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <div class="flex items-center justify-between w-full">
      <M3Button variant="text" onclick={onClose}>
        Tutup
      </M3Button>

      <div class="flex items-center gap-2">
        {#if contract}
          {#if contract.status === 'DRAFT'}
            <M3Button variant="tonal" onclick={() => { onClose(); onEdit(contract!); }}>
              <Edit3 class="w-4 h-4 mr-1" />
              <span>Edit Draft</span>
            </M3Button>
            <M3Button
              variant="filled"
              onclick={() => { onClose(); onStatusTransition(contract!, 'ACTIVATE'); }}
            >
              <CheckCircle2 class="w-4 h-4 mr-1" />
              <span>Aktivasi</span>
            </M3Button>
          {:else if contract.status === 'ACTIVE'}
            <M3Button variant="tonal" onclick={() => { onClose(); onAmend(contract!); }}>
              <FileSignature class="w-4 h-4 mr-1" />
              <span>Amandemen</span>
            </M3Button>
            <M3Button variant="tonal" onclick={() => { onClose(); onRenew(contract!); }}>
              <RotateCw class="w-4 h-4 mr-1" />
              <span>Perpanjang</span>
            </M3Button>
            <M3Button
              variant="outlined"
              onclick={() => { onClose(); onStatusTransition(contract!, 'SUSPEND'); }}
            >
              <PauseCircle class="w-4 h-4 mr-1 text-amber-600" />
              <span>Suspend</span>
            </M3Button>
          {:else if contract.status === 'SUSPENDED'}
            <M3Button
              variant="filled"
              onclick={() => { onClose(); onStatusTransition(contract!, 'REACTIVATE'); }}
            >
              <CheckCircle2 class="w-4 h-4 mr-1" />
              <span>Aktifkan Kembali</span>
            </M3Button>
            <M3Button
              variant="outlined"
              onclick={() => { onClose(); onStatusTransition(contract!, 'TERMINATE'); }}
            >
              <XCircle class="w-4 h-4 mr-1 text-rose-600" />
              <span>Putus Kontrak</span>
            </M3Button>
          {:else if contract.status === 'TERMINATED'}
            <M3Button variant="filled" onclick={() => { onClose(); onRenew(contract!); }}>
              <RotateCw class="w-4 h-4 mr-1" />
              <span>Perpanjang (Renew)</span>
            </M3Button>
          {/if}
        {/if}
      </div>
    </div>
  {/snippet}
</M3Modal>
