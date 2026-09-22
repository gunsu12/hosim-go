<script lang="ts">
  import { Search, Plus, BadgePercent, Tag, FileText, CheckCircle2, XCircle } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import type { TariffClassRecord } from '../../../types/master/tariff_class';

  let searchQuery = $state('');

  let tariffClasses = $state<TariffClassRecord[]>([
    {
      id: 'TC-001',
      code: 'VVIP',
      name: 'Kelas VVIP / President Suite',
      description: 'Ruang perawatan eksklusif dengan fasilitas penunjang lengkap dan perawat privat',
      multiplier: '2.5x Standar',
      coverageType: 'Pribadi / Asuransi Komersial',
      isActive: true
    },
    {
      id: 'TC-002',
      code: 'VIP',
      name: 'Kelas VIP',
      description: 'Satu pasien per kamar dengan kamar mandi dalam, sofa tunggu keluarga, dan AC',
      multiplier: '1.8x Standar',
      coverageType: 'Pribadi / Asuransi Komersial',
      isActive: true
    },
    {
      id: 'TC-003',
      code: 'KLS-1',
      name: 'Kelas 1 (KRIS Utama)',
      description: 'Maksimal 2 tempat tidur per ruangan sesuai standar Kemenkes',
      multiplier: '1.2x Standar',
      coverageType: 'BPJS / Umum',
      isActive: true
    },
    {
      id: 'TC-004',
      code: 'KLS-2',
      name: 'Kelas 2 (Standar)',
      description: 'Maksimal 3-4 tempat tidur per ruangan, fasilitas berstandar KRIS',
      multiplier: '1.0x Standar (Base)',
      coverageType: 'BPJS / Umum',
      isActive: true
    },
    {
      id: 'TC-005',
      code: 'KLS-3',
      name: 'Kelas 3 (Subsidi BPJS)',
      description: 'Kamar perawatan subsidi pemerintah untuk peserta BPJS PBI dan mandiri kelas 3',
      multiplier: '0.8x Standar',
      coverageType: 'BPJS Kesehatan PBI',
      isActive: true
    },
    {
      id: 'TC-006',
      code: 'NON-CLASS',
      name: 'Non-Kelas (IGD / ICU / PICU / NICU)',
      description: 'Tarif unit intensif & gawat darurat khusus tanpa pembagian strata kamar',
      multiplier: 'Sesuai Paket INA-CBGs',
      coverageType: 'Semua Penjamin',
      isActive: true
    }
  ]);

  let filteredTariffClasses = $derived(
    tariffClasses.filter(t =>
      t.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.description.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="flex flex-col gap-4">
  <!-- Header Domain -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Master Kelas Tarif (Tariff Class)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
          domain: internal/master/tariff_class
        </span>
      </div>
      <p class="text-xs text-[#444746]">Klasifikasi kelas tarif akomodasi rawat inap, rawat jalan, dan koefisien pengali tindakan medis (INA-CBGs / Fee for Service).</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="filled" onclick={() => alert('Modal Tambah Kelas Tarif Baru')}>
        <Plus class="w-4 h-4" />
        <span>Tambah Kelas Tarif</span>
      </M3Button>
    </div>
  </div>

  <!-- Search & Filter Controls -->
  <div class="flex items-center gap-3 bg-[#f8fafd] p-3 rounded-2xl border border-[#e1e5ea]">
    <div class="flex-1 relative">
      <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#747775]" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari kode kelas tarif (VVIP, VIP, KLS-1) atau deskripsi..."
        class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-white border border-[#e1e5ea] text-[#1f1f1f] focus:outline-none focus:ring-2 focus:ring-[#0b57d0]"
      />
    </div>
  </div>

  <!-- Data Table -->
  <div class="bg-white rounded-2xl border border-[#e1e5ea] overflow-hidden shadow-xs">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] border-b border-[#e1e5ea] text-[#444746] font-semibold">
            <th class="py-3 px-4">Kode Tarif</th>
            <th class="py-3 px-4">Nama Kelas Tarif</th>
            <th class="py-3 px-4">Deskripsi / Kriteria Fasilitas</th>
            <th class="py-3 px-4">Koefisien Akomodasi</th>
            <th class="py-3 px-4">Kategori Penjamin</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] text-[#1f1f1f]">
          {#each filteredTariffClasses as tc (tc.id)}
            <tr class="hover:bg-[#f8fafd] transition-colors">
              <td class="py-3.5 px-4 font-mono font-bold text-[#0b57d0]">
                {tc.code}
              </td>
              <td class="py-3.5 px-4 font-semibold text-[#1f1f1f]">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-[#fef7e0] text-[#b06000] flex items-center justify-center font-bold text-xs shrink-0">
                    <Tag class="w-4 h-4" />
                  </div>
                  <div>
                    <div>{tc.name}</div>
                    <div class="text-[11px] text-[#747775] font-mono">ID: {tc.id}</div>
                  </div>
                </div>
              </td>
              <td class="py-3.5 px-4 text-[#444746] max-w-[280px]">
                {tc.description}
              </td>
              <td class="py-3.5 px-4">
                <span class="px-2.5 py-1 rounded-full text-[11px] font-semibold bg-[#e8f0fe] text-[#0b57d0]">
                  {tc.multiplier}
                </span>
              </td>
              <td class="py-3.5 px-4 text-[#444746]">
                {tc.coverageType}
              </td>
              <td class="py-3.5 px-4">
                {#if tc.isActive}
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#e6f4ea] text-[#137333]">
                    <CheckCircle2 class="w-3 h-3" />
                    <span>Aktif</span>
                  </span>
                {:else}
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#fce8e6] text-[#c5221f]">
                    <XCircle class="w-3 h-3" />
                    <span>Nonaktif</span>
                  </span>
                {/if}
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  class="text-[#0b57d0] hover:underline font-medium text-xs cursor-pointer"
                  onclick={() => alert(`Edit kelas tarif: ${tc.name}`)}
                >
                  Edit
                </button>
              </td>
            </tr>
          {/each}
          {#if filteredTariffClasses.length === 0}
            <tr>
              <td colspan="7" class="py-8 text-center text-xs text-[#747775]">
                Tidak ada data kelas tarif yang sesuai dengan pencarian.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
