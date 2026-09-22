<script lang="ts">
  import type { Snippet } from 'svelte';

  export type ChipColor = 'primary' | 'pain' | 'injury' | 'fracture' | 'edema';

  interface Props {
    selected?: boolean;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    color?: ChipColor;
    class?: string;
  }

  let {
    selected = false,
    onclick,
    children,
    color = 'primary',
    class: className = ''
  }: Props = $props();

  let activeClass = $derived.by(() => {
    const colorStyles: Record<ChipColor, string> = {
      primary: selected
        ? 'bg-[#c2e7ff] text-[#001d35] border-[#7fcfff] font-semibold shadow-xs'
        : 'bg-white text-[#444746] border-[#c4c7c5] hover:bg-[#f8fafd] hover:text-[#1f1f1f]',
      pain: selected
        ? 'bg-[#ffe7c4] text-[#422200] border-[#ffba5a] font-semibold shadow-xs'
        : 'bg-white text-[#b06000] border-[#ffba5a]/50 hover:bg-[#fff7ee]',
      injury: selected
        ? 'bg-[#f9dedc] text-[#410e0b] border-[#f2b8b5] font-semibold shadow-xs'
        : 'bg-white text-[#b3261e] border-[#f2b8b5]/50 hover:bg-[#fdf2f2]',
      fracture: selected
        ? 'bg-[#eedeff] text-[#2c0052] border-[#d8b9ff] font-semibold shadow-xs'
        : 'bg-white text-[#7c3aed] border-[#d8b9ff]/50 hover:bg-[#f9f5ff]',
      edema: selected
        ? 'bg-[#d3e3fd] text-[#041e49] border-[#a8c7fa] font-semibold shadow-xs'
        : 'bg-white text-[#0b57d0] border-[#a8c7fa]/50 hover:bg-[#f0f4f9]'
    };
    return colorStyles[color] || colorStyles.primary;
  });
</script>

<button
  type="button"
  {onclick}
  class="inline-flex items-center gap-2 h-8 px-3.5 text-xs rounded-lg border transition-all duration-150 cursor-pointer select-none {activeClass} {className}"
>
  {@render children?.()}
</button>
