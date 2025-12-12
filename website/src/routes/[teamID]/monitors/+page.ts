import { error, redirect } from '@sveltejs/kit';
import { listTeamMonitors } from '$lib/api/monitor';
import type { PageLoad } from './$types';
import type {
	ApiMonitor,
	ApiMonitorIncident,
	ApiMonitorIncidentStatus,
	MonitorIncident,
	MonitorIncidentStatus,
	MonitorListItem,
	MonitorStatus,
	MonitorType
} from '../../../types/monitor';

const rtf = new Intl.RelativeTimeFormat('en', { numeric: 'auto' });

const parseRelativeTime = (value?: string | null) => {
	if (!value) return '—';

	const date = new Date(value);
	if (Number.isNaN(date.getTime())) return '—';

	let duration = (date.getTime() - Date.now()) / 1000;
	const divisions = [
		{ amount: 60, unit: 'second' as const },
		{ amount: 60, unit: 'minute' as const },
		{ amount: 24, unit: 'hour' as const },
		{ amount: 7, unit: 'day' as const },
		{ amount: 4.34524, unit: 'week' as const },
		{ amount: 12, unit: 'month' as const },
		{ amount: Infinity, unit: 'year' as const }
	];

	for (const division of divisions) {
		if (Math.abs(duration) < division.amount) {
			return rtf.format(Math.round(duration), division.unit);
		}
		duration /= division.amount;
	}

	return '—';
};

const formatFrequency = (interval: number | undefined | null) => {
	if (!interval || interval <= 0) return 'Paused';
	if (interval < 60) return `Every ${interval}s`;

	const minutes = interval / 60;
	if (interval % 3600 === 0) {
		const hours = interval / 3600;
		return `Every ${hours}h`;
	}

	if (minutes >= 1 && minutes < 90) {
		return `Every ${Math.round(minutes)}m`;
	}

	const hours = interval / 3600;
	return `Every ${Number(hours.toFixed(1)).toString()}h`;
};

const parseConfig = (config: unknown) => {
	if (!config) return {};
	if (typeof config === 'string') {
		try {
			return JSON.parse(config) as Record<string, unknown>;
		} catch {
			return {};
		}
	}

	if (typeof config === 'object') return config as Record<string, unknown>;
	return {};
};

const resolveTarget = (monitor: ApiMonitor) => {
	const config = parseConfig(monitor.config);

	if (monitor.type === 'http' && typeof config.url === 'string') {
		return config.url;
	}

	if (monitor.type === 'ping' && typeof config.host === 'string') {
		return config.host;
	}

	return '—';
};

const mapMonitorType = (type: ApiMonitor['type']): MonitorType => {
	const normalized = type?.toLowerCase?.() ?? '';
	if (normalized === 'http') return 'HTTP';
	if (normalized === 'ping') return 'Ping';
	return 'TCP';
};

const mapIncidentStatus = (status: ApiMonitorIncidentStatus): MonitorIncidentStatus => {
	if (status === 'monitoring') return 'monitoring';
	if (status === 'identified') return 'identified';
	if (status === 'resolved') return 'resolved';
	return 'investigating';
};

const deriveMonitorStatus = (monitor: ApiMonitor): MonitorStatus => {
	if (!monitor.interval || monitor.interval <= 0) return 'paused';

	const activeIncident = (monitor.incidents ?? []).find((incident) => !incident.resolved_at);

	if (!activeIncident) return 'operational';

	if (activeIncident.status === 'monitoring') return 'degraded';
	if (activeIncident.status === 'resolved') return 'operational';

	return 'down';
};

const buildIncidentSummary = (incident: ApiMonitorIncident) => {
	const status = mapIncidentStatus(incident.status);
	const statusLabel = `${status.slice(0, 1).toUpperCase()}${status.slice(1)}`;

	const started = parseRelativeTime(incident.started_at);
	return started === '—' ? `Incident ${statusLabel.toLowerCase()}` : `${statusLabel} ${started}`;
};

const buildIncident = (incident: ApiMonitorIncident, teamID: string): MonitorIncident => ({
	id: incident.id,
	status: mapIncidentStatus(incident.status),
	severity: 'major',
	updatedAt: parseRelativeTime(incident.updated_at),
	summary: buildIncidentSummary(incident),
	link: `/${teamID}/incidents`
});

const formatLastIncident = (incidents: ApiMonitorIncident[] = []) => {
	if (incidents.length === 0) return 'No incidents yet';

	const [latest] = incidents;

	if (latest.resolved_at) {
		return `Resolved ${parseRelativeTime(latest.resolved_at)}`;
	}

	return `Started ${parseRelativeTime(latest.started_at)}`;
};

const mapMonitor = (monitor: ApiMonitor, teamID: string): MonitorListItem => {
	const incidents = monitor.incidents ?? [];
	const activeIncident = incidents.find((incident) => !incident.resolved_at);

	return {
		id: monitor.id,
		name: monitor.name,
		target: resolveTarget(monitor),
		type: mapMonitorType(monitor.type),
		status: deriveMonitorStatus(monitor),
		regions: [],
		frequency: formatFrequency(monitor.interval),
		uptime: '—',
		responseTime: '—',
		lastChecked: parseRelativeTime(monitor.last_checked),
		lastIncident: formatLastIncident(incidents),
		incident: activeIncident ? buildIncident(activeIncident, teamID) : undefined
	};
};

export const load: PageLoad = async ({ params, url }) => {
	const { response, body } = await listTeamMonitors(params.teamID);

	if (response.status === 401) {
		const next = encodeURIComponent(url.pathname + url.search);
		redirect(302, `/auth/login?next=${next}`);
	}

	if (response.status === 404) {
		throw error(404, 'Team not found');
	}

	if (!response.ok || !body.data) {
		throw error(response.status || 500, body.message ?? 'Failed to fetch monitors');
	}

	return {
		monitors: body.data.map((monitor) => mapMonitor(monitor, params.teamID))
	};
};
