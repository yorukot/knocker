import type { PageLoad } from './$types';
import { listTeamNotifications } from '$lib/api/notification';

export const load: PageLoad = async ({ params }) => {
	const { teamID } = params;

	try {
		const { body } = await listTeamNotifications(teamID);
		const notifications = body.data ?? [];

		return {
			notifications
		};
	} catch (error) {
		console.error('Failed to load notifications:', error);
		return {
			notifications: []
		};
	}
};
