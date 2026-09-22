<script lang="ts">
  import type { Snippet } from 'svelte';

  export type CardVariant = 'outlined' | 'elevated' | 'filled';

  interface Props {
    variant?: CardVariant;
    children?: Snippet;
    class?: string;
  }

  let {
    variant = 'outlined',
    children,
    class: className = ''
  }: Props = $props();

  const variantClasses: Record<CardVariant, string> = {
    outlined: 'bg-surface-container-lowest border border-outline-variant/80',
    elevated: 'bg-surface-container-low shadow-sm border border-outline-variant/30',
    filled: 'bg-surface-container-highest border-0'
  };
</script>

<div class="rounded-2xl transition-all duration-200 {variantClasses[variant] || variantClasses.outlined} {className}">
  {@render children?.()}
</div>
