<script lang="ts">
  import {
    Plus,
    Trash2,
    RefreshCw,
    Edit3,
    Layers,
    Pill,
    Package,
    Stethoscope,
    Receipt,
    X,
  } from "@lucide/svelte";
  import M3Modal from "../../../../components/m3/M3Modal.svelte";
  import M3Button from "../../../../components/m3/M3Button.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import ItemSpecDetailTab from "./ItemSpecDetailTab.svelte";
  import { confirmDialog } from "../../../../stores/confirm.svelte";
  import {
    getItemById,
    addItemUnit,
    deleteItemUnit,
  } from "../../../../api/catalog";
  import type {
    ItemSummaryRecord,
    ItemDetailRecord,
    ItemUnitRecord,
    ItemType,
  } from "../../../../types/master/item";

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
    showToast,
  }: Props = $props();

  let detailItem = $state<ItemDetailRecord | null>(null);
  let detailUnits = $state<ItemUnitRecord[]>([]);
  let detailActiveTab = $state<"SPEC" | "UOM" | "COA">("SPEC");
  let isDetailLoading = $state(false);
  let isActionLoading = $state(false);

  // UOM Management
  let showAddUnitForm = $state(false);
  let newUnitName = $state("");
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
    detailActiveTab = "SPEC";
    try {
      const res = await getItemById(id);
      detailItem = res;
      detailUnits = res.units || [];
    } catch (err: any) {
      showToast(err.message || "Gagal memuat rincian item", true);
    } finally {
      isDetailLoading = false;
    }
  }

  async function handleAddUnit() {
    if (!newUnitName.trim() || newUnitFactor <= 0) {
      showToast("Nama satuan dan faktor konversi (> 0) wajib diisi", true);
      return;
    }
    if (!detailItem) return;

    isActionLoading = true;
    try {
      const res = await addItemUnit(detailItem.id, {
        unit_name: newUnitName.trim(),
        conversion_factor: Number(newUnitFactor),
        is_purchase_unit: newUnitIsPurchase,
        is_dispense_unit: newUnitIsDispense,
      });

      detailUnits = [...detailUnits, res];
      showToast(`Satuan ${newUnitName} berhasil ditambahkan`);
      showAddUnitForm = false;
      newUnitName = "";
      newUnitFactor = 1;
      newUnitIsPurchase = false;
      newUnitIsDispense = false;
    } catch (err: any) {
      showToast(err.message || "Gagal menambahkan satuan konversi", true);
    } finally {
      isActionLoading = false;
    }
  }

  async function handleDeleteUnit(unit: ItemUnitRecord) {
    if (!detailItem) return;

    const ok = await confirmDialog({
      title: "Hapus Satuan Alternatif",
      message: `Apakah Anda yakin ingin menghapus konversi satuan "${unit.unit_name}" (isi: ${unit.conversion_factor})? Tindakan ini tidak dapat dibatalkan.`,
      variant: "danger",
      confirmText: "Hapus",
      cancelText: "Batal",
    });
    if (!ok) return;

    isActionLoading = true;
    try {
      await deleteItemUnit(detailItem.id, unit.id);
      detailUnits = detailUnits.filter((u) => u.id !== unit.id);
      showToast(`Satuan ${unit.unit_name} berhasil dihapus`);
    } catch (err: any) {
      showToast(err.message || "Gagal menghapus satuan", true);
    } finally {
      isActionLoading = false;
    }
  }

  function getTypeBadge(type: ItemType) {
    switch (type) {
      case "MEDICATION":
        return {
          label: "Obat",
          bg: "bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border-emerald-200 dark:border-emerald-800",
          icon: Pill,
        };
      case "GENERAL":
        return {
          label: "BMHP / Umum",
          bg: "bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-400 border-amber-200 dark:border-amber-800",
          icon: Package,
        };
      case "ASSET":
        return {
          label: "Aset / Alkes",
          bg: "bg-indigo-50 dark:bg-indigo-950/40 text-indigo-700 dark:text-indigo-400 border-indigo-200 dark:border-indigo-800",
          icon: Stethoscope,
        };
      case "TARIFF":
        return {
          label: "Tarif Jasa",
          bg: "bg-blue-50 dark:bg-blue-950/40 text-blue-700 dark:text-blue-400 border-blue-200 dark:border-blue-800",
          icon: Receipt,
        };
    }
  }
</script>

<M3Modal
  {isOpen}
  title={detailItem ? `${detailItem.name} (${detailItem.code})` : "Rincian Item"}
  size="xl"
  {onClose}
>
  <div class="flex flex-col gap-5 text-xs">
    {#if isDetailLoading || !detailItem}
      <div class="p-12 flex flex-col items-center justify-center gap-3">
        <RefreshCw class="w-8 h-8 text-[#0b57d0] animate-spin" />
        <p class="text-xs text-[#747775] dark:text-[#c4c7c5]">
          Memuat rincian spesifikasi item...
        </p>
      </div>
    {:else}
      <!-- Header Badge summary -->
      {@const b = getTypeBadge(detailItem.item_type)}
      <div
        class="p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#444746] flex items-center justify-between"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-2xl {b.bg} flex items-center justify-center font-bold"
          >
            <b.icon class="w-5 h-5" />
          </div>
          <div>
            <div class="flex items-center gap-2">
              <span class="font-mono font-bold text-xs text-[#0b57d0] dark:text-blue-400">
                {detailItem.code}
              </span>
              <span class="text-xs px-2 py-0.5 rounded-full font-semibold {b.bg}">
                {b.label}
              </span>
            </div>
            <h3 class="text-sm font-bold text-[#1f1f1f] dark:text-[#e3e3e3] mt-0.5">
              {detailItem.name}
            </h3>
          </div>
        </div>
      </div>

      <!-- Sub-Tabs in Detail -->
      <div class="flex items-center gap-1 border-b border-[#e1e5ea] dark:border-[#444746] pb-2">
        <button
          type="button"
          onclick={() => (detailActiveTab = "SPEC")}
          class="px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-colors cursor-pointer {detailActiveTab ===
          'SPEC'
            ? 'bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-blue-400'
            : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
        >
          Spesifikasi Subtipe ({detailItem.item_type})
        </button>
        <button
          type="button"
          onclick={() => (detailActiveTab = "UOM")}
          class="px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-colors cursor-pointer {detailActiveTab ===
          'UOM'
            ? 'bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-blue-400'
            : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
        >
          Satuan Multi-UOM ({detailUnits.length})
        </button>
        <button
          type="button"
          onclick={() => (detailActiveTab = "COA")}
          class="px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-colors cursor-pointer {detailActiveTab ===
          'COA'
            ? 'bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-blue-400'
            : 'text-[#444746] dark:text-[#c4c7c5] hover:bg-[#f0f4f9] dark:hover:bg-[#282a2c]'}"
        >
          Pemetaan Akuntansi & COA
        </button>
      </div>

      <!-- DETAIL SUB-TAB 1: SPESIFIKASI KHUSUS SUBTIPE -->
      {#if detailActiveTab === "SPEC"}
        <ItemSpecDetailTab {detailItem} />

      <!-- DETAIL SUB-TAB 2: SATUAN MULTI-UOM -->
      {:else if detailActiveTab === "UOM"}
        <div class="flex flex-col gap-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-xs font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
                Daftar Konversi Satuan (Multi-UOM)
              </h4>
              <p class="text-[11px] text-[#747775] dark:text-[#c4c7c5]">
                Satuan dasar menjadi acuan inventori (faktor = 1.0).
              </p>
            </div>
            {#if !showAddUnitForm}
              <M3Button variant="tonal" onclick={() => (showAddUnitForm = true)}>
                <Plus class="w-3.5 h-3.5" />
                <span>Tambah Satuan</span>
              </M3Button>
            {/if}
          </div>

          <!-- Inline Add Unit Form -->
          {#if showAddUnitForm}
            <div
              class="p-4 rounded-2xl bg-[#f0f4f9] dark:bg-[#282a2c] border border-[#e1e5ea] dark:border-[#444746] flex flex-col gap-3"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
                  Tambah Satuan Alternatif Baru
                </span>
                <button
                  type="button"
                  onclick={() => (showAddUnitForm = false)}
                  class="text-[#747775] dark:text-[#c4c7c5] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3]"
                >
                  <X class="w-4 h-4" />
                </button>
              </div>

              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label
                    for="new-uom-name"
                    class="block text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
                  >
                    Nama Satuan (Box, Strip, dll)
                  </label>
                  <input
                    id="new-uom-name"
                    type="text"
                    bind:value={newUnitName}
                    placeholder="Misal: Box"
                    class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
                  />
                </div>
                <div>
                  <label
                    for="new-uom-factor"
                    class="block text-[11px] font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
                  >
                    Faktor Konversi (x Base Unit)
                  </label>
                  <input
                    id="new-uom-factor"
                    type="number"
                    step="any"
                    bind:value={newUnitFactor}
                    class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs font-mono focus:border-[#0b57d0] focus:outline-none"
                  />
                </div>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
                <M3ToggleCard
                  bind:checked={newUnitIsPurchase}
                  title="Satuan Pembelian (PO)"
                  subtitle="Digunakan pada pemesanan ke pemasok/vendor"
                />
                <M3ToggleCard
                  bind:checked={newUnitIsDispense}
                  title="Satuan Resep / Dispensing"
                  subtitle="Digunakan pada pelayanan resep dan farmasi"
                />
              </div>

              <div class="flex justify-end gap-2 mt-1">
                <M3Button variant="text" onclick={() => (showAddUnitForm = false)}>
                  Batal
                </M3Button>
                <M3Button
                  variant="filled"
                  onclick={handleAddUnit}
                  disabled={isActionLoading}
                >
                  <span>Simpan Satuan</span>
                </M3Button>
              </div>
            </div>
          {/if}

          <!-- Units Table -->
          <div
            class="border border-[#e1e5ea] dark:border-[#444746] rounded-2xl overflow-hidden"
          >
            <table class="w-full text-left text-xs border-collapse">
              <thead>
                <tr
                  class="bg-[#f0f4f9] dark:bg-[#2b2b2b] text-[#444746] dark:text-[#c4c7c5] font-semibold border-b border-[#e1e5ea] dark:border-[#444746]"
                >
                  <th class="py-2.5 px-3">Nama Satuan</th>
                  <th class="py-2.5 px-3 font-mono">Faktor Konversi</th>
                  <th class="py-2.5 px-3">Peran Satuan</th>
                  <th class="py-2.5 px-3 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3]">
                {#each detailUnits as unit (unit.id)}
                  <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#282a2c]">
                    <td class="py-3 px-3 font-semibold">
                      <div class="flex items-center gap-2">
                        <span>{unit.unit_name}</span>
                        {#if unit.is_base_unit}
                          <span
                            class="px-1.5 py-0.5 rounded text-[10px] bg-[#d3e3fd] dark:bg-[#0b57d0]/30 text-[#041e49] dark:text-blue-300 font-bold"
                          >
                            BASE UNIT
                          </span>
                        {/if}
                      </div>
                    </td>
                    <td class="py-3 px-3 font-mono font-bold text-[#0b57d0] dark:text-blue-400">
                      1 {unit.unit_name} = {unit.conversion_factor}{" "}
                      {detailUnits.find((u) => u.is_base_unit)?.unit_name ||
                        "Unit"}
                    </td>
                    <td class="py-3 px-3">
                      <div class="flex items-center gap-1.5 flex-wrap">
                        {#if unit.is_purchase_unit}
                          <span
                            class="px-2 py-0.5 rounded text-[10px] bg-amber-50 dark:bg-amber-950/40 text-amber-700 dark:text-amber-400 border border-amber-200 dark:border-amber-800"
                          >
                            Beli (PO)
                          </span>
                        {/if}
                        {#if unit.is_dispense_unit}
                          <span
                            class="px-2 py-0.5 rounded text-[10px] bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800"
                          >
                            Resep / Jual
                          </span>
                        {/if}
                        {#if !unit.is_purchase_unit && !unit.is_dispense_unit && !unit.is_base_unit}
                          <span class="text-[11px] text-[#747775] dark:text-[#c4c7c5]">
                            Alternatif
                          </span>
                        {/if}
                      </div>
                    </td>
                    <td class="py-3 px-3 text-right">
                      {#if !unit.is_base_unit}
                        <button
                          type="button"
                          onclick={() => handleDeleteUnit(unit)}
                          class="p-1 text-[#c5221f] dark:text-rose-400 hover:bg-[#c5221f]/10 rounded cursor-pointer"
                          title="Hapus Satuan Alternatif"
                        >
                          <Trash2 class="w-3.5 h-3.5" />
                        </button>
                      {:else}
                        <span class="text-[10px] text-[#747775] dark:text-[#c4c7c5] italic">
                          Terkunci
                        </span>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>

      <!-- DETAIL SUB-TAB 3: INTEGRASI AKUNTANSI & COA -->
      {:else if detailActiveTab === "COA"}
        <div class="flex flex-col gap-4 text-xs">
          <!-- Revenue COA -->
          <div
            class="p-4 rounded-2xl bg-blue-50/50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800/60 flex flex-col gap-2"
          >
            <div
              class="flex items-center gap-2 text-blue-900 dark:text-blue-300 font-bold text-xs"
            >
              <Receipt class="w-4 h-4 text-blue-700 dark:text-blue-400" />
              <span>Akun Pendapatan & Penagihan (Billing Module)</span>
            </div>
            <p class="text-xs text-[#444746] dark:text-[#c4c7c5]">
              Diambil secara otomatis dari Kategori: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]"
                >{detailItem.category?.name || "-"}</strong
              >
            </p>
            <div
              class="grid grid-cols-2 gap-3 mt-1 bg-white dark:bg-[#1e1f20] p-3 rounded-xl border border-blue-100 dark:border-blue-900"
            >
              <div>
                <span class="text-[11px] text-[#747775] dark:text-[#c4c7c5] block">
                  COA Akun Pendapatan
                </span>
                <span class="font-mono font-bold text-[#0b57d0] dark:text-blue-400">
                  {detailItem.category?.income_coa_code || "-"}
                </span>
              </div>
              <div>
                <span class="text-[11px] text-[#747775] dark:text-[#c4c7c5] block">
                  Tipe Jurnal Transaksi
                </span>
                <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
                  Kredit Pendapatan Kasir
                </span>
              </div>
            </div>
          </div>

          <!-- Inventory & COGS COA -->
          {#if detailItem.product_line}
            <div
              class="p-4 rounded-2xl bg-amber-50/50 dark:bg-amber-950/20 border border-amber-200 dark:border-amber-800/60 flex flex-col gap-2"
            >
              <div
                class="flex items-center gap-2 text-amber-900 dark:text-amber-300 font-bold text-xs"
              >
                <Layers class="w-4 h-4 text-amber-700 dark:text-amber-400" />
                <span>Akun Persediaan & Harga Pokok Penjualan (HPP)</span>
              </div>
              <p class="text-xs text-[#444746] dark:text-[#c4c7c5]">
                Diambil secara otomatis dari Lini Produk: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]"
                  >{detailItem.product_line.name}</strong
                >
              </p>
              <div
                class="grid grid-cols-2 gap-3 mt-1 bg-white dark:bg-[#1e1f20] p-3 rounded-xl border border-amber-100 dark:border-amber-900"
              >
                <div>
                  <span class="text-[11px] text-[#747775] dark:text-[#c4c7c5] block">
                    COA Akun Persediaan (Aset)
                  </span>
                  <span class="font-mono font-bold text-amber-800 dark:text-amber-400">
                    {detailItem.product_line.inventory_coa_code || "-"}
                  </span>
                </div>
                <div>
                  <span class="text-[11px] text-[#747775] dark:text-[#c4c7c5] block">
                    COA Akun Beban Pokok (HPP)
                  </span>
                  <span class="font-mono font-bold text-amber-800 dark:text-amber-400">
                    {detailItem.product_line.cogs_coa_code || "-"}
                  </span>
                </div>
              </div>
            </div>
          {/if}
        </div>
      {/if}

      <!-- Modal Footer inside content -->
      <div
        class="pt-4 border-t border-[#e1e5ea] dark:border-[#444746] flex justify-end gap-2"
      >
        <button
          type="button"
          onclick={() => {
            onClose();
            onEdit(detailItem!);
          }}
          class="px-4 py-2 rounded-xl border border-[#e1e5ea] dark:border-[#444746] bg-white dark:bg-[#2b2b2b] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#f0f4f9] dark:hover:bg-[#383b3d] flex items-center gap-1.5 cursor-pointer shadow-xs"
        >
          <Edit3 class="w-3.5 h-3.5" />
          <span>Edit Item</span>
        </button>
        <M3Button variant="tonal" onclick={onClose}>
          <span>Tutup</span>
        </M3Button>
      </div>
    {/if}
  </div>
</M3Modal>
