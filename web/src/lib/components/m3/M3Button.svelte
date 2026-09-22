<script lang="ts">
  import type { Snippet } from 'svelte';

  export type ButtonVariant = 'filled' | 'tonal' | 'outlined' | 'text' | 'elevated' | 'drive-new';
  export type ButtonType = 'button' | 'submit' | 'reset';

  interface Props {
    variant?: ButtonVariant;
    disabled?: boolean;
    type?: ButtonType;
    onclick?: (e: MouseEvent) => void;
    children?: Snippet;
    class?: string;
  }

  let {
    variant = 'filled',
    disabled = false,
    type = 'button',
    onclick,
    children,
    class: className = ''
  }: Props = $props();

  const variantClasses: Record<ButtonVariant, string> = {
    filled: 'bg-[#0b57d0] text-white hover:bg-[#0842a0] active:bg-[#073887] shadow-xs hover:shadow',
    tonal: 'bg-[#d3e3fd] text-[#041e49] hover:bg-[#c2e7ff] active:bg-[#a8c7fa]',
    outlined: 'border border-[#747775]/50 bg-transparent text-[#0b57d0] hover:bg-[#0b57d0]/8 active:bg-[#0b57d0]/12',
    text: 'bg-transparent text-[#0b57d0] hover:bg-[#0b57d0]/8 active:bg-[#0b57d0]/12',
    elevated: 'bg-white text-[#1f1f1f] shadow-xs border border-[#e1e5ea] hover:shadow hover:bg-[#f8fafd] active:bg-[#f0f4f9]',
    'drive-new': 'bg-white text-[#1f1f1f] shadow-md hover:shadow-lg hover:bg-[#f8fafd] active:bg-[#e9eef6] border border-[#e1e5ea] h-12 px-6 rounded-2xl font-semibold text-sm'
  };
</script>

<button
  {type}
  {disabled}
  {onclick}
  class="inline-flex items-center justify-center gap-2 h-10 px-5 text-sm font-medium tracking-normal rounded-full transition-all duration-150 disabled:opacity-38 disabled:pointer-events-none cursor-pointer {variantClasses[variant] || variantClasses.filled} {className}"
>
  {@render children?.()}
</button>
