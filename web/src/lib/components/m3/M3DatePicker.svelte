<script lang="ts">
  import {
    Calendar as CalendarIcon,
    ChevronLeft,
    ChevronRight,
    ChevronDown,
    X,
    Check
  } from '@lucide/svelte';

  interface Props {
    value?: string; // Format: YYYY-MM-DD
    label?: string;
    id?: string;
    name?: string;
    placeholder?: string;
    required?: boolean;
    disabled?: boolean;
    min?: string; // Format: YYYY-MM-DD
    max?: string; // Format: YYYY-MM-DD
    class?: string;
    inputClass?: string;
  }

  let {
    value = $bindable(''),
    label = '',
    id = '',
    name = '',
    placeholder = 'Pilih tanggal',
    required = false,
    disabled = false,
    min = '',
    max = '',
    class: className = '',
    inputClass = ''
  }: Props = $props();

  const MONTH_NAMES = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];

  const MONTH_NAMES_SHORT = [
    'Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun',
    'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'
  ];

  const DAY_NAMES = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'];
  const DAY_NAMES_FULL = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];

  let isOpen = $state(false);
  let viewMode = $state<'calendar' | 'year'>('calendar');

  // Internal calendar navigation state
  let viewYear = $state(new Date().getFullYear());
  let viewMonth = $state(new Date().getMonth()); // 0-11
  let tempSelectedDate = $state<string>(''); // YYYY-MM-DD

  // Format Helper: YYYY-MM-DD to Indonesian display
  function formatDisplayDate(dateStr: string): string {
    if (!dateStr) return '';
    const parts = dateStr.split('-');
    if (parts.length !== 3) return dateStr;
    const y = parseInt(parts[0], 10);
    const m = parseInt(parts[1], 10) - 1;
    const d = parseInt(parts[2], 10);
    if (isNaN(y) || isNaN(m) || isNaN(d)) return dateStr;
    return `${d} ${MONTH_NAMES_SHORT[m] || ''} ${y}`;
  }

  // Format Headline: "Rab, 11 Okt 2026"
  function formatHeadline(dateStr: string): string {
    if (!dateStr) return 'Pilih tanggal';
    const parts = dateStr.split('-');
    if (parts.length !== 3) return 'Pilih tanggal';
    const y = parseInt(parts[0], 10);
    const m = parseInt(parts[1], 10) - 1;
    const d = parseInt(parts[2], 10);
    const dateObj = new Date(y, m, d);
    if (isNaN(dateObj.getTime())) return 'Pilih tanggal';
    const dayName = DAY_NAMES[dateObj.getDay()];
    return `${dayName}, ${d} ${MONTH_NAMES_SHORT[m]} ${y}`;
  }

  function openPicker() {
    if (disabled) return;
    tempSelectedDate = value || '';
    if (value) {
      const parts = value.split('-');
      if (parts.length === 3) {
        const y = parseInt(parts[0], 10);
        const m = parseInt(parts[1], 10) - 1;
        if (!isNaN(y) && !isNaN(m)) {
          viewYear = y;
          viewMonth = m;
        }
      }
    } else {
      const now = new Date();
      viewYear = now.getFullYear();
      viewMonth = now.getMonth();
    }
    viewMode = 'calendar';
    isOpen = true;
  }

  function closePicker() {
    isOpen = false;
  }

  function handleSelectDate(day: number) {
    const mStr = String(viewMonth + 1).padStart(2, '0');
    const dStr = String(day).padStart(2, '0');
    tempSelectedDate = `${viewYear}-${mStr}-${dStr}`;
  }

  function handleConfirm() {
    value = tempSelectedDate;
    isOpen = false;
  }

  function handleSelectToday() {
    const now = new Date();
    viewYear = now.getFullYear();
    viewMonth = now.getMonth();
    const mStr = String(now.getMonth() + 1).padStart(2, '0');
    const dStr = String(now.getDate()).padStart(2, '0');
    tempSelectedDate = `${now.getFullYear()}-${mStr}-${dStr}`;
  }

  function handleClear(e: MouseEvent) {
    e.stopPropagation();
    value = '';
    tempSelectedDate = '';
  }

  function prevMonth() {
    if (viewMonth === 0) {
      viewMonth = 11;
      viewYear -= 1;
    } else {
      viewMonth -= 1;
    }
  }

  function nextMonth() {
    if (viewMonth === 11) {
      viewMonth = 0;
      viewYear += 1;
    } else {
      viewMonth += 1;
    }
  }

  function selectYear(y: number) {
    viewYear = y;
    viewMode = 'calendar';
  }

  // Calendar days grid computation
  interface CalendarDay {
    day: number;
    isCurrentMonth: boolean;
    dateStr: string;
    isToday: boolean;
    isSelected: boolean;
    isDisabled: boolean;
  }

  let calendarDays = $derived.by<CalendarDay[]>(() => {
    const today = new Date();
    const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`;

    // First day of current month (0: Sunday, 1: Monday, ...)
    const firstDayIndex = new Date(viewYear, viewMonth, 1).getDay();
    // Total days in current month
    const daysInMonth = new Date(viewYear, viewMonth + 1, 0).getDate();
    // Total days in previous month
    const daysInPrevMonth = new Date(viewYear, viewMonth, 0).getDate();

    const days: CalendarDay[] = [];

    // Leading days from previous month
    for (let i = firstDayIndex - 1; i >= 0; i--) {
      const d = daysInPrevMonth - i;
      const prevM = viewMonth === 0 ? 12 : viewMonth;
      const prevY = viewMonth === 0 ? viewYear - 1 : viewYear;
      const dateStr = `${prevY}-${String(prevM).padStart(2, '0')}-${String(d).padStart(2, '0')}`;
      days.push({
        day: d,
        isCurrentMonth: false,
        dateStr,
        isToday: dateStr === todayStr,
        isSelected: false,
        isDisabled: true
      });
    }

    // Days in current month
    for (let d = 1; d <= daysInMonth; d++) {
      const dateStr = `${viewYear}-${String(viewMonth + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`;
      const isSelected = dateStr === tempSelectedDate;
      const isToday = dateStr === todayStr;

      let isDisabled = false;
      if (min && dateStr < min) isDisabled = true;
      if (max && dateStr > max) isDisabled = true;

      days.push({
        day: d,
        isCurrentMonth: true,
        dateStr,
        isToday,
        isSelected,
        isDisabled
      });
    }

    // Trailing days from next month to complete 42 or 35 grid cells
    const remaining = (7 - (days.length % 7)) % 7;
    for (let d = 1; d <= remaining; d++) {
      const nextM = viewMonth === 11 ? 1 : viewMonth + 2;
      const nextY = viewMonth === 11 ? viewYear + 1 : viewYear;
      const dateStr = `${nextY}-${String(nextM).padStart(2, '0')}-${String(d).padStart(2, '0')}`;
      days.push({
        day: d,
        isCurrentMonth: false,
        dateStr,
        isToday: dateStr === todayStr,
        isSelected: false,
        isDisabled: true
      });
    }

    return days;
  });

  // Year list for Year Selector view
  const yearList = $derived.by<number[]>(() => {
    const currentYear = new Date().getFullYear();
    const list: number[] = [];
    for (let y = currentYear - 30; y <= currentYear + 20; y++) {
      list.push(y);
    }
    return list;
  });

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape' && isOpen) {
      closePicker();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="relative flex flex-col {className}">
  {#if label}
    <label for={id} class="block font-semibold text-[#444746] dark:text-[#c4c7c5] mb-1 text-xs">
      {label}
    </label>
  {/if}

  <!-- Trigger Field: M3 Outlined Container Style -->
  <div class="relative w-full">
    <button
      {id}
      type="button"
      {disabled}
      onclick={openPicker}
      class="flex items-center justify-between w-full h-9 px-3 pr-8 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] text-xs transition-all duration-150 cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed hover:border-[#0b57d0] dark:hover:border-[#a8c7fa] focus:border-[#0b57d0] dark:focus:border-[#a8c7fa] focus:outline-hidden {inputClass}"
      aria-haspopup="dialog"
      aria-expanded={isOpen}
    >
      <div class="flex items-center gap-2 truncate text-left">
        <CalendarIcon class="w-4 h-4 shrink-0 text-[#0b57d0] dark:text-[#a8c7fa]" />
        {#if value}
          <span class="font-medium text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
            {formatDisplayDate(value)}
          </span>
        {:else}
          <span class="text-[#747775] dark:text-[#8e918f] truncate font-normal">
            {placeholder}
          </span>
        {/if}
      </div>
    </button>

    {#if value && !disabled && !required}
      <button
        type="button"
        onclick={handleClear}
        class="absolute right-2.5 top-1/2 -translate-y-1/2 p-0.5 rounded-full hover:bg-[#e9eef6] dark:hover:bg-[#28292a] text-[#747775] dark:text-[#8e918f] hover:text-[#1f1f1f] dark:hover:text-white transition-colors cursor-pointer"
        title="Hapus tanggal"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    {/if}
  </div>

  <!-- Hidden native input to support native form validations & bindings if needed -->
  {#if name}
    <input type="hidden" {name} {value} {required} />
  {/if}
</div>

<!-- Material Design 3 Modal Date Picker Dialog -->
{#if isOpen}
  <!-- Backdrop Overlay -->
  <div
    class="fixed inset-0 z-50 bg-black/60 backdrop-blur-xs flex items-center justify-center p-4 animate-in fade-in duration-200"
    role="presentation"
    onclick={closePicker}
  >
    <!-- M3 Dialog Surface (Mobile & Desktop) -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      role="dialog"
      aria-modal="true"
      tabindex="-1"
      aria-label="Pemilih Tanggal Material Design"
      class="w-full max-w-[340px] sm:max-w-[360px] bg-white dark:bg-[#1e1f20] rounded-[28px] shadow-2xl border border-[#e1e5ea] dark:border-[#353738] overflow-hidden flex flex-col animate-in zoom-in-95 duration-200"
      onclick={(e) => e.stopPropagation()}
    >
      <!-- M3 Date Picker Header -->
      <div class="px-6 pt-5 pb-3 border-b border-[#e1e5ea] dark:border-[#2e3032] bg-[#f8fafd] dark:bg-[#18191a]">
        <div class="text-[11px] font-bold uppercase tracking-wider text-[#747775] dark:text-[#8e918f]">
          PILIH TANGGAL
        </div>
        <div class="mt-1 flex items-center justify-between">
          <span class="text-xl sm:text-2xl font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] truncate">
            {formatHeadline(tempSelectedDate)}
          </span>
          <button
            type="button"
            onclick={() => viewMode = viewMode === 'calendar' ? 'year' : 'calendar'}
            class="p-1.5 rounded-full hover:bg-[#e9eef6] dark:hover:bg-[#28292a] text-[#0b57d0] dark:text-[#a8c7fa] transition-colors cursor-pointer"
            title="Ganti mode pemilihan tahun"
          >
            <ChevronDown class="w-4.5 h-4.5 transition-transform {viewMode === 'year' ? 'rotate-180' : ''}" />
          </button>
        </div>
      </div>

      <!-- Main Body: Calendar or Year Picker -->
      <div class="p-4 flex flex-col">
        {#if viewMode === 'calendar'}
          <!-- Month & Year Navigation Header -->
          <div class="flex items-center justify-between mb-3 px-1">
            <button
              type="button"
              onclick={() => viewMode = 'year'}
              class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] transition-colors cursor-pointer"
            >
              <span>{MONTH_NAMES[viewMonth]} {viewYear}</span>
              <ChevronDown class="w-3.5 h-3.5 text-[#747775] dark:text-[#8e918f]" />
            </button>

            <div class="flex items-center gap-1">
              <button
                type="button"
                onclick={prevMonth}
                class="w-8 h-8 rounded-full flex items-center justify-center text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] active:bg-[#e1e5ea] transition-colors cursor-pointer"
                title="Bulan sebelumnya"
              >
                <ChevronLeft class="w-4.5 h-4.5" />
              </button>
              <button
                type="button"
                onclick={nextMonth}
                class="w-8 h-8 rounded-full flex items-center justify-center text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] active:bg-[#e1e5ea] transition-colors cursor-pointer"
                title="Bulan berikutnya"
              >
                <ChevronRight class="w-4.5 h-4.5" />
              </button>
            </div>
          </div>

          <!-- Day of Week Headers -->
          <div class="grid grid-cols-7 gap-1 mb-1 text-center">
            {#each DAY_NAMES as dn, i}
              <span class="text-[11px] font-semibold {i === 0 ? 'text-red-500/80 dark:text-red-400/80' : 'text-[#747775] dark:text-[#8e918f]'}">
                {dn}
              </span>
            {/each}
          </div>

          <!-- Calendar Days Grid -->
          <div class="grid grid-cols-7 gap-1 place-items-center">
            {#each calendarDays as cd}
              {#if cd.isCurrentMonth}
                <button
                  type="button"
                  disabled={cd.isDisabled}
                  onclick={() => handleSelectDate(cd.day)}
                  class="w-9 h-9 sm:w-10 sm:h-10 rounded-full flex items-center justify-center text-xs sm:text-sm font-medium transition-all duration-150 cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed
                    {cd.isSelected
                      ? 'bg-[#0b57d0] text-white dark:bg-[#a8c7fa] dark:text-[#062e6f] font-bold shadow-sm scale-105'
                      : cd.isToday
                        ? 'border border-[#0b57d0] text-[#0b57d0] dark:border-[#a8c7fa] dark:text-[#a8c7fa] font-bold hover:bg-[#e9eef6] dark:hover:bg-[#28292a]'
                        : 'text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#e9eef6] dark:hover:bg-[#28292a]'
                    }"
                >
                  {cd.day}
                </button>
              {:else}
                <div class="w-9 h-9 sm:w-10 sm:h-10 flex items-center justify-center text-xs text-[#747775]/30 dark:text-[#8e918f]/30 select-none">
                  {cd.day}
                </div>
              {/if}
            {/each}
          </div>
        {:else}
          <!-- Year Selector Grid (M3 Year Picker) -->
          <div class="h-64 overflow-y-auto pr-1 grid grid-cols-3 gap-2 py-2">
            {#each yearList as y}
              {@const isYearActive = y === viewYear}
              <button
                type="button"
                onclick={() => selectYear(y)}
                class="h-10 rounded-full text-xs sm:text-sm font-medium transition-all cursor-pointer flex items-center justify-center {isYearActive
                  ? 'bg-[#0b57d0] text-white dark:bg-[#a8c7fa] dark:text-[#062e6f] font-bold shadow-xs'
                  : 'text-[#1f1f1f] dark:text-[#e3e3e3] hover:bg-[#e9eef6] dark:hover:bg-[#28292a]'}"
              >
                {y}
              </button>
            {/each}
          </div>
        {/if}
      </div>

      <!-- M3 Actions Footer Bar -->
      <div class="px-4 py-3 bg-[#f8fafd] dark:bg-[#18191a] border-t border-[#e1e5ea] dark:border-[#2e3032] flex items-center justify-between">
        <button
          type="button"
          onclick={handleSelectToday}
          class="px-3 py-1.5 rounded-full text-xs font-semibold text-[#0b57d0] dark:text-[#a8c7fa] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] transition-colors cursor-pointer"
        >
          Hari Ini
        </button>

        <div class="flex items-center gap-2">
          <button
            type="button"
            onclick={closePicker}
            class="px-4 py-1.5 rounded-full text-xs font-medium text-[#444746] dark:text-[#c4c7c5] hover:bg-[#e9eef6] dark:hover:bg-[#28292a] transition-colors cursor-pointer"
          >
            Batal
          </button>
          <button
            type="button"
            onclick={handleConfirm}
            class="px-4 py-1.5 rounded-full text-xs font-semibold bg-[#0b57d0] text-white hover:bg-[#0842a0] dark:bg-[#a8c7fa] dark:text-[#062e6f] dark:hover:bg-[#8ab4f8] shadow-xs transition-colors cursor-pointer"
          >
            Pilih
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
