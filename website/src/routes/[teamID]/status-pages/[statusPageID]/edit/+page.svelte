<script lang="ts">
	import {
		dndzone,
		dragHandle,
		dragHandleZone,
		SHADOW_PLACEHOLDER_ITEM_ID,
		type DndEvent
	} from 'svelte-dnd-action';
	import { flip } from 'svelte/animate';
	import Icon from '@iconify/svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { toast } from 'svelte-sonner';

	import * as Card from '$lib/components/ui/card';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Field from '$lib/components/ui/field';
	import * as Select from '$lib/components/ui/select';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { deleteStatusPage, updateStatusPage } from '$lib/api/status-page';
	import type { StatusPageElementType } from '../../../../types';
	import type { StatusPageWithElements } from '../../../../types/status-page';

	const flipDurationMs = 150;

	type GroupItem = {
		id: string;
		name: string;
		type: StatusPageElementType;
	};

	type MonitorItem = {
		id: string;
		monitorId: string;
		name: string;
		type: StatusPageElementType;
		groupId: string | null;
	};

	const typeOptions = [
		{ value: 'current_status_indicator', label: 'Current status indicator' },
		{ value: 'historical_timeline', label: 'Historical timeline' }
	] as const;

	/** @type {import('./$types').PageProps} */
	let { data } = $props();

	let title = $state(data.statusPage.statusPage.title);
	let slug = $state(data.statusPage.statusPage.slug);

	let groups = $state<GroupItem[]>(normalizeGroups(data.statusPage));
	let ungroupedMonitors = $state<MonitorItem[]>(normalizeUngroupedMonitors(data.statusPage));
	let monitorsByGroup = $state<Record<string, MonitorItem[]>>(normalizeGroupedMonitors(data.statusPage));

	let selectedMonitorId = $state<string>('');
	let selectedGroupId = $state<string>('');
	let isSaving = $state(false);
	let deleteOpen = $state(false);
	let deleteLoading = $state(false);

	let tempIdCounter = $state(-1);

	function normalizeGroups(pageData: StatusPageWithElements): GroupItem[] {
		return [...pageData.groups]
			.sort((a, b) => a.sortOrder - b.sortOrder)
			.map((group) => ({
				id: group.id,
				name: group.name,
				type: group.type
			}));
	}

	function normalizeUngroupedMonitors(pageData: StatusPageWithElements): MonitorItem[] {
		return pageData.monitors
			.filter((monitor) => !monitor.groupId)
			.sort((a, b) => a.sortOrder - b.sortOrder)
			.map((monitor) => ({
				id: monitor.id,
				monitorId: monitor.monitorId,
				name: monitor.name,
				type: monitor.type,
				groupId: null
			}));
	}

	function normalizeGroupedMonitors(pageData: StatusPageWithElements): Record<string, MonitorItem[]> {
		const grouped: Record<string, MonitorItem[]> = {};
		for (const group of pageData.groups) {
			grouped[group.id] = [];
		}

		pageData.monitors
			.filter((monitor) => monitor.groupId)
			.sort((a, b) => a.sortOrder - b.sortOrder)
			.forEach((monitor) => {
				const groupId = monitor.groupId as string;
				if (!grouped[groupId]) {
					grouped[groupId] = [];
				}
				grouped[groupId].push({
					id: monitor.id,
					monitorId: monitor.monitorId,
					name: monitor.name,
					type: monitor.type,
					groupId
				});
			});

		return grouped;
	}

	function nextTempId(): string {
		const id = tempIdCounter;
		tempIdCounter -= 1;
		return String(id);
	}

	function handleGroupSort(event: CustomEvent<DndEvent<GroupItem>>) {
		groups = event.detail.items.filter((item) => item.id !== SHADOW_PLACEHOLDER_ITEM_ID);
	}

	function handleUngroupedSort(event: CustomEvent<DndEvent<MonitorItem>>) {
		ungroupedMonitors = event.detail.items
			.filter((item) => item.id !== SHADOW_PLACEHOLDER_ITEM_ID)
			.map((item) => ({ ...item, groupId: null }));
	}

	function handleGroupMonitorSort(groupId: string, event: CustomEvent<DndEvent<MonitorItem>>) {
		monitorsByGroup = {
			...monitorsByGroup,
			[groupId]: event.detail.items
				.filter((item) => item.id !== SHADOW_PLACEHOLDER_ITEM_ID)
				.map((item) => ({ ...item, groupId }))
		};
	}

	function addGroup() {
		const id = nextTempId();
		groups = [
			...groups,
			{
				id,
				name: 'New group',
				type: 'current_status_indicator'
			}
		];
		monitorsByGroup = {
			...monitorsByGroup,
			[id]: monitorsByGroup[id] ?? []
		};
	}

	function removeGroup(groupId: string) {
		const groupMonitors = monitorsByGroup[groupId] || [];
		ungroupedMonitors = [...ungroupedMonitors, ...groupMonitors.map((item) => ({ ...item, groupId: null }))];

		const updatedGroups = groups.filter((group) => group.id !== groupId);
		groups = updatedGroups;

		const updatedMonitors = { ...monitorsByGroup };
		delete updatedMonitors[groupId];
		monitorsByGroup = updatedMonitors;
	}

	function addMonitor() {
		if (!selectedMonitorId) return;
		const monitor = data.monitors.find((item) => item.id === selectedMonitorId);
		if (!monitor) return;

		const alreadyAdded =
			ungroupedMonitors.some((item) => item.monitorId === monitor.id) ||
			Object.values(monitorsByGroup).some((items) =>
				items.some((item) => item.monitorId === monitor.id)
			);

		if (alreadyAdded) {
			toast.error('Monitor already added to this status page');
			return;
		}

		const newItem: MonitorItem = {
			id: nextTempId(),
			monitorId: monitor.id,
			name: monitor.name,
			type: 'current_status_indicator',
			groupId: selectedGroupId || null
		};

		if (selectedGroupId) {
			const existing = monitorsByGroup[selectedGroupId] || [];
			monitorsByGroup = { ...monitorsByGroup, [selectedGroupId]: [...existing, newItem] };
		} else {
			ungroupedMonitors = [...ungroupedMonitors, newItem];
		}

		selectedMonitorId = '';
	}

	function removeMonitor(target: MonitorItem) {
		if (target.groupId) {
			const updated = (monitorsByGroup[target.groupId] || []).filter((item) => item.id !== target.id);
			monitorsByGroup = { ...monitorsByGroup, [target.groupId]: updated };
			return;
		}

		ungroupedMonitors = ungroupedMonitors.filter((item) => item.id !== target.id);
	}

	function buildGroupsPayload() {
		return groups.map((group, index) => ({
			id: group.id,
			name: group.name,
			type: group.type,
			sort_order: index + 1
		}));
	}

	function buildMonitorsPayload() {
		const payload: {
			id: string;
			monitor_id: string;
			group_id?: string | null;
			name: string;
			type: StatusPageElementType;
			sort_order: number;
		}[] = [];

		ungroupedMonitors.forEach((monitor, index) => {
			payload.push({
				id: monitor.id,
				monitor_id: monitor.monitorId,
				name: monitor.name,
				type: monitor.type,
				sort_order: index + 1
			});
		});

		Object.entries(monitorsByGroup).forEach(([groupId, items]) => {
			items.forEach((monitor, index) => {
				payload.push({
					id: monitor.id,
					monitor_id: monitor.monitorId,
					group_id: groupId,
					name: monitor.name,
					type: monitor.type,
					sort_order: index + 1
				});
			});
		});

		return payload;
	}

	async function handleSave() {
		const teamID = page.params.teamID;
		const statusPageID = page.params.statusPageID;
		if (!teamID || !statusPageID) {
			toast.error('Missing team or status page id');
			return;
		}

		isSaving = true;
		try {
			await updateStatusPage(teamID, statusPageID, {
				title,
				slug,
				icon: null,
				groups: buildGroupsPayload(),
				monitors: buildMonitorsPayload()
			});
			toast.success('Status page updated');
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Failed to update status page';
			toast.error(message);
		} finally {
			isSaving = false;
		}
	}

	async function handleDelete() {
		const teamID = page.params.teamID;
		const statusPageID = page.params.statusPageID;
		if (!teamID || !statusPageID) {
			toast.error('Missing team or status page id');
			return;
		}

		deleteLoading = true;
		try {
			await deleteStatusPage(teamID, statusPageID);
			toast.success('Status page deleted');
			goto(`/${teamID}/status-pages`);
		} catch (error) {
			const message = error instanceof Error ? error.message : 'Failed to delete status page';
			toast.error(message);
		} finally {
			deleteLoading = false;
			deleteOpen = false;
		}
	}

	function groupOptionLabel(groupId: string) {
		const group = groups.find((item) => item.id === groupId);
		return group ? group.name : 'Ungrouped';
	}
</script>

<div class="flex flex-col gap-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold">Edit status page</h1>
			<p class="text-sm text-muted-foreground">Organize what your public status page shows.</p>
		</div>
		<div class="flex items-center gap-2">
			<Button variant="destructive" onclick={() => (deleteOpen = true)}>
				Delete
			</Button>
			<Button onclick={handleSave} disabled={isSaving || !title.trim() || !slug.trim()}>
				{isSaving ? 'Saving…' : 'Save changes'}
			</Button>
		</div>
	</div>

	<Card.Root class="p-6">
		<Card.Header class="p-0">
			<Card.Title class="text-lg">Basics</Card.Title>
			<Card.Description>Update the title and slug for the status page.</Card.Description>
		</Card.Header>
		<Card.Content class="p-0">
			<div class="grid gap-4 md:grid-cols-2">
				<Field.Field>
					<Field.Label for="status-page-title">Title</Field.Label>
					<Input id="status-page-title" bind:value={title} />
				</Field.Field>
				<Field.Field>
					<Field.Label for="status-page-slug">Slug</Field.Label>
					<Input id="status-page-slug" bind:value={slug} />
				</Field.Field>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root class="p-6">
		<Card.Header class="p-0">
			<Card.Title class="text-lg">Add monitors</Card.Title>
			<Card.Description>Select monitors to show on your status page.</Card.Description>
		</Card.Header>
		<Card.Content class="p-0">
			<div class="flex flex-col gap-3 md:flex-row md:items-end">
				<div class="flex-1 space-y-2">
					<Field.Label>Monitor</Field.Label>
					<Select.Root type="single" bind:value={selectedMonitorId}>
						<Select.Trigger class="justify-between w-full">
							<span data-slot="select-value" class="text-sm font-medium">
								{selectedMonitorId
									? data.monitors.find((item) => item.id === selectedMonitorId)?.name
									: 'Choose a monitor'}
							</span>
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								{#each data.monitors as monitor (monitor.id)}
									<Select.Item value={monitor.id}>{monitor.name}</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
				<div class="flex-1 space-y-2">
					<Field.Label>Group</Field.Label>
					<Select.Root type="single" bind:value={selectedGroupId}>
						<Select.Trigger class="justify-between w-full">
							<span data-slot="select-value" class="text-sm font-medium">
								{selectedGroupId ? groupOptionLabel(selectedGroupId) : 'Ungrouped'}
							</span>
						</Select.Trigger>
						<Select.Content>
							<Select.Group>
								<Select.Item value="">Ungrouped</Select.Item>
								{#each groups as group (group.id)}
									<Select.Item value={group.id}>{group.name}</Select.Item>
								{/each}
							</Select.Group>
						</Select.Content>
					</Select.Root>
				</div>
				<Button onclick={addMonitor} disabled={!selectedMonitorId}>Add monitor</Button>
			</div>
		</Card.Content>
	</Card.Root>

	<Card.Root class="p-6">
		<Card.Header class="p-0">
			<Card.Title class="text-lg">Page elements</Card.Title>
			<Card.Description>Drag groups and monitors to arrange the status page.</Card.Description>
		</Card.Header>
		<Card.Content class="p-0">
			<div class="flex items-center justify-between">
				<p class="text-sm text-muted-foreground">Drag handles to reorder groups.</p>
				<Button variant="secondary" size="sm" onclick={addGroup}>
					<Icon icon="lucide:plus" /> Add group
				</Button>
			</div>

			<div
				class="mt-4 flex flex-col gap-3"
				use:dragHandleZone={{ items: groups, flipDurationMs }}
				onconsider={handleGroupSort}
				onfinalize={handleGroupSort}
			>
				<div class="rounded-md border border-dashed p-4 bg-background flex flex-col gap-2">
					<div class="flex items-center gap-2 text-sm text-muted-foreground">
						<Badge variant="secondary" class="rounded-sm">
							{ungroupedMonitors.length} monitors
						</Badge>
						<span>Ungrouped monitors</span>
					</div>
					<div
						class="flex flex-col gap-2"
						use:dndzone={{
							items: ungroupedMonitors,
							flipDurationMs,
							type: 'monitor'
						}}
						onconsider={handleUngroupedSort}
						onfinalize={handleUngroupedSort}
					>
						{#each ungroupedMonitors as monitor (monitor.id)}
							<div
								class="rounded-md border px-3 py-2 flex flex-col gap-2 bg-background"
								animate:flip={{ duration: flipDurationMs }}
							>
								<div class="flex items-center gap-2">
									<span class="text-sm font-medium flex-1">{monitor.name}</span>
									<Button variant="ghost" size="icon" onclick={() => removeMonitor(monitor)}>
										<Icon icon="lucide:x" />
									</Button>
								</div>
								<div class="grid gap-2 md:grid-cols-3">
									<div class="space-y-1 md:col-span-2">
										<Field.Label>Display name</Field.Label>
										<Input bind:value={monitor.name} />
									</div>
									<div class="space-y-1">
										<Field.Label>Type</Field.Label>
										<Select.Root type="single" bind:value={monitor.type}>
											<Select.Trigger class="justify-between w-full">
												<span data-slot="select-value" class="text-sm font-medium">
													{typeOptions.find((option) => option.value === monitor.type)?.label}
												</span>
											</Select.Trigger>
											<Select.Content>
												<Select.Group>
													{#each typeOptions as option (option.value)}
														<Select.Item value={option.value}>{option.label}</Select.Item>
													{/each}
												</Select.Group>
											</Select.Content>
										</Select.Root>
									</div>
								</div>
							</div>
						{/each}
						{#if !ungroupedMonitors.length}
							<div class="rounded-md border border-dashed p-3 text-sm text-muted-foreground">
								Drop monitors here
							</div>
						{/if}
					</div>
				</div>

				{#each groups as group (group.id)}
					<div
						class="rounded-md border p-4 bg-background flex flex-col gap-3"
						animate:flip={{ duration: flipDurationMs }}
					>
						<div class="flex items-center gap-2">
							<button
								class="text-muted-foreground hover:text-foreground"
								type="button"
								use:dragHandle
								aria-label="Drag group"
							>
								<Icon icon="lucide:grip-vertical" />
							</button>
							<div class="flex-1 grid gap-2 md:grid-cols-2">
								<div class="space-y-1">
									<Field.Label>Group name</Field.Label>
									<Input bind:value={group.name} />
								</div>
								<div class="space-y-1">
									<Field.Label>Type</Field.Label>
									<Select.Root type="single" bind:value={group.type}>
										<Select.Trigger class="justify-between w-full">
											<span data-slot="select-value" class="text-sm font-medium">
												{typeOptions.find((option) => option.value === group.type)?.label}
											</span>
										</Select.Trigger>
										<Select.Content>
											<Select.Group>
												{#each typeOptions as option (option.value)}
													<Select.Item value={option.value}>{option.label}</Select.Item>
												{/each}
											</Select.Group>
										</Select.Content>
									</Select.Root>
								</div>
							</div>
							<Button variant="ghost" size="icon" onclick={() => removeGroup(group.id)}>
								<Icon icon="lucide:trash" />
							</Button>
						</div>

						<div class="space-y-2">
							<div class="flex items-center gap-2 text-sm text-muted-foreground">
								<Badge variant="secondary" class="rounded-sm">
									{(monitorsByGroup[group.id] || []).length} monitors
								</Badge>
								<span>Drag monitors into this group.</span>
							</div>
							<div
								class="flex flex-col gap-2"
								use:dndzone={{
									items: monitorsByGroup[group.id] || [],
									flipDurationMs,
									type: 'monitor'
								}}
								onconsider={(event) => handleGroupMonitorSort(group.id, event)}
								onfinalize={(event) => handleGroupMonitorSort(group.id, event)}
							>
								{#each monitorsByGroup[group.id] || [] as monitor (monitor.id)}
									<div
										class="rounded-md border px-3 py-2 flex flex-col gap-2 bg-background"
										animate:flip={{ duration: flipDurationMs }}
									>
										<div class="flex items-center gap-2">
											<span class="text-sm font-medium flex-1">{monitor.name}</span>
											<Button variant="ghost" size="icon" onclick={() => removeMonitor(monitor)}>
												<Icon icon="lucide:x" />
											</Button>
										</div>
										<div class="grid gap-2 md:grid-cols-2">
											<div class="space-y-1">
												<Field.Label>Display name</Field.Label>
												<Input bind:value={monitor.name} />
											</div>
											<div class="space-y-1">
												<Field.Label>Type</Field.Label>
												<Select.Root type="single" bind:value={monitor.type}>
													<Select.Trigger class="justify-between w-full">
														<span data-slot="select-value" class="text-sm font-medium">
															{typeOptions.find((option) => option.value === monitor.type)?.label}
														</span>
													</Select.Trigger>
													<Select.Content>
														<Select.Group>
															{#each typeOptions as option (option.value)}
																<Select.Item value={option.value}>{option.label}</Select.Item>
															{/each}
														</Select.Group>
													</Select.Content>
												</Select.Root>
											</div>
										</div>
									</div>
								{/each}
								{#if !(monitorsByGroup[group.id] || []).length}
									<div class="rounded-md border border-dashed p-3 text-sm text-muted-foreground">
										Drop monitors here
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/each}
				{#if !groups.length}
					<div class="rounded-md border border-dashed p-4 text-sm text-muted-foreground">
						No groups yet. Add one to organize your monitors.
					</div>
				{/if}
			</div>
		</Card.Content>
	</Card.Root>
</div>

<AlertDialog.Root bind:open={deleteOpen}>
	<AlertDialog.Content class="max-w-md">
		<AlertDialog.Header>
			<AlertDialog.Title>Delete status page</AlertDialog.Title>
			<AlertDialog.Description>
				This will remove the status page and its groups and monitors.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel disabled={deleteLoading}>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action onclick={handleDelete} disabled={deleteLoading}>
				{deleteLoading ? 'Deleting…' : 'Delete'}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>
