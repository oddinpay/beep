<script lang="ts">
	import { goto } from '$app/navigation';
	import Header from '$lib/components/Header.svelte';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { Unplug, House } from 'lucide-svelte';

	let currentTab = 'tab-5';
</script>

<div class="flex min-h-screen items-start justify-center bg-black py-10">
	<div class="w-full max-w-7xl">
		<div class="w-full bg-black">
			<div class="mx-auto w-full max-w-7xl px-4 py-10">
				<Header />
			</div>
		</div>
		<Tabs value={currentTab} class="flex w-full flex-col gap-8 md:flex-row">
			<TabsList class="flex w-full gap-2 bg-transparent p-4 md:w-80 md:flex-col">
				<TabsTrigger
					value="tab-5"
					onclick={() => window.history.pushState({}, '', '/connect')}
					class="flex w-full cursor-pointer  items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<Unplug class="h-8 w-8" /> Connect
				</TabsTrigger>

				<TabsTrigger
					value="tab-0"
					onclick={() => goto('/')}
					class="flex w-full cursor-pointer  items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<House class="h-8 w-8" /> Home
				</TabsTrigger>
			</TabsList>

			<div class="grow cursor-auto rounded-lg border border-border bg-zinc-900 p-8 text-start">
				<TabsContent value="tab-5">
					<p class="text-base text-gray-200">Content for Tab 5</p>
				</TabsContent>
			</div>
		</Tabs>
	</div>
</div>
