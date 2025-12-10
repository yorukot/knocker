<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';

	type IncidentStatus = 'investigating' | 'identified' | 'monitoring' | 'resolved';
	type IncidentSeverity = 'critical' | 'major' | 'minor' | 'maintenance';

	type Incident = {
		id: string;
		title: string;
		status: IncidentStatus;
		severity: IncidentSeverity;
		affected: string[];
		startedAt: string;
		updatedAt: string;
		duration: string;
		summary: string;
	};

	const statusStyles: Record<
		IncidentStatus,
		{
			label: string;
			class: string;
		}
	> = {
		investigating: {
			label: 'Investigating',
			class:
				'border-amber-200 bg-amber-500/10 text-amber-700 dark:border-amber-400/30 dark:bg-amber-400/20 dark:text-amber-50'
		},
		identified: {
			label: 'Identified',
			class:
				'border-blue-200 bg-blue-500/10 text-blue-700 dark:border-blue-400/30 dark:bg-blue-400/20 dark:text-blue-50'
		},
		monitoring: {
			label: 'Monitoring',
			class:
				'border-purple-200 bg-purple-500/10 text-purple-700 dark:border-purple-400/30 dark:bg-purple-400/20 dark:text-purple-50'
		},
		resolved: {
			label: 'Resolved',
			class:
				'border-emerald-200 bg-emerald-500/10 text-emerald-700 dark:border-emerald-400/30 dark:bg-emerald-400/20 dark:text-emerald-50'
		}
	};

	const severityStyles: Record<
		IncidentSeverity,
		{
			label: string;
			class: string;
		}
	> = {
		critical: {
			label: 'Critical',
			class:
				'border-destructive/40 bg-destructive/10 text-destructive dark:border-destructive/30 dark:bg-destructive/20 dark:text-destructive-foreground'
		},
		major: {
			label: 'Major',
			class:
				'border-orange-200 bg-orange-500/10 text-orange-700 dark:border-orange-400/30 dark:bg-orange-400/20 dark:text-orange-50'
		},
		minor: {
			label: 'Minor',
			class:
				'border-amber-200 bg-amber-100 text-amber-700 dark:border-amber-300/30 dark:bg-amber-300/20 dark:text-amber-50'
		},
		maintenance: {
			label: 'Maintenance',
			class:
				'border-muted-foreground/30 bg-muted text-muted-foreground dark:bg-input/40 dark:border-input'
		}
	};

	const mockIncidents: Incident[] = [
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

	let { incidents = mockIncidents }: { incidents?: Incident[] } = $props();
</script>

<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
	{#each incidents as incident (incident.id)}
		<Card class="h-full shadow-xs">
			<CardHeader class="grid grid-cols-[1fr_auto] items-start gap-3">
				<div class="space-y-2">
					<div class="flex flex-wrap items-center gap-2">
						<CardTitle class="text-base">{incident.title}</CardTitle>
						<Badge variant="outline" class={`rounded-md ${severityStyles[incident.severity].class}`}>
							{severityStyles[incident.severity].label}
						</Badge>
					</div>
					<CardDescription class="flex flex-wrap items-center gap-2 text-sm">
						<Icon icon="lucide:alert-triangle" class="size-4 text-muted-foreground" />
						<span class="text-muted-foreground">{incident.affected.join(' • ')}</span>
					</CardDescription>
					<div class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
						<Badge variant="outline" class={`rounded-md ${statusStyles[incident.status].class}`}>
							{statusStyles[incident.status].label}
						</Badge>
						<div class="flex items-center gap-1 rounded-full bg-muted px-2 py-1">
							<Icon icon="lucide:shield-half" class="size-3.5 text-muted-foreground" />
							<span>{incident.severity === 'critical' ? 'Priority 0' : 'Priority 1'}</span>
						</div>
						<div class="flex items-center gap-1 rounded-full bg-muted px-2 py-1">
							<Icon icon="lucide:messages-square" class="size-3.5 text-muted-foreground" />
							<span>Updates: 3</span>
						</div>
					</div>
				</div>

				<DropdownMenu.Root>
					<DropdownMenu.Trigger>
						{#snippet child({ props })}
							<Button variant="ghost" size="icon-sm" class="-m-1" aria-label="More actions" {...props}>
								<Icon icon="lucide:ellipsis-vertical" class="size-4" />
							</Button>
						{/snippet}
					</DropdownMenu.Trigger>
					<DropdownMenu.Content class="w-44" align="end" sideOffset={6}>
						<DropdownMenu.Item>
							<Icon icon="lucide:activity" class="size-4" />
							View timeline
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<Icon icon="lucide:bell" class="size-4" />
							Notify subscribers
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<Icon icon="lucide:clock" class="size-4" />
							Add update
						</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item>
							<Icon icon="lucide:check-circle-2" class="size-4" />
							Mark resolved
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</CardHeader>

			<CardContent class="space-y-4">
				<div class="grid grid-cols-2 gap-3">
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Started</p>
						<p class="text-sm font-medium text-foreground">{incident.startedAt}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Last update</p>
						<p class="text-sm font-medium text-foreground">{incident.updatedAt}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Duration</p>
						<p class="text-lg font-semibold leading-tight">{incident.duration}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Status</p>
						<p class="text-sm font-medium text-foreground capitalize">{incident.status}</p>
					</div>
				</div>

				<Separator />

				<div class="flex items-start gap-3 rounded-xl bg-accent/50 px-3 py-2">
					{#if incident.status === 'resolved'}
						<Icon icon="lucide:check-circle-2" class="mt-0.5 size-4 text-emerald-500" />
					{:else if incident.status === 'monitoring'}
						<Icon icon="lucide:hourglass" class="mt-0.5 size-4 text-blue-500" />
					{:else if incident.status === 'identified'}
						<Icon icon="lucide:dot" class="mt-0.5 size-4 text-orange-500" />
					{:else}
						<Icon icon="lucide:alert-triangle" class="mt-0.5 size-4 text-amber-500" />
					{/if}
					<div class="space-y-1">
						<p class="text-sm font-medium text-foreground">Update</p>
						<p class="text-sm text-muted-foreground">{incident.summary}</p>
					</div>
				</div>
			</CardContent>
		</Card>
	{/each}
</div>
