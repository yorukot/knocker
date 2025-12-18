import type { Team } from '../../types';
import { apiRequest } from './utils';

export type TeamsResponse = {
  message: string;
  data: Team[];
};

export function getTeams(): Promise<TeamsResponse> {
	return apiRequest<TeamsResponse>("/teams", { defaultError: 'Failed to fetch teams' });
}

export type CreateTeamResponse = {
	message: string;
	data: Team;
};

export function createTeam(name: string): Promise<CreateTeamResponse> {
	return apiRequest<CreateTeamResponse>('/teams', {
		method: 'POST',
		body: { name },
		defaultError: 'Failed to create team'
	});
}
