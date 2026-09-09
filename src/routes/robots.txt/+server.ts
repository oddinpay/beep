import type { RequestHandler } from "@sveltejs/kit";

export const GET: RequestHandler = async ({ url }) => {
  const fullHostname = `${url.protocol}//${url.host}`;
  const sitemapUrl = `${fullHostname}/sitemap.xml`;

  const body = ["User-agent: *", "", `Sitemap: ${sitemapUrl}`].join("\n");

  return new Response(body, {
    headers: {
      "Content-Type": "text/plain",
    },
  });
};
