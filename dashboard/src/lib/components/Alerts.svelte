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
	import { useCharacterLimit } from '$lib/hooks/use-character-limit.svelte';
	import Textarea from '$lib/components/ui/textarea.svelte';

	const id = $props.id();

	const incidents = [
		{ class: 'text-emerald-600', label: 'Resolved', value: 'i1' },
		{ class: 'text-yellow-500', label: 'In Progress', value: 'i2' },
		{ class: 'text-gray-500', label: 'Investigating', value: 'i3' },
		{ class: 'text-white', label: 'Identified', value: 'i4' }
	] as const;

	const services = [
		{ value: 'HTTPS', label: 'HTTPS' },
		{ value: 'HTTP', label: 'HTTP' },
		{ value: 'TCP', label: 'TCP' },
		{ value: 'DNS', label: 'DNS' },
		{ value: 'REDIS', label: 'REDIS' },
		{ value: 'SMTP', label: 'SMTP' },
		{ value: 'PING', label: 'PING' }
	];

	let value = $state('i4');

	let types = $state('HTTPS');

	let name = $state('');

	const bioLimit = useCharacterLimit(180, '');

	function handleOnSubmit(e: Event) {
		e.preventDefault();
		console.log('Submitted form data:', { name, value, types, note: bioLimit.value });
	}

	const selected = $derived(incidents.find((i) => i.value === value));

	const triggerContent = $derived(
		services.find((f) => f.value === types)?.label ?? services[0].label
	);
</script>

{#snippet status(item: (typeof incidents)[number])}
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
								<Label class="font-bold text-gray-300" for="{id}-title">Title</Label>
								<Input
									class=" border-zinc-700 text-white"
									id="{id}-title"
									placeholder="Service Disruption"
									type="text"
									bind:value={name}
									required
								/>
							</div>

							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="{id}-logo">Service Name</Label>

								<Select.Root type="single" name="monitorType" required bind:value={types}>
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
								<Label class="font-bold text-gray-300" for="{id}-title">Status</Label>
								<Select.Root type="single" bind:value>
									<Select.Trigger
										{id}
										class="w-full cursor-pointer border-zinc-700 text-white [&_svg:not([class*='text-'])]:text-zinc-200 [&>span]:flex [&>span]:items-center [&>span]:gap-2 [&>span_svg]:shrink-0"
									>
										{#if selected}
											{@render status(selected)}
										{:else}
											Select a status
										{/if}
									</Select.Trigger>
									<Select.Content
										class="bg-zinc-800  text-white [&_*[data-select-item]]:ps-2 [&_*[data-select-item]]:pe-8 [&_*[data-select-item]>span]:start-auto [&_*[data-select-item]>span]:end-2 [&_*[data-select-item]>span]:flex [&_*[data-select-item]>span]:items-center [&_*[data-select-item]>span]:gap-2 [&_*[data-select-item]>span>svg]:shrink-0"
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
								<div class="*:not-first:mt-2">
									<Label class="font-bold text-gray-300" for="{id}-note">Note</Label>
									<Textarea
										id="{id}-note"
										class=" border-zinc-700 text-white"
										bind:value={bioLimit.value}
										maxlength={bioLimit.maxLength}
										placeholder="Write a few sentences about incident..."
										aria-describedby="{id}-left-textarea"
										required
									/>
									<p
										id="{id}-left-textarea"
										class="mt-2 text-right text-xs text-muted-foreground"
										role="status"
										aria-live="polite"
									>
										<span class="tabular-nums">{bioLimit.maxLength - bioLimit.characterCount}</span>
										characters left
									</p>
								</div>
							</div>
						</div>
						<Button class="mt-2 w-full cursor-pointer" type="submit" variant="outline"
							>Send</Button
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
