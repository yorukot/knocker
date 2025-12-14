// ============================================================================
// Monitor Types
// ============================================================================

import type { HTTPMonitorConfig, PingMonitorConfig } from './monitor-config';
import type { Incident } from './incident';

export type MonitorType = 'http' | 'ping';

export interface Monitor {
	id: string;
	teamId: string;
	name: string;
	type: MonitorType;
	config: HTTPMonitorConfig | PingMonitorConfig;
	interval: number;
	lastChecked: string;
	nextCheck: string;
	failureThreshold: number;
	recoveryThreshold: number;
	notification: string[];
	updatedAt: string;
	createdAt: string;
}

export interface MonitorWithIncidents extends Monitor {
	incidents?: Incident[];
}

export interface MonitorNotification {
	id: string;
	monitorId: string;
	notificationId: string;
}
