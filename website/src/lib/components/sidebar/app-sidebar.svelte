<script lang="ts" module>
	// This is sample data.
	const data = {
		user: {
			name: 'shadcn',
			email: 'm@example.com',
			avatar: '/avatars/shadcn.jpg'
		},
		teams: [
			{
			  id: '123',
			name: 'Acme Inc',
			logo: 'lucide:layout-panel-left',
			plan: 'Enterprise'
		},
		{
			id: '456',
			name: 'Acme Corp.',
			logo: 'lucide:audio-waveform',
			plan: 'Startup'
		},
		{
      id: '789',
			name: 'Evil Corp.',
			logo: 'lucide:terminal',
			plan: 'Free'
		}
	],
	navMain: [
		{
			title: 'Monitors',
			url: '/monitors',
			icon: 'lucide:monitor'
		},
		{
			title: 'Incidents',
			url: '/incidents',
			icon: 'lucide:alert-triangle'
		},
		{
			title: 'Notifications',
			url: '/notifications',
			icon: 'lucide:bell'
		},
		{
			title: 'Status pages',
			url: '/status-pages',
			icon: 'lucide:globe'
		},
		{
			title: 'Teams',
			url: '/teams',
			icon: 'lucide:users'
		},
		{
			title: 'Settings',
			url: '/settings',
			icon: 'lucide:settings-2'
		}
	]
};
</script>

<script lang="ts">
	import NavMain from './nav-main.svelte';
	import NavUser from './nav-user.svelte';
	import TeamSwitcher from './team-switcher.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';
	let {
		ref = $bindable(null),
		collapsible = 'icon',
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root {collapsible} {...restProps}>
	<Sidebar.Header>
		<TeamSwitcher teams={data.teams} />
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={data.user} />
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
