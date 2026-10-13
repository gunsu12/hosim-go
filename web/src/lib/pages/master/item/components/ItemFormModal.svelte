<script lang="ts">
  import { Pill, Package, Stethoscope, Receipt } from "@lucide/svelte";
  import M3Modal from "../../../../components/m3/M3Modal.svelte";
  import M3Button from "../../../../components/m3/M3Button.svelte";
  import M3Select from "../../../../components/m3/M3Select.svelte";
  import M3ToggleCard from "../../../../components/m3/M3ToggleCard.svelte";
  import ItemFormMedicationSection from "./ItemFormMedicationSection.svelte";
  import ItemFormGeneralSection from "./ItemFormGeneralSection.svelte";
  import ItemFormAssetSection from "./ItemFormAssetSection.svelte";
  import ItemFormTariffSection from "./ItemFormTariffSection.svelte";
  import {
    getItemById,
    createMedicationItem,
    updateMedicationItem,
    createGeneralItem,
    updateGeneralItem,
    createAssetItem,
    updateAssetItem,
    createTariffItem,
    updateTariffItem,
  } from "../../../../api/catalog";
  import type {
    ItemSummaryRecord,
    ItemCategoryRecord,
    ItemProductLineRecord,
    ItemType,
    MedicationType,
    GeneralType,
    SterilizationMethod,
    DepreciationMethod,
    ChargeType,
  } from "../../../../types/master/item";

  interface Props {
    isOpen: boolean;
    mode: "CREATE" | "EDIT";
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
    showToast,
  }: Props = $props();

  let isActionLoading = $state(false);
  let formItemType = $state<ItemType>("MEDICATION");
  let editingItemId = $state<string | null>(null);

  // Form Core Fields
  let formCode = $state("");
  let formName = $state("");
  let formGenericName = $state("");
  let formCategoryId = $state("");
  let formProductLineId = $state("");
  let formBaseUnitName = $state("");
  let formIsActive = $state(true);

  // Form Subtype: Medication
  let formKfaCode = $state("");
  let formBpomNie = $state("");
  let formDosageForm = $state("");
  let formStrengthAmount = $state("");
  let formStrengthUnit = $state("mg");
  let formDefaultRoute = $state("Oral");
  let formMedicationType = $state<MedicationType>("Keras");
  let formIsHighAlert = $state(false);
  let formIsLasa = $state(false);
  let formIsFornas = $state(true);
  let formIsAntibiotic = $state(false);
  let formStorageTemp = $state("Suhu Ruang (15-25C)");

  // Form Subtype: General
  let formGeneralType = $state<GeneralType>("BMHP_MEDIS");
  let formIsSterile = $state(false);
  let formIsDisposable = $state(true);
  let formIsCssdItem = $state(false);
  let formSterilizationMethod = $state<SterilizationMethod>("STEAM_AUTOCLAVE");

  // Form Subtype: Asset
  let formBrand = $state("");
  let formModelName = $state("");
  let formIsMedicalEquipment = $state(true);
  let formExpectedLifeYears = $state<number>(5);
  let formDepreciationMethod = $state<DepreciationMethod>("STRAIGHT_LINE");
  let formMaintenanceIntervalDays = $state<number>(180);

  // Form Subtype: Tariff
  let formNameAlias = $state("");
  let formChargeType = $state<ChargeType>("TINDAKAN");
  let formTariffNotes = $state("");

  let filteredCategoriesForForm = $derived(
    categories.filter((c) => c.item_type === formItemType),
  );

  let categoryOptions = $derived(
    filteredCategoriesForForm.map((c) => ({
      value: c.id,
      label: c.name,
    })),
  );

  let productLineOptions = $derived([
    { value: "", label: "-- Non-Persediaan / Tidak Terikat --" },
    ...productLines.map((pl) => ({
      value: pl.id,
      label: pl.name,
    })),
  ]);

  $effect(() => {
    if (isOpen) {
      if (mode === "CREATE") {
        initCreate(initialType);
      } else if (initialItem) {
        initEdit(initialItem);
      }
    }
  });

  function initCreate(type: ItemType) {
    editingItemId = null;
    formItemType = type || "MEDICATION";

    formCode = "";
    formName = "";
    formGenericName = "";
    formCategoryId = "";
    formProductLineId = "";
    formBaseUnitName = formItemType === "TARIFF" ? "Kali" : "";
    formIsActive = true;

    formKfaCode = "";
    formBpomNie = "";
    formDosageForm = "Tablet";
    formStrengthAmount = "500";
    formStrengthUnit = "mg";
    formDefaultRoute = "Oral";
    formMedicationType = "Keras";
    formIsHighAlert = false;
    formIsLasa = false;
    formIsFornas = true;
    formIsAntibiotic = false;
    formStorageTemp = "Suhu Ruang (15-25C)";

    formGeneralType = "BMHP_MEDIS";
    formIsSterile = false;
    formIsDisposable = true;
    formIsCssdItem = false;
    formSterilizationMethod = "STEAM_AUTOCLAVE";

    formBrand = "";
    formModelName = "";
    formIsMedicalEquipment = true;
    formExpectedLifeYears = 5;
    formDepreciationMethod = "STRAIGHT_LINE";
    formMaintenanceIntervalDays = 180;

    formNameAlias = "";
    formChargeType = "TINDAKAN";
    formTariffNotes = "";

    const firstCat = categories.find((c) => c.item_type === formItemType);
    if (firstCat) formCategoryId = firstCat.id;
  }

  async function initEdit(item: ItemSummaryRecord) {
    editingItemId = item.id;
    formItemType = item.item_type;

    formCode = item.code;
    formName = item.name;
    formGenericName = item.generic_name || "";
    formCategoryId = item.category?.id || item.category_id || "";
    formProductLineId = item.product_line?.id || item.product_line_id || "";
    formBaseUnitName = item.base_unit?.unit_name || item.base_unit_name || "";
    formIsActive = item.is_active;

    try {
      const d = await getItemById(item.id);

      if (d.medication) {
        formKfaCode = d.medication.kfa_code || "";
        formBpomNie = d.medication.bpom_nie || "";
        formDosageForm = d.medication.dosage_form || "Tablet";
        formStrengthAmount = d.medication.strength_amount || "";
        formStrengthUnit = d.medication.strength_unit || "mg";
        formDefaultRoute = d.medication.default_route || "Oral";
        formMedicationType = d.medication.medication_type || "Keras";
        formIsHighAlert = d.medication.is_high_alert || false;
        formIsLasa = d.medication.is_lasa || false;
        formIsFornas = d.medication.is_fornas ?? true;
        formIsAntibiotic = d.medication.is_antibiotic || false;
        formStorageTemp =
          d.medication.storage_temperature || "Suhu Ruang (15-25C)";
      }

      if (d.general) {
        formGeneralType = d.general.general_type || "BMHP_MEDIS";
        formIsSterile = d.general.is_sterile || false;
        formIsDisposable = d.general.is_disposable ?? true;
        formIsCssdItem = d.general.is_cssd_item || false;
        formSterilizationMethod =
          d.general.sterilization_method || "STEAM_AUTOCLAVE";
      }

      if (d.asset) {
        formBrand = d.asset.brand || "";
        formModelName = d.asset.model_name || "";
        formIsMedicalEquipment = d.asset.is_medical_equipment ?? true;
        formExpectedLifeYears = d.asset.expected_life_years || 5;
        formDepreciationMethod =
          d.asset.depreciation_method || "STRAIGHT_LINE";
        formMaintenanceIntervalDays = d.asset.maintenance_interval_days || 180;
      }

      if (d.tariff) {
        formNameAlias = d.tariff.name_alias || "";
        formChargeType = d.tariff.charge_type || "TINDAKAN";
        formTariffNotes = d.tariff.notes || "";
      }
    } catch (err: any) {
      showToast(err.message || "Gagal memuat rincian item untuk edit", true);
    }
  }

  async function handleSaveForm() {
    if (!formCode.trim() || !formName.trim() || !formCategoryId) {
      showToast("Kode SKU, Nama Item, dan Kategori wajib diisi!", true);
      return;
    }
    if (mode === "CREATE" && !formBaseUnitName.trim()) {
      showToast("Satuan Dasar (Base Unit) wajib ditentukan!", true);
      return;
    }

    isActionLoading = true;
    try {
      if (mode === "CREATE") {
        if (formItemType === "MEDICATION") {
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
            storage_temperature: formStorageTemp.trim() || undefined,
          });
        } else if (formItemType === "GENERAL") {
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
            sterilization_method: formSterilizationMethod,
          });
        } else if (formItemType === "ASSET") {
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
            expected_life_years: Number(formExpectedLifeYears) || 5,
            depreciation_method: formDepreciationMethod,
            maintenance_interval_days:
              Number(formMaintenanceIntervalDays) || undefined,
          });
        } else if (formItemType === "TARIFF") {
          await createTariffItem({
            code: formCode.trim(),
            name: formName.trim(),
            name_alias: formNameAlias.trim() || undefined,
            category_id: formCategoryId,
            charge_type: formChargeType,
            notes: formTariffNotes.trim() || undefined,
          });
        }
        showToast(`Item ${formName} berhasil didaftarkan`);
      } else if (editingItemId) {
        if (formItemType === "MEDICATION") {
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
            is_active: formIsActive,
          });
        } else if (formItemType === "GENERAL") {
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
            is_active: formIsActive,
          });
        } else if (formItemType === "ASSET") {
          await updateAssetItem(editingItemId, {
            code: formCode.trim(),
            name: formName.trim(),
            generic_name: formGenericName.trim() || undefined,
            category_id: formCategoryId,
            product_line_id: formProductLineId || undefined,
            brand: formBrand.trim() || undefined,
            model_name: formModelName.trim() || undefined,
            is_medical_equipment: formIsMedicalEquipment,
            expected_life_years: Number(formExpectedLifeYears) || 5,
            depreciation_method: formDepreciationMethod,
            maintenance_interval_days:
              Number(formMaintenanceIntervalDays) || undefined,
            is_active: formIsActive,
          });
        } else if (formItemType === "TARIFF") {
          await updateTariffItem(editingItemId, {
            code: formCode.trim(),
            name: formName.trim(),
            name_alias: formNameAlias.trim() || undefined,
            category_id: formCategoryId,
            charge_type: formChargeType,
            notes: formTariffNotes.trim() || undefined,
            is_active: formIsActive,
          });
        }
        showToast(`Item ${formName} berhasil diperbarui`);
      }

      onSaved();
    } catch (err: any) {
      showToast(err.message || "Gagal menyimpan perubahan item", true);
    } finally {
      isActionLoading = false;
    }
  }

  function getTypeBadge(type: ItemType) {
    switch (type) {
      case "MEDICATION":
        return { label: "Obat", icon: Pill };
      case "GENERAL":
        return { label: "BMHP / Umum", icon: Package };
      case "ASSET":
        return { label: "Aset / Alkes", icon: Stethoscope };
      case "TARIFF":
        return { label: "Tarif Jasa", icon: Receipt };
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === "CREATE" ? "Tambah Master Item Baru" : "Edit Master Item"}
  size="2xl"
  {onClose}
>
  <form
    onsubmit={(e) => {
      e.preventDefault();
      handleSaveForm();
    }}
    class="flex flex-col gap-4 text-xs"
  >
    <!-- Subtype Selector (Hanya saat CREATE) -->
    {#if mode === "CREATE"}
      <div class="flex flex-col gap-1.5">
        <span class="font-semibold text-[#444746] dark:text-[#c4c7c5]">
          Pilih Subtipe Item
        </span>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
          {#each ["MEDICATION", "GENERAL", "ASSET", "TARIFF"] as const as t}
            {@const b = getTypeBadge(t)}
            <button
              type="button"
              onclick={() => {
                formItemType = t;
                formCategoryId = "";
              }}
              class="p-2.5 rounded-xl border text-center transition-all cursor-pointer {formItemType ===
              t
                ? 'border-[#0b57d0] bg-[#e8f0fe] dark:bg-[#0b57d0]/20 font-bold text-[#0b57d0] dark:text-blue-400'
                : 'border-[#e1e5ea] dark:border-[#444746] hover:bg-[#f8fafd] dark:hover:bg-[#282a2c] text-[#444746] dark:text-[#c4c7c5]'}"
            >
              <b.icon class="w-4 h-4 mx-auto mb-1" />
              <span class="text-[11px] block">{b.label}</span>
            </button>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Section 1: Atribut Dasar -->
    <div
      class="p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#2b2b2b] border border-[#e1e5ea] dark:border-[#444746] flex flex-col gap-3.5"
    >
      <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
        1. Identitas Universal Item
      </span>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div>
          <label
            for="form-item-code"
            class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
          >
            Kode Item / SKU *
          </label>
          <input
            id="form-item-code"
            type="text"
            bind:value={formCode}
            placeholder="Contoh: MED-AMOX-500"
            required
            class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] font-mono text-xs focus:border-[#0b57d0] focus:outline-none"
          />
        </div>
        <div>
          <label
            for="form-item-name"
            class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
          >
            Nama Resmi Item *
          </label>
          <input
            id="form-item-name"
            type="text"
            bind:value={formName}
            placeholder="Contoh: Amoxicillin 500mg Kapsul"
            required
            class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
          />
        </div>
      </div>

      <div>
        <label
          for="form-item-generic"
          class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
        >
          Nama Generik / Zat Aktif / Deskripsi Umum
        </label>
        <input
          id="form-item-generic"
          type="text"
          bind:value={formGenericName}
          placeholder="Contoh: Amoxicillin Trihydrate"
          class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
        />
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 items-end">
        <M3Select
          label="Kategori Item *"
          bind:value={formCategoryId}
          options={categoryOptions}
          placeholder="Pilih Kategori"
        />

        {#if formItemType !== "TARIFF"}
          <M3Select
            label="Lini Produk Persediaan"
            bind:value={formProductLineId}
            options={productLineOptions}
            placeholder="Pilih Lini Produk"
          />
        {/if}

        {#if mode === "CREATE"}
          <div>
            <label
              for="form-item-base-uom"
              class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1"
            >
              Satuan Dasar (Base Unit) *
            </label>
            <input
              id="form-item-base-uom"
              type="text"
              bind:value={formBaseUnitName}
              placeholder={formItemType === "TARIFF"
                ? "Kali / Tindakan"
                : "Tablet / Pcs / Botol"}
              required
              class="w-full h-9 px-3 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#444746] text-[#1f1f1f] dark:text-[#e3e3e3] text-xs focus:border-[#0b57d0] focus:outline-none"
            />
          </div>
        {/if}
      </div>

      {#if mode === "EDIT"}
        <div class="pt-1">
          <M3ToggleCard
            title="Status Item Aktif"
            description="Item aktif dapat digunakan di resep, stok, atau billing kasir"
            bind:checked={formIsActive}
          />
        </div>
      {/if}
    </div>

    <!-- Section 2: Subtype Spesifik Fields (Decomposed Components) -->
    {#if formItemType === "MEDICATION"}
      <ItemFormMedicationSection
        bind:kfaCode={formKfaCode}
        bind:bpomNie={formBpomNie}
        bind:dosageForm={formDosageForm}
        bind:strengthAmount={formStrengthAmount}
        bind:strengthUnit={formStrengthUnit}
        bind:defaultRoute={formDefaultRoute}
        bind:medicationType={formMedicationType}
        bind:storageTemp={formStorageTemp}
        bind:isHighAlert={formIsHighAlert}
        bind:isLasa={formIsLasa}
        bind:isFornas={formIsFornas}
        bind:isAntibiotic={formIsAntibiotic}
      />
    {:else if formItemType === "GENERAL"}
      <ItemFormGeneralSection
        bind:generalType={formGeneralType}
        bind:sterilizationMethod={formSterilizationMethod}
        bind:isSterile={formIsSterile}
        bind:isDisposable={formIsDisposable}
        bind:isCssdItem={formIsCssdItem}
      />
    {:else if formItemType === "ASSET"}
      <ItemFormAssetSection
        bind:brand={formBrand}
        bind:modelName={formModelName}
        bind:expectedLifeYears={formExpectedLifeYears}
        bind:depreciationMethod={formDepreciationMethod}
        bind:maintenanceIntervalDays={formMaintenanceIntervalDays}
        bind:isMedicalEquipment={formIsMedicalEquipment}
      />
    {:else if formItemType === "TARIFF"}
      <ItemFormTariffSection
        bind:nameAlias={formNameAlias}
        bind:chargeType={formChargeType}
        bind:tariffNotes={formTariffNotes}
      />
    {/if}

    <!-- Footer Action Buttons -->
    <div
      class="flex items-center justify-end gap-2 pt-3 border-t border-[#e1e5ea] dark:border-[#444746]"
    >
      <M3Button variant="text" type="button" onclick={onClose}>
        Batal
      </M3Button>
      <M3Button variant="filled" type="submit" disabled={isActionLoading}>
        <span>{mode === "CREATE" ? "Daftarkan Item" : "Simpan Perubahan"}</span>
      </M3Button>
    </div>
  </form>
</M3Modal>
