<script lang="ts">
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import Badge from '$lib/components/ui/badge.svelte';
	import IconPlugConnected from '@tabler/icons-svelte/icons/plug-connected';

	const uid = $props.id();

	let apiUrl = $state('');

	function handleSubmit(event: Event) {
		event.preventDefault();
		console.log('API URL submitted:', apiUrl);
	}
</script>

<Empty.Root>
	<Empty.Header>
		<Empty.Media variant="icon">
			<IconPlugConnected />
		</Empty.Media>
		<Empty.Description class="text-gray-400">
			<div class="*:not-first:mt-4">
				<Badge variant="outline" class="gap-1.5">
					<span class="size-1.5 rounded-full bg-red-500" aria-hidden="true"></span>
					<Label class="text-md font-medium text-gray-300" for={uid}>Backend API Instance</Label>
				</Badge>
				<form class="flex gap-2" onsubmit={handleSubmit}>
					<Input
						id={uid}
						name="apiUrl"
						class="flex-1 border-zinc-700 text-white"
						placeholder="API URL"
						type="text"
						bind:value={apiUrl}
					/>
					<button
						type="submit"
						class="inline-flex cursor-pointer items-center rounded-lg border border-input bg-background px-3 text-sm font-medium text-foreground ring-offset-background transition-shadow hover:bg-accent hover:text-foreground focus:z-10 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
					>
						Connect
					</button>
				</form>
			</div>
		</Empty.Description>
	</Empty.Header>
</Empty.Root>
