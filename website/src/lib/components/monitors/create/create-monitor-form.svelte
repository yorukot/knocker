<script lang="ts">
	import Icon from "@iconify/svelte";

	import BaseConfigCard from "./base-config-card.svelte";
	import HttpConfigCard from "./http-config-card.svelte";
	import PingConfigCard from "./ping-config-card.svelte";

	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "$lib/components/ui/card/index.js";
	import { Textarea } from "$lib/components/ui/textarea/index.js";

	import type {
		HttpMonitorConfig,
		MonitorConfigByType,
		MonitorCreateState,
		MonitorKind,
		PingMonitorConfig,
	} from "../../../../types/monitor-create";

	let { teamId }: { teamId: string } = $props();

	const defaultConfigs: MonitorConfigByType = {
		http: {
			url: "",
			method: "GET",
			max_redirects: 5,
			request_timeout: 5,
			headers: {},
			body_encoding: "",
			body: "",
			upside_down_mode: false,
			certificate_expiry_notification: false,
			ignore_tls_error: false,
			accepted_status_codes: [200, 201, 202, 204],
		},
		ping: {
			host: "",
			timeout_seconds: 5,
			packet_size: null,
		},
	};

	const typeOptions = [
		{
			id: "http",
			label: "HTTP",
			description: "Check an HTTP(S) endpoint with method, headers, and status validation.",
		},
		{
			id: "ping",
			label: "Ping",
			description: "TCP ping a host or IP with configurable timeout and packet size.",
		},
	] satisfies { id: MonitorKind; label: string; description: string; icon: string }[];

	const clone = <T,>(value: T): T => JSON.parse(JSON.stringify(value));

	let draft = $state<MonitorCreateState>({
		name: "",
		type: "http",
		interval: 60,
		failure_threshold: 3,
		recovery_threshold: 1,
		notification: [],
		config: clone(defaultConfigs.http),
	});

	const handleTypeChange = (next: MonitorKind) => {
		if (draft.type === next) return;
		draft.type = next;
		draft.config = clone(defaultConfigs[next]);
	};

	const requestPath = $derived.by(() => `/api/teams/${teamId}/monitors`);

	const requestPreview = $derived.by(() =>
		JSON.stringify(
			{
				name: draft.name.trim(),
				type: draft.type,
				interval: draft.interval,
				failure_threshold: draft.failure_threshold,
				recovery_threshold: draft.recovery_threshold,
				notification: draft.notification,
				config: draft.config,
			},
			null,
			2,
		),
	);
</script>

<form class="grid gap-4">
	<div class="space-y-4">
		<BaseConfigCard settings={draft} type={draft.type} typeOptions={typeOptions} onTypeChange={handleTypeChange} />

		{#if draft.type === "http"}
			<HttpConfigCard bind:config={draft.config as HttpMonitorConfig} />
		{:else}
			<PingConfigCard bind:config={draft.config as PingMonitorConfig} />
		{/if}
	</div>
</form>
