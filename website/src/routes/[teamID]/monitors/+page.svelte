<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import type { MonitorsData } from './proxy+page';
	import Monitors from '$lib/components/monitor/monitors.svelte';

	let { data }: { data: MonitorsData } = $props();
</script>

<div class="flex flex-col gap-5">
	<div class="flex justify-between">
		<h1 class="text-2xl font-bold">Monitors</h1>

		<Button size="default">
			<Icon icon="lucide:plus" />
			New monitor
		</Button>
	</div>

	<div>
		<Monitors monitors={data.monitors} />
	</div>
</div>
