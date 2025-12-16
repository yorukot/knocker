import { getNotifications } from '$lib/api/notification';
import type { PageLoad } from './$types';
import type { Notification } from '../../../../types';

export type NewMonitorPageData = {
	notifications: Notification[];
};

export const load: PageLoad<NewMonitorPageData> = async ({ params }) => {
	const { teamID } = params;
	const response = await getNotifications(teamID);

	return {
		notifications: response.data
	};
};
