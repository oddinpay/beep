<script lang="ts">
	import Button, { buttonVariants } from '$lib/components/ui/button.svelte';
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';
	import * as Dialog from '$lib/components/ui/dialog';
	import { cn } from '$lib/utils';

	import * as Empty from '$lib/components/ui/empty/index.js';
	import IconFileOrientation from '@tabler/icons-svelte/icons/file-orientation';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';

	import { useCharacterLimit } from '$lib/hooks/use-character-limit.svelte';
	import { useImageUpload } from '$lib/hooks/use-image-upload.svelte';

	import Check from '@lucide/svelte/icons/check';
	import ImagePlus from '@lucide/svelte/icons/image-plus';
	import X from '@lucide/svelte/icons/x';

	const bioLimit = useCharacterLimit(
		180,
		'Hey, I am Margaret, a web developer who loves turning ideas into amazing websites!'
	);
	const profileImageHandler = useImageUpload({
		initialImage: ''
	});

	const id = $props.id();
</script>

<Empty.Root>
	<Empty.Header>
		<Empty.Media variant="icon">
			<IconFileOrientation />
		</Empty.Media>
		<Empty.Title class=" text-gray-200">Let’s Get Started</Empty.Title>
		<Empty.Description class="text-gray-400">
			You haven't created a status page yet. Get started by creating your first status page.
		</Empty.Description>
	</Empty.Header>
	<Empty.Content>
		<div class="flex gap-2">
			<Dialog.Root>
				<Dialog.Trigger class={cn('cursor-pointer', buttonVariants({ variant: 'outline' }))}
					>Create Status Page</Dialog.Trigger
				>
				<Dialog.Content class="bg-zinc-900">
					<div class="flex flex-col items-center gap-2">
						<div
							class="flex size-11 shrink-0 items-center justify-center rounded-full border border-border"
							aria-hidden="true"
						>
							{@render Avatar()}
						</div>

						<Dialog.Header>
							<Dialog.Title class="mt-10 text-gray-300 sm:text-center">Favicon</Dialog.Title>
							<Dialog.Description class="text-gray-400 sm:text-center">
								Set up and publish your status page.
							</Dialog.Description>
						</Dialog.Header>
					</div>

					<form class="space-y-5">
						<div class="space-y-4">
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="signup-name">Title</Label>
								<Input
									class="border-zinc-600 text-white"
									id="{id}-signup-name"
									placeholder="Beep"
									type="text"
									required
								/>
							</div>

							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="signup-name">Title</Label>
								<Input
									class="border-zinc-600 text-white"
									id="{id}-signup-name"
									placeholder="Beep"
									type="text"
									required
								/>
							</div>
						</div>
					</form>

					<div
						class="flex items-center gap-3 before:h-px before:flex-1 before:bg-border after:h-px after:flex-1 after:bg-border"
					></div>

					<Button variant="outline">Publish Page</Button>
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
	<div class="mt-10 px-6">
		<div
			class="relative flex size-20 items-center justify-center overflow-hidden rounded-full border-4 border-background bg-muted shadow-xs shadow-black/10"
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
				class="absolute flex size-8 cursor-pointer items-center justify-center rounded-full bg-black/60 text-white transition-[color,box-shadow] outline-none hover:bg-black/80 focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
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
	</div>
{/snippet}
