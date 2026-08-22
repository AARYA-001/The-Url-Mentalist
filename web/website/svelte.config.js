import adapterAuto from '@sveltejs/adapter-auto';
import adapterStatic from '@sveltejs/adapter-static';
import adapterVercel from '@sveltejs/adapter-vercel';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),

  kit: {
    // Docker builds set DOCKER_BUILD=1 → adapter-static (served by Caddy).
    // Vercel builds → adapter-vercel.
    // Local dev & fallback → adapter-auto.
    adapter:
      process.env.DOCKER_BUILD === '1'
        ? adapterStatic({ fallback: '200.html' })
        : process.env.VERCEL
          ? adapterVercel()
          : adapterAuto(),
  },
};

export default config;

