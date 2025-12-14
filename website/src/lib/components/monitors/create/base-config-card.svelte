<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { NativeSelect, NativeSelectOption } from '$lib/components/ui/native-select/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	import type { SuperForm } from 'sveltekit-superforms';
	import type { MonitorCreate } from '$lib/schemas/monitor';
	import type { MonitorKind } from '../../../../types/monitor-create';

	type TypeOption = {
		id: MonitorKind;
		label: string;
		description: string;
	};

	let {
		form,
		type,
		typeOptions = [] as TypeOption[],
		handleTypeChange
	}: {
		form: SuperForm<MonitorCreate>;
		type: MonitorKind;
		typeOptions?: TypeOption[];
		handleTypeChange: (next: MonitorKind) => void;
	} = $props();

	const formData = $derived(form.form);

	const intervalPresets = [
		{ label: '30s', value: 30 },
		{ label: '1m', value: 60 },
		{ label: '3m', value: 180 },
		{ label: '5m', value: 300 },
	];

	let notificationInput = $state($formData.notification.join(', '));

	$effect(() => {
		const joined = $formData.notification.join(', ');
		if (joined !== notificationInput) notificationInput = joined;
	});

	$effect(() => {
		$formData.notification = notificationInput
			.split(',')
			.map((id) => id.trim())
			.filter(Boolean);
	});
</script>

<Card>
	<CardHeader>
		<CardTitle class="flex items-center gap-2 text-xl">
			<Icon icon="lucide:settings-2" class="size-5 text-primary" />
			Monitor basics
		</CardTitle>
	</CardHeader>
	<CardContent class="space-y-4">
		<Field.Group>
			<Field.Field>
				<Field.Label>Monitor name</Field.Label>
				<Field.Content class="space-y-2">
					<Input
						placeholder="API latency (US)"
						bind:value={$formData.name}
						autocomplete="off"
						required
					/>
					<Field.Description>Visible in lists, incidents, and notifications.</Field.Description>
				</Field.Content>
			</Field.Field>

			<Field.Field class="flex-col">
				<Field.Label>Monitor type</Field.Label>
				<Field.Content>
					<div class="grid gap-3 md:grid-cols-2">
						{#each typeOptions as option (option.id)}
							<button
								type="button"
								class={`cursor-pointer border-input hover:border-primary/60 hover:shadow-sm text-left transition-colors rounded-lg border p-3 ${
									type === option.id
										? 'border-primary/80 bg-primary/5 ring-2 ring-primary/20'
										: 'bg-card/50'
								}`}
								onclick={() => handleTypeChange(option.id)}
								aria-pressed={type === option.id}
							>
								<div class="flex items-start gap-2">
									<div class="space-y-1">
										<p class="text-sm font-semibold leading-tight text-foreground">
											{option.label}
										</p>
										<p class="text-xs text-muted-foreground">{option.description}</p>
									</div>
								</div>
							</button>
						{/each}
					</div>
				</Field.Content>
			</Field.Field>
		</Field.Group>

		<Field.Group>
			<Field.Field>
				<Field.Label>Check interval</Field.Label>
				<Field.Content class="space-y-2">
					<div class="flex flex-wrap gap-2">
						<Input
							class="max-w-32"
							type="number"
							min="10"
							step="5"
							inputmode="numeric"
							bind:value={$formData.interval}
							aria-label="Interval in seconds"
						/>
						<div class="flex gap-1">
							{#each intervalPresets as preset (preset.value)}
								<Button
									type="button"
									variant="outline"
									size="sm"
									class="rounded-full"
									onclick={() => ($formData.interval = preset.value)}
								>
									{preset.label}
								</Button>
							{/each}
						</div>
					</div>
					<Field.Description>
						How often we probe the target. Uses seconds; 60s is a good default for uptime checks.
					</Field.Description>
				</Field.Content>
			</Field.Field>

			<Field.Field>
				<Field.Label>Failure &amp; recovery thresholds</Field.Label>
				<Field.Content class="grid gap-3">
					<div class="space-y-1.5">
						<NativeSelect aria-label="Failure threshold" bind:value={$formData.failure_threshold}>
							<NativeSelectOption value={1}>After 1 failure</NativeSelectOption>
							<NativeSelectOption value={2}>After 2 failures</NativeSelectOption>
							<NativeSelectOption value={3}>After 3 failures</NativeSelectOption>
							<NativeSelectOption value={5}>After 5 failures</NativeSelectOption>
						</NativeSelect>
						<Field.Description>Consecutive failures before opening an incident.</Field.Description>
					</div>
					<div class="space-y-1.5">
						<NativeSelect aria-label="Recovery threshold" bind:value={$formData.recovery_threshold}>
							<NativeSelectOption value={1}>After 1 success</NativeSelectOption>
							<NativeSelectOption value={2}>After 2 successes</NativeSelectOption>
							<NativeSelectOption value={3}>After 3 successes</NativeSelectOption>
						</NativeSelect>
						<Field.Description>Consecutive successes to auto-resolve.</Field.Description>
					</div>
				</Field.Content>
			</Field.Field>
		</Field.Group>

		<Field.Group>
			<Field.Field>
				<Field.Label>Notification channel IDs</Field.Label>
				<Field.Content class="space-y-2">
					<Textarea bind:value={notificationInput} rows={2} placeholder="12345, 67890" />
					<Field.Description>
						Comma-separated IDs matching the backend notification resources. Leave blank to wire
						later.
					</Field.Description>
				</Field.Content>
			</Field.Field>
		</Field.Group>
	</CardContent>
</Card>
