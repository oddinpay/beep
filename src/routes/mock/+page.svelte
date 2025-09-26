<script>
  import { onMount } from "svelte";
  import { writable } from "svelte/store";

  let name = "world";
  let messages = writable([]);

  onMount(() => {
    const evtSource = new EventSource("http://127.0.0.1:8976/status");
    evtSource.onmessage = function (event) {
      console.log("New message", event.data);
      var dataobj = JSON.parse(event.data);
      messages.update((arr) => arr.concat(dataobj));
    };
  });
</script>

<h1>Hello {name}!</h1>

{#each $messages as m}
  <p>
    {m.protocol} - {m.status} - {m.description}
  </p>
{/each}

