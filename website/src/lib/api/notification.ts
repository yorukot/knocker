import { apiFetch, type ApiResponse } from './client';

export type NotificationType = 'discord' | 'telegram';

export type ApiNotification = {
	id: string;
	team_id: string;
	type: NotificationType;
	name: string;
	config: Record<string, unknown>;
	updated_at: string;
	created_at: string;
};

export type NotificationListResponse = ApiResponse<ApiNotification[]>;

export const listTeamNotifications = async (teamId: string) => {
	const { response, body } = await apiFetch<ApiNotification[]>(
		`/api/teams/${teamId}/notifications`
	);

	return { response, body };
};