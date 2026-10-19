<script lang="ts">
  import { RotateCw, AlertCircle, Calendar, FileText, Check } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3DatePicker from '$lib/components/m3/M3DatePicker.svelte';
  import type {
    CustomerContractRecord,
    RenewContractDTO,
  } from '$lib/types/finance/customer_contract';

  interface Props {
    isOpen: boolean;
    contract: CustomerContractRecord | null;
    isSubmitting?: boolean;
    errorMessage?: string | null;
    onClose: () => void;
    onSubmit: (dto: RenewContractDTO) => Promise<void> | void;
  }

  let {
    isOpen,
    contract,
    isSubmitting = false,
    errorMessage = null,
    onClose,
    onSubmit,
  }: Props = $props();

  let contractNumber = $state('');
  let name = $state('');
  let effectiveDate = $state('');
  let hasExpirationDate = $state(true);
  let expirationDate = $state('');
  let notes = $state('');

  $effect(() => {
    if (isOpen && contract) {
      const nextYear = new Date().getFullYear() + 1;
      contractNumber = `${contract.contract_number}/RENEW-${nextYear}`;
      name = `${contract.name} (Tahun ${nextYear})`;

      // Suggest effective date following expiration date of old contract
      if (contract.expiration_date) {
        const d = new Date(contract.expiration_date);
        d.setDate(d.getDate() + 1);
        effectiveDate = d.toISOString().split('T')[0];

        const exp = new Date(d);
        exp.setFullYear(exp.getFullYear() + 1);
        expirationDate = exp.toISOString().split('T')[0];
        hasExpirationDate = true;
      } else {
        effectiveDate = new Date().toISOString().split('T')[0];
        expirationDate = '';
        hasExpirationDate = false;
      }

      notes = `Perpanjangan dari PKS: ${contract.contract_number}`;
    }
  });

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!contractNumber.trim() || !name.trim() || !effectiveDate) return;

    const dto: RenewContractDTO = {
      contract_number: contractNumber.trim().toUpperCase(),
      name: name.trim(),
      effective_date: effectiveDate,
      expiration_date: hasExpirationDate && expirationDate ? expirationDate : null,
      notes: notes.trim() || undefined,
    };

    onSubmit(dto);
  }
</script>

<M3Modal
  {isOpen}
  title="Perpanjang Kontrak PKS (Renewal)"
  subtitle={contract ? `Mengacu pada PKS ${contract.contract_number}` : ''}
  icon={RotateCw}
  iconClass="bg-[#0b57d0]/10 text-[#0b57d0] dark:bg-[#a8c7fa]/15 dark:text-[#a8c7fa]"
  size="xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if contract}
    <div class="space-y-4 text-xs">
      <!-- Error Message in Modal -->
      {#if errorMessage}
        <div
          class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/60 border border-rose-300 dark:border-rose-800 text-rose-800 dark:text-rose-200 text-xs flex items-start gap-2.5 shadow-2xs"
        >
          <AlertCircle class="w-4 h-4 text-rose-600 dark:text-rose-400 shrink-0 mt-0.5" />
          <span class="font-medium leading-relaxed">{errorMessage}</span>
        </div>
      {/if}

      <!-- Renewal Architecture Explanation Box -->
      <div
        class="p-3.5 rounded-2xl bg-blue-50/80 dark:bg-blue-950/40 border border-blue-200 dark:border-blue-900/60 text-blue-900 dark:text-blue-200 leading-relaxed"
      >
        <p class="font-semibold mb-1">Prinsip Immutability Perpanjangan:</p>
        <p class="text-[11px] text-blue-800/90 dark:text-blue-300/80">
          Sistem akan membuat kontrak <strong>DRAFT baru</strong> dengan referensi ke kontrak lama ini.
          Termin pembayaran (Net {contract.payment_term_days} hari), batas kredit, syarat jaminan, serta
          seluruh Price Schema terhubung akan disalin secara otomatis ke kontrak baru.
        </p>
      </div>

      <!-- Nomor Kontrak Baru -->
      <div>
        <label
          for="renew-contract-number"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Nomor Kontrak Baru (PKS) <span class="text-rose-500">*</span>
        </label>
        <input
          id="renew-contract-number"
          type="text"
          required
          bind:value={contractNumber}
          placeholder="cth: PKS/2027/RS/BPJS/001"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <!-- Nama Kontrak Baru -->
      <div>
        <label
          for="renew-contract-name"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Nama / Judul Kontrak Baru <span class="text-rose-500">*</span>
        </label>
        <input
          id="renew-contract-name"
          type="text"
          required
          bind:value={name}
          placeholder="cth: PKS Pelayanan Rawat Inap & Rawat Jalan Periode 2027"
          class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
      </div>

      <!-- Tanggal Berlaku -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
        <div>
          <M3DatePicker
            label="Tanggal Mulai Berlaku Baru"
            required
            bind:value={effectiveDate}
            placeholder="Pilih tanggal mulai"
          />
        </div>

        <div>
          <div class="flex items-center justify-between mb-1">
            <label for="renew-expiry-input" class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
              Tanggal Berakhir
            </label>
            <label class="flex items-center gap-1.5 cursor-pointer">
              <input
                type="checkbox"
                bind:checked={hasExpirationDate}
                class="rounded text-[#0b57d0] focus:ring-0"
              />
              <span class="text-[10px] text-[#747775] dark:text-[#8e918f]">Ada Batas</span>
            </label>
          </div>
          {#if hasExpirationDate}
            <M3DatePicker
              id="renew-expiry-input"
              bind:value={expirationDate}
              min={effectiveDate}
              placeholder="Pilih tanggal berakhir"
            />
          {:else}
            <div
              class="h-10 px-3.5 rounded-xl bg-gray-100 dark:bg-gray-800 text-[#747775] dark:text-[#8e918f] flex items-center font-medium"
            >
              Tanpa batas waktu
            </div>
          {/if}
        </div>
      </div>

      <!-- Catatan -->
      <div>
        <label
          for="renew-notes"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Catatan Tambahan
        </label>
        <textarea
          id="renew-notes"
          rows="2"
          bind:value={notes}
          placeholder="Keterangan perpanjangan..."
          class="w-full p-3 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        ></textarea>
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting}>
      Batal
    </M3Button>
    <M3Button
      variant="filled"
      type="submit"
      disabled={isSubmitting || !contractNumber.trim() || !name.trim() || !effectiveDate}
    >
      {#if isSubmitting}
        <div
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-1"
        ></div>
        <span>Membuat Perpanjangan...</span>
      {:else}
        <RotateCw class="w-4 h-4 mr-1" />
        <span>Buat Kontrak Baru (Renew)</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
