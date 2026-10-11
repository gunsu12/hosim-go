<script lang="ts">
  import { onMount } from 'svelte';
  import { Copy, CheckCircle2 } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3ToggleCard from '$lib/components/m3/M3ToggleCard.svelte';
  import M3DatePicker from '$lib/components/m3/M3DatePicker.svelte';
  import { getCustomers } from '$lib/api/finance/customer';
  import type { CustomerRecord } from '$lib/types/master/customer';
  import type {
    PricePlanRecord,
    ClonePricePlanDTO
  } from '$lib/types/finance/price_plan';

  interface Props {
    isOpen: boolean;
    sourcePlan: PricePlanRecord | null;
    isCloning: boolean;
    onClose: () => void;
    onClone: (payload: ClonePricePlanDTO) => Promise<void>;
  }

  let {
    isOpen,
    sourcePlan,
    isCloning,
    onClose,
    onClone
  }: Props = $props();

  let newCode = $state('');
  let newName = $state('');
  let effectiveFrom = $state('');
  let effectiveTo = $state('');
  let defaultCitoPercent = $state(25);
  let isDefault = $state(false);
  let customerId = $state<string>('');
  let description = $state('');
  let formError = $state<string | null>(null);

  let customers = $state<CustomerRecord[]>([]);

  onMount(() => {
    getCustomers({ is_active: true, limit: 100 }).then((res) => {
      customers = res.data || [];
    }).catch((err) => {
      console.error('Failed to load customers for clone modal:', err);
    });
  });

  $effect(() => {
    if (isOpen && sourcePlan) {
      const today = new Date().toISOString().split('T')[0];
      newCode = `${sourcePlan.code}-COPY`;
      newName = `${sourcePlan.name} (Salinan)`;
      effectiveFrom = today;
      effectiveTo = '';
      defaultCitoPercent = sourcePlan.default_cito_percent || 25;
      isDefault = false;
      customerId = sourcePlan.customer_id || '';
      description = `Kloning dari acuan ${sourcePlan.code}`;
      formError = null;
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!newCode.trim()) {
      formError = 'Kode buku tarif baru wajib diisi';
      return;
    }
    if (!newName.trim()) {
      formError = 'Nama buku tarif baru wajib diisi';
      return;
    }
    if (!effectiveFrom) {
      formError = 'Tanggal berlaku awal wajib diisi';
      return;
    }

    formError = null;
    await onClone({
      new_code: newCode.trim().toUpperCase(),
      new_name: newName.trim(),
      effective_from: effectiveFrom,
      effective_to: effectiveTo || undefined,
      default_cito_percent: Number(defaultCitoPercent),
      is_default: isDefault,
      customer_id: customerId || undefined,
      description: description.trim() || undefined
    });
  }
</script>

<M3Modal
  {isOpen}
  title="Gandakan (Clone) Buku Tarif"
  subtitle="Menyalin seluruh matriks item dan pecahan komponen biaya secara atomik"
  icon={Copy}
  iconClass="bg-purple-500/10 text-purple-600 dark:bg-purple-500/20 dark:text-purple-400"
  size="xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if formError}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <span class="font-bold">Gagal:</span>
      <span>{formError}</span>
    </div>
  {/if}

  {#if sourcePlan}
    <!-- Source Plan Context Banner -->
    <div class="p-3.5 rounded-2xl bg-purple-50/70 dark:bg-purple-950/30 border border-purple-200 dark:border-purple-800/60 text-xs text-purple-900 dark:text-purple-200 flex flex-col gap-1.5">
      <span class="text-[11px] font-semibold text-purple-700 dark:text-purple-300">
        Buku Tarif Sumber (Acuan Kloning):
      </span>
      <div class="flex items-center gap-2 font-bold">
        <span class="font-mono bg-white dark:bg-[#1e1f20] px-2 py-0.5 rounded border border-purple-200 dark:border-purple-700 text-purple-800 dark:text-purple-300 text-xs">
          {sourcePlan.code}
        </span>
        <span class="truncate">{sourcePlan.name}</span>
      </div>
      <span class="text-[11px] text-[#444746] dark:text-[#c4c7c5]">
        Status: <strong>{sourcePlan.status}</strong> • Total Item: <strong>{sourcePlan.total_items} tindakan</strong>
      </span>
    </div>
  {/if}

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="pp-new-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Kode Buku Tarif Baru <span class="text-rose-500">*</span>
      </label>
      <input
        id="pp-new-code"
        type="text"
        bind:value={newCode}
        placeholder="Misal: TPP-2027-DRAFT"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div>
      <label for="pp-new-cito" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Standar CITO (% Tambahan) <span class="text-rose-500">*</span>
      </label>
      <div class="relative">
        <input
          id="pp-new-cito"
          type="number"
          step="any"
          min="0"
          max="100"
          bind:value={defaultCitoPercent}
          required
          class="w-full h-10 px-3.5 pr-8 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
        <span class="absolute right-3 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f] font-bold">%</span>
      </div>
    </div>
  </div>

  <div>
    <label for="pp-new-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Nama Buku Tarif Baru <span class="text-rose-500">*</span>
    </label>
    <input
      id="pp-new-name"
      type="text"
      bind:value={newName}
      placeholder="Misal: Buku Tarif Layanan RS 2027"
      required
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    />
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <M3DatePicker
      id="pp-clone-eff-from"
      label="Berlaku Mulai *"
      bind:value={effectiveFrom}
      required
      placeholder="Pilih tanggal mulai"
    />

    <M3DatePicker
      id="pp-clone-eff-to"
      label="Berlaku Sampai (Opsional)"
      bind:value={effectiveTo}
      min={effectiveFrom || ''}
      placeholder="Pilih tanggal selesai"
    />
  </div>

  <div>
    <label for="pp-clone-customer" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Penjamin Rekanan (Opsional)
    </label>
    <select
      id="pp-clone-customer"
      bind:value={customerId}
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    >
      <option value="">-- Standar Umum (Berlaku untuk Semua Penjamin) --</option>
      {#each customers as cust (cust.id)}
        <option value={cust.id}>[{cust.code}] {cust.name}</option>
      {/each}
    </select>
  </div>

  <div>
    <label for="pp-clone-desc" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
      Catatan Kloning
    </label>
    <textarea
      id="pp-clone-desc"
      bind:value={description}
      rows="2"
      placeholder="Keterangan draf hasil duplikasi..."
      class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
    ></textarea>
  </div>

  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
    <M3ToggleCard
      title="Jadikan Buku Tarif Standar RS (Default)"
      description="Jika diaktifkan, buku tarif salinan ini akan menggantikan buku default lama."
      bind:checked={isDefault}
      colorVariant="primary"
    />
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isCloning}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isCloning}>
      {#if isCloning}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menggandakan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>Kloning Buku Tarif</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
