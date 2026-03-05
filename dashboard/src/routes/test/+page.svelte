<script lang="ts">
	import { useQuery } from 'convex-svelte';
	import { api } from '../../convex/_generated/api';

	const query = useQuery(api.status.get, {});
</script>

{#if query.isLoading}
	Loading...
{:else if query.error}
	failed to load: {query.error.toString()}
{:else}
	<ul>
		{#each query.data as status}
			<li>
				<span>Host: {status.host}</span>
			</li>
			<li>
				<span>assigned by {status.assigner}</span>
			</li>
			<li>
				<span>status: {status.protocol}</span>
			</li>
			<li>
				<span>Name: {status.name}</span>
			</li>
			<li>
				<span>Interval: {status.interval}</span>
			</li>
		{/each}
	</ul>
{/if}
