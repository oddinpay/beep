<script lang="ts">
	import Button, { buttonVariants } from '$lib/components/ui/button.svelte';
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { cn } from '$lib/utils';

	import * as Select from '$lib/components/ui/select/index.js';

	import { SquareActivity } from 'lucide-svelte';

	import * as Empty from '$lib/components/ui/empty/index.js';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';

	const id = $props.id();

	const services = [
		{ value: 'HTTPS', label: 'HTTPS' },
		{ value: 'HTTP', label: 'HTTP' },
		{ value: 'TCP', label: 'TCP' },
		{ value: 'DNS', label: 'DNS' },
		{ value: 'REDIS', label: 'REDIS' },
		{ value: 'SMTP', label: 'SMTP' },
		{ value: 'PING', label: 'PING' }
	];

	let value = $state('HTTPS');
	let name = $state();
	let interval = $state();
	let url = $state();
	let host = $state();

	function handleOnSubmit(e: Event) {
		e.preventDefault();
		console.log('Submitted form data:', { name, value, interval, url, host });
	}

	const triggerContent = $derived(
		services.find((f) => f.value === value)?.label ?? services[0].label
	);
</script>

<Empty.Root>
	<Empty.Header>
		<Empty.Media variant="icon">
			<SquareActivity />
		</Empty.Media>
		<Empty.Title class=" text-gray-200">Let’s Get Started</Empty.Title>
		<Empty.Description class="text-gray-400">
			Get started by creating an uptime monitor, and you’ll start seeing real-time updates.
		</Empty.Description>
	</Empty.Header>
	<Empty.Content>
		<div class="flex gap-2">
			<Dialog.Root>
				<Dialog.Trigger class={cn('cursor-pointer', buttonVariants({ variant: 'outline' }))}
					>Add Monitor</Dialog.Trigger
				>
				<Dialog.Content class="bg-zinc-900">
					<div class="flex flex-col items-center gap-2">
						<div
							class="flex size-10 shrink-0 items-center justify-center rounded-full border border-border"
							aria-hidden="true"
						>
							<SquareActivity class="h-10 w-10 text-white" />
						</div>

						<Dialog.Header>
							<Dialog.Title class=" text-gray-300 sm:text-center">Add New Monitor</Dialog.Title>
							<Dialog.Description class="text-gray-400 sm:text-center">
								Set up and publish your uptime monitor.
							</Dialog.Description>
						</Dialog.Header>
					</div>

					<form onsubmit={handleOnSubmit} class="space-y-5">
						<div class="space-y-4">
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="logo">Name</Label>
								<Input
									class=" border-zinc-700 text-white"
									id="{id}-logo"
									placeholder="oddinpay"
									type="text"
									bind:value={name}
									required
								/>
							</div>
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="title">Monitor Type</Label>
								<Select.Root type="single" name="monitorType" required bind:value>
									<Select.Trigger
										class="w-full cursor-pointer border-zinc-700 text-white [&_svg:not([class*='text-'])]:text-zinc-200"
									>
										{triggerContent}
									</Select.Trigger>
									<Select.Content class="bg-zinc-800 text-white">
										<Select.Group>
											<Select.Label class="text-zinc-400">services</Select.Label>
											{#each services as type (type.value)}
												<Select.Item
													id="{id}-monitorType"
													class="cursor-pointer  data-[highlighted]:bg-zinc-700 data-[highlighted]:text-white [&_svg:not([class*='text-'])]:text-gray-300"
													value={type.value}
													label={type.label}
												>
													{type.label}
												</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							</div>

							<div class="space-y-2">
								{#if value === 'PING'}
									<Label class="font-bold text-gray-300" for="slug">Host</Label>
								{:else if value === 'DNS'}
									<Label class="font-bold text-gray-300" for="slug">Host</Label>
								{:else if value === 'REDIS'}
									<Label class="font-bold text-gray-300" for="slug">Host</Label>
								{:else if value === 'SMTP'}
									<Label class="font-bold text-gray-300" for="slug">Host</Label>
								{:else if value === 'TCP'}
									<Label class="font-bold text-gray-300" for="slug">Host</Label>
								{:else}
									<Label class="font-bold text-gray-300" for="slug">URL</Label>
								{/if}

								{#if value === 'HTTP' || value === 'HTTPS'}
									<Input
										class="border-zinc-700 text-white"
										id="{id}-url"
										placeholder="https://oddinpay.com"
										type="text"
										bind:value={url}
										required
									/>
								{:else}
									<Input
										class="border-zinc-700 text-white"
										id="{id}-host"
										placeholder="IP address or domain"
										type="text"
										bind:value={host}
										required
									/>
								{/if}
							</div>

							{#if value === 'TCP' || value === 'REDIS' || value === 'SMTP'}
								<div class="space-y-2">
									<Label class="font-bold text-gray-300" for="slug">Port</Label>
									<Input
										class="border-zinc-700 text-white"
										id="{id}-description"
										placeholder="443"
										type="number"
										required
									/>
								</div>
							{/if}
							{#if value === 'REDIS' || value === 'SMTP'}
								<div class="space-y-2">
									<Label class="font-bold text-gray-300" for="slug">Username</Label>
									<Input
										class="border-zinc-700 text-white"
										id="{id}-description"
										placeholder="sachinsenal"
										type="text"
										required
									/>
								</div>

								<div class="space-y-2">
									<Label class="font-bold text-gray-300" for="slug">Password</Label>
									<Input
										class="border-zinc-700 text-white"
										id="{id}-description"
										placeholder="supersecret"
										type="text"
										required
									/>
								</div>
							{/if}
						</div>
						<div class="space-y-2">
							<Label class="font-bold text-gray-300" for="logo"
								>Interval (Check every 5 seconds)</Label
							>
							<Input
								class=" border-zinc-700 text-white"
								id="{id}-logo"
								placeholder="5s"
								type="number"
								bind:value={interval}
								required
							/>
						</div>
						<Button class="mt-2 w-full cursor-pointer" type="submit" variant="outline">Save</Button>
					</form>
				</Dialog.Content>
			</Dialog.Root>
		</div>
	</Empty.Content>
	<Button variant="link" class="text-gray-400" size="sm">
		<a href="#/">
			Learn More <ArrowUpRightIcon class="inline" />
		</a>
	</Button>
</Empty.Root>
