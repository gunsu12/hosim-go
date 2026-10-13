<script lang="ts">
  import { Search, Plus, Share2, Phone, MapPin } from '@lucide/svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import M3Select from '../../../components/m3/M3Select.svelte';
  import { alertDialog } from '../../../stores/confirm.svelte';
  import type { ReferalRecord } from '../../../types/master/referal';

  let searchQuery = $state('');
  let selectedType = $state('ALL');

  const typeOptions = [
    { value: 'ALL', label: 'Semua Jenis Faskes' },
    { value: 'PUSKESMAS', label: 'Puskesmas / Faskes 1' },
    { value: 'KLINIK_PRATAMA', label: 'Klinik Pratama' },
    { value: 'RS_TIPE_C', label: 'RS Tipe C' },
    { value: 'DOKTER_PRAKTEK', label: 'Dokter Praktek' }
  ];

  let referals = $state<ReferalRecord[]>([
    {
      id: 'REF-001',
      code: 'FAS-PKM-01',
      name: 'Puskesmas Sukajadi',
      type: 'PUSKESMAS',
      typeLabel: 'Puskesmas / Faskes 1',
      address: 'Jl. Sukajadi No. 124, Bandung',
      phone: '022-2031123',
      email: 'pkm.sukajadi@dinkes.bandung.go.id',
      status: 'ACTIVE'
    },
    {
      id: 'REF-002',
      code: 'FAS-KLN-02',
      name: 'Klinik Pratama Kimia Farma Dago',
      type: 'KLINIK_PRATAMA',
      typeLabel: 'Klinik Pratama Swasta',
      address: 'Jl. Ir. H. Juanda No. 88, Bandung',
      phone: '022-2504455',
      email: 'klinik.dago@kimiafarma.co.id',
      status: 'ACTIVE'
    },
    {
      id: 'REF-003',
      code: 'FAS-RSC-01',
      name: 'RS Tipe C AMC Cileunyi',
      type: 'RS_TIPE_C',
      typeLabel: 'Rumah Sakit Tipe C',
      address: 'Jl. Raya Cileunyi No. 1, Kab. Bandung',
      phone: '022-7798999',
      email: 'info@rs-amc.co.id',
      status: 'ACTIVE'
    },
    {
      id: 'REF-004',
      code: 'FAS-PKM-02',
      name: 'Puskesmas Garuda',
      type: 'PUSKESMAS',
      typeLabel: 'Puskesmas / Faskes 1',
      address: 'Jl. Dadali No. 8, Bandung',
      phone: '022-6012433',
      email: 'pkm.garuda@dinkes.bandung.go.id',
      status: 'ACTIVE'
    },
    {
      id: 'REF-005',
      code: 'FAS-DOC-01',
      name: 'Praktek Mandiri dr. Anwar Santoso, Sp.A',
      type: 'DOKTER_PRAKTEK',
      typeLabel: 'Dokter Praktek Perorangan',
      address: 'Jl. R.E. Martadinata No. 45, Bandung',
      phone: '0812-2134-5566',
      email: 'dr.anwar@gmail.com',
      status: 'ACTIVE'
    }
  ]);

  let filteredReferals = $derived(
    referals.filter(r => {
      const matchSearch =
        r.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        r.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
        r.address.toLowerCase().includes(searchQuery.toLowerCase()) ||
        r.phone.includes(searchQuery);
      const matchType = selectedType === 'ALL' || r.type === selectedType;
      return matchSearch && matchType;
    })
  );
</script>

<div class="flex flex-col gap-4">
  <!-- Header Domain -->
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <div class="flex items-center gap-2">
        <h2 class="text-lg font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] tracking-tight">Master Asal & Tujuan Rujukan (Referral)</h2>
        <span class="text-xs font-mono font-medium px-2 py-0.5 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa]">
          domain: internal/master/referal
        </span>
      </div>
      <p class="text-xs text-[#444746] dark:text-[#8e918f] mt-0.5">Registrasi fasilitas kesehatan pengirim rujukan BPJS (Faskes 1/2), RS Jejaring, dan Dokter Praktek Perorangan.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button variant="filled" onclick={() => alertDialog({ title: 'Tambah Fasilitas Rujukan', message: 'Fitur pendaftaran faskes rujukan baru sedang dalam pengembangan.', variant: 'info' })}>
        <Plus class="w-4 h-4" />
        <span>Tambah Faskes Rujukan</span>
      </M3Button>
    </div>
  </div>

  <!-- Search & Filter Controls -->
  <div class="flex flex-wrap items-center gap-3 bg-[#f8fafd] dark:bg-[#28292a] p-3 rounded-2xl border border-[#e1e5ea] dark:border-[#353738]">
    <div class="flex-1 min-w-[240px] relative">
      <Search class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f]" />
      <input
        type="text"
        bind:value={searchQuery}
        placeholder="Cari nama faskes, kode rujukan, alamat, atau nomor telepon..."
        class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden focus:ring-2 focus:ring-[#0b57d0] dark:focus:ring-[#a8c7fa]"
      />
    </div>

    <div class="w-56">
      <M3Select
        value={selectedType}
        onchange={(val) => selectedType = val}
        options={typeOptions}
      />
    </div>
  </div>

  <!-- Data Table -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] overflow-hidden shadow-xs">
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-[#f0f4f9] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5] font-semibold">
            <th class="py-3 px-4">Kode Faskes</th>
            <th class="py-3 px-4">Nama Faskes / Mitra Rujukan</th>
            <th class="py-3 px-4">Kategori Faskes</th>
            <th class="py-3 px-4">Alamat & Wilayah</th>
            <th class="py-3 px-4">Kontak & Email</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738] text-[#1f1f1f] dark:text-[#e3e3e3]">
          {#each filteredReferals as ref (ref.id)}
            <tr class="hover:bg-[#f8fafd] dark:hover:bg-[#28292a] transition-colors">
              <td class="py-3.5 px-4 font-mono font-semibold text-[#0b57d0] dark:text-[#a8c7fa]">
                {ref.code}
              </td>
              <td class="py-3.5 px-4 font-medium">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-[#e8f0fe] dark:bg-[#0b57d0]/20 text-[#0b57d0] dark:text-[#a8c7fa] flex items-center justify-center font-bold text-xs shrink-0">
                    <Share2 class="w-4 h-4" />
                  </div>
                  <div>
                    <div class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">{ref.name}</div>
                    <div class="text-[11px] text-[#747775] dark:text-[#8e918f] font-mono">ID: {ref.id}</div>
                  </div>
                </div>
              </td>
              <td class="py-3.5 px-4">
                <span class="px-2.5 py-1 rounded-full text-[11px] font-semibold bg-[#e9eef6] dark:bg-[#28292a] text-[#041e49] dark:text-[#c4c7c5] border border-transparent dark:border-[#353738]">
                  {ref.typeLabel}
                </span>
              </td>
              <td class="py-3.5 px-4 text-[#444746] dark:text-[#8e918f] max-w-[200px] truncate">
                <div class="flex items-center gap-1.5">
                  <MapPin class="w-3.5 h-3.5 text-[#747775] dark:text-[#8e918f] shrink-0" />
                  <span class="truncate">{ref.address}</span>
                </div>
              </td>
              <td class="py-3.5 px-4 text-[#444746] dark:text-[#8e918f]">
                <div class="flex items-center gap-1.5">
                  <Phone class="w-3.5 h-3.5 text-[#747775] dark:text-[#8e918f]" />
                  <span>{ref.phone}</span>
                </div>
                <div class="text-[11px] text-[#747775] dark:text-[#8e918f] mt-0.5">{ref.email}</div>
              </td>
              <td class="py-3.5 px-4">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#e6f4ea] dark:bg-emerald-950/40 text-[#137333] dark:text-emerald-300 border border-[#ceead6] dark:border-emerald-800">
                  Aktif
                </span>
              </td>
              <td class="py-3.5 px-4 text-right">
                <button
                  type="button"
                  class="text-[#0b57d0] dark:text-[#a8c7fa] hover:underline font-medium text-xs cursor-pointer"
                  onclick={() => alertDialog({ title: 'Edit Faskes', message: `Fitur pengubahan faskes "${ref.name}" sedang dalam pengembangan.`, variant: 'info' })}
                >
                  Edit
                </button>
              </td>
            </tr>
          {/each}
          {#if filteredReferals.length === 0}
            <tr>
              <td colspan="7" class="py-8 text-center text-xs text-[#747775] dark:text-[#8e918f]">
                Tidak ada data faskes rujukan yang sesuai dengan pencarian.
              </td>
            </tr>
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
