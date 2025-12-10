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

	type ChannelType = 'Email' | 'Slack' | 'SMS' | 'Webhook';
	type ChannelStatus = 'active' | 'paused' | 'degraded';

	type NotificationChannel = {
		id: string;
		name: string;
		type: ChannelType;
		status: ChannelStatus;
		target: string;
		description: string;
		sentThisWeek: string;
		deliveryRate: string;
		lastSent: string;
		lastIssue: string;
		primaryUse: string;
	};

	const statusStyles: Record<
		ChannelStatus,
		{
			label: string;
			class: string;
		}
	> = {
		active: {
			label: 'Active',
			class:
				'border-emerald-200 bg-emerald-500/10 text-emerald-700 dark:border-emerald-400/30 dark:bg-emerald-400/20 dark:text-emerald-50'
		},
		paused: {
			label: 'Paused',
			class:
				'border-muted-foreground/30 bg-muted text-muted-foreground dark:bg-input/40 dark:border-input'
		},
		degraded: {
			label: 'Issues',
			class:
				'border-amber-200 bg-amber-500/10 text-amber-700 dark:border-amber-400/30 dark:bg-amber-400/20 dark:text-amber-50'
		}
	};

	const typeIcons: Record<ChannelType, string> = {
		Email: 'lucide:mail',
		Slack: 'ri:slack-line',
		SMS: 'lucide:smartphone',
		Webhook: 'lucide:webhook'
	};

	const mockChannels: NotificationChannel[] = [
		{
			id: 'email-alerts',
			name: 'On-call email',
			type: 'Email',
			status: 'active',
			target: 'alerts@knocker.dev',
			description: 'Primary on-call rotation via email',
			sentThisWeek: '124',
			deliveryRate: '99.6%',
			lastSent: '2m ago',
			lastIssue: 'None in last 30d',
			primaryUse: 'Critical incidents'
		},
		{
			id: 'slack-status',
			name: 'Slack #ops-status',
			type: 'Slack',
			status: 'degraded',
			target: '#ops-status (Slack)',
			description: 'Broadcasts incident updates to ops channel',
			sentThisWeek: '87',
			deliveryRate: '97.3%',
			lastSent: '6m ago',
			lastIssue: 'Rate-limited 1h ago',
			primaryUse: 'All incidents'
		},
		{
			id: 'sms-escalation',
			name: 'SMS escalation',
			type: 'SMS',
			status: 'active',
			target: '+1 (415) •••• ••••',
			description: 'Escalation path for Sev0 outages',
			sentThisWeek: '18',
			deliveryRate: '100%',
			lastSent: '18m ago',
			lastIssue: 'None in last 90d',
			primaryUse: 'Sev0 / Sev1'
		},
		{
			id: 'webhook-statuspage',
			name: 'Status page webhook',
			type: 'Webhook',
			status: 'paused',
			target: 'https://hooks.statuspage.dev/...',
			description: 'Syncs incidents to public status page',
			sentThisWeek: '—',
			deliveryRate: '—',
			lastSent: 'Paused yesterday',
			lastIssue: 'Maintenance window',
			primaryUse: 'Status page sync'
		}
	];

	let { channels = mockChannels }: { channels?: NotificationChannel[] } = $props();
</script>

<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
	{#each channels as channel (channel.id)}
		<Card class="h-full shadow-xs">
			<CardHeader class="grid grid-cols-[1fr_auto] items-start gap-3">
				<div class="space-y-2">
					<div class="flex flex-wrap items-center gap-2">
						<CardTitle class="text-base">{channel.name}</CardTitle>
						<Badge variant="secondary" class="rounded-md">
							<Icon icon={typeIcons[channel.type]} class="size-4" />
							<span>{channel.type}</span>
						</Badge>
					</div>
					<CardDescription class="flex flex-wrap items-center gap-2 text-sm">
						<Icon icon="lucide:bell" class="size-4 text-muted-foreground" />
						<span class="text-muted-foreground">{channel.description}</span>
					</CardDescription>
					<div class="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
						<Badge variant="outline" class={`rounded-md ${statusStyles[channel.status].class}`}>
							{statusStyles[channel.status].label}
						</Badge>
						<Badge variant="outline" class="rounded-md border-dashed text-muted-foreground">
							{channel.target}
						</Badge>
						<div class="flex items-center gap-1 rounded-full bg-muted px-2 py-1">
							<Icon icon="lucide:shield-check" class="size-3.5 text-muted-foreground" />
							<span>{channel.primaryUse}</span>
						</div>
					</div>
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
					<DropdownMenu.Content class="w-48" align="end" sideOffset={6}>
						<DropdownMenu.Item>
							<Icon icon="lucide:send" class="size-4" />
							Test notification
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<Icon icon="lucide:radio-tower" class="size-4" />
							View delivery log
						</DropdownMenu.Item>
						<DropdownMenu.Item>
							<Icon icon="lucide:clock" class="size-4" />
							Schedule quiet hours
						</DropdownMenu.Item>
						<DropdownMenu.Separator />
						{#if channel.status === 'paused'}
							<DropdownMenu.Item>
								<Icon icon="lucide:check-circle-2" class="size-4" />
								Resume channel
							</DropdownMenu.Item>
						{:else}
							<DropdownMenu.Item>
								<Icon icon="lucide:pause-circle" class="size-4" />
								Pause channel
							</DropdownMenu.Item>
						{/if}
					</DropdownMenu.Content>
				</DropdownMenu.Root>
			</CardHeader>

			<CardContent class="space-y-4">
				<div class="grid grid-cols-2 gap-3">
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Sent this week</p>
						<p class="text-lg font-semibold leading-tight">{channel.sentThisWeek}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Delivery rate</p>
						<p class="text-lg font-semibold leading-tight">{channel.deliveryRate}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Last sent</p>
						<p class="text-sm font-medium text-foreground">{channel.lastSent}</p>
					</div>
					<div class="rounded-lg border bg-muted/50 px-3 py-2">
						<p class="text-xs text-muted-foreground">Last issue</p>
						<p class="text-sm font-medium text-foreground">{channel.lastIssue}</p>
					</div>
				</div>

				<Separator />

				<div class="flex items-start gap-3 rounded-xl bg-accent/50 px-3 py-2">
					{#if channel.status === 'active'}
						<Icon icon="lucide:check-circle-2" class="mt-0.5 size-4 text-emerald-500" />
					{:else if channel.status === 'degraded'}
						<Icon icon="lucide:alert-triangle" class="mt-0.5 size-4 text-amber-500" />
					{:else}
						<Icon icon="lucide:pause-circle" class="mt-0.5 size-4 text-muted-foreground" />
					{/if}
					<div class="space-y-1">
						<p class="text-sm font-medium text-foreground">Channel health</p>
						<p class="text-sm text-muted-foreground">
							{channel.status === 'active'
								? 'Delivering normally — no action needed.'
								: channel.status === 'degraded'
									? 'Recent delivery issues detected. Consider testing or rotating keys.'
									: 'Channel paused. Resume when ready to send updates.'}
						</p>
					</div>
				</div>
			</CardContent>
		</Card>
	{/each}
</div>
