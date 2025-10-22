<script lang="ts">
	import Button from '$lib/components/ui/button.svelte';
	import { SquareActivity } from 'lucide-svelte';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import ArrowUpRightIcon from '@lucide/svelte/icons/arrow-up-right';
</script>

<Empty.Root>
	<Empty.Header>
		<Empty.Media variant="icon">
			<SquareActivity />
		</Empty.Media>
		<Empty.Title class=" text-gray-200">Monitor Not Found</Empty.Title>
		<Empty.Description class="text-gray-400">
			The monitor you are looking for doesn't exist or has been moved.
		</Empty.Description>
	</Empty.Header>
	<Button variant="link" class="text-gray-400" size="sm">
		<a href="#/">
			Learn More <ArrowUpRightIcon class="inline" />
		</a>
	</Button>
</Empty.Root>
