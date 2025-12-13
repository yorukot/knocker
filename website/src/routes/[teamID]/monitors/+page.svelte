<script lang="ts">
	import Icon from '@iconify/svelte';

	import Monitors from '$lib/components/monitors/monitors.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import type { PageData } from './$types';

	import { page } from '$app/state';

	let { data }: { data: PageData } = $props();

	const teamId = $derived(page.params.teamID);
</script>

<section class="space-y-6">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="space-y-1">
			<h1 class="text-2xl font-semibold leading-tight text-foreground">Monitors</h1>
		</div>
		<Button class="ms-auto" size="lg" href={`/${teamId}/monitors/create`}>
			<Icon icon="ri:add-line" class="size-4" />
			Create monitor
		</Button>
	</div>

	<Monitors monitors={data.monitors} />
</section>
