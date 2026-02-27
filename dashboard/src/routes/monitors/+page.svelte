<script lang="ts">
	import { goto } from '$app/navigation';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { SquareActivity, Siren, ShieldAlert, House, CalendarCheck } from 'lucide-svelte';
	import Header from '$lib/components/Header.svelte';
	import { Gauge } from '$lib/components/ui/gauge';
	import Monitor from '$lib/components/Monitor.svelte';
	import NotMonitor from '$lib/components/NotMonitor.svelte';

	let currentTab = 'tab-1';
</script>

<div class="flex min-h-screen items-start justify-center overflow-hidden bg-black">
	<div class="fixed mx-auto w-full max-w-7xl py-4">
		<div class="px-4 py-4">
			<Header />
		</div>

		<Tabs
			value={currentTab}
			class="flex max-h-[calc(100dvh-3rem)]  w-full flex-col gap-8 md:flex-row"
		>
			<TabsList class="sticky top-6 flex w-80 flex-col gap-2 self-start bg-transparent p-4">
				<TabsTrigger
					value="tab-0"
					onclick={() => goto('/')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:px-6 xl:text-3xl"
				>
					<House class="h-8 w-8" /> Home
				</TabsTrigger>

				<TabsTrigger
					value="tab-1"
					onclick={() => goto('/monitors')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:px-6 xl:text-3xl"
				>
					<SquareActivity class="h-8 w-8" /> Monitors
				</TabsTrigger>

				<TabsTrigger
					value="tab-2"
					onclick={() => goto('/incidents')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:px-6 xl:text-3xl"
				>
					<ShieldAlert class="h-8 w-8" /> Incidents
				</TabsTrigger>

				<TabsTrigger
					value="tab-3"
					onclick={() => goto('/alerts')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:px-6 xl:text-3xl"
				>
					<Siren class="h-8 w-8" /> Alerts
				</TabsTrigger>

				<TabsTrigger
					value="tab-4"
					onclick={() => goto('/schedules')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md py-4 text-2xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white xl:px-6 xl:text-3xl"
				>
					<CalendarCheck class="h-8 w-8" /> Schedules
				</TabsTrigger>
			</TabsList>

			<div class="grid w-full grid-cols-1 gap-4 overflow-y-auto px-10 pb-20 md:grid-cols-2">
				<div
					class="relative col-span-1 rounded-lg border border-border bg-zinc-900 p-8 md:col-span-2 md:min-h-[480px] md:overflow-y-hidden"
				>
					<TabsContent value="tab-1" class="h-auto min-h-[300px]">
						<Tabs value="tab-1" class=" items-center">
							<TabsList
								class="h-auto gap-2 rounded-full border-b border-border bg-zinc-800 px-10 py-2 text-zinc-400"
							>
								<TabsTrigger
									value="tab-1"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
									Overview
								</TabsTrigger>
								<TabsTrigger
									value="tab-2"
									class="relative cursor-pointer after:absolute  after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5  hover:text-white hover:after:bg-white data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=active]:shadow-none data-[state=active]:after:bg-white data-[state=active]:hover:text-white"
								>
									Create Monitor
								</TabsTrigger>
							</TabsList>
							<TabsContent value="tab-1">
								<NotMonitor />
							</TabsContent>
							<TabsContent value="tab-2">
								<Monitor />
							</TabsContent>
						</Tabs>
					</TabsContent>
				</div>
				<div
					class="min-h-[250px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-1" class="h-auto min-h-[150px]">
						<p class="text-base font-semibold text-zinc-200">Up</p>
						<Gauge
							colors={{ primary: 'stroke-green-600', secondary: 'stroke-green-200' }}
							class="text-white"
							show_value
							size="lg"
							value={100}
						/>
					</TabsContent>
				</div>

				<div
					class="min-h-[250px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-1" class="h-auto min-h-[150px]">
						<p class="text-base font-semibold text-zinc-200">Down</p>
						<Gauge
							colors={{ primary: 'stroke-red-700', secondary: 'stroke-red-200' }}
							class="text-white"
							show_value
							size="lg"
							value={100}
						/>
					</TabsContent>
				</div>
			</div>
		</Tabs>
	</div>
</div>
