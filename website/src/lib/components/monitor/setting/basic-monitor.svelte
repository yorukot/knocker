<script lang="ts">
	import { createForm } from 'felte';
	import { validator } from '@felte/validator-zod';
	import { z } from 'zod';

	import * as Card from '$lib/components/ui/card/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as RadioGroup from '$lib/components/ui/radio-group';
	import MultiSelect, { type MultiSelectOption } from '$lib/components/ui/multi-select';
	import type { MonitorType } from '../../../../types';
	import type { Notification } from '../../../../types';
	import { Slider } from '$lib/components/ui/slider';
	import { decidedNotificationIcon } from '../utils';
	import * as Select from '$lib/components/ui/select';
	import HttpMonitor from './http-monitor.svelte';
	import PingMonitor from './ping-monitor.svelte';
	import { Button } from '$lib/components/ui/button';
	import {
		intervalOptions,
		monitorTypeSelectData,
		thresholdOptions,
		httpMethods
	} from './setting';

	const httpConfigSchema = z.object({
		url: z.url('Must be a valid URL'),
		method: z.enum(httpMethods),
		maxRedirects: z.coerce.number().int().min(0).max(1000),
		requestTimeoutSeconds: z.coerce.number().int().min(0).max(120),
		headers: z.string().optional().default(''),
		bodyEncoding: z.enum(['json', 'xml', '']).default(''),
		body: z.string().optional().default(''),
		acceptedStatusCodes: z.array(z.string()).min(1, 'Select at least one status'),
		upsideDownMode: z.coerce.boolean(),
		certificateExpiryNotification: z.coerce.boolean(),
		ignoreTlsError: z.coerce.boolean()
	});

	const pingConfigSchema = z.object({
		host: z.string().min(1, 'Host is required'),
		timeoutSeconds: z.coerce.number().int().min(0).max(120),
		packetSize: z
			.union([z.coerce.number().int().min(1).max(65000), z.literal('')])
			.optional()
			.default('')
	});

	const baseSchema = z.object({
		name: z.string().min(1, 'Monitor name is required').max(255),
		type: z.enum(['http', 'ping']),
		interval: z.coerce.number().int().min(10, 'Interval must be at least 10 seconds'),
		failureThreshold: z.coerce.number().int().min(1).max(10),
		recoveryThreshold: z.coerce.number().int().min(1).max(10),
		notification: z.array(z.string())
	});

	const monitorFormSchema = z.discriminatedUnion('type', [
		baseSchema.extend({
			type: z.literal('http'),
			config: httpConfigSchema
		}),
		baseSchema.extend({
			type: z.literal('ping'),
			config: pingConfigSchema
		})
	]);

	type HttpConfig = z.infer<typeof httpConfigSchema>;
	type PingConfig = z.infer<typeof pingConfigSchema>;
	type MonitorFormValues = z.infer<typeof monitorFormSchema>;

	const defaultHttpConfig = (): HttpConfig => ({
		url: '',
		method: 'GET' as const,
		maxRedirects: 5,
		requestTimeoutSeconds: 10,
		headers: '',
		bodyEncoding: '' as const,
		body: '',
		acceptedStatusCodes: ['2xx'],
		upsideDownMode: false,
		certificateExpiryNotification: true,
		ignoreTlsError: false
	});

	const defaultPingConfig = (): PingConfig => ({
		host: '',
		timeoutSeconds: 5,
		packetSize: ''
	});

	const initialIntervalIndex = 3;
	const initialFailureThreshold = '1';
	const initialRecoveryThreshold = '1';

	let { notifications }: { notifications: Notification[] } = $props();
	let selectedMonitorType = $state<MonitorType>('http');
	let intervalIndex = $state<number>(initialIntervalIndex);
	let failureThresholdValue = $state<string>(initialFailureThreshold);
	let recoveryThresholdValue = $state<string>(initialRecoveryThreshold);
	let selectedNotificationIds = $state<string[]>([]);

	const notificationOptions: MultiSelectOption[] = $derived.by(() =>
		notifications.map((notification) => ({
			label: notification.name,
			value: notification.id,
			keywords: [notification.type, notification.name],
			icon: decidedNotificationIcon(notification)
		}))
	);

	const initialValues: MonitorFormValues = {
		name: '',
		type: 'http',
		interval: intervalOptions[initialIntervalIndex].seconds,
		failureThreshold: Number(initialFailureThreshold),
		recoveryThreshold: Number(initialRecoveryThreshold),
		notification: [] as string[],
		config: defaultHttpConfig()
	};

	const { form, errors, setFields } = createForm<MonitorFormValues>({
		initialValues,
		extend: validator({ schema: monitorFormSchema }),
		onSubmit: async (values) => {
			// TODO: wire to API
			console.log('submit monitor payload', values);
		}
	});

	$effect(() => {
		setFields('interval', intervalOptions[intervalIndex].seconds);
	});

	$effect(() => {
		setFields('failureThreshold', Number(failureThresholdValue));
	});

	$effect(() => {
		setFields('recoveryThreshold', Number(recoveryThresholdValue));
	});

	$effect(() => {
		setFields('notification', selectedNotificationIds);
	});

	function handleTypeChange(next: MonitorType) {
		selectedMonitorType = next;
		setFields('type', next);
		const configValue: MonitorFormValues['config'] =
			next === 'http' ? defaultHttpConfig() : defaultPingConfig();
		setFields('config', configValue);
	}

	const failureThresholdLabel = $derived.by(() => {
		const match = thresholdOptions.find(
			(option) => option.value.toString() === failureThresholdValue
		);
		return match?.label ?? 'Select threshold';
	});

	const recoveryThresholdLabel = $derived.by(() => {
		const match = thresholdOptions.find(
			(option) => option.value.toString() === recoveryThresholdValue
		);
		return match?.label ?? 'Select threshold';
	});

	const failureThresholdHelper = $derived.by(() => {
		if (failureThresholdValue === '1') return 'Immediately after the first failed check.';
		return `After ${failureThresholdValue} consecutive failed checks.`;
	});

	const recoveryThresholdHelper = $derived.by(() => {
		if (recoveryThresholdValue === '1') return 'After a single successful check.';
		return `After ${recoveryThresholdValue} consecutive successful checks.`;
	});
</script>

<form use:form class="space-y-4 w-full">

<Card.Root class="mx-auto w-full">
	<Card.Content>
			<Field.Set>
				<div class="space-y-2">
					<Field.Label for="monitor-name">Monitor name</Field.Label>
					<Input id="monitor-name" name="name" type="text" placeholder="My Monitor" required />
					<Field.Description>The name of your monitor.</Field.Description>
					{#if $errors.name}
						<Field.Description class="text-destructive">
							{$errors.name[0]}
						</Field.Description>
					{/if}
				</div>

				<div class="space-y-2">
					<Field.Label>Monitor type</Field.Label>
					<Field.Description>Select the type of monitor you want to create.</Field.Description>
					<RadioGroup.Root
						bind:value={selectedMonitorType}
						onValueChange={(value: string) => handleTypeChange(value as MonitorType)}
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
				</div>

				<div class="space-y-2">
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
					<input type="hidden" name="interval" value={intervalOptions[intervalIndex].seconds} />

					<div class="mt-2 flex justify-between text-xs text-muted-foreground">
						{#each intervalOptions as option (option.seconds)}
							<span class="w-6 text-center">
								{option.label}
							</span>
						{/each}
					</div>
				</div>

				<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
					<div class="space-y-2">
						<Field.Label>Failure threshold</Field.Label>
						<Field.Description>{failureThresholdHelper}</Field.Description>

						<Select.Root type="single" bind:value={failureThresholdValue}>
							<Select.Trigger class="w-full justify-between">
								<span data-slot="select-value" class="text-sm font-medium">
									{failureThresholdLabel}
								</span>
							</Select.Trigger>
							<Select.Content>
								<Select.Group>
									{#each thresholdOptions as option (option.value)}
										<Select.Item value={option.value.toString()}>{option.label}</Select.Item>
									{/each}
								</Select.Group>
							</Select.Content>
						</Select.Root>
						<input type="hidden" name="failureThreshold" value={failureThresholdValue} />
						{#if $errors.failureThreshold}
							<Field.Description class="text-destructive">
								{$errors.failureThreshold[0]}
							</Field.Description>
						{/if}
					</div>

					<div class="space-y-2">
						<Field.Label>Recovery threshold</Field.Label>
						<Field.Description>{recoveryThresholdHelper}</Field.Description>

						<Select.Root type="single" bind:value={recoveryThresholdValue}>
							<Select.Trigger class="w-full justify-between">
								<span data-slot="select-value" class="text-sm font-medium">
									{recoveryThresholdLabel}
								</span>
							</Select.Trigger>
							<Select.Content>
								<Select.Group>
									{#each thresholdOptions as option (option.value)}
										<Select.Item value={option.value.toString()}>{option.label}</Select.Item>
									{/each}
								</Select.Group>
							</Select.Content>
						</Select.Root>
						<input type="hidden" name="recoveryThreshold" value={recoveryThresholdValue} />
						{#if $errors.recoveryThreshold}
							<Field.Description class="text-destructive">
								{$errors.recoveryThreshold[0]}
							</Field.Description>
						{/if}
					</div>
				</div>

				<div class="space-y-2">
					<Field.Label>Notifications</Field.Label>
					<Field.Description>Choose where alerts should be sent.</Field.Description>
					<MultiSelect
						name="notification"
						bind:value={selectedNotificationIds}
						options={notificationOptions}
						placeholder="Select notifications"
						emptyMessage="No matching notification channels"
						maxBadges={4}
					/>
					{#if $errors.notification}
						<Field.Description class="text-destructive">
							{$errors.notification[0]}
						</Field.Description>
					{/if}
				</div>
			</Field.Set>
	</Card.Content>
</Card.Root>

{#if selectedMonitorType === 'http'}
	<HttpMonitor errors={$errors} />
{:else}
	<PingMonitor errors={$errors} />
{/if}

<div class="flex gap-2 justify-end">
	<Button size="default" variant="secondary">Cancel</Button>
	<Button size="default" variant="default">Create</Button>
</div>

</form>
