<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import Icon from '@iconify/svelte';
	import StatusPages from '$lib/components/status-page/status-pages.svelte';
	import { page } from '$app/state';

	/** @type {import('./$types').PageProps} */
	let { data } = $props();
</script>

<div class="flex flex-col gap-5">
	<div class="flex justify-between">
		<h1 class="text-2xl font-bold">Status pages</h1>

		<Button size="default" href="status-pages/new">
			<Icon icon="lucide:plus" />
			New status page
		</Button>
	</div>

	<div>
		<StatusPages statusPages={data.statusPages} teamID={page.params.teamID} />
	</div>
</div>
