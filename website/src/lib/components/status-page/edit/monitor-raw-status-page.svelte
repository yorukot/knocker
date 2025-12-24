<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import * as InputGroup from '$lib/components/ui/input-group';
	import * as Select from '$lib/components/ui/select';
	import Icon from '@iconify/svelte';
	import type { StatusPageMonitor, StatusPageElementType } from '$lib/types';

	const {
		monitor,
		index,
		onDelete
	}: {
		monitor: StatusPageMonitor;
		index: number;
		onDelete?: (monitorId: string) => void;
	} = $props();

	const typeLabel = (t: StatusPageElementType) =>
		t === 'historical_timeline' ? 'Historical Timeline' : 'Only Current Status';

	const typeIcon = (t: StatusPageElementType) =>
		t === 'historical_timeline' ? 'lucide:chart-line' : 'lucide:circle-small';
</script>

<div class="flex justify-between items-center p-2 gap-2">
	<InputGroup.Root class="w-full">
		<InputGroup.Input value={monitor.name} placeholder="Please enter element name" />
		<InputGroup.Addon class="hidden sm:block">
			<Icon icon="lucide:activity" />
		</InputGroup.Addon>
	</InputGroup.Root>

	<div class="flex items-center gap-2">
		<Select.Root type="single" name={`monitors.${index}.type`}  value={monitor.type}>
			<Select.Trigger class="lg:w-51">
				<Icon icon={typeIcon(monitor.type)} />
				<p class="hidden lg:block">{typeLabel(monitor.type)}</p>
			</Select.Trigger>
			<Select.Content>
				<Select.Group>
					<Select.Item value="historical_timeline" label="Historical Timeline">
						<Icon icon="lucide:chart-line" /> Historical Timeline
					</Select.Item>
					<Select.Item value="current_status_indicator" label="Only Current Status">
						<Icon icon="lucide:circle-small" /> Only Current Status
					</Select.Item>
				</Select.Group>
			</Select.Content>
		</Select.Root>

		<Button size="icon" variant="destructive" onclick={() => onDelete?.(monitor.id)}>
			<Icon icon="lucide:trash" />
		</Button>
	</div>
</div>
