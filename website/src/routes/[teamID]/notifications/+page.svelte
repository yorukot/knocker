<script lang="ts">
	import Icon from '@iconify/svelte';

	import Notifications from '$lib/components/notifications/notifications.svelte';
	import { Button } from '$lib/components/ui/button/index.js';

	const channels = [
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
</script>

<section class="space-y-6">
	<div class="flex flex-wrap items-center justify-between gap-3">
		<div class="space-y-1">
			<p class="text-sm font-medium text-primary">Notifications</p>
			<h1 class="text-2xl font-semibold leading-tight text-foreground">Notification channels</h1>
			<p class="text-sm text-muted-foreground">
				Manage how alerts reach your team. Test, pause, or tweak delivery in one place.
			</p>
		</div>
		<Button class="ms-auto" size="lg">
			<Icon icon="ri:add-line" class="size-4" />
			Create channel
		</Button>
	</div>

	<Notifications channels={channels} />
</section>
