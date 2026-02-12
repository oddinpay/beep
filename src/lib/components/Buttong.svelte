<script lang="ts">
  import { ArrowRightIcon } from "lucide-svelte";
  import Button from "./Button.svelte";
  import { goto } from "$app/navigation";

  export let label = "Get started";


</script>

<Button
  onclick={() => goto("https://oddinpay.com/signup")}
  class="group transition-transform  hover:bg-green-600 focus:bg-green-600  bg-green-600 text-white px-[14px] py-[20px] rounded-full flex items-center gap-2 active:scale-99 font-medium cursor-pointer hover:opacity-90 focus:opacity-90"
>
  <span>{label}</span>
  <ArrowRightIcon
    size={16}
    class="transition-transform duration-200 group-hover:translate-x-1"
  />
</Button>