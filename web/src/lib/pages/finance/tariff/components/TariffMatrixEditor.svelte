<script lang="ts">
  import { onMount } from 'svelte';
  import {
    BookOpen,
    Save,
    CheckCircle2,
    Sliders
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import TariffLookupInput from '$lib/components/TariffLookupInput.svelte';
  import TariffMatrixClassCard, {
    type ClassMatrixRowState,
    type ComponentRowState
  } from './TariffMatrixClassCard.svelte';
  import { confirmDialog } from '$lib/stores/confirm.svelte';
  import {
    getPricePlans,
    getPricePlanItems,
    batchUpsertPricePlanItems
  } from '$lib/api/finance/price_plan';
  import { getTariffClasses, getTariffComponents } from '$lib/api/finance/tariff';
  import type {
    PricePlanRecord,
    PricePlanItemRecord,
    AddPricePlanItemDTO,
    ItemComponentDTO
  } from '$lib/types/finance/price_plan';
  import type { TariffClassRecord, TariffComponentRecord } from '$lib/types/finance/tariff';
  import type { ItemSummaryRecord } from '$lib/types/master/item';

  interface Props {
    showToast: (msg: string, isError?: boolean) => void;
  }

  let { showToast }: Props = $props();

  // Master Data
  let pricePlans = $state<PricePlanRecord[]>([]);
  let tariffClasses = $state<TariffClassRecord[]>([]);
  let tariffComponents = $state<TariffComponentRecord[]>([]);

  // Selection State
  let selectedPlanId = $state<string>('');
  let selectedItemId = $state<string>('');
  let selectedItemObj = $state<ItemSummaryRecord | null>(null);

  // Loading States
  let isInitialLoading = $state(true);
  let isLoadingMatrix = $state(false);
  let isSavingMatrix = $state(false);

  // Selected Plan Object
  let selectedPlan = $derived(pricePlans.find((p) => p.id === selectedPlanId) || null);
  let selectedItemTitle = $derived(selectedItemObj?.name || 'Tindakan Medis');

  let matrixRows = $state<ClassMatrixRowState[]>([]);

  // Load foundational master records (fast & lightweight)
  onMount(async () => {
    isInitialLoading = true;
    try {
      const [plansRes, classesRes, compsRes] = await Promise.all([
        getPricePlans({ limit: 100 }),
        getTariffClasses({ is_active: true, limit: 100 }),
        getTariffComponents({ is_active: true, limit: 100 })
      ]);

      pricePlans = plansRes.data || [];
      tariffClasses = classesRes.data || [];
      tariffComponents = compsRes.data || [];

      // Auto select first active or default plan if available
      if (pricePlans.length > 0) {
        const defaultPlan = pricePlans.find((p) => p.is_default && p.status === 'ACTIVE') || pricePlans[0];
        selectedPlanId = defaultPlan.id;
      }
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat master data tarif', true);
    } finally {
      isInitialLoading = false;
    }
  });

  // Watch for changes on selected plan or item to load matrix
  $effect(() => {
    if (selectedPlanId && selectedItemId) {
      loadMatrixForSelectedItem(selectedPlanId, selectedItemId);
    } else {
      matrixRows = [];
    }
  });

  async function loadMatrixForSelectedItem(planId: string, itemId: string) {
    isLoadingMatrix = true;
    try {
      // 1. Ambil existing konfigurasi item pada buku tarif ini untuk tindakan tersebut
      const existingRes = await getPricePlanItems(planId, {
        item_id: itemId,
        limit: 100
      });
      const existingItems = existingRes.data || [];

      // 2. Buat map existing berdasarkan tariff_class_id
      const existingMap = new Map<string, PricePlanItemRecord>();
      for (const it of existingItems) {
        existingMap.set(it.tariff_class_id, it);
      }

      // 3. Bangun baris matriks untuk setiap kelas rawat
      const newRows: ClassMatrixRowState[] = tariffClasses.map((cls) => {
        const existing = existingMap.get(cls.id);
        if (existing) {
          // Pre-fill dari data yang sudah tersimpan
          return {
            tariff_class_id: cls.id,
            tariff_class_code: cls.code,
            tariff_class_name: cls.name,
            existing_item_id: existing.id,
            total_base_price: Number(existing.total_base_price),
            total_cito_price: existing.total_cito_price ? Number(existing.total_cito_price) : null,
            is_active: existing.is_active,
            is_enabled: true,
            components: (existing.components || []).map((c) => ({
              component_id: c.component_id,
              base_amount: Number(c.base_amount),
              cito_amount: c.cito_amount ? Number(c.cito_amount) : null,
              coa_code: c.coa_code || null
            }))
          };
        }

        // Default template jika belum pernah dikonfigurasi:
        // Siapkan komponen default RS jika ada (misal Sarana RS dan Jasa Medis Dokter)
        const defaultComponents: ComponentRowState[] = [];
        if (tariffComponents.length > 0) {
          const sarana = tariffComponents.find((c) => c.component_type === 'SARANA') || tariffComponents[0];
          defaultComponents.push({
            component_id: sarana.id,
            base_amount: 0,
            cito_amount: null,
            coa_code: sarana.default_coa_code || null
          });

          const medis = tariffComponents.find((c) => c.component_type === 'MEDIS_DOKTER');
          if (medis && medis.id !== sarana.id) {
            defaultComponents.push({
              component_id: medis.id,
              base_amount: 0,
              cito_amount: null,
              coa_code: medis.default_coa_code || null
            });
          }
        }

        return {
          tariff_class_id: cls.id,
          tariff_class_code: cls.code,
          tariff_class_name: cls.name,
          existing_item_id: null,
          total_base_price: 0,
          total_cito_price: null,
          is_active: true,
          is_enabled: false, // belum diaktifkan sebelum staf mengisi harga
          components: defaultComponents
        };
      });

      matrixRows = newRows;
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat konfigurasi tarif kelas', true);
    } finally {
      isLoadingMatrix = false;
    }
  }

  // Auto calculate CITO price for a class row based on plan default CITO %
  function handleCalculateCito(row: ClassMatrixRowState) {
    const citoPercent = selectedPlan?.default_cito_percent || 25;
    const base = Number(row.total_base_price) || 0;
    row.total_cito_price = Math.round(base * (1 + citoPercent / 100));

    // Calculate components CITO as well
    for (const comp of row.components) {
      const cBase = Number(comp.base_amount) || 0;
      comp.cito_amount = Math.round(cBase * (1 + citoPercent / 100));
    }
  }

  // Add component to a class row
  function handleAddComponent(row: ClassMatrixRowState) {
    // Cari komponen yang belum dipakai di kelas ini
    const usedIds = new Set(row.components.map((c) => c.component_id));
    const available = tariffComponents.find((c) => !usedIds.has(c.id)) || tariffComponents[0];
    if (available) {
      row.components.push({
        component_id: available.id,
        base_amount: 0,
        cito_amount: null,
        coa_code: available.default_coa_code || null
      });
    }
  }

  // Remove component from a class row
  function handleRemoveComponent(row: ClassMatrixRowState, compIdx: number) {
    if (row.components.length <= 1) {
      showToast('Minimal satu rincian komponen biaya harus dipertahankan', true);
      return;
    }
    row.components.splice(compIdx, 1);
  }

  // Auto-balance components for a class row
  function handleAutoBalance(row: ClassMatrixRowState) {
    if (row.components.length === 0) return;
    const target = Number(row.total_base_price) || 0;
    let otherSum = 0;
    for (let i = 0; i < row.components.length - 1; i++) {
      otherSum += Number(row.components[i].base_amount) || 0;
    }
    const remainder = target - otherSum;
    row.components[row.components.length - 1].base_amount = Math.max(0, remainder);
  }

  // Copy this class row configuration to all other classes below it
  async function handleCopyToAllClasses(sourceRow: ClassMatrixRowState) {
    const ok = await confirmDialog({
      title: 'Salin Pengaturan Kelas',
      message: `Terapkan susunan tarif dan komponen kelas "${sourceRow.tariff_class_name}" ke seluruh kelas lainnya?`,
      confirmText: 'Terapkan',
      variant: 'primary'
    });
    if (!ok) return;

    for (const row of matrixRows) {
      if (row.tariff_class_id === sourceRow.tariff_class_id) continue;
      row.total_base_price = sourceRow.total_base_price;
      row.total_cito_price = sourceRow.total_cito_price;
      row.is_active = sourceRow.is_active;
      row.is_enabled = true;
      row.components = sourceRow.components.map((c) => ({
        component_id: c.component_id,
        base_amount: c.base_amount,
        cito_amount: c.cito_amount,
        coa_code: c.coa_code
      }));
    }

    showToast(`Konfigurasi kelas ${sourceRow.tariff_class_name} berhasil diterapkan ke semua kelas`);
  }

  // Save entire matrix (Batch Upsert)
  async function handleSaveAllClasses() {
    if (!selectedPlanId || !selectedItemId) {
      showToast('Pilih buku tarif dan tindakan medis terlebih dahulu', true);
      return;
    }

    // Filter baris yang di-enable (atau yang memiliki harga > 0)
    const activeRowsToSave = matrixRows.filter((r) => r.is_enabled || r.total_base_price > 0);

    if (activeRowsToSave.length === 0) {
      showToast('Tidak ada kelas yang diaktifkan atau diberi nominal tarif untuk disimpan', true);
      return;
    }

    // Validasi balancing komponen untuk setiap baris
    for (const row of activeRowsToSave) {
      if (row.total_base_price < 0) {
        showToast(`Nominal tarif pada kelas ${row.tariff_class_name} tidak valid`, true);
        return;
      }

      if (row.components.length === 0) {
        showToast(`Kelas ${row.tariff_class_name} wajib memiliki minimal satu komponen biaya`, true);
        return;
      }

      let sum = 0;
      for (const comp of row.components) {
        sum += Number(comp.base_amount) || 0;
      }

      if (Math.abs(sum - row.total_base_price) > 0.01) {
        showToast(
          `Pecahan komponen kelas "${row.tariff_class_name}" tidak seimbang! Total: Rp ${row.total_base_price.toLocaleString('id-ID')}, Pecahan: Rp ${sum.toLocaleString('id-ID')}. Klik tombol Seimbangkan.`,
          true
        );
        return;
      }
    }

    isSavingMatrix = true;
    try {
      const payloadItems: AddPricePlanItemDTO[] = activeRowsToSave.map((r) => ({
        item_id: selectedItemId,
        tariff_class_id: r.tariff_class_id,
        total_base_price: Number(r.total_base_price),
        total_cito_price: r.total_cito_price ? Number(r.total_cito_price) : null,
        is_active: r.is_active,
        components: r.components.map((c) => ({
          component_id: c.component_id,
          base_amount: Number(c.base_amount),
          cito_amount: c.cito_amount ? Number(c.cito_amount) : null,
          coa_code: c.coa_code || null
        }))
      }));

      await batchUpsertPricePlanItems(selectedPlanId, payloadItems);
      showToast(
        `Tarif tindakan "${selectedItemTitle}" untuk ${payloadItems.length} kelas berhasil disimpan ke buku tarif!`
      );

      // Reload matrix to show updated persisted state
      loadMatrixForSelectedItem(selectedPlanId, selectedItemId);
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan tarif multi-kelas', true);
    } finally {
      isSavingMatrix = false;
    }
  }
</script>

<div class="flex flex-col gap-5">
  <!-- Header Card & Intro -->
  <div class="bg-gradient-to-r from-blue-50/80 via-indigo-50/60 to-purple-50/50 dark:from-[#1b2234] dark:via-[#1e1f2b] dark:to-[#1a1f2c] p-5 rounded-3xl border border-blue-200/80 dark:border-indigo-900/50 shadow-xs">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2.5">
          <div class="w-10 h-10 rounded-2xl bg-blue-600 dark:bg-[#0b57d0] text-white flex items-center justify-center font-bold shadow-xs shrink-0">
            <Sliders class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-bold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">
              Editor Matriks Tarif & Komponen Multi-Kelas
            </h3>
            <div class="flex items-center gap-2 mt-0.5">
              <span class="text-[11px] font-mono font-medium px-2.5 py-0.5 rounded-full bg-blue-100 text-blue-800 border border-blue-300 dark:bg-blue-950/70 dark:text-blue-300 dark:border-blue-800/60">
                Service-Centric Batch Matrix
              </span>
              <span class="text-xs text-[#5f6368] dark:text-[#8e918f]">• Penetapan tarif serentak untuk seluruh kelas rawat</span>
            </div>
          </div>
        </div>
        <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-2.5 max-w-3xl leading-relaxed">
          Pilih buku tarif dan tindakan medis untuk memuat seluruh kelas rawat inap/jalan (VVIP s/d Kelas 3). Tentukan nominal dasar, override CITO darurat, dan pecahan pendapatan (jasa dokter & sarana RS) secara dinamis dalam satu layar kerja.
        </p>
      </div>

      {#if selectedPlan && selectedItemId && matrixRows.length > 0}
        <div class="flex items-center gap-2">
          <M3Button
            variant="filled"
            onclick={handleSaveAllClasses}
            disabled={isSavingMatrix || isLoadingMatrix}
            class="!h-10 !px-4 !rounded-xl !bg-[#0b57d0] !text-white text-xs font-semibold flex items-center gap-2 shadow-sm"
          >
            {#if isSavingMatrix}
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              <span>Menyimpan...</span>
            {:else}
              <Save class="w-4 h-4" />
              <span>Simpan Tarif Seluruh Kelas</span>
            {/if}
          </M3Button>
        </div>
      {/if}
    </div>

    <!-- Parameter Selector Bar -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4 pt-4 border-t border-blue-200/60 dark:border-[#353738]">
      <!-- 1. Select Buku Tarif -->
      <div class="flex flex-col gap-1.5">
        <label for="matrix-plan-select" class="text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] flex items-center gap-1.5">
          <BookOpen class="w-3.5 h-3.5 text-indigo-600 dark:text-indigo-400" />
          <span>1. Pilih Buku Tarif Target <span class="text-red-500">*</span></span>
        </label>
        <div class="relative">
          <select
            id="matrix-plan-select"
            bind:value={selectedPlanId}
            class="w-full h-10 px-3 pr-8 rounded-xl border border-[#c4c7c5] dark:border-[#353738] bg-white dark:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-none focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:ring-1 focus:ring-[#0b57d0]"
          >
            <option value="" disabled>-- Pilih Buku Tarif --</option>
            {#each pricePlans as plan (plan.id)}
              <option value={plan.id}>
                [{plan.status}] {plan.code} — {plan.name} {plan.is_default ? '(DEFAULT RS)' : ''}
              </option>
            {/each}
          </select>
        </div>
        {#if selectedPlan}
          <div class="flex items-center gap-2 text-[11px] text-[#5f6368] dark:text-[#8e918f] mt-0.5">
            <span class="font-medium text-indigo-700 dark:text-indigo-300">Berlaku: {selectedPlan.effective_from} {selectedPlan.effective_to ? `s/d ${selectedPlan.effective_to}` : '(Aktif)'}</span>
            <span>• Default CITO: +{selectedPlan.default_cito_percent}%</span>
            {#if selectedPlan.status === 'ACTIVE'}
              <span class="px-2 py-0.5 rounded-full bg-emerald-100 text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border dark:border-emerald-800/60 font-semibold text-[10px]">AKTIF DI PELAYANAN</span>
            {/if}
          </div>
        {/if}
      </div>

      <!-- 2. Select Tindakan Medis / Layanan via General Lookup -->
      <div class="flex flex-col gap-1.5">
        <TariffLookupInput
          id="matrix-tariff-lookup"
          label="2. Pilih Layanan / Tindakan Medis"
          required={true}
          bind:value={selectedItemId}
          placeholder="Cari ribuan tindakan medis katalog (misal: Konsultasi, USG, Darah)..."
          onSelect={(itm) => {
            selectedItemObj = itm;
          }}
        />
      </div>
    </div>
  </div>

  <!-- Loading State -->
  {#if isLoadingMatrix}
    <div class="py-16 flex flex-col items-center justify-center gap-3 bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738]">
      <div class="w-8 h-8 border-3 border-[#0b57d0] border-t-transparent rounded-full animate-spin"></div>
      <span class="text-xs font-medium text-[#444746] dark:text-[#c4c7c5]">Memuat konfigurasi kelas tindakan...</span>
    </div>

  <!-- Empty State (Belum memilih tindakan) -->
  {:else if !selectedPlanId || !selectedItemId}
    <div class="py-16 px-4 flex flex-col items-center justify-center text-center gap-3 bg-white dark:bg-[#1e1f20] rounded-3xl border border-[#e1e5ea] dark:border-[#353738]">
      <div class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-600 dark:bg-blue-950/60 dark:text-blue-300 flex items-center justify-center">
        <Sliders class="w-6 h-6" />
      </div>
      <div>
        <h4 class="text-sm font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">Pilih Buku Tarif dan Tindakan Medis</h4>
        <p class="text-xs text-[#747775] dark:text-[#8e918f] mt-1 max-w-md">
          Pilih buku tarif target dan nama tindakan medis di atas untuk memunculkan tabel konfigurasi matriks seluruh kelas perawatan.
        </p>
      </div>
    </div>

  <!-- Matrix Loaded: Daftar Kelas Perawatan -->
  {:else}
    <div class="flex flex-col gap-4">
      <div class="flex items-center justify-between text-xs text-[#444746] dark:text-[#c4c7c5] px-1">
        <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
          Konfigurasi Tarif: <span class="text-blue-700 dark:text-blue-400">{selectedItemTitle}</span> ({matrixRows.length} Kelas Tersedia)
        </span>
        <span class="text-[11px] text-[#747775] dark:text-[#8e918f]">
          * Pastikan jumlah pecahan komponen biaya sama dengan total tarif dasar kelas.
        </span>
      </div>

      <!-- Iterasi Seluruh Kelas Perawatan -->
      {#each matrixRows as row, rIdx (row.tariff_class_id)}
        <TariffMatrixClassCard
          {row}
          {rIdx}
          {selectedPlan}
          {tariffComponents}
          onAutoBalance={handleAutoBalance}
          onCopyToAllClasses={handleCopyToAllClasses}
          onCalculateCito={handleCalculateCito}
          onAddComponent={handleAddComponent}
          onRemoveComponent={handleRemoveComponent}
        />
      {/each}

      <!-- Bottom Save Action Bar -->
      <div class="sticky bottom-4 p-4 rounded-2xl bg-white/95 dark:bg-[#1e1f20]/95 backdrop-blur-md border border-[#c4c7c5] dark:border-[#353738] shadow-lg flex items-center justify-between gap-4 z-10">
        <div class="flex items-center gap-2 text-xs text-[#444746] dark:text-[#c4c7c5]">
          <CheckCircle2 class="w-4 h-4 text-emerald-600 dark:text-emerald-400 shrink-0" />
          <span>
            Siap menyimpan seluruh pengaturan tarif kelas untuk <b class="text-[#1f1f1f] dark:text-[#e3e3e3]">{selectedItemTitle}</b> ke buku tarif <b class="text-[#1f1f1f] dark:text-[#e3e3e3]">{selectedPlan?.name}</b>.
          </span>
        </div>

        <div class="flex items-center gap-2">
          <M3Button
            variant="filled"
            onclick={handleSaveAllClasses}
            disabled={isSavingMatrix || isLoadingMatrix}
            class="!h-10 !px-5 !rounded-xl !bg-[#0b57d0] !text-white text-xs font-semibold flex items-center gap-2 shadow-sm"
          >
            {#if isSavingMatrix}
              <div class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              <span>Menyimpan Seluruh Kelas...</span>
            {:else}
              <Save class="w-4 h-4" />
              <span>Simpan Seluruh Kelas (Batch Upsert)</span>
            {/if}
          </M3Button>
        </div>
      </div>
    </div>
  {/if}
</div>
