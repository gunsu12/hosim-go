<script lang="ts">
  import { AlertTriangle, HeartPulse } from '@lucide/svelte';
  import type { Patient } from '../types';

  interface Props {
    patient?: Patient;
  }

  let {
    patient = {
      name: 'Budi Santoso',
      mrn: 'RM-2026-08492',
      nik: '3273151405810003',
      birthDate: '14 Mei 1981 (45 th)',
      gender: 'Laki-laki',
      bloodType: 'O+',
      payer: 'BPJS Kesehatan (JKN PBI)',
      allergies: ['Amoxicillin', 'Ibuprofen'],
      vitals: {
        bp: '120/80',
        hr: '82',
        rr: '18',
        temp: '36.8',
        spo2: '98'
      }
    }
  }: Props = $props();
</script>

<div class="p-4 bg-white rounded-2xl border border-[#e1e5ea] shadow-xs">
  <div class="flex flex-col lg:flex-row items-start lg:items-center justify-between gap-4">
    <!-- Identitas Pasien -->
    <div class="flex items-center gap-3.5">
      <div class="w-12 h-12 rounded-2xl bg-[#0b57d0] text-white flex items-center justify-center font-bold text-base shadow-xs">
        BS
      </div>
      <div>
        <div class="flex flex-wrap items-center gap-2">
          <h2 class="text-base font-semibold text-[#1f1f1f]">{patient.name}</h2>
          <span class="text-xs font-mono font-medium px-2.5 py-0.5 rounded-full bg-[#f0f4f9] text-[#444746] border border-[#e1e5ea]">
            {patient.mrn}
          </span>
          <span class="text-xs font-medium px-2.5 py-0.5 rounded-full bg-[#e8f0fe] text-[#0b57d0]">
            {patient.gender}
          </span>
          <span class="text-xs font-medium px-2.5 py-0.5 rounded-full bg-[#e6f4ea] text-[#137333]">
            {patient.payer}
          </span>
        </div>
        <div class="flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-[#444746] mt-1">
          <span>Lahir: <strong class="text-[#1f1f1f] font-normal">{patient.birthDate}</strong></span>
          <span>•</span>
          <span>Gol. Darah: <strong class="text-[#1f1f1f] font-semibold">{patient.bloodType}</strong></span>
          <span>•</span>
          <span>NIK: <span class="font-mono">{patient.nik}</span></span>
        </div>
      </div>
    </div>

    <!-- Tanda Vital & Alergi -->
    <div class="flex flex-wrap items-center gap-2.5 self-stretch lg:self-auto justify-start lg:justify-end">
      <!-- Badge Alergi -->
      <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-[#fce8e6] border border-[#fad2cf] text-[#c5221f] text-xs font-medium">
        <AlertTriangle class="w-3.5 h-3.5 text-[#c5221f] shrink-0" />
        <span class="font-semibold">Alergi:</span>
        <span>{patient.allergies.join(', ')}</span>
      </div>

      <!-- Vital Signs Pills -->
      <div class="flex items-center gap-2 bg-[#f8fafd] px-3.5 py-1.5 rounded-xl border border-[#e1e5ea] text-xs text-[#1f1f1f]">
        <HeartPulse class="w-4 h-4 text-[#0b57d0]" />
        <span>TD <strong>{patient.vitals.bp}</strong></span>
        <span class="text-[#c4c7c5]">|</span>
        <span>HR <strong>{patient.vitals.hr}</strong> bpm</span>
        <span class="text-[#c4c7c5]">|</span>
        <span>T <strong>{patient.vitals.temp}</strong>°C</span>
        <span class="text-[#c4c7c5]">|</span>
        <span>SpO2 <strong>{patient.vitals.spo2}%</strong></span>
      </div>
    </div>
  </div>
</div>
