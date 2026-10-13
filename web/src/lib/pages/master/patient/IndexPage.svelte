<script lang="ts">
  import { Search, Plus } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { alertDialog } from '../../../stores/confirm.svelte';
  import type { PatientRecord } from '../../../types/master/patient';

  let searchQuery = $state('');

  // Sample data awal mencerminkan schema internal/master/patient
  let patients = $state<PatientRecord[]>([
    {
      id: 1,
      mrn: 'RM-2026-08492',
      nik: '3273151405810003',
      name: 'Budi Santoso',
      gender: 'L',
      birthDate: '1981-05-14',
      phone: '0812-3456-7890',
      bloodType: 'O+',
      payer: 'BPJS Kesehatan (PBI)',
      status: 'ACTIVE'
    },
    {
      id: 2,
      mrn: 'RM-2026-08493',
      nik: '3273152208880004',
      name: 'Siti Rahmawati',
      gender: 'P',
      birthDate: '1988-08-22',
      phone: '0813-9876-5432',
      bloodType: 'A+',
      payer: 'BPJS Kesehatan (Non-PBI)',
      status: 'ACTIVE'
    },
    {
      id: 3,
      mrn: 'RM-2026-08494',
      nik: '3273151011950002',
      name: 'Ahmad Fauzi',
      gender: 'L',
      birthDate: '1995-11-10',
      phone: '0857-1122-3344',
      bloodType: 'B+',
      payer: 'Asuransi Mandiri Inhealth',
      status: 'ACTIVE'
    },
    {
      id: 4,
      mrn: 'RM-2026-08495',
      nik: '3273150502010001',
      name: 'Dewi Lestari',
      gender: 'P',
      birthDate: '2001-02-05',
      phone: '0878-5566-7788',
      bloodType: 'AB+',
      payer: 'Umum / Mandiri',
      status: 'ACTIVE'
    }
  ]);

  let filteredPatients = $derived(
    patients.filter(p =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.mrn.toLowerCase().includes(searchQuery.toLowerCase()) ||
      p.nik.includes(searchQuery)
    )
  );
</script>

<div class="flex flex-col gap-4">
  <!-- Header Domain -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">Master Data Pasien</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa]">
          domain: internal/master/patient
        </span>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#8e918f] mt-0.5">Pengelolaan database pasien rumah sakit, identitas NIK/KTP, data kontak, dan riwayat penjamin.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="filled" onclick={() => alertDialog({ title: 'Tambah Pasien', message: 'Formulir pendaftaran pasien baru sedang dalam pengembangan.', variant: 'info' })}>
        <Plus class="w-4 h-4" />
        <span>Tambah Pasien</span>
      </M3Button>
    </div>
  </div>

  <!-- Search & Filter Toolbar -->
  <div class="p-3 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] dark:text-[#8e918f] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama pasien, No. RM, atau NIK..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="flex items-center gap-2 text-xs text-[#444746] dark:text-[#8e918f]">
      <span>Total: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]">{filteredPatients.length}</strong> pasien</span>
    </div>
  </div>

  <!-- Tabel Data Pasien -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold">No. RM</th>
            <th class="py-3 px-4 font-semibold">Nama Pasien</th>
            <th class="py-3 px-4 font-semibold">NIK</th>
            <th class="py-3 px-4 font-semibold">L/P</th>
            <th class="py-3 px-4 font-semibold">Tgl. Lahir</th>
            <th class="py-3 px-4 font-semibold">Gol. Darah</th>
            <th class="py-3 px-4 font-semibold">Penjamin</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#each filteredPatients as p}
            <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa]">{p.mrn}</td>
              <td class="py-3.5 px-4 font-semibold">{p.name}</td>
              <td class="py-3.5 px-4 font-mono text-[#444746] dark:text-[#8e918f]">{p.nik}</td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-md font-semibold text-[11px] {p.gender === 'L' ? 'bg-blue-100 dark:bg-blue-950/40 text-blue-800 dark:text-blue-300' : 'bg-pink-100 dark:bg-pink-950/40 text-pink-800 dark:text-pink-300'}">
                  {p.gender}
                </span>
              </td>
              <td class="py-3.5 px-4 text-[#444746] dark:text-[#8e918f]">{p.birthDate}</td>
              <td class="py-3.5 px-4 font-semibold">{p.bloodType}</td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-full text-[11px] font-medium bg-[#e6f4ea] dark:bg-emerald-950/40 text-[#137333] dark:text-emerald-300 border border-[#ceead6] dark:border-emerald-800">
                  {p.payer}
                </span>
              </td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 dark:bg-emerald-950/60 text-emerald-800 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800">
                  {p.status}
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  onclick={() => alertDialog({ title: 'Rekam Medis Elektronik', message: `Membuka Rekam Medis Pasien: ${p.name} (No. RM: ${p.mrn}). Fitur integrasi berkas RME sedang disiapkan.`, variant: 'info' })}
                  class="text-xs text-[#0b57d0] dark:text-[#a8c7fa] hover:underline font-semibold cursor-pointer"
                >
                  Detail RME
                </button>
              </td>
            </tr>
          {/each}
          {#if filteredPatients.length === 0}
            <tr>
              <td colspan="9" class="py-8 text-center text-xs text-[#747775] dark:text-[#8e918f]">
                Tidak ada data pasien yang cocok dengan pencarian.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
