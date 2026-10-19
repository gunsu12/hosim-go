<script lang="ts">
  import { Building2, Edit, FileText } from '@lucide/svelte';
  import M3Modal from '$lib/components/m3/M3Modal.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { CustomerRecord } from '$lib/types/finance';

  interface Props {
    isOpen: boolean;
    customer?: CustomerRecord | null;
    onClose: () => void;
    onEdit: (customer: CustomerRecord) => void;
  }

  let {
    isOpen = false,
    customer = null,
    onClose = () => {},
    onEdit = () => {},
  }: Props = $props();
</script>

<M3Modal
  {isOpen}
  title="Detail Debitur / Penjamin"
  subtitle={customer ? `${customer.code} • ${customer.name}` : ''}
  icon={Building2}
  size="xl"
  {onClose}
>
  {#if customer}
    <div class="space-y-4 text-xs">
      <!-- Title & Status Banner -->
      <div class="p-4 rounded-2xl bg-[#f0f4f9] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] flex flex-wrap items-center justify-between gap-3">
        <div>
          <h4 class="text-sm font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">{customer.name}</h4>
          <p class="text-xs text-[#444746] dark:text-[#c4c7c5] mt-0.5">
            Kategori: <strong>{customer.customer_type?.name || 'Umum / Lainnya'}</strong>
          </p>
        </div>
        <div class="flex items-center gap-2">
          {#if customer.is_active}
            <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300">
              AKTIF
            </span>
          {:else}
            <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-gray-200 dark:bg-gray-800 text-gray-700 dark:text-gray-300">
              NON-AKTIF
            </span>
          {/if}
        </div>
      </div>

      <!-- Info Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
        <div class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]">
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Person In Charge (PIC)</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5">{customer.contact_person || '-'}</span>
        </div>

        <div class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]">
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Telepon / Call Center</span>
          <span class="font-mono font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5">{customer.phone || '-'}</span>
        </div>

        <div class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]">
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Email Korespondensi</span>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5">{customer.email || '-'}</span>
        </div>

        <div class="p-3 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]">
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block">Situs Resmi / Portal</span>
          {#if customer.website && (customer.website.startsWith('http://') || customer.website.startsWith('https://'))}
            <a href={customer.website} target="_blank" rel="noopener noreferrer" class="text-[#0b57d0] dark:text-[#a8c7fa] hover:underline block mt-0.5 font-semibold truncate">
              {customer.website}
            </a>
          {:else if customer.website}
            <span class="text-[#1f1f1f] dark:text-[#e3e3e3] block mt-0.5 font-mono truncate">{customer.website}</span>
          {:else}
            <span class="text-[#747775] dark:text-[#8e918f] block mt-0.5">-</span>
          {/if}
        </div>
      </div>

      <!-- Alamat -->
      <div class="p-3.5 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]">
        <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block font-semibold mb-1">Alamat Kantor / Instansi</span>
        <p class="text-[#1f1f1f] dark:text-[#e3e3e3]">{customer.address || 'Belum ada data alamat tercatat.'}</p>
      </div>

      <!-- Keterangan -->
      <div class="p-3.5 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-white dark:bg-[#1e1f20]">
        <span class="text-[11px] text-[#747775] dark:text-[#8e918f] block font-semibold mb-1">Catatan Tambahan & PKS</span>
        <p class="text-[#444746] dark:text-[#c4c7c5] whitespace-pre-wrap">{customer.description || 'Tidak ada catatan khusus.'}</p>
      </div>

      <!-- Aturan Peserta -->
      <div class="p-3.5 rounded-xl border border-[#e1e5ea] dark:border-[#353738] bg-[#f8fafd] dark:bg-[#28292a] flex items-center justify-between">
        <div>
          <span class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] block">Verifikasi Kartu / SEP Peserta</span>
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f]">Wajib diverifikasi saat pasien mendaftar di unit layanan</span>
        </div>
        <div>
          {#if customer.require_card}
            <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-amber-100 dark:bg-amber-950/60 text-amber-800 dark:text-amber-300">
              Wajib Kartu / SEP
            </span>
          {:else}
            <span class="px-2.5 py-1 rounded-full text-[10px] font-medium bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300">
              Tanpa Kartu
            </span>
          {/if}
        </div>
      </div>

      <!-- Audit Meta -->
      <div class="pt-2 text-[11px] text-[#747775] dark:text-[#8e918f] flex flex-wrap justify-between gap-2 border-t border-[#e1e5ea] dark:border-[#353738]">
        <span>ID Sistem: <strong class="font-mono">{customer.id}</strong></span>
        {#if customer.created_at}
          <span>Dibuat: {new Date(customer.created_at).toLocaleString('id-ID')}</span>
        {/if}
      </div>
    </div>
  {/if}

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose}>
      Tutup
    </M3Button>
    {#if customer}
      <M3Button
        variant="tonal"
        onclick={() => {
          onClose();
          window.location.hash = '#finance-customer-contract';
        }}
      >
        <FileText class="w-4 h-4 mr-1" />
        <span>Kontrak PKS</span>
      </M3Button>
      <M3Button variant="filled" onclick={() => { onClose(); onEdit(customer!); }}>
        <Edit class="w-4 h-4 mr-1" />
        <span>Edit Penjamin</span>
      </M3Button>
    {/if}
  {/snippet}
</M3Modal>
