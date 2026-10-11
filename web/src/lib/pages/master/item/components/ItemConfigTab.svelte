<script lang="ts">
  import { Tag, Layers, Plus, Trash2, Edit3 } from '@lucide/svelte';
  import M3Button from '../../../../components/m3/M3Button.svelte';
  import ItemCategoryModal from './ItemCategoryModal.svelte';
  import ProductLineModal from './ProductLineModal.svelte';
  import { confirmDialog } from '../../../../stores/confirm.svelte';
  import {
    createCategory,
    updateCategory,
    deleteCategory,
    createProductLine,
    updateProductLine,
    deleteProductLine
  } from '../../../../api/catalog';
  import type {
    ItemCategoryRecord,
    ItemProductLineRecord,
    ItemType
  } from '../../../../types/master/item';

  interface Props {
    categories: ItemCategoryRecord[];
    productLines: ItemProductLineRecord[];
    onRefresh: () => void;
    showToast: (msg: string, isError?: boolean) => void;
  }

  let {
    categories,
    productLines,
    onRefresh,
    showToast
  }: Props = $props();

  let isActionLoading = $state(false);

  // Category Modal State
  let showCategoryModal = $state(false);
  let catModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let editingCatId = $state<string | null>(null);
  let catFormCode = $state('');
  let catFormName = $state('');
  let catFormType = $state<ItemType>('MEDICATION');
  let catFormIncomeCoa = $state('');
  let catFormIsActive = $state(true);

  // Product Line Modal State
  let showProductLineModal = $state(false);
  let plModalMode = $state<'CREATE' | 'EDIT'>('CREATE');
  let editingPlId = $state<string | null>(null);
  let plFormCode = $state('');
  let plFormName = $state('');
  let plFormInventoryCoa = $state('');
  let plFormCogsCoa = $state('');
  let plFormIsActive = $state(true);

  // ==========================================
  // CATEGORY ACTIONS
  // ==========================================

  function openCreateCategory() {
    catModalMode = 'CREATE';
    editingCatId = null;
    catFormCode = '';
    catFormName = '';
    catFormType = 'MEDICATION';
    catFormIncomeCoa = '';
    catFormIsActive = true;
    showCategoryModal = true;
  }

  function openEditCategory(cat: ItemCategoryRecord) {
    catModalMode = 'EDIT';
    editingCatId = cat.id;
    catFormCode = cat.code;
    catFormName = cat.name;
    catFormType = cat.item_type;
    catFormIncomeCoa = cat.income_coa_code || '';
    catFormIsActive = cat.is_active;
    showCategoryModal = true;
  }

  async function handleSaveCategory() {
    if (!catFormCode.trim() || !catFormName.trim()) {
      showToast('Kode dan nama kategori wajib diisi', true);
      return;
    }

    isActionLoading = true;
    try {
      if (catModalMode === 'CREATE') {
        await createCategory({
          code: catFormCode.trim(),
          name: catFormName.trim(),
          item_type: catFormType,
          income_coa_code: catFormIncomeCoa.trim() || undefined,
          is_active: true
        });
        showToast(`Kategori ${catFormName} berhasil ditambahkan`);
      } else if (editingCatId) {
        await updateCategory(editingCatId, {
          code: catFormCode.trim(),
          name: catFormName.trim(),
          item_type: catFormType,
          income_coa_code: catFormIncomeCoa.trim() || undefined,
          is_active: catFormIsActive
        });
        showToast(`Kategori ${catFormName} berhasil diperbarui`);
      }

      showCategoryModal = false;
      onRefresh();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan kategori', true);
    } finally {
      isActionLoading = false;
    }
  }

  async function handleDeleteCategory(id: string, name: string) {
    const ok = await confirmDialog({
      title: 'Hapus Kategori',
      message: `Apakah Anda yakin ingin menghapus kategori "${name}"? Data yang sudah dihapus tidak dapat dikembalikan.`,
      variant: 'danger',
      confirmText: 'Hapus',
      cancelText: 'Batal'
    });
    if (!ok) return;

    try {
      await deleteCategory(id);
      showToast(`Kategori ${name} berhasil dihapus`);
      onRefresh();
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus kategori', true);
    }
  }

  // ==========================================
  // PRODUCT LINE ACTIONS
  // ==========================================

  function openCreateProductLine() {
    plModalMode = 'CREATE';
    editingPlId = null;
    plFormCode = '';
    plFormName = '';
    plFormInventoryCoa = '';
    plFormCogsCoa = '';
    plFormIsActive = true;
    showProductLineModal = true;
  }

  function openEditProductLine(pl: ItemProductLineRecord) {
    plModalMode = 'EDIT';
    editingPlId = pl.id;
    plFormCode = pl.code;
    plFormName = pl.name;
    plFormInventoryCoa = pl.inventory_coa_code || '';
    plFormCogsCoa = pl.cogs_coa_code || '';
    plFormIsActive = pl.is_active;
    showProductLineModal = true;
  }

  async function handleSaveProductLine() {
    if (!plFormCode.trim() || !plFormName.trim()) {
      showToast('Kode dan nama lini produk wajib diisi', true);
      return;
    }

    isActionLoading = true;
    try {
      if (plModalMode === 'CREATE') {
        await createProductLine({
          code: plFormCode.trim(),
          name: plFormName.trim(),
          inventory_coa_code: plFormInventoryCoa.trim() || undefined,
          cogs_coa_code: plFormCogsCoa.trim() || undefined,
          is_active: true
        });
        showToast(`Lini produk ${plFormName} berhasil ditambahkan`);
      } else if (editingPlId) {
        await updateProductLine(editingPlId, {
          code: plFormCode.trim(),
          name: plFormName.trim(),
          inventory_coa_code: plFormInventoryCoa.trim() || undefined,
          cogs_coa_code: plFormCogsCoa.trim() || undefined,
          is_active: plFormIsActive
        });
        showToast(`Lini produk ${plFormName} berhasil diperbarui`);
      }

      showProductLineModal = false;
      onRefresh();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan lini produk', true);
    } finally {
      isActionLoading = false;
    }
  }

  async function handleDeleteProductLine(id: string, name: string) {
    const ok = await confirmDialog({
      title: 'Hapus Lini Produk',
      message: `Apakah Anda yakin ingin menghapus lini produk "${name}"? Data yang sudah dihapus tidak dapat dikembalikan.`,
      variant: 'danger',
      confirmText: 'Hapus',
      cancelText: 'Batal'
    });
    if (!ok) return;

    try {
      await deleteProductLine(id);
      showToast(`Lini produk ${name} berhasil dihapus`);
      onRefresh();
    } catch (err: any) {
      showToast(err.message || 'Gagal menghapus lini produk', true);
    }
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
  <!-- Panel 1: Master Kategori Item -->
  <div class="bg-white rounded-2xl border border-[#e1e5ea] p-5 flex flex-col gap-4 shadow-xs">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-sm font-bold text-[#1f1f1f] flex items-center gap-2">
          <Tag class="w-4 h-4 text-[#0b57d0]" />
          <span>Master Kategori Item ({categories.length})</span>
        </h3>
        <p class="text-xs text-[#747775]">Memetakan kelompok barang ke Akun Pendapatan (COA Revenue).</p>
      </div>
      <M3Button variant="tonal" onclick={openCreateCategory}>
        <Plus class="w-3.5 h-3.5" />
        <span>Tambah Kategori</span>
      </M3Button>
    </div>

    <div class="divide-y divide-[#e1e5ea] border border-[#e1e5ea] rounded-xl overflow-hidden max-h-[460px] overflow-y-auto">
      {#each categories as cat (cat.id)}
        <div class="p-3 flex items-center justify-between hover:bg-[#f8fafd] transition-colors">
          <div>
            <div class="flex items-center gap-2">
              <span class="font-mono text-xs font-bold text-[#0b57d0]">{cat.code}</span>
              <span class="text-xs font-semibold text-[#1f1f1f]">{cat.name}</span>
              {#if !cat.is_active}
                <span class="px-1.5 py-0.2 rounded text-[10px] bg-rose-50 text-rose-700 border border-rose-200">Nonaktif</span>
              {/if}
            </div>
            <div class="flex items-center gap-2 mt-1 text-[11px] text-[#747775]">
              <span class="px-1.5 py-0.5 rounded bg-[#f0f4f9] font-mono">{cat.item_type}</span>
              {#if cat.income_coa_code}
                <span>COA Pendapatan: <strong class="font-mono">{cat.income_coa_code}</strong></span>
              {/if}
            </div>
          </div>
          <div class="flex items-center gap-1">
            <button
              type="button"
              onclick={() => openEditCategory(cat)}
              class="p-1.5 text-[#444746] hover:bg-[#f0f4f9] rounded-lg transition-colors cursor-pointer"
              title="Edit Kategori"
            >
              <Edit3 class="w-3.5 h-3.5" />
            </button>
            <button
              type="button"
              onclick={() => handleDeleteCategory(cat.id, cat.name)}
              class="p-1.5 text-[#c5221f] hover:bg-[#c5221f]/10 rounded-lg transition-colors cursor-pointer"
              title="Hapus Kategori"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <!-- Panel 2: Master Lini Produk Persediaan -->
  <div class="bg-white rounded-2xl border border-[#e1e5ea] p-5 flex flex-col gap-4 shadow-xs">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-sm font-bold text-[#1f1f1f] flex items-center gap-2">
          <Layers class="w-4 h-4 text-amber-600" />
          <span>Lini Produk Persediaan ({productLines.length})</span>
        </h3>
        <p class="text-xs text-[#747775]">Memetakan kelompok fisik barang ke Akun Persediaan & HPP/Beban.</p>
      </div>
      <M3Button variant="tonal" onclick={openCreateProductLine}>
        <Plus class="w-3.5 h-3.5" />
        <span>Tambah Lini Produk</span>
      </M3Button>
    </div>

    <div class="divide-y divide-[#e1e5ea] border border-[#e1e5ea] rounded-xl overflow-hidden max-h-[460px] overflow-y-auto">
      {#each productLines as pl (pl.id)}
        <div class="p-3 flex items-center justify-between hover:bg-[#f8fafd] transition-colors">
          <div>
            <div class="flex items-center gap-2">
              <span class="font-mono text-xs font-bold text-amber-700">{pl.code}</span>
              <span class="text-xs font-semibold text-[#1f1f1f]">{pl.name}</span>
              {#if !pl.is_active}
                <span class="px-1.5 py-0.2 rounded text-[10px] bg-rose-50 text-rose-700 border border-rose-200">Nonaktif</span>
              {/if}
            </div>
            <div class="flex items-center gap-2 mt-1 text-[11px] text-[#747775] flex-wrap">
              {#if pl.inventory_coa_code}
                <span>Persediaan: <strong class="font-mono">{pl.inventory_coa_code}</strong></span>
              {/if}
              {#if pl.cogs_coa_code}
                <span>• HPP: <strong class="font-mono">{pl.cogs_coa_code}</strong></span>
              {/if}
            </div>
          </div>
          <div class="flex items-center gap-1">
            <button
              type="button"
              onclick={() => openEditProductLine(pl)}
              class="p-1.5 text-[#444746] hover:bg-[#f0f4f9] rounded-lg transition-colors cursor-pointer"
              title="Edit Lini Produk"
            >
              <Edit3 class="w-3.5 h-3.5" />
            </button>
            <button
              type="button"
              onclick={() => handleDeleteProductLine(pl.id, pl.name)}
              class="p-1.5 text-[#c5221f] hover:bg-[#c5221f]/10 rounded-lg transition-colors cursor-pointer"
              title="Hapus Lini Produk"
            >
              <Trash2 class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>

<!-- Modal Tambah / Edit Kategori -->
<ItemCategoryModal
  isOpen={showCategoryModal}
  mode={catModalMode}
  bind:code={catFormCode}
  bind:name={catFormName}
  bind:type={catFormType}
  bind:incomeCoa={catFormIncomeCoa}
  bind:isActive={catFormIsActive}
  isLoading={isActionLoading}
  onClose={() => (showCategoryModal = false)}
  onSubmit={handleSaveCategory}
/>

<!-- Modal Tambah / Edit Lini Produk -->
<ProductLineModal
  isOpen={showProductLineModal}
  mode={plModalMode}
  bind:code={plFormCode}
  bind:name={plFormName}
  bind:inventoryCoa={plFormInventoryCoa}
  bind:cogsCoa={plFormCogsCoa}
  bind:isActive={plFormIsActive}
  isLoading={isActionLoading}
  onClose={() => (showProductLineModal = false)}
  onSubmit={handleSaveProductLine}
/>
