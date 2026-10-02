<script lang="ts">
  import { Save } from '@lucide/svelte';
  import BodyDiagram from '$lib/components/BodyDiagram.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { BodyFinding } from '$lib/types';

  interface Props {
    findings: BodyFinding[];
    onSave?: () => void;
  }

  let { findings = $bindable(), onSave = () => {} }: Props = $props();
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Status Lokalis & Diagram Tubuh</h2>
      <p class="text-xs text-[#444746]">Tandai lokasi keluhan fisik, rasa nyeri, luka robek, fraktur, atau hematoma pasien secara interaktif.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button
        variant="filled"
        onclick={onSave}
      >
        <Save class="w-4 h-4" />
        <span>Simpan Rekam Medis</span>
      </M3Button>
    </div>
  </div>

  <!-- Komponen Interaktif Body Diagram -->
  <BodyDiagram bind:findings />
</div>
