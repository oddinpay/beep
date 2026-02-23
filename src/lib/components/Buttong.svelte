<script lang="ts">
  import { ArrowRightIcon } from "lucide-svelte";
  import Button from "./Button.svelte";
  export let label = "Get started";
  export let url = "";
</script>

<Button
  onclick={() => window.open(url, "_blank")}
  class="group transition-transform  hover:bg-green-300 focus:bg-green-300  bg-green-700 text-white px-[14px] py-[20px] rounded-full flex items-center gap-2 active:scale-99 font-medium cursor-pointer"
>
  <span>{label}</span>
  <ArrowRightIcon
    size={16}
    class="transition-transform duration-200 group-hover:translate-x-1"
  />
</Button>
