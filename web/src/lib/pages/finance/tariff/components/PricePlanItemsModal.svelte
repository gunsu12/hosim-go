<script lang="ts">
  import { onMount } from "svelte";
  import {
    Plus,
    Search,
    Trash2,
    Edit3,
    Coins,
    Layers,
    Lock,
  } from "@lucide/svelte";
  import M3Modal from "$lib/components/m3/M3Modal.svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3Select from "$lib/components/m3/M3Select.svelte";
  import PricePlanItemFormModal from "./PricePlanItemFormModal.svelte";
  import { confirmDialog } from "$lib/stores/confirm.svelte";
  import {
    getPricePlanItems,
    deletePricePlanItem,
  } from "$lib/api/finance/price_plan";
  import {
    getTariffClasses,
    getTariffComponents,
  } from "$lib/api/finance/tariff";
  import { getItems } from "$lib/api/catalog/item";
  import type {
    PricePlanRecord,
    PricePlanItemRecord,
  } from "$lib/types/finance/price_plan";
  import type {
    TariffClassRecord,
    TariffComponentRecord,
  } from "$lib/types/finance/tariff";
  import type { ItemSummaryRecord } from "$lib/types/master/item";

  interface Props {
    isOpen: boolean;
    plan: PricePlanRecord | null;
    onClose: () => void;
    showToast: (msg: string, isError?: boolean) => void;
  }

  let { isOpen, plan, onClose, showToast }: Props = $props();

  let items = $state<PricePlanItemRecord[]>([]);
  let isLoading = $state(false);
  let searchQuery = $state("");
  let filterClassId = $state("");
  let currentPage = $state(1);
  let totalPages = $state(1);
  let totalCount = $state(0);

  // Reference Options
  let tariffClasses = $state<TariffClassRecord[]>([]);
  let tariffComponents = $state<TariffComponentRecord[]>([]);
  let tariffCatalogItems = $state<ItemSummaryRecord[]>([]);

  let classFilterOptions = $derived([
    { value: "", label: "Semua Kelas Tarif" },
    ...tariffClasses.map((tc) => ({
      value: tc.id,
      label: tc.name,
    })),
  ]);

  // Item Form Modal State
  let showItemForm = $state(false);
  let itemFormMode = $state<"ADD" | "EDIT">("ADD");
  let itemToEdit = $state<PricePlanItemRecord | null>(null);

  onMount(() => {
    Promise.all([
      getTariffClasses({ limit: 100 }),
      getTariffComponents({ limit: 100 }),
      getItems({ item_type: "TARIFF", limit: 100 }),
    ])
      .then(([classRes, compRes, itemRes]) => {
        tariffClasses = classRes.data || [];
        tariffComponents = compRes.data || [];
        tariffCatalogItems = itemRes.data || [];
      })
      .catch((err) => {
        console.error("Failed to load tariff references:", err);
      });
  });

  async function loadPlanItems() {
    if (!plan) return;
    isLoading = true;
    try {
      const res = await getPricePlanItems(plan.id, {
        page: currentPage,
        limit: 20,
        search: searchQuery || undefined,
        tariff_class_id: filterClassId || undefined,
      });
      items = res.data;
      totalCount = res.meta.total_items;
      totalPages = res.meta.total_pages;
    } catch (err: any) {
      showToast(err.message || "Gagal memuat item tarif", true);
    } finally {
      isLoading = false;
    }
  }

  $effect(() => {
    if (isOpen && plan) {
      loadPlanItems();
    }
  });

  function openAddItem() {
    itemFormMode = "ADD";
    itemToEdit = null;
    showItemForm = true;
  }

  function openEditItem(item: PricePlanItemRecord) {
    itemFormMode = "EDIT";
    itemToEdit = item;
    showItemForm = true;
  }

  function handleItemSaved(msg: string) {
    showItemForm = false;
    showToast(msg);
    loadPlanItems();
  }

  async function handleDeleteItem(item: PricePlanItemRecord) {
    if (!plan) return;
    const ok = await confirmDialog({
      title: "Hapus Item Tarif",
      message: `Hapus item tarif ${item.item_name} pada kelas ${item.tariff_class_name}?`,
      confirmText: "Hapus",
      variant: "danger",
    });
    if (!ok) return;

    try {
      await deletePricePlanItem(plan.id, item.id);
      showToast("Item tarif berhasil dihapus");
      loadPlanItems();
    } catch (err: any) {
      showToast(err.message || "Gagal menghapus item tarif", true);
    }
  }
</script>

<M3Modal
  isOpen={isOpen && !!plan}
  title={plan?.name || "Item Buku Tarif"}
  subtitle="Kelola matriks tindakan medis, kelas layanan, dan pecahan komponen biaya"
  icon={Coins}
  size="4xl"
  {onClose}
>
  {#snippet headerExtra()}
    {#if plan}
      <div class="flex items-center gap-2 mr-2">
        <span
          class="font-mono font-bold text-xs text-[#0b57d0] dark:text-[#a8c7fa] bg-[#e8f0fe] dark:bg-[#0b57d0]/20 px-2 py-0.5 rounded"
        >
          {plan.code}
        </span>
        <span
          class="text-xs px-2.5 py-0.5 rounded-full font-bold {plan.status ===
          'ACTIVE'
            ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300'
            : 'bg-amber-100 text-amber-800 dark:bg-amber-950/60 dark:text-amber-300'}"
        >
          {plan.status}
        </span>
        {#if plan.status !== "DRAFT"}
          <span
            class="inline-flex items-center gap-1 text-[11px] text-[#747775] dark:text-[#8e918f] italic"
          >
            <Lock class="w-3 h-3" />
            Terkunci
          </span>
        {/if}
      </div>
    {/if}
  {/snippet}

  {#if plan}
    <!-- Controls & Filter Bar -->
    <div
      class="flex flex-wrap items-center justify-between gap-3 p-3 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738]"
    >
      <div class="flex items-center gap-2 flex-1 min-w-[240px]">
        <div class="relative flex-1">
          <Search
            class="w-3.5 h-3.5 text-[#747775] dark:text-[#8e918f] absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none"
          />
          <input
            type="text"
            bind:value={searchQuery}
            onkeydown={(e) => e.key === "Enter" && loadPlanItems()}
            placeholder="Cari tindakan atau kelas (tekan Enter)..."
            class="w-full h-9 pl-9 pr-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:outline-hidden"
          />
        </div>

        <div class="w-48">
          <M3Select
            bind:value={filterClassId}
            options={classFilterOptions}
            onchange={loadPlanItems}
            placeholder="Semua Kelas Tarif"
          />
        </div>
      </div>

      {#if plan.status === "DRAFT"}
        <M3Button variant="filled" onclick={openAddItem}>
          <Plus class="w-4 h-4 mr-1" />
          <span>Tambah Item Tarif</span>
        </M3Button>
      {/if}
    </div>

    <!-- Table Content Area -->
    <div>
      {#if isLoading}
        <div
          class="py-16 text-center text-[#747775] dark:text-[#8e918f] flex flex-col items-center justify-center gap-2"
        >
          <div
            class="w-6 h-6 border-2 border-[#0b57d0] dark:border-[#a8c7fa] border-t-transparent rounded-full animate-spin"
          ></div>
          <span class="text-xs">Memuat item tarif...</span>
        </div>
      {:else if items.length === 0}
        <div
          class="py-16 text-center text-[#747775] dark:text-[#8e918f] flex flex-col items-center justify-center gap-2"
        >
          <Layers class="w-10 h-10 text-[#b0b8c4] dark:text-[#555]" />
          <span class="font-bold text-xs text-[#1f1f1f] dark:text-[#e3e3e3]"
            >Belum ada item tarif dalam buku tarif ini</span
          >
          <p class="text-[11px] text-[#747775] dark:text-[#8e918f]">
            {plan.status === "DRAFT"
              ? 'Klik tombol "+ Tambah Item Tarif" di atas untuk memasukkan tindakan medis dan pecahan komponen biayanya.'
              : "Buku tarif ini tidak memiliki item."}
          </p>
        </div>
      {:else}
        <div
          class="border border-[#e1e5ea] dark:border-[#353738] rounded-2xl overflow-hidden shadow-2xs"
        >
          <table class="w-full text-left text-xs border-collapse">
            <thead>
              <tr
                class="bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] font-semibold border-b border-[#e1e5ea] dark:border-[#353738]"
              >
                <th class="py-3 px-3 w-10 text-center">No</th>
                <th class="py-3 px-3">Tindakan Medis</th>
                <th class="py-3 px-3">Kelas Tarif</th>
                <th class="py-3 px-3 text-right">Tarif Dasar</th>
                <th class="py-3 px-3 text-right"
                  >Tarif CITO (+{plan.default_cito_percent}%)</th
                >
                <th class="py-3 px-3">Pecahan Komponen Biaya</th>
                {#if plan.status === "DRAFT"}
                  <th class="py-3 px-3 w-20 text-right">Aksi</th>
                {/if}
              </tr>
            </thead>
            <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
              {#each items as item, idx (item.id)}
                <tr
                  class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/50 transition-colors"
                >
                  <td
                    class="py-3 px-3 text-center font-mono text-[#747775] dark:text-[#8e918f]"
                  >
                    {idx + 1}
                  </td>
                  <td class="py-3 px-3">
                    <div
                      class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]"
                    >
                      {item.item_name}
                    </div>
                    <span
                      class="font-mono text-[10px] text-[#0b57d0] dark:text-[#a8c7fa]"
                      >{item.item_code}</span
                    >
                  </td>
                  <td class="py-3 px-3">
                    <span
                      class="px-2 py-0.5 rounded bg-blue-50 dark:bg-blue-950/50 text-blue-700 dark:text-blue-300 font-bold border border-blue-200 dark:border-blue-800 text-[11px]"
                    >
                      {item.tariff_class_name || item.tariff_class_code}
                    </span>
                  </td>
                  <td
                    class="py-3 px-3 text-right font-mono font-bold text-[#1f1f1f] dark:text-[#e3e3e3]"
                  >
                    Rp {item.total_base_price.toLocaleString("id-ID")}
                  </td>
                  <td
                    class="py-3 px-3 text-right font-mono text-[#747775] dark:text-[#8e918f]"
                  >
                    {#if item.total_cito_price}
                      Rp {item.total_cito_price.toLocaleString("id-ID")}
                    {:else}
                      Rp {Math.round(
                        item.total_base_price *
                          (1 + plan.default_cito_percent / 100),
                      ).toLocaleString("id-ID")}
                    {/if}
                  </td>
                  <td class="py-3 px-3">
                    <div class="flex flex-wrap gap-1">
                      {#each item.components as c (c.id)}
                        <span
                          class="px-1.5 py-0.5 rounded bg-[#f0f4f9] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] text-[10px] text-[#333] dark:text-[#ccc]"
                        >
                          <strong
                            >{c.component_name || c.component_code}:</strong
                          >
                          Rp {c.base_amount.toLocaleString("id-ID")}
                        </span>
                      {/each}
                    </div>
                  </td>
                  {#if plan.status === "DRAFT"}
                    <td class="py-3 px-3 text-right">
                      <div class="flex items-center justify-end gap-1">
                        <button
                          type="button"
                          onclick={() => openEditItem(item)}
                          class="p-1 rounded text-[#0b57d0] dark:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#0b57d0]/20 cursor-pointer"
                          title="Edit Item & Komponen"
                        >
                          <Edit3 class="w-3.5 h-3.5" />
                        </button>
                        <button
                          type="button"
                          onclick={() => handleDeleteItem(item)}
                          class="p-1 rounded text-[#c5221f] dark:text-rose-400 hover:bg-[#fce8e6] dark:hover:bg-rose-950/50 cursor-pointer"
                          title="Hapus Item"
                        >
                          <Trash2 class="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </td>
                  {/if}
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  {/if}

  {#snippet footer()}
    <div class="w-full flex items-center justify-between">
      <span class="text-xs text-[#747775] dark:text-[#8e918f]">
        Total <strong>{totalCount}</strong> kombinasi item tindakan per kelas
      </span>
      <M3Button variant="text" onclick={onClose}>Tutup</M3Button>
    </div>
  {/snippet}
</M3Modal>

<!-- SUB-MODAL FORM: TAMBAH / EDIT ITEM TARIF -->
{#if showItemForm && plan}
  <PricePlanItemFormModal
    isOpen={showItemForm}
    mode={itemFormMode}
    planId={plan.id}
    {itemToEdit}
    {tariffClasses}
    {tariffComponents}
    {tariffCatalogItems}
    onClose={() => (showItemForm = false)}
    onSuccess={handleItemSaved}
  />
{/if}
