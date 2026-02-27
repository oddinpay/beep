import { z } from 'zod/v4';

export const formSchema = z.object({
	navbar: z
		.string()
		.trim()
		.min(2, 'Navbar must be at least 2 characters long')
		.max(50, 'Navbar must not exceed 50 characters')
		.optional(),

	signup: z
		.string()
		.trim()
		.optional()
		.refine((val) => !val || z.url().safeParse(val).success, {
			message: 'Sign up URL must be a valid URL'
		}),

	signin: z
		.string()
		.trim()
		.optional()
		.refine((val) => !val || z.url().safeParse(val).success, {
			message: 'Sign in URL must be a valid URL'
		}),

	title: z
		.string()
		.trim()
		.min(2, 'Title must be at least 2 characters long')
		.max(50, 'Title must not exceed 50 characters')
		.optional(),

	description: z
		.string()
		.trim()
		.min(2, 'Description must be at least 2 characters long')
		.max(100, 'Description must not exceed 100 characters')
		.optional()
});
