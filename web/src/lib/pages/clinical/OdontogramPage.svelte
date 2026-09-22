<script lang="ts">
  // Odontogram Gigi 32 Gigi FDI System
  const maxillaTeeth = [18, 17, 16, 15, 14, 13, 12, 11, 21, 22, 23, 24, 25, 26, 27, 28];
  const mandibulaTeeth = [48, 47, 46, 45, 44, 43, 42, 41, 31, 32, 33, 34, 35, 36, 37, 38];

  let selectedTooth = $state<number | null>(null);
</script>

<div class="p-6 bg-[#f8fafd] rounded-2xl border border-[#e1e5ea]">
  <div class="flex items-center justify-between mb-4">
    <div>
      <h3 class="text-base font-semibold text-[#1f1f1f]">Pemeriksaan Odontogram Gigi</h3>
      <p class="text-xs text-[#444746]">Peta kondisi 32 gigi pasien dewasa (Karies, Tumpatan, Fraktur, Gigi Hilang).</p>
    </div>
    <span class="text-xs font-semibold px-3 py-1 rounded-full bg-[#d3e3fd] text-[#041e49]">Odontogram FDI System</span>
  </div>

  <!-- Sketsa Interaktif Odontogram 32 Gigi Google Style -->
  <div class="p-6 bg-white rounded-2xl border border-[#e1e5ea] flex flex-col items-center gap-6">
    <!-- Rahang Atas (Maxilla 18-28) -->
    <div>
      <div class="text-xs font-bold text-center text-[#747775] mb-2">Rahang Atas (Maxilla)</div>
      <div class="flex flex-wrap justify-center gap-1.5">
        {#each maxillaTeeth as tooth}
          <button
            type="button"
            onclick={() => selectedTooth = tooth}
            class="w-10 h-12 rounded-xl border transition-all cursor-pointer group flex flex-col items-center justify-center {selectedTooth === tooth ? 'border-[#0b57d0] bg-[#e8f0fe] ring-2 ring-[#0b57d0]/30' : 'border-[#c4c7c5] hover:border-[#0b57d0] hover:bg-[#e8f0fe]'}"
          >
            <span class="text-[11px] font-bold text-[#1f1f1f] group-hover:text-[#0b57d0]">{tooth}</span>
            <span class="w-2 h-2 rounded-full {tooth === 16 ? 'bg-[#d93025]' : tooth === 24 ? 'bg-[#0b57d0]' : 'bg-[#e9eef6]'} mt-1"></span>
          </button>
        {/each}
      </div>
    </div>

    <!-- Rahang Bawah (Mandibula 48-38) -->
    <div>
      <div class="flex flex-wrap justify-center gap-1.5">
        {#each mandibulaTeeth as tooth}
          <button
            type="button"
            onclick={() => selectedTooth = tooth}
            class="w-10 h-12 rounded-xl border transition-all cursor-pointer group flex flex-col items-center justify-center {selectedTooth === tooth ? 'border-[#0b57d0] bg-[#e8f0fe] ring-2 ring-[#0b57d0]/30' : 'border-[#c4c7c5] hover:border-[#0b57d0] hover:bg-[#e8f0fe]'}"
          >
            <span class="text-[11px] font-bold text-[#1f1f1f] group-hover:text-[#0b57d0]">{tooth}</span>
            <span class="w-2 h-2 rounded-full {tooth === 36 ? 'bg-[#e37400]' : 'bg-[#e9eef6]'} mt-1"></span>
          </button>
        {/each}
      </div>
      <div class="text-xs font-bold text-center text-[#747775] mt-2">Rahang Bawah (Mandibula)</div>
    </div>

    {#if selectedTooth}
      <div class="w-full max-w-md p-3 bg-[#f8fafd] rounded-xl border border-[#e1e5ea] text-xs flex items-center justify-between animate-in fade-in">
        <div>
          Gigi Terpilih: <strong>Elemen {selectedTooth}</strong>
        </div>
        <div class="flex items-center gap-1.5">
          <button type="button" class="px-2 py-1 rounded bg-[#e6f4ea] text-[#137333] font-medium text-[11px] cursor-pointer">Sehat / Sound</button>
          <button type="button" class="px-2 py-1 rounded bg-[#fce8e6] text-[#c5221f] font-medium text-[11px] cursor-pointer">Karies</button>
          <button type="button" class="px-2 py-1 rounded bg-[#e8f0fe] text-[#0b57d0] font-medium text-[11px] cursor-pointer">Tambalan</button>
        </div>
      </div>
    {/if}
  </div>
</div>
