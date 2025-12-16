import type { BodyEncoding, HTTPMethod } from '../../../../types/monitor-config';
import type { MonitorType } from '../../../../types';
import type { MultiSelectOption } from '$lib/components/ui/multi-select';

export type MonitorTypeSelect = {
	title: string;
	description: string;
	value: MonitorType;
};

export const monitorTypeSelectData: readonly MonitorTypeSelect[] = [
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
] as const;

export type IntervalOption = {
	label: string;
	seconds: number;
};

export const intervalOptions: readonly IntervalOption[] = [
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
] as const;

export type ThresholdOption = {
	label: string;
	value: number;
};

export const thresholdOptions: readonly ThresholdOption[] = [
	{ label: 'Immediate (after 1 check)', value: 1 },
	{ label: 'After 2 checks', value: 2 },
	{ label: 'After 3 checks', value: 3 },
	{ label: 'After 4 checks', value: 4 },
	{ label: 'After 5 checks', value: 5 }
] as const;

export const httpMethods: readonly HTTPMethod[] = [
	'GET',
	'POST',
	'PUT',
	'DELETE',
	'PATCH',
	'HEAD',
	'OPTIONS'
] as const;

export const bodyEncodingOptions: Array<{ label: string; value: BodyEncoding | '' }> = [
	{ label: 'None', value: '' },
	{ label: 'JSON', value: 'json' },
	{ label: 'XML', value: 'xml' }
];

const statusRangeOptions: MultiSelectOption[] = [
	{ label: 'Any 2xx', value: '2xx', keywords: ['2xx', 'success'] },
	{ label: 'Any 3xx', value: '3xx', keywords: ['3xx', 'redirect'] },
	{ label: 'Any 4xx', value: '4xx', keywords: ['4xx', 'client'] },
	{ label: 'Any 5xx', value: '5xx', keywords: ['5xx', 'server'] }
];

const commonStatusOptions: MultiSelectOption[] = [
	{ label: '200 OK', value: '200' },
	{ label: '201 Created', value: '201' },
	{ label: '202 Accepted', value: '202' },
	{ label: '204 No Content', value: '204' },
	{ label: '301 Moved Permanently', value: '301' },
	{ label: '302 Found', value: '302' },
	{ label: '400 Bad Request', value: '400' },
	{ label: '401 Unauthorized', value: '401' },
	{ label: '403 Forbidden', value: '403' },
	{ label: '404 Not Found', value: '404' },
	{ label: '429 Too Many Requests', value: '429' },
	{ label: '500 Internal Server Error', value: '500' },
	{ label: '502 Bad Gateway', value: '502' },
	{ label: '503 Service Unavailable', value: '503' }
];

export const acceptedStatusOptions: MultiSelectOption[] = [...statusRangeOptions, ...commonStatusOptions];