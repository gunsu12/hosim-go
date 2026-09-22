<script lang="ts">
  import { Search, Plus, Building2, Layers } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import type { DepartementRecord } from '../../../types/master/departement';

  let searchQuery = $state('');

  let departments = $state<DepartementRecord[]>([
    {
      id: 'DEP-001',
      code: 'INST-RJ',
      name: 'Instalasi Rawat Jalan (Poliklinik)',
      type: 'CLINICAL',
      phone: '022-710001',
      ihsOrgId: 'ORG-10009823',
      isActive: true
    },
    {
      id: 'DEP-002',
      code: 'INST-RI',
      name: 'Instalasi Rawat Inap',
      type: 'CLINICAL',
      phone: '022-710002',
      ihsOrgId: 'ORG-10009824',
      isActive: true
    },
    {
      id: 'DEP-003',
      code: 'INST-BEDAH',
      name: 'Instalasi Bedah Sentral (IBS / OK)',
      type: 'SURGICAL',
      phone: '022-710003',
      ihsOrgId: 'ORG-10009825',
      isActive: true
    },
    {
      id: 'DEP-004',
      code: 'INST-IGD',
      name: 'Instalasi Gawat Darurat (IGD 24 Jam)',
      type: 'EMERGENCY',
      phone: '022-710004',
      ihsOrgId: 'ORG-10009826',
      isActive: true
    },
    {
      id: 'DEP-005',
      code: 'INST-LAB',
      name: 'Instalasi Laboratorium Patologi',
      type: 'DIAGNOSTIC',
      phone: '022-710005',
      ihsOrgId: 'ORG-10009827',
      isActive: true
    }
  ]);

  let filtered = $derived(
    departments.filter(d =>
      d.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      d.code.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Departemen / Instalasi</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          domain: internal/master/departement
        </span>
      </div>
      <p class="text-xs text-[#444746]">Manajemen instalasi utama rumah sakit (Rawat Jalan, Rawat Inap, Kamar Operasi, IGD, Lab).</p>
    </div>

    <M3Button variant="filled" onclick={() => alert('Tambah Departemen Baru')}>
      <Plus class="w-4 h-4" />
      <span>Tambah Departemen</span>
    </M3Button>
  </div>

  <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama instalasi atau kode departemen..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] border border-transparent focus:border-[#0b57d0] focus:bg-white text-xs text-[#1f1f1f] focus:outline-none transition-all"
      />
    </div>
    <span class="text-xs text-[#444746]">Total: <strong>{filtered.length}</strong> instalasi</span>
  </div>

  <div class="bg-white rounded-2xl border border-[#e1e5ea] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746]">
            <th class="py-3 px-4 font-semibold">Kode</th>
            <th class="py-3 px-4 font-semibold">Nama Instalasi / Departemen</th>
            <th class="py-3 px-4 font-semibold">Kategori Layanan</th>
            <th class="py-3 px-4 font-semibold">Telepon Internal</th>
            <th class="py-3 px-4 font-semibold">SATUSEHAT Org ID</th>
            <th class="py-3 px-4 font-semibold text-center">Status</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#each filtered as d}
            <tr class="hover:bg-[#f8fafd] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0]">{d.code}</td>
              <td class="py-3.5 px-4 font-semibold text-[#1f1f1f]">{d.name}</td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-md font-medium text-[11px] bg-[#f0f4f9] text-[#444746] border border-[#e1e5ea]">
                  {d.type}
                </span>
              </td>
              <td class="py-3.5 px-4 font-mono text-[#444746]">{d.phone}</td>
              <td class="py-3.5 px-4 font-mono text-[#0b57d0]">{d.ihsOrgId}</td>
              <td class="py-3.5 px-4 text-center">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-100 text-emerald-800">
                  AKTIF
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  onclick={() => alert(`Edit Departemen: ${d.name}`)}
                  class="text-xs text-[#0b57d0] hover:underline font-semibold cursor-pointer"
                >
                  Kelola Unit
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>
