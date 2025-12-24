<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Card from '$lib/components/ui/card';
	import * as Select from '$lib/components/ui/select';
	import Icon from '@iconify/svelte';

	import type {
		MonitorWithIncidents,
		StatusPageGroup,
		StatusPageMonitor,
		StatusPageWithElements
	} from '$lib/types';

	import StatusPageMonitorCard from './monitor-raw-status-page.svelte';
	import StatusPageGroupCard from './group-raw-status-page.svelte';

	let monitorID = $state('');
	let groupID = $state('');

	let {
		monitors,
		statusPage
	}: { monitors: MonitorWithIncidents[]; statusPage: StatusPageWithElements } = $props();

	const monitorTriggerContent = $derived(
		monitors.find((m) => m.id === monitorID)?.name ?? 'Select a monitor'
	);

	const groupTriggerContent = $derived(
		statusPage.groups.find((g) => g.id === groupID)?.name ?? 'Ungrouped'
	);

	export type DndFatherElement =
		| (StatusPageGroup & { kind: 'group' })
		| (StatusPageMonitor & { kind: 'monitor' });

	const dndFatherElements = $derived.by((): DndFatherElement[] => {
		const groupItems: DndFatherElement[] = statusPage.groups.map((g) => ({
			...g,
			kind: 'group' as const
		}));

		const ungroupedMonitorItems: DndFatherElement[] = statusPage.monitors
			.filter((m) => m.groupId == null || m.groupId === '')
			.map((m) => ({ ...m, kind: 'monitor' as const }));

		return [...groupItems, ...ungroupedMonitorItems].sort((a, b) => {
			const d = a.sortOrder - b.sortOrder;
			if (d !== 0) return d;
			return a.id.localeCompare(b.id);
		});
	});

	function isMonitor(el: DndFatherElement): el is StatusPageMonitor & { kind: 'monitor' } {
		return el.kind === 'monitor';
	}

	function getMonitorsInGroup(groupId: string) {
		return statusPage.monitors
			.filter((m) => m.groupId === groupId)
			.slice()
			.sort((a, b) => a.sortOrder - b.sortOrder);
	}

	// TODO: 下面這些你自己接 API 或更新 state
	const onDeleteMonitor = (id: string) => {
		console.log('delete monitor element', id);
	};

	const onDeleteGroup = (id: string) => {
		console.log('delete group', id);
	};

	const onDeleteMonitorInGroup = (id: string) => {
		console.log('delete monitor in group', id);
	};

	const onAddElement = () => {
		// 你上面選了 monitorID + groupID 之後點 + 會新增 element
		console.log('add element', { monitorID, groupID });
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
						{#each monitors as monitor (monitor.id)}
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
							{#each statusPage.groups as group (group.id)}
								<Select.Item value={group.id} label={group.name}>
									{group.name}
								</Select.Item>
							{/each}
						</Select.Group>
					</Select.Content>
				</Select.Root>

				<Button size="icon" onclick={onAddElement}>
					<Icon icon="lucide:plus" />
				</Button>
			</div>
		</div>

		<!-- elements list -->
		<div class="flex flex-col gap-2">
			{#each dndFatherElements as fatherElement, i (fatherElement.id)}
				{#if isMonitor(fatherElement)}
					<Card.Root class="bg-muted p-0">
						<StatusPageMonitorCard monitor={fatherElement} index={i} onDelete={onDeleteMonitor} />
					</Card.Root>
				{:else}
					<StatusPageGroupCard
						group={fatherElement}
						index={i}
						monitors={getMonitorsInGroup(fatherElement.id)}
						{onDeleteGroup}
						onDeleteMonitor={onDeleteMonitorInGroup}
					/>
				{/if}
			{/each}
		</div>

		<div class="flex justify-end">
			<Button>
				<Icon icon="lucide:plus" />
				New Group
			</Button>
		</div>
	</Card.Content>
</Card.Root>
