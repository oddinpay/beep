<script lang="ts">
  import { Diamonds } from "svelte-loading-spinners";
  import { onMount } from "svelte";

  let showLoading = true;
  let timer: ReturnType<typeof setTimeout> | undefined;

  onMount(() => {
    timer = setTimeout(() => {
      showLoading = false;
    }, 400);

    return () => {
      if (timer) clearTimeout(timer);
    };
  });
</script>

{#if showLoading}
  <div
    class="fixed inset-0 bg-[#DCD9E7] bg-opacity-50 backdrop-blur-sm flex justify-center pt-15 z-[9999]"
  >
    <Diamonds size={60} color="green" unit="px" duration="1s" />
  </div>
{/if}
