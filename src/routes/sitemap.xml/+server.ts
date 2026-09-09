// /src/routes/sitemap.xml/+server.ts
import * as sitemap from "super-sitemap/sveltekit";
import type { RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ url }) => {
  const fullHostname = `${url.protocol}//${url.host}`;
  return await sitemap.response({
    origin: fullHostname,
    excludeRoutePatterns: [/^\/unsubscribe.*/],
  });
};
