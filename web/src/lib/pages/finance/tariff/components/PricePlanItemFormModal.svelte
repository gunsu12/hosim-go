<script lang="ts">
  import {
    Plus,
    Trash2,
    CheckCircle2,
    AlertTriangle,
    Calculator
  } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select from '$lib/components/m3/M3Select.svelte';
  import {
    addPricePlanItem,
    updatePricePlanItem
  } from '$lib/api/finance/price_plan';
  import type {
    PricePlanItemRecord,
    ItemComponentDTO
  } from '$lib/types/finance/price_plan';
  import type { TariffClassRecord, TariffComponentRecord } from '$lib/types/finance/tariff';
  import type { ItemSummaryRecord } from '$lib/types/master/item';

  interface Props {
    isOpen: boolean;
    mode: 'ADD' | 'EDIT';
    planId: string;
    itemToEdit?: PricePlanItemRecord | null;
    tariffClasses?: TariffClassRecord[];
    tariffComponents?: TariffComponentRecord[];
    tariffCatalogItems?: ItemSummaryRecord[];
    onClose: () => void;
    onSuccess: (msg: string) => void;
  }

  let {
    isOpen,
    mode = 'ADD',
    planId,
    itemToEdit = null,
    tariffClasses = [],
    tariffComponents = [],
    tariffCatalogItems = [],
    onClose,
    onSuccess
  }: Props = $props();

  let selectedCatalogItemId = $state('');
  let selectedClassId = $state('');
  let formBasePrice = $state<number>(0);
  let formCitoPrice = $state<number | null>(null);
  let formComponents = $state<ItemComponentDTO[]>([]);
  let isSavingItem = $state(false);
  let formError = $state<string | null>(null);

  let catalogItemOptions = $derived([
    { value: '', label: '-- Pilih Tindakan Medis --' },
    ...tariffCatalogItems.map(i => ({ value: i.id, label: `[${i.code}] ${i.name}` }))
  ]);

  let tariffClassOptions = $derived(
    tariffClasses.map(tc => ({ value: tc.id, label: `[${tc.code}] ${tc.name}` }))
  );

  let tariffComponentOptions = $derived(
    tariffComponents.map(tc => ({ value: tc.id, label: `[${tc.code}] ${tc.name}` }))
  );

  // Balancing Calculation
  let sumComponentBase = $derived(
    formComponents.reduce((acc, c) => acc + (Number(c.base_amount) || 0), 0)
  );
  let balanceDiff = $derived(
    Math.round((Number(formBasePrice) - sumComponentBase) * 100) / 100
  );
  let isBalanced = $derived(balanceDiff === 0);

  $effect(() => {
    if (isOpen) {
      formError = null;
      if (mode === 'EDIT' && itemToEdit) {
        selectedCatalogItemId = itemToEdit.item_id;
        selectedClassId = itemToEdit.tariff_class_id;
        formBasePrice = itemToEdit.total_base_price;
        formCitoPrice = itemToEdit.total_cito_price || null;
        formComponents = itemToEdit.components.map((c) => ({
          component_id: c.component_id,
          base_amount: c.base_amount,
          cito_amount: c.cito_amount || null,
          coa_code: c.coa_code || null
        }));
      } else {
        selectedCatalogItemId = '';
        selectedClassId = tariffClasses[0]?.id || '';
        formBasePrice = 0;
        formCitoPrice = null;
        if (tariffComponents.length > 0) {
          formComponents = [{
            component_id: tariffComponents[0].id,
            base_amount: 0,
            cito_amount: null,
            coa_code: tariffComponents[0].default_coa_code || null
          }];
        } else {
          formComponents = [];
        }
      }
    }
  });

  function addComponentRow() {
    const firstComp = tariffComponents[0];
    formComponents = [
      ...formComponents,
      {
        component_id: firstComp?.id || '',
        base_amount: 0,
        cito_amount: null,
        coa_code: firstComp?.default_coa_code || null
      }
    ];
  }

  function removeComponentRow(idx: number) {
    formComponents = formComponents.filter((_, i) => i !== idx);
  }

  function handleAutoBalance() {
    if (formComponents.length === 0) return;
    const lastIdx = formComponents.length - 1;
    const currentLastAmount = formComponents[lastIdx].base_amount || 0;
    const adjusted = currentLastAmount + balanceDiff;
    if (adjusted >= 0) {
      formComponents[lastIdx].base_amount = adjusted;
    }
  }

  async function handleSaveItem() {
    if (!selectedCatalogItemId) {
      formError = 'Pilih tindakan medis dari katalog item tarif';
      return;
    }
    if (!selectedClassId) {
      formError = 'Pilih kelas tarif layanan';
      return;
    }
    if (Number(formBasePrice) <= 0) {
      formError = 'Total tarif dasar harus lebih dari 0';
      return;
    }
    if (formComponents.length === 0) {
      formError = 'Minimal satu rincian komponen biaya wajib disertakan';
      return;
    }
    if (!isBalanced) {
      formError = `Total pecahan komponen (Rp ${sumComponentBase.toLocaleString('id-ID')}) belum seimbang dengan total tarif dasar (Rp ${Number(formBasePrice).toLocaleString('id-ID')}). Selisih: Rp ${balanceDiff.toLocaleString('id-ID')}`;
      return;
    }

    formError = null;
    isSavingItem = true;
    try {
      if (mode === 'ADD') {
        await addPricePlanItem(planId, {
          item_id: selectedCatalogItemId,
          tariff_class_id: selectedClassId,
          total_base_price: Number(formBasePrice),
          total_cito_price: formCitoPrice ? Number(formCitoPrice) : undefined,
          is_active: true,
          components: formComponents
        });
        onSuccess('Item tarif berhasil ditambahkan ke buku tarif');
      } else if (itemToEdit) {
        await updatePricePlanItem(planId, itemToEdit.id, {
          total_base_price: Number(formBasePrice),
          total_cito_price: formCitoPrice ? Number(formCitoPrice) : undefined,
          is_active: true,
          components: formComponents
        });
        onSuccess('Item tarif berhasil diperbarui');
      }
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan item tarif';
    } finally {
      isSavingItem = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'ADD' ? 'Tambah Item Tindakan Tarif' : 'Edit Item & Pecahan Komponen Biaya'}
  subtitle="Wajib memenuhi invariant balancing: Total Komponen = Total Tarif Dasar"
  icon={Calculator}
  size="2xl"
  onsubmit={(e) => { e.preventDefault(); handleSaveItem(); }}
  {onClose}
>
  {#if formError}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <AlertTriangle class="w-4 h-4 shrink-0 mt-0.5" />
      <span>{formError}</span>
    </div>
  {/if}

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <M3Select
      label="Tindakan Medis (Katalog TARIFF) *"
      bind:value={selectedCatalogItemId}
      options={catalogItemOptions}
      disabled={mode === 'EDIT'}
      placeholder="Pilih Tindakan Medis"
    />

    <M3Select
      label="Kelas Layanan Perawatan *"
      bind:value={selectedClassId}
      options={tariffClassOptions}
      disabled={mode === 'EDIT'}
      placeholder="Pilih Kelas Layanan"
    />
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="f-base" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Total Tarif Dasar (Rp) <span class="text-rose-500">*</span>
      </label>
      <input
        id="f-base"
        type="number"
        min="0"
        step="any"
        bind:value={formBasePrice}
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] font-mono font-bold text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div>
      <label for="f-cito" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Override Tarif CITO (Rp) (Opsional)
      </label>
      <input
        id="f-cito"
        type="number"
        min="0"
        step="any"
        bind:value={formCitoPrice}
        placeholder="Otomatis dari % CITO jika kosong"
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] font-mono text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <!-- Pecahan Komponen Biaya Table -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738] flex flex-col gap-2">
    <div class="flex items-center justify-between">
      <div>
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs block">Pecahan Komponen Biaya Layanan</span>
        <span class="text-[11px] text-[#747775] dark:text-[#8e918f]">Distribusikan total tarif ke masing-masing komponen penagihan</span>
      </div>
      <button
        type="button"
        onclick={addComponentRow}
        class="px-3 py-1.5 rounded-xl bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa] font-bold text-xs hover:bg-[#d2e3fc] dark:hover:bg-[#0b57d0]/30 transition-colors flex items-center gap-1 cursor-pointer"
      >
        <Plus class="w-3.5 h-3.5" />
        <span>Tambah Komponen</span>
      </button>
    </div>

    <div class="border border-[#e1e5ea] dark:border-[#353738] rounded-xl overflow-hidden">
      <table class="w-full text-left text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] font-semibold border-b border-[#e1e5ea] dark:border-[#353738]">
            <th class="py-2.5 px-3">Komponen Biaya</th>
            <th class="py-2.5 px-3 w-40 text-right">Nominal Dasar (Rp)</th>
            <th class="py-2.5 px-3 w-28 text-center">Akun COA</th>
            <th class="py-2.5 px-2 w-10 text-center">Hapus</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
          {#each formComponents as comp, idx}
            <tr class="bg-white dark:bg-[#1e1f20]">
              <td class="p-2">
                <M3Select
                  bind:value={comp.component_id}
                  options={tariffComponentOptions}
                  placeholder="Pilih Komponen"
                />
              </td>
              <td class="p-2">
                <input
                  type="number"
                  min="0"
                  step="any"
                  bind:value={comp.base_amount}
                  required
                  class="w-full h-8 px-2 text-right font-mono font-bold rounded-lg bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] text-xs text-[#1f1f1f] dark:text-[#e3e3e3]"
                />
              </td>
              <td class="p-2">
                <input
                  type="text"
                  bind:value={comp.coa_code}
                  placeholder="Auto COA"
                  class="w-full h-8 px-2 font-mono text-center rounded-lg bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] text-[11px] text-[#1f1f1f] dark:text-[#e3e3e3]"
                />
              </td>
              <td class="p-2 text-center">
                <button
                  type="button"
                  onclick={() => removeComponentRow(idx)}
                  class="p-1 rounded text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/50 cursor-pointer"
                >
                  <Trash2 class="w-3.5 h-3.5" />
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Balancing Status Box -->
    <div class="p-3 rounded-xl border flex items-center justify-between text-xs transition-colors {isBalanced ? 'bg-emerald-50 dark:bg-emerald-950/30 border-emerald-200 dark:border-emerald-800 text-emerald-800 dark:text-emerald-300' : 'bg-rose-50 dark:bg-rose-950/30 border-rose-200 dark:border-rose-800 text-rose-800 dark:text-rose-300'}">
      <div class="flex items-center gap-2">
        {#if isBalanced}
          <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
          <span class="font-bold">Komponen Seimbang (100%): Rp {sumComponentBase.toLocaleString('id-ID')}</span>
        {:else}
          <AlertTriangle class="w-4 h-4 text-rose-600 shrink-0" />
          <div>
            <span class="font-bold">Belum Seimbang!</span>
            <span class="text-[11px] block">
              Total Komponen: Rp {sumComponentBase.toLocaleString('id-ID')} • Selisih: Rp {balanceDiff.toLocaleString('id-ID')}
            </span>
          </div>
        {/if}
      </div>

      {#if !isBalanced && formComponents.length > 0}
        <button
          type="button"
          onclick={handleAutoBalance}
          class="px-2.5 py-1 rounded-lg bg-white dark:bg-[#1e1f20] border border-rose-300 dark:border-rose-700 text-rose-700 dark:text-rose-300 font-bold text-[11px] hover:bg-rose-100 dark:hover:bg-rose-950/50 cursor-pointer transition-colors"
        >
          Seimbangkan Otomatis
        </button>
      {/if}
    </div>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSavingItem}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSavingItem || !isBalanced}>
      {#if isSavingItem}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>Simpan Item Tarif</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
