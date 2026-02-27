<script lang="ts">
  import "../app.css";
  import favicon from "$lib/assets/favicon.svg";
  import "@fontsource-variable/inter";
  import { page } from "$app/state";
  import "default-passive-events";

  let { children } = $props();

  const canonicalUrl = $derived(
    (() => {
      const protocol = page.url.protocol;
      const origin = page.url.host;
      const path = page.url.pathname;

      const base = `${protocol}//${origin}`;

      if (path === "/") return base;
      if (/^\/[^/]+$/.test(path)) return base + path;
      return base + path;
    })(),
  );
</script>

<svelte:head>
  <link rel="icon" href={favicon} />

  <!-- Canonical URL -->
  <link rel="canonical" href={canonicalUrl} />
</svelte:head>

{@render children?.()}
