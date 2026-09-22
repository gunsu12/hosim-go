<script lang="ts">
  import { Search, Plus, Filter, UserCheck, Shield, Phone, MapPin, Calendar, FileText } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
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
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Data Pasien</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          domain: internal/master/patient
        </span>
      </div>
      <p class="text-xs text-[#444746]">Pengelolaan database pasien rumah sakit, identitas NIK/KTP, data kontak, dan riwayat penjamin.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="filled" onclick={() => alert('Fitur Tambah Pasien Baru (Modal / Form)')}>
        <Plus class="w-4 h-4" />
        <span>Tambah Pasien</span>
      </M3Button>
    </div>
  </div>

  <!-- Search & Filter Toolbar -->
  <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama pasien, No. RM, atau NIK..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all"
      />
    </div>

    <div class="flex items-center gap-2 text-xs text-[#444746]">
      <span>Total: <strong>{filteredPatients.length}</strong> pasien</span>
    </div>
  </div>

  <!-- Tabel Data Pasien Google Drive Style -->
  <div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746]">
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
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#each filteredPatients as p}
            <tr class="hover:bg-[#f8fafd] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0]">{p.mrn}</td>
              <td class="py-3.5 px-4 font-semibold">{p.name}</td>
              <td class="py-3.5 px-4 font-mono text-[#444746]">{p.nik}</td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-md font-semibold text-[11px] {p.gender === 'L' ? 'bg-blue-100 text-blue-800' : 'bg-pink-100 text-pink-800'}">
                  {p.gender}
                </span>
              </td>
              <td class="py-3.5 px-4 text-[#444746]">{p.birthDate}</td>
              <td class="py-3.5 px-4 font-semibold">{p.bloodType}</td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-full text-[11px] font-medium bg-[#e6f4ea] text-[#137333]">
                  {p.payer}
                </span>
              </td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                  {p.status}
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  onclick={() => alert(`Buka Rekam Medis: ${p.name} (${p.mrn})`)}
                  class="text-xs text-[#0b57d0] hover:underline font-semibold cursor-pointer"
                >
                  Detail RME
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>
