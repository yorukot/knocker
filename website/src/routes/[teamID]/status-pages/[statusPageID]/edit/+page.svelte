<script lang="ts">
	import DeleteStatusPage from '$lib/components/status-page/edit/delete-status-page.svelte';
	import { Button } from '$lib/components/ui/button';
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import {
		statusPageUpsertRequestSchema,
		type StatusPageUpsertValues
	} from '$lib/components/status-page/edit/schema';
	import { reporter } from '@felte/reporter-svelte';
	import GeneralStatusPage from '$lib/components/status-page/edit/general-status-page.svelte';
	import type { StatusPage, StatusPageElement } from '$lib/types/status-page.js';
	import DndStatusPage from '$lib/components/status-page/edit/dnd-status-page.svelte';

	const schema = statusPageUpsertRequestSchema;

	/** @type {import('$types').PageProps} */
	let { data } = $props();

	// svelte-ignore state_referenced_locally
	let statusPage: StatusPage = data.statusPage.statusPage;
	// svelte-ignore state_referenced_locally
	const elements: StatusPageElement[] = data.statusPage.elements ?? [];

	const { form, setFields } = createForm<StatusPageUpsertValues>({
		extend: [validator({ schema }), reporter()],
		initialValues: {
			name: statusPage.title,
			slug: statusPage.slug,
			elements: elements ?? []
		}
	});
</script>

<form use:form>
	<div class="flex flex-col gap-6">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold">Edit status page</h1>
				<p class="text-sm text-muted-foreground">Organize what your public status page shows.</p>
			</div>
		</div>
		<GeneralStatusPage />
		<DndStatusPage
			availableMonitors={data.monitors}
			statusPageId={statusPage.id}
			{elements}
			{setFields}
		/>
		<div class="flex items-center gap-2 justify-end">
			<DeleteStatusPage />
			<Button type="submit">Save changes</Button>
		</div>
	</div>
</form>
