// ============================================================================
// Incident Types
// ============================================================================

export type IncidentStatus =
	| 'detected'
	| 'investigating'
	| 'identified'
	| 'monitoring'
	| 'resolved';

export type IncidentEventType =
	| 'detected'
	| 'notification_sent'
	| 'manually_resolved'
	| 'auto_resolved'
	| 'unpublished'
	| 'published'
	| 'investigating'
	| 'identified'
	| 'update'
	| 'monitoring';

export interface Incident {
	id: string;
	monitorId: string;
	status: IncidentStatus;
	startedAt: string;
	resolvedAt?: string;
	createdAt: string;
	updatedAt: string;
}

export interface IncidentEvent {
	id: string;
	incidentId: string;
	createdBy?: string;
	message: string;
	eventType: IncidentEventType;
	public: boolean;
	createdAt: string;
	updatedAt: string;
}
