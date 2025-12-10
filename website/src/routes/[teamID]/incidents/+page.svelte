<script lang="ts">
	import Icon from '@iconify/svelte';

	import Incidents from '$lib/components/incidents/incidents.svelte';
	import { Button } from '$lib/components/ui/button/index.js';

	const incidents = [
		{
			id: 'inc-001',
			title: 'Elevated latency on Public API',
			status: 'investigating',
			severity: 'major',
			affected: ['Public API', 'Status page'],
			startedAt: 'Today, 09:12 UTC',
			updatedAt: '5m ago',
			duration: '38m',
			summary: 'Increased p95 latency after deploy. Rolling back and monitoring error budgets.'
		},
		{
			id: 'inc-002',
			title: 'Worker queue backlog in EU-West',
			status: 'monitoring',
			severity: 'minor',
			affected: ['Background jobs'],
			startedAt: 'Yesterday, 19:43 UTC',
			updatedAt: '1h ago',
			duration: '2h 14m',
			summary: 'Scaled up workers; backlog drained. Watching metrics for recurrence.'
		},
		{
			id: 'inc-003',
			title: 'Scheduled database maintenance',
			status: 'resolved',
			severity: 'maintenance',
			affected: ['Postgres cluster', 'Admin panel'],
			startedAt: 'Mar 14, 02:00 UTC',
			updatedAt: 'Mar 14, 02:45 UTC',
			duration: '45m',
			summary: 'Routine patching completed; read replicas back in rotation.'
		},
		{
			id: 'inc-004',
			title: 'Packet loss to edge POP (AP-South-1)',
			status: 'identified',
			severity: 'critical',
			affected: ['Edge POP', 'Health checks'],
			startedAt: 'Today, 07:58 UTC',
			updatedAt: '12m ago',
			duration: '1h 21m',
			summary: 'Upstream provider outage; rerouting traffic to alternate region.'
		}
	];
</script>

<section class="space-y-6">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="space-y-1">
			<p class="text-sm font-medium text-primary">Incidents</p>
			<h1 class="text-2xl font-semibold leading-tight text-foreground">Incident overview</h1>
			<p class="text-sm text-muted-foreground">
				Track live incidents, communicate updates, and close out postmortems in one place.
			</p>
		</div>
		<Button class="ms-auto" size="lg">
			<Icon icon="ri:add-line" class="size-4" />
			Create incident
		</Button>
	</div>

	<Incidents {incidents} />
</section>
