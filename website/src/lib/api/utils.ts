import { PUBLIC_API_BASE } from '$env/static/public';
import { refreshToken } from './auth';
import camelcaseKeys from 'camelcase-keys';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

type ApiErrorBody = {
	message?: string;
};

export type ApiDefaultBody = {
	message: string;
};

type ApiOptions = {
	method?: HttpMethod;
	body?: unknown;
	defaultError?: string;
	headers?: HeadersInit;
};

async function parseJson(res: Response): Promise<unknown> {
	try {
		return await res.json();
	} catch {
		return null;
	}
}

function normalizeResponse<T>(data: unknown): T {
	if (data && typeof data === 'object') {
		return camelcaseKeys(data as object, { deep: true }) as T;
	}
	return data as T;
}

export async function apiRequest<T>(url: string, options: ApiOptions = {}): Promise<T> {
	const { method = 'GET', body, defaultError = 'Request failed', headers } = options;

	const requestInit: RequestInit = {
		method,
		credentials: 'include',
		headers: {
			...(body ? { 'Content-Type': 'application/json' } : {}),
			...headers
		},
		body: body ? JSON.stringify(body) : undefined
	};

	const res = await fetch(`${PUBLIC_API_BASE}${url}`, requestInit);

	const responseBody = await parseJson(res);

	/* ---------- auth retry ---------- */

	if (res.status === 401) {
		const refreshed = await refreshToken();

		if (!refreshed) {
			throw new Error('AUTH_EXPIRED');
		}

		const retryRes = await fetch(`${PUBLIC_API_BASE}${url}`, requestInit);

		if (!retryRes.ok) {
			throw new Error('AUTH_EXPIRED');
		}

		const retryBody = await parseJson(retryRes);
		return normalizeResponse<T>(retryBody);
	}

	/* ---------- error handling ---------- */

	if (!res.ok) {
		const errorBody = responseBody as ApiErrorBody | null;
		const message =
			typeof errorBody?.message === 'string'
				? `${defaultError}: ${errorBody.message}`
				: defaultError;

		throw new Error(message);
	}

	/* ---------- success ---------- */

	return normalizeResponse<T>(responseBody);
}
