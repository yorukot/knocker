import { z } from 'zod';

export const Int64String = z.string().regex(/^-?\d+$/, 'Expected int64 string');

export const statusPageElementTypeSchema = z.enum([
	'historical_timeline',
	'current_status_indicator'
]);

export const statusPageGroupInputSchema = z.object({
	id: Int64String.optional(),
	name: z.string().min(1).max(255),
	type: statusPageElementTypeSchema,
	sort_order: z.coerce.number().int().min(1)
});

export const statusPageMonitorInputSchema = z.object({
	id: Int64String.optional(),
	monitor_id: z.coerce.number().int(),
	group_id: z.coerce.number().int().optional(),
	name: z.string().min(1).max(255),
	type: statusPageElementTypeSchema,
	sort_order: z.coerce.number().int().min(1)
});

export const statusPageUpsertRequestSchema = z.object({
	name: z.string().min(1).max(255),
	slug: z.string().min(3).max(255),

	groups: z.array(statusPageGroupInputSchema).default([]),
	monitors: z.array(statusPageMonitorInputSchema).default([])
});

export type StatusPageUpsertRequest = z.infer<typeof statusPageUpsertRequestSchema>;
export type StatusPageGroupInput = z.infer<typeof statusPageGroupInputSchema>;
export type StatusPageMonitorInput = z.infer<typeof statusPageMonitorInputSchema>;
export type StatusPageElementType = z.infer<typeof statusPageElementTypeSchema>;
