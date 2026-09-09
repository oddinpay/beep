// /src/routes/sitemap.xml/+server.ts
import { env } from "$env/dynamic/private";
import * as sitemap from "super-sitemap/sveltekit";
import type { RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ url }) => {
  const fullHostname = url.hostname;
  const domain = fullHostname.split(".").slice(-2).join(".");

  return await sitemap.response({
    origin: "https://status." + domain,
    excludeRoutePatterns: [/^\/unsubscribe.*/],
  });
};
