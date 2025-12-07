<script lang="ts" module>
	import AlertTriangleIcon from '@lucide/svelte/icons/alert-triangle';
	import AudioWaveformIcon from '@lucide/svelte/icons/audio-waveform';
	import BellIcon from '@lucide/svelte/icons/bell';
	import CommandIcon from '@lucide/svelte/icons/command';
	import GalleryVerticalEndIcon from '@lucide/svelte/icons/gallery-vertical-end';
	import Globe2Icon from '@lucide/svelte/icons/globe-2';
	import MonitorIcon from '@lucide/svelte/icons/monitor';
	import Settings2Icon from '@lucide/svelte/icons/settings-2';
	import UsersIcon from '@lucide/svelte/icons/users';
	// This is sample data.
	const data = {
		user: {
			name: 'shadcn',
			email: 'm@example.com',
			avatar: '/avatars/shadcn.jpg'
		},
		teams: [
			{
				name: 'Acme Inc',
				logo: GalleryVerticalEndIcon,
				plan: 'Enterprise'
			},
			{
				name: 'Acme Corp.',
				logo: AudioWaveformIcon,
				plan: 'Startup'
			},
			{
				name: 'Evil Corp.',
				logo: CommandIcon,
				plan: 'Free'
			}
		],
		navMain: [
			{
				title: 'Monitors',
				url: '/monitors',
				icon: MonitorIcon
			},
			{
				title: 'Incidents',
				url: '/incidents',
				icon: AlertTriangleIcon
			},
			{
				title: 'Notifications',
				url: '/notifications',
				icon: BellIcon
			},
			{
				title: 'Status pages',
				url: '/status-pages',
				icon: Globe2Icon
			},
			{
				title: 'Teams',
				url: '/teams',
				icon: UsersIcon
			},
			{
				title: 'Settings',
				url: '/settings',
				icon: Settings2Icon
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
