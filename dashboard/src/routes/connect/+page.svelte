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

			<div class="grid w-[1100px] grid-cols-1 gap-4 px-10 md:grid-cols-2">
				<div
					class="relative col-span-1 rounded-lg border border-border bg-zinc-900 p-8 md:col-span-2 md:h-[calc(50dvh-3rem)] md:overflow-y-hidden"
				>
					<TabsContent value="tab-5" class="h-auto min-h-[300px]">
						<Tabs value="tab-5" class=" items-center">
							<TabsList
								class="h-auto gap-2 rounded-full border-b border-border bg-zinc-800 px-10 py-2 text-zinc-400"
							>
								<TabsTrigger
									value="tab-5"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
									Connect
								</TabsTrigger>
								<TabsTrigger
									value="tab-6"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
									API Keys
								</TabsTrigger>
							</TabsList>
							<TabsContent value="tab-5">
								<p class="p-4 text-center text-xs text-white">No status page created</p>
							</TabsContent>
							<TabsContent value="tab-6">
								<p class="p-4 text-center text-xs text-white">Content for Tab 2</p>
							</TabsContent>
						</Tabs>
					</TabsContent>
				</div>
			</div>
		</Tabs>
	</div>
</div>
