<script lang="ts">
  import {
    X,
    Check,
    AlertCircle,
    Building2,
    Hospital,
    Activity,
    ShieldCheck,
    Search,
    ChevronDown,
    ChevronUp
  } from '@lucide/svelte';
  import type {
    PriceSchemaTreeNode,
    UpsertRuleDTO,
    AdjustmentType
  } from '$lib/types/finance/price_schema';
  import type { ItemSummaryRecord } from '$lib/types/master/item';

  interface Props {
    isOpen: boolean;
    targetNode: PriceSchemaTreeNode | null;
    isSaving: boolean;
    availableItems: ItemSummaryRecord[];
    onClose: () => void;
    onSave: (payload: UpsertRuleDTO) => Promise<void>;
  }

  let {
    isOpen,
    targetNode,
    isSaving,
    availableItems,
    onClose,
    onSave
  }: Props = $props();

  // Form State
  let adjType = $state<AdjustmentType>('NONE');
  let adjValue = $state<number>(0);
  let patientPercent = $state<number>(0);
  let customerPercent = $state<number>(100);
  let taxPercent = $state<number>(11);

  // Selected item if adding a new Item rule under Unit
  let selectedItemId = $state<string | null>(null);
  let itemSearch = $state('');

  // Medication (Fase 2)
  let showMedicationSection = $state(false);
  let medAdjType = $state<AdjustmentType>('NONE');
  let medAdjValue = $state<number>(0);

  let validationError = $state<string | null>(null);

  // Filtered items for autocomplete
  const filteredItems = $derived.by(() => {
    if (!itemSearch.trim()) return availableItems.slice(0, 20);
    const q = itemSearch.toLowerCase();
    return availableItems.filter(
      it => it.name.toLowerCase().includes(q) || it.code.toLowerCase().includes(q)
    ).slice(0, 25);
  });

  $effect(() => {
    if (isOpen && targetNode) {
      validationError = null;
      adjType = targetNode.tariffAdjustmentType || 'NONE';
      adjValue = Number(targetNode.tariffAdjustmentValue) || 0;
      patientPercent = Number(targetNode.payByPatientPercent) || 0;
      customerPercent = Number(targetNode.payByCustomerPercent) || 100;
      taxPercent = Number(targetNode.taxPercent) || 11;

      medAdjType = targetNode.medicationAdjustmentType || 'NONE';
      medAdjValue = Number(targetNode.medicationAdjustmentValue) || 0;

      selectedItemId = targetNode.itemId || null;
      itemSearch = '';
    }
  });

  function handlePatientPctChange(val: number) {
    patientPercent = Math.max(0, Math.min(100, val || 0));
    customerPercent = 100 - patientPercent;
  }

  function handleCustomerPctChange(val: number) {
    customerPercent = Math.max(0, Math.min(100, val || 0));
    patientPercent = 100 - customerPercent;
  }

  function applyCostSharePreset(pat: number, cust: number) {
    patientPercent = pat;
    customerPercent = cust;
  }

  const isItemLevel = $derived(
    targetNode?.nodeType === 'ITEM' ||
    (targetNode?.nodeType === 'UNIT' && selectedItemId !== null)
  );

  async function handleSubmit(e?: Event) {
    e?.preventDefault();
    validationError = null;

    if (!targetNode) return;

    if (patientPercent + customerPercent !== 100) {
      validationError = 'Total pembagian tanggungan (Pasien + Penjamin) harus tepat 100%';
      return;
    }

    if (adjType === 'PERCENT' && adjValue < -100) {
      validationError = 'Nilai diskon persentase minimal adalah -100%';
      return;
    }

    if (adjType === 'FIXED_PRICE' && !isItemLevel) {
      validationError = 'Harga Tetap (FIXED_PRICE) hanya dapat diterapkan pada tingkat Item Tindakan';
      return;
    }

    if (adjType === 'FIXED_PRICE' && adjValue < 0) {
      validationError = 'Harga Tetap (FIXED_PRICE) tidak boleh bernilai negatif';
      return;
    }

    // Determine IDs according to scope
    let deptId: string | null = null;
    let unitId: string | null = null;
    let itemId: string | null = null;

    if (targetNode.scopeType === 'SCHEMA_DEFAULT') {
      deptId = null;
      unitId = null;
      itemId = null;
    } else if (targetNode.nodeType === 'DEPT') {
      deptId = targetNode.departementId || null;
    } else if (targetNode.nodeType === 'UNIT') {
      if (selectedItemId) {
        // Adding item under unit
        unitId = targetNode.serviceUnitId || null;
        itemId = selectedItemId;
      } else {
        unitId = targetNode.serviceUnitId || null;
      }
    } else if (targetNode.nodeType === 'ITEM') {
      deptId = targetNode.departementId || null;
      unitId = targetNode.serviceUnitId || null;
      itemId = targetNode.itemId || null;
    }

    const payload: UpsertRuleDTO = {
      departement_id: deptId,
      service_unit_id: unitId,
      item_id: itemId,
      tariff_adjustment_type: adjType,
      tariff_adjustment_value: Number(adjValue) || 0,
      medication_adjustment_type: medAdjType,
      medication_adjustment_value: Number(medAdjValue) || 0,
      pay_by_patient_percent: patientPercent,
      pay_by_customer_percent: customerPercent,
      tax_percent: taxPercent
    };

    await onSave(payload);
  }
</script>

{#if isOpen && targetNode}
  <!-- Backdrop -->
  <div class="fixed inset-0 z-50 overflow-hidden bg-black/40 backdrop-blur-xs flex justify-end animate-in fade-in duration-200">
    <!-- Slide-over Drawer Canvas -->
    <div
      class="bg-white dark:bg-surface-container w-full max-w-md h-full shadow-2xl flex flex-col border-l border-outline-variant/60 animate-in slide-in-from-right duration-250"
    >
      <!-- Drawer Header -->
      <div class="px-6 py-4 bg-[#f8fafd] dark:bg-surface-container-high border-b border-outline-variant/60 flex items-center justify-between">
        <div>
          <h3 class="text-sm font-bold text-on-surface">
            {targetNode.isInherited ? 'Tetapkan Aturan Override' : 'Ubah Aturan Penyesuaian'}
          </h3>
          <p class="text-xs text-on-surface-variant mt-0.5">
            Konfigurasi penyesuaian tarif, tanggungan, dan pajak
          </p>
        </div>
        <button
          type="button"
          onclick={onClose}
          class="p-1.5 rounded-full text-on-surface-variant hover:text-on-surface hover:bg-surface-container cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Drawer Form Content -->
      <form onsubmit={handleSubmit} class="flex-1 overflow-y-auto p-6 flex flex-col gap-4 text-xs">
        {#if validationError}
          <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800/60 text-rose-700 dark:text-rose-300 flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{validationError}</span>
          </div>
        {/if}

        <!-- Scope Context Info Card -->
        <div class="p-4 rounded-2xl bg-surface-container-high/60 border border-outline-variant/60 flex items-center gap-3">
          <div class="w-9 h-9 rounded-xl bg-primary text-on-primary flex items-center justify-center shrink-0">
            {#if targetNode.nodeType === 'DEPT'}
              <Building2 class="w-4 h-4" />
            {:else if targetNode.nodeType === 'UNIT'}
              <Hospital class="w-4 h-4" />
            {:else if targetNode.nodeType === 'ITEM'}
              <Activity class="w-4 h-4" />
            {:else}
              <ShieldCheck class="w-4 h-4" />
            {/if}
          </div>
          <div class="min-w-0">
            <div class="text-[11px] font-bold text-primary uppercase">
              Target Scope: {targetNode.nodeType}
            </div>
            <div class="text-sm font-bold text-on-surface truncate">
              {targetNode.name}
            </div>
            {#if targetNode.parentRuleSummary}
              <div class="text-[10px] text-outline mt-0.5">
                {targetNode.parentRuleSummary}
              </div>
            {/if}
          </div>
        </div>

        <!-- Optional Item Selection (if user clicked "+ Item" on a Service Unit) -->
        {#if targetNode.nodeType === 'UNIT'}
          <div class="flex flex-col gap-1.5 pt-1">
            <label for="item_search_input" class="font-semibold text-on-surface flex items-center justify-between">
              <span>Item Tindakan (Opsional)</span>
              <span class="text-[11px] font-normal text-on-surface-variant">Kosongkan jika untuk unit secara keseluruhan</span>
            </label>

            {#if selectedItemId}
              {@const selectedItem = availableItems.find(i => i.id === selectedItemId)}
              <div class="p-2.5 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-between">
                <div class="min-w-0">
                  <div class="font-bold text-primary truncate">{selectedItem?.name || selectedItemId}</div>
                  <div class="font-mono text-[10px] text-primary/80">{selectedItem?.code || ''}</div>
                </div>
                <button
                  type="button"
                  onclick={() => (selectedItemId = null)}
                  class="p-1 text-primary hover:bg-primary/20 rounded-md cursor-pointer"
                  title="Batalkan pilihan item"
                >
                  <X class="w-3.5 h-3.5" />
                </button>
              </div>
            {:else}
              <div class="relative">
                <Search class="w-3.5 h-3.5 absolute left-3 top-2.5 text-on-surface-variant" />
                <input
                  id="item_search_input"
                  type="text"
                  bind:value={itemSearch}
                  placeholder="Cari item tindakan spesifik di unit ini..."
                  class="w-full pl-8 pr-3 py-2 rounded-xl border border-outline-variant bg-surface text-xs focus:outline-hidden focus:border-primary"
                />
              </div>

              {#if itemSearch.trim() && filteredItems.length > 0}
                <div class="max-h-36 overflow-y-auto rounded-xl border border-outline-variant/60 bg-surface divide-y divide-outline-variant/30 shadow-xs">
                  {#each filteredItems as item (item.id)}
                    <button
                      type="button"
                      onclick={() => {
                        selectedItemId = item.id;
                        itemSearch = '';
                      }}
                      class="w-full px-3 py-2 text-left hover:bg-surface-container flex items-center justify-between cursor-pointer"
                    >
                      <span class="truncate font-medium text-on-surface">{item.name}</span>
                      <span class="font-mono text-[10px] text-outline shrink-0 ml-2">{item.code}</span>
                    </button>
                  {/each}
                </div>
              {/if}
            {/if}
          </div>
        {/if}

        <!-- 1. Tipe Penyesuaian Tarif -->
        <div class="flex flex-col gap-2 pt-2 border-t border-outline-variant/40">
          <span class="font-semibold text-on-surface">Tipe Penyesuaian Tarif</span>
          <div class="grid grid-cols-2 gap-2">
            <!-- NONE -->
            <button
              type="button"
              onclick={() => { adjType = 'NONE'; adjValue = 0; }}
              class="p-2.5 rounded-xl border text-left transition-all cursor-pointer {adjType === 'NONE' ? 'border-primary bg-primary/10 text-primary font-bold' : 'border-outline-variant hover:bg-surface-container text-on-surface'}"
            >
              <div class="text-xs">NONE</div>
              <div class="text-[10px] font-normal opacity-75">Buku Tarif Normal</div>
            </button>

            <!-- PERCENT -->
            <button
              type="button"
              onclick={() => (adjType = 'PERCENT')}
              class="p-2.5 rounded-xl border text-left transition-all cursor-pointer {adjType === 'PERCENT' ? 'border-primary bg-primary/10 text-primary font-bold' : 'border-outline-variant hover:bg-surface-container text-on-surface'}"
            >
              <div class="text-xs">PERCENT (%)</div>
              <div class="text-[10px] font-normal opacity-75">Markup (+) atau Diskon (-)</div>
            </button>

            <!-- FIXED_AMOUNT -->
            <button
              type="button"
              onclick={() => (adjType = 'FIXED_AMOUNT')}
              class="p-2.5 rounded-xl border text-left transition-all cursor-pointer {adjType === 'FIXED_AMOUNT' ? 'border-primary bg-primary/10 text-primary font-bold' : 'border-outline-variant hover:bg-surface-container text-on-surface'}"
            >
              <div class="text-xs">FIXED_AMOUNT (Rp)</div>
              <div class="text-[10px] font-normal opacity-75">Nominal Tetap (+/- Rp)</div>
            </button>

            <!-- FIXED_PRICE -->
            <button
              type="button"
              disabled={!isItemLevel}
              onclick={() => (adjType = 'FIXED_PRICE')}
              class="p-2.5 rounded-xl border text-left transition-all cursor-pointer {!isItemLevel ? 'opacity-40 cursor-not-allowed bg-surface-container/50 border-outline-variant/40' : adjType === 'FIXED_PRICE' ? 'border-primary bg-primary/10 text-primary font-bold' : 'border-outline-variant hover:bg-surface-container text-on-surface'}"
              title={!isItemLevel ? 'Harga Tetap hanya berlaku pada Item Tindakan' : ''}
            >
              <div class="text-xs">FIXED_PRICE (Rp)</div>
              <div class="text-[10px] font-normal opacity-75">Harga PKS Khusus Item</div>
            </button>
          </div>
        </div>

        <!-- 2. Nilai Penyesuaian Tarif -->
        {#if adjType !== 'NONE'}
          <div class="flex flex-col gap-1.5">
            <label for="adj_val_input" class="font-semibold text-on-surface flex items-center justify-between">
              <span>Nilai Penyesuaian</span>
              <span class="text-[11px] text-on-surface-variant">
                {adjType === 'PERCENT' ? 'Gunakan minus (-) untuk diskon' : 'Nominal Rupiah'}
              </span>
            </label>
            <div class="relative">
              {#if adjType === 'FIXED_AMOUNT' || adjType === 'FIXED_PRICE'}
                <span class="absolute left-3.5 top-2.5 font-bold text-on-surface-variant">Rp</span>
              {/if}
              <input
                id="adj_val_input"
                type="number"
                step={adjType === 'PERCENT' ? '0.01' : '1'}
                bind:value={adjValue}
                placeholder={adjType === 'PERCENT' ? '-10 untuk Diskon 10%' : '150000'}
                class="w-full {adjType === 'FIXED_AMOUNT' || adjType === 'FIXED_PRICE' ? 'pl-10' : 'pl-3.5'} pr-10 py-2.5 rounded-xl border border-outline-variant bg-surface text-on-surface font-semibold text-xs focus:outline-hidden focus:border-primary"
                required
              />
              {#if adjType === 'PERCENT'}
                <span class="absolute right-3.5 top-2.5 font-bold text-on-surface-variant">%</span>
              {/if}
            </div>
          </div>
        {/if}

        <!-- 3. Cost Sharing (Pasien vs Penjamin) -->
        <div class="flex flex-col gap-2 pt-2 border-t border-outline-variant/40">
          <div class="flex items-center justify-between">
            <span class="font-semibold text-on-surface">Porsi Tanggungan Tagihan</span>
            <div class="flex items-center gap-1">
              <button
                type="button"
                onclick={() => applyCostSharePreset(0, 100)}
                class="px-2 py-0.5 rounded text-[10px] font-semibold bg-surface-container hover:bg-surface-container-high border border-outline-variant/60 cursor-pointer"
              >
                100% Penjamin
              </button>
              <button
                type="button"
                onclick={() => applyCostSharePreset(100, 0)}
                class="px-2 py-0.5 rounded text-[10px] font-semibold bg-surface-container hover:bg-surface-container-high border border-outline-variant/60 cursor-pointer"
              >
                100% Pasien
              </button>
              <button
                type="button"
                onclick={() => applyCostSharePreset(20, 80)}
                class="px-2 py-0.5 rounded text-[10px] font-semibold bg-surface-container hover:bg-surface-container-high border border-outline-variant/60 cursor-pointer"
              >
                Co-Pay 20:80
              </button>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="flex flex-col gap-1">
              <label for="drawer_pat_pct" class="text-[11px] font-medium text-on-surface-variant">Pasien (%)</label>
              <div class="relative">
                <input
                  id="drawer_pat_pct"
                  type="number"
                  min="0"
                  max="100"
                  value={patientPercent}
                  oninput={(e) => handlePatientPctChange(Number(e.currentTarget.value))}
                  class="w-full px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface font-semibold text-xs"
                />
                <span class="absolute right-3 top-2 text-on-surface-variant font-bold">%</span>
              </div>
            </div>

            <div class="flex flex-col gap-1">
              <label for="drawer_cust_pct" class="text-[11px] font-medium text-on-surface-variant">Penjamin / Rekanan (%)</label>
              <div class="relative">
                <input
                  id="drawer_cust_pct"
                  type="number"
                  min="0"
                  max="100"
                  value={customerPercent}
                  oninput={(e) => handleCustomerPctChange(Number(e.currentTarget.value))}
                  class="w-full px-3 py-2 rounded-xl border border-outline-variant bg-surface text-on-surface font-semibold text-xs"
                />
                <span class="absolute right-3 top-2 text-on-surface-variant font-bold">%</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 4. Pajak (PPN %) -->
        <div class="flex items-center justify-between pt-2 border-t border-outline-variant/40">
          <div>
            <div class="font-semibold text-on-surface">Pajak Penambahan Nilai (PPN)</div>
            <div class="text-[11px] text-on-surface-variant">Dihitung dari tarif setelah penyesuaian</div>
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              onclick={() => (taxPercent = 11)}
              class="px-2 py-1 rounded-lg text-[10px] font-bold {taxPercent === 11 ? 'bg-primary text-on-primary' : 'bg-surface-container text-on-surface-variant'} cursor-pointer"
            >
              11%
            </button>
            <button
              type="button"
              onclick={() => (taxPercent = 0)}
              class="px-2 py-1 rounded-lg text-[10px] font-bold {taxPercent === 0 ? 'bg-primary text-on-primary' : 'bg-surface-container text-on-surface-variant'} cursor-pointer"
            >
              0%
            </button>
            <input
              id="drawer_tax_pct"
              type="number"
              min="0"
              max="100"
              bind:value={taxPercent}
              class="w-14 px-2 py-1 text-right rounded-lg border border-outline-variant bg-surface text-xs font-semibold"
            />
            <span class="font-bold text-xs">%</span>
          </div>
        </div>

        <!-- 5. Collapsible: Farmasi & Obat (Fase 2) -->
        <div class="pt-2 border-t border-outline-variant/40">
          <button
            type="button"
            onclick={() => (showMedicationSection = !showMedicationSection)}
            class="w-full py-1.5 flex items-center justify-between text-on-surface-variant hover:text-on-surface cursor-pointer"
          >
            <div class="flex items-center gap-1.5 font-semibold text-[11px]">
              <span>Pengaturan Farmasi & Obat</span>
              <span class="px-1.5 py-0.2 rounded text-[9px] font-bold bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300">
                Fase 2
              </span>
            </div>
            {#if showMedicationSection}
              <ChevronUp class="w-3.5 h-3.5" />
            {:else}
              <ChevronDown class="w-3.5 h-3.5" />
            {/if}
          </button>

          {#if showMedicationSection}
            <div class="p-3 mt-1.5 rounded-xl bg-surface-container-high/50 border border-outline-variant/40 flex flex-col gap-2">
              <p class="text-[10px] text-on-surface-variant">
                Kolom penyesuaian obat telah disiapkan pada basis data, namun kalkulasinya belum aktif di engine fase 1.
              </p>
              <div class="grid grid-cols-2 gap-2">
                <select
                  bind:value={medAdjType}
                  class="px-2 py-1.5 rounded-lg border border-outline-variant bg-surface text-xs"
                >
                  <option value="NONE">NONE (0%)</option>
                  <option value="PERCENT">PERCENT (%)</option>
                </select>
                <input
                  type="number"
                  bind:value={medAdjValue}
                  placeholder="0"
                  class="px-2 py-1.5 rounded-lg border border-outline-variant bg-surface text-xs"
                />
              </div>
            </div>
          {/if}
        </div>
      </form>

      <!-- Drawer Footer Actions -->
      <div class="p-4 px-6 bg-[#f8fafd] dark:bg-surface-container-high border-t border-outline-variant/60 flex items-center justify-end gap-2 shrink-0">
        <button
          type="button"
          onclick={onClose}
          class="px-4 py-2 text-xs font-semibold rounded-xl text-on-surface-variant hover:bg-surface-container transition-colors cursor-pointer"
          disabled={isSaving}
        >
          Batal
        </button>
        <button
          type="button"
          onclick={handleSubmit}
          class="px-5 py-2 text-xs font-semibold rounded-xl bg-primary text-on-primary hover:bg-primary/90 transition-colors flex items-center gap-1.5 cursor-pointer shadow-xs"
          disabled={isSaving}
        >
          {#if isSaving}
            <span class="loading loading-spinner loading-xs"></span>
            <span>Menyimpan...</span>
          {:else}
            <Check class="w-3.5 h-3.5" />
            <span>Simpan Aturan</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
