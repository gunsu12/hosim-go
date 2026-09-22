<script lang="ts">
  import { Search, Plus, ShieldCheck, CreditCard } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import type { PayerRecord } from '../../../types/master/payer';

  let searchQuery = $state('');

  let payers = $state<PayerRecord[]>([
    {
      id: 'PAY-001',
      code: 'BPJS-PBI',
      name: 'BPJS Kesehatan - Penerima Bantuan Iuran (PBI)',
      typeName: 'Pemerintah / JKN',
      phone: '165',
      requireCard: true,
      isActive: true
    },
    {
      id: 'PAY-002',
      code: 'BPJS-NONPBI',
      name: 'BPJS Kesehatan - Pekerja Penerima Upah / Mandiri',
      typeName: 'Pemerintah / JKN',
      phone: '165',
      requireCard: true,
      isActive: true
    },
    {
      id: 'PAY-003',
      code: 'ASR-MANDIRI',
      name: 'Asuransi Mandiri Inhealth',
      typeName: 'Asuransi Swasta',
      phone: '021-5250988',
      requireCard: true,
      isActive: true
    },
    {
      id: 'PAY-004',
      code: 'ASR-PRUDENTIAL',
      name: 'Prudential Life Assurance (PRUPrime Healthcare)',
      typeName: 'Asuransi Swasta',
      phone: '1500085',
      requireCard: true,
      isActive: true
    },
    {
      id: 'PAY-005',
      code: 'UMUM',
      name: 'Pasien Umum / Bayar Pribadi (Out of Pocket)',
      typeName: 'Pribadi / Cash',
      phone: '-',
      requireCard: false,
      isActive: true
    }
  ]);

  let filtered = $derived(
    payers.filter(p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.typeName.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Penjamin / Asuransi (Payer)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          domain: internal/master/payer
        </span>
      </div>
      <p class="text-xs text-[#444746]">Konfigurasi penjamin pembiayaan pasien (BPJS Kesehatan VClaim, Asuransi Swasta, Perusahaan, & Umum).</p>
    </div>

    <M3Button variant="filled" onclick={() => alert('Tambah Penjamin Baru')}>
      <Plus class="w-4 h-4" />
      <span>Tambah Penjamin</span>
    </M3Button>
  </div>

  <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama penjamin atau tipe asuransi..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all"
      />
    </div>
    <span class="text-xs text-[#444746]">Total: <strong>{filtered.length}</strong> penjamin</span>
  </div>

  <div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746]">
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Penjamin / Asuransi</th>
            <th class="py-3 px-4 font-semibold">Kategori Penjamin</th>
            <th class="py-3 px-4 font-semibold">Call Center</th>
            <th class="py-3 px-4 font-semibold text-center">Wajib Kartu</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#each filtered as p}
            <tr class="hover:bg-[#f8fafd] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0]">{p.code}</td>
              <td class="py-3.5 px-4 font-semibold text-[#1f1f1f]">{p.name}</td>
              <td class="py-3.5 px-4 text-[#444746]">{p.typeName}</td>
              <td class="py-3.5 px-4 font-mono text-[#444746]">{p.phone}</td>
              <td class="py-3.5 px-4 text-center">
                {#if p.requireCard}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-100 text-amber-800">
                    Wajib Kartu/SEP
                  </span>
                {:else}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600">
                    Tanpa Kartu
                  </span>
                {/if}
              </td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                  AKTIF
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  onclick={() => alert(`Edit Penjamin: ${p.name}`)}
                  class="text-xs text-[#0b57d0] hover:underline font-semibold cursor-pointer"
                >
                  Edit
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>
