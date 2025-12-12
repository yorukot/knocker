import { apiFetch, type ApiResponse } from './client';
import type { ApiMonitor } from '../../types/monitor';

export type MonitorListResponse = ApiResponse<ApiMonitor[]>;

export const listTeamMonitors = async (teamId: string) => {
	const { response, body } = await apiFetch<ApiMonitor[]>(`/api/teams/${teamId}/monitors`);

	return { response, body };
};
