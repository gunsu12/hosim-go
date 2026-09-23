<script lang="ts">
  import { CheckCircle } from '@lucide/svelte';
  import TopBar from '$lib/components/TopBar.svelte';
  import PatientBanner from '$lib/components/PatientBanner.svelte';
  import LoginPage from '$lib/components/LoginPage.svelte';
  import NavigationRail from '$lib/components/NavigationRail.svelte';
  import { auth } from '$lib/stores/auth.svelte';
  import type { BodyFinding, ClinicalNotes } from '$lib/types';
  import { getNavFromCurrentHash, setHashFromNav } from '$lib/router';

  // Domain Clinical Pages
  import {
    PhysicalExamPage,
    SoapPage,
    OdontogramPage,
    LabPage,
    RadiologyPage,
    PharmacyPage
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
    TariffClassPage
  } from '$lib/pages/master';

  // Inisialisasi langsung dari window.location.hash browser agar persisten saat di-refresh (F5)
  let activeNav = $state<string>(getNavFromCurrentHash());
  let isSidebarExpanded = $state<boolean>(false);
  let saveSuccessAlert = $state<boolean>(false);

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
    <!-- Top Bar (Google Drive Header dengan Hamburger Menu & Profil) -->
    <TopBar
      onToggleSidebar={() => isSidebarExpanded = !isSidebarExpanded}
      user={auth.user}
      onLogout={() => auth.logout()}
    />

    <!-- Main Layout (Compact Side Rail + Spacious White Canvas) -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Google M3 Navigation Rail dengan Dukungan Sub-Menu (Flyout & Accordion) -->
      <NavigationRail
        bind:activeNav
        bind:isExpanded={isSidebarExpanded}
        onNewEncounter={handleSaveExamination}
      />

      <!-- Google Drive Spacious White Content Island -->
      <main class="flex-1 bg-white md:rounded-3xl border border-[#e1e5ea] shadow-xs md:mr-3 md:mb-3 overflow-y-auto p-5 md:p-6 flex flex-col gap-5">
        <!-- Patient Summary Banner (Hanya ditampilkan pada modul rekam medis / tindakan klinis pasien) -->
        {#if !activeNav.startsWith('master-')}
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

        <!-- Clinical Modules -->
        {#if activeNav === 'physical'}
          <PhysicalExamPage
            bind:findings
            {clinicalNotes}
            onSave={handleSaveExamination}
          />

        {:else if activeNav === 'anamnesis'}
          <SoapPage
            bind:clinicalNotes
            onSave={handleSaveExamination}
          />

        {:else if activeNav === 'odontogram'}
          <OdontogramPage />

        {:else if activeNav === 'lab'}
          <LabPage />

        {:else if activeNav === 'radiology'}
          <RadiologyPage />

        {:else if activeNav === 'pharmacy'}
          <PharmacyPage />

        <!-- Master Data Domain Modules -->
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
