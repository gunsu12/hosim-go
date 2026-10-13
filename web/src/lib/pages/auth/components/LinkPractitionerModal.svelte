<script lang="ts">
  import {
    Stethoscope,
    Search,
    CheckCircle2,
    Unlink,
    AlertCircle,
    UserCheck,
    Loader2,
    ShieldAlert,
    ChevronDown,
    X
  } from '@lucide/svelte';
  import M3Modal from '../../../components/m3/M3Modal.svelte';
  import M3Button from '../../../components/m3/M3Button.svelte';
  import { linkPractitioner, unlinkPractitioner } from '$lib/api/auth';
  import { getPractitioners, getPractitionerById } from '$lib/api/practitioners';
  import type { UserRecord } from '../../../types/auth';
  import type { PractitionerRecord } from '../../../types/master/practitioner';

  interface Props {
    isOpen: boolean;
    user: UserRecord | null;
    practitioners?: PractitionerRecord[]; // Initial fallback cache dari parent
    onClose: () => void;
    onSuccess: (updatedUser: UserRecord) => void;
  }

  let {
    isOpen = false,
    user = null,
    practitioners = [],
    onClose = () => {},
    onSuccess = () => {}
  }: Props = $props();

  const PAGE_SIZE = 20;

  // Selected State
  let selectedId = $state<string>('');
  let selectedPractitioner = $state<PractitionerRecord | null>(null);
  let currentLinkedPractitioner = $state<PractitionerRecord | null>(null);

  // Search & Pagination State
  let searchQuery = $state<string>('');
  let currentPage = $state<number>(1);
  let totalCount = $state<number>(0);
  let practitionerList = $state<PractitionerRecord[]>([]);

  // Loading & Submission State
  let isLoadingList = $state<boolean>(false);
  let isLoadingMore = $state<boolean>(false);
  let isSubmitting = $state<boolean>(false);
  let isUnlinking = $state<boolean>(false);
  let errorMessage = $state<string | null>(null);

  let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null;

  let hasMore = $derived(practitionerList.length < totalCount);

  // Sync saat modal dibuka
  $effect(() => {
    if (isOpen && user) {
      selectedId = user.practitioner_id || '';
      searchQuery = '';
      currentPage = 1;
      errorMessage = null;
      practitionerList = [];
      totalCount = 0;

      // Ambil detail practitioner yang sedang aktif tertaut
      resolveCurrentLinkedPractitioner();

      // Pemuatan data halaman pertama (dibatasi 20 item untuk efisiensi DOM & memori)
      fetchPractitioners(1, '', false);
    }
  });

  async function resolveCurrentLinkedPractitioner() {
    if (!user?.practitioner_id) {
      currentLinkedPractitioner = null;
      selectedPractitioner = null;
      return;
    }

    // Cek cache props terlebih dahulu
    const cached = practitioners.find((p) => p.id === user.practitioner_id);
    if (cached) {
      currentLinkedPractitioner = cached;
      selectedPractitioner = cached;
      return;
    }

    // Jika tidak ada di cache awal, fetch langsung by ID dari backend
    try {
      const fetched = await getPractitionerById(user.practitioner_id);
      currentLinkedPractitioner = fetched;
      selectedPractitioner = fetched;
    } catch {
      currentLinkedPractitioner = null;
    }
  }

  async function fetchPractitioners(page: number, query: string, append: boolean = false) {
    if (append) {
      isLoadingMore = true;
    } else {
      isLoadingList = true;
    }

    try {
      const res = await getPractitioners({
        search: query.trim() || undefined,
        page,
        limit: PAGE_SIZE
      });

      const items = res.data || [];
      totalCount = res.count || items.length;
      currentPage = page;

      if (append) {
        // Gabungkan dan deduplikasi ID
        const existingIds = new Set(practitionerList.map((p) => p.id));
        const newUnique = items.filter((p) => !existingIds.has(p.id));
        practitionerList = [...practitionerList, ...newUnique];
      } else {
        practitionerList = items;
      }
    } catch (err: any) {
      errorMessage = err?.message || 'Gagal memuat katalog tenaga medis dari server.';
    } finally {
      isLoadingList = false;
      isLoadingMore = false;
    }
  }

  function handleSearchInput(e: Event) {
    const val = (e.target as HTMLInputElement).value;
    searchQuery = val;

    if (searchDebounceTimer) clearTimeout(searchDebounceTimer);
    searchDebounceTimer = setTimeout(() => {
      fetchPractitioners(1, val, false);
    }, 300);
  }

  function clearSearch() {
    searchQuery = '';
    fetchPractitioners(1, '', false);
  }

  function handleLoadMore() {
    if (isLoadingMore || !hasMore) return;
    fetchPractitioners(currentPage + 1, searchQuery, true);
  }

  function selectPractitioner(p: PractitionerRecord) {
    selectedId = p.id;
    selectedPractitioner = p;
  }

  async function handleLink() {
    if (!user) return;
    if (!selectedId) {
      errorMessage = 'Pilih salah satu tenaga medis terlebih dahulu.';
      return;
    }
    if (selectedId === user.practitioner_id) {
      onClose();
      return;
    }

    isSubmitting = true;
    errorMessage = null;
    try {
      const updated = await linkPractitioner(user.id, selectedId);
      onSuccess(updated);
      onClose();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal menghubungkan tenaga medis ke akun pengguna.';
    } finally {
      isSubmitting = false;
    }
  }

  async function handleUnlink() {
    if (!user || !user.practitioner_id) return;
    isUnlinking = true;
    errorMessage = null;
    try {
      const updated = await unlinkPractitioner(user.id);
      onSuccess(updated);
      onClose();
    } catch (err: any) {
      errorMessage = err.message || 'Gagal memutuskan tautan tenaga medis.';
    } finally {
      isUnlinking = false;
    }
  }
</script>

<M3Modal
  {isOpen}
  title="Tautkan Tenaga Medis (Practitioner)"
  subtitle={user ? `Akun: ${user.name} (@${user.username}) • Peran: ${user.role_name || user.role}` : ''}
  icon={Stethoscope}
  iconClass="bg-teal-50 text-teal-600 dark:bg-teal-950/60 dark:text-teal-400"
  size="xl"
  {onClose}
>
  <div class="flex flex-col gap-4 text-xs">
    {#if errorMessage}
      <div class="p-3.5 rounded-2xl bg-rose-50 dark:bg-rose-950/40 border border-rose-200 dark:border-rose-800 text-rose-700 dark:text-rose-300 flex items-start gap-2.5">
        <AlertCircle class="w-4 h-4 shrink-0 mt-0.5" />
        <div>
          <span class="font-bold">Gagal: </span>
          <span>{errorMessage}</span>
        </div>
      </div>
    {/if}

    <!-- Info banner pertanggungjawaban medis -->
    <div class="p-3.5 rounded-2xl bg-[#f0f4f9] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] flex items-start gap-3">
      <ShieldAlert class="w-4 h-4 text-[#0b57d0] dark:text-[#a8c7fa] shrink-0 mt-0.5" />
      <div class="text-[#444746] dark:text-[#c4c7c5] leading-relaxed">
        <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">Aturan Hukum & Rekam Medis: </span>
        Pemetaan ini mengikat akun login dengan identitas resmi nakes. Seluruh catatan rekam medis (EHR), asesmen, dan resep yang dibuat akan dipertanggungjawabkan atas nama tenaga medis terpilih serta disinkronkan ke SATUSEHAT.
      </div>
    </div>

    <!-- Status saat ini jika sudah tertaut -->
    {#if currentLinkedPractitioner}
      <div class="p-4 rounded-2xl bg-teal-50/70 dark:bg-teal-950/30 border border-teal-200 dark:border-teal-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-teal-100 dark:bg-teal-900/60 text-teal-700 dark:text-teal-300 flex items-center justify-center shrink-0">
            <UserCheck class="w-5 h-5" />
          </div>
          <div>
            <div class="text-[11px] font-bold uppercase tracking-wider text-teal-700 dark:text-teal-300">
              Saat Ini Tertaut ke:
            </div>
            <div class="text-sm font-bold text-[#1f1f1f] dark:text-[#e3e3e3]">
              {currentLinkedPractitioner.name}
            </div>
            <div class="text-[11px] text-[#444746] dark:text-[#c4c7c5] font-mono">
              SIP: {currentLinkedPractitioner.sip || '-'} • STR: {currentLinkedPractitioner.str || '-'}
            </div>
          </div>
        </div>
        <M3Button
          variant="outlined"
          onclick={handleUnlink}
          disabled={isUnlinking || isSubmitting}
          class="!border-rose-200 dark:!border-rose-800 !text-rose-600 dark:!text-rose-400 hover:!bg-rose-50 dark:hover:!bg-rose-950/50 shrink-0"
        >
          {#if isUnlinking}
            <Loader2 class="w-3.5 h-3.5 animate-spin mr-1" />
          {:else}
            <Unlink class="w-3.5 h-3.5 mr-1" />
          {/if}
          <span>Putuskan Tautan</span>
        </M3Button>
      </div>
    {/if}

    <!-- Search box dengan Server-Side Debounce -->
    <div>
      <div class="flex items-center justify-between mb-1">
        <label for="search-practitioner" class="font-semibold text-[#1f1f1f] dark:text-[#e3e3e3]">
          Pilih Tenaga Medis dari Katalog
        </label>
        {#if !isLoadingList && totalCount > 0}
          <span class="text-[11px] text-[#747775] dark:text-[#8e918f]">
            Menampilkan {practitionerList.length} dari {totalCount} tenaga medis
          </span>
        {/if}
      </div>

      <div class="relative">
        <Search class="w-4 h-4 text-[#747775] dark:text-[#8e918f] absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          id="search-practitioner"
          type="text"
          value={searchQuery}
          oninput={handleSearchInput}
          placeholder="Cari berdasarkan nama dokter/perawat, SIP, STR, atau profesi..."
          class="w-full h-10 pl-10 pr-9 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
        />
        {#if searchQuery}
          <button
            type="button"
            onclick={clearSearch}
            class="absolute right-3 top-1/2 -translate-y-1/2 text-[#747775] hover:text-[#1f1f1f] dark:hover:text-[#e3e3e3] p-0.5"
            title="Hapus pencarian"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        {/if}
      </div>
    </div>

    <!-- Daftar Tenaga Medis (Dibatasi 20 per request untuk mencegah membludak di DOM) -->
    <div class="max-h-64 overflow-y-auto rounded-2xl border border-[#e1e5ea] dark:border-[#353738] divide-y divide-[#e1e5ea] dark:divide-[#353738] bg-white dark:bg-[#1e1f20]">
      {#if isLoadingList}
        <div class="p-8 text-center text-[#747775] dark:text-[#8e918f] flex flex-col items-center justify-center gap-2">
          <Loader2 class="w-6 h-6 animate-spin text-[#0b57d0] dark:text-[#a8c7fa]" />
          <span>Mencari tenaga medis di server...</span>
        </div>
      {:else if practitionerList.length === 0}
        <div class="p-8 text-center text-[#747775] dark:text-[#8e918f]">
          {searchQuery
            ? `Tidak ada data tenaga medis yang cocok dengan kata kunci "${searchQuery}".`
            : 'Belum ada data tenaga medis terdaftar.'}
        </div>
      {:else}
        {#each practitionerList as p (p.id)}
          {@const isSelected = selectedId === p.id}
          {@const professionTitle = typeof p.profession === 'object' && p.profession ? p.profession.name : (p.profession || 'Tenaga Medis')}
          {@const ihs = p.ihs_practitioner_id || p.ihs_number || p.ihsId}
          <button
            type="button"
            onclick={() => selectPractitioner(p)}
            class="w-full text-left p-3 flex items-center justify-between gap-3 transition-colors cursor-pointer {isSelected ? 'bg-teal-50 dark:bg-teal-950/50 text-teal-900 dark:text-teal-100' : 'hover:bg-[#f8fafd] dark:hover:bg-[#28292a]/60 text-[#1f1f1f] dark:text-[#e3e3e3]'}"
          >
            <div class="flex items-center gap-3 min-w-0">
              <div class="w-8 h-8 rounded-xl flex items-center justify-center shrink-0 {isSelected ? 'bg-teal-600 text-white' : 'bg-[#e8f0fe] dark:bg-[#1e2738] text-[#0b57d0] dark:text-[#a8c7fa]'}">
                <Stethoscope class="w-4 h-4" />
              </div>
              <div class="min-w-0">
                <div class="font-bold truncate text-xs">{p.name}</div>
                <div class="text-[11px] text-[#747775] dark:text-[#8e918f] truncate">
                  {professionTitle} {p.sip ? `• SIP: ${p.sip}` : ''} {ihs ? `• IHS: ${ihs}` : ''}
                </div>
              </div>
            </div>
            <div class="shrink-0 flex items-center gap-2">
              {#if isSelected}
                <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-teal-100 dark:bg-teal-900 text-teal-800 dark:text-teal-200">
                  <CheckCircle2 class="w-3 h-3" />
                  Dipilih
                </span>
              {/if}
            </div>
          </button>
        {/each}

        <!-- Pagination Load More jika total ratusan -->
        {#if hasMore}
          <div class="p-2.5 text-center bg-[#f8fafd] dark:bg-[#28292a]/50">
            <button
              type="button"
              onclick={handleLoadMore}
              disabled={isLoadingMore}
              class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold text-[#0b57d0] dark:text-[#a8c7fa] hover:bg-[#e8f0fe] dark:hover:bg-[#1e2738] transition-colors cursor-pointer disabled:opacity-50"
            >
              {#if isLoadingMore}
                <Loader2 class="w-3.5 h-3.5 animate-spin" />
                <span>Memuat data berikutnya...</span>
              {:else}
                <ChevronDown class="w-3.5 h-3.5" />
                <span>Muat {PAGE_SIZE} Tenaga Medis Lainnya (Sisa {totalCount - practitionerList.length})</span>
              {/if}
            </button>
          </div>
        {/if}
      {/if}
    </div>
  </div>

  {#snippet footer()}
    <M3Button variant="text" onclick={onClose} disabled={isSubmitting || isUnlinking}>
      Batal
    </M3Button>
    <M3Button
      variant="filled"
      onclick={handleLink}
      disabled={isSubmitting || isUnlinking || !selectedId || selectedId === user?.practitioner_id}
      class="!bg-teal-600 hover:!bg-teal-700 !text-white"
    >
      {#if isSubmitting}
        <Loader2 class="w-4 h-4 animate-spin mr-1.5" />
        <span>Menyimpan...</span>
      {:else}
        <CheckCircle2 class="w-4 h-4 mr-1.5" />
        <span>Simpan Tautan</span>
      {/if}
    </M3Button>
  {/snippet}
</M3Modal>
