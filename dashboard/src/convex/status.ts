import { query } from './_generated/server';
import { v } from 'convex/values';

export const get = query({
	args: {
		apiKey: v.string()
	},
	handler: async (ctx, args) => {
		if (args.apiKey !== process.env.API_KEY) {
			throw new Error('Unauthorized: Wrong API Key');
		}

		const status = await ctx.db.query('status').collect();
		return status.map((status) => ({ ...status, assigner: 'oddin' }));
	}
});
