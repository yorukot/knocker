import { getMonitor } from '$lib/api/monitor';
import { getNotifications } from '$lib/api/notification';
import type { Monitor, Notification } from '../../../../../types';
import type { PageLoad } from './$types';
import { error } from '@sveltejs/kit';

export type EditMonitorData = {
	monitor: Monitor;
	notifications: Notification[];
};

export const load: PageLoad<EditMonitorData> = async ({ params }) => {
	const { teamID, monitorID } = params;

	if (!monitorID) {
		error(400, 'monitorId query parameter is required');
	}

	const [monitorRes, notificationsRes] = await Promise.all([
		getMonitor(teamID, monitorID),
		getNotifications(teamID)
	]);

	return {
		monitor: monitorRes.data,
		notifications: notificationsRes.data
	};
};
