import { query } from './_generated/server';

export const get = query({
	args: {},
	handler: async (ctx) => {
		const status = await ctx.db.query('status').collect();
		return status.map((status) => ({ ...status, assigner: 'oddin'}));
	}
});
