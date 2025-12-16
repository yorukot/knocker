import type { Monitor, MonitorWithIncidents } from '../../types';
import { apiRequest } from './utils';

export type MonitorResponse = {
	message: string;
	data: MonitorWithIncidents[];
};

export type MonitorCreateRequest = {
	name: string;
	type: 'http' | 'ping';
	interval: number;
	failure_threshold: number;
	recovery_threshold: number;
	notification: string[];
	config:
		| {
			url: string;
			method: string;
			max_redirects: number;
			request_timeout: number;
			headers?: Record<string, string>;
			body_encoding?: 'json' | 'xml';
			body?: string;
			upside_down_mode: boolean;
			certificate_expiry_notification: boolean;
			ignore_tls_error: boolean;
			accepted_status_codes: number[];
		}
		| {
			host: string;
			timeout_seconds: number;
			packet_size?: number;
		};
};

export type MonitorCreateResponse = {
	message: string;
	data: Monitor;
};

export type MonitorDeleteResponse = {
	message: string;
};

export function getMonitors(teamID: string): Promise<MonitorResponse> {
	return apiRequest<MonitorResponse>(`/teams/${teamID}/monitors`, {
		defaultError: 'Failed to fetch monitors'
	});
}

export function createMonitor(
	teamID: string,
	payload: MonitorCreateRequest
): Promise<MonitorCreateResponse> {
	return apiRequest<MonitorCreateResponse>(`/teams/${teamID}/monitors`, {
		method: 'POST',
		body: payload,
		defaultError: 'Failed to create monitor'
	});
}

export function deleteMonitor(teamID: string, monitorID: string): Promise<MonitorDeleteResponse> {
	return apiRequest<MonitorDeleteResponse>(`/teams/${teamID}/monitors/${monitorID}`, {
		method: 'DELETE',
		defaultError: 'Failed to delete monitor'
	});
}
