<script lang="ts">
  import {
    Search,
    Stethoscope,
    Plus,
    Edit2,
    Trash2,
    Loader2
  } from '@lucide/svelte';
  import M3Select from '../../../components/m3/M3Select.svelte';
  import type { UserRecord, RoleRecord } from '../../../types/auth';
  import type { PractitionerRecord } from '../../../types/master/practitioner';

  interface Props {
    users: UserRecord[];
    roles: RoleRecord[];
    practitioners: PractitionerRecord[];
    isLoading: boolean;
    userSearch: string;
    filterRole: string;
    filterStatus: string;
    onSearchChange: (search: string) => void;
    onFilterRoleChange: (roleId: string) => void;
    onFilterStatusChange: (status: string) => void;
    onOpenEditUser: (user: UserRecord) => void;
    onOpenDeleteUser: (user: UserRecord) => void;
    onOpenLinkPractitioner: (user: UserRecord) => void;
  }

  let {
    users = [],
    roles = [],
    practitioners = [],
    isLoading = false,
    userSearch = '',
    filterRole = '',
    filterStatus = '',
    onSearchChange = () => {},
    onFilterRoleChange = () => {},
    onFilterStatusChange = () => {},
    onOpenEditUser = () => {},
    onOpenDeleteUser = () => {},
    onOpenLinkPractitioner = () => {}
  }: Props = $props();

  let practitionerMap = $derived(new Map(practitioners.map((p) => [p.id, p])));
</script>

<div class="space-y-4">
  <!-- Search & Filter Bar -->
  <div class="flex flex-wrap items-center gap-3">
    <div class="relative flex-1 min-w-[240px]">
      <Search class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#747775] dark:text-[#8e918f]" />
      <input
        type="text"
        value={userSearch}
        oninput={(e) => onSearchChange(e.currentTarget.value)}
        placeholder="Cari username, nama, atau email..."
        class="w-full h-10 pl-9 pr-4 text-xs rounded-xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] focus:bg-white dark:focus:bg-[#1e1f20] focus:ring-2 focus:ring-[#0b57d0] dark:focus:ring-[#a8c7fa] text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="w-48">
      <M3Select
        value={filterRole}
        onchange={onFilterRoleChange}
        options={[
          { value: '', label: 'Semua Peran' },
          ...roles.map((r) => ({ value: r.id, label: r.name, badge: r.code }))
        ]}
      />
    </div>

    <div class="w-36">
      <M3Select
        value={filterStatus}
        onchange={onFilterStatusChange}
        options={[
          { value: '', label: 'Semua Status' },
          { value: 'true', label: 'Aktif' },
          { value: 'false', label: 'Nonaktif' }
        ]}
      />
    </div>
  </div>

  <!-- Users Table -->
  <div class="bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] overflow-hidden shadow-xs">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-xs border-collapse">
        <thead>
          <tr class="bg-[#f8fafd] dark:bg-[#28292a] border-b border-[#e1e5ea] dark:border-[#353738] text-[#444746] dark:text-[#c4c7c5] font-semibold">
            <th class="py-3 px-4">Pengguna</th>
            <th class="py-3 px-4">Email</th>
            <th class="py-3 px-4">Peran (Role)</th>
            <th class="py-3 px-4">Tenaga Medis (Link)</th>
            <th class="py-3 px-4">Status</th>
            <th class="py-3 px-4">Terdaftar</th>
            <th class="py-3 px-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#e1e5ea] dark:divide-[#353738]">
          {#if isLoading}
            <tr>
              <td colspan="7" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                <div class="flex items-center justify-center gap-2">
                  <Loader2 class="w-5 h-5 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
                  <span>Memuat data pengguna...</span>
                </div>
              </td>
            </tr>
          {:else if users.length === 0}
            <tr>
              <td colspan="7" class="py-12 text-center text-[#747775] dark:text-[#8e918f]">
                Tidak ada data pengguna yang sesuai dengan filter.
              </td>
            </tr>
          {:else}
            {#each users as u (u.id)}
              <tr class="hover:bg-[#f8fafd]/80 dark:hover:bg-[#28292a]/50 transition-colors">
                <td class="py-3 px-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa] font-bold flex items-center justify-center text-xs shrink-0">
                      {u.name ? u.name.charAt(0).toUpperCase() : u.username.charAt(0).toUpperCase()}
                    </div>
                    <div class="min-w-0">
                      <div class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">{u.name}</div>
                      <div class="text-[11px] text-[#747775] dark:text-[#8e918f] font-mono truncate">@{u.username}</div>
                    </div>
                  </div>
                </td>
                <td class="py-3 px-4 text-[#444746] dark:text-[#c4c7c5] font-mono text-[11px]">
                  {u.email}
                </td>
                <td class="py-3 px-4">
                  <span class="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-semibold {u.role === 'ADMIN' ? 'bg-purple-100 dark:bg-purple-950/60 text-purple-800 dark:text-purple-300' : 'bg-blue-100 dark:bg-blue-950/60 text-blue-800 dark:text-blue-300'}">
                    {u.role_name || u.role || 'Staf'}
                  </span>
                </td>
                <td class="py-3 px-4">
                  {#if u.practitioner_id}
                    {@const prac = practitionerMap.get(u.practitioner_id)}
                    <button
                      type="button"
                      onclick={() => onOpenLinkPractitioner(u)}
                      class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-xl bg-teal-50 dark:bg-teal-950/60 hover:bg-teal-100 dark:hover:bg-teal-900/60 border border-teal-200 dark:border-teal-800 text-teal-800 dark:text-teal-200 text-[11px] font-semibold transition-colors cursor-pointer text-left"
                      title="Klik untuk mengubah atau memutuskan tautan nakes"
                    >
                      <Stethoscope class="w-3.5 h-3.5 text-teal-600 dark:text-teal-400 shrink-0" />
                      <span class="truncate max-w-[130px] font-bold">{prac?.name || 'Tertaut'}</span>
                    </button>
                  {:else}
                    <button
                      type="button"
                      onclick={() => onOpenLinkPractitioner(u)}
                      class="inline-flex items-center gap-1 px-2 py-0.5 rounded-lg border border-dashed border-[#c4c7c5] dark:border-[#5e5e5e] hover:border-teal-500 hover:text-teal-700 dark:hover:text-teal-300 text-[#747775] dark:text-[#8e918f] text-[10px] font-medium transition-colors cursor-pointer"
                      title="Tautkan akun ini ke data tenaga medis"
                    >
                      <Plus class="w-3 h-3" />
                      <span>Tautkan Nakes</span>
                    </button>
                  {/if}
                </td>
                <td class="py-3 px-4">
                  {#if u.is_active}
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#e6f4ea] dark:bg-emerald-950/60 text-[#137333] dark:text-emerald-300">
                      <span class="w-1.5 h-1.5 rounded-full bg-[#137333] dark:bg-emerald-400"></span>
                      Aktif
                    </span>
                  {:else}
                    <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#fce8e6] dark:bg-rose-950/60 text-[#c5221f] dark:text-rose-300">
                      <span class="w-1.5 h-1.5 rounded-full bg-[#c5221f] dark:bg-rose-400"></span>
                      Nonaktif
                    </span>
                  {/if}
                </td>
                <td class="py-3 px-4 text-[#747775] dark:text-[#8e918f] text-[11px]">
                  {new Date(u.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}
                </td>
                <td class="py-3 px-4 text-right">
                  <div class="flex items-center justify-end gap-1">
                    <button
                      type="button"
                      onclick={() => onOpenLinkPractitioner(u)}
                      class="p-1.5 text-[#444746] dark:text-[#c4c7c5] hover:text-teal-600 dark:hover:text-teal-400 hover:bg-teal-50 dark:hover:bg-teal-950/50 rounded-lg transition-colors cursor-pointer"
                      title={u.practitioner_id ? 'Kelola Tautan Tenaga Medis' : 'Tautkan ke Tenaga Medis'}
                    >
                      <Stethoscope class="w-4 h-4" />
                    </button>
                    <button
                      type="button"
                      onclick={() => onOpenEditUser(u)}
                      class="p-1.5 text-[#444746] dark:text-[#c4c7c5] hover:text-[#0b57d0] dark:hover:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] rounded-lg transition-colors cursor-pointer"
                      title="Ubah Data Pengguna"
                    >
                      <Edit2 class="w-4 h-4" />
                    </button>
                    {#if u.username !== 'admin'}
                      <button
                        type="button"
                        onclick={() => onOpenDeleteUser(u)}
                        class="p-1.5 text-[#444746] dark:text-[#c4c7c5] hover:text-[#c5221f] dark:hover:text-rose-400 hover:bg-[#fce8e6] dark:hover:bg-rose-950/50 rounded-lg transition-colors cursor-pointer"
                        title="Hapus / Nonaktifkan Pengguna"
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
