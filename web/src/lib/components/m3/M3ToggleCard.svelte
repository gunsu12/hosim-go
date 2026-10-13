<script lang="ts">
  interface Props {
    title: string;
    description?: string;
    subtitle?: string;
    icon?: any;
    checked?: boolean;
    disabled?: boolean;
    colorVariant?: 'primary' | 'emerald' | 'amber' | 'rose';
    onchange?: (checked: boolean) => void;
  }

  let {
    title,
    description = '',
    subtitle = '',
    icon = null,
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

  let descText = $derived(description || subtitle);

  function handleChange(e: Event) {
    const target = e.target as HTMLInputElement;
    checked = target.checked;
    onchange?.(target.checked);
  }
</script>

<div
  class="flex items-center justify-between p-3.5 rounded-2xl bg-[#f8fafd] dark:bg-[#28292a] border border-[#e1e5ea] dark:border-[#353738] transition-colors"
>
  <div class="pr-3 min-w-0 flex items-start gap-3">
    {#if icon}
      {@const IconComponent = icon}
      <div class="w-8 h-8 rounded-xl bg-white dark:bg-[#1e1f20] border border-[#e1e5ea] dark:border-[#353738] flex items-center justify-center shrink-0 mt-0.5 text-[#0b57d0] dark:text-[#a8c7fa]">
        <IconComponent class="w-4 h-4" />
      </div>
    {/if}
    <div class="min-w-0">
      <span class="font-bold text-[#1f1f1f] dark:text-[#e3e3e3] block text-xs">
        {title}
      </span>
      {#if descText}
        <span class="text-[11px] text-[#444746] dark:text-[#c4c7c5] block mt-0.5 leading-relaxed">
          {descText}
        </span>
      {/if}
    </div>
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
