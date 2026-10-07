<script lang="ts">
  import {
    Stethoscope,
    Database,
    History,
    Calendar,
    ChevronDown,
    ChevronRight,
    Plus,
    ArrowLeft,
    LayoutGrid,
    ShieldCheck,
    BookOpen,
  } from "@lucide/svelte";
  import { auth } from "../stores/auth.svelte";

  export interface SubMenuItem {
    id: string;
    label: string;
    badge?: string | null;
    permission?: string;
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
    activeModule?: string | null;
    isExpanded?: boolean;
    onNewEncounter?: () => void;
    onNavigateHome?: () => void;
  }

  let {
    activeNav = $bindable("physical"),
    activeModule = null,
    isExpanded = $bindable(false),
    onNewEncounter = () => {},
    onNavigateHome = () => {},
  }: Props = $props();

  let openParentIds = $state<string[]>(["clinical", "master", "accounting"]);
  let flyoutMenuId = $state<string | null>(null);

  const navigationStructure: MenuItem[] = [
    {
      id: "clinical",
      label: "Pelayanan Medis",
      shortLabel: "Klinis",
      icon: Stethoscope,
      children: [
        {
          id: "physical",
          label: "Status Lokalis & Anatomi",
          badge: "Body",
          permission: "patient:read",
        },
        {
          id: "anamnesis",
          label: "Anamnesis & SOAP",
          badge: null,
          permission: "patient:read",
        },
      ],
    },
    {
      id: "master",
      label: "Master Data RS",
      shortLabel: "Master",
      icon: Database,
      children: [
        {
          id: "master-patient",
          label: "Pasien (Patient)",
          badge: null,
          permission: "patient:read",
        },
        {
          id: "master-practitioner",
          label: "Tenaga Medis (Practitioner)",
          badge: null,
          permission: "practitioner:read",
        },
        {
          id: "master-departement",
          label: "Departemen & Instalasi",
          badge: null,
          permission: "department:read",
        },
        {
          id: "master-service-unit",
          label: "Unit Layanan & Poli",
          badge: null,
          permission: "service_unit:read",
        },
        {
          id: "master-room",
          label: "Ruangan & Bed (Room)",
          badge: null,
          permission: "room:read",
        },
        {
          id: "master-customer",
          label: "Debitur & Penjamin (Customer)",
          badge: null,
          permission: "customer:read",
        },
        {
          id: "master-referal",
          label: "Faskes Rujukan (Referal)",
          badge: null,
          permission: "referal:read",
        },
        {
          id: "master-tariff-class",
          label: "Kelas Tarif (Tariff Class)",
          badge: null,
          permission: "tariff_class:read",
        },
        {
          id: "master-item",
          label: "Katalog Item & Tarif",
          badge: "Universal",
          permission: "item:read",
        },
      ],
    },
    {
      id: "accounting",
      label: "Chart of Account",
      shortLabel: "Coa",
      icon: BookOpen,
      children: [
        {
          id: "accounting-coa",
          label: "Bagan Akun (COA)",
          badge: "Aktif",
          permission: "accounting:read",
        },
        {
          id: "accounting-journals",
          label: "Jurnal Umum Transaksi",
          badge: "Tahap 2",
          permission: "accounting:read",
        },
        {
          id: "accounting-ledger",
          label: "Buku Besar & Kas/Bank",
          badge: "Tahap 3",
          permission: "accounting:read",
        },
        {
          id: "accounting-reports",
          label: "Laporan Keuangan & Neraca",
          badge: "Tahap 4",
          permission: "accounting:read",
        },
      ],
    },
    {
      id: "auth",
      label: "Pengguna & Akses",
      shortLabel: "Akses",
      icon: ShieldCheck,
      children: [
        {
          id: "auth-users",
          label: "Kelola Pengguna",
          badge: null,
          permission: "user:read",
        },
        {
          id: "auth-roles",
          label: "Peran & Hak Akses",
          badge: null,
          permission: "role:read",
        },
        {
          id: "auth-permissions",
          label: "Katalog Izin",
          badge: null,
          permission: "role:read",
        },
      ],
    },
    {
      id: "activity",
      label: "Aktivitas Pasien",
      shortLabel: "Aktivitas",
      icon: History,
      children: [
        {
          id: "history",
          label: "Riwayat Kunjungan RME",
          badge: null,
          permission: "patient:read",
        },
        {
          id: "schedule",
          label: "Jadwal Kontrol Poliklinik",
          badge: null,
          permission: "appointment:view",
        },
      ],
    },
  ];

  // Filter navigasi berdasarkan konteks modul aktif (Focused Contextual Sidebar)
  const visibleNavigation = $derived(
    navigationStructure
      .filter((item) => {
        if (!activeModule) return true;
        if (activeModule === "master") return item.id === "master";
        if (activeModule === "clinical")
          return item.id === "clinical" || item.id === "activity";
        if (activeModule === "auth") return item.id === "auth";
        if (activeModule === "accounting") return item.id === "accounting";
        return item.id === activeModule;
      })
      .map((item) => ({
        ...item,
        children: item.children.filter(
          (sub) => !sub.permission || auth.hasPermission(sub.permission),
        ),
      }))
      .filter((item) => item.children.length > 0),
  );

  function toggleAccordion(parentId: string): void {
    if (openParentIds.includes(parentId)) {
      openParentIds = openParentIds.filter((id) => id !== parentId);
    } else {
      openParentIds = [...openParentIds, parentId];
    }
  }

  function handleSubmenuClick(subId: string): void {
    activeNav = subId;
    flyoutMenuId = null;
  }

  function isParentActive(parent: MenuItem): boolean {
    return parent.children?.some((c) => c.id === activeNav);
  }

  // Otomatis buka accordion parent jika activeNav berada di dalamnya
  $effect(() => {
    const parent = visibleNavigation.find((p) =>
      p.children?.some((c) => c.id === activeNav),
    );
    if (parent && !openParentIds.includes(parent.id)) {
      openParentIds = [...openParentIds, parent.id];
    }
  });
</script>

<aside
  class="py-3 flex flex-col gap-2 select-none shrink-0 transition-all duration-200 relative {isExpanded
    ? 'w-64 px-3'
    : 'w-[72px] px-2 items-center'}"
>
  <!-- Tombol Navigasi Kembali ke Desktop Launcher -->
  {#if isExpanded}
    <button
      type="button"
      onclick={onNavigateHome}
      class="flex items-center gap-2.5 px-3.5 py-2.5 rounded-2xl bg-[#e8f0fe] text-[#0b57d0] hover:bg-[#d3e3fd] text-xs font-semibold transition-all cursor-pointer w-full border border-[#d3e3fd]/60 mb-1"
      title="Kembali ke Beranda Modul (Desktop Launcher)"
    >
      <ArrowLeft class="w-4 h-4" />
      <span>Beranda Modul</span>
    </button>
  {:else}
    <button
      type="button"
      onclick={onNavigateHome}
      class="flex items-center justify-center w-12 h-10 rounded-2xl bg-[#e8f0fe] text-[#0b57d0] hover:bg-[#d3e3fd] transition-all cursor-pointer mb-1"
      title="Kembali ke Beranda Modul (Desktop Launcher)"
    >
      <ArrowLeft class="w-5 h-5" />
    </button>
  {/if}

  <!-- Tombol Aksi Simpan / Tambah Cepat jika berada di modul Klinis -->
  {#if activeModule === "clinical"}
    {#if isExpanded}
      <button
        type="button"
        onclick={onNewEncounter}
        class="flex items-center gap-3 h-11 px-4 rounded-2xl bg-white text-[#1f1f1f] shadow-xs hover:shadow-md hover:bg-[#f8fafd] active:bg-[#e9eef6] border border-[#e1e5ea] transition-all duration-150 cursor-pointer w-full font-medium text-xs tracking-wide mb-2"
      >
        <Plus class="w-5 h-5 text-[#0b57d0]" />
        <span>Simpan RME</span>
      </button>
    {:else}
      <button
        type="button"
        onclick={onNewEncounter}
        class="flex items-center justify-center w-12 h-11 rounded-2xl bg-white text-[#0b57d0] shadow-xs hover:shadow-md hover:bg-[#f8fafd] active:bg-[#e9eef6] border border-[#e1e5ea] transition-all duration-150 cursor-pointer group mb-2"
        title="Simpan RME Pasien"
      >
        <Plus class="w-5 h-5" />
      </button>
    {/if}
  {/if}

  <!-- List Menu Item dengan Sub-Menu Terfokus -->
  <nav class="flex flex-col gap-1 w-full">
    {#each visibleNavigation as item}
      {@const parentActive = isParentActive(item)}
      {@const isOpen = openParentIds.includes(item.id)}
      {@const isFlyoutOpen = flyoutMenuId === item.id}

      {#if isExpanded}
        <!-- MODE DIPERLUAS: ACCORDION TREE -->
        <div class="flex flex-col">
          <button
            type="button"
            onclick={() => toggleAccordion(item.id)}
            class="flex items-center justify-between px-3 py-2 rounded-full transition-colors cursor-pointer {parentActive
              ? 'bg-[#c2e7ff]/60 text-[#001d35] font-semibold'
              : 'text-[#444746] hover:bg-[#e9eef6] hover:text-[#1f1f1f]'}"
          >
            <div class="flex items-center gap-3">
              <item.icon
                class="w-4.5 h-4.5 {parentActive
                  ? 'text-[#0b57d0]'
                  : 'text-[#444746]'}"
              />
              <span class="text-xs">{item.label}</span>
            </div>
            {#if isOpen}
              <ChevronDown class="w-4 h-4 text-[#747775]" />
            {:else}
              <ChevronRight class="w-4 h-4 text-[#747775]" />
            {/if}
          </button>

          {#if isOpen}
            <div
              class="flex flex-col gap-0.5 pl-6 pr-1 mt-1 border-l-2 border-[#e1e5ea] ml-5 py-1 animate-in fade-in duration-150"
            >
              {#each item.children as sub}
                {@const isSubActive = activeNav === sub.id}
                <button
                  type="button"
                  onclick={() => handleSubmenuClick(sub.id)}
                  class="flex items-center justify-between px-3 py-2 text-xs rounded-full transition-colors cursor-pointer text-left {isSubActive
                    ? 'bg-[#0b57d0] text-white font-semibold shadow-xs'
                    : 'text-[#444746] hover:bg-[#e9eef6] hover:text-[#1f1f1f]'}"
                >
                  <span class="truncate">{sub.label}</span>
                  {#if sub.badge}
                    <span
                      class="text-[10px] px-1.5 py-0.2 rounded-md font-semibold {isSubActive
                        ? 'bg-white/20 text-white'
                        : 'bg-[#e9eef6] text-[#444746]'}"
                    >
                      {sub.badge}
                    </span>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {:else}
        <!-- MODE MINIMAL: COMPACT ICON RAIL DENGAN FLYOUT POPOVER -->
        <div class="relative w-full flex justify-center">
          <button
            type="button"
            onclick={() => (flyoutMenuId = isFlyoutOpen ? null : item.id)}
            class="flex flex-col items-center justify-center w-12 h-12 rounded-2xl transition-all cursor-pointer group {parentActive
              ? 'bg-[#c2e7ff] text-[#001d35]'
              : 'text-[#444746] hover:bg-[#e9eef6] hover:text-[#1f1f1f]'}"
            title={item.label}
          >
            <item.icon
              class="w-5 h-5 group-hover:scale-110 transition-transform {parentActive
                ? 'text-[#0b57d0]'
                : 'text-[#444746]'}"
            />
            <span
              class="text-[10px] font-medium tracking-tighter mt-1 truncate max-w-[56px] text-center"
            >
              {item.shortLabel}
            </span>
          </button>

          <!-- Flyout Menu Popover saat hover/klik icon rail -->
          {#if isFlyoutOpen}
            <div
              class="absolute left-14 top-0 w-56 bg-white rounded-2xl border border-[#e1e5ea] shadow-lg p-2 z-50 flex flex-col gap-1 animate-in fade-in duration-150"
            >
              <div
                class="px-3 py-1.5 text-xs font-bold text-[#1f1f1f] border-b border-[#e1e5ea] mb-1 flex items-center justify-between"
              >
                <span>{item.label}</span>
                <span class="text-[10px] text-[#747775] font-normal"
                  >{item.children.length} Menu</span
                >
              </div>
              {#each item.children as sub}
                {@const isSubActive = activeNav === sub.id}
                <button
                  type="button"
                  onclick={() => handleSubmenuClick(sub.id)}
                  class="flex items-center justify-between px-3 py-2 text-xs rounded-xl transition-colors cursor-pointer text-left {isSubActive
                    ? 'bg-[#0b57d0] text-white font-medium shadow-xs'
                    : 'text-[#444746] hover:bg-[#e9eef6] hover:text-[#1f1f1f]'}"
                >
                  <span class="truncate">{sub.label}</span>
                  {#if sub.badge}
                    <span
                      class="text-[10px] px-1.5 py-0.2 rounded-md font-semibold {isSubActive
                        ? 'bg-white/20 text-white'
                        : 'bg-[#e9eef6] text-[#444746]'}"
                    >
                      {sub.badge}
                    </span>
                  {/if}
                </button>
              {/each}
            </div>
          {/if}
        </div>
      {/if}
    {/each}
  </nav>
</aside>
