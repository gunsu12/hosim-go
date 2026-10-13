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

<div class="p-6 bg-white dark:bg-[#1e1f20] rounded-2xl border border-[#e1e5ea] dark:border-[#353738] shadow-xs">
  <h3 class="text-base font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] mb-4">Catatan Anamnesis & Resume Medis (SOAP)</h3>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
    <div>
      <label for="chief-complaint-textarea" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1.5">
        S (Subjective) - Keluhan Utama
      </label>
      <textarea
        id="chief-complaint-textarea"
        bind:value={clinicalNotes.chiefComplaint}
        rows="3"
        class="w-full p-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
      ></textarea>
    </div>

    <div>
      <label for="diagnosis-input" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1.5">
        A (Assessment) - Diagnosis Utama ICD-10
      </label>
      <input
        id="diagnosis-input"
        type="text"
        bind:value={clinicalNotes.diagnosis}
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>

    <div class="md:col-span-2">
      <label for="anamnesis-textarea" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1.5">
        O (Objective) - Riwayat Penyakit Sekarang
      </label>
      <textarea
        id="anamnesis-textarea"
        bind:value={clinicalNotes.anamnesis}
        rows="3"
        class="w-full p-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all resize-none"
      ></textarea>
    </div>

    <div class="md:col-span-2">
      <label for="disposition-input" class="block font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] text-xs mb-1.5">
        P (Plan) - Rencana Tindakan & Disposisi
      </label>
      <input
        id="disposition-input"
        type="text"
        bind:value={clinicalNotes.disposition}
        class="w-full h-10 px-3.5 rounded-xl bg-[#f0f4f9] dark:bg-[#28292a] border border-transparent focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:bg-white dark:focus:bg-[#1e1f20] text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] focus:outline-hidden transition-all"
      />
    </div>
  </div>

  <div class="mt-6 flex justify-end">
    <M3Button variant="filled" onclick={onSave}>
      <Save class="w-4 h-4 mr-1.5" />
      <span>Simpan Catatan SOAP</span>
    </M3Button>
  </div>
</div>
