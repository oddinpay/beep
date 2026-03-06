import { query } from './_generated/server';

export const get = query({
	args: {},
	handler: async (ctx, args) => {
		const identity = await ctx.auth.getUserIdentity();
		if (identity === null) {
			throw new Error('Unauthenticated call to mutation');
		}
		const status = await ctx.db.query('status').collect();
		return status.map((status) => ({ ...status, assigner: 'oddin' }));
	}
});
