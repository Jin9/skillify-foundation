<!DOCTYPE html><html class="h-full antialiased" lang="en"><head><meta charSet="utf-8" data-next-head=""/><meta name="viewport" content="width=device-width" data-next-head=""/><title data-next-head="">Implementing CLAUDE.md and Agent Skills In Your Repository - Matthew Groff</title><meta name="description" content="A practical guide to the 3-tier documentation architecture that makes AI coding agents work: root CLAUDE.md, task-specific skills, and agent guides with progressive disclosure." data-next-head=""/><link rel="alternate" type="application/rss+xml" href="undefined/rss/feed.xml"/><link rel="alternate" type="application/feed+json" href="undefined/rss/feed.json"/><link rel="preload" href="/_next/static/css/5fa54e54744d7cde.css" as="style"/><link rel="preload" as="image" imageSrcSet="/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=16&amp;q=75 16w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=32&amp;q=75 32w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=48&amp;q=75 48w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=64&amp;q=75 64w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=96&amp;q=75 96w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=128&amp;q=75 128w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=256&amp;q=75 256w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=384&amp;q=75 384w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=640&amp;q=75 640w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=750&amp;q=75 750w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=828&amp;q=75 828w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1080&amp;q=75 1080w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1200&amp;q=75 1200w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1920&amp;q=75 1920w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=2048&amp;q=75 2048w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=3840&amp;q=75 3840w" imageSizes="2.25rem" data-next-head=""/><script>
  let darkModeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

  updateMode()
  darkModeMediaQuery.addEventListener('change', updateModeWithoutTransitions)
  window.addEventListener('storage', updateModeWithoutTransitions)

  function updateMode() {
    let isSystemDarkMode = darkModeMediaQuery.matches
    let isDarkMode = window.localStorage.isDarkMode === 'true' || (!('isDarkMode' in window.localStorage) && isSystemDarkMode)

    if (isDarkMode) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }

    if (isDarkMode === isSystemDarkMode) {
      delete window.localStorage.isDarkMode
    }
  }

  function disableTransitionsTemporarily() {
    document.documentElement.classList.add('[&_*]:!transition-none')
    window.setTimeout(() => {
      document.documentElement.classList.remove('[&_*]:!transition-none')
    }, 0)
  }

  function updateModeWithoutTransitions() {
    disableTransitionsTemporarily()
    updateMode()
  }
</script><link rel="stylesheet" href="/_next/static/css/5fa54e54744d7cde.css" data-n-g=""/><noscript data-n-css=""></noscript><script defer="" noModule="" src="/_next/static/chunks/polyfills-42372ed130431b0a.js"></script><script src="/_next/static/chunks/webpack-2da1769ff35eb1ed.js" defer=""></script><script src="/_next/static/chunks/framework-4a99af1472046e21.js" defer=""></script><script src="/_next/static/chunks/main-43f8e147a5e63233.js" defer=""></script><script src="/_next/static/chunks/pages/_app-1205d05b845c50e5.js" defer=""></script><script src="/_next/static/chunks/pages/blog/implementing-claude-md-agent-skills-7501dd9ab09e1ca8.js" defer=""></script><script src="/_next/static/GDxMezF22_T_gZTJytwox/_buildManifest.js" defer=""></script><script src="/_next/static/GDxMezF22_T_gZTJytwox/_ssgManifest.js" defer=""></script></head><body class="flex h-full flex-col bg-zinc-50 dark:bg-black"><link rel="preload" as="image" imageSrcSet="/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=16&amp;q=75 16w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=32&amp;q=75 32w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=48&amp;q=75 48w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=64&amp;q=75 64w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=96&amp;q=75 96w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=128&amp;q=75 128w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=256&amp;q=75 256w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=384&amp;q=75 384w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=640&amp;q=75 640w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=750&amp;q=75 750w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=828&amp;q=75 828w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1080&amp;q=75 1080w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1200&amp;q=75 1200w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1920&amp;q=75 1920w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=2048&amp;q=75 2048w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=3840&amp;q=75 3840w" imageSizes="2.25rem"/><link rel="preload" as="image" href="/images/photos/3TierArchitecture.png"/><div id="__next"><div class="fixed inset-0 flex justify-center sm:px-8"><div class="flex w-full max-w-7xl lg:px-8"><div class="w-full bg-white ring-1 ring-zinc-100 dark:bg-zinc-900 dark:ring-zinc-300/20"></div></div></div><div class="relative"><header class="pointer-events-none relative z-50 flex flex-col" style="height:var(--header-height);margin-bottom:var(--header-mb)"><div class="top-0 z-10 h-16 pt-6" style="position:var(--header-position)"><div class="sm:px-8 top-[var(--header-top,theme(spacing.6))] w-full" style="position:var(--header-inner-position)"><div class="mx-auto max-w-7xl lg:px-8"><div class="relative px-4 sm:px-8 lg:px-12"><div class="mx-auto max-w-2xl lg:max-w-5xl"><div class="relative flex gap-4"><div class="flex flex-1"><div class="h-10 w-10 rounded-full bg-white/90 p-0.5 shadow-lg shadow-zinc-800/5 ring-1 ring-zinc-900/5 backdrop-blur dark:bg-zinc-800/90 dark:ring-white/10"><a aria-label="Home" class="pointer-events-auto" href="/"><img alt="" width="512" height="512" decoding="async" data-nimg="1" class="rounded-full bg-zinc-100 object-cover dark:bg-zinc-800 h-9 w-9" style="color:transparent" sizes="2.25rem" srcSet="/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=16&amp;q=75 16w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=32&amp;q=75 32w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=48&amp;q=75 48w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=64&amp;q=75 64w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=96&amp;q=75 96w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=128&amp;q=75 128w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=256&amp;q=75 256w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=384&amp;q=75 384w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=640&amp;q=75 640w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=750&amp;q=75 750w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=828&amp;q=75 828w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1080&amp;q=75 1080w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1200&amp;q=75 1200w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=1920&amp;q=75 1920w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=2048&amp;q=75 2048w, /_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=3840&amp;q=75 3840w" src="/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Favatar.1fa39701.jpg&amp;w=3840&amp;q=75"/></a></div></div><div class="flex flex-1 justify-end md:justify-center"><div class="pointer-events-auto md:hidden" data-headlessui-state=""><button class="group flex items-center rounded-full bg-white/90 px-4 py-2 text-sm font-medium text-zinc-800 shadow-lg shadow-zinc-800/5 ring-1 ring-zinc-900/5 backdrop-blur dark:bg-zinc-800/90 dark:text-zinc-200 dark:ring-white/10 dark:hover:ring-white/20" type="button" aria-expanded="false" data-headlessui-state="">Menu<svg viewBox="0 0 8 6" aria-hidden="true" class="ml-3 h-auto w-2 stroke-zinc-500 group-hover:stroke-zinc-700 dark:group-hover:stroke-zinc-400"><path d="M1.75 1.75 4 4.25l2.25-2.5" fill="none" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg></button></div><span hidden="" style="position:fixed;top:1px;left:1px;width:1px;height:0;padding:0;margin:-1px;overflow:hidden;clip:rect(0, 0, 0, 0);white-space:nowrap;border-width:0;display:none"></span><nav class="pointer-events-auto hidden md:block"><ul class="flex rounded-full bg-white/90 px-3 text-sm font-medium text-zinc-800 shadow-lg shadow-zinc-800/5 ring-1 ring-zinc-900/5 backdrop-blur dark:bg-zinc-800/90 dark:text-zinc-200 dark:ring-white/10"><li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400" href="/blog">Blog</a></li><li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400" href="/projects">Projects</a></li><li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400" href="/uses">Uses</a></li><li><a class="relative block px-3 py-2 transition hover:text-teal-500 dark:hover:text-teal-400" href="/publications">Publications</a></li></ul></nav></div><div class="flex justify-end md:flex-1"><div class="pointer-events-auto"><button type="button" aria-label="Toggle dark mode" class="group rounded-full bg-white/90 px-3 py-2 shadow-lg shadow-zinc-800/5 ring-1 ring-zinc-900/5 backdrop-blur transition dark:bg-zinc-800/90 dark:ring-white/10 dark:hover:ring-white/20"><svg viewBox="0 0 24 24" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" class="h-6 w-6 fill-zinc-100 stroke-zinc-500 transition group-hover:fill-zinc-200 group-hover:stroke-zinc-700 dark:hidden [@media(prefers-color-scheme:dark)]:fill-teal-50 [@media(prefers-color-scheme:dark)]:stroke-teal-500 [@media(prefers-color-scheme:dark)]:group-hover:fill-teal-50 [@media(prefers-color-scheme:dark)]:group-hover:stroke-teal-600"><path d="M8 12.25A4.25 4.25 0 0 1 12.25 8v0a4.25 4.25 0 0 1 4.25 4.25v0a4.25 4.25 0 0 1-4.25 4.25v0A4.25 4.25 0 0 1 8 12.25v0Z"></path><path d="M12.25 3v1.5M21.5 12.25H20M18.791 18.791l-1.06-1.06M18.791 5.709l-1.06 1.06M12.25 20v1.5M4.5 12.25H3M6.77 6.77 5.709 5.709M6.77 17.73l-1.061 1.061" fill="none"></path></svg><svg viewBox="0 0 24 24" aria-hidden="true" class="hidden h-6 w-6 fill-zinc-700 stroke-zinc-500 transition dark:block [@media(prefers-color-scheme:dark)]:group-hover:stroke-zinc-400 [@media_not_(prefers-color-scheme:dark)]:fill-teal-400/10 [@media_not_(prefers-color-scheme:dark)]:stroke-teal-500"><path d="M17.25 16.22a6.937 6.937 0 0 1-9.47-9.47 7.451 7.451 0 1 0 9.47 9.47ZM12.75 7C17 7 17 2.75 17 2.75S17 7 21.25 7C17 7 17 11.25 17 11.25S17 7 12.75 7Z" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path></svg></button></div></div></div></div></div></div></div></div></header><main><div class="sm:px-8 mt-16 lg:mt-32"><div class="mx-auto max-w-7xl lg:px-8"><div class="relative px-4 sm:px-8 lg:px-12"><div class="mx-auto max-w-2xl lg:max-w-5xl"><div class="xl:relative"><div class="mx-auto max-w-2xl"><article><header class="flex flex-col"><h1 class="mt-6 text-4xl font-bold tracking-tight text-zinc-800 dark:text-zinc-100 sm:text-5xl">Implementing CLAUDE.md and Agent Skills In Your Repository</h1><time dateTime="2026-02-10" class="order-first flex items-center text-base text-zinc-400 dark:text-zinc-500"><span class="h-4 w-0.5 rounded-full bg-zinc-200 dark:bg-zinc-500"></span><span class="ml-3">February 10, 2026</span></time></header><div class="mt-8 prose dark:prose-invert"><h2>Your Rules File Is The Product</h2>
<p>In my <a href="/blog/early-2026-agentic-coding-update">Early 2026 Agentic Coding Update</a>, I talked about the loop: manager agents, coding agents, PRIME_DIRECTIVEs, and review gates.</p>
<p>That post covered how to run the loop. This one covers the thing that makes the loop actually work: the documentation architecture.</p>
<p>Every session with Claude Code or OpenCode starts stateless. No memory of your conventions. No knowledge of your test commands. No awareness that you use bun instead of npm or that migrations go through Alembic and not raw SQL.</p>
<p>The only thing that reliably onboards each new session is your rules file. <code>CLAUDE.md</code> for Claude Code. <code>AGENTS.md</code> for OpenCode. Whatever you put there shapes every decision the agent makes.</p>
<p>Most repos either have nothing, or they have a bloated auto-generated file that the model quietly ignores. Both fail for the same reason: the agent does not have the right context at the right time.</p>
<p>This post is the practical guide to fixing that.</p>
<h2>The Core Problem: Context At The Wrong Time</h2>
<p>The <a href="https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices">official Claude best practices for agent skills</a> make this clear: the context window is a shared public good. Your <code>CLAUDE.md</code> competes with the system prompt, conversation history, and every other piece of context the model needs.</p>
<p><a href="https://www.humanlayer.dev/blog/writing-a-good-claude-md">HumanLayer&#x27;s guide on writing a good CLAUDE.md</a> puts it more bluntly: Claude often deprioritizes <code>CLAUDE.md</code> content entirely. The system injects a reminder saying the context &quot;may or may not be relevant.&quot; The more instructions you include that are not universally applicable, the more likely the model dismisses the whole file.</p>
<p>This means you cannot solve the problem by writing a bigger rules file. You solve it with architecture.</p>
<h2>The 3-Tier Architecture</h2>
<p>After running this pattern across multiple production repos, I have landed on three tiers:</p>
<p><strong>Tier 1: Root CLAUDE.md</strong> - Universal rules that apply to every task. Under 100 lines. Loaded automatically every session.</p>
<p><strong>Tier 2: Skills (.claude/skills/name/SKILL.md)</strong> - Task-specific behavior loaded on demand. The agent reads a skill only when the task matches. Think of these as specialized playbooks: one for commits, one for PRs, one for migrations.</p>
<p><strong>Tier 3: Agent Guides (docs/agent-guides/name.md)</strong> - Deep reference material. Build commands, architecture docs, convention details. Skills point to these. The agent reads them only when it needs the full picture.</p>
<p>The key principle is progressive disclosure. The root file is a table of contents. Skills are chapters. Agent guides are appendices. The agent loads only what the current task requires.</p>
<p><img src="/images/photos/3TierArchitecture.png" alt="3-Tier Documentation Architecture"/></p>
<pre><code>your-repo/
  CLAUDE.md                          # Tier 1: Universal (&lt; 100 lines)
  .claude/
    skills/
      build-test-verify/SKILL.md     # Tier 2: On-demand task playbooks
      create-pull-request/SKILL.md
      git-commit/SKILL.md
      core-conventions/SKILL.md
      self-review-checklist/SKILL.md
  docs/
    agent-guides/
      build-test-verify.md           # Tier 3: Deep reference
      core-conventions.md
  backend/CLAUDE.md                  # Directory-level overrides
  frontend/CLAUDE.md
</code></pre>
<h2>Step 1: Audit What You Have</h2>
<p>Before writing anything, inventory what already exists. Most repos have some combination of:</p>
<ul>
<li>An existing <code>AGENTS.md</code> or <code>CLAUDE.md</code> (possibly auto-generated)</li>
<li>A <code>README.md</code> with build/test/run instructions</li>
<li>CI workflow files that encode the real validation commands</li>
<li>Scattered docs about architecture or conventions</li>
</ul>
<p>Do not throw any of this away. You are going to absorb it into the right tier.</p>
<p>Map each piece of existing documentation to where it belongs:</p>
<table><thead><tr><th>Existing Content</th><th>Target Tier</th></tr></thead><tbody><tr><td>&quot;Use bun, not npm&quot;</td><td>Tier 1 (root CLAUDE.md)</td></tr><tr><td>&quot;Run pytest -x for backend tests&quot;</td><td>Tier 2 (build-test-verify skill) or Tier 3 (agent guide)</td></tr><tr><td>&quot;Our agents follow this directory contract...&quot;</td><td>Tier 2 (domain-specific skill)</td></tr><tr><td>&quot;Here is the full project directory tree&quot;</td><td>Tier 3 (project-map agent guide)</td></tr><tr><td>&quot;Always use Conventional Commits&quot;</td><td>Tier 2 (git-commit skill)</td></tr></tbody></table>
<p>The rule: if it applies to every task, it goes in Tier 1. If it applies only when doing a specific kind of work, it goes in Tier 2 or 3.</p>
<h2>Step 2: Write the Root CLAUDE.md</h2>
<p>This is the highest-leverage file in your entire repo for AI-assisted development. Every token competes for attention.</p>
<p>Use the Why / What / How / Progressive Disclosure structure:</p>
<pre class="language-md"><code class="language-md"><span class="token title important"><span class="token punctuation">#</span> Project Name Agent Guide</span>

Use this file as the default onboarding context for this repo.

<span class="token title important"><span class="token punctuation">##</span> Why</span>

Brief description of what this project is and what matters.
One to three sentences. No filler.

<span class="token title important"><span class="token punctuation">##</span> What (project map)</span>

<span class="token list punctuation">-</span> <span class="token code-snippet code keyword">`frontend/`</span> - React app (Vite, Bun)
<span class="token list punctuation">-</span> <span class="token code-snippet code keyword">`backend/`</span> - FastAPI + Pydantic AI
<span class="token list punctuation">-</span> <span class="token code-snippet code keyword">`docs/agent-guides/`</span> - task-specific guidance loaded on demand

Read <span class="token code-snippet code keyword">`docs/agent-guides/project-map.md`</span> for the full map.

<span class="token title important"><span class="token punctuation">##</span> How (always apply)</span>

<span class="token list punctuation">-</span> Use existing project patterns before introducing new abstractions.
<span class="token list punctuation">-</span> Backend: use <span class="token code-snippet code keyword">`uv`</span> for Python deps, <span class="token code-snippet code keyword">`ruff`</span> for linting.
<span class="token list punctuation">-</span> Frontend: use <span class="token code-snippet code keyword">`bun`</span> for deps and scripts.
<span class="token list punctuation">-</span> Do not hardcode API URLs.
<span class="token list punctuation">-</span> Validate changes with the smallest relevant command set first.

<span class="token title important"><span class="token punctuation">##</span> Progressive Disclosure</span>

Do not load every guide for every task. Read only what is relevant:

<span class="token list punctuation">-</span> Build/test/lint: <span class="token code-snippet code keyword">`docs/agent-guides/build-test-verify.md`</span>
<span class="token list punctuation">-</span> Conventions: <span class="token code-snippet code keyword">`docs/agent-guides/core-conventions.md`</span>
<span class="token list punctuation">-</span> Migrations: <span class="token code-snippet code keyword">`docs/agent-guides/alembic-migrations.md`</span>

Use skills in <span class="token code-snippet code keyword">`.claude/skills/`</span> for task-specific behavior:

<span class="token list punctuation">-</span> Core: <span class="token code-snippet code keyword">`build-test-verify`</span>, <span class="token code-snippet code keyword">`core-conventions`</span>
<span class="token list punctuation">-</span> Workflow: <span class="token code-snippet code keyword">`create-pull-request`</span>, <span class="token code-snippet code keyword">`git-commit`</span>, <span class="token code-snippet code keyword">`self-review-checklist`</span>

<span class="token title important"><span class="token punctuation">##</span> Local Overrides</span>

Directory-level <span class="token code-snippet code keyword">`CLAUDE.md`</span> files may add stricter rules.
Apply the nearest file in addition to this one.

<span class="token title important"><span class="token punctuation">##</span> PR and Branching</span>

Use a feature branch. Never push directly to <span class="token code-snippet code keyword">`main`</span>.
Use the <span class="token code-snippet code keyword">`create-pull-request`</span> skill for standards.
</code></pre>
<p>Target: under 100 lines. HumanLayer keeps theirs under 60. Anthropic recommends under 300 but less is better. In my experience, 60 to 100 lines is the sweet spot for a real production repo.</p>
<p>What to leave out of Tier 1:</p>
<ul>
<li>Detailed build commands (belongs in a skill or agent guide)</li>
<li>Code style rules (use a linter, not the LLM)</li>
<li>Architecture deep dives (belongs in agent guides)</li>
<li>Anything that only applies to one directory (use local overrides)</li>
</ul>
<h2>Step 3: Write Your Skills</h2>
<p>Skills are the biggest upgrade over a flat <code>CLAUDE.md</code>. They give you task-specific behavior that loads only when relevant.</p>
<p>Every skill needs YAML frontmatter with <code>name</code> and <code>description</code>. The description is critical because the agent uses it to decide whether to load the skill.</p>
<pre class="language-yaml"><code class="language-yaml"><span class="token punctuation">---</span>
<span class="token key atrule">name</span><span class="token punctuation">:</span> build<span class="token punctuation">-</span>test<span class="token punctuation">-</span>verify
<span class="token key atrule">description</span><span class="token punctuation">:</span> Run lint<span class="token punctuation">,</span> test<span class="token punctuation">,</span> and build verification commands
  for the project. Use when validating changes<span class="token punctuation">,</span> running tests<span class="token punctuation">,</span>
  or checking builds.
<span class="token punctuation">---</span>
</code></pre>
<p>Start with these five skills. They cover the most common AI-assisted workflows in any repo:</p>
<ol>
<li><strong>build-test-verify</strong> - Your lint/test/build commands, organized by stack. The agent loads this whenever it needs to validate work.</li>
<li><strong>git-commit</strong> - Your commit message format, branch naming, what to check before committing. Prevents the agent from writing vague commit messages.</li>
<li><strong>create-pull-request</strong> - PR title format, description template, base branch, what checks must pass. Prevents sloppy PRs.</li>
<li><strong>core-conventions</strong> - Code style, import ordering, naming patterns, file organization. Only loaded when writing or reviewing code.</li>
<li><strong>self-review-checklist</strong> - Quality gate the agent runs before finishing work. Catches convention drift and missing tests.</li>
</ol>
<p>Then add domain-specific skills for what makes your repo unique. For example:</p>
<ul>
<li>A repo with database migrations needs an <code>alembic-migrations</code> or <code>drizzle-migrations</code> skill</li>
<li>A repo with AI agents needs a skill that codifies the agent creation contract</li>
<li>A repo with a complex auth system needs a skill for authorization patterns</li>
</ul>
<p>Skill authoring principles:</p>
<ul>
<li>Keep <code>SKILL.md</code> under 500 lines (Anthropic recommendation). Under 150 is better for most skills.</li>
<li>Point to agent guides for deep content. The skill says what to do. The agent guide explains how it all works.</li>
<li>Assume Claude is smart. Do not explain what Python is. Do not explain what a PR is. Only add context Claude genuinely lacks: your specific commands, your specific conventions, your specific architecture.</li>
<li>Set the right degree of freedom. Fragile operations (migrations, deployments) need exact commands. Flexible operations (code review, refactoring) need principles and heuristics.</li>
</ul>
<h2>Step 4: Write Agent Guides</h2>
<p>Agent guides live in <code>docs/agent-guides/</code> and serve as the deep reference layer. Skills point to them. The agent reads them only when it needs the full context.</p>
<p>Minimum set:</p>
<ul>
<li><strong>build-test-verify.md</strong> - Every lint, test, and build command with expected output. Include Docker and CI commands if relevant.</li>
<li><strong>core-conventions.md</strong> - The expanded version of your code conventions. Import patterns, naming rules, file organization, error handling patterns.</li>
</ul>
<p>Add more as complexity demands. Migration workflows, API contracts, auth patterns, streaming protocols. Whatever would take a new human engineer more than five minutes to figure out from the code alone.</p>
<p>Do not try to maintain a project-map or directory tree guide. These go stale immediately and the maintenance cost is not worth it. The agent can explore the filesystem directly when it needs to understand structure.</p>
<h2>Step 5: Add Directory-Level Overrides</h2>
<p>If your repo has distinct subsystems (like a <code>backend/</code> and <code>frontend/</code>), add a <code>CLAUDE.md</code> in each directory with rules that only apply there.</p>
<pre class="language-md"><code class="language-md"><span class="token title important"><span class="token punctuation">#</span> Backend Rules</span>

Apply the root <span class="token code-snippet code keyword">`CLAUDE.md`</span> first, then this file.

<span class="token title important"><span class="token punctuation">##</span> Non-Negotiable</span>

<span class="token list punctuation">-</span> All database changes go through Alembic migrations. No raw DDL.
<span class="token list punctuation">-</span> Use <span class="token code-snippet code keyword">`uv`</span> for all Python dependency management.
<span class="token list punctuation">-</span> All database operations use async SQLAlchemy sessions.
<span class="token list punctuation">-</span> Auth is MSAL-based. Do not introduce alternative auth patterns.
</code></pre>
<p>Keep these short. 20 to 30 lines. They exist to prevent mistakes in a specific area, not to repeat what is already in the root file.</p>
<h2>What Not To Do</h2>
<p><strong>Do not auto-generate your CLAUDE.md.</strong> Running <code>/init</code> produces generic output that wastes your highest-leverage file. Write it yourself.</p>
<p><strong>Do not use CLAUDE.md as a linter.</strong> HumanLayer is right about this: style rules in <code>CLAUDE.md</code> are expensive and unreliable. Use Ruff, Biome, ESLint, or Prettier. Set up a pre-commit hook or a Claude Code stop hook. LLMs are the wrong tool for deterministic formatting.</p>
<p><strong>Do not duplicate content across tiers.</strong> A skill should point to an agent guide, not copy its content. If the same instructions exist in two places, they will drift apart.</p>
<p><strong>Do not stuff everything into Tier 1.</strong> Every line in your root <code>CLAUDE.md</code> competes for attention. If it does not apply to literally every task, push it down to a skill or agent guide.</p>
<p><strong>Do not create skills for things that do not have enough complexity.</strong> A single CI workflow file does not need a <code>github-actions</code> skill. A simple UI library does not need a <code>frontend-design</code> skill. Skill count should match actual decision surface, not aspiration.</p>
<h2>Skills Are Powerful, But Not Automatically Safe</h2>
<p>One place I look for skill ideas is <a href="https://www.skills.sh/">skills.sh</a>, Vercel&#x27;s skills leaderboard. Discovery is useful. Discovery is not trust.</p>
<p>Skills are just markdown files, but they shape agent behavior in real ways. Treat every installed skill like executable influence over your workflow. A few failure modes are common:</p>
<p><strong>Skills can quietly bias recommendations toward products and services.</strong> A &quot;database optimization&quot; skill might steer every recommendation toward a specific managed service. You would not notice unless you read the raw markdown.</p>
<p><strong>Skills can contain dangerous instructions.</strong> Exfiltrating secrets with <code>curl</code>, writing to unexpected paths, or running destructive commands. The skill does not need to be malicious on purpose. A poorly scoped skill can still cause real damage.</p>
<p><strong>Skills can burn context when they are no longer useful.</strong> Install-heavy skills linger in your <code>.claude/skills/</code> directory and get loaded into sessions where they add nothing. Every loaded skill competes for context window space with the work you actually need done.</p>
<p><strong>Skills can conflict with your codebase conventions.</strong> A generic &quot;React best practices&quot; skill might recommend patterns that contradict your existing architecture. When skills and repo conventions disagree, the agent gets confused and output quality drops.</p>
<p>I have had the best results by taking inspiration from public skills, then rewriting them for my own repo standards and review gates.</p>
<h2>Safe Skill Adoption Checklist</h2>
<ol>
<li>Read the raw <code>SKILL.md</code> yourself before enabling it.</li>
<li>Ask an agent to audit the skill for security risks and exfiltration patterns.</li>
<li>Check that the skill matches your repo&#x27;s actual conventions and architecture.</li>
<li>Restrict permissions so the skill cannot overreach.</li>
<li>Remove or disable stale skills that no longer earn their context cost.</li>
</ol>
<h2>Cross-Tool Compatibility</h2>
<p>This architecture is not Claude-only. OpenCode reads <code>.agents/skills/</code> and <code>AGENTS.md</code> but also supports Claude-compatible fallbacks. I keep <code>.claude/skills/</code> as the canonical location and it works across both tools.</p>
<p>If you use OpenCode, you can add a compatibility note to your skill frontmatter:</p>
<pre class="language-yaml"><code class="language-yaml"><span class="token punctuation">---</span>
<span class="token key atrule">name</span><span class="token punctuation">:</span> build<span class="token punctuation">-</span>test<span class="token punctuation">-</span>verify
<span class="token key atrule">description</span><span class="token punctuation">:</span> Run lint<span class="token punctuation">,</span> test<span class="token punctuation">,</span> and build verification.
<span class="token key atrule">compatibility</span><span class="token punctuation">:</span> claude<span class="token punctuation">-</span>code<span class="token punctuation">,</span> opencode
<span class="token punctuation">---</span>
</code></pre>
<p>Other tools like Cursor and Windsurf have their own rules file conventions. The agent guide layer (<code>docs/agent-guides/</code>) works with anything because it is just markdown in a standard location.</p>
<h2>The Payoff</h2>
<p>After implementing this across multiple repos, the difference is measurable:</p>
<p><strong>Before:</strong> The agent guesses at test commands, uses the wrong package manager, creates PRs with no description format, writes commit messages that say &quot;update code.&quot;</p>
<p><strong>After:</strong> The agent runs the right commands, follows your conventions, creates PRs that match your team&#x27;s standards, and loads deep context only when the task actually needs it.</p>
<p>The root <code>CLAUDE.md</code> is the highest-leverage file in your repo for AI-assisted development. The skills and agent guides turn one good file into a system. Treat it like infrastructure, not documentation.</p>
<h2>Getting Started Checklist</h2>
<p>If you want to add this to your repo today:</p>
<ol>
<li><strong>Audit</strong> - Inventory your existing docs, README, CI files, and any <code>AGENTS.md</code></li>
<li><strong>Root CLAUDE.md</strong> - Write the Why / What / How / Progressive Disclosure structure. Under 100 lines.</li>
<li><strong>Five starter skills</strong> - <code>build-test-verify</code>, <code>git-commit</code>, <code>create-pull-request</code>, <code>core-conventions</code>, <code>self-review-checklist</code></li>
<li><strong>Two starter agent guides</strong> - <code>build-test-verify</code>, <code>core-conventions</code></li>
<li><strong>Directory overrides</strong> - Add <code>CLAUDE.md</code> in directories with unique constraints</li>
<li><strong>Test it</strong> - Start a fresh Claude Code or OpenCode session and try a real task. Watch what the agent loads and where it stumbles. Iterate.</li>
</ol>
<p>Do not try to get it perfect on the first pass. The best <code>CLAUDE.md</code> files are iterated over weeks based on watching real agent behavior.</p>
<p>The best source of improvements is code review. Every comment a reviewer leaves on an AI-assisted PR is a signal that the agent lacked context. Wrong import pattern? Add it to <code>core-conventions</code>. Missed a test command? Update <code>build-test-verify</code>. Used the wrong migration tool? Add a line to the directory override. I treat every review comment as an opportunity to add a new rule, a new checklist item, a new guardrail. Over time, the same mistakes stop appearing because the agent has the context it was missing. These files are living documents, not write-once artifacts.</p>
<h2>Further Reading</h2>
<ul>
<li><a href="https://platform.claude.com/docs/en/agents-and-tools/agent-skills/best-practices">Claude Agent Skills Best Practices</a> - The official guide on writing effective skills. Covers progressive disclosure, YAML frontmatter, and the checklist.</li>
<li><a href="https://www.humanlayer.dev/blog/writing-a-good-claude-md">Writing a Good CLAUDE.md (HumanLayer)</a> - Why less is more, and why your <code>CLAUDE.md</code> should not be a linter.</li>
<li><a href="/blog/early-2026-agentic-coding-update">Early 2026 Agentic Coding Update</a> - My broader agentic workflow post covering the manager/coding agent loop, PRIME_DIRECTIVEs, and skills safety.</li>
</ul></div></article></div></div></div></div></div></div></main><footer class="mt-32"><div class="sm:px-8"><div class="mx-auto max-w-7xl lg:px-8"><div class="border-t border-zinc-100 pb-16 pt-10 dark:border-zinc-700/40"><div class="relative px-4 sm:px-8 lg:px-12"><div class="mx-auto max-w-2xl lg:max-w-5xl"><div class="flex flex-col items-center justify-between gap-6 sm:flex-row"><div class="flex flex-wrap justify-center gap-x-6 gap-y-1 text-sm font-medium text-zinc-800 dark:text-zinc-200"><a class="transition hover:text-teal-500 dark:hover:text-teal-400" href="/blog">Blog</a><a class="transition hover:text-teal-500 dark:hover:text-teal-400" href="/projects">Projects</a><a class="transition hover:text-teal-500 dark:hover:text-teal-400" href="/uses">Uses</a><a class="transition hover:text-teal-500 dark:hover:text-teal-400" href="/publications">Publications</a></div><p class="cursor-pointer text-sm text-zinc-400 dark:text-zinc-500"><a target="_blank" href="https://icons8.com/icon/91234/home">Home</a> <!-- -->icon by<!-- --> <a target="_blank" href="https://icons8.com">Icons8</a></p><p class="text-sm text-zinc-400 dark:text-zinc-500">© <!-- -->2026<!-- --> Matthew Groff. All rights reserved.</p></div></div></div></div></div></div></footer></div></div><script id="__NEXT_DATA__" type="application/json">{"props":{"pageProps":{}},"page":"/blog/implementing-claude-md-agent-skills","query":{},"buildId":"GDxMezF22_T_gZTJytwox","nextExport":true,"autoExport":true,"isFallback":false,"scriptLoader":[]}</script></body></html>