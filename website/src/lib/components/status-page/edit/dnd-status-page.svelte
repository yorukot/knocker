<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import Icon from '@iconify/svelte';
	import type { createForm } from 'felte';

	import type { MonitorWithIncidents, StatusPageElement, StatusPageMonitor } from '$lib/types';

	import StatusPageMonitorCard from './monitor-raw-status-page.svelte';
	import StatusPageGroupCard from './group-raw-status-page.svelte';
	import type { StatusPageUpsertValues } from './schema';

	type FelteReturn = ReturnType<typeof createForm<StatusPageUpsertValues>>;
	type SetFields = FelteReturn['setFields'];

	let monitorID = $state('');
	let groupID = $state('');

	let {
		availableMonitors,
		statusPageId,
		elements,
		setFields
	}: {
		availableMonitors: MonitorWithIncidents[];
		statusPageId: string;
		elements: StatusPageElement[];
		setFields: SetFields;
	} = $props();

	const sortByOrder = (
		a: { sortOrder: number; id: string },
		b: { sortOrder: number; id: string }
	) => {
		const diff = a.sortOrder - b.sortOrder;
		if (diff !== 0) return diff;
		return a.id.localeCompare(b.id);
	};

	function normalizeElements(list: StatusPageElement[]) {
		return list
			.map((element) => ({
				...element,
				monitors: (element.monitors ?? []).slice().sort(sortByOrder)
			}))
			.sort(sortByOrder);
	}

	let editableElements = $derived(normalizeElements(structuredClone(elements ?? [])));

	const monitorTriggerContent = $derived(
		availableMonitors.find((m) => m.id === monitorID)?.name ?? 'Select a monitor'
	);

	const groupTriggerContent = $derived(
		editableElements.find((g) => !g.monitor && g.id === groupID)?.name ?? 'Ungrouped'
	);

	function updateElements(next: StatusPageElement[]) {
		const normalized = normalizeElements(next);
		editableElements = normalized;
		setFields('elements', normalized, true);
	}

	const onDeleteMonitor = (id: string) => {
		console.log('delete monitor element', id);
	};

	const onDeleteGroup = (id: string) => {
		console.log('delete group', id);
	};

	const onDeleteMonitorInGroup = (id: string) => {
		console.log('delete monitor in group', id);
	};

	function nextTopLevelSortOrder() {
		const max = editableElements.reduce((acc, element) => Math.max(acc, element.sortOrder ?? 0), 0);
		return max + 1;
	}

	function nextSortOrderInGroup(groupId: string) {
		const group = editableElements.find((element) => !element.monitor && element.id === groupId);
		const list = group?.monitors ?? [];
		const max = list.reduce((acc, m) => Math.max(acc, m.sortOrder ?? 0), 0);
		return max + 1;
	}

	const onAddMonitor = () => {
		if (!monitorID) return;

		const targetGroupId = groupID === '0' || groupID === '' ? null : groupID;

		if (!targetGroupId) {
			const newElement: StatusPageElement = {
				id: crypto.randomUUID(),
				statusPageId,
				name: availableMonitors.find((m) => m.id === monitorID)?.name ?? 'New Monitor',
				type: 'historical_timeline',
				sortOrder: nextTopLevelSortOrder(),
				monitor: true,
				monitorId: monitorID,
				monitors: []
			};

			updateElements([...editableElements, newElement]);
		} else {
			updateElements(
				editableElements.map((element) => {
					if (element.monitor || element.id !== targetGroupId) return element;
					const nextMonitor: StatusPageMonitor = {
						id: crypto.randomUUID(),
						statusPageId,
						monitorId: monitorID,
						groupId: targetGroupId,
						name: availableMonitors.find((m) => m.id === monitorID)?.name ?? 'New Monitor',
						type: 'historical_timeline',
						sortOrder: nextSortOrderInGroup(targetGroupId)
					};
					return {
						...element,
						monitors: [...(element.monitors ?? []), nextMonitor]
					};
				})
			);
		}

		monitorID = '';
		groupID = '';
	};

	const onAddGroup = () => {
		const newGroup: StatusPageElement = {
			id: crypto.randomUUID(),
			statusPageId,
			name: 'New Group',
			type: 'historical_timeline',
			sortOrder: nextTopLevelSortOrder(),
			monitor: false,
			monitorId: null,
			monitors: []
		};

		updateElements([...editableElements, newGroup]);
	};
	
	
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Status page elements</Card.Title>
		<Card.Description>Configure the status page's elements.</Card.Description>
	</Card.Header>

	<Card.Content class="flex flex-col gap-4">
		<!-- top selectors -->
		<div class="flex gap-2 flex-col sm:flex-row warp">
			<Select.Root type="single" name="monitor" bind:value={monitorID}>
				<Select.Trigger class="max-w-lg w-full">{monitorTriggerContent}</Select.Trigger>
				<Select.Content>
					<Select.Group>
						<Select.Label>Monitors</Select.Label>
						{#each availableMonitors as monitor (monitor.id)}
							<Select.Item value={monitor.id} label={monitor.name}>
								{monitor.name}
							</Select.Item>
						{/each}
					</Select.Group>
				</Select.Content>
			</Select.Root>

			<div class="flex gap-2 w-full">
				<Select.Root type="single" name="group" bind:value={groupID}>
					<Select.Trigger class="max-w-lg w-full truncate">{groupTriggerContent}</Select.Trigger>
					<Select.Content>
						<Select.Group>
							<Select.Label>Groups</Select.Label>
							<Select.Item value="0" label="Ungrouped">Ungrouped</Select.Item>
							{#each editableElements.filter((element) => !element.monitor) as group (group.id)}
								<Select.Item value={group.id} label={group.name}>
									{group.name}
								</Select.Item>
							{/each}
						</Select.Group>
					</Select.Content>
				</Select.Root>

				<Button size="icon" onclick={onAddMonitor}>
					<Icon icon="lucide:plus" />
				</Button>
			</div>
		</div>

		<!-- elements list -->
		<div class="flex flex-col gap-2">
			{#each editableElements as element, i (element.id)}
				{#if element.monitor}
					<Card.Root class="bg-muted p-0">
						<StatusPageMonitorCard
							monitor={element}
							namePrefix={`elements.${i}`}
							isElement
							onDelete={onDeleteMonitor}
						/>
					</Card.Root>
				{:else}
					<StatusPageGroupCard
						group={element}
						namePrefix={`elements.${i}`}
						{onDeleteGroup}
						onDeleteMonitor={onDeleteMonitorInGroup}
					/>
				{/if}
			{/each}
		</div>

			<div class="flex justify-end">
				<Button onclick={onAddGroup}>
					<Icon icon="lucide:plus" />
					New Group
				</Button>
		</div>
	</Card.Content>
</Card.Root>
