import type { Notification } from '../../types';
import { apiRequest } from './utils';

export type NotificationsResponse = {
	message: string;
	data: Notification[];
};

export function getNotifications(teamID: string): Promise<NotificationsResponse> {
	return apiRequest<NotificationsResponse>(`/teams/${teamID}/notifications`, {
		defaultError: 'Failed to fetch notifications'
	});
}
