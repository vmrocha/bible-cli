# ADR 0003: Publish the project website with VitePress and GitHub Pages

- Status: proposed
- Date: 2026-08-03

## Context

Bible CLI needs a public project website that provides a clear first impression
for prospective users and a durable home for installation, usage, command
reference, configuration, and release documentation. The site should be simple
for Go contributors to maintain: documentation changes should be ordinary
Markdown pull requests reviewed alongside the code they describe.

The project does not need server-side rendering, a database, accounts, or other
dynamic application infrastructure. Maintaining a separate hosting system or
deploying generated site files to a branch would add operational work without
helping users.

## Decision

Use [VitePress](https://vitepress.dev/) to build a static project website from
Markdown files stored in this repository. The initial site will include:

- a landing page explaining Bible CLI, its offline-first approach, and supported
  platforms;
- installation and quick-start guides;
- guides for reading, searching, configuration, shell completion, and Bible
  translation attribution; and
- a command reference generated or maintained from the CLI's documented command
  surface.

Host the generated site on GitHub Pages. Start at the repository Pages URL
(`https://vmrocha.github.io/bible-cli/`) and attach a project-owned custom
domain later, when one is available. Require HTTPS for the custom domain.

Keep VitePress source and configuration under `website/`. Configure its public
base path as `/bible-cli/` while using the GitHub Pages project URL; change it
to `/` when moving to a custom domain.

Deploy with a GitHub Actions workflow that runs on pushes to `main` affecting
the website or its workflow, and supports manual dispatch. The workflow will:

1. check out the source;
2. install the pinned Node.js and package-manager versions;
3. install dependencies using the lockfile;
4. build the static site;
5. upload the generated artifact using `actions/upload-pages-artifact`; and
6. publish it using `actions/deploy-pages` with only the `pages: write` and
   `id-token: write` permissions required for deployment.

Pull requests should run the website build as a validation check. Deployment is
limited to `main`, so a documentation pull request cannot publish production
content before review and merge.

## Consequences

- Website content, its deployment workflow, and application code stay in one
  repository and follow the same review, ownership, and history.
- Contributors primarily write Markdown; VitePress supplies navigation, search,
  responsive documentation layouts, syntax highlighting, and a small amount of
  optional Vue-based customization for the landing page.
- The production site is static, fast, inexpensive, and has no application
  server or hosting credentials to operate beyond GitHub Pages configuration.
- The project adds a small Node.js toolchain and lockfile alongside its Go
  toolchain. The website build must be kept reproducible and updated through
  normal dependency-review practices.
- GitHub Pages is appropriate while requirements remain static. Reconsider
  Cloudflare Pages if the project later needs edge controls, advanced redirects
  and headers, stronger preview-deployment controls, or broader hosting
  requirements.

## Alternatives considered

### Docusaurus

Docusaurus is a capable documentation platform with an integrated landing-page
and blog model. It is a good choice for a large documentation program, but it
introduces a heavier React-based toolchain than Bible CLI needs initially.

### Hugo

Hugo aligns with the project's Go ecosystem and produces excellent static
sites. VitePress is preferred for its focused documentation experience and
straightforward Markdown authoring, without making contributors learn Go
templates for ordinary documentation changes.

### Cloudflare Pages

Cloudflare Pages is a strong static-hosting option, particularly when advanced
edge configuration or preview environments are important. GitHub Pages is
preferred initially because source, review, CI, and publishing stay in the
repository's existing GitHub workflow.

### Generated files committed to a `gh-pages` branch

Do not commit build output. Publishing an Actions artifact keeps generated
files out of review history and makes each deployment traceable to its source
commit and workflow run.
