<script lang="ts">
	import Button, { buttonVariants } from '$lib/components/ui/button.svelte';
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { cn } from '$lib/utils';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Siren } from 'lucide-svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';

	const id = $props.id();

	const incidents = [
		{ class: 'text-emerald-600', label: 'Resolved', value: 'i1' },
		{ class: 'text-yellow-500', label: 'In Progress', value: 'i2' },
		{ class: 'text-gray-500', label: 'Investigating', value: 'i3' },
		{ class: 'text-white', label: 'Identified', value: 'i4' }
	] as const;

	let value = $state('i4');

	let name = $state('');

	let note = $state('');

	function handleOnSubmit(e: Event) {
		e.preventDefault();
		console.log('Submitted form data:', { name, value, note });
	}

	const selected = $derived(incidents.find((i) => i.value === value));
</script>

{#snippet status(item: (typeof incidents)[number])}
	<span class="incidents-center flex gap-2">
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
			<Siren />
		</Empty.Media>
		<Empty.Title class=" text-gray-200">Let’s Get Started</Empty.Title>
		<Empty.Description class="text-gray-400">
			Get started by sending an alert, and you’ll receive real-time updates.
		</Empty.Description>
	</Empty.Header>
	<Empty.Content>
		<div class="flex gap-2">
			<Dialog.Root>
				<Dialog.Trigger class={cn('cursor-pointer', buttonVariants({ variant: 'outline' }))}
					>Send Alert</Dialog.Trigger
				>
				<Dialog.Content class="bg-zinc-900">
					<div class="flex flex-col items-center gap-2">
						<div
							class="flex size-10 shrink-0 items-center justify-center rounded-full border border-border"
							aria-hidden="true"
						>
							<Siren class="h-10 w-10 text-white" />
						</div>

						<Dialog.Header>
							<Dialog.Title class=" text-gray-300 sm:text-center">Send New Alert</Dialog.Title>
							<Dialog.Description class="text-gray-400 sm:text-center">
								Send an alert to notify your team.
							</Dialog.Description>
						</Dialog.Header>
					</div>

					<form onsubmit={handleOnSubmit} class="space-y-5">
						<div class="space-y-4">
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="logo">Title</Label>
								<Input
									class=" border-zinc-700 text-white"
									id="{id}-logo"
									placeholder="Service Disruption"
									type="text"
									bind:value={name}
									required
								/>
							</div>

							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="logo">Service Name</Label>
								<Input
									class=" border-zinc-700 text-white"
									id="{id}-logo"
									placeholder="HTTP"
									type="text"
									bind:value={name}
									required
								/>
							</div>
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="{id}-title">Status</Label>
								<Select.Root type="single" bind:value>
									<Select.Trigger
										{id}
										class="[&>span]:incidents-center w-full cursor-pointer border-zinc-700 text-white [&_svg:not([class*='text-'])]:text-zinc-200 [&>span]:flex [&>span]:gap-2 [&>span_svg]:shrink-0"
									>
										{#if selected}
											{@render status(selected)}
										{:else}
											Select a status
										{/if}
									</Select.Trigger>
									<Select.Content
										class="[&_*[data-select-item]>span]:incidents-center  bg-zinc-800 text-white [&_*[data-select-item]]:ps-2 [&_*[data-select-item]]:pe-8 [&_*[data-select-item]>span]:start-auto [&_*[data-select-item]>span]:end-2 [&_*[data-select-item]>span]:flex [&_*[data-select-item]>span]:gap-2 [&_*[data-select-item]>span>svg]:shrink-0"
									>
										{#each incidents as item (item.value)}
											<Select.Item
												class="cursor-pointer data-[highlighted]:bg-zinc-700 data-[highlighted]:text-white [&_svg:not([class*='text-'])]:text-gray-300"
												value={item.value}
											>
												{@render status(item)}
											</Select.Item>
										{/each}
									</Select.Content>
								</Select.Root>
							</div>
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="logo">Note</Label>
								<Input
									class=" border-zinc-700 text-white"
									id="{id}-note"
									placeholder="Investigating the issue"
									type="text"
									bind:value={note}
									required
								/>
							</div>
						</div>
						<Button class="mt-2 w-full cursor-pointer" type="submit" variant="outline"
							>Create</Button
						>
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
