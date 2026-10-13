<script lang="ts">
  import { Search, Plus } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { alertDialog } from '../../../stores/confirm.svelte';
  import type { RoomRecord } from '../../../types/master/room';

  let searchQuery = $state('');

  let rooms = $state<RoomRecord[]>([
    {
      id: 'RM-001',
      code: 'OK-01',
      name: 'Kamar Operasi Bedah Mayor 1',
      capacity: 1,
      location: 'Gedung Bedah Lt. 3',
      roomType: 'OPERATION_ROOM',
      unitName: 'Kamar Bedah Sentral'
    },
    {
      id: 'RM-002',
      code: 'VIP-201',
      name: 'Kamar Perawatan VIP Melati 201',
      capacity: 1,
      location: 'Gedung Rawat Inap B Lt. 2',
      roomType: 'INPATIENT_VIP',
      unitName: 'Instalasi Rawat Inap'
    },
    {
      id: 'RM-003',
      code: 'KLS1-305',
      name: 'Bangsal Kelas 1 Dahlia 305 (2 Bed)',
      capacity: 2,
      location: 'Gedung Rawat Inap A Lt. 3',
      roomType: 'INPATIENT_CLASS_1',
      unitName: 'Instalasi Rawat Inap'
    },
    {
      id: 'RM-004',
      code: 'ICU-BED-04',
      name: 'Intensive Care Unit Bed 04 (Ventilator)',
      capacity: 1,
      location: 'Gedung Utama Lt. 2',
      roomType: 'ICU',
      unitName: 'Instalasi Rawat Intensif'
    }
  ]);

  let filtered = $derived(
    rooms.filter(r =>
      r.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      r.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
      r.location.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">Master Ruangan & Bed (Room)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa]">
          domain: internal/master/room
        </span>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#8e918f] mt-0.5">Pengelolaan kamar rawat inap, kapasitas tempat tidur, ruang bedah (OK), dan unit intensif ICU/ICCU.</p>
    </div>

    <M3Button variant="filled" onclick={() => alertDialog({ title: 'Tambah Ruangan', message: 'Fitur registrasi kamar atau ruangan baru sedang dalam pengembangan.', variant: 'info' })}>
      <Plus class="w-4 h-4" />
      <span>Tambah Ruangan</span>
    </M3Button>
  </div>

  <div class="p-3 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-between gap-3 shadow-xs">
    <div class="relative flex-1 max-w-md">
      <Search class="w-4 h-4 text-[#747775] dark:text-[#8e918f] absolute left-3.5 top-1/2 -translate-y-1/2 pointer-events-none" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama ruangan, kode, atau lokasi gedung..."
        class="w-full h-10 pl-10 pr-4 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
    <span class="text-xs text-[#444746] dark:text-[#8e918f]">Total: <strong class="text-[#1f1f1f] dark:text-[#e3e3e3]">{filtered.length}</strong> ruangan</span>
  </div>

  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5]">
            <th class="py-3 px-4 font-semibold">Kode Ruangan</th>
            <th class="py-3 px-4 font-semibold">Nama Ruangan / Bangsal</th>
            <th class="py-3 px-4 font-semibold">Tipe Ruangan</th>
            <th class="py-3 px-4 font-semibold">Lokasi Gedung</th>
            <th class="py-3 px-4 font-semibold text-center">Kapasitas Bed</th>
            <th class="py-3 px-4 font-semibold">Unit Layanan</th>
            <th class="py-3 px-4 font-semibold text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#each filtered as r}
            <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa]">{r.code}</td>
              <td class="py-3.5 px-4 font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{r.name}</td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-md font-medium text-[11px] bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa]">
                  {r.roomType}
                </span>
              </td>
              <td class="py-3.5 px-4 text-[#444746] dark:text-[#8e918f]">{r.location}</td>
              <td class="py-3.5 px-4 text-center font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">{r.capacity} Bed</td>
              <td class="py-3.5 px-4 text-[#444746] dark:text-[#8e918f]">{r.unitName}</td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  onclick={() => alertDialog({ title: 'Status Bed Ruangan', message: `Detail tempat tidur untuk ruangan "${r.name}" (${r.unitName}) sedang disiapkan.`, variant: 'info' })}
                  class="text-xs text-[#0b57d0] dark:text-[#a8c7fa] hover:underline font-semibold cursor-pointer"
                >
                  Detail Bed
                </button>
              </td>
            </tr>
          {/each}
          {#if filtered.length === 0}
            <tr>
              <td colspan="7" class="py-8 text-center text-xs text-[#747775] dark:text-[#8e918f]">
                Tidak ada data ruangan yang cocok dengan pencarian.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
