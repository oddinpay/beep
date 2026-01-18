import { mdsvex } from "mdsvex";
import cf from "@sveltejs/adapter-cloudflare";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";
// import path from "node:path";
// import { fileURLToPath } from "node:url";

// const dirname = path.resolve(fileURLToPath(import.meta.url), "../");

/** @type {import('@sveltejs/kit').Config} */

const adapterOptions = { precompress: true };

const config = {
  // Consult https://svelte.dev/docs/kit/integrations
  // for more information about preprocessors

  extensions: [".svelte", ".svx", ".md"],
  preprocess: [
    vitePreprocess(),
    mdsvex({
      extensions: [".svx", ".md"],
      // layout: {
      //   status: path.join(
      //     dirname,
      //     "./src/lib/components/layout/_status.svelte"
      //   ),
      // },
    }),
  ],
  kit: {
    // adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
    // If your environment is not supported, or you settled on a specific environment, switch out the adapter.
    // See https://svelte.dev/docs/kit/adapters for more information about adapters.

    adapter: cf(adapterOptions),

    alias: {
      $assets: "./src/lib/assets",
      $data: "./src/lib/data",
      $helpers: "./src/lib/helpers",
      $lib: "./src/lib",
    },
  },
};

export default config;
