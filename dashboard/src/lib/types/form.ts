import { z } from 'zod/v4';

export const formSchema = z.object({
	logo: z
		.string()
		.trim()
		.superRefine((val, ctx) => {
			if (!val) {
				ctx.addIssue({ code: 'custom', message: 'Logo is required' });
				return;
			}
			if (val.length < 2) {
				ctx.addIssue({ code: 'custom', message: 'Logo must be at least 2 characters long' });
				return;
			}
			if (val.length > 50) {
				ctx.addIssue({ code: 'custom', message: 'Logo must not exceed 50 characters' });
				return;
			}
		}),
	title: z
		.string()
		.trim()
		.superRefine((val, ctx) => {
			if (!val) {
				ctx.addIssue({ code: 'custom', message: 'Title is required' });
				return;
			}
			if (val.length < 2) {
				ctx.addIssue({ code: 'custom', message: 'Title must be at least 2 characters long' });
				return;
			}
			if (val.length > 50) {
				ctx.addIssue({ code: 'custom', message: 'Title must not exceed 50 characters' });
				return;
			}
		}),
	description: z
		.string()
		.trim()
		.superRefine((val, ctx) => {
			if (!val) {
				ctx.addIssue({ code: 'custom', message: 'Description is required' });
				return;
			}
			if (val.length < 2) {
				ctx.addIssue({
					code: 'custom',
					message: 'Description must be at least 2 characters long'
				});
				return;
			}
			if (val.length > 100) {
				ctx.addIssue({ code: 'custom', message: 'Description must not exceed 100 characters' });
				return;
			}
		})
});
