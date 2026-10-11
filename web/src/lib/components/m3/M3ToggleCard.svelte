<script lang="ts">
  interface Props {
    title: string;
    description?: string;
    checked?: boolean;
    disabled?: boolean;
    colorVariant?: 'primary' | 'emerald' | 'amber' | 'rose';
    onchange?: (checked: boolean) => void;
  }

  let {
    title,
    description = '',
    checked = $bindable(false),
    disabled = false,
    colorVariant = 'primary',
    onchange
  }: Props = $props();

  const colorClasses: Record<string, string> = {
    primary: 'peer-checked:bg-[#0b57d0]',
    emerald: 'peer-checked:bg-emerald-600',
    amber: 'peer-checked:bg-amber-600',
    rose: 'peer-checked:bg-rose-600'
  };

  function handleChange(e: Event) {
    const target = e.target as HTMLInputElement;
    checked = target.checked;
    onchange?.(target.checked);
  }
</script>

<div
  class="flex items-center justify-between p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] transition-colors"
>
  <div class="pr-3 min-w-0">
    <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] block text-xs">
      {title}
    </span>
    {#if description}
      <span class="text-[11px] text-[#444746] dark:text-[#c4c7c5] block mt-0.5 leading-relaxed">
        {description}
      </span>
    {/if}
  </div>
  <label class="relative inline-flex items-center cursor-pointer shrink-0">
    <input
      type="checkbox"
      {checked}
      {disabled}
      onchange={handleChange}
      class="sr-only peer"
    />
    <div
      class="w-11 h-6 bg-gray-200 dark:bg-gray-700 peer-focus:outline-hidden rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 dark:after:border-gray-600 after:border after:rounded-full after:h-5 after:w-5 after:transition-all {colorClasses[colorVariant] || colorClasses.primary} disabled:opacity-50"
    ></div>
  </label>
</div>
