<script lang="ts">
  export interface SegmentOption {
    value: string;
    label: string;
    icon?: any;
  }

  interface Props {
    options?: SegmentOption[];
    selected: string;
    class?: string;
  }

  let {
    options = [],
    selected = $bindable(),
    class: className = ''
  }: Props = $props();
</script>

<div class="inline-flex rounded-full p-1 bg-[#e9eef6] border border-[#e1e5ea] {className}">
  {#each options as opt}
    {@const isSelected = selected === opt.value}
    <button
      type="button"
      onclick={() => selected = opt.value}
      class="flex items-center gap-2 px-4 py-1.5 text-xs font-medium rounded-full transition-all duration-150 cursor-pointer {isSelected ? 'bg-white text-[#0b57d0] shadow-xs font-semibold' : 'text-[#444746] hover:text-[#1f1f1f] hover:bg-white/50'}"
    >
      {#if opt.icon}
        <opt.icon class="w-3.5 h-3.5" />
      {/if}
      <span>{opt.label}</span>
    </button>
  {/each}
</div>
