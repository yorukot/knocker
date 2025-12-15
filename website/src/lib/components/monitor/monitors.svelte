<script lang="ts">
	import type { Monitor } from '../../../types';
	import { Card } from '$lib/components/ui/card';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import Icon from '@iconify/svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { monitorTarget } from './utils';

	let { monitors }: { monitors: Monitor[] } = $props();
</script>

<div>
	<div class="flex flex-col gap-2">
		{#each monitors as monitor, monitorIndex (monitor.id)}
			<Card id={'monitor-card-' + monitorIndex} class="py-2 px-4">
				<div class="flex justify-between items-center gap-2">
					<div class="flex-1 min-w-0">
						<div class="flex items-center gap-2 min-w-0">
							<span class="w-3 h-3 rounded-full bg-success shrink-0"></span>

							<h2 class="text-md font-semibold truncate min-w-0 flex-1">
								{monitor.name}
							</h2>
						</div>

						<div class="min-w-0">
							<p class="text-sm text-muted-foreground truncate">
								<Badge variant="secondary" class="rounded-sm shrink-0">
									{monitor.type.toUpperCase()}
								</Badge>
								{monitorTarget(monitor)}
							</p>
						</div>
					</div>

					<DropdownMenu.Root>
						<DropdownMenu.Trigger class="shrink-0">
							<Icon icon="lucide:more-vertical" />
						</DropdownMenu.Trigger>
						<DropdownMenu.Content>
							<DropdownMenu.Group>
								<DropdownMenu.Item>
									<Icon icon="lucide:eye" /> View details
								</DropdownMenu.Item>
								<DropdownMenu.Item>
									<Icon icon="lucide:edit" /> Edit
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item variant="destructive">
									<Icon icon="lucide:trash" />
									Delete
								</DropdownMenu.Item>
							</DropdownMenu.Group>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>
			</Card>
		{/each}
	</div>
</div>
