<script lang="ts">
	import Input from '$lib/components/ui/input.svelte';
	import Label from '$lib/components/ui/label.svelte';

	const uid = $props.id();
</script>

<div class="*:not-first:mt-2">
	<Label class="text-gray-300" for={uid}>Connect to API</Label>
	<div class="flex rounded-lg shadow-xs shadow-black/5">
		<span
			class="inline-flex items-center rounded-s-lg border border-zinc-700 bg-zinc-900 px-3 text-sm text-muted-foreground"
		>
			https://
		</span>
		<Input
			id={uid}
			class="-ms-px rounded-s-none border-zinc-700 text-white shadow-none"
			placeholder="Hostname"
			type="text"
		/>
	</div>
</div>
