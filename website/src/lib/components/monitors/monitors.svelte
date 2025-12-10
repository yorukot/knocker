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

	type MonitorStatus = 'operational' | 'degraded' | 'down' | 'paused';
	type MonitorType = 'HTTP' | 'Ping' | 'TCP';

	type Monitor = {
		id: string;
		name: string;
		target: string;
		type: MonitorType;
		status: MonitorStatus;
		regions: string[];
		frequency: string;
		uptime: string;
		responseTime: string;
		lastChecked: string;
		lastIncident: string;
	};

	const statusStyles: Record<
		MonitorStatus,
		{
			label: string;
			class: string;
		}
	> = {
		operational: {
			label: 'Operational',
			class:
				'border-emerald-200 bg-emerald-500/10 text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/20 dark:text-emerald-100'
		},
		degraded: {
			label: 'Degraded',
			class:
				'border-amber-200 bg-amber-500/10 text-amber-700 dark:border-amber-400/30 dark:bg-amber-400/15 dark:text-amber-50'
		},
		down: {
			label: 'Down',
			class:
				'border-destructive/40 bg-destructive/10 text-destructive dark:border-destructive/30 dark:bg-destructive/20 dark:text-destructive-foreground'
		},
		paused: {
			label: 'Paused',
			class:
				'border-muted-foreground/30 bg-muted text-muted-foreground dark:bg-input/40 dark:border-input'
		}
	};

	const mockMonitors: Monitor[] = [
		{
			id: 'api-edge',
			name: 'Public API',
			target: 'https://api.knocker.dev/health',
			type: 'HTTP',
			status: 'operational',
			regions: ['us-east-1', 'eu-west-2', 'ap-south-1'],
			frequency: 'Every 30s',
			uptime: '99.99%',
			responseTime: '182 ms',
			lastChecked: '2m ago',
			lastIncident: 'Resolved 2d ago'
		},
		{
			id: 'marketing-site',
			name: 'Marketing site',
			target: 'https://knocker.dev',
			type: 'HTTP',
			status: 'degraded',
			regions: ['us-east-1', 'eu-west-2'],
			frequency: 'Every 1m',
			uptime: '99.80%',
			responseTime: '438 ms',
			lastChecked: 'Just now',
			lastIncident: 'Slow TTFB detected'
		},
		{
			id: 'bg-worker',
			name: 'Worker heartbeat',
			target: 'redis://internal:6379/ping',
			type: 'TCP',
			status: 'operational',
			regions: ['us-east-1'],
			frequency: 'Every 15s',
			uptime: '100%',
			responseTime: '41 ms',
			lastChecked: '45s ago',
			lastIncident: 'None in the last 7d'
		},
		{
			id: 'edge-latency',
			name: 'Edge latency',
			target: 'ping knocker-edge',
			type: 'Ping',
			status: 'paused',
			regions: ['eu-west-2', 'ap-south-1'],
			frequency: 'Paused',
			uptime: '—',
			responseTime: '—',
			lastChecked: 'Paused 1h ago',
			lastIncident: 'Maintenance window'
		}
	];

	let { monitors = mockMonitors }: { monitors?: Monitor[] } = $props();
</script>

<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
	{#each monitors as monitor (monitor.id)}
		<Card class="h-full shadow-xs">
			<CardHeader class="grid grid-cols-[1fr_auto] items-start gap-3">
				<div class="space-y-2">
					<div class="flex flex-wrap items-center gap-2">
						<CardTitle class="text-base">{monitor.name}</CardTitle>
						<Badge variant="secondary" class="rounded-md">
							{monitor.type}
						</Badge>
					</div>
					<CardDescription class="flex items-center gap-2 text-sm">
						<Icon icon="lucide:link" class="size-4 text-muted-foreground" />
						<span class="truncate text-muted-foreground">{monitor.target}</span>
					</CardDescription>
					<div class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
						<Badge variant="outline" class={`rounded-md ${statusStyles[monitor.status].class}`}>
							{statusStyles[monitor.status].label}
						</Badge>
						<Badge variant="outline" class="rounded-md border-dashed text-muted-foreground">
							{monitor.frequency}
						</Badge>
						<div class="flex items-center gap-1 rounded-full bg-muted px-2 py-1">
							<Icon icon="lucide:globe" class="size-3.5 text-muted-foreground" />
							<span>{monitor.regions.length} regions</span>
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
							<Icon icon="lucide:route" class="size-4" />
							View details
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<Icon icon="lucide:zap" class="size-4" />
							Run check now
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<Icon icon="lucide:radio-tower" class="size-4" />
							Edit monitor
						</DropdownMenu.Item>
						<DropdownMenu.Separator />
						<DropdownMenu.Item>
							<Icon icon="lucide:clock" class="size-4" />
							Schedule downtime
						</DropdownMenu.Item>
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</CardHeader>

			<CardContent class="space-y-4">
				<div class="grid grid-cols-2 gap-3">
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Uptime (30d)</p>
						<p class="text-lg font-semibold leading-tight">{monitor.uptime}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Avg. response</p>
						<p class="text-lg font-semibold leading-tight">{monitor.responseTime}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Regions</p>
						<p class="text-sm font-medium text-foreground">{monitor.regions.join(' • ')}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Last check</p>
						<p class="text-sm font-medium text-foreground">{monitor.lastChecked}</p>
					</div>
				</div>

				<Separator />

				<div class="flex items-start gap-3 rounded-xl bg-accent/50 px-3 py-2">
					<Icon icon="lucide:activity" class="mt-0.5 size-4 text-primary" />
					<div class="space-y-1">
						<p class="text-sm font-medium text-foreground">Last incident</p>
						<p class="text-sm text-muted-foreground">{monitor.lastIncident}</p>
					</div>
				</div>
			</CardContent>
		</Card>
	{/each}
</div>
