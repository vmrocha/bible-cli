import { defineConfig } from "vitepress";

export default defineConfig({
  lang: "en-US",
  title: "Bible CLI",
  description: "Read and search the Bible from your terminal, offline.",
  base: "/bible-cli/",
  cleanUrls: true,
  themeConfig: {
    nav: [
      { text: "Guide", link: "/guide/installation" },
      { text: "Reference", link: "/reference/commands" },
      { text: "GitHub", link: "https://github.com/vmrocha/bible-cli" }
    ],
    sidebar: {
      "/guide/": [
        {
          text: "Guide",
          items: [
            { text: "Installation", link: "/guide/installation" },
            { text: "Reading and search", link: "/guide/reading-and-search" },
            { text: "Configuration", link: "/guide/configuration" }
          ]
        }
      ],
      "/reference/": [
        {
          text: "Reference",
          items: [{ text: "Commands", link: "/reference/commands" }]
        }
      ]
    },
    socialLinks: [{ icon: "github", link: "https://github.com/vmrocha/bible-cli" }],
    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026 Bible CLI contributors"
    },
    search: { provider: "local" }
  }
});
