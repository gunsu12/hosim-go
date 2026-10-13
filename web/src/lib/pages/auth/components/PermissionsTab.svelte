<script lang="ts">
  import { Search, Loader2 } from '@lucide/svelte';
  import M3Select from '../../../components/m3/M3Select.svelte';
  import type { PermissionRecord } from '../../../types/auth';

  interface Props {
    permissions: PermissionRecord[];
    isLoading: boolean;
    permSearch: string;
    selectedPermModule: string;
    onSearchChange: (search: string) => void;
    onModuleChange: (mod: string) => void;
  }

  let {
    permissions = [],
    isLoading = false,
    permSearch = '',
    selectedPermModule = '',
    onSearchChange = () => {},
    onModuleChange = () => {}
  }: Props = $props();

  let permissionModules = $derived(
    Array.from(new Set(permissions.map((p) => p.module || 'Umum'))).sort()
  );

  let filteredPermissions = $derived(
    permissions.filter((p) => {
      const q = permSearch.toLowerCase().trim();
      const matchSearch =
        !q ||
        p.name.toLowerCase().includes(q) ||
        p.code.toLowerCase().includes(q) ||
        (p.description || '').toLowerCase().includes(q);
      const matchModule = !selectedPermModule || p.module === selectedPermModule;
      return matchSearch && matchModule;
    })
  );
</script>

<div class="space-y-4">
  <!-- Search & Filter Bar -->
  <div class="flex flex-wrap items-center gap-3">
    <div class="relative flex-1 min-w-[240px]">
      <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f]" />
      <input
        type="text"
        value={permSearch}
        oninput={(e) => onSearchChange(e.currentTarget.value)}
        placeholder="Cari kode izin atau nama wewenang..."
        class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] focus:bg-white dark:focus:bg-[#1e1f20] focus:ring-2 focus:ring-[#0b57d0] dark:focus:ring-[#a8c7fa] text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="w-48">
      <M3Select
        value={selectedPermModule}
        onchange={onModuleChange}
        options={[
          { value: '', label: 'Semua Modul' },
          ...permissionModules.map((m) => ({ value: m, label: m }))
        ]}
      />
    </div>
  </div>

  <!-- Permissions Table -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] overflow-hidden shadow-xs">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="bg-[#f8fafd] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5] font-semibold">
            <th class="py-3 px-4">Kode Izin (System Code)</th>
            <th class="py-3 px-4">Nama Izin / Wewenang</th>
            <th class="py-3 px-4">Modul Sistem</th>
            <th class="py-3 px-4">Deskripsi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
          {#if isLoading}
            <tr>
              <td colspan="4" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                <div class="flex items-center justify-center gap-2">
                  <Loader2 class="w-5 h-5 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
                  <span>Memuat katalog izin...</span>
                </div>
              </td>
            </tr>
          {:else if filteredPermissions.length === 0}
            <tr>
              <td colspan="4" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                Tidak ada izin yang sesuai dengan pencarian atau modul yang dipilih.
              </td>
            </tr>
          {:else}
            {#each filteredPermissions as p (p.id)}
              <tr class="hover:bg-[#f8fafd]/80 dark:hover:bg-[#28292a]/50 transition-colors">
                <td class="py-3 px-4 font-mono font-bold text-[#0b57d0] dark:text-[#a8c7fa]">
                  {p.code}
                </td>
                <td class="py-3 px-4 font-medium text-[#1f1f1f] dark:text-[#e3e3e3]">
                  {p.name}
                </td>
                <td class="py-3 px-4">
                  <span class="px-2.5 py-0.5 rounded-md bg-[#f0f4f9] dark:bg-[#28292a] text-[#444746] dark:text-[#c4c7c5] text-[11px] font-semibold border border-[#e1e5ea] dark:border-[#353738]">
                    {p.module}
                  </span>
                </td>
                <td class="py-3 px-4 text-[#747775] dark:text-[#8e918f]">
                  {p.description || '-'}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
