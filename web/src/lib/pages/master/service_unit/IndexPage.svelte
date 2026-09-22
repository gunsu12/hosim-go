<script lang="ts">
  import { Search, Plus, Stethoscope, MapPin } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import type { ServiceUnitRecord } from '../../../types/master/service_unit';

  let searchQuery = $state('');

  let serviceUnits = $state<ServiceUnitRecord[]>([
    {
      id: 'SU-001',
      code: 'POLI-BEDAH',
      name: 'Poliklinik Bedah Umum',
      departmentName: 'Instalasi Rawat Jalan',
      phone: 'Ext. 101',
      ihsLocationId: 'LOC-3273-0101',
      isRegistrationTarget: true,
      isActive: true
    },
    {
      id: 'SU-002',
      code: 'POLI-INTERNA',
      name: 'Poliklinik Penyakit Dalam',
      departmentName: 'Instalasi Rawat Jalan',
      phone: 'Ext. 102',
      ihsLocationId: 'LOC-3273-0102',
      isRegistrationTarget: true,
      isActive: true
    },
    {
      id: 'SU-003',
      code: 'POLI-GIGI',
      name: 'Poliklinik Gigi & Mulut',
      departmentName: 'Instalasi Rawat Jalan',
      phone: 'Ext. 103',
      ihsLocationId: 'LOC-3273-0103',
      isRegistrationTarget: true,
      isActive: true
    },
    {
      id: 'SU-004',
      code: 'UNIT-OK-CENTRAL',
      name: 'Kamar Bedah Sentral',
      departmentName: 'Instalasi Bedah Sentral',
      phone: 'Ext. 201',
      ihsLocationId: 'LOC-3273-0201',
      isRegistrationTarget: false,
      isActive: true
    },
    {
      id: 'SU-005',
      code: 'UNIT-IGD-TRIAGE',
      name: 'Zona Triage & Resusitasi',
      departmentName: 'Instalasi Gawat Darurat',
      phone: 'Ext. 118',
      ihsLocationId: 'LOC-3273-0118',
      isRegistrationTarget: true,
      isActive: true
    }
  ]);

  let filtered = $derived(
    serviceUnits.filter(u =>
      u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
      u.departmentName.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Unit Layanan (Service Unit)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          domain: internal/master/service_unit
        </span>
      </div>
      <p class="text-xs text-[#444746]">Daftar poliklinik, unit pemeriksaan, dan tujuan pendaftaran pasien rawat jalan/inap.</p>
    </div>

    <M3Button variant="filled" onclick={() => alert('Tambah Unit Layanan Baru')}>
      <Plus class="w-4 h-4" />
      <span>Tambah Unit Layanan</span>
    </M3Button>
  </div>

  <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama poliklinik atau instalasi..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all"
      />
    </div>
    <span class="text-xs text-[#444746]">Total: <strong>{filtered.length}</strong> unit</span>
  </div>

  <div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746]">
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Unit Layanan / Poli</th>
            <th class="py-3 px-4 font-semibold">Induk Departemen</th>
            <th class="py-3 px-4 font-semibold">Telepon / Ext</th>
            <th class="py-3 px-4 font-semibold">SATUSEHAT Location ID</th>
            <th class="py-3 px-4 font-semibold text-center">Tujuan Pendaftaran</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#each filtered as u}
            <tr class="hover:bg-[#f8fafd] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0]">{u.code}</td>
              <td class="py-3.5 px-4 font-semibold text-[#1f1f1f]">{u.name}</td>
              <td class="py-3.5 px-4 text-[#444746]">{u.departmentName}</td>
              <td class="py-3.5 px-4 font-mono text-[#444746]">{u.phone}</td>
              <td class="py-3.5 px-4 font-mono text-[#0b57d0]">{u.ihsLocationId}</td>
              <td class="py-3.5 px-4 text-center">
                {#if u.isRegistrationTarget}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-blue-100 text-blue-800">
                    Bisa Didaftar
                  </span>
                {:else}
                  <span class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-gray-100 text-gray-600">
                    Internal
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
                  onclick={() => alert(`Edit Unit: ${u.name}`)}
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
