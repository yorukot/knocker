<script lang="ts">
	import { formatSli, statusMeta } from '$lib/styles/status';
	import type { PublicStatusPageMonitor } from '../../../types';
	import HistoricalTimeline from './historical-timeline.svelte';

	let { monitor, days = 90 }: { monitor: PublicStatusPageMonitor; days?: number } = $props();

	function selectSli() {
		if (days <= 30) return { value: monitor.uptimeSli30 };
		if (days <= 60) return { value: monitor.uptimeSli60 };
		return { value: monitor.uptimeSli90 };
	}
</script>

<div class="rounded-md border bg-background/70 p-3">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="flex items-center gap-2">
			<span class="text-lg font-semibold">{monitor.name}</span>
		</div>
		{#if monitor.type === 'current_status_indicator'}
			<span class={`text-sm font-medium ${statusMeta[monitor.status ?? 'up'].tone}`}>
				{statusMeta[monitor.status ?? 'up'].label}
			</span>
		{:else}
			{@const sli = selectSli()}
			<span class="text-xs text-muted-foreground">
				Uptime {formatSli(sli.value)}
			</span>
		{/if}
	</div>
	{#if monitor.type === 'historical_timeline' && monitor.timeline?.length}
		<HistoricalTimeline timeline={monitor.timeline} {days} />
	{/if}
</div>
