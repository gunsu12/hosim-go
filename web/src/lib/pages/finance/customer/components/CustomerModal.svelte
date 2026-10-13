<script lang="ts">
  import {
    ShieldCheck,
    CheckCircle2,
  } from "@lucide/svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
  import M3Modal from "$lib/components/m3/M3Modal.svelte";
  import M3Select from "$lib/components/m3/M3Select.svelte";
  import M3ToggleCard from "$lib/components/m3/M3ToggleCard.svelte";
  import { createCustomer, updateCustomer } from "$lib/api/finance";
  import type {
    CustomerRecord,
    CustomerTypeRecord,
    CreateCustomerDTO,
    UpdateCustomerDTO,
  } from "$lib/types/finance";

  interface Props {
    isOpen: boolean;
    mode: "CREATE" | "EDIT";
    customer?: CustomerRecord | null;
    customerTypes?: CustomerTypeRecord[];
    onClose: () => void;
    onSuccess: (saved: CustomerRecord) => void;
  }

  let {
    isOpen = false,
    mode = "CREATE",
    customer = null,
    customerTypes = [],
    onClose = () => {},
    onSuccess = () => {},
  }: Props = $props();

  let code = $state("");
  let name = $state("");
  let customerTypeId = $state("");
  let contactPerson = $state("");
  let phone = $state("");
  let email = $state("");
  let website = $state("");
  let address = $state("");
  let description = $state("");
  let requireCard = $state(true);
  let isActive = $state(true);

  let isSubmitting = $state(false);
  let errorMessage = $state<string | null>(null);

  $effect(() => {
    if (isOpen) {
      errorMessage = null;
      if (mode === "EDIT" && customer) {
        code = customer.code || "";
        name = customer.name || "";
        customerTypeId = customer.customer_type_id || "";
        contactPerson = customer.contact_person || "";
        phone = customer.phone || "";
        email = customer.email || "";
        website = customer.website || "";
        address = customer.address || "";
        description = customer.description || "";
        requireCard = customer.require_card ?? true;
        isActive = customer.is_active;
      } else {
        code = "";
        name = "";
        customerTypeId = customerTypes.length > 0 ? customerTypes[0].id : "";
        contactPerson = "";
        phone = "";
        email = "";
        website = "";
        address = "";
        description = "";
        requireCard = true;
        isActive = true;
      }
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    errorMessage = null;

    if (!code.trim()) {
      errorMessage = "Kode penjamin wajib diisi";
      return;
    }
    if (!name.trim()) {
      errorMessage = "Nama penjamin wajib diisi";
      return;
    }

    isSubmitting = true;

    try {
      if (mode === "CREATE") {
        const dto: CreateCustomerDTO = {
          code: code.trim().toUpperCase(),
          name: name.trim(),
          customer_type_id: customerTypeId ? customerTypeId : null,
          contact_person: contactPerson.trim() || null,
          phone: phone.trim() || null,
          email: email.trim() || null,
          website: website.trim() || null,
          address: address.trim() || null,
          description: description.trim() || null,
          require_card: requireCard,
          is_active: isActive,
        };
        const res = await createCustomer(dto);
        onSuccess(res);
      } else if (mode === "EDIT" && customer) {
        const dto: UpdateCustomerDTO = {
          code: code.trim().toUpperCase(),
          name: name.trim(),
          customer_type_id: customerTypeId ? customerTypeId : null,
          contact_person: contactPerson.trim() || null,
          phone: phone.trim() || null,
          email: email.trim() || null,
          website: website.trim() || null,
          address: address.trim() || null,
          description: description.trim() || null,
          require_card: requireCard,
          is_active: isActive,
        };
        const res = await updateCustomer(customer.id, dto);
        onSuccess(res);
      }
    } catch (err: any) {
      console.error("Error saving customer:", err);
      errorMessage =
        err?.message || "Terjadi kesalahan saat menyimpan data penjamin";
    } finally {
      isSubmitting = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === "CREATE"
    ? "Tambah Debitur / Penjamin Baru"
    : "Ubah Data Debitur / Penjamin"}
  subtitle={mode === "CREATE"
    ? "Konfigurasi penjamin pembiayaan (BPJS, Asuransi Swasta, Korporat)"
    : `Kode: ${customer?.code || "-"}`}
  icon={ShieldCheck}
  size="2xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if errorMessage}
    <div
      class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5"
    >
      <span class="font-bold">Gagal:</span>
      <span>{errorMessage}</span>
    </div>
  {/if}

  <!-- Section: Identitas Utama -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label
        for="cust-code"
        class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
      >
        Kode Penjamin <span class="text-rose-500">*</span>
      </label>
      <input
        id="cust-code"
        type="text"
        bind:value={code}
        placeholder="Contoh: BPJS-JKN, ASR-PRUD"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
      <p class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1">
        Harus unik di database sistem RS.
      </p>
    </div>

    <div>
      <M3Select
        id="cust-type"
        label="Kategori Penjamin"
        bind:value={customerTypeId}
        placeholder="-- Pilih Kategori / Tipe --"
        options={customerTypes.map((ct) => ({
          value: ct.id,
          label: ct.name,
          badge: ct.code,
        }))}
      />
    </div>
  </div>

  <div>
    <label
      for="cust-name"
      class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
    >
      Nama Instansi / Perusahaan Penjamin <span class="text-rose-500">*</span>
    </label>
    <input
      id="cust-name"
      type="text"
      bind:value={name}
      placeholder="Contoh: BPJS Kesehatan KC Jakarta, PT Prudential Life Assurance"
      required
      class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
    />
  </div>

  <!-- Section: Kontak & PIC -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
    <h4
      class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider mb-3"
    >
      Kontak & Perwakilan (PIC)
    </h4>
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div>
        <label
          for="cust-contact-person"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Nama PIC / Koordinator Klaim
        </label>
        <input
          id="cust-contact-person"
          type="text"
          bind:value={contactPerson}
          placeholder="Contoh: Budi Santoso (Head of Claims)"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <label
          for="cust-phone"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Telepon / Hotline Call Center
        </label>
        <input
          id="cust-phone"
          type="text"
          bind:value={phone}
          placeholder="Contoh: 1500085 atau 021-29958888"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all font-mono"
        />
      </div>

      <div>
        <label
          for="cust-email"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Email Korespondensi / Klaim
        </label>
        <input
          id="cust-email"
          type="email"
          bind:value={email}
          placeholder="Contoh: claim.hospital@insurance.co.id"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <label
          for="cust-website"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Situs Resmi / Portal Klaim
        </label>
        <input
          id="cust-website"
          type="text"
          bind:value={website}
          placeholder="Contoh: https://portal.asuransi.co.id"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>
    </div>
  </div>

  <!-- Section: Alamat & Deskripsi -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738]">
    <h4
      class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider mb-3"
    >
      Alamat & Catatan Kontrak
    </h4>
    <div class="space-y-3">
      <div>
        <label
          for="cust-address"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Alamat Kantor Penjamin
        </label>
        <input
          id="cust-address"
          type="text"
          bind:value={address}
          placeholder="Gedung / Jalan / Kota penjamin..."
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <div>
        <label
          for="cust-desc"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Keterangan Tambahan / Ketentuan PKS
        </label>
        <textarea
          id="cust-desc"
          bind:value={description}
          rows="2"
          placeholder="Ketentuan tarif khusus, plafon, atau nomor PKS kerja sama..."
          class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        ></textarea>
      </div>
    </div>
  </div>

  <!-- Section: Kebijakan & Aturan Sistem -->
  <div class="pt-2 border-t border-[#e1e5ea] dark:border-[#353738] space-y-3">
    <h4
      class="text-xs font-bold text-[#444746] dark:text-[#c4c7c5] uppercase tracking-wider mb-2"
    >
      Konfigurasi Aturan Transaksi
    </h4>

    <M3ToggleCard
      title="Wajib Kartu Fisik / Surat Eligibilitas Peserta (SEP)"
      description="Petugas pendaftaran wajib memverifikasi kartu peserta atau nomor SEP sebelum pelayanan."
      bind:checked={requireCard}
      colorVariant="primary"
    />

    <M3ToggleCard
      title="Status Operasional Penjamin"
      description="Nonaktifkan jika kerja sama PKS telah berakhir atau ditangguhkan sementara."
      bind:checked={isActive}
      colorVariant="emerald"
    />
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSubmitting}>
      {#if isSubmitting}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>{mode === "CREATE" ? "Simpan Penjamin" : "Perbarui Penjamin"}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
