<script lang="ts">
	import { goto } from '$app/navigation';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { SquareActivity, Siren, ShieldAlert, House, CalendarCheck } from 'lucide-svelte';
	import Header from '$lib/components/Header.svelte';

	let currentTab = 'tab-2';
</script>

<div class="flex min-h-screen items-start justify-center bg-black py-4">
	<div class="w-full max-w-7xl">
		<div class="mx-auto w-full max-w-7xl px-4 py-4">
			<Header />
		</div>

		<Tabs value={currentTab} class="flex w-full flex-col gap-8 md:flex-row">
			<TabsList
				class="flex w-full gap-2 self-start bg-transparent p-4 md:sticky md:top-6
				       md:max-h-[calc(100dvh-3rem)] md:w-80 md:flex-col md:overflow-y-hidden"
			>
				<TabsTrigger
					value="tab-0"
					onclick={() => goto('/')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:text-3xl"
				>
					<House class="h-8 w-8" /> Home
				</TabsTrigger>

				<TabsTrigger
					value="tab-1"
					onclick={() => goto('/monitors')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:text-3xl"
				>
					<SquareActivity class="h-8 w-8" /> Monitors
				</TabsTrigger>

				<TabsTrigger
					value="tab-2"
					onclick={() => goto('/incidents')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:text-3xl"
				>
					<ShieldAlert class="h-8 w-8" /> Incidents
				</TabsTrigger>

				<TabsTrigger
					value="tab-3"
					onclick={() => goto('/alerts')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:text-3xl"
				>
					<Siren class="h-8 w-8" /> Alerts
				</TabsTrigger>

				<TabsTrigger
					value="tab-4"
					onclick={() => goto('/schedules')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:text-3xl"
				>
					<CalendarCheck class="h-8 w-8" /> Schedules
				</TabsTrigger>
			</TabsList>

			<div class="grid w-[1100px] grid-cols-1 gap-4 px-10 md:grid-cols-2">
				<div
					class="relative col-span-1 rounded-lg border border-border bg-zinc-900 p-8 md:col-span-2 md:h-[calc(50dvh-3rem)] md:overflow-y-hidden"
				>
					<TabsContent value="tab-2" class="h-auto min-h-[300px]">
						<Tabs value="tab-2" class=" items-center">
							<TabsList
								class="h-auto gap-2 rounded-full border-b border-border bg-zinc-800 px-10 py-2 text-zinc-400"
							>
								<TabsTrigger
									value="tab-2"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
								 Create Incident
								</TabsTrigger>
								<TabsTrigger
									value="tab-3"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
								  Active
								</TabsTrigger>
								<TabsTrigger
									value="tab-4"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
								 Pending  
								</TabsTrigger>
								<TabsTrigger
									value="tab-5"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
								  Resolved  
								</TabsTrigger>

							</TabsList>
							<TabsContent value="tab-2">
								<p class="p-4 text-center text-xs text-white">No status page created</p>
							</TabsContent>
							<TabsContent value="tab-3">
								<p class="p-4 text-center text-xs text-white">Content for Tab 2</p>
							</TabsContent>
							<TabsContent value="tab-4">
								<p class="p-4 text-center text-xs text-white">Content for Tab 3</p>
							</TabsContent>
							<TabsContent value="tab-5">
								<p class="p-4 text-center text-xs text-white">Content for Tab 4</p>
							</TabsContent>
						</Tabs>
					</TabsContent>
				</div>
				<div
					class="min-h-[250px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-2" class="h-auto min-h-[150px]">
						<p class="text-base text-gray-200">Active Incidents</p>
					</TabsContent>
				</div>

				<div
					class="min-h-[250px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-2" class="h-auto min-h-[150px]">
						<p class="text-base text-gray-200">Total Incidents</p>
					</TabsContent>
				</div>
			</div>
		</Tabs>
	</div>
</div>
