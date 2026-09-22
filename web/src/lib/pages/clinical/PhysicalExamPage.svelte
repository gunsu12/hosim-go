<script lang="ts">
  import { Code2, Save } from '@lucide/svelte';
  import BodyDiagram from '$lib/components/BodyDiagram.svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { BodyFinding, ClinicalNotes } from '$lib/types';

  interface Props {
    findings: BodyFinding[];
    clinicalNotes: ClinicalNotes;
    onSave: () => void;
  }

  let { findings = $bindable(), clinicalNotes, onSave }: Props = $props();
  let showJsonDebug = $state(false);
</script>

<div class="flex flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <div>
      <h2 class="text-lg font-semibold text-[#1f1f1f] tracking-tight">Status Lokalis & Diagram Tubuh</h2>
      <p class="text-xs text-[#444746]">Tandai lokasi keluhan fisik, rasa nyeri, luka robek, fraktur, atau hematoma pasien secara interaktif.</p>
    </div>

    <div class="flex items-center gap-2">
      <M3Button
        variant="outlined"
        onclick={() => showJsonDebug = !showJsonDebug}
      >
        <Code2 class="w-4 h-4" />
        <span>{showJsonDebug ? 'Tutup Payload JSON' : 'Lihat Data JSON API'}</span>
      </M3Button>

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

  <!-- JSON Payload Inspector -->
  {#if showJsonDebug}
    <div class="p-4 bg-[#f8fafd] rounded-2xl border border-[#e1e5ea] font-mono text-xs animate-in fade-in duration-150">
      <div class="flex items-center justify-between pb-2 mb-2 border-b border-[#e1e5ea] text-[#444746]">
        <span class="font-bold">Format Payload JSON yang Dikirim ke Go Gin API (`POST /api/v1/clinical/physical-exam`)</span>
        <span>{findings.length} Titik Temuan</span>
      </div>
      <pre class="bg-white p-3 rounded-xl border border-[#e1e5ea] overflow-x-auto text-[11px] text-[#1f1f1f] leading-relaxed max-h-56">{JSON.stringify({
        patient_mrn: "RM-2026-08492",
        practitioner_id: 1,
        encounter_type: "OUTPATIENT",
        body_diagram_findings: findings,
        clinical_notes: clinicalNotes
      }, null, 2)}</pre>
    </div>
  {/if}
</div>
