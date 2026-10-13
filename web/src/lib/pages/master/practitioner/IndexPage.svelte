<script lang="ts">
  import { Search, Plus } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { alertDialog } from '../../../stores/confirm.svelte';
  import type { PractitionerRecord, SpecialtyRecord, ProfessionRecord } from '../../../types/master/practitioner';

  let searchQuery = $state('');

  function getSpecialtyText(s?: SpecialtyRecord | string): string {
    if (!s) return '';
    return typeof s === 'string' ? s : s.name;
  }

  function getProfessionText(p?: ProfessionRecord | string): string {
    if (!p) return '';
    return typeof p === 'string' ? p : p.name;
  }

  let practitioners = $state<PractitionerRecord[]>([
    {
      id: 'DOC-001',
      name: 'dr. Hendra Wijaya, Sp.B',
      sip: 'SIP.446/012/DINKES/2024',
      str: 'STR-327310029381',
      specialty: 'Spesialis Bedah Umum',
      profession: 'Dokter Spesialis',
      phone: '0812-8877-6655',
      ihsId: '10001234567',
      is_active: true
    },
    {
      id: 'DOC-002',
      name: 'dr. Amanda Putri, Sp.PD',
      sip: 'SIP.446/045/DINKES/2023',
      str: 'STR-327310088219',
      specialty: 'Spesialis Penyakit Dalam',
      profession: 'Dokter Spesialis',
      phone: '0813-1122-3344',
      ihsId: '10007654321',
      is_active: true
    },
    {
      id: 'DOC-003',
      name: 'drg. Kevin Pratama, Sp.KG',
      sip: 'SIP.446/089/DINKES/2024',
      str: 'STR-327310055412',
      specialty: 'Konservasi Gigi (Endodontik)',
      profession: 'Dokter Gigi Spesialis',
      phone: '0857-4433-2211',
      ihsId: '10009988776',
      is_active: true
    },
    {
      id: 'NRS-001',
      name: 'Ns. Ratna Sari, S.Kep',
      sip: 'SIPK.446/102/2022',
      str: 'STR-327320011456',
      specialty: 'Perawat Kamar Bedah',
      profession: 'Perawat (Nurse)',
      phone: '0877-9988-1122',
      ihsId: '10003344556',
      is_active: true
    }
  ]);

  let filtered = $derived(
    practitioners.filter(p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      getSpecialtyText(p.specialty).toLowerCase().includes(searchQuery.toLowerCase()) ||
      (p.sip || '').toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">Master Tenaga Medis (Practitioner)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa]">
          domain: internal/master/practitioner
        </span>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#8e918f] mt-0.5">Registrasi dokter spesialis, dokter umum, perawat, nomor izin praktek (SIP/STR), dan SATUSEHAT IHS ID.</p>
    </div>

    <M3Button variant="filled" onclick={() => alertDialog({ title: 'Tambah Tenaga Medis', message: 'Fitur registrasi tenaga medis baru (dokter/perawat) sedang dalam pengembangan.', variant: 'info' })}>
      <Plus class="w-4 h-4" />
      <span>Tambah Tenaga Medis</span>
    </M3Button>
  </div>

  <div class="p-3 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] dark:text-[#8e918f] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama dokter, spesialisasi, atau SIP..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
    <span class="text-xs text-[#444746] dark:text-[#8e918f]">Total: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]">{filtered.length}</strong> nakes</span>
  </div>

  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold">Nama Tenaga Medis</th>
            <th class="py-3 px-4 font-semibold">Profesi & Spesialisasi</th>
            <th class="py-3 px-4 font-semibold">Nomor SIP</th>
            <th class="py-3 px-4 font-semibold">Nomor STR</th>
            <th class="py-3 px-4 font-semibold">IHS Practitioner ID</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#each filtered as p}
            <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a] transition-colors">
              <td class="py-3.5 px-4">
                <div class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{p.name}</div>
                <div class="text-[11px] text-[#444746] dark:text-[#8e918f] font-mono">{p.phone}</div>
              </td>
              <td class="py-3.5 px-4">
                <div class="font-semibold text-[#0b57d0] dark:text-[#a8c7fa]">{getSpecialtyText(p.specialty)}</div>
                <div class="text-[11px] text-[#444746] dark:text-[#8e918f]">{getProfessionText(p.profession)}</div>
              </td>
              <td class="py-3.5 px-4 font-mono text-[#1f1f1f] dark:text-[#e3e3e3]">{p.sip || '-'}</td>
              <td class="py-3.5 px-4 font-mono text-[#444746] dark:text-[#8e918f]">{p.str || '-'}</td>
              <td class="py-3.5 px-4">
                <span class="font-mono text-[11px] px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa] font-semibold">
                  {p.ihsId || p.ihs_practitioner_id || '-'}
                </span>
              </td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800">
                  AKTIF
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  onclick={() => alertDialog({ title: 'Jadwal Praktek', message: `Jadwal praktek dokter "${p.name}" (${getSpecialtyText(p.specialty)}) sedang disiapkan.`, variant: 'info' })}
                  class="text-xs text-[#0b57d0] dark:text-[#a8c7fa] hover:underline font-semibold cursor-pointer"
                >
                  Jadwal Praktek
                </button>
              </td>
            </tr>
          {/each}
          {#if filtered.length === 0}
            <tr>
              <td colspan="7" class="py-8 text-center text-xs text-[#747775] dark:text-[#8e918f]">
                Tidak ada data nakes yang cocok dengan pencarian.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
