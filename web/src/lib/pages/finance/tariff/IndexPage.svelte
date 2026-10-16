<script lang="ts">
  import PricePlanPage from "./PricePlanPage.svelte";
  import TariffMatrixPage from "./TariffMatrixPage.svelte";
  import TariffLookupPage from "./TariffLookupPage.svelte";
  import TariffClassPage from "./TariffClassPage.svelte";
  import TariffComponentPage from "./TariffComponentPage.svelte";

  interface Props {
    initialTab?: "PLAN" | "MATRIX" | "SIMULATOR" | "CLASS" | "COMPONENT";
  }

  let { initialTab = "PLAN" }: Props = $props();
</script>

{#if initialTab === "MATRIX"}
  <TariffMatrixPage />
{:else if initialTab === "SIMULATOR"}
  <TariffLookupPage />
{:else if initialTab === "CLASS"}
  <TariffClassPage />
{:else if initialTab === "COMPONENT"}
  <TariffComponentPage />
{:else}
  <PricePlanPage />
{/if}
