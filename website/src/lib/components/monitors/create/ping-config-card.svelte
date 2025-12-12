<script lang="ts">
	import Icon from "@iconify/svelte";

	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import * as Field from "$lib/components/ui/field/index.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";

	import type { PingMonitorConfig } from "../../../../types/monitor-create";

	let {
		config = $bindable<PingMonitorConfig>({
			host: "",
			timeout_seconds: 5,
			packet_size: null,
		}),
	}: { config: PingMonitorConfig } = $props();

	let timeoutInput = $derived(config.timeout_seconds.toString());
	$effect(() => {
		timeoutInput = config.timeout_seconds.toString();
	});
	$effect(() => {
		const parsed = Number.parseInt(timeoutInput, 10);
		config.timeout_seconds = Number.isFinite(parsed) && parsed >= 0 ? parsed : 0;
	});

	let packetSizeInput = $derived(config.packet_size?.toString() ?? "");
	$effect(() => {
		packetSizeInput = config.packet_size ? config.packet_size.toString() : "";
	});
	$effect(() => {
		const parsed = Number.parseInt(packetSizeInput, 10);
		config.packet_size = Number.isFinite(parsed) && parsed > 0 ? parsed : null;
	});
</script>

<Card>
	<CardHeader class="gap-1 pb-4">
		<CardTitle class="flex items-center gap-2 text-xl">
			<Icon icon="lucide:radio-tower" class="size-5 text-primary" />
			Ping settings
		</CardTitle>
		<CardDescription>Host, timeout, and packet sizing for TCP ping monitors.</CardDescription>
	</CardHeader>
	<CardContent class="space-y-6">
		<Field.Group>
			<Field.Field>
				<Field.Label>Target host</Field.Label>
				<Field.Content class="space-y-2">
					<Input placeholder="api.internal.local" bind:value={config.host} required />
					<Field.Description>Hostname or IP to ping from the worker.</Field.Description>
				</Field.Content>
			</Field.Field>
		</Field.Group>

		<Field.Group>
			<Field.Field orientation="responsive">
				<Field.Label>Timeout &amp; packet size</Field.Label>
				<Field.Content class="grid gap-3 md:grid-cols-2">
					<div class="space-y-1.5">
						<Input
							type="number"
							min="0"
							max="120"
							step="1"
							inputmode="numeric"
							bind:value={timeoutInput}
							placeholder="5"
						/>
						<Field.Description>Timeout in seconds before marking a failed attempt.</Field.Description>
					</div>
					<div class="space-y-1.5">
						<Input
							type="number"
							min="1"
							max="65000"
							step="1"
							inputmode="numeric"
							bind:value={packetSizeInput}
							placeholder="Auto"
						/>
						<Field.Description>
							Optional packet size in bytes. Leave blank to use the default for the runtime.
						</Field.Description>
					</div>
				</Field.Content>
			</Field.Field>
		</Field.Group>

		<Field.Group>
			<Field.Field>
				<Field.Label>Notes</Field.Label>
				<Field.Content class="space-y-2">
					<Textarea
						rows={3}
						placeholder="Document firewall expectations, VPN access, or maintenance windows."
						readonly
						aria-readonly="true"
						class="opacity-80"
					/>
					<Field.Description>
						Use this space to capture operational notes when we wire editing and storage.
					</Field.Description>
				</Field.Content>
			</Field.Field>
		</Field.Group>
	</CardContent>
</Card>
