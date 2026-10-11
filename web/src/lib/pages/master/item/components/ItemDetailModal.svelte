<script lang="ts">
  import {
    X,
    Plus,
    Trash2,
    RefreshCw,
    Edit3,
    Layers,
    Pill,
    Package,
    Stethoscope,
    Receipt
  } from '@lucide/svelte';
  import M3Button from '../../../../components/m3/M3Button.svelte';
  import ItemSpecDetailTab from './ItemSpecDetailTab.svelte';
  import { confirmDialog } from '../../../../stores/confirm.svelte';
  import {
    getItemById,
    addItemUnit,
    deleteItemUnit
  } from '../../../../api/catalog';
  import type {
    ItemSummaryRecord,
    ItemDetailRecord,
    ItemUnitRecord,
    ItemType
  } from '../../../../types/master/item';

  interface Props {
    isOpen: boolean;
    itemSummary: ItemSummaryRecord | null;
    onClose: () => void;
    onEdit: (item: ItemSummaryRecord) => void;
    showToast: (msg: string, isError?: boolean) => void;
  }

  let {
    isOpen,
    itemSummary,
    onClose,
    onEdit,
    showToast
  }: Props = $props();

  let detailItem = $state<ItemDetailRecord | null>(null);
  let detailUnits = $state<ItemUnitRecord[]>([]);
  let detailActiveTab = $state<'SPEC' | 'UOM' | 'COA'>('SPEC');
  let isDetailLoading = $state(false);
  let isActionLoading = $state(false);

  // UOM Management
  let showAddUnitForm = $state(false);
  let newUnitName = $state('');
  let newUnitFactor = $state<number>(1);
  let newUnitIsPurchase = $state(false);
  let newUnitIsDispense = $state(false);

  $effect(() => {
    if (isOpen && itemSummary) {
      loadDetail(itemSummary.id);
    } else {
      detailItem = null;
      detailUnits = [];
      showAddUnitForm = false;
    }
  });

  async function loadDetail(id: string) {
    isDetailLoading = true;
    detailActiveTab = 'SPEC';
    try {
      const res = await getItemById(id);
      detailItem = res;
      detailUnits = res.units || [];
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat rincian item', true);
      onClose();
    } finally {
      isDetailLoading = false;
    }
  }

  async function handleAddUnit() {
    if (!detailItem || !newUnitName.trim() || newUnitFactor <= 0) {
      showToast('Nama satuan dan faktor konversi (> 0) wajib diisi', true);
      return;
    }

    isActionLoading = true;
    try {
      const created = await addItemUnit(detailItem.id, {
        unit_name: newUnitName.trim(),
        conversion_factor: Number(newUnitFactor),
        is_purchase_unit: newUnitIsPurchase,
        is_dispense_unit: newUnitIsDispense
      });
      detailUnits = [...detailUnits, created];
      newUnitName = '';
      newUnitFactor = 1;
      newUnitIsPurchase = false;
      newUnitIsDispense = false;
      showAddUnitForm = false;
      showToast('Satuan alternatif baru berhasil ditambahkan');
    } catch (err: any) {
      showToast(err.message || 'Gagal menambahkan satuan alternatif', true);
    } finally {
      isActionLoading = false;
    }
  }

  async function handleDeleteUnit(unit: ItemUnitRecord) {
    if (!detailItem) return;
    if (unit.is_base_unit) {
      showToast('Satuan dasar (Base Unit) tidak dapat dihapus', true);
      return;
    }

    const ok = await confirmDialog({
      title: 'Hapus Satuan Alternatif',
      message: `Apakah Anda yakin ingin menghapus konversi satuan "${unit.unit_name}" (isi: ${unit.conversion_factor})? Tindakan ini tidak dapat dibatalkan.`,
      variant: 'danger',
      confirmText: 'Hapus',
      cancelText: 'Batal'
    });
    if (!ok) return;

    isActionLoading = true;
    try {
      await deleteItemUnit(detailItem.id, unit.id);
      detailUnits = detailUnits.filter(u => u.id !== unit.id);
      showToast(`Satuan ${unit.unit_name} berhasil dihapus`);
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus satuan', true);
    } finally {
      isActionLoading = false;
    }
  }

  function getTypeBadge(type: ItemType) {
    switch (type) {
      case 'MEDICATION':
        return { label: 'Obat', bg: 'bg-emerald-50 text-emerald-700 border-emerald-200', icon: Pill };
      case 'GENERAL':
        return { label: 'BMHP / Umum', bg: 'bg-amber-50 text-amber-700 border-amber-200', icon: Package };
      case 'ASSET':
        return { label: 'Aset / Alkes', bg: 'bg-indigo-50 text-indigo-700 border-indigo-200', icon: Stethoscope };
      case 'TARIFF':
        return { label: 'Tarif Jasa', bg: 'bg-blue-50 text-blue-700 border-blue-200', icon: Receipt };
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col my-auto max-h-[90vh]">
      <!-- Modal Header -->
      <div class="p-5 border-b border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between">
        <div class="flex items-center gap-3">
          {#if detailItem}
            {@const b = getTypeBadge(detailItem.item_type)}
            <div class="w-10 h-10 rounded-2xl {b.bg} flex items-center justify-center font-bold">
              <b.icon class="w-5 h-5" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-xs text-[#0b57d0]">{detailItem.code}</span>
                <span class="text-xs px-2 py-0.5 rounded-full font-semibold {b.bg}">{b.label}</span>
              </div>
              <h3 class="text-sm font-bold text-[#1f1f1f] mt-0.5">{detailItem.name}</h3>
            </div>
          {:else}
            <div class="text-sm font-bold text-[#1f1f1f]">Memuat Rincian Item...</div>
          {/if}
        </div>
        <button
          type="button"
          onclick={onClose}
          class="w-8 h-8 rounded-full text-[#444746] hover:bg-[#e1e5ea] flex items-center justify-center transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 overflow-y-auto flex-1 flex flex-col gap-5">
        {#if isDetailLoading || !detailItem}
          <div class="p-12 flex flex-col items-center justify-center gap-3">
            <RefreshCw class="w-8 h-8 text-[#0b57d0] animate-spin" />
            <p class="text-xs text-[#747775]">Memuat rincian spesifikasi item...</p>
          </div>
        {:else}
          <!-- Sub-Tabs in Detail -->
          <div class="flex items-center gap-1 border-b border-[#e1e5ea] pb-2">
            <button
              type="button"
              onclick={() => detailActiveTab = 'SPEC'}
              class="px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-colors cursor-pointer {detailActiveTab === 'SPEC' ? 'bg-[#e8f0fe] text-[#0b57d0]' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
            >
              Spesifikasi Subtipe ({detailItem.item_type})
            </button>
            <button
              type="button"
              onclick={() => detailActiveTab = 'UOM'}
              class="px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-colors cursor-pointer {detailActiveTab === 'UOM' ? 'bg-[#e8f0fe] text-[#0b57d0]' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
            >
              Satuan Multi-UOM ({detailUnits.length})
            </button>
            <button
              type="button"
              onclick={() => detailActiveTab = 'COA'}
              class="px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-colors cursor-pointer {detailActiveTab === 'COA' ? 'bg-[#e8f0fe] text-[#0b57d0]' : 'text-[#444746] hover:bg-[#f0f4f9]'}"
            >
              Pemetaan Akuntansi & COA
            </button>
          </div>

          <!-- DETAIL SUB-TAB 1: SPESIFIKASI KHUSUS SUBTIPE -->
          {#if detailActiveTab === 'SPEC'}
            <ItemSpecDetailTab {detailItem} />

          <!-- DETAIL SUB-TAB 2: SATUAN MULTI-UOM -->
          {:else if detailActiveTab === 'UOM'}
            <div class="flex flex-col gap-4">
              <div class="flex items-center justify-between">
                <div>
                  <h4 class="text-xs font-bold text-[#1f1f1f]">Daftar Konversi Satuan (Multi-UOM)</h4>
                  <p class="text-[11px] text-[#747775]">Satuan dasar menjadi acuan inventori (faktor = 1.0).</p>
                </div>
                {#if !showAddUnitForm}
                  <M3Button variant="tonal" onclick={() => showAddUnitForm = true}>
                    <Plus class="w-3.5 h-3.5" />
                    <span>Tambah Satuan</span>
                  </M3Button>
                {/if}
              </div>

              <!-- Inline Add Unit Form -->
              {#if showAddUnitForm}
                <div class="p-4 rounded-2xl bg-[#f0f4f9] border border-[#e1e5ea] flex flex-col gap-3">
                  <div class="flex items-center justify-between">
                    <span class="text-xs font-bold text-[#1f1f1f]">Tambah Satuan Alternatif Baru</span>
                    <button type="button" onclick={() => showAddUnitForm = false} class="text-[#747775] hover:text-[#1f1f1f]">
                      <X class="w-4 h-4" />
                    </button>
                  </div>

                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <label for="new-uom-name" class="block text-[11px] font-semibold text-[#444746] mb-1">Nama Satuan (Box, Strip, dll)</label>
                      <input
                        id="new-uom-name"
                        type="text"
                        bind:value={newUnitName}
                        placeholder="Misal: Box"
                        class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs"
                      />
                    </div>
                    <div>
                      <label for="new-uom-factor" class="block text-[11px] font-semibold text-[#444746] mb-1">Faktor Konversi (x Base Unit)</label>
                      <input
                        id="new-uom-factor"
                        type="number"
                        step="any"
                        bind:value={newUnitFactor}
                        class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs font-mono"
                      />
                    </div>
                  </div>

                  <div class="flex items-center gap-4 text-xs">
                    <label class="flex items-center gap-2 cursor-pointer">
                      <input type="checkbox" bind:checked={newUnitIsPurchase} class="rounded text-[#0b57d0]" />
                      <span>Satuan Pembelian (PO)</span>
                    </label>
                    <label class="flex items-center gap-2 cursor-pointer">
                      <input type="checkbox" bind:checked={newUnitIsDispense} class="rounded text-[#0b57d0]" />
                      <span>Satuan Resep / Dispensing</span>
                    </label>
                  </div>

                  <div class="flex justify-end gap-2 mt-1">
                    <button
                      type="button"
                      onclick={() => showAddUnitForm = false}
                      class="px-3 py-1.5 rounded-lg border border-[#e1e5ea] bg-white text-xs text-[#444746]"
                    >
                      Batal
                    </button>
                    <M3Button variant="filled" onclick={handleAddUnit} disabled={isActionLoading}>
                      <span>Simpan Satuan</span>
                    </M3Button>
                  </div>
                </div>
              {/if}

              <!-- Units Table -->
              <div class="border border-[#e1e5ea] rounded-2xl overflow-hidden">
                <table class="w-full text-left text-xs border-collapse">
                  <thead>
                    <tr class="bg-[#f0f4f9] text-[#444746] font-semibold border-b border-[#e1e5ea]">
                      <th class="py-2.5 px-3">Nama Satuan</th>
                      <th class="py-2.5 px-3 font-mono">Faktor Konversi</th>
                      <th class="py-2.5 px-3">Peran Satuan</th>
                      <th class="py-2.5 px-3 text-right">Aksi</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-[#e1e5ea]">
                    {#each detailUnits as unit (unit.id)}
                      <tr class="hover:bg-[#f8fafd]">
                        <td class="py-3 px-3 font-semibold text-[#1f1f1f]">
                          <div class="flex items-center gap-2">
                            <span>{unit.unit_name}</span>
                            {#if unit.is_base_unit}
                              <span class="px-1.5 py-0.5 rounded text-[10px] bg-[#d3e3fd] text-[#041e49] font-bold">
                                BASE UNIT
                              </span>
                            {/if}
                          </div>
                        </td>
                        <td class="py-3 px-3 font-mono font-bold text-[#0b57d0]">
                          1 {unit.unit_name} = {unit.conversion_factor} {detailUnits.find(u => u.is_base_unit)?.unit_name || 'Unit'}
                        </td>
                        <td class="py-3 px-3">
                          <div class="flex items-center gap-1.5 flex-wrap">
                            {#if unit.is_purchase_unit}
                              <span class="px-2 py-0.5 rounded text-[10px] bg-amber-50 text-amber-700 border border-amber-200">
                                Beli (PO)
                              </span>
                            {/if}
                            {#if unit.is_dispense_unit}
                              <span class="px-2 py-0.5 rounded text-[10px] bg-emerald-50 text-emerald-700 border border-emerald-200">
                                Resep / Jual
                              </span>
                            {/if}
                            {#if !unit.is_purchase_unit && !unit.is_dispense_unit && !unit.is_base_unit}
                              <span class="text-[11px] text-[#747775]">Alternatif</span>
                            {/if}
                          </div>
                        </td>
                        <td class="py-3 px-3 text-right">
                          {#if !unit.is_base_unit}
                            <button
                              type="button"
                              onclick={() => handleDeleteUnit(unit)}
                              class="p-1 text-[#c5221f] hover:bg-[#c5221f]/10 rounded cursor-pointer"
                              title="Hapus Satuan Alternatif"
                            >
                              <Trash2 class="w-3.5 h-3.5" />
                            </button>
                          {:else}
                            <span class="text-[10px] text-[#747775] italic">Terkunci</span>
                          {/if}
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>
            </div>

          <!-- DETAIL SUB-TAB 3: INTEGRASI AKUNTANSI & COA -->
          {:else if detailActiveTab === 'COA'}
            <div class="flex flex-col gap-4 text-xs">
              <!-- Revenue COA -->
              <div class="p-4 rounded-2xl bg-blue-50/50 border border-blue-200 flex flex-col gap-2">
                <div class="flex items-center gap-2 text-blue-900 font-bold text-xs">
                  <Receipt class="w-4 h-4 text-blue-700" />
                  <span>Akun Pendapatan & Penagihan (Billing Module)</span>
                </div>
                <p class="text-xs text-[#444746]">
                  Diambil secara otomatis dari Kategori: <strong>{detailItem.category?.name || '-'}</strong>
                </p>
                <div class="grid grid-cols-2 gap-3 mt-1 bg-white p-3 rounded-xl border border-blue-100">
                  <div>
                    <span class="text-[11px] text-[#747775] block">COA Akun Pendapatan</span>
                    <span class="font-mono font-bold text-[#0b57d0]">{detailItem.category?.income_coa_code || '-'}</span>
                  </div>
                  <div>
                    <span class="text-[11px] text-[#747775] block">Tipe Jurnal Transaksi</span>
                    <span class="font-semibold text-[#1f1f1f]">Kredit Pendapatan Kasir</span>
                  </div>
                </div>
              </div>

              <!-- Inventory & COGS COA -->
              {#if detailItem.product_line}
                <div class="p-4 rounded-2xl bg-amber-50/50 border border-amber-200 flex flex-col gap-2">
                  <div class="flex items-center gap-2 text-amber-900 font-bold text-xs">
                    <Layers class="w-4 h-4 text-amber-700" />
                    <span>Akun Persediaan & Harga Pokok Penjualan (HPP)</span>
                  </div>
                  <p class="text-xs text-[#444746]">
                    Diambil secara otomatis dari Lini Produk: <strong>{detailItem.product_line.name}</strong>
                  </p>
                  <div class="grid grid-cols-2 gap-3 mt-1 bg-white p-3 rounded-xl border border-amber-100">
                    <div>
                      <span class="text-[11px] text-[#747775] block">COA Akun Persediaan (Aset)</span>
                      <span class="font-mono font-bold text-amber-800">{detailItem.product_line.inventory_coa_code || '-'}</span>
                    </div>
                    <div>
                      <span class="text-[11px] text-[#747775] block">COA Akun Beban Pokok (HPP)</span>
                      <span class="font-mono font-bold text-amber-800">{detailItem.product_line.cogs_coa_code || '-'}</span>
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          {/if}
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="p-4 border-t border-[#e1e5ea] bg-[#f8fafd] flex justify-end gap-2">
        {#if detailItem}
          <button
            type="button"
            onclick={() => { onClose(); onEdit(detailItem!); }}
            class="px-4 py-2 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#1f1f1f] hover:bg-[#f0f4f9] flex items-center gap-1.5 cursor-pointer shadow-xs"
          >
            <Edit3 class="w-3.5 h-3.5" />
            <span>Edit Item</span>
          </button>
        {/if}
        <M3Button variant="tonal" onclick={onClose}>
          <span>Tutup</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}
