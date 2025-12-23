// ============================================================================
// Team Types
// ============================================================================

export type MemberRole = 'owner' | 'admin' | 'member' | 'viewer';

export interface Team {
	id: string;
	name: string;
	updatedAt: string;
	createdAt: string;
}

export interface TeamMember {
	id: string;
	teamId: string;
	userId: string;
	role: MemberRole;
	updatedAt: string;
	createdAt: string;
}

export interface TeamInvite {
	id: string;
	teamId: string;
	invitedBy: string;
	invitedTo: string;
	updatedAt: string;
	createdAt: string;
}

export interface TeamWithRole extends Team {
	role: MemberRole;
}
