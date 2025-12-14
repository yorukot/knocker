import type { User } from '../../types';
import { apiRequest } from './utils';

export type UserResponse = {
	message: string;
	data: User;
};

export function getUser(): Promise<UserResponse> {
	return apiRequest<UserResponse>("/users/me", { defaultError: 'Failed to fetch user' });
}
