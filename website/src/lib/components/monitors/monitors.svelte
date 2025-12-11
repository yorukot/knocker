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
			dot: string;
		}
	> = {
		operational: {
			label: 'Operational',
			class:
				'border-emerald-200 bg-emerald-500/10 text-emerald-700 dark:border-emerald-500/30 dark:bg-emerald-500/20 dark:text-emerald-100',
			dot: 'bg-emerald-500 dark:bg-emerald-300'
		},
		degraded: {
			label: 'Degraded',
			class:
				'border-amber-200 bg-amber-500/10 text-amber-700 dark:border-amber-400/30 dark:bg-amber-400/15 dark:text-amber-50',
			dot: 'bg-amber-500 dark:bg-amber-300'
		},
		down: {
			label: 'Down',
			class:
				'border-destructive/40 bg-destructive/10 text-destructive dark:border-destructive/30 dark:bg-destructive/20 dark:text-destructive-foreground',
			dot: 'bg-destructive'
		},
		paused: {
			label: 'Paused',
			class:
				'border-muted-foreground/30 bg-muted text-muted-foreground dark:bg-input/40 dark:border-input',
			dot: 'bg-muted-foreground'
		}
	};

	let { monitors }: { monitors?: Monitor[] } = $props();
</script>

<svelte:head>
	<title>Monitors</title>
</svelte:head>

<div class="space-y-4">
	{#each monitors as monitor (monitor.id)}
		<Card class="py-2 px-4">
			<CardContent class="p-0">
				<div class="flex gap-3 flex-row tems-start justify-between">
					<div class="space-y-2">
						<div class="flex flex-wrap items-center gap-2">
							<span
								class={`inline-block size-2.5 rounded-full ${statusStyles[monitor.status].dot}`}
								aria-hidden="true"
							></span>
							<CardTitle class="text-base">{monitor.name}</CardTitle>
							<Badge variant="secondary" class="rounded-md">{monitor.type}</Badge>
						</div>
						<div class="flex items-center gap-2 text-sm text-muted-foreground">
							<Icon icon="lucide:link" class="size-4" />
							<span class="truncate">{monitor.target}</span>
						</div>
					</div>

					<div class="flex items-center gap-2 justify-end">

					  <div class="text-left hidden md:block">
							<p class="text-xs text-muted-foreground">Uptime</p>
							<p class="text-sm font-semibold leading-tight">{monitor.uptime}</p>
						</div>

						<div class="text-left hidden md:block">
							<p class="text-xs text-muted-foreground">Last check</p>
							<p class="text-sm font-semibold leading-tight">{monitor.lastChecked}</p>
						</div>

						<DropdownMenu.Root>
							<DropdownMenu.Trigger>
								{#snippet child({ props })}
									<Button
										variant="ghost"
										size="icon-sm"
										class="-m-1"
										aria-label="More actions"
										{...props}
									>
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
									<Icon icon="lucide:radio-tower" class="size-4" />
									Edit monitor
								</DropdownMenu.Item>
								<DropdownMenu.Separator />
								<DropdownMenu.Item variant="destructive">
									<Icon icon="lucide:trash" class="size-4" />
									Delete
								</DropdownMenu.Item>
							</DropdownMenu.Content>
						</DropdownMenu.Root>
					</div>
				</div>
			</CardContent>
		</Card>
	{/each}
</div>
