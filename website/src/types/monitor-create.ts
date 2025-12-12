export type MonitorKind = 'http' | 'ping';

export type HttpMethod =
	| 'GET'
	| 'POST'
	| 'PUT'
	| 'DELETE'
	| 'PATCH'
	| 'HEAD'
	| 'OPTIONS';

export type BodyEncoding = 'json' | 'xml' | '';

export type HttpMonitorConfig = {
	url: string;
	method: HttpMethod;
	max_redirects: number;
	request_timeout: number;
	headers: Record<string, string>;
	body_encoding?: BodyEncoding;
	body?: string;
	upside_down_mode: boolean;
	certificate_expiry_notification: boolean;
	ignore_tls_error: boolean;
	accepted_status_codes: number[];
};

export type PingMonitorConfig = {
	host: string;
	timeout_seconds: number;
	packet_size?: number | null;
};

export type MonitorBaseSettings = {
	name: string;
	interval: number;
	failure_threshold: number;
	recovery_threshold: number;
	notification: string[];
};

export type MonitorConfigByType = {
	http: HttpMonitorConfig;
	ping: PingMonitorConfig;
};

export type MonitorCreateState<T extends MonitorKind = MonitorKind> = MonitorBaseSettings & {
	type: T;
	config: MonitorConfigByType[T];
};
