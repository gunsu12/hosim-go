<script lang="ts">
  import {
    FolderTree,
    Landmark,
    CheckCircle2,
    AlertCircle,
    Info
  } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type {
    AccountRecord,
    AccountType,
    AccountPosition,
    CreateAccountDTO,
    UpdateAccountDTO
  } from '$lib/types/accounting/account';

  interface Props {
    isOpen: boolean;
    mode: 'create' | 'edit';
    account?: AccountRecord | null;
    parentForNewChild?: AccountRecord | null;
    flatAccounts?: AccountRecord[];
    isSaving: boolean;
    onClose: () => void;
    onSave: (payload: CreateAccountDTO | UpdateAccountDTO) => Promise<void>;
  }

  let {
    isOpen,
    mode = 'create',
    account = null,
    parentForNewChild = null,
    flatAccounts = [],
    isSaving = false,
    onClose,
    onSave
  }: Props = $props();

  let formCode = $state('');
  let formName = $state('');
  let formParentId = $state('');
  let formType = $state<AccountType>('ASSET');
  let formPosition = $state<AccountPosition>('DEBIT');
  let formDescription = $state('');
  let formIsPostable = $state(true);
  let formIsTreasury = $state(false);
  let formBankName = $state('');
  let formBankAccountNumber = $state('');
  let formCurrency = $state('IDR');
  let formIsActive = $state(true);
  let formError = $state<string | null>(null);

  const defaultPositionMap: Record<AccountType, AccountPosition> = {
    ASSET: 'DEBIT',
    LIABILITY: 'CREDIT',
    EQUITY: 'CREDIT',
    REVENUE: 'CREDIT',
    EXPENSE: 'DEBIT'
  };

  $effect(() => {
    if (isOpen) {
      formError = null;
      if (mode === 'edit' && account) {
        formCode = account.code;
        formName = account.name;
        formParentId = account.parent_id || '';
        formType = account.type;
        formPosition = account.position;
        formDescription = account.description || '';
        formIsPostable = account.is_postable;
        formIsTreasury = account.is_treasury_account;
        formBankName = account.bank_name || '';
        formBankAccountNumber = account.bank_account_number || '';
        formCurrency = account.currency || 'IDR';
        formIsActive = account.is_active;
      } else {
        formCode = '';
        formName = '';
        formDescription = '';
        formIsPostable = true;
        formIsTreasury = false;
        formBankName = '';
        formBankAccountNumber = '';
        formCurrency = 'IDR';
        formIsActive = true;

        if (parentForNewChild) {
          formParentId = parentForNewChild.id;
          formType = parentForNewChild.type;
          formPosition = parentForNewChild.position;
        } else {
          formParentId = '';
          formType = 'ASSET';
          formPosition = 'DEBIT';
        }
      }
    }
  });

  function handleTypeChange(newType: AccountType) {
    formType = newType;
    formPosition = defaultPositionMap[newType] || 'DEBIT';
  }

  function handleParentSelect(parentId: string) {
    formParentId = parentId;
    if (parentId) {
      const parentAcc = flatAccounts.find(a => a.id === parentId);
      if (parentAcc) {
        formType = parentAcc.type;
        formPosition = parentAcc.position;
      }
    }
  }

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    formError = null;

    if (!formCode.trim()) {
      formError = 'Kode akun wajib diisi';
      return;
    }
    if (!formName.trim()) {
      formError = 'Nama akun wajib diisi';
      return;
    }
    if (formIsTreasury && !formBankName.trim()) {
      formError = 'Nama Bank / Kas wajib diisi untuk Akun Kas & Bank';
      return;
    }

    try {
      if (mode === 'create') {
        const payload: CreateAccountDTO = {
          code: formCode.trim(),
          name: formName.trim(),
          parent_id: formParentId.trim() ? formParentId : null,
          type: formType,
          position: formPosition,
          description: formDescription.trim() || null,
          is_postable: formIsPostable,
          is_treasury_account: formIsTreasury,
          bank_name: formIsTreasury && formBankName.trim() ? formBankName.trim() : null,
          bank_account_number: formIsTreasury && formBankAccountNumber.trim() ? formBankAccountNumber.trim() : null,
          currency: formCurrency.trim() || 'IDR',
          is_active: formIsActive
        };
        await onSave(payload);
      } else {
        const payload: UpdateAccountDTO = {
          name: formName.trim(),
          position: formPosition,
          description: formDescription.trim() || null,
          is_treasury_account: formIsTreasury,
          bank_name: formIsTreasury && formBankName.trim() ? formBankName.trim() : null,
          bank_account_number: formIsTreasury && formBankAccountNumber.trim() ? formBankAccountNumber.trim() : null,
          currency: formCurrency.trim() || 'IDR',
          is_active: formIsActive
        };
        await onSave(payload);
      }
    } catch (err: any) {
      formError = err.message || 'Gagal menyimpan data akun';
    }
  }
</script>

<M3Modal
  {isOpen}
  title={mode === 'create' ? 'Tambah Akun Buku Besar Baru' : `Edit Akun: ${formCode}`}
  subtitle={mode === 'create' ? 'Konfigurasi akun baru pada bagan akun standar RS' : 'Perbarui nama, posisi saldo, atau informasi perbankan'}
  icon={FolderTree}
  size="2xl"
  onsubmit={handleSubmit}
  {onClose}
>
  {#if formError}
    <div class="p-3 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5 text-xs">
      <AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
      <span>{formError}</span>
    </div>
  {/if}

  <!-- Mode Tambah Sub-Akun Banner Info -->
  {#if mode === 'create' && parentForNewChild}
    <div class="p-3 rounded-2xl bg-blue-50 dark:bg-blue-950/40 border border-blue-200 dark:border-blue-800 text-blue-900 dark:text-blue-200 flex items-center justify-between text-xs">
      <div class="flex items-center gap-2">
        <Info class="w-4 h-4 text-blue-700 dark:text-blue-400 shrink-0" />
        <span>Membuat sub-akun di bawah induk: <strong>{parentForNewChild.code} - {parentForNewChild.name}</strong></span>
      </div>
      <span class="text-[10px] font-mono px-2 py-0.5 rounded bg-blue-100 dark:bg-blue-900/60 text-blue-800 dark:text-blue-300 font-bold">
        Level {parentForNewChild.account_level + 1}
      </span>
    </div>
  {/if}

  <!-- Section 1: Kode & Nama Akun -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="coa-form-code" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Kode Akun <span class="text-rose-500">*</span>
      </label>
      <input
        id="coa-form-code"
        type="text"
        bind:value={formCode}
        disabled={mode === 'edit'}
        placeholder="Contoh: 11101 atau 111.01"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-mono font-bold uppercase text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all disabled:opacity-60"
      />
      {#if mode === 'edit'}
        <span class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1 block">
          Kode akun bersifat tetap (immutable) untuk menjaga konsistensi buku besar.
        </span>
      {/if}
    </div>

    <div>
      <label for="coa-form-name" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Nama Akun <span class="text-rose-500">*</span>
      </label>
      <input
        id="coa-form-name"
        type="text"
        bind:value={formName}
        placeholder="Contoh: Kas Kasir Rawat Jalan"
        required
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <!-- Section 2: Parent Akun & Tipe Akun -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
    <div>
      <label for="coa-form-parent" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Akun Induk (Parent)
      </label>
      <select
        id="coa-form-parent"
        value={formParentId}
        onchange={(e) => handleParentSelect((e.target as HTMLSelectElement).value)}
        disabled={mode === 'edit'}
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all disabled:opacity-60"
      >
        <option value="">(Tanpa Induk - Akun Kepala Level 1)</option>
        {#each flatAccounts as pAcc}
          {#if mode === 'create' || pAcc.id !== account?.id}
            <option value={pAcc.id}>
              [{pAcc.code}] {pAcc.name} (L{pAcc.account_level} • {pAcc.type})
            </option>
          {/if}
        {/each}
      </select>
      {#if mode === 'edit'}
        <span class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1 block">
          Struktur hierarki induk tidak dapat dipindahkan setelah akun memiliki buku besar.
        </span>
      {/if}
    </div>

    <div>
      <label for="coa-form-type" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Klasifikasi Tipe Akun <span class="text-rose-500">*</span>
      </label>
      <select
        id="coa-form-type"
        value={formType}
        onchange={(e) => handleTypeChange((e.target as HTMLSelectElement).value as AccountType)}
        disabled={mode === 'edit' || (mode === 'create' && !!formParentId)}
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all disabled:opacity-60"
      >
        <option value="ASSET">ASSET (1xxx - Aset / Aktiva)</option>
        <option value="LIABILITY">LIABILITY (2xxx - Kewajiban / Hutang)</option>
        <option value="EQUITY">EQUITY (3xxx - Ekuitas / Modal)</option>
        <option value="REVENUE">REVENUE (4xxx - Pendapatan)</option>
        <option value="EXPENSE">EXPENSE (5xxx - Beban Operasional)</option>
      </select>
    </div>
  </div>

  <!-- Section 3: Posisi Saldo Normal & Sifat Postable -->
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 p-4 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738]">
    <div>
      <span class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1.5 text-xs">
        Posisi Saldo Normal <span class="text-rose-500">*</span>
      </span>
      <div class="flex items-center gap-4">
        <label class="flex items-center gap-2 cursor-pointer">
          <input
            type="radio"
            name="position"
            value="DEBIT"
            checked={formPosition === 'DEBIT'}
            onchange={() => formPosition = 'DEBIT'}
            class="text-[#0b57d0] focus:ring-[#0b57d0]"
          />
          <span class="font-bold text-xs text-[#1f1f1f] dark:text-[#e3e3e3]">DEBIT</span>
        </label>
        <label class="flex items-center gap-2 cursor-pointer">
          <input
            type="radio"
            name="position"
            value="CREDIT"
            checked={formPosition === 'CREDIT'}
            onchange={() => formPosition = 'CREDIT'}
            class="text-[#0b57d0] focus:ring-[#0b57d0]"
          />
          <span class="font-bold text-xs text-[#1f1f1f] dark:text-[#e3e3e3]">KREDIT</span>
        </label>
      </div>
      <span class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-1.5 block">
        Standar: Aset & Beban di Debit; Kewajiban, Ekuitas & Pendapatan di Kredit.
      </span>
    </div>

    <div>
      <span class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1.5 text-xs">
        Sifat Akun Transaksi
      </span>
      <label class="flex items-start gap-2.5 cursor-pointer mt-1">
        <input
          type="checkbox"
          bind:checked={formIsPostable}
          disabled={mode === 'edit'}
          class="mt-0.5 rounded border-[#c4c7c5] dark:border-[#555] text-[#0b57d0] focus:ring-[#0b57d0] disabled:opacity-50"
        />
        <div>
          <span class="font-bold text-xs text-[#1f1f1f] dark:text-[#e3e3e3]">Akun Postable (Dapat Dijurnal)</span>
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block mt-0.5">
            Hilangkan centang jika akun ini hanya sebagai Header / Induk ringkasan.
          </span>
        </div>
      </label>
    </div>
  </div>

  <!-- Section 4: Kas & Bank / Treasury Flag -->
  <div class="p-4 rounded-2xl border border-teal-200/80 dark:border-teal-800/60 bg-teal-50/50 dark:bg-teal-950/25 flex flex-col gap-3 transition-colors">
    <div class="flex items-start justify-between gap-3">
      <label class="flex items-start gap-2.5 cursor-pointer flex-1">
        <div class="w-8 h-8 rounded-xl bg-teal-100 dark:bg-teal-900/60 text-teal-700 dark:text-teal-300 flex items-center justify-center shrink-0 mt-0.5">
          <Landmark class="w-4 h-4" />
        </div>
        <div class="flex-1">
          <span class="font-bold text-teal-950 dark:text-teal-200 block text-xs">
            Merupakan Akun Kas / Bank (Treasury Account)
          </span>
          <span class="text-[11px] text-teal-800/90 dark:text-teal-300/80 block mt-0.5 leading-relaxed">
            Centang opsi ini jika akun ini mewakili kas tunai, rekening giro operasional, rekening penampung BPJS, atau kas kecil RS.
          </span>
        </div>
      </label>
      <label class="relative inline-flex items-center cursor-pointer shrink-0 mt-1">
        <input
          type="checkbox"
          bind:checked={formIsTreasury}
          class="sr-only peer"
        />
        <div
          class="w-11 h-6 bg-gray-200 dark:bg-gray-700 peer-focus:outline-hidden rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 dark:after:border-gray-600 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-teal-600 dark:peer-checked:bg-teal-500"
        ></div>
      </label>
    </div>

    {#if formIsTreasury}
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-3 border-t border-teal-200 dark:border-teal-800/60 animate-fadeIn">
        <div>
          <label for="coa-form-bank-name" class="block font-semibold text-teal-950 dark:text-teal-200 mb-1">
            Nama Bank / Entitas Kas <span class="text-rose-500">*</span>
          </label>
          <input
            id="coa-form-bank-name"
            type="text"
            bind:value={formBankName}
            placeholder="Contoh: Bank Mandiri, BCA, atau Kas Induk RS"
            class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-teal-300 dark:border-teal-700/60 focus:border-teal-600 dark:focus:border-teal-400 text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
          />
        </div>

        <div>
          <label for="coa-form-bank-account" class="block font-semibold text-teal-950 dark:text-teal-200 mb-1">
            Nomor Rekening Bank
          </label>
          <input
            id="coa-form-bank-account"
            type="text"
            bind:value={formBankAccountNumber}
            placeholder="Contoh: 137-00-1234567-8"
            class="w-full h-10 px-3.5 rounded-xl bg-white dark:bg-[#1e1f20] border border-teal-300 dark:border-teal-700/60 focus:border-teal-600 dark:focus:border-teal-400 text-xs font-mono font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
          />
        </div>
      </div>
    {/if}
  </div>

  <!-- Section 5: Keterangan & Status -->
  <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
    <div class="sm:col-span-2">
      <label for="coa-form-desc" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
        Keterangan / Deskripsi Akun
      </label>
      <textarea
        id="coa-form-desc"
        bind:value={formDescription}
        rows="2"
        placeholder="Penjelasan fungsi akun buku besar ini..."
        class="w-full p-3 text-xs bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-[#1f1f1f] dark:text-[#e3e3e3] rounded-xl focus:outline-hidden transition-all resize-none"
      ></textarea>
    </div>

    <div class="flex flex-col justify-between">
      <div>
        <label for="coa-form-currency" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-1">
          Mata Uang
        </label>
        <input
          id="coa-form-currency"
          type="text"
          bind:value={formCurrency}
          class="w-full h-10 px-3 text-xs bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-[#1f1f1f] dark:text-[#e3e3e3] rounded-xl font-mono text-center font-bold focus:outline-hidden transition-all"
        />
      </div>

      <label class="flex items-center gap-2 cursor-pointer mt-3">
        <input
          type="checkbox"
          bind:checked={formIsActive}
          class="rounded border-[#c4c7c5] dark:border-[#555] text-[#0b57d0] focus:ring-[#0b57d0]"
        />
        <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3] text-xs">Akun Aktif Digunakan</span>
      </label>
    </div>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSaving}>
      Batal
    </M3Button>
    <M3Button variant="filled" type="submit" disabled={isSaving}>
      {#if isSaving}
        <span class="inline-block animate-spin mr-1">⏳</span>
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1" />
        <span>{mode === 'create' ? 'Buat Akun' : 'Simpan Perubahan'}</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
