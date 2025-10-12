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
	const bannerImageHandler = useImageUpload({ initialImage: '/profile-bg.jpg' });
	const profileImageHandler = useImageUpload({ initialImage: '/avatar-72-01.jpg' });

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
							<svg
								xmlns="http://www.w3.org/2000/svg"
								aria-hidden="true"
								width="28"
								class="size-6"
								height="28"
								viewBox="0 0 512 512"
								><path
									fill="white"
									d="m39.93 149.25l197.4 95.32c5.14 2.45 12 3.73 18.79 3.73s13.65-1.28 18.78-3.73l197.4-95.32c10.38-5 10.38-13.18 0-18.2L274.9 35.73c-5.13-2.45-12-3.73-18.78-3.73s-13.65 1.28-18.79 3.73l-197.4 95.32c-10.38 5.02-10.38 13.18 0 18.2"
								/><path
									fill="white"
									d="M472.3 246.9s-36.05-17.38-40.83-19.72s-6.07-2.21-11.09.12s-145.6 70.23-145.6 70.23a45.7 45.7 0 0 1-18.78 3.74c-6.77 0-13.65-1.29-18.78-3.74c0 0-136.85-66-143.27-69.18C87 225 85 225 78.67 228l-39 18.78c-10.38 5-10.38 13.19 0 18.2L237.1 360.3c5.13 2.45 12 3.73 18.78 3.73s13.65-1.28 18.79-3.73l197.4-95.3c10.61-4.92 10.61-13.08.23-18.1"
								/><path
									fill="white"
									d="M472.3 362.75s-36.05-17.38-40.83-19.75s-6.07-2.21-11.09.12S274.9 413.5 274.9 413.5a45.7 45.7 0 0 1-18.78 3.73c-6.77 0-13.65-1.28-18.79-3.73c0 0-136.85-66-143.26-69.18c-7-3.39-9-3.39-15.29-.35l-39 18.78c-10.39 5-10.39 13.18 0 18.2l197.4 95.32c5.13 2.56 12 3.73 18.78 3.73s13.65-1.28 18.78-3.73L472.18 381c10.5-5.07 10.5-13.23.12-18.25"
								/></svg
							>
						</div>
						<Dialog.Header>
							<Dialog.Title class="text-gray-300 sm:text-center">New Status Page</Dialog.Title>
							<Dialog.Description class="text-gray-400 sm:text-center">
								Set up and publish your status page.
							</Dialog.Description>
						</Dialog.Header>
					</div>

					<form class="space-y-5">
						<div class="space-y-4">
							<div class="space-y-2">
								<Label class="font-bold text-gray-300" for="signup-name">Full name</Label>
								<Input id="signup-name" placeholder="Matt Welsh" type="text" required />
							</div>
							<div class="space-y-2">
								<Label for="signup-email">Email</Label>
								<Input id="signup-email" placeholder="hi@yourcompany.com" type="email" required />
							</div>
							<div class="space-y-2">
								<Label for="signup-password">Password</Label>
								<Input
									id="signup-password"
									placeholder="Enter your password"
									type="password"
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
