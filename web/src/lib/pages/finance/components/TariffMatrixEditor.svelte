<script lang="ts">
  import { onMount } from 'svelte';
  import {
    BookOpen,
    Layers,
    Tag,
    Calculator,
    Save,
    RefreshCw,
    Plus,
    Trash2,
    Copy,
    CheckCircle2,
    AlertTriangle,
    Coins,
    Sliders,
    Building2,
    UserCheck,
    CheckCheck,
    ArrowRight,
    HelpCircle
  } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import TariffLookupInput from '$lib/components/TariffLookupInput.svelte';
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

  // ==========================================
  // MATRIX ROW DATA STRUCTURE
  // ==========================================
  interface ComponentRowState {
    component_id: string;
    base_amount: number;
    cito_amount: number | null;
    coa_code: string | null;
  }

  interface ClassMatrixRowState {
    tariff_class_id: string;
    tariff_class_code: string;
    tariff_class_name: string;
    existing_item_id: string | null; // null jika belum ada di database
    total_base_price: number;
    total_cito_price: number | null;
    is_active: boolean;
    is_enabled: boolean; // apakah kelas ini akan disimpan/disertakan
    components: ComponentRowState[];
  }

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
  function handleCopyToAllClasses(sourceRow: ClassMatrixRowState) {
    if (!confirm(`Terapkan susunan tarif dan komponen kelas "${sourceRow.tariff_class_name}" ke seluruh kelas lainnya?`)) {
      return;
    }

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

  function formatIDR(val: number): string {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0
    }).format(val || 0);
  }
</script>

<div class="flex flex-col gap-5">
  <!-- Header Card & Intro -->
  <div class="bg-gradient-to-r from-blue-50/70 via-indigo-50/50 to-purple-50/40 p-4.5 rounded-3xl border border-blue-200/80 shadow-xs">
    <div class="flex flex-wrap items-center justify-between gap-4">
      <div>
        <div class="flex items-center gap-2.5">
          <div class="w-9 h-9 rounded-2xl bg-blue-600 text-white flex items-center justify-center font-bold shadow-xs">
            <Sliders class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-bold text-[#1f1f1f] tracking-tight">
              Editor Matriks Tarif & Komponen Multi-Kelas
            </h3>
            <div class="flex items-center gap-2 mt-0.5">
              <span class="text-[11px] font-mono font-medium px-2 py-0.5 rounded-full bg-blue-100 text-blue-800 border border-blue-300">
                Service-Centric Batch Matrix
              </span>
              <span class="text-xs text-[#5f6368]">• Penetapan tarif serentak untuk seluruh kelas rawat</span>
            </div>
          </div>
        </div>
        <p class="text-xs text-[#444746] mt-2 max-w-3xl leading-relaxed">
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
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4 pt-4 border-t border-blue-200/60">
      <!-- 1. Select Buku Tarif -->
      <div class="flex flex-col gap-1.5">
        <label for="matrix-plan-select" class="text-xs font-semibold text-[#1f1f1f] flex items-center gap-1.5">
          <BookOpen class="w-3.5 h-3.5 text-indigo-600" />
          <span>1. Pilih Buku Tarif Target <span class="text-red-500">*</span></span>
        </label>
        <div class="relative">
          <select
            id="matrix-plan-select"
            bind:value={selectedPlanId}
            class="w-full h-10 px-3 pr-8 rounded-xl border border-[#c4c7c5] bg-white text-xs text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0]"
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
          <div class="flex items-center gap-2 text-[11px] text-[#5f6368] mt-0.5">
            <span class="font-medium text-indigo-700">Berlaku: {selectedPlan.effective_from} {selectedPlan.effective_to ? `s/d ${selectedPlan.effective_to}` : '(Aktif)'}</span>
            <span>• Default CITO: +{selectedPlan.default_cito_percent}%</span>
            {#if selectedPlan.status === 'ACTIVE'}
              <span class="px-1.5 py-0.2 rounded bg-emerald-100 text-emerald-800 font-semibold text-[10px]">AKTIF DI PELAYANAN</span>
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
    <div class="py-16 flex flex-col items-center justify-center gap-3 bg-white rounded-3xl border border-[#e1e5ea]">
      <div class="w-8 h-8 border-3 border-[#0b57d0] border-t-transparent rounded-full animate-spin"></div>
      <span class="text-xs font-medium text-[#444746]">Memuat konfigurasi kelas tindakan...</span>
    </div>

  <!-- Empty State (Belum memilih tindakan) -->
  {:else if !selectedPlanId || !selectedItemId}
    <div class="py-16 px-4 flex flex-col items-center justify-center text-center gap-3 bg-white rounded-3xl border border-[#e1e5ea]">
      <div class="w-12 h-12 rounded-2xl bg-blue-50 text-blue-600 flex items-center justify-center">
        <Sliders class="w-6 h-6" />
      </div>
      <div>
        <h4 class="text-sm font-bold text-[#1f1f1f]">Pilih Buku Tarif dan Tindakan Medis</h4>
        <p class="text-xs text-[#747775] mt-1 max-w-md">
          Pilih buku tarif target dan nama tindakan medis di atas untuk memunculkan tabel konfigurasi matriks seluruh kelas perawatan.
        </p>
      </div>
    </div>

  <!-- Matrix Loaded: Daftar Kelas Perawatan -->
  {:else}
    <div class="flex flex-col gap-4">
      <div class="flex items-center justify-between text-xs text-[#444746] px-1">
        <span class="font-semibold text-[#1f1f1f]">
          Konfigurasi Tarif: <span class="text-blue-700">{selectedItemTitle}</span> ({matrixRows.length} Kelas Tersedia)
        </span>
        <span class="text-[11px] text-[#747775]">
          * Pastikan jumlah pecahan komponen biaya sama dengan total tarif dasar kelas.
        </span>
      </div>

      <!-- Iterasi Seluruh Kelas Perawatan -->
      {#each matrixRows as row, rIdx (row.tariff_class_id)}
        {@const sumComponents = row.components.reduce((acc, c) => acc + (Number(c.base_amount) || 0), 0)}
        {@const isBalanced = Math.abs(sumComponents - (Number(row.total_base_price) || 0)) <= 0.01}
        {@const diff = Math.abs(sumComponents - (Number(row.total_base_price) || 0))}

        <div class="bg-white rounded-3xl border transition-all duration-200 shadow-2xs overflow-hidden {row.is_enabled ? 'border-blue-300 ring-1 ring-blue-100' : 'border-[#e1e5ea] opacity-95'}">
          <!-- Class Header Bar -->
          <div class="p-3.5 bg-[#f8fafd] border-b border-[#e1e5ea] flex flex-wrap items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <!-- Checkbox Enable / Sertakan Kelas Ini -->
              <input
                type="checkbox"
                bind:checked={row.is_enabled}
                class="w-4 h-4 rounded text-blue-600 border-[#c4c7c5] focus:ring-blue-500 cursor-pointer"
                title="Sertakan kelas ini dalam penetapan tarif"
              />

              <div class="flex items-center gap-2">
                <span class="font-mono font-bold text-xs px-2 py-0.5 rounded bg-blue-100 text-blue-800 border border-blue-200">
                  {row.tariff_class_code}
                </span>
                <span class="text-xs font-bold text-[#1f1f1f]">{row.tariff_class_name}</span>
                {#if row.existing_item_id}
                  <span class="text-[10px] font-semibold px-2 py-0.2 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200">
                    SUDAH TERKONFIGURASI
                  </span>
                {:else}
                  <span class="text-[10px] font-semibold px-2 py-0.2 rounded-full bg-amber-50 text-amber-800 border border-amber-200">
                    BELUM DISET
                  </span>
                {/if}
              </div>
            </div>

            <!-- Action Buttons on Header -->
            <div class="flex items-center gap-2">
              {#if row.is_enabled}
                <!-- Live Balancing Badge -->
                {#if isBalanced && row.total_base_price > 0}
                  <span class="text-[11px] font-semibold px-2 py-0.5 rounded-full bg-emerald-50 text-emerald-700 border border-emerald-200 flex items-center gap-1">
                    <CheckCircle2 class="w-3 h-3 text-emerald-600" />
                    <span>Komponen Seimbang ({formatIDR(row.total_base_price)})</span>
                  </span>
                {:else if row.total_base_price > 0}
                  <button
                    type="button"
                    onclick={() => handleAutoBalance(row)}
                    class="text-[11px] font-semibold px-2.5 py-1 rounded-full bg-red-50 text-red-700 border border-red-200 flex items-center gap-1.5 hover:bg-red-100 cursor-pointer transition-colors"
                    title="Klik untuk seimbangkan otomatis"
                  >
                    <AlertTriangle class="w-3 h-3 text-red-600" />
                    <span>Selisih {formatIDR(diff)} • Klik Seimbangkan</span>
                  </button>
                {/if}

                <!-- Copy To Other Classes Button -->
                <button
                  type="button"
                  onclick={() => handleCopyToAllClasses(row)}
                  class="h-7 px-2.5 rounded-lg border border-[#c4c7c5] bg-white text-[#444746] hover:bg-[#f0f4f9] text-[11px] font-medium flex items-center gap-1 cursor-pointer transition-colors shadow-2xs"
                  title="Salin nominal dan struktur komponen kelas ini ke semua kelas lain"
                >
                  <Copy class="w-3 h-3 text-blue-600" />
                  <span>Salin ke Kelas Lain</span>
                </button>
              {/if}
            </div>
          </div>

          <!-- Class Body Form -->
          <div class="p-4 flex flex-col gap-4">
            <!-- Row Pricing Inputs -->
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <!-- Total Base Price -->
              <div class="flex flex-col gap-1">
                <label class="text-[11px] font-semibold text-[#444746]">
                  Total Tarif Dasar (Reguler) <span class="text-red-500">*</span>
                </label>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-xs font-semibold text-[#747775]">Rp</span>
                  <input
                    type="number"
                    min="0"
                    step="500"
                    bind:value={row.total_base_price}
                    oninput={() => {
                      row.is_enabled = true;
                    }}
                    placeholder="0"
                    class="w-full h-9 pl-9 pr-3 rounded-xl border border-[#c4c7c5] bg-white text-xs font-semibold text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0]"
                  />
                </div>
              </div>

              <!-- Total Cito Price (Override) -->
              <div class="flex flex-col gap-1">
                <div class="flex items-center justify-between">
                  <label class="text-[11px] font-semibold text-[#444746]">
                    Tarif CITO (Override Nominal)
                  </label>
                  <button
                    type="button"
                    onclick={() => handleCalculateCito(row)}
                    class="text-[10px] text-blue-600 hover:underline font-semibold flex items-center gap-0.5 cursor-pointer"
                  >
                    <Calculator class="w-2.5 h-2.5" />
                    <span>Auto +{selectedPlan?.default_cito_percent || 25}%</span>
                  </button>
                </div>
                <div class="relative">
                  <span class="absolute left-3 top-2.5 text-xs font-semibold text-[#747775]">Rp</span>
                  <input
                    type="number"
                    min="0"
                    step="500"
                    bind:value={row.total_cito_price}
                    placeholder="Kosongkan jika formula"
                    class="w-full h-9 pl-9 pr-3 rounded-xl border border-[#c4c7c5] bg-white text-xs text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0] focus:ring-1 focus:ring-[#0b57d0]"
                  />
                </div>
              </div>

              <!-- Is Active Checkbox -->
              <div class="flex flex-col justify-end">
                <label class="h-9 px-3 rounded-xl border border-[#e1e5ea] bg-[#f8fafd] flex items-center gap-2 cursor-pointer hover:bg-[#f0f4f9] transition-colors">
                  <input
                    type="checkbox"
                    bind:checked={row.is_active}
                    class="w-3.5 h-3.5 rounded text-blue-600 border-[#c4c7c5] focus:ring-blue-500"
                  />
                  <span class="text-xs font-medium text-[#1f1f1f]">Aktif Digunakan di Kasir</span>
                </label>
              </div>
            </div>

            <!-- Pecahan Komponen Biaya Matrix Section -->
            <div class="bg-[#f8fafd] p-3 rounded-2xl border border-[#e1e5ea] flex flex-col gap-2.5">
              <div class="flex items-center justify-between text-xs">
                <div class="flex items-center gap-1.5 font-semibold text-[#1f1f1f]">
                  <Layers class="w-3.5 h-3.5 text-indigo-600" />
                  <span>Pecahan Komponen Biaya (Jasa Dokter, Sarana RS, dll.)</span>
                </div>
                <div class="flex items-center gap-2">
                  <button
                    type="button"
                    onclick={() => handleAddComponent(row)}
                    class="px-2 py-1 rounded-lg bg-blue-50 text-blue-700 hover:bg-blue-100 border border-blue-200 text-[11px] font-semibold flex items-center gap-1 cursor-pointer transition-colors"
                  >
                    <Plus class="w-3 h-3" />
                    <span>Tambah Komponen</span>
                  </button>
                </div>
              </div>

              <!-- Table Komponen Biaya -->
              <div class="overflow-x-auto">
                <table class="w-full text-left text-xs border-collapse">
                  <thead>
                    <tr class="text-[11px] font-semibold text-[#747775] border-b border-[#e1e5ea]">
                      <th class="py-1.5 px-2">Komponen Biaya</th>
                      <th class="py-1.5 px-2 w-36">Nominal Reguler (Rp)</th>
                      <th class="py-1.5 px-2 w-36">Nominal CITO (Rp)</th>
                      <th class="py-1.5 px-2 w-32">Akun COA</th>
                      <th class="py-1.5 px-2 w-10 text-center">Aksi</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-[#e1e5ea]">
                    {#each row.components as comp, cIdx}
                      {@const compDef = tariffComponents.find((tc) => tc.id === comp.component_id)}
                      <tr>
                        <!-- Select Component -->
                        <td class="py-1.5 px-2">
                          <select
                            bind:value={comp.component_id}
                            onchange={(e: any) => {
                              const found = tariffComponents.find((tc) => tc.id === e.target.value);
                              if (found) comp.coa_code = found.default_coa_code || null;
                            }}
                            class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] bg-white text-xs text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0]"
                          >
                            {#each tariffComponents as tc (tc.id)}
                              <option value={tc.id}>
                                [{tc.component_type}] {tc.name} ({tc.code})
                              </option>
                            {/each}
                          </select>
                        </td>

                        <!-- Base Amount -->
                        <td class="py-1.5 px-2">
                          <input
                            type="number"
                            min="0"
                            step="500"
                            bind:value={comp.base_amount}
                            placeholder="0"
                            class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] bg-white text-xs font-medium text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0]"
                          />
                        </td>

                        <!-- Cito Amount -->
                        <td class="py-1.5 px-2">
                          <input
                            type="number"
                            min="0"
                            step="500"
                            bind:value={comp.cito_amount}
                            placeholder="Opsional"
                            class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] bg-white text-xs text-[#1f1f1f] focus:outline-none focus:border-[#0b57d0]"
                          />
                        </td>

                        <!-- COA Code Override -->
                        <td class="py-1.5 px-2">
                          <input
                            type="text"
                            bind:value={comp.coa_code}
                            placeholder={compDef?.default_coa_code || 'Akun COA'}
                            class="w-full h-8 px-2 rounded-lg border border-[#c4c7c5] bg-white text-xs font-mono text-[#444746] focus:outline-none focus:border-[#0b57d0]"
                          />
                        </td>

                        <!-- Delete Button -->
                        <td class="py-1.5 px-2 text-center">
                          <button
                            type="button"
                            onclick={() => handleRemoveComponent(row, cIdx)}
                            class="p-1 rounded text-red-600 hover:bg-red-50 cursor-pointer transition-colors"
                            title="Hapus baris komponen ini"
                          >
                            <Trash2 class="w-3.5 h-3.5" />
                          </button>
                        </td>
                      </tr>
                    {/each}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      {/each}

      <!-- Bottom Save Action Bar -->
      <div class="sticky bottom-4 p-4 rounded-2xl bg-white/95 backdrop-blur-md border border-[#c4c7c5] shadow-lg flex items-center justify-between gap-4 z-10">
        <div class="flex items-center gap-2 text-xs text-[#444746]">
          <CheckCircle2 class="w-4 h-4 text-emerald-600 shrink-0" />
          <span>
            Siap menyimpan seluruh pengaturan tarif kelas untuk <b>{selectedItemTitle}</b> ke buku tarif <b>{selectedPlan?.name}</b>.
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
