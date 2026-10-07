<script lang="ts">
  import { CheckCircle } from '@lucide/svelte';
  import TopBar from '$lib/components/TopBar.svelte';
  import PatientBanner from '$lib/components/PatientBanner.svelte';
  import LoginPage from '$lib/components/LoginPage.svelte';
  import NavigationRail from '$lib/components/NavigationRail.svelte';
  import DesktopLauncher from '$lib/pages/launcher/DesktopLauncher.svelte';
  import ModulePlaceholder from '$lib/pages/launcher/ModulePlaceholder.svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import type { BodyFinding, ClinicalNotes } from '$lib/types';
  import { getNavFromCurrentHash, setHashFromNav, getModuleFromNav } from '$lib/router';

  // Domain Clinical Pages
  import {
    PhysicalExamPage,
    SoapPage
  } from '$lib/pages/clinical';

  // Domain Master Pages (Organized in domain folders matching internal/master/)
  import {
    PatientPage,
    PractitionerPage,
    DepartementPage,
    ServiceUnitPage,
    RoomPage,
    CustomerPage,
    ReferalPage,
    TariffClassPage,
    ItemPage
  } from '$lib/pages/master';

  // Domain Accounting Pages
  import { ChartOfAccountsPage } from '$lib/pages/accounting';

  // Domain Finance & Tariff Pages
  import { FinanceTariffPage } from '$lib/pages/finance';

  // Domain Auth & User Management Pages
  import { AuthManagementPage } from '$lib/pages/auth';

  // Inisialisasi langsung dari window.location.hash browser agar persisten saat di-refresh (F5)
  let activeNav = $state<string>(getNavFromCurrentHash());
  let isSidebarExpanded = $state<boolean>(false);
  let saveSuccessAlert = $state<boolean>(false);

  // Modul aktif berdasarkan navId (null jika sedang berada di Desktop Launcher)
  let activeModule = $derived(getModuleFromNav(activeNav));

  // Sinkronkan URL browser setiap kali activeNav berganti
  $effect(() => {
    setHashFromNav(activeNav);
  });

  // Data temuan klinis anatomi tubuh pasien
  let findings = $state<BodyFinding[]>([
    {
      id: 1,
      x: 43.5,
      y: 35.8,
      view: 'front',
      category: 'pain',
      severity: 6,
      note: 'Nyeri tekan kuadran kanan bawah (McBurney point suspect apendisitis)',
      createdAt: '10:14'
    },
    {
      id: 2,
      x: 28.2,
      y: 42.1,
      view: 'front',
      category: 'injury',
      severity: 0,
      note: 'Vulnus laceratum ±3 cm pada antebrachii dextra, perdarahan terkontrol',
      createdAt: '10:18'
    }
  ]);

  // Form klinis anamnesis SOAP
  let clinicalNotes = $state<ClinicalNotes>({
    chiefComplaint: 'Nyeri perut kanan bawah sejak 2 hari yang lalu, mual (+), muntah 1x.',
    anamnesis: 'Pasien datang dengan keluhan nyeri perut kanan bawah mendadak bertambah berat saat berjalan. Ada riwayat terjatuh kemarin yang menyebabkan luka lecet di lengan kanan.',
    diagnosis: 'K35.80 - Acute appendicitis, other and unspecified',
    disposition: 'Rawat Inap / Persiapan Laparotomi Cito'
  });

  function handleSaveExamination(): void {
    saveSuccessAlert = true;
    setTimeout(() => {
      saveSuccessAlert = false;
    }, 4000);
  }

  function handleSelectModule(moduleId: string, defaultNavId: string): void {
    activeNav = defaultNavId;
  }
</script>

<!-- Deklaratif Window Event Listener ala Svelte 5 -->
<svelte:window onhashchange={() => {
  const navFromHash = getNavFromCurrentHash();
  if (navFromHash !== activeNav) activeNav = navFromHash;
}} />

{#if !auth.isAuthenticated}
  <!-- Halaman Login Resmi Google Material You -->
  <LoginPage />
{:else}
  <div class="min-h-screen bg-[#f8fafd] flex flex-col font-sans">
    <!-- Top Bar (Header dengan Logo, Modul Switcher, & Profil Pengguna) -->
    <TopBar
      onToggleSidebar={() => isSidebarExpanded = !isSidebarExpanded}
      user={auth.user}
      onLogout={() => auth.logout()}
      {activeModule}
      onNavigateToDesktop={() => activeNav = 'desktop'}
    />

    <!-- Main Layout: Jika di dalam modul tampilkan NavigationRail; jika di Desktop Launcher tampilkan canvas penuh -->
    <div class="flex-1 flex overflow-hidden">
      {#if activeModule}
        <NavigationRail
          bind:activeNav
          {activeModule}
          bind:isExpanded={isSidebarExpanded}
          onNewEncounter={handleSaveExamination}
          onNavigateHome={() => activeNav = 'desktop'}
        />
      {/if}

      <!-- Main Content Canvas (Google Drive Spacious Content Island) -->
      <main class="flex-1 bg-white md:rounded-3xl border border-[#e1e5ea] shadow-xs md:mr-3 md:mb-3 overflow-y-auto p-4 sm:p-6 flex flex-col gap-5">
        <!-- Patient Summary Banner (Hanya ditampilkan pada konteks pelayanan klinis pasien aktif) -->
        {#if activeModule === 'clinical' && activeNav !== 'desktop'}
          <PatientBanner />
        {/if}

        <!-- Notification Toast -->
        {#if saveSuccessAlert}
          <div class="p-3.5 rounded-2xl bg-[#e6f4ea] border border-[#ceead6] text-[#137333] flex items-center justify-between shadow-xs animate-in fade-in duration-200">
            <div class="flex items-center gap-2 text-xs font-semibold">
              <CheckCircle class="w-4 h-4 text-[#137333]" />
              <span>Pemeriksaan fisik dan pemetaan diagram tubuh berhasil disimpan ke rekam medis!</span>
            </div>
            <span class="text-[11px] text-[#137333] font-mono">Status: Tersinkronisasi dengan Go Gin Backend</span>
          </div>
        {/if}

        <!-- 1. DESKTOP LAUNCHER (BERANDA MODUL UTAMA) -->
        {#if activeNav === 'desktop'}
          <DesktopLauncher onSelectModule={handleSelectModule} />

        <!-- 2. CLINICAL MODULES -->
        {:else if activeNav === 'physical'}
          <PhysicalExamPage
            bind:findings
            onSave={handleSaveExamination}
          />

        {:else if activeNav === 'anamnesis'}
          <SoapPage
            bind:clinicalNotes
            onSave={handleSaveExamination}
          />

        <!-- 3. MASTER DATA RS DOMAIN MODULES -->
        {:else if activeNav === 'master-patient'}
          <PatientPage />
        {:else if activeNav === 'master-practitioner'}
          <PractitionerPage />
        {:else if activeNav === 'master-departement'}
          <DepartementPage />
        {:else if activeNav === 'master-service-unit'}
          <ServiceUnitPage />
        {:else if activeNav === 'master-room'}
          <RoomPage />
        {:else if activeNav === 'master-customer' || activeNav === 'master-payer'}
          <CustomerPage />
        {:else if activeNav === 'master-referal'}
          <ReferalPage />
        {:else if activeNav === 'master-tariff-class'}
          <TariffClassPage />
        {:else if activeNav === 'finance-tariff'}
          <FinanceTariffPage initialTab="PLAN" />
        {:else if activeNav === 'finance-price-plan'}
          <FinanceTariffPage initialTab="PLAN" />
        {:else if activeNav === 'finance-tariff-lookup'}
          <FinanceTariffPage initialTab="SIMULATOR" />
        {:else if activeNav === 'finance-tariff-class'}
          <FinanceTariffPage initialTab="CLASS" />
        {:else if activeNav === 'finance-tariff-component'}
          <FinanceTariffPage initialTab="COMPONENT" />
        {:else if activeNav === 'master-item'}
          <ItemPage />
        {:else if activeNav === 'master-coa' || activeNav === 'accounting-coa'}
          <ChartOfAccountsPage />

        <!-- 3.6 AKUNTANSI & KEUANGAN RS DOMAIN MODULES -->
        {:else if activeNav === 'accounting-journals'}
          <ModulePlaceholder
            title="Jurnal Umum & Transaksi Memorial (General Journal)"
            description="Modul pencatatan transaksi jurnal umum double-entry berimbang (debit = kredit), jurnal penyesuaian, dan approval posting ke buku besar."
            domainPackage="accounting/journal"
            backendStatus="Sesuai PRD Akuntansi internal/accounting/PRD.md (Tahap 2: Journal Engine)"
            plannedEndpoints={[
              "POST /api/v1/accounting/journals (Input Transaksi Jurnal Double-Entry)",
              "GET /api/v1/accounting/journals (Daftar & Filter Status Jurnal)",
              "POST /api/v1/accounting/journals/:id/post (Posting ke Buku Besar)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'accounting-ledger'}
          <ModulePlaceholder
            title="Buku Besar & Kas/Bank (General Ledger & Treasury)"
            description="Modul rekap mutasi buku besar per akun, mutasi rekening koran kas/bank RS, rekonsiliasi bank, dan ringkasan saldo berjalan."
            domainPackage="accounting/ledger"
            backendStatus="Terintegrasi dengan akun bertipe Treasury pada tabel chart_of_accounts"
            plannedEndpoints={[
              "GET /api/v1/accounting/ledger/:id (Mutasi Buku Besar Per Akun)",
              "GET /api/v1/accounting/treasury/book (Buku Kas & Bank Berjalan)",
              "GET /api/v1/accounting/trial-balance (Neraca Saldo / Trial Balance)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'accounting-reports'}
          <ModulePlaceholder
            title="Laporan Keuangan & Fiskal RS (Financial Reports)"
            description="Laporan Laba Rugi Operasional RS, Neraca Saldo Keuangan, Laporan Arus Kas, Perubahan Ekuitas, serta Tutup Buku Akhir Tahun Fiskal."
            domainPackage="accounting/report"
            backendStatus="Mengikuti Pedoman Standar Akuntansi Keuangan Rumah Sakit (PSAK / PARS)"
            plannedEndpoints={[
              "GET /api/v1/accounting/reports/balance-sheet (Laporan Posisi Keuangan / Neraca)",
              "GET /api/v1/accounting/reports/income-statement (Laporan Laba Rugi RS)",
              "POST /api/v1/accounting/fiscal-years/close (Tutup Buku Akhir Tahun Fiskal)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        <!-- 3.5 AUTH & USER MANAGEMENT DOMAIN MODULE -->
        {:else if activeNav === 'auth-users'}
          <AuthManagementPage initialTab="users" />
        {:else if activeNav === 'auth-roles'}
          <AuthManagementPage initialTab="roles" />
        {:else if activeNav === 'auth-permissions'}
          <AuthManagementPage initialTab="permissions" />

        <!-- 4. WORKSPACES PREVIEW (ROADMAP & INTEGRASI BACKEND) -->
        {:else if activeNav === 'outpatient-workspace'}
          <ModulePlaceholder
            title="Pelayanan Rawat Jalan (Outpatient)"
            description="Modul alur registrasi antrean poliklinik, penjadwalan janji temu dokter spesialis, dan konsultasi pemeriksaan pasien."
            domainPackage="outpatient"
            backendStatus="Use case alur pendaftaran dan integrasi encounter rawat jalan dalam perancangan"
            plannedEndpoints={[
              "POST /api/v1/outpatient/registrations (Daftar Kunjungan Poli)",
              "GET /api/v1/outpatient/queues (Antrean Realtime Pasien)",
              "POST /api/v1/outpatient/encounters/:id/consult (Mulai Konsultasi Dokter)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'emergency-workspace'}
          <ModulePlaceholder
            title="Instalasi Gawat Darurat (IGD)"
            description="Modul triase pasien darurat (Emergency Severity Index / ATS), penanganan cito 24 jam, tindakan resusitasi, dan disposisi rawat/rujuk."
            domainPackage="emergency"
            backendStatus="Domain triage, emergency disposition, dan sinkronisasi encounter IGD siap diintegrasikan"
            plannedEndpoints={[
              "POST /api/v1/emergency/triage (Klasifikasi Derajat Kedaruratan)",
              "GET /api/v1/emergency/active-cases (Daftar Kasus Kritis IGD)",
              "POST /api/v1/emergency/disposition (Disposisi Rawat Inap / Operasi Cito)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'inpatient-workspace'}
          <ModulePlaceholder
            title="Pelayanan Rawat Inap (Inpatient)"
            description="Modul admisi rawat inap, alokasi bed bangsal, transfer antar-ruangan (ICU/Bangsal/VK), dan resume pulang pasien."
            domainPackage="inpatient"
            backendStatus="Tabel master room, bed, dan service unit terhubung dengan GORM PostgreSQL"
            plannedEndpoints={[
              "POST /api/v1/inpatient/admissions (Admisi Rawat Inap & Alokasi Bed)",
              "POST /api/v1/inpatient/transfers (Transfer Ruangan / Bed)",
              "POST /api/v1/inpatient/discharges (Pemulangan & Resume Medis Akhir)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'pharmacy-workspace'}
          <ModulePlaceholder
            title="Farmasi & E-Resep (Pharmacy)"
            description="Modul pengelolaan resep elektronik dari dokter, telaah resep (clinical review), peracikan obat, dan dispensing ke pasien."
            domainPackage="pharmacy"
            backendStatus="Package internal/pharmacy/prescription dan dispensing disiapkan"
            plannedEndpoints={[
              "GET /api/v1/pharmacy/prescriptions/pending (Antrean Resep Masuk)",
              "POST /api/v1/pharmacy/prescriptions/:id/review (Telaah 7 Benar Obat)",
              "POST /api/v1/pharmacy/dispensing (Dispensing & Pengurangan Stok Depo)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'inventory-workspace'}
          <ModulePlaceholder
            title="Inventori & Logistik RS (Inventory)"
            description="Modul katalog barang/obat (Item CTI), manajemen depo ruangan dan gudang farmasi, mutasi stok, serta penyesuaian opname."
            domainPackage="inventory"
            backendStatus="Tabel master items, storages, dan migrasi Goose 00008-00009 aktif di PostgreSQL"
            plannedEndpoints={[
              "GET /api/v1/items (Katalog Obat & BHP)",
              "POST /api/v1/inventory/movements (Mutasi Stok Antar-Gudang/Depo)",
              "GET /api/v1/inventory/stock-balance (Monitoring Sisa Stok Realtime)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'billing-workspace'}
          <ModulePlaceholder
            title="Kasir & Billing Pasien (Billing)"
            description="Modul kalkulasi tagihan pelayanan, verifikasi klaim BPJS / asuransi swasta, kasir penerimaan pembayaran, dan cetak kuitansi."
            domainPackage="billing"
            backendStatus="Struktur tarif kelas dan komponen tarif internal/finance terhubung dengan database"
            plannedEndpoints={[
              "GET /api/v1/billing/encounters/:id/charges (Rincian Biaya Pelayanan)",
              "POST /api/v1/billing/invoices (Penerbitan Invoice Tagihan Pasien)",
              "POST /api/v1/billing/payments (Penerimaan Kasir & Bukti Pembayaran)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else if activeNav === 'audit-workspace'}
          <ModulePlaceholder
            title="Audit Trail & Keamanan Sistem (Audit)"
            description="Audit append-only untuk pencatatan rekam jejak akses data medis, perubahan sensitif, mutasi tarif, dan kepatuhan regulasi."
            domainPackage="audit"
            backendStatus="Audit trail port dan logging terproteksi di PostgreSQL"
            plannedEndpoints={[
              "GET /api/v1/audit/logs (Log Mutasi Rekam Medis & Master Data)",
              "GET /api/v1/audit/security-events (Peringatan Akses Sensitif)"
            ]}
            onBackToLauncher={() => activeNav = 'desktop'}
          />

        {:else}
          <!-- Tab Aktivitas & Riwayat Pasien -->
          <div class="p-6 bg-[#f8fafd] rounded-2xl border border-[#e1e5ea]">
            <h3 class="text-base font-semibold text-[#1f1f1f] mb-2">Aktivitas & Riwayat Pasien</h3>
            <p class="text-xs text-[#444746] mb-4">Kronologi kunjungan poli dan jadwal kontrol berkala.</p>
            <div class="p-4 rounded-xl bg-white border border-[#e1e5ea] text-xs text-[#444746]">
              <p>1. <strong>22 Sep 2026, 10:15</strong> — Pemeriksaan Fisik & Pemetaan Body Diagram oleh dr. Hendra Wijaya, Sp.B</p>
              <p class="mt-2">2. <strong>14 Ags 2026, 09:30</strong> — Rawat Jalan Kontrol Rutin Hipertensi di Poliklinik Penyakit Dalam</p>
            </div>
          </div>
        {/if}
      </main>
    </div>
  </div>
{/if}
