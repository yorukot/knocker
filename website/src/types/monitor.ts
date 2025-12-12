export type ApiMonitorIncidentStatus =
	| 'detected'
	| 'investigating'
	| 'identified'
	| 'monitoring'
	| 'resolved';

export type ApiMonitorIncident = {
	id: string;
	monitor_id: string;
	status: ApiMonitorIncidentStatus;
	started_at: string;
	resolved_at?: string | null;
	created_at: string;
	updated_at: string;
};

export type ApiMonitor = {
	id: string;
	team_id: string;
	name: string;
	type: 'http' | 'ping' | string;
	config: unknown;
	interval: number;
	last_checked: string;
	next_check: string;
	failure_threshold: number;
	recovery_threshold: number;
	notification: string[];
	incidents?: ApiMonitorIncident[];
	updated_at: string;
	created_at: string;
};

export type MonitorStatus = 'operational' | 'degraded' | 'down' | 'paused';
export type MonitorType = 'HTTP' | 'Ping' | 'TCP';
export type MonitorIncidentStatus = 'investigating' | 'identified' | 'monitoring' | 'resolved';
export type MonitorIncidentSeverity = 'critical' | 'major' | 'minor' | 'maintenance';

export type MonitorIncident = {
	id: string;
	status: MonitorIncidentStatus;
	severity: MonitorIncidentSeverity;
	updatedAt: string;
	summary: string;
	link: string;
};

export type MonitorListItem = {
	id: string;
	name: string;
	target: string;
	type: MonitorType;
	status: MonitorStatus;
	regions: string[];
	frequency: string;
	uptime: string;
	responseTime: string;
	lastChecked: string;
	lastIncident: string;
	incident?: MonitorIncident;
};
