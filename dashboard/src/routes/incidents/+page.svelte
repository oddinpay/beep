<script lang="ts">
	import { goto } from '$app/navigation';
	import { Tabs, TabsContent, TabsList, TabsTrigger } from '$lib/components/ui/tabs';
	import { SquareActivity, Siren, ShieldAlert, House, CalendarCheck } from 'lucide-svelte';
	import Header from '$lib/components/Header.svelte';
	import Incidents from '$lib/components/Incidents.svelte';

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
								<Incidents />
							</TabsContent>
							<TabsContent value="tab-3">
								<svg
									class="p-4 text-center text-white"
									xmlns="http://www.w3.org/2000/svg"
									width="128"
									height="128"
									viewBox="0 0 32 32"
									><path
										fill="currentColor"
										d="M30.48 4.575H32v1.52h-1.52Zm0 10.66h-3.05v1.53h1.52v4.57H25.9v-1.52h-1.52v3.04h4.57v4.58h1.53v-4.58H32v-1.52h-1.52zm0-6.09h-1.53v1.52h-1.52v1.53h1.52v1.52h1.53v-1.52H32v-1.53h-1.52z"
									/><path
										fill="currentColor"
										d="M25.9 16.765h1.53v3.05H25.9Zm-1.52-10.67h1.52v1.53h-1.52Zm-1.52 6.1h1.52v3.04h-1.52Zm-10.67 13.71v-7.62h-1.52v4.57H9.14v1.53H7.62v3.05h1.52v1.52h13.72v-1.52h1.52v-3.05h-1.52v-1.53h-1.53v-4.57h-1.52v7.62zm9.14-10.67h1.53v1.53h-1.53Zm0-4.57h1.53v1.53h-1.53Zm-1.52-7.62h1.52v1.53h-1.52Zm-1.52 13.72h1.52v1.52h-1.52Zm0-6.1h1.52v1.53h-1.52Z"
									/><path
										fill="currentColor"
										d="M16.76 19.815h1.53v1.52h-1.53Zm-3.05 3.04h4.58v1.53h-4.58Zm0-7.62h4.58v1.53h-4.58Zm0 4.58h1.53v1.52h-1.53Zm0-12.19v1.52h-1.52v1.52h1.52v1.53h1.53v-1.53h1.52v-1.52h-1.52v-1.52zm-1.52 9.14h1.52v1.52h-1.52Zm-3.05-4.57h3.05v1.52H9.14Zm0 3.04h1.53v1.53H9.14Zm0-12.19h1.53v1.53H9.14Zm-1.52 10.67h1.52v1.52H7.62Zm-4.57 1.52v1.53h1.52v4.57H1.52v-1.52H0v3.04h4.57v3.05H6.1v-3.05h1.52v-1.52H6.1v-6.1zm3.05-6.09h1.52v1.52H6.1Z"
									/><path
										fill="currentColor"
										d="M1.52 16.765h1.53v3.05H1.52Zm0-12.19v1.52H0v1.53h1.52v1.52h1.53v-1.52h1.52v-1.53H3.05v-1.52z"
									/></svg
								>
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
						<p class="text-base font-semibold text-zinc-200">Active Incidents</p>
					</TabsContent>
				</div>

				<div
					class="min-h-[250px] rounded-lg border border-border bg-zinc-900 p-8 md:overflow-y-hidden"
				>
					<TabsContent value="tab-2" class="h-auto min-h-[150px]">
						<p class="text-base font-semibold text-zinc-200">Total Incidents</p>
					</TabsContent>
				</div>
			</div>
		</Tabs>
	</div>
</div>
