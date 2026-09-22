<script lang="ts">
  import { Save } from '@lucide/svelte';
  import M3Button from '$lib/components/m3/M3Button.svelte';
  import type { ClinicalNotes } from '$lib/types';

  interface Props {
    clinicalNotes: ClinicalNotes;
    onSave: () => void;
  }

  let { clinicalNotes = $bindable(), onSave }: Props = $props();
</script>

<div class="p-6 bg-[#f8fafd] rounded-2xl border border-[#e1e5ea]">
  <h3 class="text-base font-semibold text-[#1f1f1f] mb-4">Catatan Anamnesis & Resume Medis (SOAP)</h3>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
    <div>
      <label for="chief-complaint-textarea" class="block text-xs font-semibold text-[#444746] mb-1.5">S (Subjective) - Keluhan Utama</label>
      <textarea
        id="chief-complaint-textarea"
        bind:value={clinicalNotes.chiefComplaint}
        rows="3"
        class="w-full p-3 text-xs rounded-xl bg-white border border-[#e1e5ea] text-[#1f1f1f] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
      ></textarea>
    </div>

    <div>
      <label for="diagnosis-input" class="block text-xs font-semibold text-[#444746] mb-1.5">A (Assessment) - Diagnosis Utama ICD-10</label>
      <input
        id="diagnosis-input"
        type="text"
        bind:value={clinicalNotes.diagnosis}
        class="w-full h-11 px-3 text-xs rounded-xl bg-white border border-[#e1e5ea] text-[#1f1f1f] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
      />
    </div>

    <div class="md:col-span-2">
      <label for="anamnesis-textarea" class="block text-xs font-semibold text-[#444746] mb-1.5">O (Objective) - Riwayat Penyakit Sekarang</label>
      <textarea
        id="anamnesis-textarea"
        bind:value={clinicalNotes.anamnesis}
        rows="3"
        class="w-full p-3 text-xs rounded-xl bg-white border border-[#e1e5ea] text-[#1f1f1f] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
      ></textarea>
    </div>

    <div class="md:col-span-2">
      <label for="disposition-input" class="block text-xs font-semibold text-[#444746] mb-1.5">P (Plan) - Rencana Tindakan & Disposisi</label>
      <input
        id="disposition-input"
        type="text"
        bind:value={clinicalNotes.disposition}
        class="w-full h-11 px-3 text-xs rounded-xl bg-white border border-[#e1e5ea] text-[#1f1f1f] focus:ring-2 focus:ring-[#0b57d0] focus:outline-none"
      />
    </div>
  </div>

  <div class="mt-6 flex justify-end">
    <M3Button variant="filled" onclick={onSave}>
      <Save class="w-4 h-4" />
      <span>Simpan Catatan SOAP</span>
    </M3Button>
  </div>
</div>
