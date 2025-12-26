import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
// import path from 'node:path';
// import { fileURLToPath } from 'node:url';

// const dirname = path.resolve(fileURLToPath(import.meta.url), '../');

const adapterOptions = { precompress: true };

const config = {
	extensions: [".svelte", ".svx", ".md"],
	preprocess: [vitePreprocess()],

	kit: {
		adapter: adapter(adapterOptions),
		alias: {
			$assets: './src/lib/assets',
			$data: './src/lib/data',
			$helpers: './src/lib/helpers',
			$lib: './src/lib'
		}
	}
};

export default config;
