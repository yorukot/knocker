<script lang="ts">
	import Icon from '@iconify/svelte';

	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { NativeSelect, NativeSelectOption } from '$lib/components/ui/native-select/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';

	import type { SuperForm } from 'sveltekit-superforms';
	import type { MonitorCreate } from '$lib/schemas/monitor';

	let {
		form
	}: {
		form: SuperForm<MonitorCreate>;
	} = $props();

	const formData = $derived(form.form);

	const methods = ['GET', 'HEAD', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'] as const;

	// Type guard helper - only works when type is 'http'
	const httpConfig = $derived('url' in $formData.config ? $formData.config : null);

	let statusInput = $state('200, 201, 202, 204');
	$effect(() => {
		if (httpConfig?.accepted_status_codes) {
			statusInput = httpConfig.accepted_status_codes.join(', ');
		}
	});
	$effect(() => {
		if (!httpConfig) return;
		const parsed = statusInput
			.split(',')
			.map((value) => Number.parseInt(value.trim(), 10))
			.filter((value) => Number.isFinite(value) && value >= 100 && value <= 599);
		httpConfig.accepted_status_codes =
			parsed.length > 0 ? Array.from(new Set(parsed)) : [200, 201, 202, 204];
	});

	const formatHeaders = (headers: Record<string, string>) =>
		Object.entries(headers ?? {})
			.map(([key, value]) => `${key}: ${value}`)
			.join('\n');

	const parseHeaders = (input: string) => {
		const entries = input
			.split('\n')
			.map((line) => line.trim())
			.filter(Boolean)
			.map((line) => {
				const [key, ...rest] = line.split(':');
				const value = rest.join(':').trim();
				return key && value ? [key.trim(), value] : null;
			})
			.filter(Boolean) as [string, string][];

		return Object.fromEntries(entries);
	};

	let headersInput = $state('');
	$effect(() => {
		if (httpConfig?.headers) {
			headersInput = formatHeaders(httpConfig.headers);
		}
	});
	$effect(() => {
		if (httpConfig) {
			httpConfig.headers = parseHeaders(headersInput);
		}
	});
</script>

{#if httpConfig}
	<Card>
		<CardHeader class="gap-1 pb-4">
			<CardTitle class="flex items-center gap-2 text-xl">
				<Icon icon="lucide:globe-2" class="size-5 text-primary" />
				HTTP settings
			</CardTitle>
			<CardDescription>URL, method, validation, and TLS controls for HTTP monitors.</CardDescription
			>
		</CardHeader>
		<CardContent class="space-y-6">
			<Field.Group>
				<Field.Field>
					<Field.Label>Request target</Field.Label>
					<Field.Content class="space-y-2">
						<Input
							placeholder="https://status.yourapp.com/health"
							bind:value={httpConfig.url}
							required
						/>
						<Field.Description>Public endpoint we should check from the worker.</Field.Description>
					</Field.Content>
				</Field.Field>

				<Field.Field>
					<Field.Label>Method &amp; timeout</Field.Label>
					<Field.Content class="grid gap-3 md:grid-cols-2">
						<div class="space-y-1.5">
							<NativeSelect aria-label="HTTP method" bind:value={httpConfig.method}>
								{#each methods as method (method)}
									<NativeSelectOption value={method}>{method}</NativeSelectOption>
								{/each}
							</NativeSelect>
							<Field.Description>Verb used for each probe.</Field.Description>
						</div>
						<div class="space-y-1.5">
							<Input
								type="number"
								min="0"
								max="120"
								step="1"
								inputmode="numeric"
								bind:value={httpConfig.request_timeout}
								placeholder="5"
							/>
							<Field.Description>Request timeout in seconds.</Field.Description>
						</div>
					</Field.Content>
				</Field.Field>
			</Field.Group>

			<Field.Group>
				<Field.Field>
					<Field.Label>Accepted status codes</Field.Label>
					<Field.Content class="space-y-2">
						<Input placeholder="200, 201, 202, 204" bind:value={statusInput} autocomplete="off" />
						<Field.Description>
							Comma-separated list (100-599). Any other code marks the monitor as failed.
						</Field.Description>
					</Field.Content>
				</Field.Field>

				<Field.Field>
					<Field.Label>Redirects &amp; TLS</Field.Label>
					<Field.Content class="grid gap-3 md:grid-cols-3">
						<div class="space-y-1.5">
							<Input
								type="number"
								min="0"
								max="1000"
								step="1"
								inputmode="numeric"
								bind:value={httpConfig.max_redirects}
								placeholder="5"
							/>
							<Field.Description>Max redirects to follow.</Field.Description>
						</div>
						<label
							class="flex items-start gap-3 rounded-lg border border-input/70 bg-card/40 px-3 py-2"
						>
							<Switch bind:checked={httpConfig.ignore_tls_error} aria-label="Ignore TLS errors" />
							<div class="space-y-0.5">
								<p class="text-sm font-medium leading-tight">Ignore TLS errors</p>
								<p class="text-xs text-muted-foreground">
									Useful for self-signed internal endpoints.
								</p>
							</div>
						</label>
						<label
							class="flex items-start gap-3 rounded-lg border border-input/70 bg-card/40 px-3 py-2"
						>
							<Switch
								bind:checked={httpConfig.certificate_expiry_notification}
								aria-label="Warn on certificate expiry"
							/>
							<div class="space-y-0.5">
								<p class="text-sm font-medium leading-tight">Warn on certificate expiry</p>
								<p class="text-xs text-muted-foreground">
									Receive a heads-up before TLS certs lapse.
								</p>
							</div>
						</label>
					</Field.Content>
				</Field.Field>
			</Field.Group>

			<Field.Group>
				<Field.Field>
					<Field.Label>Headers</Field.Label>
					<Field.Content class="space-y-2">
						<Textarea rows={4} placeholder="Authorization: Bearer ****" bind:value={headersInput} />
						<Field.Description>One header per line using the format Key: Value.</Field.Description>
					</Field.Content>
				</Field.Field>

				<Field.Field>
					<Field.Label>Body (optional)</Field.Label>
					<Field.Content class="space-y-2">
						<div class="grid gap-3 md:grid-cols-5">
							<NativeSelect
								class="md:col-span-2"
								aria-label="Body encoding"
								bind:value={httpConfig.body_encoding}
							>
								<NativeSelectOption value="">No body</NativeSelectOption>
								<NativeSelectOption value="json">JSON</NativeSelectOption>
								<NativeSelectOption value="xml">XML</NativeSelectOption>
							</NativeSelect>
							<Textarea
								class="md:col-span-3"
								rows={3}
								placeholder="status payload (JSON, XML, or plain text)"
								bind:value={httpConfig.body}
							/>
						</div>
						<Field.Description>Optional payload for POST/PUT/PATCH requests.</Field.Description>
					</Field.Content>
				</Field.Field>
			</Field.Group>

			<Field.Group>
				<Field.Field>
					<Field.Label>Inverted success</Field.Label>
					<Field.Content>
						<label
							class="flex items-start gap-3 rounded-lg border border-input/70 bg-card/40 px-3 py-2"
						>
							<Switch bind:checked={httpConfig.upside_down_mode} aria-label="Upside down mode" />
							<div class="space-y-0.5">
								<p class="text-sm font-medium leading-tight">Expect failure (maintenance)</p>
								<p class="text-xs text-muted-foreground">
									Treat non-2xx responses as success—handy during controlled tests.
								</p>
							</div>
						</label>
					</Field.Content>
				</Field.Field>
			</Field.Group>
		</CardContent>
	</Card>
{/if}
