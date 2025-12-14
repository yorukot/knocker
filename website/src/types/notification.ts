// ============================================================================
// Notification Types
// ============================================================================

export type NotificationType = 'discord' | 'telegram' | 'email';

export interface Notification {
	id: string;
	teamId: string;
	type: NotificationType;
	name: string;
	config: DiscordNotificationConfig | TelegramNotificationConfig;
	updatedAt: string;
	createdAt: string;
}

export interface DiscordNotificationConfig {
	webhookUrl: string;
}

export interface TelegramNotificationConfig {
	botToken: string;
	chatId: string;
}
