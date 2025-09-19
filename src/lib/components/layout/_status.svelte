<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button } from "$lib/components/ui/button/index.js"; 
  import Buttong from "$lib/components/Buttong.svelte";
  import { Label } from "$lib/components/ui/label";
  import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
  import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
  import { MediaQuery } from "svelte/reactivity";
  import * as Pagination from "$lib/components/ui/pagination/index.js";
  import { page } from "$app/state";

  const isDesktop = new MediaQuery("(min-width: 768px)");

  type AccordionItem = {
    value: string;
    date: string;
    title: string;
    content: string;
    PageTitle: string;
    cover: string;
    active: boolean;
    features: string[];
  };

  // --- types ---
  type AccordionItemNoVal = Omit<AccordionItem, "value">;

  interface RoadmapProps {
    sections?: AccordionItemNoVal[];
    badge?: string;
    status?: string;
    logo?: string;
    slug?: string;
    cover?: string;
    title?: string;
    description?: string;
    features?: string[];
  }

  const allProps = $props() as RoadmapProps & Record<string, unknown>;

  const {
    sections: inputSections,
    badge,
    status,
    logo,
    features,
    cover,
    slug,
    title,
    description,
    ..._
  } = allProps;

  const count = (inputSections ?? []).length;

  console.log("Count:", count);
  const perPage = $derived(isDesktop.current ? 3 : 3);
  const siblingCount = $derived(isDesktop.current ? 1 : 1);
  const totalPages = $derived(Math.ceil(count / perPage));

  function getValidPage() {
    const param = page.url.searchParams.get("page");
    const parsed = Number(param);

    // Fallback logic: if not a valid number or out of range → 1
    if (!param || isNaN(parsed) || parsed < 1 || parsed > totalPages) {
      return 1;
    }
    return parsed;
  }

  const currentPage = $derived(getValidPage());

  function setPage(p: number) {
    const pathname = page.url.pathname;
    goto(`${pathname}?page=${p}`, { replaceState: true });
  }

  $effect(() => {
    console.log("Total pages:", totalPages);
    console.log("Current page:", currentPage);
  });

  // --- styles ---
  const COLOR_STYLES: Record<
    string,
    {
      badgeBg: string;
      badgeText: string;
      dividerBg: string;
      numberText: string;
      hoverText: string;
      caretText: string;
    }
  > = {
    orange: {
      badgeBg: "bg-orange-100",
      badgeText: "text-orange-700",
      dividerBg: "via-orange-400",
      numberText: "text-orange-700",
      hoverText: "hover:text-orange-700",
      caretText: "text-orange-700",
    },
    cyan: {
      badgeBg: "bg-cyan-100",
      badgeText: "text-cyan-700",
      dividerBg: "via-cyan-400",
      numberText: "text-cyan-700",
      hoverText: "hover:text-cyan-700",
      caretText: "text-cyan-700",
    },
    lime: {
      badgeBg: "bg-lime-100",
      badgeText: "text-lime-700",
      dividerBg: "via-lime-400",
      numberText: "text-lime-700",
      hoverText: "hover:text-lime-700",
      caretText: "text-lime-700",
    },

    gray: {
      badgeBg: "bg-gray-100",
      badgeText: "text-gray-700",
      dividerBg: "via-gray-400",
      numberText: "text-gray-700",
      hoverText: "hover:text-gray-700",
      caretText: "text-gray-700",
    },
    red: {
      badgeBg: "bg-red-100",
      badgeText: "text-red-700",
      dividerBg: "via-red-400",
      numberText: "text-red-700",
      hoverText: "hover:text-red-700",
      caretText: "text-red-700",
    },
    amber: {
      badgeBg: "bg-amber-100",
      badgeText: "text-amber-700",
      dividerBg: "via-amber-400",
      numberText: "text-amber-700",
      hoverText: "hover:text-amber-700",
      caretText: "text-amber-700",
    },
    yellow: {
      badgeBg: "bg-yellow-100",
      badgeText: "text-yellow-700",
      dividerBg: "via-yellow-400",
      numberText: "text-yellow-700",
      hoverText: "hover:text-yellow-700",
      caretText: "text-yellow-700",
    },
    green: {
      badgeBg: "bg-green-100",
      badgeText: "text-green-700",
      dividerBg: "via-green-400",
      numberText: "text-green-700",
      hoverText: "hover:text-green-700",
      caretText: "text-green-700",
    },
    emerald: {
      badgeBg: "bg-emerald-100",
      badgeText: "text-emerald-700",
      dividerBg: "via-emerald-400",
      numberText: "text-emerald-700",
      hoverText: "hover:text-emerald-700",
      caretText: "text-emerald-700",
    },
    teal: {
      badgeBg: "bg-teal-100",
      badgeText: "text-teal-700",
      dividerBg: "via-teal-400",
      numberText: "text-teal-700",
      hoverText: "hover:text-teal-700",
      caretText: "text-teal-700",
    },
    blue: {
      badgeBg: "bg-blue-100",
      badgeText: "text-blue-700",
      dividerBg: "via-blue-400",
      numberText: "text-blue-700",
      hoverText: "hover:text-blue-700",
      caretText: "text-blue-700",
    },
    indigo: {
      badgeBg: "bg-indigo-100",
      badgeText: "text-indigo-700",
      dividerBg: "via-indigo-400",
      numberText: "text-indigo-700",
      hoverText: "hover:text-indigo-700",
      caretText: "text-indigo-700",
    },
    purple: {
      badgeBg: "bg-purple-100",
      badgeText: "text-purple-700",
      dividerBg: "via-purple-400",
      numberText: "text-purple-700",
      hoverText: "hover:text-purple-700",
      caretText: "text-purple-700",
    },
    pink: {
      badgeBg: "bg-pink-100",
      badgeText: "text-pink-700",
      dividerBg: "via-pink-400",
      numberText: "text-pink-700",
      hoverText: "hover:text-pink-700",
      caretText: "text-pink-700",
    },
    rose: {
      badgeBg: "bg-rose-100",
      badgeText: "text-rose-700",
      dividerBg: "via-rose-400",
      numberText: "text-rose-700",
      hoverText: "hover:text-rose-700",
      caretText: "text-rose-700",
    },

    black: {
      badgeBg: "bg-black",
      badgeText: "text-white",
      dividerBg: "via-black",
      numberText: "text-black",
      hoverText: "hover:text-white",
      caretText: "text-white",
    },
    white: {
      badgeBg: "bg-white",
      badgeText: "text-gray-800",
      dividerBg: "via-gray-200",
      numberText: "text-gray-800",
      hoverText: "hover:text-gray-800",
      caretText: "text-gray-800",
    },
    zinc: {
      badgeBg: "bg-zinc-100",
      badgeText: "text-zinc-700",
      dividerBg: "via-zinc-400",
      numberText: "text-zinc-700",
      hoverText: "hover:text-zinc-700",
      caretText: "text-zinc-700",
    },
    slate: {
      badgeBg: "bg-slate-100",
      badgeText: "text-slate-700",
      dividerBg: "via-slate-400",
      numberText: "text-slate-700",
      hoverText: "hover:text-slate-700",
      caretText: "text-slate-700",
    },
    stone: {
      badgeBg: "bg-stone-100",
      badgeText: "text-stone-700",
      dividerBg: "via-stone-400",
      numberText: "text-stone-700",
      hoverText: "hover:text-stone-700",
      caretText: "text-stone-700",
    },
    sky: {
      badgeBg: "bg-sky-100",
      badgeText: "text-sky-700",
      dividerBg: "via-sky-400",
      numberText: "text-sky-700",
      hoverText: "hover:text-sky-700",
      caretText: "text-sky-700",
    },
    violet: {
      badgeBg: "bg-violet-100",
      badgeText: "text-violet-700",
      dividerBg: "via-violet-400",
      numberText: "text-violet-700",
      hoverText: "hover:text-violet-700",
      caretText: "text-violet-700",
    },

    fuchsia: {
      badgeBg: "bg-fuchsia-100",
      badgeText: "text-fuchsia-700",
      dividerBg: "via-fuchsia-400",
      numberText: "text-fuchsia-700",
      hoverText: "hover:text-fuchsia-700",
      caretText: "text-fuchsia-700",
    },
    neutral: {
      badgeBg: "bg-neutral-100",
      badgeText: "text-neutral-700",
      dividerBg: "via-neutral-400",
      numberText: "text-neutral-700",
      hoverText: "hover:text-neutral-700",
      caretText: "text-neutral-700",
    },
  };

  function makeBase(color: string) {
    return COLOR_STYLES[color] ?? COLOR_STYLES.gray; // fallback to gray
  }
  const sections = $derived(
    Array.isArray(inputSections) && inputSections.length
      ? inputSections.map((s: any, i: number) => {
          const key = `${s.color ?? "section"}-${i}`;

          // Build Tailwind classes dynamically from markdown color
          const style = makeBase(s.color);

          const rawItems: AccordionItemNoVal[] = Array.isArray(s.items)
            ? s.items
            : [
                {
                  date: s.date ?? "",
                  color: s.color ?? "",
                  title: s.title ?? s.PageTitle ?? s.heading ?? "",
                  content: s.content ?? "",
                  PageTitle: s.PageTitle ?? s.title ?? "",
                  active: s.active ?? false,
                  features: s.features ?? [],
                  cover: s.cover ?? s.covers?.[0] ?? "",
                },
              ].filter(
                (item) =>
                  item.title || item.PageTitle || item.content || item.date,
              );

          return {
            key,
            items: rawItems,
            ...style,
          };
        })
      : [],
  );

  const highlightPhrases = [
    { phrase: "Learn more", href: "/docs", class: "highlight-indigo-500" }, 
  ];

  function escapeRegex(s: string) {
    return s.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  }

  function normalizeHref(url: string) {
    if (!url) return url;
    if (url.startsWith("/")) return url;
    return url.startsWith("http://") || url.startsWith("https://")
      ? url
      : `https://${url}`;
  }

  type Part =
    | { type: "text"; text: string }
    | {
        type: "link";
        text: string;
        href: string;
        class?: string;
        external?: boolean;
      };

  function highlightInPlainParts(
    text: string,
    phrases: { phrase: string; href: string; class: string }[],
  ): Part[] {
    if (!text) return [];
    const sorted = [...phrases].sort(
      (a, b) => b.phrase.length - a.phrase.length,
    );
    const pattern = sorted.map((p) => escapeRegex(p.phrase)).join("|");
    if (!pattern) return [{ type: "text", text }];

    const re = new RegExp(pattern, "gi");
    let lastIndex = 0;
    let result: Part[] = [];
    let m: RegExpExecArray | null;
    while ((m = re.exec(text)) !== null) {
      if (m.index > lastIndex) {
        result.push({ type: "text", text: text.slice(lastIndex, m.index) });
      }
      const matched = m[0];
      const cfg = sorted.find(
        (p) => p.phrase.toLowerCase() === matched.toLowerCase(),
      );
      const href = cfg ? cfg.href : "/pricing";
      const normalized = normalizeHref(href);
      result.push({
        type: "link",
        text: matched,
        href: normalized,
        class: cfg ? cfg.class : "highlight-indigo",
        external: normalized.startsWith("http"),
      });
      lastIndex = m.index + matched.length;
    }
    if (lastIndex < text.length) {
      result.push({ type: "text", text: text.slice(lastIndex) });
    }
    return result;
  }

  function renderContentParts(
    text: string,
    phrases: { phrase: string; href: string; class: string }[],
  ): Part[] {
    if (!text) return [];
    const urlRe = /(https?:\/\/[^\s<]+|www\.[^\s<]+|\/[^\s<]+)/g;
    let lastIndex = 0;
    let result: Part[] = [];
    let match: RegExpExecArray | null;

    while ((match = urlRe.exec(text)) !== null) {
      const before = text.slice(lastIndex, match.index);
      result.push(...highlightInPlainParts(before, phrases));

      const rawUrl = match[0];
      const href = normalizeHref(rawUrl);
      result.push({
        type: "link",
        text: rawUrl,
        href,
        class: "link-underline",
        external: href.startsWith("http"),
      });

      lastIndex = match.index + rawUrl.length;
    }

    const tail = text.slice(lastIndex);
    result.push(...highlightInPlainParts(tail, phrases));

    return result;
  }

   function setOverflow(_node: HTMLElement) {
    if (window.matchMedia && window.matchMedia("(min-width: 1279px)").matches) {
      document.documentElement.style.overflow = "unset";
    }
  }
</script>

<svelte:head>
  <title>{title}</title>
  <meta name="description" content={description} />

  <style>
    #navToggle {
      display: none !important;
    }

    #content h1 {
      font-size: 2.2em;
      line-height: 1.5;
      font-weight: 600;
    }
    #content h2 {
      font-size: 1.6em;
      line-height: 1.5;
      font-weight: 600;
    }
    #content h3 {
      font-size: 1.2em;
      line-height: 1.5;
      font-weight: 600;
    }

    #content p {
      font-size: 1.1em;
      line-height: 1.75;
    }

    #content h1 {
      margin-bottom: 24px;
    }

    #content h2 {
      margin-bottom: 24px;
      margin-top: 24px;
    }

    #content h3 {
      margin-bottom: 16px;
      margin-top: 32px;
    }

    @media (min-width: 768px) {
      #content h1 {
        font-size: 2.8em;
        line-height: 1.5;
        font-weight: bold;
      }

      #content h2 {
        font-size: 2em;
        line-height: 1.5;
        font-weight: 600;
      }

      #content h3 {
        font-size: 1.4em;
        line-height: 1.5;
        font-weight: 600;
      }

      #content p {
        font-size: 1.2em;
        line-height: 1.7;
        text-align: justify;
      }
    }

    @layer base {
      /* Sticky-header offset so anchors don't hide under navbar */
      #content h1,
      #content h2,
      #content h3 {
        scroll-margin-top: 84px;
      }

      #content h1,
      #content h2 {
        position: relative; /* make them positioning parents */
      }
    }
    main {
      display: flex;
      justify-content: center; /* Horizontally center the content */
      align-items: center; /* Vertically center the content */
    }
  </style>
</svelte:head>


<div use:setOverflow class="min-h-screen bg-zinc-50 text-zinc-900">
  <!-- Navbar -->
  <header
    class="fixed w-full top-0 z-40 h-14 border-b border-black/5 backdrop-blur bg-white/50"
  >
    <div
      class="mx-auto max-w-screen-2xl px-4 sm:px-6 lg:px-8 h-full flex items-center gap-3"
    >
      <button
        id="navToggle"
        class="xl:hidden p-2 rounded-lg hover:bg-black/5 dark:hover:bg:white/10"
        aria-label="Toggle nav"
      >
        <span class="block w-5 h-0.5 bg-current mb-1"></span>
        <span class="block w-5 h-0.5 bg-current mb-1"></span>
        <span class="block w-5 h-0.5 bg-current"></span>
      </button>
      <a href="/{slug}" class="hover:opacity-50 font-semibold tracking-tight"
        >{logo}</a
      >
      <div id="themeBtn" class="ml-auto"></div>

      <Button
        id="change"
        onclick={() => goto("/signin")}
        class="text-black  hidden  stm:block cursor-pointer hover:text-green-700"
        variant="ghost"
      >
        Sign in
      </Button>

      <Buttong />
    </div>
  </header>
  <div id="navBackdrop" class="hidden fixed inset-0 bg-black/40 z-40"></div>
  <div
    class="mx-auto max-w-screen-2xl grid grid-cols-1 xl:grid-cols-[180px_minmax(0,1fr)_200px] items-start"
  >
    <!-- LEFT NAV (sections -> H2 -> optional H3 subpanel) -->
    <aside
      id="leftNav"
      class="relative xl:sticky xl:top-14 xl:h-[885px] overflow-auto border-gray-100 px-4 py-6 hidden xl:block"
    >
      <!-- Border container with indicator inside -->
      <div class="absolute right-0 top-0 h-full w-1 bg-gray-100">
        <!-- Indicator sits inside gray border -->
        <span
          aria-hidden="true"
          class="absolute bottom-[17em] -translate-y-1/2 h-10 w-full bg-gradient-to-b from-green-400 via-green-600 to-green-500 rounded-md shadow-md"
        ></span>
      </div>
    </aside>

    <!-- MAIN CONTENT (sample) -->
    <main>
      <div class="relative">
        <Pagination.Root page={currentPage} {count} {perPage} {siblingCount}>
          {#snippet children({ pages, currentPage })}
            {@const start = (currentPage - 1) * perPage}
            {@const end = currentPage * perPage}
            {@const pageitems = sections.slice(start, end)}
            <article id="content" class="markdown-body  max-w-3xl">
              <div class="flex flex-col justify-center">
                <div class="py-25">
                  <div class="wrapper-ui">
                    <div class="child-wrapper-ui">
                      <div class="headline-container">
                       <svg 
                          class="w-20 h-20 mx-auto"
                          xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                            <path fill="#21ba45" fill-rule="evenodd" d="M22 12c0 5.523-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2s10 4.477 10 10m-5.97-3.03a.75.75 0 0 1 0 1.06l-5 5a.75.75 0 0 1-1.06 0l-2-2a.75.75 0 1 1 1.06-1.06l1.47 1.47l2.235-2.235L14.97 8.97a.75.75 0 0 1 1.06 0" clip-rule="evenodd" />
                      </svg>
                      <h1
                          id="content"
                          class="headline8"
                          style="font-size: clamp(2.5rem, 3vh, 5rem);"
                        >
                         No problems detected.
 
                        </h1>

                        <p
                          class="text-base  text-gray-500 font-bold text-center sm:text-left"
                        >
                            <span class="text-lg">
                            {new Date().toLocaleString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' })}
                            </span>
                         </p>

                        <span class="text-base mt-5 text-center sm:text-left">
                          <span
                          class="inline-flex pointer-events-none items-center px-4 py-0.5 rounded-full bg-green-100 text-green-700 text-sm font-semibold no-underline"
                          style="text-decoration: none;"
                          >
                          {badge}
                          </span>
                        </span>
                        </div>

                    </div>
                  </div>
                <div class="legend mt-2">
                  <strong class="text-sm text-gray-400">Legend:</strong>
                  <span class="text-green-600 "><span class="dot operational"></span>Operational</span>
                  <span class="text-amber-600"><span class="dot partial"></span>Partial degradation</span>
                  <span class="text-red-600"><span class="dot severe"></span>Severe degradation</span>
               </div>
                </div>

                <!-- Full-width gradient divider -->
                <div aria-hidden="true">
                  <div class="relative left-1/2 -translate-x-1/2 w-screen">
                    <div
                      class="h-[1px] w-full bg-gradient-to-r from-transparent via-gray-200 to-transparent"
                    ></div>
                  </div>
                </div>

                <div class="px-10">
                  <div class="flex flex-col gap-6">
                    <div class="wrapper">
                      <div class="accordion-wrapper">
                        {#each pageitems as section (section.key)}
                          <!-- Section heading (render once per section) -->
                          {#if section.items && section.items.length}
                            <h1 class="font-semibold text-black mb-0 mt-[4rem]">
                              <span class="block text-xl mb-4 font-semibold">
                                <Label
                                  class={`inline-flex items-center px-2 py-0.5 rounded-full ${section.badgeBg} ${section.badgeText} text-sm font-semibold no-underline`}
                                >
                                  {#each section.items as it, idx (idx)}
                                    {it.date}{#if idx < section.items.length - 1},
                                    {/if}
                                  {/each}
                                </Label>
                              </span>

                              {#if section.items[0].PageTitle}
                                {section.items[0].PageTitle}
                              {:else}
                                {section.key}
                              {/if}
                            </h1>

                            <span class="block text-xl mb-4 font-semibold">
                              <div
                                class={`inline-flex w-full aspect-video items-center px-3 py-3 sm:px-5 sm:py-5  rounded-lg ${section.badgeBg} ${section.badgeText} text-sm font-semibold no-underline`}
                              >
                                {#each section.items as item, idx (idx)}
                                  <img
                                    src={item.cover}
                                    alt={item.PageTitle}
                                    class="object-cover rounded-lg mx-auto"
                                  />
                                {/each}
                              </div>
                            </span>
                          {/if}
                          {#if section.items[0].content}
                            <div class="mt-10">
                              <p
                                class="flex-1 font-semibold text-left title-clamp"
                              >
                                {#each renderContentParts(String(section.items[0].content), highlightPhrases) as part, i}
                                  {#if part.type === "text"}
                                    {part.text}
                                  {:else}
                                    <a
                                      href={part.href}
                                      class={part.class}
                                      target={part.external
                                        ? "_blank"
                                        : "_self"}
                                      rel={part.external
                                        ? "noopener noreferrer"
                                        : undefined}>{part.text}</a
                                    >
                                  {/if}
                                {/each}
                              </p>
                            </div>
                          {/if}

                          <!-- Section items (one block per item) -->
                          {#each section.items as item, idx (idx)}
                            <div class="flex items-center gap-5 mt-10 w-full">
                              <h2
                                class="flex-1 font-semibold text-left title-clamp"
                              >
                                {item.title}
                              </h2>
                            </div>

                            {#if item.features && item.features.length}
                              <ul class="mt-5 space-y-5">
                                {#each item.features as feat, i (i)}
                                  <li class="flex items-center gap-5 w-full">
                                    <span
                                      class={`${section.numberText} text-xl font-bold`}
                                      >─</span
                                    >
                                    <p
                                      class="flex-1 font-semibold text-left title-clamp"
                                    >
                                      {#each renderContentParts(String(feat), highlightPhrases) as part, i}
                                        {#if part.type === "text"}
                                          <span> {part.text}</span>
                                        {:else}
                                          <a
                                            href={part.href}
                                            class={part.class}
                                            target={part.external
                                              ? "_blank"
                                              : "_self"}
                                            rel={part.external
                                              ? "noopener noreferrer"
                                              : undefined}>{part.text}</a
                                          >
                                        {/if}
                                      {/each}
                                    </p>
                                  </li>
                                {/each}
                              </ul>
                            {/if}
                          {/each}

                          <!-- Dynamic divider: uses section.key as color token -->
                          <div aria-hidden="true" class="my-30">
                            <div class="relative left-1/2 -translate-x-1/2">
                              <div
                                class={`h-[1px] w-full bg-gradient-to-r from-transparent ${section.dividerBg} to-transparent`}
                              ></div>
                            </div>
                          </div>
                        {/each}
                      </div>
                    </div>
                  </div>
                </div>
              </div>

            </article>
            <!-- Pagination -->
            <Pagination.Content
              class="absolute bottom-0 left-1/2 transform -translate-x-1/2 flex pb-40 "
            >
              <Pagination.Item>
                <Pagination.PrevButton
                  onclick={() => setPage(currentPage - 1)}
                  class="cursor-pointer flex items-center space-x-1"
                >
                  <ChevronLeftIcon class="text-green-700 w-5 h-5" />
                  <span class="text-green-700 font-bold hidden sm:block"
                    >Previous</span
                  >
                </Pagination.PrevButton>
              </Pagination.Item>

              {#each pages as page (page.key)}
                {#if page.type === "ellipsis"}
                  <Pagination.Item>
                    <Pagination.Ellipsis />
                  </Pagination.Item>
                {:else}
                  <Pagination.Item>
                    <Pagination.Link
                      {page}
                      onclick={() => setPage(page.value)}
                      class="cursor-pointer px-3 py-1 font-bold text-green-700 rounded-md hover:bg-gray-100"
                      isActive={currentPage === page.value}
                    >
                      {page.value}
                    </Pagination.Link>
                  </Pagination.Item>
                {/if}
              {/each}

              <Pagination.Item>
                <Pagination.NextButton
                  onclick={() => setPage(currentPage + 1)}
                  class="cursor-pointer flex  items-center space-x-1"
                >
                  <span class="hidden text-green-700 font-bold sm:block"
                    >Next</span
                  >
                  <ChevronRightIcon class="text-green-700 w-5 h-5" />
                </Pagination.NextButton>
              </Pagination.Item>
            </Pagination.Content>
          {/snippet}
        </Pagination.Root>

        <div aria-hidden="true" class="h-[150px]"></div>
      </div>
    </main>

    <!-- RIGHT ToC -->
    <aside
      class="xl:sticky xl:h-[calc(100vh-56px)] top-14 px-6 py-8 border-l border-black/5 z-30"
    >
      <div class="text-xs uppercase tracking-wider text-gray-500 mb-4">
        <!-- On this page -->
      </div>
      <div
        id="toc-scroll"
        class="relative h-[calc(100%-1rem)] overflow-auto pr-2"
      >
        <span
          id="toc-rail"
          class="pointer-events-none absolute left-5 w-[2px] bg-gradient-to-b from-black/5 via-black/5 to-black/5"
        ></span>

        <ul id="toc-list" class="relative pl-7 space-y-2"></ul>
      </div>
    </aside>
  </div>
</div>



<style>
  :root {
    --wrapper-ui-max-width: 1800px;
    --wrapper-ui-min-width: 320px;
    --wrapper-ui-radius: 2rem;
    --wrapper-ui-bg-color: #f4f4f5;
  }

  .headline8 {
    text-align: center;
  }

  /* reduce gap between headline and paragraph (set to 0 to remove) */
  .headline-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    gap: 0;
    max-width: 640px; /* keep content width reasonable for center layout */
    text-align: center;
  }

  .wrapper-ui {
    display: flex;
    flex-direction: column;
    padding: 1rem;
    width: 100%;
    border-radius: var(--wrapper-ui-radius);
    background-color: var(--wrapper-ui-bg-color);
    max-width: var(--wrapper-ui-max-width);
    min-width: var(--wrapper-ui-min-width);
    margin: 0 auto;
    align-items: center;
    box-sizing: border-box;
  }

  .child-wrapper-ui {
    min-height: 300px; /* ensure vertical space */
    width: 100%;
    min-width: 200px;
    max-width: var(--wrapper-ui-max-width);
    border-radius: var(--wrapper-ui-radius);
    flex: 1;
    display: flex;
    justify-content: center; /* center horizontally */
    align-items: center; /* center vertically */
    padding: 1rem;
    padding-top: 0;
  }

  .highlight-indigo-500 {
    color: #6366f1; /* indigo-500 */
    font-weight: 600;
    background: rgba(99, 102, 241, 0.06);
    padding: 0 4px;
    border-radius: 4px;
  }

  .link-underline {
    text-decoration: underline;
    color: #4f46e5;
    font-weight: 600;
  }

  .title-clamp {
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    word-break: break-word;
    line-height: 1.4;
  }


  .legend {
      display: flex;
      align-items: center;
      gap: 20px;
      color: #d9d9d9;
      padding: 20px;
    }

  .legend span {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 14px;
    }

  .dot {
      width: 14px;
      height: 14px;
      border-radius: 50%;
      display: inline-block;
    }

  .operational {
      background-color: #21ba45; /* Green */
    }

  .partial {
      background-color: oklch(76.9% 0.188 70.08); /* Yellow */
    }

  .severe {
      background-color: #db2828; /* Red/Orange */
    }

  @media (max-width: 600px) {
    .legend {
      flex-direction: column;     /* stack vertically */
      align-items: flex-start;    /* align to left */
      gap: 8px;
    }

    .legend strong {
      margin-bottom: 4px;
    }
  }
</style>
