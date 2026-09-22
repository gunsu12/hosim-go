<script lang="ts">
  import {
    Stethoscope,
    FlaskConical,
    Database,
    History,
    Calendar,
    ChevronDown,
    ChevronRight,
    Plus
  } from '@lucide/svelte';

  export interface SubMenuItem {
    id: string;
    label: string;
    badge?: string | null;
  }

  export interface MenuItem {
    id: string;
    label: string;
    shortLabel: string;
    icon: any;
    children: SubMenuItem[];
  }

  interface Props {
    activeNav?: string;
    isExpanded?: boolean;
    onNewEncounter?: () => void;
  }

  let {
    activeNav = $bindable('physical'),
    isExpanded = $bindable(false),
    onNewEncounter = () => {}
  }: Props = $props();

  let openParentIds = $state<string[]>(['clinical']);
  let flyoutMenuId = $state<string | null>(null);

  const navigationStructure: MenuItem[] = [
    {
      id: 'clinical',
      label: 'Pelayanan Medis',
      shortLabel: 'Klinis',
      icon: Stethoscope,
      children: [
        { id: 'physical', label: 'Status Lokalis & Anatomi', badge: 'Body' },
        { id: 'anamnesis', label: 'Anamnesis & SOAP', badge: null },
        { id: 'odontogram', label: 'Odontogram Gigi', badge: 'Gigi' }
      ]
    },
    {
      id: 'diagnostics',
      label: 'Penunjang Medis',
      shortLabel: 'Penunjang',
      icon: FlaskConical,
      children: [
        { id: 'lab', label: 'Laboratorium Cito', badge: 'Cito' },
        { id: 'radiology', label: 'Radiologi & Imaging', badge: null },
        { id: 'pharmacy', label: 'Farmasi & E-Resep', badge: null }
      ]
    },
    {
      id: 'master',
      label: 'Master Data RS',
      shortLabel: 'Master',
      icon: Database,
      children: [
        { id: 'master-patient', label: 'Pasien (Patient)', badge: null },
        { id: 'master-practitioner', label: 'Tenaga Medis (Practitioner)', badge: null },
        { id: 'master-departement', label: 'Departemen & Instalasi', badge: null },
        { id: 'master-service-unit', label: 'Unit Layanan & Poli', badge: null },
        { id: 'master-room', label: 'Ruangan & Bed (Room)', badge: null },
        { id: 'master-payer', label: 'Penjamin & Asuransi (Payer)', badge: null },
        { id: 'master-referal', label: 'Faskes Rujukan (Referal)', badge: null },
        { id: 'master-tariff-class', label: 'Kelas Tarif (Tariff Class)', badge: null }
      ]
    },
    {
      id: 'activity',
      label: 'Aktivitas Pasien',
      shortLabel: 'Aktivitas',
      icon: History,
      children: [
        { id: 'history', label: 'Riwayat Kunjungan RME', badge: null },
        { id: 'schedule', label: 'Jadwal Kontrol Poliklinik', badge: null }
      ]
    }
  ];

  function toggleAccordion(parentId: string): void {
    if (openParentIds.includes(parentId)) {
      openParentIds = openParentIds.filter(id => id !== parentId);
    } else {
      openParentIds = [...openParentIds, parentId];
    }
  }

  function handleSubmenuClick(subId: string): void {
    activeNav = subId;
    flyoutMenuId = null;
  }

  function isParentActive(parent: MenuItem): boolean {
    return parent.children?.some(c => c.id === activeNav);
  }

  // Otomatis buka accordion parent jika activeNav berada di dalamnya (misal saat refresh/direct link URL)
  $effect(() => {
    const parent = navigationStructure.find(p => p.children?.some(c => c.id === activeNav));
    if (parent && !openParentIds.includes(parent.id)) {
      openParentIds = [...openParentIds, parent.id];
    }
  });
</script>

<aside
  class="py-3 flex flex-col gap-3 select-none shrink-0 transition-all duration-200 relative {isExpanded ? 'w-64 px-3' : 'w-[72px] px-2 items-center'}"
>
  <!-- Tombol "+ Simpan / New" ala Google Drive -->
  {#if isExpanded}
    <button
      type="button"
      onclick={onNewEncounter}
      class="flex items-center gap-3 h-12 px-4 rounded-2xl bg-white text-[#1f1f1f] shadow-md hover:shadow-lg hover:bg-[#f8fafd] active:bg-[#e9eef6] border border-[#e1e5ea] transition-all duration-150 cursor-pointer w-full font-medium text-xs tracking-wide"
    >
      <Plus class="w-5 h-5 text-[#0b57d0]" />
      <span>Simpan RME</span>
    </button>
  {:else}
    <button
      type="button"
      onclick={onNewEncounter}
      class="flex items-center justify-center w-12 h-12 rounded-2xl bg-white text-[#0b57d0] shadow-md hover:shadow-lg hover:bg-[#f8fafd] active:bg-[#e9eef6] border border-[#e1e5ea] transition-all duration-150 cursor-pointer group"
      title="Simpan RME Pasien"
    >
      <Plus class="w-6 h-6" />
    </button>
  {/if}

  <!-- List Menu Item dengan Sub-Menu -->
  <nav class="flex flex-col gap-1.5 w-full">
    {#each navigationStructure as item}
      {@const parentActive = isParentActive(item)}
      {@const isOpen = openParentIds.includes(item.id)}
      {@const isFlyoutOpen = flyoutMenuId === item.id}

      {#if isExpanded}
        <!-- MODE DIPERLUAS: ACCORDION TREE -->
        <div class="flex flex-col">
          <button
            type="button"
            onclick={() => toggleAccordion(item.id)}
            class="flex items-center justify-between px-3 py-2.5 rounded-full transition-colors cursor-pointer {parentActive ? 'bg-[#c2e7ff]/60 text-[#001d35] font-semibold' : 'text-[#444746] hover:bg-[#e9eef6] hover:text-[#1f1f1f]'}"
          >
            <div class="flex items-center gap-3">
              <item.icon class="w-5 h-5 {parentActive ? 'text-[#0b57d0]' : 'text-[#444746]'}" />
              <span class="text-xs">{item.label}</span>
            </div>
            {#if isOpen}
              <ChevronDown class="w-4 h-4 text-[#747775]" />
            {:else}
              <ChevronRight class="w-4 h-4 text-[#747775]" />
            {/if}
          </button>

          {#if isOpen}
            <div class="flex flex-col gap-1 pl-7 pr-1 mt-1 border-l-2 border-[#e1e5ea] ml-5 py-1 animate-in fade-in duration-150">
              {#each item.children as sub}
                {@const isSubActive = activeNav === sub.id}
                <button
                  type="button"
                  onclick={() => handleSubmenuClick(sub.id)}
                  class="flex items-center justify-between px-3 py-2 text-xs rounded-full transition-colors cursor-pointer text-left {isSubActive ? 'bg-[#0b57d0] text-white font-semibold shadow-xs' : 'text-[#444746] hover:bg-[#e9eef6] hover:text-[#1f1f1f]'}"
                >
                  <span class="truncate">{sub.label}</span>
                  {#if sub.badge}
                    <span class="text-[10px] px-1.5 py-0.2 rounded-md font-semibold {isSubActive ? 'bg-white/20 text-white' : 'bg-[#e9eef6] text-[#444746]'}">
                      {sub.badge}
                    </span>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>

      {:else}
        <!-- MODE RINGKAS (72px): FLYOUT POPOVER -->
        <div class="relative flex flex-col items-center w-full">
          <button
            type="button"
            onclick={() => {
              flyoutMenuId = isFlyoutOpen ? null : item.id;
            }}
            class="flex flex-col items-center gap-1 w-full group cursor-pointer"
            title="{item.label} (Klik untuk sub-menu)"
          >
            <div
              class="w-12 h-8 rounded-full flex items-center justify-center transition-all duration-150 {parentActive ? 'bg-[#c2e7ff] text-[#001d35]' : 'text-[#444746] group-hover:bg-[#e9eef6]'}"
            >
              <item.icon class="w-5 h-5 {parentActive ? 'text-[#001d35]' : 'text-[#444746]'}" />
            </div>
            <span class="text-[10px] text-center leading-tight line-clamp-1 {parentActive ? 'text-[#001d35] font-semibold' : 'text-[#444746]'}">
              {item.shortLabel}
            </span>
          </button>

          {#if isFlyoutOpen}
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
              class="fixed inset-0 z-40"
              onclick={() => flyoutMenuId = null}
            ></div>

            <div
              class="absolute left-full top-0 ml-3 w-60 bg-white rounded-2xl border border-[#e1e5ea] shadow-xl p-2.5 z-50 animate-in fade-in zoom-in-95 duration-150"
            >
              <div class="px-3 py-1.5 border-b border-[#e1e5ea] mb-1.5">
                <div class="text-xs font-bold text-[#1f1f1f] flex items-center gap-2">
                  <item.icon class="w-4 h-4 text-[#0b57d0]" />
                  <span>{item.label}</span>
                </div>
                <div class="text-[10px] text-[#747775]">Pilih modul pemeriksaan</div>
              </div>

              <div class="flex flex-col gap-1">
                {#each item.children as sub}
                  {@const isSubActive = activeNav === sub.id}
                  <button
                    type="button"
                    onclick={() => handleSubmenuClick(sub.id)}
                    class="flex items-center justify-between px-3 py-2 rounded-xl text-xs transition-colors cursor-pointer text-left {isSubActive ? 'bg-[#c2e7ff] text-[#001d35] font-semibold' : 'text-[#444746] hover:bg-[#f0f4f9] hover:text-[#1f1f1f]'}"
                  >
                    <span>{sub.label}</span>
                    {#if sub.badge}
                      <span class="text-[10px] px-1.5 py-0.5 rounded-md font-semibold {isSubActive ? 'bg-white text-[#001d35]' : 'bg-[#e9eef6] text-[#444746]'}">
                        {sub.badge}
                      </span>
                    {/if}
                  </button>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      {/if}
    {/each}
  </nav>
</aside>
