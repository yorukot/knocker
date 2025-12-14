<script lang="ts">
	import { superForm } from 'sveltekit-superforms';
	import { zodClient } from 'sveltekit-superforms/adapters';
	import { monitorCreateSchema } from '$lib/schemas/monitor';
	import type { MonitorCreate } from '$lib/schemas/monitor';

	import BaseConfigCard from './base-config-card.svelte';
	import HttpConfigCard from './http-config-card.svelte';
	import PingConfigCard from './ping-config-card.svelte';

	import type { MonitorKind } from '../../../../types/monitor-create';

	let { teamId }: { teamId: string } = $props();
	// teamId will be used for API submission: POST to `/api/teams/${teamId}/monitors`

	const initialData: MonitorCreate = {
		name: '',
		type: 'http',
		interval: 60,
		failure_threshold: 3,
		recovery_threshold: 1,
		notification: [],
		config: {
			url: '',
			method: 'GET',
			max_redirects: 5,
			request_timeout: 5,
			headers: {},
			body_encoding: '',
			body: '',
			upside_down_mode: false,
			certificate_expiry_notification: false,
			ignore_tls_error: false,
			accepted_status_codes: [200, 201, 202, 204]
		}
	};

	const typeOptions = [
		{
			id: 'http',
			label: 'HTTP',
			description: 'Check an HTTP(S) endpoint with method, headers, and status validation.'
		},
		{
			id: 'ping',
			label: 'Ping',
			description: 'TCP ping a host or IP with configurable timeout and packet size.'
		}
	] satisfies { id: MonitorKind; label: string; description: string }[];

	const superFormInstance = superForm(initialData, {
		// Type assertion needed due to Zod's discriminated union type not matching zodClient's expected type
		// The schema is valid and will work correctly at runtime
		validators: zodClient(monitorCreateSchema as unknown as Parameters<typeof zodClient>[0]),
		SPA: true,
		dataType: 'json',
		onUpdate: ({ form: formResult }) => {
			if (formResult.valid) {
				console.log('Form is valid for team:', teamId, formResult.data);
				// TODO: Submit to API at `/api/teams/${teamId}/monitors`
			}
		}
	});

	const { form: formData, enhance } = superFormInstance;

	const handleTypeChange = (next: MonitorKind) => {
		if ($formData.type === next) return;

		$formData.type = next;

		if (next === 'http') {
			$formData.config = {
				url: '',
				method: 'GET',
				max_redirects: 5,
				request_timeout: 5,
				headers: {},
				body_encoding: '',
				body: '',
				upside_down_mode: false,
				certificate_expiry_notification: false,
				ignore_tls_error: false,
				accepted_status_codes: [200, 201, 202, 204]
			};
		} else {
			$formData.config = {
				host: '',
				timeout_seconds: 5,
				packet_size: null
			};
		}
	};
</script>

<form method="POST" use:enhance class="grid gap-4">
	<div class="space-y-4">
		<BaseConfigCard
			form={superFormInstance}
			type={$formData.type}
			{typeOptions}
			{handleTypeChange}
		/>

		{#if $formData.type === 'http'}
			<HttpConfigCard form={superFormInstance} />
		{:else}
			<PingConfigCard form={superFormInstance} />
		{/if}
	</div>
</form>
