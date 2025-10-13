<script lang="ts">
	import Button, { buttonVariants } from '$lib/components/ui/button.svelte';
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { cn } from '$lib/utils';

	import * as Select from '$lib/components/ui/select/index.js';

	import { ShieldAlert } from 'lucide-svelte';

	import * as Empty from '$lib/components/ui/empty/index.js';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';

	const items = [
		{ class: 'text-emerald-600', label: 'Resolved', value: 's1' },
		{ class: 'text-amber-500', label: 'In Progress', value: 's3' },
		{ class: 'text-gray-500', label: 'Investigating', value: 's4' },
		{ class: 'text-red-500', label: 'Identified', value: 's5' }
	] as const;

	let value = $state('s1');

	let name = $state('');

	const selected = $derived(items.find((i) => i.value === value));

	const uid = $props.id();

	function handleOnSubmit(e: Event) {
		e.preventDefault();
		console.log('Submitted form data:', { name, value });
	}
</script>

{#snippet status(item: (typeof items)[number])}
	<span class="flex items-center gap-2">
		<svg
			width="8"
			height="8"
			fill="currentColor"
			viewBox="0 0 8 8"
			xmlns="http://www.w3.org/2000/svg"
			class={item.class}
			aria-hidden="true"
		>
			<circle cx="4" cy="4" r="4" />
		</svg>
		<span class="truncate">{item.label}</span>
	</span>
{/snippet}

<Empty.Root>
	<Empty.Header>
		<Empty.Media variant="icon">
			<ShieldAlert />
		</Empty.Media>
		<Empty.Title class=" text-gray-200">Let’s Get Started</Empty.Title>
		<Empty.Description class="text-gray-400">
			Get started by creating an incident, and you’ll start seeing updates.
		</Empty.Description>
	</Empty.Header>
	<Empty.Content>
		<div class="flex gap-2">
			<Dialog.Root>
				<Dialog.Trigger class={cn('cursor-pointer', buttonVariants({ variant: 'outline' }))}
					>Create Incident</Dialog.Trigger
				>
				<Dialog.Content class="bg-zinc-900">
					<div class="flex flex-col items-center gap-2">
						<div
							class="flex size-10 shrink-0 items-center justify-center rounded-full border border-border"
							aria-hidden="true"
						>
							<ShieldAlert class="h-10 w-10 text-white" />
						</div>

						<Dialog.Header>
							<Dialog.Title class=" text-gray-300 sm:text-center">Add New Incident</Dialog.Title>
							<Dialog.Description class="text-gray-400 sm:text-center">
								Set up and publish your incident.
							</Dialog.Description>
						</Dialog.Header>
					</div>

					<form onsubmit={handleOnSubmit} class="space-y-5">
						<div class="space-y-4">
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="logo">Title</Label>
								<Input
									class=" border-zinc-700 text-white"
									id="{uid}-logo"
									placeholder="Payment errors"
									type="text"
									bind:value={name}
									required
								/>
							</div>
							<div class="space-y-2">
								<Label for={uid}>Status select</Label>
								<Select.Root type="single" bind:value>
									<Select.Trigger
										id={uid}
										class="w-full border-zinc-700 text-white [&>span]:flex [&>span]:items-center [&>span]:gap-2 [&>span_svg]:shrink-0"
									>
										{#if selected}
											{@render status(selected)}
										{:else}
											Select a status
										{/if}
									</Select.Trigger>
									<Select.Content
										class="[&_*[data-select-item]]:ps-2 [&_*[data-select-item]]:pe-8 [&_*[data-select-item]>span]:start-auto [&_*[data-select-item]>span]:end-2 [&_*[data-select-item]>span]:flex [&_*[data-select-item]>span]:items-center [&_*[data-select-item]>span]:gap-2 [&_*[data-select-item]>span>svg]:shrink-0 [&_*[data-select-item]>span>svg]:text-muted-foreground/80"
									>
										{#each items as item (item.value)}
											<Select.Item value={item.value}>
												{@render status(item)}
											</Select.Item>
										{/each}
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

								<Input
									class="border-zinc-700 text-white"
									id="{uid}-description"
									placeholder={value === 'HTTP' || value === 'HTTPS'
										? 'https://oddinpay.com'
										: 'IP address or domain'}
									type="text"
									required
								/>
							</div>

							{#if value === 'TCP' || value === 'REDIS' || value === 'SMTP'}
								<div class="space-y-2">
									<Label class="font-bold text-gray-300" for="slug">Port</Label>
									<Input
										class="border-zinc-700 text-white"
										id="{uid}-description"
										placeholder="443"
										type="number"
										required
									/>
								</div>
							{/if}
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
