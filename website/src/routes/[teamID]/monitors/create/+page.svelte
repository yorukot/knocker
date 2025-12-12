<script lang="ts">
	import Icon from "@iconify/svelte";

	import CreateMonitorForm from "$lib/components/monitors/create/create-monitor-form.svelte";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Separator } from "$lib/components/ui/separator/index.js";

	import { page } from "$app/state";

	const teamId = $derived(page.params.teamID);
</script>

<svelte:head>
	<title>Create monitor</title>
</svelte:head>

<section class="space-y-4">
	<div class="flex flex-wrap items-start justify-between gap-3">
		<div>
			<h1 class="text-2xl font-semibold leading-tight text-foreground">Create a monitor</h1>
		</div>
		<Button variant="ghost" href={`/${teamId}/monitors`} size="sm" class="gap-2">
			<Icon icon="lucide:arrow-left" class="size-4" />
			Back to list
		</Button>
	</div>

	<Separator />

	<CreateMonitorForm {teamId} />
</section>
