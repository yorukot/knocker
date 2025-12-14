<script lang="ts">
	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button/index.js';
	import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { NativeSelect, NativeSelectOption } from '$lib/components/ui/native-select/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';

	import type { SuperForm } from 'sveltekit-superforms';
	import type { MonitorCreate } from '$lib/schemas/monitor';
	import type { MonitorKind } from '../../../../types/monitor-create';
	import type { ApiNotification } from '$lib/api/notification';

	type TypeOption = {
		id: MonitorKind;
		label: string;
		description: string;
	};

	let {
		form,
		type,
		typeOptions = [] as TypeOption[],
		handleTypeChange,
		notifications = []
	}: {
		form: SuperForm<MonitorCreate>;
		type: MonitorKind;
		typeOptions?: TypeOption[];
		handleTypeChange: (next: MonitorKind) => void;
		notifications?: ApiNotification[];
	} = $props();

	const formData = $derived(form.form);

	const intervalPresets = [
		{ label: '30s', value: 30 },
		{ label: '1m', value: 60 },
		{ label: '3m', value: 180 },
		{ label: '5m', value: 300 }
	];

	let notificationPopoverOpen = $state(false);

	const selectedNotifications = $derived(
		notifications.filter((n) => $formData.notification.includes(n.id))
	);

	function toggleNotification(notificationId: string) {
		if ($formData.notification.includes(notificationId)) {
			$formData.notification = $formData.notification.filter((id) => id !== notificationId);
		} else {
			$formData.notification = [...$formData.notification, notificationId];
		}
	}

	function removeNotification(notificationId: string) {
		$formData.notification = $formData.notification.filter((id) => id !== notificationId);
	}

	function getNotificationIcon(type: string) {
		switch (type) {
			case 'discord':
				return 'ri:discord-fill';
			case 'telegram':
				return 'ri:telegram-fill';
			default:
				return 'lucide:bell';
		}
	}
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
				<Field.Label>Notification channels</Field.Label>
				<Field.Content class="space-y-2">
					<Popover.Root bind:open={notificationPopoverOpen}>
						<Popover.Trigger
							class="flex h-9 w-full items-center justify-between gap-2 whitespace-nowrap rounded-md border border-input bg-background px-4 py-2 text-sm shadow-xs outline-none transition-all hover:bg-accent hover:text-accent-foreground focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 dark:border-input dark:bg-input/30 dark:hover:bg-input/50"
							role="combobox"
							aria-expanded={notificationPopoverOpen}
						>
							<span class="truncate">
								{selectedNotifications.length > 0
									? `${selectedNotifications.length} selected`
									: 'Select channels...'}
							</span>
							<Icon icon="lucide:chevrons-up-down" class="ml-2 size-4 shrink-0 opacity-50" />
						</Popover.Trigger>
						<Popover.Content class="w-[--bits-popover-trigger-width] p-2" align="start">
							{#if notifications.length === 0}
								<div class="py-6 text-center text-sm">
									<Icon icon="lucide:bell-off" class="mx-auto mb-2 size-8 text-muted-foreground" />
									<p class="text-muted-foreground">No notification channels found</p>
									<p class="mt-1 text-xs text-muted-foreground">
										Create one in your team settings first
									</p>
								</div>
							{:else}
								<div class="max-h-[300px] space-y-1 overflow-y-auto">
									{#each notifications as notification (notification.id)}
										<button
											type="button"
											class="flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-2 text-sm transition-colors hover:bg-accent"
											onclick={() => toggleNotification(notification.id)}
										>
											<div
												class={`flex size-4 shrink-0 items-center justify-center rounded-sm border border-primary ${
													$formData.notification.includes(notification.id)
														? 'bg-primary text-primary-foreground'
														: 'opacity-50 [&_svg]:invisible'
												}`}
											>
												<Icon icon="lucide:check" class="size-3" />
											</div>
											<Icon
												icon={getNotificationIcon(notification.type)}
												class="size-4 shrink-0 text-muted-foreground"
											/>
											<span class="flex-1 truncate text-left">{notification.name}</span>
											<Badge variant="secondary" class="text-xs capitalize">
												{notification.type}
											</Badge>
										</button>
									{/each}
								</div>
							{/if}
						</Popover.Content>
					</Popover.Root>

					{#if selectedNotifications.length > 0}
						<div class="flex flex-wrap gap-2">
							{#each selectedNotifications as notification (notification.id)}
								<Badge variant="secondary" class="gap-1.5 pr-1">
									<Icon icon={getNotificationIcon(notification.type)} class="size-3" />
									<span class="truncate">{notification.name}</span>
									<button
										type="button"
										class="ml-1 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
										onclick={() => removeNotification(notification.id)}
										aria-label={`Remove ${notification.name}`}
									>
										<Icon icon="lucide:x" class="size-3" />
									</button>
								</Badge>
							{/each}
						</div>
					{/if}

					<Field.Description>
						Select notification channels to alert when this monitor fails or recovers. You can add
						multiple channels.
					</Field.Description>
				</Field.Content>
			</Field.Field>
		</Field.Group>
	</CardContent>
</Card>
