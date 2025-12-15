<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import type { MonitorType } from '../../../../types';
	import { Slider } from '$lib/components/ui/slider';

	let selectedMonitorType = $state<MonitorType>('http');

	type MonitorTypeSelect = {
		title: string;
		description: string;
		value: MonitorType;
	};

	const monitorTypeSelectData: readonly MonitorTypeSelect[] = [
		{
			title: 'HTTP',
			description: 'Monitor an HTTP or HTTPS endpoint and check response status and latency.',
			value: 'http'
		},
		{
			title: 'Ping',
			description: 'Monitor a host using ICMP ping to check network availability and latency.',
			value: 'ping'
		}
	];

	type IntervalOption = {
		label: string;
		seconds: number;
	};

	const intervalOptions: readonly IntervalOption[] = [
		{ label: '30s', seconds: 30 },
		{ label: '45s', seconds: 45 },
		{ label: '1m', seconds: 60 },
		{ label: '3m', seconds: 180 },
		{ label: '5m', seconds: 300 },
		{ label: '10m', seconds: 600 },
		{ label: '15m', seconds: 900 },
		{ label: '30m', seconds: 1800 },
		{ label: '1h', seconds: 3600 },
		{ label: '2h', seconds: 7200 }
	];
	let intervalIndex = $state<number>(3);
	const intervalSeconds = $derived(() => intervalOptions[intervalIndex].seconds);
</script>

<Card.Root class="mx-auto w-full">
	<Card.Content>
		<form class="space-y-4">
			<Field.Set>
				<Field.Label for="monitor-name">Monitor name</Field.Label>
				<Input id="monitor-name" name="monitor-name" type="text" placeholder="My Monitor" />
				<Field.Description>The name of your monitor.</Field.Description>

				<Field.Label>Monitor type</Field.Label>
				<Field.Description>Select the type of monitor you want to create.</Field.Description>
				<RadioGroup.Root
					bind:value={selectedMonitorType}
					class="grid gap-4 grid-cols-[repeat(auto-fit,minmax(240px,1fr))]"
				>
					{#each monitorTypeSelectData as option (option.value)}
						<Field.Label for={`monitor-type-${option.value}`}>
							<Field.Field orientation="horizontal">
								<Field.Content>
									<Field.Title>{option.title}</Field.Title>
									<Field.Description>
										{option.description}
									</Field.Description>
								</Field.Content>

								<RadioGroup.Item id={`monitor-type-${option.value}`} value={option.value} />
							</Field.Field>
						</Field.Label>
					{/each}
				</RadioGroup.Root>

				<Field.Label>Interval</Field.Label>
				<Field.Description>
					Monitor will run every
					<span class="font-medium">
						{intervalOptions[intervalIndex].label}
					</span>
				</Field.Description>

				<Slider
					type="single"
					min={0}
					max={intervalOptions.length - 1}
					step={1}
					bind:value={intervalIndex}
					class="mt-4"
				/>
				<div class="mt-2 flex justify-between text-xs text-muted-foreground">
					{#each intervalOptions as option (option.seconds)}
						<span class="w-6 text-center">
							{option.label}
						</span>
					{/each}
				</div>
			</Field.Set>
		</form>
	</Card.Content>
</Card.Root>
