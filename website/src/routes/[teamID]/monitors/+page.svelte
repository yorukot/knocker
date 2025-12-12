<script lang="ts">
	import Icon from '@iconify/svelte';

	import Monitors from '$lib/components/monitors/monitors.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<section class="space-y-6">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="space-y-1">
			<h1 class="text-2xl font-semibold leading-tight text-foreground">Monitors</h1>
		</div>
		<Button class="ms-auto" size="lg">
			<Icon icon="ri:add-line" class="size-4" />
			Create monitor
		</Button>
	</div>

	<Monitors monitors={data.monitors} />
</section>
