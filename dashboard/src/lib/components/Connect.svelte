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

	let status = $state('error');
</script>

<Empty.Root>
	<Empty.Header>
		<Empty.Media variant="icon">
			<IconPlugConnected />
		</Empty.Media>
		<Empty.Description class="text-gray-400">
			<div class="*:not-first:mt-4">
				<Badge variant="outline" class="gap-1.5  border-zinc-700">
					<span
						class={`size-1.5 rounded-full ${
							status === 'ok' ? 'bg-green-500' : status === 'warn' ? 'bg-yellow-500' : 'bg-red-500'
						}`}
						aria-hidden="true"
					></span>
					<Label class="text-md font-medium text-gray-300" for={uid}>Backend Connection</Label>
				</Badge>
				<form class="flex gap-2" onsubmit={handleSubmit}>
					<Input
						id={uid}
						name="apiUrl"
						class="min-w-55 flex-1 border-zinc-700 text-white"
						placeholder="e.g https://myapp.com:5532"
						type="text"
						bind:value={apiUrl}
					/>
					<button
						type="submit"
						class="inline-flex cursor-pointer items-center rounded-lg border border-input bg-background px-4 text-sm font-medium text-foreground ring-offset-background transition-shadow hover:bg-accent hover:text-foreground focus:z-10 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30 focus-visible:ring-offset-2 focus-visible:outline-hidden disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50"
					>
						{status === 'ok'
							? 'Connected'
							: apiUrl === ''
								? 'Connect'
								: status === 'warn'
									? 'Connecting...'
									: 'Connect'}
					</button>
				</form>
			</div>
		</Empty.Description>
	</Empty.Header>
</Empty.Root>
