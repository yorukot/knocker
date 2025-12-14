import type { Team } from '../../types';
import { apiRequest } from './utils';

export type TeamsResponse = {
  message: string;
  data: Team[];
};

export function getTeams(): Promise<TeamsResponse> {
	return apiRequest<TeamsResponse>("/teams", { defaultError: 'Failed to fetch teams' });
}
