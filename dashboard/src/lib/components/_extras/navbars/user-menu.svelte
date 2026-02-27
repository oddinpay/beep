<script lang="ts">
	import Button from '$lib/components/ui/button.svelte';

	import LogOutIcon from '@lucide/svelte/icons/log-out';
	import { Unplug, BookOpenText } from 'lucide-svelte';
	import { Avatar, AvatarFallback, AvatarImage } from '$lib/components/ui/avatar';
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuGroup,
		DropdownMenuItem,
		DropdownMenuLabel,
		DropdownMenuSeparator,
		DropdownMenuTrigger
	} from '$lib/components/ui/dropdowns';
	import { goto } from '$app/navigation';
</script>

<DropdownMenu>
	<DropdownMenuTrigger>
		{#snippet child({ props })}
			<Button
				variant="ghost"
				class="h-auto cursor-pointer p-0 hover:bg-transparent hover:opacity-80"
				{...props}
			>
				<Avatar>
					<AvatarImage src="" alt="Profile image" />
					<AvatarFallback>KK</AvatarFallback>
				</Avatar>
			</Button>
		{/snippet}
	</DropdownMenuTrigger>
	<DropdownMenuContent class="max-w-64 bg-zinc-700" align="end">
		<DropdownMenuLabel class="flex min-w-0  flex-col">
			<span class="truncate text-sm font-medium text-zinc-100">Keith Kennedy</span>
			<span class="truncate text-xs font-normal text-zinc-300">
				k.kennedy@originui-svelte.com
			</span>
		</DropdownMenuLabel>
		<DropdownMenuSeparator class="bg-zinc-600" />
		<DropdownMenuGroup>
			<DropdownMenuSeparator class="bg-zinc-600" />
			<DropdownMenuItem
				onclick={() => goto('/alerts')}
				class="cursor-pointer text-zinc-100 focus:bg-zinc-600 focus:text-white"
			>
				<Unplug size={16} class="opacity-80" aria-hidden="true" />
				<span>Connect</span>
			</DropdownMenuItem>
		</DropdownMenuGroup>
		<DropdownMenuGroup>
			<DropdownMenuSeparator class="bg-zinc-600" />
			<DropdownMenuItem class="cursor-pointer text-zinc-100 focus:bg-zinc-600 focus:text-white">
				<BookOpenText size={16} class="opacity-80" aria-hidden="true" />
				<span>Docs</span>
			</DropdownMenuItem>
		</DropdownMenuGroup>
		<DropdownMenuSeparator class="bg-zinc-600" />
		<DropdownMenuItem class="cursor-pointer text-zinc-100 focus:bg-zinc-600 focus:text-white">
			<LogOutIcon size={16} class="opacity-80" aria-hidden="true" />
			<span>Sign out</span>
		</DropdownMenuItem>
	</DropdownMenuContent>
</DropdownMenu>
