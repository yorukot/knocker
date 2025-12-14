import { z } from 'zod';

// HTTP Monitor Schema
export const httpMonitorConfigSchema = z.object({
	url: z.string().url('Must be a valid URL').min(1, 'URL is required'),
	method: z.enum(['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']),
	max_redirects: z.number().int().min(0).max(1000),
	request_timeout: z.number().int().min(0).max(120),
	headers: z.record(z.string(), z.string()),
	body_encoding: z.enum(['json', 'xml', '']).optional().default(''),
	body: z.string().optional().default(''),
	upside_down_mode: z.boolean(),
	certificate_expiry_notification: z.boolean(),
	ignore_tls_error: z.boolean(),
	accepted_status_codes: z
		.array(z.number().int().min(100).max(599))
		.min(1, 'At least one status code required')
});

// Ping Monitor Schema
export const pingMonitorConfigSchema = z.object({
	host: z.string().min(1, 'Host is required'),
	timeout_seconds: z.number().int().min(0).max(120),
	packet_size: z.number().int().min(1).max(65000).nullable().optional()
});

// Base Monitor Settings Schema
export const monitorBaseSettingsSchema = z.object({
	name: z.string().min(1, 'Monitor name is required').max(255),
	interval: z.number().int().min(10, 'Interval must be at least 10 seconds'),
	failure_threshold: z.number().int().min(1).max(10),
	recovery_threshold: z.number().int().min(1).max(10),
	notification: z.array(z.string())
});

// Flexible schema that allows either config type
export const monitorCreateSchema = z
	.object({
		name: z.string().min(1, 'Monitor name is required').max(255),
		type: z.enum(['http', 'ping']),
		interval: z.number().int().min(10, 'Interval must be at least 10 seconds'),
		failure_threshold: z.number().int().min(1).max(10),
		recovery_threshold: z.number().int().min(1).max(10),
		notification: z.array(z.string()),
		config: z.union([httpMonitorConfigSchema, pingMonitorConfigSchema])
	})
	.superRefine((data, ctx) => {
		// Runtime validation based on type
		if (data.type === 'http') {
			const result = httpMonitorConfigSchema.safeParse(data.config);
			if (!result.success) {
				result.error.issues.forEach((issue) => {
					ctx.addIssue({
						...issue,
						path: ['config', ...issue.path]
					});
				});
			}
		} else if (data.type === 'ping') {
			const result = pingMonitorConfigSchema.safeParse(data.config);
			if (!result.success) {
				result.error.issues.forEach((issue) => {
					ctx.addIssue({
						...issue,
						path: ['config', ...issue.path]
					});
				});
			}
		}
	});

export type HttpMonitorConfig = z.infer<typeof httpMonitorConfigSchema>;
export type PingMonitorConfig = z.infer<typeof pingMonitorConfigSchema>;
export type MonitorBaseSettings = z.infer<typeof monitorBaseSettingsSchema>;
export type MonitorCreate = z.infer<typeof monitorCreateSchema>;
