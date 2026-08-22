<script lang="ts">
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import { theme } from "$lib/theme";
  import { onMount } from "svelte";
  import "../app.css";

  const schemaWebSite = {
    "@context": "https://schema.org",
    "@type": "WebSite",
    name: "The Url Mentalist",
    alternateName: ["Url Mentalist", "T-H-E U-R-L M-E-N-T-A-L-I-S-T"],
    url: "Live_deployment_link",
    description:
      "Free real-time URL scanner. Paste any link and get an instant phishing and safety verdict, no signup needed.",
    potentialAction: {
      "@type": "SearchAction",
      target: {
        "@type": "EntryPoint",
        urlTemplate: "Live_deployment_link/?q={search_term_string}",
      },
      "query-input": "required name=search_term_string",
    },
  };

  const schemaOrganization = {
    "@context": "https://schema.org",
    "@type": "Organization",
    name: "The Url Mentalist",
    alternateName: ["Url Mentalist", "T-H-E U-R-L M-E-N-T-A-L-I-S-T"],
    url: "Live_deployment_link",
    logo: "Live_deployment_link/favicon.ico",
    sameAs: ["https://github.com/AARYA-001/The-Url-Mentalist"],
  };

  onMount(() => {
    theme.init();
  });
</script>

<svelte:head>
  {@html `<script type="application/ld+json">${JSON.stringify(schemaWebSite)}</script>`}
  {@html `<script type="application/ld+json">${JSON.stringify(schemaOrganization)}</script>`}
</svelte:head>

<div class="relative min-h-screen flex flex-col">
  <ThemeToggle />
  <main class="flex-1">
    <slot />
  </main>

  <footer
    class="bg-gray-100 dark:bg-gray-900 text-gray-600 dark:text-gray-300 border-t border-gray-200 dark:border-gray-800 py-6"
  >
    <div
      class="max-w-5xl mx-auto px-6 flex flex-col md:flex-row md:justify-between items-center gap-4 text-base"
    >
      <!-- Navigation links -->
      <div class="flex flex-wrap items-center justify-center md:justify-start gap-4">
        <a
          href="/how-it-works"
          class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
        >
          How it works
        </a>
        <span class="text-gray-300 dark:text-gray-700 hidden md:inline">·</span>
        <a
          href="/about"
          class="text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
        >
          About
        </a>
      </div>
    </div>
  </footer>
</div>
