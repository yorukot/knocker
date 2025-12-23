import type { PublicIncident } from './incident';

// ============================================================================
// Status Page Types
// ============================================================================

export type StatusPageElementType = 'historical_timeline' | 'current_status_indicator';

export interface StatusPage {
	id: string;
	teamId: string;
	title: string;
	slug: string;
	icon?: string; // base64 or URL depending on API usage
	createdAt: string;
	updatedAt: string;
}

export interface StatusPageGroup {
	id: string;
	statusPageId: string;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
}

export interface StatusPageMonitor {
	id: string;
	statusPageId: string;
	monitorId: string;
	groupId?: string | null;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
}

export interface StatusPageWithElements {
	statusPage: StatusPage;
	groups: StatusPageGroup[];
	monitors: StatusPageMonitor[];
}

export interface PublicTimelinePoint {
	day: string;
	success: number;
	fail: number;
}

export interface PublicStatusPageGroup {
	id: string;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
	status?: 'up' | 'down';
	uptimeSli30?: number;
	uptimeSli60?: number;
	uptimeSli90?: number;
	timeline?: PublicTimelinePoint[];
}

export interface PublicStatusPageMonitor {
	id: string;
	monitorId: string;
	groupId?: string | null;
	name: string;
	type: StatusPageElementType;
	sortOrder: number;
	status?: 'up' | 'down';
	uptimeSli30?: number;
	uptimeSli60?: number;
	uptimeSli90?: number;
	timeline?: PublicTimelinePoint[];
}

export interface PublicStatusPageData {
	statusPage: StatusPage;
	groups: PublicStatusPageGroup[];
	monitors: PublicStatusPageMonitor[];
	incidents: PublicIncident[];
}
