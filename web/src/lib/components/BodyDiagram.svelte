<script lang="ts">
  import {
    Activity,
    AlertCircle,
    Flame,
    Eye,
    Plus,
    Trash2,
    RotateCcw,
    CheckCircle2,
    Sparkles,
    User,
    ChevronRight
  } from '@lucide/svelte';
  import M3Card from './m3/M3Card.svelte';
  import M3Chip from './m3/M3Chip.svelte';
  import M3Button from './m3/M3Button.svelte';
  import M3SegmentedButton, { type SegmentOption } from './m3/M3SegmentedButton.svelte';
  import type { BodyFinding, FindingCategory } from '../types';

  interface Props {
    findings?: BodyFinding[];
    onSave?: (findings: BodyFinding[]) => void;
  }

  let { findings = $bindable([]), onSave = () => {} }: Props = $props();

  let activeView = $state<'front' | 'back'>('front');
  let selectedCategory = $state<FindingCategory>('pain');
  let selectedMarkerId = $state<number | null>(null);
  let isEditorOpen = $state<boolean>(false);

  // Form input untuk titik baru
  let tempMarker = $state<BodyFinding>({
    id: 0,
    x: 0,
    y: 0,
    view: 'front',
    category: 'pain',
    severity: 5,
    note: '',
    createdAt: ''
  });

  export interface CategoryOption {
    id: FindingCategory;
    label: string;
    color: 'primary' | 'pain' | 'injury' | 'fracture' | 'edema';
    hex: string;
    badge: string;
  }

  const categories: CategoryOption[] = [
    { id: 'pain', label: 'Nyeri (Pain)', color: 'pain', hex: '#E37400', badge: 'bg-[#E37400]' },
    { id: 'injury', label: 'Luka Robek / Vulnus', color: 'injury', hex: '#D93025', badge: 'bg-[#D93025]' },
    { id: 'fracture', label: 'Fraktur / Deformitas', color: 'fracture', hex: '#9334E6', badge: 'bg-[#9334E6]' },
    { id: 'edema', label: 'Memar / Hematoma', color: 'edema', hex: '#1A73E8', badge: 'bg-[#1A73E8]' },
    { id: 'mass', label: 'Benjolan / Massa', color: 'primary', hex: '#137333', badge: 'bg-[#137333]' }
  ];

  const viewOptions: SegmentOption[] = [
    { value: 'front', label: 'Tampak Depan (Anterior)' },
    { value: 'back', label: 'Tampak Belakang (Posterior)' }
  ];

  function handleSvgClick(e: MouseEvent & { currentTarget: HTMLElement }): void {
    const rect = e.currentTarget.getBoundingClientRect();
    const x = ((e.clientX - rect.left) / rect.width) * 100;
    const y = ((e.clientY - rect.top) / rect.height) * 100;

    tempMarker = {
      id: Date.now(),
      x: Math.round(x * 10) / 10,
      y: Math.round(y * 10) / 10,
      view: activeView,
      category: selectedCategory,
      severity: selectedCategory === 'pain' ? 5 : 0,
      note: '',
      createdAt: new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
    };

    selectedMarkerId = tempMarker.id;
    isEditorOpen = true;
  }

  function saveCurrentMarker(): void {
    if (!tempMarker.note.trim()) {
      tempMarker.note = getCategoryDetails(tempMarker.category).label;
    }
    
    const index = findings.findIndex(f => f.id === tempMarker.id);
    if (index >= 0) {
      findings[index] = { ...tempMarker };
    } else {
      findings = [...findings, { ...tempMarker }];
    }
    isEditorOpen = false;
    if (onSave) onSave(findings);
  }

  function deleteMarker(id: number, e?: MouseEvent): void {
    if (e) e.stopPropagation();
    findings = findings.filter(f => f.id !== id);
    if (selectedMarkerId === id) {
      selectedMarkerId = null;
      isEditorOpen = false;
    }
    if (onSave) onSave(findings);
  }

  function selectExistingMarker(marker: BodyFinding, e: MouseEvent): void {
    e.stopPropagation();
    selectedMarkerId = marker.id;
    tempMarker = { ...marker };
    isEditorOpen = true;
  }

  function getCategoryDetails(catId: FindingCategory): CategoryOption {
    return categories.find(c => c.id === catId) || categories[0];
  }

  let activeViewFindings = $derived(
    findings.filter(f => f.view === activeView)
  );
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-5">
  <!-- Area Kiri: Diagram Anatomi Tubuh -->
  <div class="lg:col-span-8 flex flex-col gap-4">
    <!-- Toolbar Kontrol Diagram Google Drive Style -->
    <div class="p-3 bg-white rounded-2xl border border-[#e1e5ea] flex flex-wrap items-center justify-between gap-3 shadow-xs">
      <div class="flex items-center gap-2">
        <M3SegmentedButton options={viewOptions} bind:selected={activeView} />
      </div>

      <div class="flex items-center gap-2.5">
        <span class="text-xs text-[#444746]">Klik anatomi untuk menandai</span>
        <span class="inline-flex items-center justify-center px-2 py-0.5 rounded-full text-xs font-semibold bg-[#d3e3fd] text-[#041e49]">
          {findings.length} Titik
        </span>
      </div>
    </div>

    <!-- Filter Chips Kategori (Google Drive Style) -->
    <div class="flex flex-wrap items-center gap-2 px-0.5">
      <span class="text-xs font-semibold uppercase tracking-wider text-[#444746] mr-1">Tipe:</span>
      {#each categories as cat}
        <M3Chip
          selected={selectedCategory === cat.id}
          color={cat.color}
          onclick={() => selectedCategory = cat.id}
        >
          <span class="w-2.5 h-2.5 rounded-full {cat.badge}"></span>
          <span>{cat.label}</span>
        </M3Chip>
      {/each}
    </div>

    <!-- Kanvas Diagram Anatomi SVG -->
    <div class="relative p-6 bg-white rounded-2xl border border-[#e1e5ea] flex items-center justify-center overflow-hidden min-h-[560px] shadow-xs">
      
      <!-- Label Sudut Pandang -->
      <div class="absolute top-4 left-4 z-10 flex items-center gap-2 bg-[#f0f4f9] px-3.5 py-1.5 rounded-full border border-[#e1e5ea] text-xs font-medium text-[#1f1f1f]">
        <User class="w-3.5 h-3.5 text-[#0b57d0]" />
        <span>{activeView === 'front' ? 'Tampak Depan (Anterior)' : 'Tampak Belakang (Posterior)'}</span>
      </div>

      <!-- Container Interaktif SVG -->
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="relative w-full max-w-[340px] aspect-[1/2] cursor-crosshair select-none"
        onclick={handleSvgClick}
      >
        <!-- SVG Siluet Tubuh Anatomi Bersih & Elegan -->
        <svg
          viewBox="0 0 200 400"
          class="w-full h-full drop-shadow-xs transition-opacity duration-200 pointer-events-none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <defs>
            <linearGradient id="driveBodyGrad" x1="0%" y1="0%" x2="0%" y2="100%">
              <stop offset="0%" stop-color="#f1f3f4" />
              <stop offset="100%" stop-color="#e3e7ea" />
            </linearGradient>
          </defs>

          {#if activeView === 'front'}
            <!-- === TAMPAK DEPAN (ANTERIOR) === -->
            <!-- Kepala -->
            <path
              d="M100,20 C85,20 82,35 82,48 C82,62 88,72 100,72 C112,72 118,62 118,48 C118,35 115,20 100,20 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Leher -->
            <path
              d="M93,69 L93,82 L107,82 L107,69 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Dada, Perut & Panggul -->
            <path
              d="M93,82 C75,83 62,92 56,102 C50,112 48,135 48,160 C48,172 50,185 54,198 C58,212 62,225 68,235 C74,242 82,248 100,248 C118,248 126,242 132,235 C138,225 142,212 146,198 C150,185 152,172 152,160 C152,135 150,112 144,102 C138,92 125,83 107,82 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Lengan Kanan Pasien -->
            <path
              d="M56,102 C48,110 38,130 35,155 C32,178 30,205 28,225 C26,238 24,248 22,255 C20,260 22,266 26,266 C30,266 33,258 36,248 C39,235 41,215 44,195 C46,180 48,168 50,155 C52,142 54,120 56,102 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Lengan Kiri Pasien -->
            <path
              d="M144,102 C152,110 162,130 165,155 C168,178 170,205 172,225 C174,238 176,248 178,255 C180,260 178,266 174,266 C170,266 167,258 164,248 C161,235 159,215 156,195 C154,180 152,168 150,155 C148,142 146,120 144,102 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Kaki Kanan Pasien -->
            <path
              d="M72,244 C68,258 66,280 65,305 C64,325 63,350 63,375 C63,382 60,387 55,389 C52,390 52,394 58,394 C68,394 72,390 73,384 C76,365 77,340 79,315 C81,290 84,268 87,247 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Kaki Kiri Pasien -->
            <path
              d="M128,244 C132,258 134,280 135,305 C136,325 137,350 137,375 C137,382 140,387 145,389 C148,390 148,394 142,394 C132,394 128,390 127,384 C124,365 123,340 121,315 C119,290 116,268 113,247 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />

            <!-- Garis Anatomi Ringan -->
            <path d="M85,92 C93,96 100,96 107,92" stroke="#bdc1c6" stroke-width="1" fill="none" stroke-linecap="round" />
            <path d="M72,125 C82,130 92,128 97,125" stroke="#dadce0" stroke-width="1" fill="none" />
            <path d="M128,125 C118,130 108,128 103,125" stroke="#dadce0" stroke-width="1" fill="none" />
            <circle cx="100" cy="188" r="2.5" fill="#bdc1c6" />
          {:else}
            <!-- === TAMPAK BELAKANG (POSTERIOR) === -->
            <!-- Kepala Belakang -->
            <path
              d="M100,20 C85,20 82,35 82,48 C82,62 88,72 100,72 C112,72 118,62 118,48 C118,35 115,20 100,20 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Leher & Punggung Atas -->
            <path
              d="M93,69 L93,82 L107,82 L107,69 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Punggung & Pantat -->
            <path
              d="M93,82 C75,83 62,92 56,102 C50,112 48,135 48,160 C48,172 50,185 54,198 C58,212 62,225 68,235 C74,242 82,248 100,248 C118,248 126,242 132,235 C138,225 142,212 146,198 C150,185 152,172 152,160 C152,135 150,112 144,102 C138,92 125,83 107,82 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Lengan Kanan Posterior -->
            <path
              d="M56,102 C48,110 38,130 35,155 C32,178 30,205 28,225 C26,238 24,248 22,255 C20,260 22,266 26,266 C30,266 33,258 36,248 C39,235 41,215 44,195 C46,180 48,168 50,155 C52,142 54,120 56,102 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Lengan Kiri Posterior -->
            <path
              d="M144,102 C152,110 162,130 165,155 C168,178 170,205 172,225 C174,238 176,248 178,255 C180,260 178,266 174,266 C170,266 167,258 164,248 C161,235 159,215 156,195 C154,180 152,168 150,155 C148,142 146,120 144,102 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Kaki Kanan Posterior -->
            <path
              d="M72,244 C68,258 66,280 65,305 C64,325 63,350 63,375 C63,382 60,387 55,389 C52,390 52,394 58,394 C68,394 72,390 73,384 C76,365 77,340 79,315 C81,290 84,268 87,247 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />
            <!-- Kaki Kiri Posterior -->
            <path
              d="M128,244 C132,258 134,280 135,305 C136,325 137,350 137,375 C137,382 140,387 145,389 C148,390 148,394 142,394 C132,394 128,390 127,384 C124,365 123,340 121,315 C119,290 116,268 113,247 Z"
              fill="url(#driveBodyGrad)"
              stroke="#9aa0a6"
              stroke-width="1.5"
            />

            <!-- Garis Tulang Belakang & Skapula -->
            <line x1="100" y1="85" x2="100" y2="230" stroke="#bdc1c6" stroke-width="1.2" stroke-dasharray="2 3" />
            <path d="M75,108 C80,120 84,130 76,140" stroke="#dadce0" stroke-width="1.2" fill="none" />
            <path d="M125,108 C120,120 116,130 124,140" stroke="#dadce0" stroke-width="1.2" fill="none" />
          {/if}
        </svg>

        <!-- Titik Marker Klinis -->
        {#each activeViewFindings as marker, idx}
          {@const cat = getCategoryDetails(marker.category)}
          {@const isSelected = selectedMarkerId === marker.id}

          <!-- svelte-ignore a11y_click_events_have_key_events -->
          <!-- svelte-ignore a11y_no_static_element_interactions -->
          <div
            class="absolute -translate-x-1/2 -translate-y-1/2 z-20 transition-transform duration-150 cursor-pointer group {isSelected ? 'scale-125 z-30' : 'hover:scale-115'}"
            style="left: {marker.x}%; top: {marker.y}%;"
            onclick={(e) => selectExistingMarker(marker, e)}
          >
            <!-- Google Style Pulsing Ring -->
            {#if isSelected}
              <span class="absolute -inset-2 rounded-full animate-ping opacity-35 bg-[#0b57d0]"></span>
            {/if}

            <div
              class="relative flex items-center justify-center w-6 h-6 rounded-full text-white text-[11px] font-bold shadow-md border-2 border-white {cat.badge}"
            >
              {idx + 1}
            </div>

            <!-- Tooltip -->
            <div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-1.5 hidden group-hover:flex flex-col items-center pointer-events-none z-40 whitespace-nowrap">
              <div class="bg-[#1f1f1f] text-white text-[11px] font-medium px-2.5 py-1 rounded-md shadow-lg">
                <span class="font-bold">#{idx + 1} {cat.label}</span>
                {#if marker.note}
                  <span class="block text-[10px] text-[#e3e3e3] max-w-[150px] truncate">{marker.note}</span>
                {/if}
              </div>
              <div class="w-1.5 h-1.5 bg-[#1f1f1f] rotate-45 -mt-0.5"></div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>

  <!-- Area Kanan: Editor & Daftar Temuan -->
  <div class="lg:col-span-4 flex flex-col gap-4">
    {#if isEditorOpen}
      <!-- Card Editor Google Drive Style -->
      <div class="p-5 bg-white rounded-2xl border border-[#0b57d0]/40 shadow-sm animate-in fade-in duration-150">
        <div class="flex items-center justify-between pb-3 border-b border-[#e1e5ea]">
          <div class="flex items-center gap-2">
            <span class="w-3 h-3 rounded-full {getCategoryDetails(tempMarker.category).badge}"></span>
            <h3 class="font-semibold text-sm text-[#1f1f1f]">Detail Temuan Fisik</h3>
          </div>
          <span class="text-[11px] text-[#444746] font-mono">X:{tempMarker.x}% Y:{tempMarker.y}%</span>
        </div>

        <div class="mt-4 flex flex-col gap-3.5">
          <div>
            <label for="marker-category-select" class="block text-xs font-semibold text-[#444746] mb-1.5">Klasifikasi Temuan</label>
            <select
              id="marker-category-select"
              bind:value={tempMarker.category}
              class="w-full h-10 px-3 text-xs rounded-xl bg-[#f0f4f9] border border-[#e1e5ea] text-[#1f1f1f] focus:outline-none focus:ring-2 focus:ring-[#0b57d0] focus:bg-white"
            >
              {#each categories as cat}
                <option value={cat.id}>{cat.label}</option>
              {/each}
            </select>
          </div>

          {#if tempMarker.category === 'pain'}
            <div>
              <div class="flex justify-between items-center mb-1">
                <label for="marker-severity-range" class="text-xs font-semibold text-[#444746]">Skala Nyeri (NRS)</label>
                <span class="text-xs font-bold text-[#b06000] bg-[#fff7ee] border border-[#ffba5a]/50 px-2.5 py-0.5 rounded-full">{tempMarker.severity} / 10</span>
              </div>
              <input
                id="marker-severity-range"
                type="range"
                min="1"
                max="10"
                bind:value={tempMarker.severity}
                class="w-full accent-[#e37400] cursor-pointer"
              />
              <div class="flex justify-between text-[10px] text-[#747775]">
                <span>Ringan (1-3)</span>
                <span>Sedang (4-6)</span>
                <span>Berat (7-10)</span>
              </div>
            </div>
          {/if}

          <div>
            <label for="marker-note-textarea" class="block text-xs font-semibold text-[#444746] mb-1.5">Deskripsi Klinis / Catatan Dokter</label>
            <textarea
              id="marker-note-textarea"
              bind:value={tempMarker.note}
              rows="3"
              placeholder="Contoh: Nyeri tekan kuadran kanan bawah, eritema 2x3 cm..."
              class="w-full p-2.5 text-xs rounded-xl bg-[#f0f4f9] border border-[#e1e5ea] text-[#1f1f1f] placeholder:text-[#747775] focus:outline-none focus:ring-2 focus:ring-[#0b57d0] focus:bg-white resize-none"
            ></textarea>
          </div>

          <div class="flex items-center justify-between gap-2 pt-2 border-t border-[#e1e5ea]">
            <button
              type="button"
              onclick={(e) => deleteMarker(tempMarker.id, e)}
              class="p-2 text-[#b3261e] hover:bg-[#f9dedc]/40 rounded-lg transition-colors cursor-pointer"
              title="Hapus Titik"
            >
              <Trash2 class="w-4 h-4" />
            </button>

            <div class="flex items-center gap-2">
              <M3Button variant="text" onclick={() => isEditorOpen = false}>
                Batal
              </M3Button>
              <M3Button variant="filled" onclick={saveCurrentMarker}>
                Simpan Titik
              </M3Button>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <!-- Daftar Semua Temuan Rekam Medis -->
    <div class="p-5 bg-white rounded-2xl border border-[#e1e5ea] shadow-xs flex-1 flex flex-col">
      <div class="flex items-center justify-between pb-3 border-b border-[#e1e5ea]">
        <div>
          <h3 class="font-semibold text-sm text-[#1f1f1f]">Daftar Temuan Fisik</h3>
          <p class="text-xs text-[#444746]">Tersinkronisasi dengan form rekam medis</p>
        </div>
        {#if findings.length > 0}
          <button
            type="button"
            onclick={() => findings = []}
            class="text-[11px] text-[#b3261e] hover:underline cursor-pointer"
          >
            Reset Semua
          </button>
        {/if}
      </div>

      <div class="mt-3 flex flex-col gap-2.5 overflow-y-auto max-h-[420px] pr-1">
        {#if findings.length === 0}
          <div class="py-12 flex flex-col items-center justify-center text-center text-[#444746]">
            <User class="w-8 h-8 text-[#c4c7c5] mb-2" />
            <p class="text-xs font-medium">Belum ada titik temuan</p>
            <p class="text-[11px] text-[#747775] mt-0.5">Klik area tubuh pada diagram untuk menambahkan catatan fisik.</p>
          </div>
        {:else}
          {#each findings as f, i}
            {@const cat = getCategoryDetails(f.category)}
            {@const isSelected = selectedMarkerId === f.id}

            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div
              class="p-3 rounded-xl border transition-all duration-150 cursor-pointer {isSelected ? 'bg-[#e8f0fe] border-[#0b57d0]' : 'bg-[#f8fafd] border-[#e1e5ea] hover:bg-[#f0f4f9]'}"
              onclick={() => {
                activeView = f.view;
                selectedMarkerId = f.id;
                tempMarker = { ...f };
                isEditorOpen = true;
              }}
            >
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center gap-2">
                  <span class="flex items-center justify-center w-5 h-5 rounded-full text-[10px] font-bold text-white {cat.badge}">
                    {i + 1}
                  </span>
                  <div>
                    <h4 class="text-xs font-semibold text-[#1f1f1f]">{cat.label}</h4>
                    <span class="text-[10px] text-[#444746] font-mono">
                      {f.view === 'front' ? 'Anterior' : 'Posterior'} • {f.createdAt || 'Baru saja'}
                    </span>
                  </div>
                </div>

                {#if f.category === 'pain'}
                  <span class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-[#fff7ee] text-[#b06000] border border-[#ffba5a]/40">
                    NRS {f.severity}/10
                  </span>
                {/if}
              </div>

              {#if f.note}
                <p class="text-xs text-[#444746] mt-2 pl-7 border-l-2 border-[#dadce0]">
                  {f.note}
                </p>
              {/if}
            </div>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</div>
