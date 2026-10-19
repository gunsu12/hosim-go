<script lang="ts">
  import {
    FileSignature,
    AlertCircle,
    Calendar,
    Clock,
    DollarSign,
    Shield,
    Check,
  } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import M3Select, { type SelectOption } from '$lib/components/m3/M3Select.svelte';
  import M3DatePicker from '$lib/components/m3/M3DatePicker.svelte';
  import type {
    CustomerContractRecord,
    AmendContractTermsDTO,
    GuaranteeRequirement,
  } from '$lib/types/finance/customer_contract';

  interface Props {
    isOpen: boolean;
    contract: CustomerContractRecord | null;
    isSubmitting?: boolean;
    errorMessage?: string | null;
    onClose: () => void;
    onSubmit: (dto: AmendContractTermsDTO) => Promise<void> | void;
  }

  let {
    isOpen,
    contract,
    isSubmitting = false,
    errorMessage = null,
    onClose,
    onSubmit,
  }: Props = $props();

  // Form states
  let paymentTermDays = $state<number>(30);
  let hasCreditLimit = $state<boolean>(true);
  let creditLimitAmount = $state<number | string>(0);
  let hasExpirationDate = $state<boolean>(true);
  let expirationDate = $state<string>('');
  let guaranteeRequirement = $state<GuaranteeRequirement>('NONE');
  let reason = $state<string>('');

  const guaranteeOptions: SelectOption[] = [
    { value: 'NONE', label: 'Tanpa Jaminan Khusus (NONE)' },
    { value: 'SPPK', label: 'Wajib SPPK / Surat Eligibilitas Peserta' },
    { value: 'LETTER_OF_GUARANTEE', label: 'Wajib Surat Jaminan / Letter of Guarantee (GL)' },
  ];

  $effect(() => {
    if (isOpen && contract) {
      paymentTermDays = contract.payment_term_days || 30;
      if (contract.credit_limit !== null && contract.credit_limit !== undefined) {
        hasCreditLimit = true;
        creditLimitAmount = contract.credit_limit;
      } else {
        hasCreditLimit = false;
        creditLimitAmount = 0;
      }

      if (contract.expiration_date) {
        hasExpirationDate = true;
        expirationDate = contract.expiration_date;
      } else {
        hasExpirationDate = false;
        expirationDate = '';
      }

      guaranteeRequirement = contract.guarantee_requirement || 'NONE';
      reason = '';
    }
  });

  function handleFormSubmit(e: SubmitEvent) {
    e.preventDefault();
    if (!reason.trim()) return;

    const dto: AmendContractTermsDTO = {
      payment_term_days: Number(paymentTermDays) || 30,
      credit_limit: hasCreditLimit ? Number(creditLimitAmount) || 0 : null,
      expiration_date: hasExpirationDate && expirationDate ? expirationDate : null,
      guarantee_requirement: guaranteeRequirement,
      reason: reason.trim(),
    };

    onSubmit(dto);
  }
</script>

<M3Modal
  {isOpen}
  title="Amandemen / Addendum Kontrak PKS"
  subtitle={contract ? `${contract.contract_number} • ${contract.name}` : ''}
  icon={FileSignature}
  iconClass="bg-[#0b57d0]/10 text-[#0b57d0] dark:bg-[#a8c7fa]/15 dark:text-[#a8c7fa]"
  size="xl"
  onsubmit={handleFormSubmit}
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

      <!-- Warning Box -->
      <div
        class="p-3.5 rounded-2xl bg-amber-50 dark:bg-amber-950/40 border border-amber-300 dark:border-amber-800 text-amber-900 dark:text-amber-200 leading-relaxed"
      >
        <strong>Klausul Finansial:</strong> Perubahan termin pembayaran, batas kredit, tanggal berakhir,
        dan syarat jaminan pada kontrak yang sudah aktif memerlukan alasan formal (Addendum PKS) untuk
        dicatat dalam audit keuangan.
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <!-- Termin Pembayaran -->
        <div>
          <label
            for="amend-payment-term"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Termin Pembayaran (Hari) <span class="text-rose-500">*</span>
          </label>
          <div class="relative flex items-center">
            <input
              id="amend-payment-term"
              type="number"
              min="1"
              max="365"
              required
              bind:value={paymentTermDays}
              class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
            />
            <span class="absolute right-3 text-[#747775] dark:text-[#8e918f] font-medium text-[11px]">
              Net {paymentTermDays}
            </span>
          </div>
          <span class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-0.5 block">
            Sebelumnya: Net {contract.payment_term_days} hari
          </span>
        </div>

        <!-- Syarat Jaminan -->
        <div>
          <label
            for="amend-guarantee"
            class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
          >
            Syarat Jaminan Per Kunjungan <span class="text-rose-500">*</span>
          </label>
          <M3Select
            id="amend-guarantee"
            options={guaranteeOptions}
            bind:value={guaranteeRequirement}
            placeholder="Pilih syarat jaminan..."
          />
        </div>
      </div>

      <!-- Batas Kredit / Plafon -->
      <div
        class="p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] space-y-2.5"
      >
        <div class="flex items-center justify-between">
          <label for="amend-limit-toggle" class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
            Batas Kredit / Plafon Tagihan
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              id="amend-limit-toggle"
              type="checkbox"
              bind:checked={hasCreditLimit}
              class="rounded text-[#0b57d0] focus:ring-0"
            />
            <span class="text-[11px] text-[#444746] dark:text-[#c4c7c5]">Batasi Plafon Nominal</span>
          </label>
        </div>

        {#if hasCreditLimit}
          <div>
            <div class="relative flex items-center">
              <span class="absolute left-3.5 text-[#747775] dark:text-[#8e918f] font-semibold text-xs">
                Rp
              </span>
              <input
                type="number"
                min="0"
                step="1000000"
                bind:value={creditLimitAmount}
                placeholder="cth: 500000000"
                class="w-full h-10 pl-10 pr-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden"
              />
            </div>
            <p class="text-[10px] text-[#747775] dark:text-[#8e918f] mt-1">
              Limit Rp {(Number(creditLimitAmount) || 0).toLocaleString('id-ID')}
            </p>
          </div>
        {:else}
          <p class="text-[11px] text-emerald-600 dark:text-emerald-400 font-medium">
            Tanpa batas kredit (Unlimited Credit / Bebas Plafon)
          </p>
        {/if}
      </div>

      <!-- Tanggal Berakhir -->
      <div
        class="p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] space-y-2.5"
      >
        <div class="flex items-center justify-between">
          <label for="amend-expiry-toggle" class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
            Masa Berlaku Kontrak
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              id="amend-expiry-toggle"
              type="checkbox"
              bind:checked={hasExpirationDate}
              class="rounded text-[#0b57d0] focus:ring-0"
            />
            <span class="text-[11px] text-[#444746] dark:text-[#c4c7c5]">Memiliki Batas Waktu</span>
          </label>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 items-center">
          <div>
            <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block mb-1">
              Tanggal Mulai Berlaku (Tetap)
            </span>
            <div
              class="h-10 px-3.5 rounded-xl bg-gray-100 dark:bg-gray-800 text-[#747775] dark:text-[#8e918f] flex items-center font-mono font-semibold"
            >
              {contract.effective_date}
            </div>
          </div>

          <div>
            {#if hasExpirationDate}
              <M3DatePicker
                label="Tanggal Berakhir (Baru)"
                bind:value={expirationDate}
                min={contract.effective_date}
                placeholder="Pilih tanggal berakhir"
              />
            {:else}
              <div
                class="h-10 px-3.5 rounded-xl bg-emerald-50 dark:bg-emerald-950/40 text-emerald-700 dark:text-emerald-300 flex items-center font-medium mt-4"
              >
                Tanpa batas waktu (Otomatis Diperpanjang)
              </div>
            {/if}
          </div>
        </div>
      </div>

      <!-- Alasan Amandemen Wajib -->
      <div>
        <label
          for="amend-reason-input"
          class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1"
        >
          Nomor Dokumen & Alasan Addendum / Amandemen <span class="text-rose-500">*</span>
        </label>
        <textarea
          id="amend-reason-input"
          rows="3"
          required
          bind:value={reason}
          placeholder="Contoh: Sesuai Addendum I No. 042/ADD/2026 tanggal 15 Mei 2026 perihal penyesuaian batas kredit dan syarat penjaminan..."
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
      disabled={isSubmitting || !reason.trim()}
    >
      {#if isSubmitting}
        <div
          class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-1"
        ></div>
        <span>Menyimpan Amandemen...</span>
      {:else}
        <Check class="w-4 h-4 mr-1" />
        <span>Simpan Amandemen (Addendum)</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
