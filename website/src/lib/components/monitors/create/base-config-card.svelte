<script lang="ts">
	import Icon from "@iconify/svelte";

	import { Badge } from "$lib/components/ui/badge/index.js";
	import { Button } from "$lib/components/ui/button/index.js";
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import * as Field from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { NativeSelect, NativeSelectOption } from "$lib/components/ui/native-select/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";

	import type { MonitorBaseSettings, MonitorKind } from "../../../../types/monitor-create";

	type TypeOption = {
		id: MonitorKind;
		label: string;
		description: string;
	};

	let {
		settings = $bindable<MonitorBaseSettings>({
			name: "",
			interval: 60,
			failure_threshold: 3,
			recovery_threshold: 1,
			notification: [],
		}),
		type,
		typeOptions = [] as TypeOption[],
		onTypeChange = (next: MonitorKind) => next,
	}: {
		settings: MonitorBaseSettings;
		type: MonitorKind;
		typeOptions?: TypeOption[];
		onTypeChange?: (next: MonitorKind) => void;
	} = $props();

	const intervalPresets = [
		{ label: "30s", value: 30 },
		{ label: "1m", value: 60 },
		{ label: "5m", value: 300 },
	];

	let notificationInput = $state(settings.notification.join(", "));

	$effect(() => {
		const joined = settings.notification.join(", ");
		if (joined !== notificationInput) notificationInput = joined;
	});

	$effect(() => {
		settings.notification = notificationInput
			.split(",")
			.map((id) => id.trim())
			.filter(Boolean);
	});

	let intervalInput = $derived(settings.interval.toString());
	$effect(() => {
		intervalInput = settings.interval.toString();
	});
	$effect(() => {
		const parsed = Number.parseInt(intervalInput, 10);
		settings.interval = Number.isFinite(parsed) && parsed > 0 ? parsed : 0;
	});

	let failureInput = $derived(settings.failure_threshold.toString());
	$effect(() => {
		failureInput = settings.failure_threshold.toString();
	});
	$effect(() => {
		const parsed = Number.parseInt(failureInput, 10);
		settings.failure_threshold = Number.isFinite(parsed) && parsed > 0 ? parsed : 1;
	});

	let recoveryInput = $derived(settings.recovery_threshold.toString());
	$effect(() => {
		recoveryInput = settings.recovery_threshold.toString();
	});
	$effect(() => {
		const parsed = Number.parseInt(recoveryInput, 10);
		settings.recovery_threshold = Number.isFinite(parsed) && parsed > 0 ? parsed : 1;
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
						bind:value={settings.name}
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
								class={`border-input hover:border-primary/60 hover:shadow-sm text-left transition-colors rounded-lg border p-3 ${
									type === option.id
										? "border-primary/80 bg-primary/5 ring-2 ring-primary/20"
										: "bg-card/50"
								}`}
								onclick={() => onTypeChange(option.id)}
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
							bind:value={intervalInput}
							aria-label="Interval in seconds"
						/>
						<div class="flex gap-1">
							{#each intervalPresets as preset (preset.value)}
								<Button
									type="button"
									variant="outline"
									size="sm"
									class="rounded-full"
									onclick={() => (intervalInput = preset.value.toString())}
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
						<NativeSelect aria-label="Failure threshold" bind:value={failureInput}>
							<NativeSelectOption value="1">After 1 failure</NativeSelectOption>
							<NativeSelectOption value="2">After 2 failures</NativeSelectOption>
							<NativeSelectOption value="3">After 3 failures</NativeSelectOption>
							<NativeSelectOption value="5">After 5 failures</NativeSelectOption>
						</NativeSelect>
						<Field.Description>Consecutive failures before opening an incident.</Field.Description>
					</div>
					<div class="space-y-1.5">
						<NativeSelect aria-label="Recovery threshold" bind:value={recoveryInput}>
							<NativeSelectOption value="1">After 1 success</NativeSelectOption>
							<NativeSelectOption value="2">After 2 successes</NativeSelectOption>
							<NativeSelectOption value="3">After 3 successes</NativeSelectOption>
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
					<Textarea
						bind:value={notificationInput}
						rows={2}
						placeholder="12345, 67890"
					/>
					<Field.Description>
						Comma-separated IDs matching the backend notification resources. Leave blank to wire later.
					</Field.Description>
				</Field.Content>
			</Field.Field>
		</Field.Group>
	</CardContent>
</Card>
