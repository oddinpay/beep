<script lang="ts">
	import { goto } from '$app/navigation';
	import Header from '$lib/components/Header.svelte';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { SquareActivity, Siren, ShieldAlert, CalendarCheck, House } from 'lucide-svelte';

	let currentTab = 'tab-0';
</script>

<div class="flex min-h-screen items-start justify-center bg-black py-4">
	<div class="w-full max-w-7xl">
		<div class="mx-auto w-full max-w-7xl px-4 py-4">
			<Header />
		</div>

		<Tabs value={currentTab} class="flex w-full flex-col gap-8 md:flex-row">
			<TabsList
				class="flex w-full gap-2 self-start bg-transparent p-4 md:sticky md:top-6
				       md:max-h-[calc(100vh-3rem)] md:w-80 md:flex-col md:overflow-y-hidden"
			>
				<TabsTrigger
					value="tab-0"
					onclick={() => goto('/')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<House class="h-8 w-8" /> Home
				</TabsTrigger>

				<TabsTrigger
					value="tab-1"
					onclick={() => goto('/monitors')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<SquareActivity class="h-8 w-8" /> Monitors
				</TabsTrigger>

				<TabsTrigger
					value="tab-2"
					onclick={() => goto('/incidents')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<ShieldAlert class="h-8 w-8" /> Incidents
				</TabsTrigger>

				<TabsTrigger
					value="tab-3"
					onclick={() => goto('/alerts')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<Siren class="h-8 w-8" /> Alerts
				</TabsTrigger>

				<TabsTrigger
					value="tab-4"
					onclick={() => goto('/schedules')}
					class="flex w-full cursor-pointer items-center justify-start gap-3 rounded-md px-6 py-4 text-3xl font-bold text-white transition data-[state=active]:bg-transparent data-[state=active]:text-white data-[state=inactive]:text-zinc-400 data-[state=inactive]:hover:bg-zinc-800 data-[state=inactive]:hover:text-white"
				>
					<CalendarCheck class="h-8 w-8" /> Schedules
				</TabsTrigger>
			</TabsList>

			<div class="grid w-full grid-cols-1 gap-4 md:grid-cols-2">
				<div
					class="col-span-1 rounded-lg border border-border bg-zinc-900 p-8 md:col-span-2 md:h-[calc(40dvh-3rem)] md:overflow-y-hidden"
				>
					<TabsContent value="tab-0" class="h-auto min-h-[300px]">
						<p class="text-base text-gray-200">Content for Tab 0</p>
					</TabsContent>
				</div>

				<div
					class="min-h-[200px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-0" class="h-auto min-h-[150px]">
						<p class="text-base text-gray-200">Content for Tab 1</p>
					</TabsContent>
				</div>

				<div
					class="min-h-[200px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-0" class="h-auto min-h-[150px]">
						<p class="text-base text-gray-200">Content for Tab 2</p>
					</TabsContent>
				</div>

				<div
					class="min-h-[200px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-0" class="h-auto min-h-[150px]">
						<p class="text-base text-gray-200">Content for Tab 3</p>
					</TabsContent>
				</div>

				<div
					class="min-h-[200px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-0" class="h-auto min-h-[150px]">
						<p class="text-base text-gray-200">Content for Tab 4</p>
					</TabsContent>
				</div>
			</div>
		</Tabs>
	</div>
</div>
