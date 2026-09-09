export const GET = async ({ url }) => {
  const host = url.hostname;
  const sitemapUrl = `${host}/sitemap.xml`;

  const body = ["User-agent: *", "", `Sitemap: ${sitemapUrl}`].join("\n");

  return new Response(body, {
    headers: {
      "Content-Type": "text/plain",
    },
  });
};
