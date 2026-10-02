<script lang="ts">
  import {
    Stethoscope,
    Ambulance,
    Bed,
    Pill,
    Boxes,
    CreditCard,
    Database,
    Shield,
    ShieldCheck,
    Search,
    ArrowLeft,
    ChevronRight,
    Sparkles,
    Calendar,
    Clock,
    Activity,
    Layers,
    FileText,
    Users
  } from '@lucide/svelte';
  import { auth } from '../../stores/auth.svelte';

  interface Props {
    onSelectModule: (moduleId: string, defaultNavId: string) => void;
  }

  let { onSelectModule }: Props = $props();

  let searchQuery = $state('');

  // Tanggal dan waktu saat ini
  const currentDateFormatted = new Intl.DateTimeFormat('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  }).format(new Date());

  export interface ModuleDefinition {
    id: string;
    title: string;
    shortTitle: string;
    category: 'clinical' | 'pharmacy_logistics' | 'admin_finance';
    categoryLabel: string;
    description: string;
    icon: any;
    color: {
      bg: string;
      text: string;
      border: string;
      hoverBg: string;
      badge: string;
    };
    defaultNav: string;
    status: 'ready' | 'in_progress';
    statusText: string;
    permission?: string;
    features: string[];
    quickStats?: string;
  }

  const modules: ModuleDefinition[] = [
    // 1. Pelayanan Klinis & Asuhan Pasien
    {
      id: 'clinical',
      title: 'Rekam Medis & Status Lokalis',
      shortTitle: 'Klinis EHR',
      category: 'clinical',
      categoryLabel: 'Pelayanan Pasien & Asuhan Medis',
      description: 'Pemeriksaan status lokalis anatomi tubuh interaktif, catatan anamnesis, dan form SOAP pasien.',
      icon: Stethoscope,
      color: {
        bg: 'bg-blue-50',
        text: 'text-blue-700',
        border: 'border-blue-200',
        hoverBg: 'hover:border-blue-500 hover:shadow-blue-100',
        badge: 'bg-blue-100 text-blue-800'
      },
      defaultNav: 'physical',
      status: 'ready',
      statusText: 'Aktif & Terintegrasi',
      permission: 'patient:read',
      features: ['Body Diagram Interaktif', 'Resume SOAP', 'Tanda Vital'],
      quickStats: 'Sesuai SNOMED / ICD-10'
    },
    {
      id: 'outpatient',
      title: 'Pelayanan Rawat Jalan',
      shortTitle: 'Rawat Jalan',
      category: 'clinical',
      categoryLabel: 'Pelayanan Pasien & Asuhan Medis',
      description: 'Manajemen antrean poliklinik, registrasi kunjungan harian, dan konsultasi dokter spesialis.',
      icon: Users,
      color: {
        bg: 'bg-teal-50',
        text: 'text-teal-700',
        border: 'border-teal-200',
        hoverBg: 'hover:border-teal-500 hover:shadow-teal-100',
        badge: 'bg-teal-100 text-teal-800'
      },
      defaultNav: 'outpatient-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'outpatient:register',
      features: ['Antrean Poli', 'Pendaftaran Kunjungan', 'Konsultasi Nakes']
    },
    {
      id: 'emergency',
      title: 'Instalasi Gawat Darurat (IGD)',
      shortTitle: 'IGD Cito',
      category: 'clinical',
      categoryLabel: 'Pelayanan Pasien & Asuhan Medis',
      description: 'Alur triase kedaruratan ESI, resusitasi, tindakan akut, dan observasi penanganan cito 24 jam.',
      icon: Ambulance,
      color: {
        bg: 'bg-rose-50',
        text: 'text-rose-700',
        border: 'border-rose-200',
        hoverBg: 'hover:border-rose-500 hover:shadow-rose-100',
        badge: 'bg-rose-100 text-rose-800'
      },
      defaultNav: 'emergency-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'patient:read',
      features: ['Triase ESI 1-5', 'Observasi Kritis', 'Disposisi IGD Cito']
    },
    {
      id: 'inpatient',
      title: 'Pelayanan Rawat Inap',
      shortTitle: 'Rawat Inap',
      category: 'clinical',
      categoryLabel: 'Pelayanan Pasien & Asuhan Medis',
      description: 'Admisi rawat inap, pemetaan ketersediaan bed bangsal, transfer pasien, dan rencana pemulangan.',
      icon: Bed,
      color: {
        bg: 'bg-indigo-50',
        text: 'text-indigo-700',
        border: 'border-indigo-200',
        hoverBg: 'hover:border-indigo-500 hover:shadow-indigo-100',
        badge: 'bg-indigo-100 text-indigo-800'
      },
      defaultNav: 'inpatient-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'patient:read',
      features: ['Admisi Bangsal', 'Alokasi Bed Realtime', 'Resume Medis Pulang']
    },

    // 2. Farmasi & Logistik
    {
      id: 'pharmacy',
      title: 'Farmasi & E-Resep',
      shortTitle: 'Farmasi',
      category: 'pharmacy_logistics',
      categoryLabel: 'Farmasi & Rantai Pasok RS',
      description: 'E-Prescribing terintegrasi, telaah resep dokter, dispensing obat farmasi, dan penyerahan etiket.',
      icon: Pill,
      color: {
        bg: 'bg-emerald-50',
        text: 'text-emerald-700',
        border: 'border-emerald-200',
        hoverBg: 'hover:border-emerald-500 hover:shadow-emerald-100',
        badge: 'bg-emerald-100 text-emerald-800'
      },
      defaultNav: 'pharmacy-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'patient:read',
      features: ['Telaah 7 Benar Resep', 'Dispensing Otomatis', 'Etiket Obat']
    },
    {
      id: 'inventory',
      title: 'Inventori & Logistik RS',
      shortTitle: 'Inventori',
      category: 'pharmacy_logistics',
      categoryLabel: 'Farmasi & Rantai Pasok RS',
      description: 'Katalog master barang/obat, depo ruangan & gudang farmasi, mutasi stok, serta penyesuaian opname.',
      icon: Boxes,
      color: {
        bg: 'bg-amber-50',
        text: 'text-amber-700',
        border: 'border-amber-200',
        hoverBg: 'hover:border-amber-500 hover:shadow-amber-100',
        badge: 'bg-amber-100 text-amber-800'
      },
      defaultNav: 'inventory-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'patient:read',
      features: ['Depo & Gudang RS', 'Stok BHP & Obat', 'Mutasi Antar-Unit']
    },

    // 3. Administrasi & Keuangan RS
    {
      id: 'master',
      title: 'Master Data Rumah Sakit',
      shortTitle: 'Master Data',
      category: 'admin_finance',
      categoryLabel: 'Administrasi, Keuangan & Tata Kelola',
      description: 'Konfigurasi master pasien, nakes, departemen, unit layanan poli, ruangan/bed, debitur, dan tarif.',
      icon: Database,
      color: {
        bg: 'bg-purple-50',
        text: 'text-purple-700',
        border: 'border-purple-200',
        hoverBg: 'hover:border-purple-500 hover:shadow-purple-100',
        badge: 'bg-purple-100 text-purple-800'
      },
      defaultNav: 'master-patient',
      status: 'ready',
      statusText: 'Aktif & Terintegrasi',
      permission: 'department:read',
      features: ['Pasien & Nakes', 'Organisasi & Kamar', 'Payer & Tarif RS'],
      quickStats: '9 Sub-Domain Lengkap'
    },
    {
      id: 'billing',
      title: 'Kasir & Billing Pasien',
      shortTitle: 'Billing',
      category: 'admin_finance',
      categoryLabel: 'Administrasi, Keuangan & Tata Kelola',
      description: 'Kalkulasi tarif tindakan, invoice penjaminan BPJS/Asuransi, kasir pembayaran, dan bukti setor.',
      icon: CreditCard,
      color: {
        bg: 'bg-green-50',
        text: 'text-green-700',
        border: 'border-green-200',
        hoverBg: 'hover:border-green-500 hover:shadow-green-100',
        badge: 'bg-green-100 text-green-800'
      },
      defaultNav: 'billing-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'customer:read',
      features: ['Billing Tagihan Pasien', 'Klaim BPJS & Asuransi', 'Kasir Pembayaran']
    },
    {
      id: 'auth',
      title: 'Pengguna & Hak Akses (RBAC)',
      shortTitle: 'Akses & User',
      category: 'admin_finance',
      categoryLabel: 'Administrasi, Keuangan & Tata Kelola',
      description: 'Manajemen pengguna RS, pembuatan akun dokter/nakes, penugasan peran (role), dan konfigurasi matriks perizinan modul.',
      icon: ShieldCheck,
      color: {
        bg: 'bg-violet-50',
        text: 'text-violet-700',
        border: 'border-violet-200',
        hoverBg: 'hover:border-violet-500 hover:shadow-violet-100',
        badge: 'bg-violet-100 text-violet-800'
      },
      defaultNav: 'auth-users',
      status: 'ready',
      statusText: 'Aktif & Terintegrasi',
      permission: 'user:read',
      features: ['Kelola Pengguna', 'Peran & Matriks Izin', 'Reset Password'],
      quickStats: 'Terintegrasi Go Auth'
    },
    {
      id: 'audit',
      title: 'Audit Trail & Keamanan',
      shortTitle: 'Audit Sistem',
      category: 'admin_finance',
      categoryLabel: 'Administrasi, Keuangan & Tata Kelola',
      description: 'Append-only audit log aktivitas pengguna, riwayat mutasi data medis, dan rekam jejak kepatuhan.',
      icon: Shield,
      color: {
        bg: 'bg-slate-50',
        text: 'text-slate-700',
        border: 'border-slate-200',
        hoverBg: 'hover:border-slate-500 hover:shadow-slate-100',
        badge: 'bg-slate-100 text-slate-800'
      },
      defaultNav: 'audit-workspace',
      status: 'in_progress',
      statusText: 'Tahap Integrasi',
      permission: 'patient:read',
      features: ['Immutable Audit Trail', 'Pelacakan Mutasi RME', 'Kepatuhan Regulasi']
    }
  ];

  // Filter modul berdasarkan pencarian dan role/permission
  let filteredModules = $derived(
    modules.filter(m => {
      // Permission check (jika dispesifikasi)
      if (m.permission && !auth.hasPermission(m.permission)) {
        return false;
      }
      if (!searchQuery.trim()) return true;
      const q = searchQuery.toLowerCase();
      return (
        m.title.toLowerCase().includes(q) ||
        m.shortTitle.toLowerCase().includes(q) ||
        m.description.toLowerCase().includes(q) ||
        m.features.some(f => f.toLowerCase().includes(q))
      );
    })
  );

  // Pengelompokan modul per kategori
  let categories = $derived([
    {
      id: 'clinical',
      label: 'Pelayanan Pasien & Asuhan Medis',
      modules: filteredModules.filter(m => m.category === 'clinical')
    },
    {
      id: 'pharmacy_logistics',
      label: 'Farmasi & Rantai Pasok RS',
      modules: filteredModules.filter(m => m.category === 'pharmacy_logistics')
    },
    {
      id: 'admin_finance',
      label: 'Administrasi, Keuangan & Tata Kelola',
      modules: filteredModules.filter(m => m.category === 'admin_finance')
    }
  ].filter(c => c.modules.length > 0));
</script>

<div class="max-w-7xl mx-auto w-full py-4 px-2 sm:px-4 flex flex-col gap-6 animate-in fade-in duration-200">
  <!-- Hero Greeting Banner -->
  <div class="relative overflow-hidden bg-gradient-to-br from-[#0b57d0] to-[#042f77] text-white rounded-3xl p-6 sm:p-8 shadow-sm">
    <!-- Decorative Ambient Glow -->
    <div class="absolute -right-12 -bottom-12 w-64 h-64 bg-white/10 rounded-full blur-2xl pointer-events-none"></div>
    <div class="absolute right-24 top-0 w-48 h-48 bg-blue-400/20 rounded-full blur-xl pointer-events-none"></div>

    <div class="relative z-10 flex flex-col md:flex-row md:items-center justify-between gap-6">
      <div class="flex flex-col gap-2">
        <div class="flex items-center gap-2 text-xs font-medium text-blue-200 tracking-wide uppercase">
          <Calendar class="w-3.5 h-3.5" />
          <span>{currentDateFormatted}</span>
          <span>•</span>
          <span class="inline-flex items-center gap-1">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            HOSIM-GO Modular Monolith Active
          </span>
        </div>

        <h1 class="text-2xl sm:text-3xl font-bold tracking-tight text-white">
          Selamat Datang, {auth.user?.name || 'Dokter Spesialis'}!
        </h1>
        <p class="text-xs sm:text-sm text-blue-100/90 max-w-2xl leading-relaxed">
          Pusat kendali aplikasi rumah sakit berbasis domain. Pilih modul kerja di bawah untuk membuka alur layanan klinis, operasional, atau administrasi terfokus.
        </p>
      </div>

      <!-- Quick Search Bar Modul -->
      <div class="w-full md:w-80 shrink-0">
        <div class="relative flex items-center bg-white/15 backdrop-blur-md rounded-2xl px-3.5 py-2.5 border border-white/20 focus-within:bg-white focus-within:text-[#1f1f1f] text-white transition-all shadow-inner">
          <Search class="w-4 h-4 mr-2.5 text-blue-200 shrink-0 focus-within:text-[#444746]" />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Cari modul atau fitur..."
            class="w-full bg-transparent text-xs placeholder:text-blue-200/80 focus:text-[#1f1f1f] focus:outline-none"
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Sections of Domains -->
  {#each categories as category}
    <div class="flex flex-col gap-3">
      <div class="flex items-center justify-between border-b border-[#e1e5ea] pb-2">
        <div class="flex items-center gap-2">
          <div class="w-2.5 h-2.5 rounded-full bg-[#0b57d0]"></div>
          <h2 class="text-sm font-bold text-[#1f1f1f] tracking-tight">{category.label}</h2>
        </div>
        <span class="text-[11px] font-medium text-[#747775]">{category.modules.length} Modul Terdaftar</span>
      </div>

      <!-- Module Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {#each category.modules as mod}
          <div
            role="button"
            tabindex="0"
            onclick={() => onSelectModule(mod.id, mod.defaultNav)}
            onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') onSelectModule(mod.id, mod.defaultNav); }}
            class="group relative text-left bg-white rounded-2xl p-5 border border-[#e1e5ea] hover:shadow-md transition-all duration-200 cursor-pointer flex flex-col justify-between gap-4 {mod.color.hoverBg}"
          >
            <!-- Header Card (Icon & Badges) -->
            <div class="flex items-start justify-between gap-3">
              <div class="w-12 h-12 rounded-2xl {mod.color.bg} {mod.color.text} flex items-center justify-center shrink-0 shadow-xs group-hover:scale-105 transition-transform duration-200">
                <mod.icon class="w-6 h-6" />
              </div>

              <div class="flex flex-col items-end gap-1.5">
                <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full {mod.status === 'ready' ? 'bg-[#e6f4ea] text-[#137333] border border-[#ceead6]' : 'bg-[#f1f3f4] text-[#5f6368] border border-[#dadce0]'}">
                  {mod.statusText}
                </span>
                {#if mod.quickStats}
                  <span class="text-[9px] font-mono text-[#747775]">{mod.quickStats}</span>
                {/if}
              </div>
            </div>

            <!-- Content Card -->
            <div class="flex flex-col gap-1.5 flex-1">
              <h3 class="text-base font-bold text-[#1f1f1f] group-hover:text-[#0b57d0] transition-colors flex items-center justify-between">
                <span>{mod.title}</span>
                <ChevronRight class="w-4 h-4 text-[#747775] opacity-0 group-hover:opacity-100 group-hover:translate-x-0.5 transition-all" />
              </h3>
              <p class="text-xs text-[#444746] leading-relaxed line-clamp-2">
                {mod.description}
              </p>
            </div>

            <!-- Features Pills & Entry Action -->
            <div class="pt-3 border-t border-[#f0f4f9] flex flex-wrap items-center gap-1.5">
              {#each mod.features as feat}
                <span class="text-[10px] px-2 py-0.5 rounded-md bg-[#f8fafd] border border-[#e1e5ea] text-[#444746]">
                  {feat}
                </span>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/each}

  {#if filteredModules.length === 0}
    <div class="p-12 text-center bg-white rounded-3xl border border-[#e1e5ea] flex flex-col items-center justify-center gap-2 text-[#747775]">
      <Search class="w-8 h-8 text-[#bdc1c6]" />
      <div class="font-semibold text-sm text-[#1f1f1f]">Modul Tidak Ditemukan</div>
      <p class="text-xs max-w-sm">Tidak ada modul yang cocok dengan kata kunci "{searchQuery}". Silakan coba kata kunci lain.</p>
    </div>
  {/if}
</div>
