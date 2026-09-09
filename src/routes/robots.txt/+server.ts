export const GET = async ({ url }) => {
  const fullHostname = url.hostname;
  const domain = fullHostname.split(".").slice(-2).join(".");
  const sitemapUrl = `status.${domain}/sitemap.xml`;

  const body = ["User-agent: *", "", `Sitemap: ${sitemapUrl}`].join("\n");

  return new Response(body, {
    headers: {
      "Content-Type": "text/plain",
    },
  });
};
