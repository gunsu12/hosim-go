<script lang="ts">
  import {
    X,
    ShieldCheck,
    Building2,
    Phone,
    Mail,
    Globe,
    MapPin,
    FileText,
    CheckCircle2,
  } from "@lucide/svelte";
  import M3Button from "$lib/components/m3/M3Button.svelte";
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

{#if isOpen}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4 overflow-y-auto"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-2xl max-h-[92vh] flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Header -->
      <div
        class="px-6 py-4 border-b border-[#e1e5ea] flex items-center justify-between bg-[#f8fafd]"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-2xl bg-[#0b57d0]/10 text-[#0b57d0] flex items-center justify-center font-bold"
          >
            <ShieldCheck class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-bold text-[#1f1f1f]">
              {mode === "CREATE"
                ? "Tambah Debitur / Penjamin Baru"
                : "Ubah Data Debitur / Penjamin"}
            </h3>
            <p class="text-xs text-[#444746]">
              {mode === "CREATE"
                ? "Konfigurasi penjamin pembiayaan (BPJS, Asuransi Swasta, Korporat)"
                : `Kode: ${customer?.code || "-"}`}
            </p>
          </div>
        </div>
        <button
          onclick={onClose}
          type="button"
          class="text-[#747775] hover:text-[#1f1f1f] p-1.5 rounded-full hover:bg-[#e1e5ea] transition-colors cursor-pointer"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Scrollable Form Body -->
      <form
        onsubmit={handleSubmit}
        class="flex-1 overflow-y-auto p-6 space-y-4 text-xs"
      >
        {#if errorMessage}
          <div
            class="p-3 rounded-2xl bg-rose-50 border border-rose-200 text-rose-700 flex items-start gap-2.5"
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
              class="block font-semibold text-[#1f1f1f] mb-1"
            >
              Kode Penjamin <span class="text-rose-500">*</span>
            </label>
            <input
              id="cust-code"
              type="text"
              bind:value={code}
              placeholder="Contoh: BPJS-JKN, ASR-PRUD"
              required
              class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs font-mono font-bold uppercase text-[#1f1f1f] focus:outline-hidden transition-all"
            />
            <p class="text-[11px] text-[#747775] mt-1">
              Harus unik di database sistem RS.
            </p>
          </div>

          <div>
            <label
              for="cust-type"
              class="block font-semibold text-[#1f1f1f] mb-1"
            >
              Kategori Penjamin
            </label>
            <select
              id="cust-type"
              bind:value={customerTypeId}
              class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all"
            >
              <option value="">-- Pilih Kategori / Tipe --</option>
              {#each customerTypes as ct (ct.id)}
                <option value={ct.id}>{ct.name}</option>
              {/each}
            </select>
          </div>
        </div>

        <div>
          <label
            for="cust-name"
            class="block font-semibold text-[#1f1f1f] mb-1"
          >
            Nama Instansi / Perusahaan Penjamin <span class="text-rose-500"
              >*</span
            >
          </label>
          <input
            id="cust-name"
            type="text"
            bind:value={name}
            placeholder="Contoh: BPJS Kesehatan KC Jakarta, PT Prudential Life Assurance"
            required
            class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs font-semibold text-[#1f1f1f] focus:outline-hidden transition-all"
          />
        </div>

        <!-- Section: Kontak & PIC -->
        <div class="pt-2 border-t border-[#e1e5ea]">
          <h4
            class="text-xs font-bold text-[#444746] uppercase tracking-wider mb-3"
          >
            Kontak & Perwakilan (PIC)
          </h4>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label
                for="cust-contact-person"
                class="block font-semibold text-[#1f1f1f] mb-1"
              >
                Nama PIC / Koordinator Klaim
              </label>
              <input
                id="cust-contact-person"
                type="text"
                bind:value={contactPerson}
                placeholder="Contoh: Budi Santoso (Head of Claims)"
                class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all"
              />
            </div>

            <div>
              <label
                for="cust-phone"
                class="block font-semibold text-[#1f1f1f] mb-1"
              >
                Telepon / Hotline Call Center
              </label>
              <input
                id="cust-phone"
                type="text"
                bind:value={phone}
                placeholder="Contoh: 1500085 atau 021-29958888"
                class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all font-mono"
              />
            </div>

            <div>
              <label
                for="cust-email"
                class="block font-semibold text-[#1f1f1f] mb-1"
              >
                Email Korespondensi / Klaim
              </label>
              <input
                id="cust-email"
                type="email"
                bind:value={email}
                placeholder="Contoh: claim.hospital@insurance.co.id"
                class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all"
              />
            </div>

            <div>
              <label
                for="cust-website"
                class="block font-semibold text-[#1f1f1f] mb-1"
              >
                Situs Resmi / Portal Klaim
              </label>
              <input
                id="cust-website"
                type="text"
                bind:value={website}
                placeholder="Contoh: https://portal.asuransi.co.id"
                class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all"
              />
            </div>
          </div>
        </div>

        <!-- Section: Alamat & Deskripsi -->
        <div class="pt-2 border-t border-[#e1e5ea]">
          <h4
            class="text-xs font-bold text-[#444746] uppercase tracking-wider mb-3"
          >
            Alamat & Catatan Kontrak
          </h4>
          <div class="space-y-3">
            <div>
              <label
                for="cust-address"
                class="block font-semibold text-[#1f1f1f] mb-1"
              >
                Alamat Kantor Penjamin
              </label>
              <input
                id="cust-address"
                type="text"
                bind:value={address}
                placeholder="Gedung / Jalan / Kota penjamin..."
                class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all"
              />
            </div>

            <div>
              <label
                for="cust-desc"
                class="block font-semibold text-[#1f1f1f] mb-1"
              >
                Keterangan Tambahan / Ketentuan PKS
              </label>
              <textarea
                id="cust-desc"
                bind:value={description}
                rows="2"
                placeholder="Ketentuan tarif khusus, plafon, atau nomor PKS kerja sama..."
                class="w-full p-3 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-hidden transition-all"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Section: Kebijakan & Aturan Sistem -->
        <div class="pt-2 border-t border-[#e1e5ea] space-y-3">
          <h4
            class="text-xs font-bold text-[#444746] uppercase tracking-wider mb-2"
          >
            Konfigurasi Aturan Transaksi
          </h4>

          <div
            class="flex items-center justify-between p-3 rounded-2xl bg-[#f8fafd] border border-[#e1e5ea]"
          >
            <div>
              <span class="font-bold text-[#1f1f1f] block"
                >Wajib Kartu Fisik / Surat Eligibilitas Peserta (SEP)</span
              >
              <span class="text-[11px] text-[#444746]">
                Petugas pendaftaran wajib memverifikasi kartu peserta atau nomor
                SEP sebelum pelayanan.
              </span>
            </div>
            <label
              class="relative inline-flex items-center cursor-pointer ml-4"
            >
              <input
                type="checkbox"
                bind:checked={requireCard}
                class="sr-only peer"
              />
              <div
                class="w-11 h-6 bg-gray-200 peer-focus:outline-hidden rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#0b57d0]"
              ></div>
            </label>
          </div>

          <div
            class="flex items-center justify-between p-3 rounded-2xl bg-[#f8fafd] border border-[#e1e5ea]"
          >
            <div>
              <span class="font-bold text-[#1f1f1f] block"
                >Status Operasional Penjamin</span
              >
              <span class="text-[11px] text-[#444746]">
                Nonaktifkan jika kerja sama PKS telah berakhir atau ditangguhkan
                sementara.
              </span>
            </div>
            <label
              class="relative inline-flex items-center cursor-pointer ml-4"
            >
              <input
                type="checkbox"
                bind:checked={isActive}
                class="sr-only peer"
              />
              <div
                class="w-11 h-6 bg-gray-200 peer-focus:outline-hidden rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-emerald-600"
              ></div>
            </label>
          </div>
        </div>

        <!-- Actions -->
        <div
          class="pt-4 border-t border-[#e1e5ea] flex items-center justify-end gap-3"
        >
          <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
            Batal
          </M3Button>
          <M3Button variant="filled" type="submit" disabled={isSubmitting}>
            {#if isSubmitting}
              <span class="inline-block animate-spin mr-1">⏳</span>
              <span>Menyimpan...</span>
            {:else}
              <CheckCircle2 class="w-4 h-4 mr-1" />
              <span
                >{mode === "CREATE"
                  ? "Simpan Penjamin"
                  : "Perbarui Penjamin"}</span
              >
            {/if}
          </M3Button>
        </div>
      </form>
    </div>
  </div>
{/if}
