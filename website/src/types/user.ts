// ============================================================================
// User Types
// ============================================================================

export interface User {
	id: string;
	passwordHash?: string;
	displayName: string;
	avatar?: string;
	createdAt: string;
	updatedAt: string;
}
