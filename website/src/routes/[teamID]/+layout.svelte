<script lang="ts">
	import '../layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import AppSidebar from '$lib/components/sidebar/app-sidebar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { page } from '$app/state';
	import { ModeWatcher } from 'mode-watcher';

	let { children, data }: import('./$types').LayoutProps = $props();

	const formatSegment = (value: string) =>
		value
			.replace(/-/g, ' ')
			.replace(/\b\w/g, (c) => c.toUpperCase());

	const crumbs = $derived.by(() => {
		const segments = page.url.pathname.split('/').filter(Boolean);
		if (segments.length <= 1) return [];

		const [teamId, ...rest] = segments;
		const items: { label: string; href: string }[] = [];

		const section = rest[0];
		const matchedNav = data.navItems?.find((item) => item.url.endsWith(`/${section}`));

		items.push({
			label: matchedNav?.title ?? formatSegment(section),
			href: `/${[teamId, section].join('/')}`
		});

		rest.slice(1).forEach((segment, index) => {
			const href = `/${[teamId, ...rest.slice(0, index + 2)].join('/')}`;
			items.push({ label: formatSegment(segment), href });
		});

		return items;
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<ModeWatcher />

<Sidebar.Provider>
	<AppSidebar
		user={data.user}
		teams={[data.team]}
		currentTeamId={data.team.id}
		navItems={data.navItems}
	/>
	<Sidebar.Inset>
		<header
			class="group-has-data-[collapsible=icon]/sidebar-wrapper:h-12 flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear"
		>
			<div class="flex items-center gap-2 px-4">
				<Sidebar.Trigger class="-ms-1" />
				<Separator orientation="vertical" class="me-2 data-[orientation=vertical]:h-4" />
					{#if crumbs.length}
					<Breadcrumb.Root>
						<Breadcrumb.List>
							{#each crumbs as crumb, index (crumb.href)}
								<Breadcrumb.Item>
									{#if index < crumbs.length - 1}
										<Breadcrumb.Link href={crumb.href}>{crumb.label}</Breadcrumb.Link>
									{:else}
										<Breadcrumb.Page>{crumb.label}</Breadcrumb.Page>
									{/if}
								</Breadcrumb.Item>
								{#if index < crumbs.length - 1}
									<Breadcrumb.Separator class="hidden md:block" />
								{/if}
							{/each}
						</Breadcrumb.List>
					</Breadcrumb.Root>
				{/if}
			</div>
		</header>
		<main class="flex flex-1 flex-col gap-4 p-4 pt-0">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
