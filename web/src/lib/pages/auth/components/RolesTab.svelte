<script lang="ts">
  import { Search, Shield, Edit2, Trash2, Loader2 } from '@lucide/svelte';
  import type { RoleRecord } from '../../../types/auth';

  interface Props {
    roles: RoleRecord[];
    isLoading: boolean;
    roleSearch: string;
    onSearchChange: (search: string) => void;
    onOpenEditRole: (role: RoleRecord) => void;
    onOpenDeleteRole: (role: RoleRecord) => void;
  }

  let {
    roles = [],
    isLoading = false,
    roleSearch = '',
    onSearchChange = () => {},
    onOpenEditRole = () => {},
    onOpenDeleteRole = () => {}
  }: Props = $props();

  let filteredRoles = $derived(
    roles.filter((r) => {
      const q = roleSearch.toLowerCase().trim();
      if (!q) return true;
      return (
        r.name.toLowerCase().includes(q) ||
        r.code.toLowerCase().includes(q) ||
        (r.description || '').toLowerCase().includes(q)
      );
    })
  );
</script>

<div class="space-y-4">
  <!-- Search Bar -->
  <div class="flex items-center gap-3">
    <div class="relative flex-1 min-w-[240px]">
      <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f]" />
      <input
        type="text"
        value={roleSearch}
        oninput={(e) => onSearchChange(e.currentTarget.value)}
        placeholder="Cari kode peran atau nama peran..."
        class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] focus:bg-white dark:focus:bg-[#1e1f20] focus:ring-2 focus:ring-[#0b57d0] dark:focus:ring-[#a8c7fa] text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <!-- Roles Table -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] overflow-hidden shadow-xs">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="bg-[#f8fafd] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5] font-semibold">
            <th class="py-3 px-4">Kode Peran</th>
            <th class="py-3 px-4">Nama Peran</th>
            <th class="py-3 px-4">Deskripsi</th>
            <th class="py-3 px-4">Jumlah Wewenang</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
          {#if isLoading}
            <tr>
              <td colspan="5" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                <div class="flex items-center justify-center gap-2">
                  <Loader2 class="w-5 h-5 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
                  <span>Memuat data peran...</span>
                </div>
              </td>
            </tr>
          {:else if filteredRoles.length === 0}
            <tr>
              <td colspan="5" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                Tidak ada peran pengguna yang sesuai dengan pencarian.
              </td>
            </tr>
          {:else}
            {#each filteredRoles as r (r.id)}
              <tr class="hover:bg-[#f8fafd]/80 dark:hover:bg-[#28292a]/50 transition-colors">
                <td class="py-3 px-4 font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa]">
                  {r.code}
                </td>
                <td class="py-3 px-4 font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">
                  {r.name}
                </td>
                <td class="py-3 px-4 text-[#747775] dark:text-[#8e918f] max-w-xs truncate">
                  {r.description || '-'}
                </td>
                <td class="py-3 px-4">
                  <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa]">
                    <Shield class="w-3 h-3" />
                    {r.permissions ? r.permissions.length : 0} Izin Modul
                  </span>
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <button
                      type="button"
                      onclick={() => onOpenEditRole(r)}
                      class="p-1.5 text-[#444746] dark:text-[#c4c7c5] hover:text-[#0b57d0] dark:hover:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] rounded-lg transition-colors cursor-pointer"
                      title="Ubah Konfigurasi Peran"
                    >
                      <Edit2 class="w-4 h-4" />
                    </button>
                    {#if r.code !== 'ADMIN'}
                      <button
                        type="button"
                        onclick={() => onOpenDeleteRole(r)}
                        class="p-1.5 text-[#444746] dark:text-[#c4c7c5] hover:text-[#c5221f] dark:hover:text-rose-400 hover:bg-[#fce8e6] dark:hover:bg-rose-950/50 rounded-lg transition-colors cursor-pointer"
                        title="Hapus Peran"
                      >
                        <Trash2 class="w-4 h-4" />
                      </button>
                    {/if}
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
