import { PUBLIC_API_BASE } from "$env/static/public";

const API_BASE = (PUBLIC_API_BASE ?? "").replace(/\/+$/, "");
const AUTH_BASE_PATH = `${API_BASE}/api/auth`;

type SuccessResponse<T> = {
	data?: T;
	message?: string;
};

const parseResponse = async (response: Response) => {
	try {
		return await response.json();
	} catch {
		return {};
	}
};

export const refreshAccessToken = async () => {
	const response = await fetch(`${AUTH_BASE_PATH}/refresh`, {
		method: "POST",
		credentials: "include",
	});

	const payload = (await parseResponse(response)) as SuccessResponse<never>;

  return payload;
};

export const login = async (email: string, password: string) => {
	const response = await fetch(`${AUTH_BASE_PATH}/login`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ email, password }),
		credentials: "include",
	});

	const payload = (await parseResponse(response)) as SuccessResponse<never>;

  return payload;
};

export const registerUser = async (displayName: string, email: string, password: string) => {
	const response = await fetch(`${AUTH_BASE_PATH}/register`, {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ display_name: displayName, email, password }),
		credentials: "include",
	});

	const payload = (await parseResponse(response)) as SuccessResponse<never>;

	return payload;
};

export const buildOAuthUrl = (provider: string, nextPath = "/") => {
	const params = new URLSearchParams({ next: nextPath });
	return `${AUTH_BASE_PATH}/oauth/${provider}?${params.toString()}`;
};
