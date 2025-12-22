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
