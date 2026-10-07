<script lang="ts">
  import {
    X,
    Pill,
    Package,
    Stethoscope,
    Receipt
  } from '@lucide/svelte';
  import M3Button from '../../../../components/m3/M3Button.svelte';
  import {
    getItemById,
    createMedicationItem,
    updateMedicationItem,
    createGeneralItem,
    updateGeneralItem,
    createAssetItem,
    updateAssetItem,
    createTariffItem,
    updateTariffItem
  } from '../../../../api/catalog';
  import type {
    ItemSummaryRecord,
    ItemCategoryRecord,
    ItemProductLineRecord,
    ItemType,
    MedicationType,
    GeneralType,
    SterilizationMethod,
    DepreciationMethod,
    ChargeType
  } from '../../../../types/master/item';

  interface Props {
    isOpen: boolean;
    mode: 'CREATE' | 'EDIT';
    initialItem: ItemSummaryRecord | null;
    initialType: ItemType;
    categories: ItemCategoryRecord[];
    productLines: ItemProductLineRecord[];
    onClose: () => void;
    onSaved: () => void;
    showToast: (msg: string, isError?: boolean) => void;
  }

  let {
    isOpen,
    mode,
    initialItem,
    initialType,
    categories,
    productLines,
    onClose,
    onSaved,
    showToast
  }: Props = $props();

  let isActionLoading = $state(false);
  let formItemType = $state<ItemType>('MEDICATION');
  let editingItemId = $state<string | null>(null);

  // Form Core Fields
  let formCode = $state('');
  let formName = $state('');
  let formGenericName = $state('');
  let formCategoryId = $state('');
  let formProductLineId = $state('');
  let formBaseUnitName = $state('');
  let formIsActive = $state(true);

  // Form Subtype: Medication
  let formKfaCode = $state('');
  let formBpomNie = $state('');
  let formDosageForm = $state('');
  let formStrengthAmount = $state('');
  let formStrengthUnit = $state('mg');
  let formDefaultRoute = $state('Oral');
  let formMedicationType = $state<MedicationType>('Keras');
  let formIsHighAlert = $state(false);
  let formIsLasa = $state(false);
  let formIsFornas = $state(true);
  let formIsAntibiotic = $state(false);
  let formStorageTemp = $state('Suhu Ruang (15-25C)');

  // Form Subtype: General
  let formGeneralType = $state<GeneralType>('BMHP_MEDIS');
  let formIsSterile = $state(false);
  let formIsDisposable = $state(true);
  let formIsCssdItem = $state(false);
  let formSterilizationMethod = $state<SterilizationMethod>('STEAM_AUTOCLAVE');

  // Form Subtype: Asset
  let formBrand = $state('');
  let formModelName = $state('');
  let formIsMedicalEquipment = $state(true);
  let formExpectedLifeYears = $state<number>(5);
  let formDepreciationMethod = $state<DepreciationMethod>('STRAIGHT_LINE');
  let formMaintenanceIntervalDays = $state<number>(180);

  // Form Subtype: Tariff
  let formNameAlias = $state('');
  let formChargeType = $state<ChargeType>('TINDAKAN');
  let formTariffNotes = $state('');

  let filteredCategoriesForForm = $derived(
    categories.filter(c => c.item_type === formItemType)
  );

  $effect(() => {
    if (isOpen) {
      if (mode === 'CREATE') {
        initCreate(initialType);
      } else if (initialItem) {
        initEdit(initialItem);
      }
    }
  });

  function initCreate(type: ItemType) {
    editingItemId = null;
    formItemType = type || 'MEDICATION';

    formCode = '';
    formName = '';
    formGenericName = '';
    formCategoryId = '';
    formProductLineId = '';
    formBaseUnitName = formItemType === 'TARIFF' ? 'Kali' : '';
    formIsActive = true;

    formKfaCode = '';
    formBpomNie = '';
    formDosageForm = 'Tablet';
    formStrengthAmount = '500';
    formStrengthUnit = 'mg';
    formDefaultRoute = 'Oral';
    formMedicationType = 'Keras';
    formIsHighAlert = false;
    formIsLasa = false;
    formIsFornas = true;
    formIsAntibiotic = false;
    formStorageTemp = 'Suhu Ruang (15-25C)';

    formGeneralType = 'BMHP_MEDIS';
    formIsSterile = false;
    formIsDisposable = true;
    formIsCssdItem = false;
    formSterilizationMethod = 'STEAM_AUTOCLAVE';

    formBrand = '';
    formModelName = '';
    formIsMedicalEquipment = true;
    formExpectedLifeYears = 5;
    formDepreciationMethod = 'STRAIGHT_LINE';
    formMaintenanceIntervalDays = 180;

    formNameAlias = '';
    formChargeType = 'TINDAKAN';
    formTariffNotes = '';

    const firstCat = categories.find(c => c.item_type === formItemType);
    if (firstCat) formCategoryId = firstCat.id;

    const firstPl = productLines[0];
    if (firstPl && formItemType !== 'TARIFF') formProductLineId = firstPl.id;
  }

  async function initEdit(item: ItemSummaryRecord) {
    editingItemId = item.id;
    formItemType = item.item_type;
    isActionLoading = true;

    try {
      const full = await getItemById(item.id);
      formCode = full.code;
      formName = full.name;
      formGenericName = full.generic_name || '';
      formCategoryId = full.category?.id || '';
      formProductLineId = full.product_line?.id || '';
      formBaseUnitName = full.units?.find(u => u.is_base_unit)?.unit_name || 'Unit';
      formIsActive = full.is_active;

      if (full.medication) {
        formKfaCode = full.medication.kfa_code || '';
        formBpomNie = full.medication.bpom_nie || '';
        formDosageForm = full.medication.dosage_form || '';
        formStrengthAmount = full.medication.strength_amount || '';
        formStrengthUnit = full.medication.strength_unit || 'mg';
        formDefaultRoute = full.medication.default_route || 'Oral';
        formMedicationType = full.medication.medication_type || 'Keras';
        formIsHighAlert = full.medication.is_high_alert;
        formIsLasa = full.medication.is_lasa;
        formIsFornas = full.medication.is_fornas;
        formIsAntibiotic = full.medication.is_antibiotic;
        formStorageTemp = full.medication.storage_temperature || '';
      }

      if (full.general) {
        formGeneralType = full.general.general_type;
        formIsSterile = full.general.is_sterile;
        formIsDisposable = full.general.is_disposable;
        formIsCssdItem = full.general.is_cssd_item;
        formSterilizationMethod = full.general.sterilization_method || 'STEAM_AUTOCLAVE';
      }

      if (full.asset) {
        formBrand = full.asset.brand || '';
        formModelName = full.asset.model_name || '';
        formIsMedicalEquipment = full.asset.is_medical_equipment;
        formExpectedLifeYears = full.asset.expected_life_years || 5;
        formDepreciationMethod = full.asset.depreciation_method || 'STRAIGHT_LINE';
        formMaintenanceIntervalDays = full.asset.maintenance_interval_days || 180;
      }

      if (full.tariff) {
        formNameAlias = full.tariff.name_alias || '';
        formChargeType = full.tariff.charge_type;
        formTariffNotes = full.tariff.notes || '';
      }
    } catch (err: any) {
      showToast(err.message || 'Gagal memuat data item untuk diedit', true);
    } finally {
      isActionLoading = false;
    }
  }

  async function handleSaveForm() {
    if (!formCode.trim() || !formName.trim() || !formCategoryId) {
      showToast('Kode item, nama item, dan kategori wajib dipilih', true);
      return;
    }

    if (mode === 'CREATE' && !formBaseUnitName.trim() && formItemType !== 'TARIFF') {
      showToast('Satuan dasar (Base Unit) wajib diisi', true);
      return;
    }

    isActionLoading = true;
    try {
      if (mode === 'CREATE') {
        if (formItemType === 'MEDICATION') {
          await createMedicationItem({
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            base_unit_name: formBaseUnitName.trim(),
            kfa_code: formKfaCode.trim() || undefined,
            bpom_nie: formBpomNie.trim() || undefined,
            dosage_form: formDosageForm.trim() || undefined,
            strength_amount: formStrengthAmount.trim() || undefined,
            strength_unit: formStrengthUnit.trim() || undefined,
            default_route: formDefaultRoute.trim() || undefined,
            medication_type: formMedicationType,
            is_high_alert: formIsHighAlert,
            is_lasa: formIsLasa,
            is_fornas: formIsFornas,
            is_antibiotic: formIsAntibiotic,
            storage_temperature: formStorageTemp.trim() || undefined
          });
        } else if (formItemType === 'GENERAL') {
          await createGeneralItem({
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            base_unit_name: formBaseUnitName.trim(),
            general_type: formGeneralType,
            is_sterile: formIsSterile,
            is_disposable: formIsDisposable,
            is_cssd_item: formIsCssdItem,
            sterilization_method: formSterilizationMethod
          });
        } else if (formItemType === 'ASSET') {
          await createAssetItem({
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            base_unit_name: formBaseUnitName.trim(),
            brand: formBrand.trim() || undefined,
            model_name: formModelName.trim() || undefined,
            is_medical_equipment: formIsMedicalEquipment,
            expected_life_years: Number(formExpectedLifeYears) || undefined,
            depreciation_method: formDepreciationMethod,
            maintenance_interval_days: Number(formMaintenanceIntervalDays) || undefined
          });
        } else if (formItemType === 'TARIFF') {
          await createTariffItem({
            code: formCode.trim(),
            name: formName.trim(),
            name_alias: formNameAlias.trim() || undefined,
            category_id: formCategoryId,
            base_unit_name: formBaseUnitName.trim() || 'Kali',
            charge_type: formChargeType,
            notes: formTariffNotes.trim() || undefined
          });
        }
        showToast(`Item ${formName} berhasil didaftarkan`);
      } else if (editingItemId) {
        if (formItemType === 'MEDICATION') {
          await updateMedicationItem(editingItemId, {
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            kfa_code: formKfaCode.trim() || undefined,
            bpom_nie: formBpomNie.trim() || undefined,
            dosage_form: formDosageForm.trim() || undefined,
            strength_amount: formStrengthAmount.trim() || undefined,
            strength_unit: formStrengthUnit.trim() || undefined,
            default_route: formDefaultRoute.trim() || undefined,
            medication_type: formMedicationType,
            is_high_alert: formIsHighAlert,
            is_lasa: formIsLasa,
            is_fornas: formIsFornas,
            is_antibiotic: formIsAntibiotic,
            storage_temperature: formStorageTemp.trim() || undefined,
            is_active: formIsActive
          });
        } else if (formItemType === 'GENERAL') {
          await updateGeneralItem(editingItemId, {
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            general_type: formGeneralType,
            is_sterile: formIsSterile,
            is_disposable: formIsDisposable,
            is_cssd_item: formIsCssdItem,
            sterilization_method: formSterilizationMethod,
            is_active: formIsActive
          });
        } else if (formItemType === 'ASSET') {
          await updateAssetItem(editingItemId, {
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            brand: formBrand.trim() || undefined,
            model_name: formModelName.trim() || undefined,
            is_medical_equipment: formIsMedicalEquipment,
            expected_life_years: Number(formExpectedLifeYears) || undefined,
            depreciation_method: formDepreciationMethod,
            maintenance_interval_days: Number(formMaintenanceIntervalDays) || undefined,
            is_active: formIsActive
          });
        } else if (formItemType === 'TARIFF') {
          await updateTariffItem(editingItemId, {
            code: formCode.trim(),
            name: formName.trim(),
            name_alias: formNameAlias.trim() || undefined,
            category_id: formCategoryId,
            charge_type: formChargeType,
            notes: formTariffNotes.trim() || undefined,
            is_active: formIsActive
          });
        }
        showToast(`Item ${formName} berhasil diperbarui`);
      }

      onSaved();
    } catch (err: any) {
      showToast(err.message || 'Gagal menyimpan perubahan item', true);
    } finally {
      isActionLoading = false;
    }
  }

  function getTypeBadge(type: ItemType) {
    switch (type) {
      case 'MEDICATION':
        return { label: 'Obat', icon: Pill };
      case 'GENERAL':
        return { label: 'BMHP / Umum', icon: Package };
      case 'ASSET':
        return { label: 'Aset / Alkes', icon: Stethoscope };
      case 'TARIFF':
        return { label: 'Tarif Jasa', icon: Receipt };
    }
  }
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto animate-in fade-in duration-150">
    <div class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col my-auto max-h-[92vh]">
      <!-- Header -->
      <div class="p-5 border-b border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between">
        <div>
          <h3 class="text-sm font-bold text-[#1f1f1f]">
            {mode === 'CREATE' ? 'Tambah Master Item Baru' : 'Edit Master Item'}
          </h3>
          <p class="text-xs text-[#747775]">
            Lengkapi atribut katalog sesuai subtipe: <strong>{formItemType}</strong>
          </p>
        </div>
        <button
          type="button"
          onclick={onClose}
          class="w-8 h-8 rounded-full text-[#444746] hover:bg-[#e1e5ea] flex items-center justify-center transition-colors cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Body Form -->
      <form onsubmit={(e) => { e.preventDefault(); handleSaveForm(); }} class="p-6 overflow-y-auto flex-1 flex flex-col gap-4 text-xs">
        <!-- Subtype Selector (Hanya saat CREATE) -->
        {#if mode === 'CREATE'}
          <div class="flex flex-col gap-1.5">
            <span class="font-semibold text-[#444746]">Pilih Subtipe Item</span>
            <div class="grid grid-cols-4 gap-2">
              {#each (['MEDICATION', 'GENERAL', 'ASSET', 'TARIFF'] as const) as t}
                {@const b = getTypeBadge(t)}
                <button
                  type="button"
                  onclick={() => { formItemType = t; formCategoryId = ''; }}
                  class="p-2.5 rounded-xl border text-center transition-all cursor-pointer {formItemType === t ? 'border-[#0b57d0] bg-[#e8f0fe] font-bold text-[#0b57d0]' : 'border-[#e1e5ea] hover:bg-[#f8fafd] text-[#444746]'}"
                >
                  <b.icon class="w-4 h-4 mx-auto mb-1" />
                  <span class="text-[11px] block">{b.label}</span>
                </button>
              {/each}
            </div>
          </div>
        {/if}

        <!-- Section 1: Atribut Dasar -->
        <div class="p-4 rounded-2xl bg-[#f8fafd] border border-[#e1e5ea] flex flex-col gap-3">
          <span class="font-bold text-[#1f1f1f]">1. Identitas Universal Item</span>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="form-item-code" class="block font-semibold text-[#444746] mb-1">Kode Item / SKU *</label>
              <input
                id="form-item-code"
                type="text"
                bind:value={formCode}
                placeholder="Contoh: MED-AMOX-500"
                required
                class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] font-mono text-xs focus:border-[#0b57d0]"
              />
            </div>
            <div>
              <label for="form-item-name" class="block font-semibold text-[#444746] mb-1">Nama Resmi Item *</label>
              <input
                id="form-item-name"
                type="text"
                bind:value={formName}
                placeholder="Contoh: Amoxicillin 500mg Kaplet"
                required
                class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label for="form-item-generic" class="block font-semibold text-[#444746] mb-1">Nama Generik / Zat Aktif</label>
              <input
                id="form-item-generic"
                type="text"
                bind:value={formGenericName}
                placeholder="Contoh: Amoxicillin Trihydrate"
                class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
              />
            </div>
            <div>
              <label for="form-item-cat" class="block font-semibold text-[#444746] mb-1">Kategori Item *</label>
              <select
                id="form-item-cat"
                bind:value={formCategoryId}
                required
                class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
              >
                <option value="">Pilih Kategori</option>
                {#each filteredCategoriesForForm as cat (cat.id)}
                  <option value={cat.id}>{cat.name} ({cat.code})</option>
                {/each}
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            {#if formItemType !== 'TARIFF'}
              <div>
                <label for="form-item-pl" class="block font-semibold text-[#444746] mb-1">Lini Produk Persediaan</label>
                <select
                  id="form-item-pl"
                  bind:value={formProductLineId}
                  class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
                >
                  <option value="">Pilih Lini Produk</option>
                  {#each productLines as pl (pl.id)}
                    <option value={pl.id}>{pl.name}</option>
                  {/each}
                </select>
              </div>
            {/if}

            {#if mode === 'CREATE'}
              <div>
                <label for="form-item-base-uom" class="block font-semibold text-[#444746] mb-1">Satuan Dasar (Base Unit) *</label>
                <input
                  id="form-item-base-uom"
                  type="text"
                  bind:value={formBaseUnitName}
                  placeholder={formItemType === 'TARIFF' ? 'Kali / Tindakan' : 'Tablet / Pcs / Botol'}
                  required
                  class="w-full h-9 px-3 rounded-xl bg-white border border-[#e1e5ea] text-xs focus:border-[#0b57d0]"
                />
              </div>
            {/if}
          </div>

          {#if mode === 'EDIT'}
            <div class="flex items-center gap-2 mt-1">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsActive} class="rounded text-[#0b57d0]" />
                <span class="font-semibold text-[#1f1f1f]">Item Aktif Digunakan di Rumah Sakit</span>
              </label>
            </div>
          {/if}
        </div>

        <!-- Section 2: Subtype Spesifik Fields -->
        {#if formItemType === 'MEDICATION'}
          <div class="p-4 rounded-2xl bg-emerald-50/40 border border-emerald-200 flex flex-col gap-3">
            <span class="font-bold text-emerald-900 flex items-center gap-1.5">
              <Pill class="w-4 h-4 text-emerald-600" />
              <span>2. Atribut Farmasi & Regulasi Obat</span>
            </span>

            <div class="grid grid-cols-3 gap-3">
              <div>
                <label for="form-med-kfa" class="block font-semibold text-[#444746] mb-1">Kode KFA SATUSEHAT</label>
                <input
                  id="form-med-kfa"
                  type="text"
                  bind:value={formKfaCode}
                  placeholder="93000123"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs font-mono"
                />
              </div>
              <div>
                <label for="form-med-bpom" class="block font-semibold text-[#444746] mb-1">Nomor Izin BPOM (NIE)</label>
                <input
                  id="form-med-bpom"
                  type="text"
                  bind:value={formBpomNie}
                  placeholder="DKL1234567809A1"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs font-mono"
                />
              </div>
              <div>
                <label for="form-med-dosage" class="block font-semibold text-[#444746] mb-1">Bentuk Sediaan</label>
                <input
                  id="form-med-dosage"
                  type="text"
                  bind:value={formDosageForm}
                  placeholder="Kaplet / Sirup / Injeksi"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs"
                />
              </div>
            </div>

            <div class="grid grid-cols-3 gap-3">
              <div>
                <label for="form-med-strength" class="block font-semibold text-[#444746] mb-1">Jumlah Kekuatan</label>
                <input
                  id="form-med-strength"
                  type="text"
                  bind:value={formStrengthAmount}
                  placeholder="500"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs"
                />
              </div>
              <div>
                <label for="form-med-str-unit" class="block font-semibold text-[#444746] mb-1">Satuan Kekuatan</label>
                <input
                  id="form-med-str-unit"
                  type="text"
                  bind:value={formStrengthUnit}
                  placeholder="mg / mcg / ml"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs"
                />
              </div>
              <div>
                <label for="form-med-type" class="block font-semibold text-[#444746] mb-1">Golongan Obat</label>
                <select
                  id="form-med-type"
                  bind:value={formMedicationType}
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs"
                >
                  <option value="Bebas">Bebas</option>
                  <option value="Keras">Keras</option>
                  <option value="Narkotika">Narkotika</option>
                  <option value="Psikotropika">Psikotropika</option>
                  <option value="Prekursor">Prekursor</option>
                  <option value="Lainnya">Lainnya</option>
                </select>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="form-med-route" class="block font-semibold text-[#444746] mb-1">Rute Pemberian Default</label>
                <input
                  id="form-med-route"
                  type="text"
                  bind:value={formDefaultRoute}
                  placeholder="Oral / IV / IM / Topikal"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs"
                />
              </div>
              <div>
                <label for="form-med-storage" class="block font-semibold text-[#444746] mb-1">Suhu Penyimpanan</label>
                <input
                  id="form-med-storage"
                  type="text"
                  bind:value={formStorageTemp}
                  placeholder="Suhu Ruang (15-25C) / Kulkas (2-8C)"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-emerald-200 text-xs"
                />
              </div>
            </div>

            <!-- Flags -->
            <div class="flex items-center gap-4 flex-wrap pt-1">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsHighAlert} class="rounded text-rose-600" />
                <span class="font-medium text-rose-800">High-Alert (HAM)</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsLasa} class="rounded text-amber-600" />
                <span class="font-medium text-amber-800">LASA / NORUM</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsFornas} class="rounded text-blue-600" />
                <span class="font-medium text-blue-800">Fornas</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsAntibiotic} class="rounded text-purple-600" />
                <span class="font-medium text-purple-800">Antibiotik</span>
              </label>
            </div>
          </div>

        {:else if formItemType === 'GENERAL'}
          <div class="p-4 rounded-2xl bg-amber-50/40 border border-amber-200 flex flex-col gap-3">
            <span class="font-bold text-amber-900 flex items-center gap-1.5">
              <Package class="w-4 h-4 text-amber-600" />
              <span>2. Atribut BMHP & Siklus CSSD</span>
            </span>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="form-gen-type" class="block font-semibold text-[#444746] mb-1">Subtipe Barang Umum *</label>
                <select
                  id="form-gen-type"
                  bind:value={formGeneralType}
                  class="w-full h-9 px-3 rounded-xl bg-white border border-amber-200 text-xs"
                >
                  <option value="BMHP_MEDIS">BMHP Medis</option>
                  <option value="INSTRUMEN_MEDIS">Instrumen Medis</option>
                  <option value="ATK">ATK</option>
                  <option value="LINEN">Linen</option>
                  <option value="KEBERSIHAN">Kebersihan</option>
                  <option value="DAPUR">Dapur</option>
                  <option value="LAINNYA">Lainnya</option>
                </select>
              </div>
              <div>
                <label for="form-gen-sterilization" class="block font-semibold text-[#444746] mb-1">Metode Sterilisasi</label>
                <select
                  id="form-gen-sterilization"
                  bind:value={formSterilizationMethod}
                  class="w-full h-9 px-3 rounded-xl bg-white border border-amber-200 text-xs"
                >
                  <option value="STEAM_AUTOCLAVE">Steam Autoclave</option>
                  <option value="EO_GAS">EO Gas</option>
                  <option value="PLASMA">Plasma</option>
                  <option value="DRY_HEAT">Dry Heat</option>
                </select>
              </div>
            </div>

            <div class="flex items-center gap-4 flex-wrap pt-1">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsSterile} class="rounded text-emerald-600" />
                <span class="font-medium text-[#1f1f1f]">Produk Steril</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsDisposable} class="rounded text-[#0b57d0]" />
                <span class="font-medium text-[#1f1f1f]">Sekali Pakai (Disposable)</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsCssdItem} class="rounded text-purple-600" />
                <span class="font-medium text-[#1f1f1f]">Dikelola CSSD</span>
              </label>
            </div>
          </div>

        {:else if formItemType === 'ASSET'}
          <div class="p-4 rounded-2xl bg-indigo-50/40 border border-indigo-200 flex flex-col gap-3">
            <span class="font-bold text-indigo-900 flex items-center gap-1.5">
              <Stethoscope class="w-4 h-4 text-indigo-600" />
              <span>2. Atribut Barang Modal, Depresiasi & Kalibrasi</span>
            </span>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="form-ast-brand" class="block font-semibold text-[#444746] mb-1">Merk / Pabrikan</label>
                <input
                  id="form-ast-brand"
                  type="text"
                  bind:value={formBrand}
                  placeholder="Misal: GE Healthcare, Philips"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-indigo-200 text-xs"
                />
              </div>
              <div>
                <label for="form-ast-model" class="block font-semibold text-[#444746] mb-1">Tipe / Model Mesin</label>
                <input
                  id="form-ast-model"
                  type="text"
                  bind:value={formModelName}
                  placeholder="Misal: Voluson E8"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-indigo-200 text-xs"
                />
              </div>
            </div>

            <div class="grid grid-cols-3 gap-3">
              <div>
                <label for="form-ast-life" class="block font-semibold text-[#444746] mb-1">Estimasi Umur (Tahun)</label>
                <input
                  id="form-ast-life"
                  type="number"
                  bind:value={formExpectedLifeYears}
                  min="1"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-indigo-200 text-xs font-mono"
                />
              </div>
              <div>
                <label for="form-ast-deprec" class="block font-semibold text-[#444746] mb-1">Metode Depresiasi</label>
                <select
                  id="form-ast-deprec"
                  bind:value={formDepreciationMethod}
                  class="w-full h-9 px-3 rounded-xl bg-white border border-indigo-200 text-xs"
                >
                  <option value="STRAIGHT_LINE">Straight Line</option>
                  <option value="DOUBLE_DECLINING">Double Declining</option>
                  <option value="SUM_OF_YEARS">Sum of Years</option>
                </select>
              </div>
              <div>
                <label for="form-ast-interval" class="block font-semibold text-[#444746] mb-1">Interval Servis (Hari)</label>
                <input
                  id="form-ast-interval"
                  type="number"
                  bind:value={formMaintenanceIntervalDays}
                  min="1"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-indigo-200 text-xs font-mono"
                />
              </div>
            </div>

            <div class="pt-1">
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="checkbox" bind:checked={formIsMedicalEquipment} class="rounded text-indigo-600" />
                <span class="font-medium text-indigo-900">Alat Kesehatan Medis / Elektromedik</span>
              </label>
            </div>
          </div>

        {:else if formItemType === 'TARIFF'}
          <div class="p-4 rounded-2xl bg-blue-50/40 border border-blue-200 flex flex-col gap-3">
            <span class="font-bold text-blue-900 flex items-center gap-1.5">
              <Receipt class="w-4 h-4 text-blue-600" />
              <span>2. Atribut Tarif Layanan & Jasa Medis</span>
            </span>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label for="form-tar-alias" class="block font-semibold text-[#444746] mb-1">Nama Alias / English Name</label>
                <input
                  id="form-tar-alias"
                  type="text"
                  bind:value={formNameAlias}
                  placeholder="Misal: Outpatient Specialist Consultation"
                  class="w-full h-9 px-3 rounded-xl bg-white border border-blue-200 text-xs"
                />
              </div>
              <div>
                <label for="form-tar-charge" class="block font-semibold text-[#444746] mb-1">Klasifikasi Beban Tagihan *</label>
                <select
                  id="form-tar-charge"
                  bind:value={formChargeType}
                  class="w-full h-9 px-3 rounded-xl bg-white border border-blue-200 text-xs"
                >
                  <option value="TINDAKAN">Tindakan</option>
                  <option value="ADMINISTRASI">Administrasi</option>
                  <option value="AKOMODASI">Akomodasi</option>
                  <option value="PENUNJANG">Penunjang</option>
                  <option value="LAINNYA">Lainnya</option>
                </select>
              </div>
            </div>

            <div>
              <label for="form-tar-notes" class="block font-semibold text-[#444746] mb-1">Catatan Operasional / Deskripsi Penagihan</label>
              <textarea
                id="form-tar-notes"
                bind:value={formTariffNotes}
                rows="2"
                placeholder="Deskripsi layanan atau syarat penagihan..."
                class="w-full p-2.5 rounded-xl bg-white border border-blue-200 text-xs focus:border-[#0b57d0]"
              ></textarea>
            </div>
          </div>
        {/if}

        <!-- Footer Buttons -->
        <div class="flex items-center justify-end gap-2 pt-3 border-t border-[#e1e5ea]">
          <button
            type="button"
            onclick={onClose}
            class="px-4 py-2 rounded-xl border border-[#e1e5ea] bg-white text-xs font-semibold text-[#444746] hover:bg-[#f0f4f9] cursor-pointer"
          >
            Batal
          </button>
          <M3Button variant="filled" type="submit" disabled={isActionLoading}>
            <span>{mode === 'CREATE' ? 'Daftarkan Item' : 'Simpan Perubahan'}</span>
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
