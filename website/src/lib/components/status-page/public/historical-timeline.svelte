<script lang="ts">
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { formatDate, formatUpTo2Decimals, timelineTone } from '$lib/styles/status';
	import type { PublicTimelinePoint } from '../../../types';

	let { timeline, days = 90 }: { timeline: PublicTimelinePoint[]; days?: number } = $props();

	const visibleTimeline = $derived.by(() => (days ? timeline.slice(-days) : timeline));
</script>

<div class="mt-3 grid grid-flow-col auto-cols-fr gap-1">
	{#each visibleTimeline as point (point.day)}
		<HoverCard.Root openDelay={0} closeDelay={0}>
			<HoverCard.Trigger>
				<div
					class={`h-6 rounded-none ${timelineTone(point)}`}
					title={`${formatDate(point.day)} · ${point.success} success / ${point.fail} fail`}
				></div>
			</HoverCard.Trigger>
			<HoverCard.Content class="p-0">
				<div class="p-4 flex flex-col space-y-2">
					<div class="flex justify-between space-x-4">
						<span class="font-medium text-foreground">{formatDate(point.day)}</span>
						<span class="font-medium text-foreground">
							{#if point.fail + point.success === 0}
								<span class="font-medium text-foreground/50">No data</span>
							{:else}
								<span class="font-medium text-foreground/50">
									{formatUpTo2Decimals((point.success / (point.success + point.fail)) * 100)}%
									uptime
								</span>
							{/if}
						</span>
					</div>
					<div>
						{#if point.fail === 0}
							<span class="font-medium text-foreground/50">No down recoard found</span>
						{:else}
							<span class="font-medium text-foreground/50">
								{Math.round((point.success / (point.success + point.fail)) * 100)}% uptime
							</span>
						{/if}
					</div>
				</div>
			</HoverCard.Content>
		</HoverCard.Root>
	{/each}
</div>
