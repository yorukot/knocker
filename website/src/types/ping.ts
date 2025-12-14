// ============================================================================
// Ping Types
// ============================================================================

export type PingStatus = 'successful' | 'failed' | 'timeout';

export interface Ping {
	time: string;
	monitorId: string;
	latency: number;
	status: PingStatus;
	region: string;
}
