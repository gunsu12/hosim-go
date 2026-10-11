<script lang="ts">
  import { X, ShieldCheck, Building2, Phone, Mail, Globe, MapPin, FileText, CheckCircle2, Calendar, CreditCard, Edit } from '@lucide/svelte';
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

{#if isOpen && customer}
  <div
    class="fixed inset-0 z-50 bg-black/40 backdrop-blur-xs flex items-center justify-center p-4"
    role="dialog"
    aria-modal="true"
  >
    <div
      class="bg-white rounded-3xl border border-[#e1e5ea] shadow-2xl w-full max-w-xl max-h-[90vh] flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-150"
    >
      <!-- Header -->
      <div class="px-6 py-4 border-b border-[#e1e5ea] flex items-center justify-between bg-[#f8fafd]">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-2xl bg-[#0b57d0]/10 text-[#0b57d0] flex items-center justify-center font-bold">
            <Building2 class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-base font-bold text-[#1f1f1f]">
              Detail Debitur / Penjamin
            </h3>
            <span class="font-mono text-xs text-[#0b57d0] font-semibold">{customer.code}</span>
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

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6 space-y-4 text-xs">
        <!-- Title & Status Banner -->
        <div class="p-4 rounded-2xl bg-[#f0f4f9] border border-[#e1e5ea] flex flex-wrap items-center justify-between gap-3">
          <div>
            <h4 class="text-sm font-bold text-[#1f1f1f]">{customer.name}</h4>
            <p class="text-xs text-[#444746] mt-0.5">
              Kategori: <strong>{customer.customer_type?.name || 'Umum / Lainnya'}</strong>
            </p>
          </div>
          <div class="flex items-center gap-2">
            {#if customer.is_active}
              <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                AKTIF
              </span>
            {:else}
              <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-gray-200 text-gray-700">
                NON-AKTIF
              </span>
            {/if}
          </div>
        </div>

        <!-- Info Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <div class="p-3 rounded-xl border border-[#e1e5ea] bg-white">
            <span class="text-[11px] text-[#747775] block">Person In Charge (PIC)</span>
            <span class="font-semibold text-[#1f1f1f] block mt-0.5">{customer.contact_person || '-'}</span>
          </div>

          <div class="p-3 rounded-xl border border-[#e1e5ea] bg-white">
            <span class="text-[11px] text-[#747775] block">Telepon / Call Center</span>
            <span class="font-mono font-semibold text-[#1f1f1f] block mt-0.5">{customer.phone || '-'}</span>
          </div>

          <div class="p-3 rounded-xl border border-[#e1e5ea] bg-white">
            <span class="text-[11px] text-[#747775] block">Email Korespondensi</span>
            <span class="font-semibold text-[#1f1f1f] block mt-0.5">{customer.email || '-'}</span>
          </div>

          <div class="p-3 rounded-xl border border-[#e1e5ea] bg-white">
            <span class="text-[11px] text-[#747775] block">Situs Resmi / Portal</span>
            {#if customer.website && (customer.website.startsWith('http://') || customer.website.startsWith('https://'))}
              <a href={customer.website} target="_blank" rel="noopener noreferrer" class="text-[#0b57d0] hover:underline block mt-0.5 font-semibold truncate">
                {customer.website}
              </a>
            {:else if customer.website}
              <span class="text-[#1f1f1f] block mt-0.5 font-mono truncate">{customer.website}</span>
            {:else}
              <span class="text-[#747775] block mt-0.5">-</span>
            {/if}
          </div>
        </div>

        <!-- Alamat -->
        <div class="p-3.5 rounded-xl border border-[#e1e5ea] bg-white">
          <span class="text-[11px] text-[#747775] block font-semibold mb-1">Alamat Kantor / Instansi</span>
          <p class="text-[#1f1f1f]">{customer.address || 'Belum ada data alamat tercatat.'}</p>
        </div>

        <!-- Keterangan -->
        <div class="p-3.5 rounded-xl border border-[#e1e5ea] bg-white">
          <span class="text-[11px] text-[#747775] block font-semibold mb-1">Catatan Tambahan & PKS</span>
          <p class="text-[#444746] whitespace-pre-wrap">{customer.description || 'Tidak ada catatan khusus.'}</p>
        </div>

        <!-- Aturan Peserta -->
        <div class="p-3.5 rounded-xl border border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-between">
          <div>
            <span class="font-semibold text-[#1f1f1f] block">Verifikasi Kartu / SEP Peserta</span>
            <span class="text-[11px] text-[#747775]">Wajib diverifikasi saat pasien mendaftar di unit layanan</span>
          </div>
          <div>
            {#if customer.require_card}
              <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-amber-100 text-amber-800">
                Wajib Kartu / SEP
              </span>
            {:else}
              <span class="px-2.5 py-1 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600">
                Tanpa Kartu
              </span>
            {/if}
          </div>
        </div>

        <!-- Audit Meta -->
        <div class="pt-2 text-[11px] text-[#747775] flex flex-wrap justify-between gap-2 border-t border-[#e1e5ea]">
          <span>ID Sistem: <strong class="font-mono">{customer.id}</strong></span>
          {#if customer.created_at}
            <span>Dibuat: {new Date(customer.created_at).toLocaleString('id-ID')}</span>
          {/if}
        </div>
      </div>

      <!-- Footer -->
      <div class="px-6 py-3.5 border-t border-[#e1e5ea] bg-[#f8fafd] flex items-center justify-end gap-3">
        <M3Button variant="text" onclick={onClose}>
          Tutup
        </M3Button>
        <M3Button variant="filled" onclick={() => { onClose(); onEdit(customer); }}>
          <Edit class="w-4 h-4 mr-1" />
          <span>Edit Penjamin</span>
        </M3Button>
      </div>
    </div>
  </div>
{/if}
