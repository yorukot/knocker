import type { MonitorWithIncidents } from '../../types';
import { apiRequest } from './utils';

export type MonitorResponse = {
	message: string;
	data: MonitorWithIncidents[];
};

export function getMonitors(teamID: string): Promise<MonitorResponse> {
	return apiRequest<MonitorResponse>(`/teams/${teamID}/monitors`, {
		defaultError: 'Failed to fetch monitors'
	});
}
