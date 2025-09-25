<script lang="ts">
  import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
  } from "$lib/components/ui/tabs";

  let activeTab = "tab-1";
  let direction: "left" | "right" = "left";

  const tabsOrder = ["tab-1", "tab-3"];

  function handleChange(newValue: string) {
    const oldIndex = tabsOrder.indexOf(activeTab);
    const newIndex = tabsOrder.indexOf(newValue);

    direction = newIndex > oldIndex ? "left" : "right";
    activeTab = newValue;
  }
</script>

<Tabs value={activeTab} class="items-center" onValueChange={handleChange}>
  <TabsList
    class="border-border text-foreground h-auto gap-2 rounded-none border-b bg-transparent px-0 py-1"
  >
    {#each tabsOrder as t, i}
      <TabsTrigger
        value={t}
        class={`cursor-pointer hover:bg-accent hover:text-foreground transition-colors duration-150 ease-in-out data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 after:transform after:scale-x-0 after:transition-transform after:duration-200 after:ease-in-out data-[state=active]:after:scale-x-100 data-[state=active]:after:bg-primary data-[state=active]:bg-transparent data-[state=active]:shadow-none ${
          direction === "left" ? "after:origin-left" : "after:origin-right"
        }`}
      >
        Tab {i + 1}
      </TabsTrigger>
    {/each}
  </TabsList>
  {#each tabsOrder as t, i}
    <TabsContent value={t}>
      {#if i === 0}
        <div class="p-4 text-center">
          <h3 class="text-sm font-medium">Overview (Index: {i})</h3>
          <p class="text-muted-foreground mt-2 text-xs">
            Content for Tab {i + 1}: summary and quick stats.
          </p>
        </div>
      {:else if i === 1}
        <div class="p-4 text-center">
          <h3 class="text-sm font-medium">Details (Index: {i})</h3>
          <ul class="mt-2 text-xs list-disc list-inside text-muted-foreground">
            <li>Detail A</li>
            <li>Detail B</li>
            <li>Detail C</li>
          </ul>
        </div>
      {:else}
        <div class="p-4 text-center">
          <h3 class="text-sm font-medium">Settings (Index: {i})</h3>
          <p class="text-muted-foreground mt-2 text-xs">
            Content for Tab {i + 1}: configuration and toggles.
          </p>
        </div>
      {/if}
    </TabsContent>
  {/each}
</Tabs>
