import { getMonitor } from '$lib/api/monitor';
import { getNotifications } from '$lib/api/notification';
import { getRegions } from '$lib/api/region';
import type { Monitor, Notification, Region } from '../../../../../types';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export type EditMonitorData = {
	monitor: Monitor;
	notifications: Notification[];
	regions: Region[];
};

export const load: PageLoad<EditMonitorData> = async ({ params }) => {
	const { teamID, monitorID } = params;

	if (!monitorID) {
		error(400, 'monitorId query parameter is required');
	}

	const [monitorRes, notificationsRes, regionsRes] = await Promise.all([
		getMonitor(teamID, monitorID),
		getNotifications(teamID),
		getRegions()
	]);

	const monitorRaw = monitorRes.data as any;
	const monitor: Monitor = {
		...monitorRaw,
		teamId: monitorRaw.team_id ?? monitorRaw.teamId,
		lastChecked: monitorRaw.last_checked ?? monitorRaw.lastChecked,
		nextCheck: monitorRaw.next_check ?? monitorRaw.nextCheck,
		failureThreshold: monitorRaw.failure_threshold ?? monitorRaw.failureThreshold,
		recoveryThreshold: monitorRaw.recovery_threshold ?? monitorRaw.recoveryThreshold,
		updatedAt: monitorRaw.updated_at ?? monitorRaw.updatedAt,
		createdAt: monitorRaw.created_at ?? monitorRaw.createdAt
	};

	return {
		monitor,
		notifications: notificationsRes.data,
		regions: regionsRes.data
	};
};
