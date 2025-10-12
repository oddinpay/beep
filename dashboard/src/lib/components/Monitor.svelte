<script lang="ts">
	import Button, { buttonVariants } from '$lib/components/ui/button.svelte';
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { cn } from '$lib/utils';

	import * as Select from '$lib/components/ui/select/index.js';

	import { SquareActivity, Siren, ShieldAlert, House, CalendarCheck } from 'lucide-svelte';

	import * as Empty from '$lib/components/ui/empty/index.js';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';

	import { useImageUpload } from '$lib/hooks/use-image-upload.svelte';

	import ImagePlus from '@lucide/svelte/icons/image-plus';

	const profileImageHandler = useImageUpload({
		initialImage: ''
	});

	const id = $props.id();

	const fruits = [
		{ value: 'HTTPS', label: 'HTTPS' },
		{ value: 'HTTP', label: 'HTTP' },
		{ value: 'TCP', label: 'TCP' },
		{ value: 'DNS', label: 'DNS' },
		{ value: 'REDIS', label: 'REDIS' }
	];

	let value = $state('HTTPS');

	let name = $state('');

	function handleOnSubmit(e: Event) {
		e.preventDefault();
		console.log('Submitted form data:', { name, value });
	}

	const triggerContent = $derived(fruits.find((f) => f.value === value)?.label ?? fruits[0].label);
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
									<Select.Trigger class="w-full border-zinc-700 text-white">
										{triggerContent}
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											<Select.Label>Protocol</Select.Label>
											{#each fruits as fruit (fruit.value)}
												<Select.Item
													id="{id}-monitorType"
													class="cursor-pointer"
													value={fruit.value}
													label={fruit.label}
												>
													{fruit.label}
												</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							</div>

							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="slug">URL</Label>
								<Input
									class="border-zinc-700 text-white"
									id="{id}-description"
									placeholder="Status page"
									type="text"
									required
								/>
							</div>
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

{#snippet Avatar()}
	<label class="mt-10 cursor-pointer px-6" aria-label="Upload profile picture">
		<div
			class="relative flex size-20 items-center justify-center overflow-hidden rounded-full border-4 border-zinc-600 bg-zinc-700 shadow-xs shadow-black/10"
		>
			{#if profileImageHandler.previewUrl}
				<img
					src={profileImageHandler.previewUrl}
					class="size-full object-cover"
					width={80}
					height={80}
					alt="Profile avatar"
				/>
			{/if}
			<button
				type="button"
				class="absolute inset-0 flex cursor-pointer items-center justify-center rounded-full bg-zinc-800 text-white backdrop-blur-sm transition-[color,box-shadow,background-color,backdrop-filter] outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 {profileImageHandler.previewUrl
					? 'hidden hover:flex'
					: 'flex'} hover:opacity-50"
				onclick={profileImageHandler.handleThumbnailClick}
				aria-label="Change profile picture"
			>
				<ImagePlus size={16} aria-hidden="true" />
			</button>
			<input
				type="file"
				bind:this={profileImageHandler.fileInput}
				bind:files={profileImageHandler.files}
				class="hidden"
				accept="image/*"
				aria-label="Upload profile picture"
			/>
		</div>
	</label>
{/snippet}
