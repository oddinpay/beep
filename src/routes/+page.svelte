<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import Buttong from "$lib/components/Buttong.svelte";
  import Loader from "$lib/components/Loader.svelte";
  import Footer from "$lib/components/Footer.svelte";
  import {
    Tabs,
    TabsContent,
    TabsList,
    TabsTrigger,
  } from "$lib/components/ui/tabs";
  import { source } from "sveltekit-sse";
  // import { localStore } from '$lib/storage';
  import { onMount } from "svelte";
  import { browser } from "$app/environment";
  import timer from "$lib/timer";
  import { env } from "$env/dynamic/public";

  const slug = "https://oddinpay.com";
  const logo = "oddin status";
  const title = "OddinPay System - Status";
  const description =
    "Real-time and historical data on OddinPay system performance.";
  const badge = "Last updated";

  let ready = $state(false);
  onMount(() => (ready = true));
  const clock = timer();
  const TOTAL_DAYS = 90;
  const today = new Date();
  const end = new Date(
    Date.UTC(today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate()),
  );

  const start = new Date(end);
  start.setUTCDate(end.getUTCDate() - 89);

  type StatusType = "up" | "down" | "warn" | "default";

  interface StatusEntry {
    date: Date;
    status: StatusType;
  }

  interface ApiData {
    id?: string;
    name?: string;
    date?: string | Date;
    state?: string;
    statuses: StatusEntry[];
    uptime15: string;
    uptime30: string;
    uptime60: string;
    uptime90: string;
  }

  const beepHost = env.PUBLIC_BEEP_HOST;
  const json = source(`https://${beepHost}/v1/sse`).select("").json<ApiData>();

  type Buffered = { probe: ApiData; sla?: any; index?: number };

  const pending = new Map<string, Buffered>();
  let flushTimer: ReturnType<typeof setTimeout> | null = null;
  const FLUSH_DELAY = 50;

  function scheduleFlush() {
    if (flushTimer) return;
    flushTimer = setTimeout(() => {
      flushTimer = null;
      flushPending();
    }, FLUSH_DELAY);
  }

  function flushPending() {
    if (!pending.size) return;

    const nextMap: Record<string, ApiData> = { ...(probeMap as ProbeMap) };

    for (const [id, { probe, sla, index }] of pending) {
      const existing = nextMap[id];
      const order = Number.isFinite(index)
        ? index
        : ((existing as any)?.__order ?? Number.POSITIVE_INFINITY);

      const nextProbe = {
        ...(existing ?? {}),
        ...probe,
        uptime90: sla?.uptime90,
        __order: order,
      };

      nextMap[id] = nextProbe;
    }

    pending.clear();

    const sortedEntries = Object.entries(nextMap).sort(
      ([, a], [, b]) =>
        ((a as any).__order ?? Number.POSITIVE_INFINITY) -
        ((b as any).__order ?? Number.POSITIVE_INFINITY),
    );

    probeMap = Object.fromEntries(sortedEntries) as ProbeMap;
  }

  json.subscribe((msg: any) => {
    const probe = msg?.payload?.probe;
    const sla = msg?.payload?.sla;
    const index = msg?.index;
    if (!probe?.id) return;

    pending.set(probe.id, { probe, sla, index });

    scheduleFlush();
  });

  type ProbeMap = Record<string, ApiData>;

  let probeMap = $state<ProbeMap>({});

  // const statusStore = localStore<StatusType[]>('status', []);

  // --- utils ---

  function coerceStatus(s?: StatusType): StatusType {
    return s === "up" || s === "down" || s === "warn" ? s : "warn";
  }

  function asStatus(s: string): StatusType {
    return s === "up" || s === "down" || s === "warn" ? s : "default";
  }

  function currentStatusFor(x: {
    statuses?: StatusEntry[];
    status?: StatusType;
  }): StatusType {
    const arr = Array.isArray(x?.statuses) ? x!.statuses! : [];

    if (arr.length) {
      const idx = Math.max(0, Math.min(dayIndex, arr.length - 1));
      const candidate = arr[idx]?.status as StatusType | undefined;
      return coerceStatus(candidate);
    }

    return coerceStatus(x.status);
  }

  function monitorStatus(x: any): StatusType {
    const dates = Array.isArray(x?.date) ? x.date : x?.date ? [x.date] : [];
    const states = Array.isArray(x?.state)
      ? x.state
      : x?.state
        ? [x.state]
        : [];

    if (dates.length && states.length) {
      let newestIdx = 0;
      let newestTime = -Infinity;
      for (let i = 0; i < dates.length; i++) {
        const t = +parseDate(dates[i]);
        if (t > newestTime) {
          newestTime = t;
          newestIdx = i;
        }
      }
      return asStatus(states[newestIdx] ?? states[0]);
    }

    if (Array.isArray(x?.state) && x.state.length) return asStatus(x.state[0]);
    if (typeof x?.state === "string") return asStatus(x.state);
    if (typeof x?.status === "string") return asStatus(x.status);

    return "warn";
  }

  function parseDate(dateString: string | Date): Date {
    if (dateString instanceof Date) return dateString;
    if (typeof dateString !== "string") return new Date(String(dateString));
    const parts = dateString.split("/");
    if (parts.length !== 3) return new Date(dateString);
    const [day, month, year] = parts;
    return new Date(`${year}-${month}-${day}`);
  }

  function normalizeDates(rawDate: unknown): string[] {
    if (!rawDate) return [];

    if (Array.isArray(rawDate)) {
      return rawDate.map((d) => parseDate(d).toLocaleDateString());
    }

    if (typeof rawDate === "string") {
      return rawDate
        .split(",")
        .map((s) => s.trim())
        .filter(Boolean)
        .map((s) => parseDate(s).toLocaleDateString());
    }

    return [parseDate(rawDate as Date | string).toLocaleDateString()];
  }

  let mockData = $derived.by(() => {
    const probes = Object.values(probeMap) as ApiData[];

    const unique = new Map<string, ApiData>();
    for (const p of probes) {
      if (!p || !p.id) continue;
      unique.set(String(p.name), p);
    }
    const uniqueProbes = Array.from(unique.values());

    const limitedProbes = uniqueProbes.slice(0, 3);

    return limitedProbes.map((probe) => {
      const datesList = normalizeDates(probe?.date ?? new Date());
      const statesList = Array.isArray(probe?.state)
        ? probe.state
        : [probe?.state];

      const datesMap = new Map<string, StatusType>();
      datesList.forEach((date, index) => {
        const state = statesList[index] ?? "default";
        datesMap.set(date, state as StatusType);
      });

      const statuses: StatusEntry[] = Array.from(
        { length: TOTAL_DAYS },
        (_, i) => {
          const tempDate = new Date(start);
          tempDate.setUTCDate(start.getUTCDate() + i);
          const key = tempDate.toLocaleDateString();
          const resolved = datesMap.get(key) ?? "default";
          return {
            date: tempDate,
            status: resolved,
          };
        },
      );

      const api: ApiData = {
        name: String(probe?.name ?? ""),
        statuses,
        uptime15: String(probe?.uptime90 ?? "00.000%"),
        uptime30: String(probe?.uptime90 ?? "00.000%"),
        uptime60: String(probe?.uptime90 ?? "00.000%"),
        uptime90: String(probe?.uptime90 ?? "00.000%"),
      };

      return api;
    });
  });

  let monitors = $derived.by(() => {
    const apiMap = new Map<
      string,
      { title: string; description: string; status: string }
    >();
    const probes = Object.values(probeMap) as any[];
    const unique = new Map<string, any>();

    for (const p of probes) {
      if (!p) continue;
      unique.set(String(p.name), p);
      let id = String((p as any).id ?? "");
      if (!id) continue;
    }

    const uniqueProbes = Array.from(unique.values());
    const dataToShow = uniqueProbes.slice(3);
    const nameCounts = new Map<string, number>();

    for (const api of dataToShow) {
      if (!api) continue;
      let id = String((api as any).id ?? "");
      if (!id) continue;

      if (nameCounts.has(id)) {
        const count = nameCounts.get(id)! + 1;
        nameCounts.set(id, count);
        id = `${id}+${count}`;
      } else {
        nameCounts.set(id, 0);
      }

      apiMap.set(id, {
        title: String((api as any).name ?? id),
        description: "API",
        status: monitorStatus(api),
      });
    }

    const apiMonitors = Array.from(apiMap.values());

    const merged = new Map<
      string,
      { title: string; description: string; status: string }
    >();
    for (const m of [...apiMonitors]) {
      merged.set(m.title, m);
    }

    return Array.from(merged.values());
  });

  const pickStatus = (xs: string[]): StatusType | undefined => {
    if (xs.includes("down")) return "down";
    if (xs.includes("warn")) return "warn";
    if (xs.includes("up")) return "up";
  };

  let overallStatus = $derived.by<StatusType | undefined>(() => {
    const allStatuses = [
      ...mockData.map((api) => currentStatusFor(api)),
      ...monitors.map((monitor) => asStatus(monitor.status)),
    ];

    if (browser && allStatuses.length === 0) {
      // const stored = statusStore.get();
      // statusStore.set([]);
      // return stored ? stored[0] : undefined;
    }

    const computed = pickStatus(allStatuses);
    // const storedStatuses = statusStore.get() ?? [];
    // const currentStoredStatus = storedStatuses[0];

    if (computed !== computed) {
      // && currentStoredStatus
      // statusStore.set([computed]);
      return computed;
    }

    // ?? currentStoredStatus

    return computed;
  });

  const dayIndex = Math.floor(
    (today.getTime() - start.getTime()) / (1000 * 60 * 60 * 24),
  );

  function getStartLabelForDays(days: number): string {
    const tempStart = new Date(end);
    tempStart.setDate(end.getDate() - days);
    return tempStart.toLocaleString(undefined, {
      month: "short",
      day: "numeric",
    });
  }

  const startLabel90 = getStartLabelForDays(89);
  const startLabel60 = getStartLabelForDays(59);
  const startLabel30 = getStartLabelForDays(29);
  const startLabel15 = getStartLabelForDays(14);
  const endLabel = end.toLocaleString(undefined, {
    month: "short",
    day: "numeric",
  });

  const Indicators = {
    Resolved: {
      badge: "resolved",
      statusLabel: "Resolved",
    },
    Investigating: {
      badge: "investigating",
      statusLabel: "Investigating",
    },
    Scheduled: {
      badge: "scheduled",
      statusLabel: "Scheduled",
    },
    Identified: {
      badge: "identified",
      statusLabel: "Identified",
    },
    Inprogress: {
      badge: "inprogress",
      statusLabel: "In Progress",
    },
    Completed: {
      badge: "completed",
      statusLabel: "Completed",
    },
  } as const;

  // Each value inside Indicators
  type Indicator = (typeof Indicators)[keyof typeof Indicators];

  interface IncidentEntry {
    time: string;
    description: string;
    status: Exclude<
      Indicator,
      typeof Indicators.Completed | typeof Indicators.Scheduled
    >;
  }

  interface MaintenanceEntry {
    time: string;
    description: string;
    status: Exclude<
      Indicator,
      typeof Indicators.Resolved | typeof Indicators.Investigating
    >;
  }

  interface Incident {
    title: string;
    entries: IncidentEntry[];
  }

  interface Maintenance {
    service: string;
    entries: MaintenanceEntry[];
  }

  let incidents: Incident[] = [
    {
      title: "Elevated iDeal errors",
      entries: [
        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Investigating,
          description:
            "We are investigating reports of increased errors on iDeal payments.",
        },
        {
          time: "Sep 22, 2025 20:14 UTC",
          status: Indicators.Identified,
          description:
            "From 13:05–19:15 UTC, we saw elevated errors on iDeal payments. This is now resolved.",
        },
        {
          time: "Sep 22, 2025 12:45 UTC",
          status: Indicators.Inprogress,
          description:
            "We are investigating reports of increased errors on iDeal payments.",
        },
        {
          time: "Sep 22, 2025 19:15 UTC",
          status: Indicators.Resolved,
          description:
            "From 13:05–19:15 UTC, we saw elevated errors on iDeal payments. This is now resolved.",
        },
      ],
    },
    {
      title: "Payment errors",
      entries: [
        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Investigating,
          description:
            "We are investigating reports of increased errors on iDeal payments.",
        },
        {
          time: "Sep 22, 2025 12:45 UTC",
          status: Indicators.Inprogress,
          description:
            "We are investigating reports of increased errors on iDeal payments.",
        },

        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Identified,
          description:
            "We are identifying reports of increased errors on iDeal payments.",
        },

        {
          time: "Sep 22, 2025 20:14 UTC",
          status: Indicators.Resolved,
          description:
            "From 13:05–19:15 UTC, we saw elevated errors on iDeal payments. This is now resolved.",
        },
      ],
    },
  ];

  const statusPriority = new Map<Indicator, number>([
    [Indicators.Completed, 0],
    [Indicators.Resolved, 0],
    [Indicators.Inprogress, 1],
    [Indicators.Identified, 2],
    [Indicators.Investigating, 3],
    [Indicators.Scheduled, 3],
  ]);

  incidents.forEach((incident) => {
    incident.entries.sort((a, b) => {
      return (
        (statusPriority.get(a.status) ?? Infinity) -
        (statusPriority.get(b.status) ?? Infinity)
      );
    });
  });

  let maintenances: Maintenance[] = [
    {
      service: "API errors",
      entries: [
        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Scheduled,
          description:
            "We are investigating reports of increased errors on API.",
        },
        {
          time: "Sep 22, 2025 12:45 UTC",
          status: Indicators.Inprogress,
          description:
            "We are investigating reports of increased errors on API.",
        },

        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Identified,
          description: "We are identifying reports of increased errors on API.",
        },

        {
          time: "Sep 22, 2025 20:14 UTC",
          status: Indicators.Completed,
          description:
            "From 13:05–19:15 UTC, we saw elevated errors on API. This is now resolved.",
        },
      ],
    },
    {
      service: "Shop errors",
      entries: [
        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Scheduled,
          description:
            "We are investigating reports of increased errors on Shop.",
        },
        {
          time: "Sep 22, 2025 12:45 UTC",
          status: Indicators.Inprogress,
          description:
            "We are investigating reports of increased errors on Shop.",
        },

        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Identified,
          description: "We are identifying reports of increased errors on Shop.",
        },
      ],
    },
  ];

  maintenances.forEach((maintenance) => {
    maintenance.entries.sort((a, b) => {
      return (
        (statusPriority.get(a.status) ?? Infinity) -
        (statusPriority.get(b.status) ?? Infinity)
      );
    });
  });

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
    data?: { monitors?: boolean; miniMonitors?: boolean };
  }

  const allProps = $props() as RoadmapProps & Record<string, unknown>;

  const {
    sections: inputSections,
    data,
    status,
    features,
    cover,
    ..._
  } = allProps;

  const submitMonitors = $derived(data?.monitors);
  const miniMonitors = $derived(data?.miniMonitors);

  $effect(() => {
    console.log("Data prop:", submitMonitors);
    console.log("Mini monitors prop:", miniMonitors);
  });

  function setOverflow(_node: HTMLElement) {
    if (window.matchMedia && window.matchMedia("(min-width: 1279px)").matches) {
      document.documentElement.style.overflow = "unset";
    }
  }

  let activeTab = $state("tab-1");
  let direction = $state<"right" | "left">("right");

  $effect(() => {
    console.log("Active tab is now:", activeTab);
    console.log("Direction is now:", direction);
  });

  const tabsOrder = ["tab-1", "tab-2"];

  function handleChange(newValue: string) {
    const oldIndex = tabsOrder.indexOf(activeTab);
    const newIndex = tabsOrder.indexOf(newValue);

    direction = newIndex > oldIndex ? "left" : "right";
    activeTab = newValue;
  }

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
</script>

<Loader />

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

    .monitors {
      -webkit-box-orient: vertical;
      -webkit-line-clamp: 1;
      display: -webkit-box;
      overflow: hidden;
    }

    #content p {
      font-size: 15px;
      line-height: 1.7;
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
        font-size: 15px;
        line-height: 1.7;
        text-align: balance;
      }
    }

    @layer base {
      #content h1,
      #content h2,
      #content h3 {
        scroll-margin-top: 84px;
      }

      #content h1,
      #content h2 {
        position: relative;
      }
    }

    main {
      display: flex;
      justify-content: center;
      align-items: center;
    }

    :root {
      --wrapper-ui-max-width: 1800px;
      --wrapper-ui-min-width: 320px;
      --wrapper-ui-radius: 2rem;
      --wrapper-ui-bg-color: #f9fafb;
      --bg: #ffffff;
      --text: #000000;
      --up: #4ce04c;
      --warn: #f2a900;
      --down: #f05d5e;
      --inactive: #6b7280;
      --inactive-service: #6b7280;
      --active-service: black;
      --default: #e5e7eb;
      --chip-radius: 1px;
      --today-ring: rgba(0, 0, 0, 0.25);
    }

    .headline8 {
      text-align: center;
    }

    .headline-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 100%;
      gap: 0;
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
      min-height: 300px;
      width: 100%;
      min-width: 200px;
      max-width: var(--wrapper-ui-max-width);
      border-radius: var(--wrapper-ui-radius);
      flex: 1;
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 1rem;
      padding-top: 0;
    }

    .highlight-indigo-500 {
      color: #6366f1;
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
      gap: 4px;
      font-size: 14px;
    }

    @media (max-width: 600px) {
      .legend {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
      }

      .legend strong {
        margin-bottom: 4px;
      }
    }

    .status-page {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 50px;
      max-width: 1200px;
      margin: 40px auto;
      padding: 0 20px;
    }

    .status-card {
      padding: 15px;
      border-radius: 8px;
      margin-bottom: 15px;
      background: #fff;
      color: #222;
      box-shadow:
        rgba(0, 0, 0, 0.05) 0 6px 24px,
        rgba(0, 0, 0, 0.08) 0 0 0 1px;
    }

    .incident-card {
      background: #f9fafb;
      color: black;
      padding: 30px;
      border-radius: 8px;
      margin-bottom: 20px;
      box-shadow: rgba(0, 0, 0, 0.05) 0px 0px 0px 1px;
    }

    .status-entry {
      margin-bottom: 25px;
    }

    .time {
      font-size: 0.9rem;
      color: var(--inactive);
      margin-right: 10px;
    }

    .badge.resolved {
      background: #d7f7c2;
      color: #006908;
      font-weight: 600;
      border: 1px solid #a6eb84;
    }

    .badge.investigating {
      background: #ebeef1;
      color: #545969;
      font-weight: 600;
      border: 1px solid #d5dbe1;
    }

    .badge.inprogress {
      background: #fff4e5;
      color: #b45309;
      font-weight: 600;
      border: 1px solid #ffddb3;
    }

    .badge.completed {
      background: #d7f7c2;
      color: #006908;
      font-weight: 600;
      border: 1px solid #a6eb84;
    }

    .badge.identified {
      background: white;
      color: #4b5563;
      font-weight: 600;
      border: 1px solid #cecece;
    }

    .badge.scheduled {
      background: white;
      color: #4b5563;
      font-weight: 600;
      border: 1px solid #cecece;
    }

    .badge2.completed {
      background: #d7f7c2;
      color: #006908;
      font-weight: 600;
      border: 1px solid #a6eb84;
    }

    .badge2.inprogress {
      background: #fff4e5;
      color: #b45309;
      font-weight: 600;
      border: 1px solid #ffddb3;
    }

    .badge2.scheduled {
      background: white;
      color: #4b5563;
      font-weight: 600;
      border: 1px solid #cecece;
    }

    .badge {
      font-size: 0.8rem;
      padding: 0px 6px;
      display: inline-block;
      white-space: nowrap;
      border-radius: 4px;
    }

    .badge2 {
      font-size: 0.8rem;
      padding: 2px 6px;
      display: inline-block;
      white-space: nowrap;
      border-radius: 4px;
    }

    .maintenance-list {
      display: flex;
      flex-direction: column;
      gap: 15px;
    }

    .maintenance-card {
      padding: 15px;
      border-radius: 6px;
      color: black;
    }

    .header {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 5px;
    }

    .service {
      font-weight: 600;
      color: var(--inactive-service);
    }

    @media (max-width: 750px) {
      .status-page {
        grid-template-columns: 1fr;
      }
    }

    .card {
      max-width: 950px;
      margin: 40px auto;
      padding: 40px;
      background: var(--bg);
      border-radius: 10px;
      color: var(--text);
      box-shadow:
        rgba(0, 0, 0, 0.05) 0 6px 24px,
        rgba(0, 0, 0, 0.08) 0 0 0 1px;
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      margin-bottom: 12px;
      font-weight: 600;
      flex-wrap: wrap;
      gap: 6px;
    }

    .bar {
      display: flex;
      gap: 2px;
      justify-content: space-between;
    }

    .chip {
      flex: 1 1 0;
      display: inline-block;
      height: 24px;
      background: var(--default);
      border-radius: var(--chip-radius);
    }

    .chip:hover {
      transform: scaleY(1.15);
    }

    .chip.up {
      background: var(--up);
    }
    .chip.warn {
      background: var(--warn);
    }
    .chip.down {
      background: var(--down);
    }
    .chip.today {
      box-shadow: inset 0 0 0 2px var(--today-ring);
    }
    .chip {
      display: none;
    }

    /* >=901px: last 90 (all) */
    @media (min-width: 901px) {
      .chip:nth-last-child(-n + 90) {
        display: block;
      }
      .uptimes .uptime90 {
        display: inline;
      }
      .uptimes .uptime60,
      .uptimes .uptime30,
      .uptimes .uptime15 {
        display: none;
      }
      .timeline .label90 {
        display: inline;
      }
      .timeline .label60,
      .timeline .label30,
      .timeline .label15 {
        display: none;
      }
    }

    /* 601–900px: last 60 */
    @media (min-width: 601px) and (max-width: 900px) {
      .chip:nth-last-child(-n + 60) {
        display: block;
      }
      .uptimes .uptime60 {
        display: inline;
      }
      .uptimes .uptime90,
      .uptimes .uptime30,
      .uptimes .uptime15 {
        display: none;
      }
      .timeline .label60 {
        display: inline;
      }
      .timeline .label90,
      .timeline .label30,
      .timeline .label15 {
        display: none;
      }
    }

    /* 311–600px: last 30 */
    @media (min-width: 311px) and (max-width: 600px) {
      .chip:nth-last-child(-n + 30) {
        display: block;
      }
      .uptimes .uptime30 {
        display: inline;
      }
      .uptimes .uptime90,
      .uptimes .uptime60,
      .uptimes .uptime15 {
        display: none;
      }
      .timeline .label30 {
        display: inline;
      }
      .timeline .label90,
      .timeline .label60,
      .timeline .label15 {
        display: none;
      }
    }

    /* <=310px: last 15 */
    @media (max-width: 310px) {
      .chip:nth-last-child(-n + 15) {
        display: block;
      }
      .uptimes .uptime15 {
        display: inline;
      }
      .uptimes .uptime90,
      .uptimes .uptime60,
      .uptimes .uptime30 {
        display: none;
      }
      .timeline .label15 {
        display: inline;
      }
      .timeline .label90,
      .timeline .label60,
      .timeline .label30 {
        display: none;
      }
    }

    @media (max-width: 165px) {
      .bar {
        width: 90px;
        height: 20px;
      }
    }

    .timeline {
      display: flex;
      justify-content: space-between;
      margin-top: 8px;
      font-size: 0.85rem;
      color: #9ea0a3;
    }

    .profile {
      aspect-ratio: 1;
      animation-delay: 1s;
      border-radius: 99px;
      margin-inline: auto;
    }

    .shimmer {
      background: linear-gradient(90deg, #ededed 30%, #dcdcdc 50%, #ededed 70%);
      background-size: 400%;
      animation: shimmer 1.5s infinite linear;
    }

    @keyframes shimmer {
      0% {
        background-position: 100% 100%;
      }
      100% {
        background-position: 0 0;
      }
    }
  </style>
</svelte:head>

{#if ready}
  <div use:setOverflow class="min-h-screen bg-white text-zinc-900">
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

        <div class="oddin-status hover:opacity-50">
          <a href={slug} target="_blank" rel="noopener noreferrer">
            {logo}
          </a>
        </div>

        <div id="themeBtn" class="ml-auto"></div>
        <Button
          id="change"
          onclick={() => window.open("https://oddinpay.com/signin", "_blank")}
          class="text-black hidden stm:block cursor-pointer hover:text-green-700"
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
      <!-- LEFT NAV -->
      <aside
        id="leftNav"
        class="relative xl:sticky xl:top-14 xl:h-[885px] overflow-auto border-gray-100 px-4 py-6 hidden xl:block"
      >
        <div class="absolute right-0 top-0 h-full w-1 bg-gray-100">
          <span
            aria-hidden="true"
            class="absolute bottom-[20em] -translate-y-1/2 h-10 w-full bg-gradient-to-b from-green-400 via-green-600 to-green-500 rounded-md shadow-md"
          ></span>
        </div>
      </aside>

      <!-- MAIN CONTENT -->
      <main class="min-w-0">
        <div class="relative p-5">
          <article id="content" class="markdown-body p-5 max-w-5xl">
            <div class="flex flex-col justify-center">
              <div class="py-25">
                <div class="wrapper-ui">
                  <div class="child-wrapper-ui">
                    <div class="headline-container">
                      {#if overallStatus === "up"}
                        <svg
                          class="w-20 h-20 mx-auto"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 24 24"
                        >
                          <path
                            fill="#21ba45"
                            fill-rule="evenodd"
                            d="M22 12c0 5.523-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2s10 4.477 10 10m-5.97-3.03a.75.75 0 0 1 0 1.06l-5 5a.75.75 0 0 1-1.06 0l-2-2a.75.75 0 1 1 1.06-1.06l1.47 1.47l2.235-2.235L14.97 8.97a.75.75 0 0 1 1.06 0"
                            clip-rule="evenodd"
                          />
                        </svg>
                      {:else if overallStatus === "warn"}
                        <svg
                          class="w-20 h-20 mx-auto"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 20 20"
                        >
                          <path
                            fill="#d97706"
                            d="M10 2c4.42 0 8 3.58 8 8s-3.58 8-8 8s-8-3.58-8-8s3.58-8 8-8m1.13 9.38l.35-6.46H8.52l.35 6.46zm-.09 3.36c.24-.23.37-.55.37-.96c0-.42-.12-.74-.36-.97s-.59-.35-1.06-.35s-.82.12-1.07.35s-.37.55-.37.97c0 .41.13.73.38.96c.26.23.61.34 1.06.34s.8-.11 1.05-.34"
                          />
                        </svg>
                      {:else if overallStatus === "down"}
                        <svg
                          class="w-20 h-20 mx-auto"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 24 24"
                        >
                          <path
                            fill="#db2828"
                            d="M12 2c5.53 0 10 4.47 10 10s-4.47 10-10 10S2 17.53 2 12S6.47 2 12 2m3.59 5L12 10.59L8.41 7L7 8.41L10.59 12L7 15.59L8.41 17L12 13.41L15.59 17L17 15.59L13.41 12L17 8.41z"
                          />
                        </svg>
                      {:else}
                        <div
                          role="status"
                          class="profile shimmer w-20 h-20 bg-gray-300"
                        ></div>
                      {/if}

                      <h1
                        id="content"
                        class="headline8"
                        style="font-size: clamp(2.5rem, 3vh, 5rem);"
                      >
                        {overallStatus === "up"
                          ? "All Systems Operational"
                          : overallStatus === "warn"
                            ? "System Outage Detected"
                            : overallStatus === "down"
                              ? "Critical Issues Detected"
                              : "No Status Available"}
                      </h1>
                      <p
                        class="text-base text-gray-500 font-bold text-center sm:text-left"
                      >
                        <span class="text-lg">
                          {$clock}
                        </span>
                      </p>
                      <span class="text-base mt-5 text-center sm:text-left">
                        <span
                          class="inline-flex pointer-events-none items-center px-4 py-0.5 rounded-full text-sm font-semibold no-underline"
                          style=" 
  border: 1px solid {overallStatus === 'up'
                            ? '#a6eb84'
                            : overallStatus === 'warn'
                              ? '#ffddb3'
                              : overallStatus === 'down'
                                ? '#f05d5e'
                                : '#d1d5db'}; 
  text-decoration: none; 
  background-color: {overallStatus === 'up'
                            ? '#d7f7c2'
                            : overallStatus === 'warn'
                              ? '#fff4e5'
                              : overallStatus === 'down'
                                ? '#fddede'
                                : '#f3f4f6'}; 
  color: {overallStatus === 'up'
                            ? '#006908'
                            : overallStatus === 'warn'
                              ? '#b45309'
                              : overallStatus === 'down'
                                ? '#db2828'
                                : '#4b5563'}; 
                        "
                        >
                          {badge}
                        </span>
                      </span>
                    </div>
                  </div>
                </div>

                <div class="legend mt-2">
                  <strong class="text-sm text-gray-400">Legend:</strong>
                  <span class="text-green-600">
                    <svg
                      class="w-5 h-5 inline-block"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 24 24"
                    >
                      <path
                        fill="#21ba45"
                        fill-rule="evenodd"
                        d="M22 12c0 5.523-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2s10 4.477 10 10m-5.97-3.03a.75.75 0 0 1 0 1.06l-5 5a.75.75 0 0 1-1.06 0l-2-2a.75.75 0 1 1 1.06-1.06l1.47 1.47l2.235-2.235L14.97 8.97a.75.75 0 0 1 1.06 0"
                        clip-rule="evenodd"
                      />
                    </svg>
                    Operational
                  </span>
                  <span class="text-amber-600">
                    <svg
                      class="w-5 h-5 inline-block"
                      xmlns="http://www.w3.org/2000/svg"
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fill="#d97706"
                        d="M10 2c4.42 0 8 3.58 8 8s-3.58 8-8 8s-8-3.58-8-8s3.58-8 8-8m1.13 9.38l.35-6.46H8.52l.35 6.46zm-.09 3.36c.24-.23.37-.55.37-.96c0-.42-.12-.74-.36-.97s-.59-.35-1.06-.35s-.82.12-1.07.35s-.37.55-.37.97c0 .41.13.73.38.96c.26.23.61.34 1.06.34s.8-.11 1.05-.34"
                      />
                    </svg>
                    Partial degradation
                  </span>
                  <span class="text-red-600">
                    <svg
                      class="w-5 h-5 inline-block"
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 24 24"
                    >
                      <path
                        fill="#db2828"
                        d="M12 2c5.53 0 10 4.47 10 10s-4.47 10-10 10S2 17.53 2 12S6.47 2 12 2m3.59 5L12 10.59L8.41 7L7 8.41L10.59 12L7 15.59L8.41 17L12 13.41L15.59 17L17 15.59L13.41 12L17 8.41z"
                      />
                    </svg>
                    Severe degradation
                  </span>
                </div>

                <div aria-hidden="true">
                  <div class="relative left-1/2 -translate-x-1/2 w-screen">
                    <div
                      class="h-[1px] w-full bg-gradient-to-r from-transparent via-gray-200 to-transparent"
                    ></div>
                  </div>
                </div>

                <Tabs
                  value={activeTab}
                  class="items-left p-5"
                  onValueChange={handleChange}
                >
                  <TabsList
                    class="border-border text-foreground h-auto gap-2 rounded-none border-b bg-transparent px-0 py-1"
                  >
                    {#each tabsOrder as t, i}
                      <TabsTrigger
                        value={t}
                        class={`cursor-pointer hover:bg-accent hover:text-foreground transition-colors duration-150 ease-in-out data-[state=active]:hover:bg-accent relative after:absolute after:inset-x-0 after:bottom-0 after:-mb-1 after:h-0.5 after:transform after:scale-x-0 after:transition-transform after:duration-200 after:ease-in-out data-[state=active]:after:scale-x-100 data-[state=active]:after:bg-primary data-[state=active]:bg-transparent data-[state=active]:shadow-none ${
                          direction === "left"
                            ? "after:origin-left"
                            : "after:origin-right"
                        }`}
                      >
                        {#if i === 0}
                          <span
                            class={activeTab === t
                              ? "text-black font-semibold"
                              : "text-gray-400 font-semibold"}>Live status</span
                          >
                        {:else if i === 1}
                          <span
                            class={activeTab === t
                              ? "text-black font-semibold"
                              : "text-gray-400 font-semibold"}>History</span
                          >
                        {/if}
                      </TabsTrigger>
                    {/each}
                  </TabsList>
                  {#each tabsOrder as t, i}
                    <TabsContent value={t}>
                      {#if i === 0}
                        {#if submitMonitors && mockData.length === 0}
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            width="48"
                            height="48"
                            viewBox="0 0 24 24"
                            ><circle cx="18" cy="12" r="0" fill="currentColor"
                              ><animate
                                attributeName="r"
                                begin=".67"
                                calcMode="spline"
                                dur="1.5s"
                                keySplines="0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8"
                                repeatCount="indefinite"
                                values="0;2;0;0"
                              /></circle
                            ><circle cx="12" cy="12" r="0" fill="currentColor"
                              ><animate
                                attributeName="r"
                                begin=".33"
                                calcMode="spline"
                                dur="1.5s"
                                keySplines="0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8"
                                repeatCount="indefinite"
                                values="0;2;0;0"
                              /></circle
                            ><circle cx="6" cy="12" r="0" fill="currentColor"
                              ><animate
                                attributeName="r"
                                begin="0"
                                calcMode="spline"
                                dur="1.5s"
                                keySplines="0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8;0.2 0.2 0.4 0.8"
                                repeatCount="indefinite"
                                values="0;2;0;0"
                              /></circle
                            ></svg
                          >
                        {:else}
                          {#each mockData as api, index}
                            <section class="card" style="margin-bottom: 2px;">
                              <div class="card-header">
                                <div
                                  style="display: flex; align-items: center; gap: 5px;"
                                >
                                  {#if currentStatusFor(api) === "up"}
                                    <svg
                                      class="w-5 h-5"
                                      xmlns="http://www.w3.org/2000/svg"
                                      viewBox="0 0 24 24"
                                    >
                                      <path
                                        fill="#21ba45"
                                        fill-rule="evenodd"
                                        d="M22 12c0 5.523-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2s10 4.477 10 10m-5.97-3.03a.75.75 0 0 1 0 1.06l-5 5a.75.75 0 0 1-1.06 0l-2-2a.75.75 0 1 1 1.06-1.06l1.47 1.47l2.235-2.235L14.97 8.97a.75.75 0 0 1 1.06 0"
                                        clip-rule="evenodd"
                                      />
                                    </svg>
                                  {:else if currentStatusFor(api) === "warn"}
                                    <svg
                                      class="w-5 h-5"
                                      xmlns="http://www.w3.org/2000/svg"
                                      viewBox="0 0 20 20"
                                    >
                                      <path
                                        fill="#d97706"
                                        d="M10 2c4.42 0 8 3.58 8 8s-3.58 8-8 8s-8-3.58-8-8s3.58-8 8-8m1.13 9.38l.35-6.46H8.52l.35 6.46zm-.09 3.36c.24-.23.37-.55.37-.96c0-.42-.12-.74-.36-.97s-.59-.35-1.06-.35s-.82.12-1.07.35s-.37.55-.37.97c0 .41.13.73.38.96c.26.23.61.34 1.06.34s.8-.11 1.05-.34"
                                      />
                                    </svg>
                                  {:else if currentStatusFor(api) === "down"}
                                    <svg
                                      class="w-5 h-5"
                                      xmlns="http://www.w3.org/2000/svg"
                                      viewBox="0 0 24 24"
                                    >
                                      <path
                                        fill="#db2828"
                                        d="M12 2c5.53 0 10 4.47 10 10s-4.47 10-10 10S2 17.53 2 12S6.47 2 12 2m3.59 5L12 10.59L8.41 7L7 8.41L10.59 12L7 15.59L8.41 17L12 13.41L15.59 17L17 15.59L13.41 12L17 8.41z"
                                      />
                                    </svg>
                                  {/if}
                                  <div class="text-lg">{api.name}</div>
                                </div>
                                <div class="uptimes">
                                  <span class="uptime15"
                                    >{api.uptime15} uptime</span
                                  >
                                  <span class="uptime30"
                                    >{api.uptime30} uptime</span
                                  >
                                  <span class="uptime60"
                                    >{api.uptime60} uptime</span
                                  >
                                  <span class="uptime90"
                                    >{api.uptime90} uptime</span
                                  >
                                </div>
                              </div>
                              <div class="bar">
                                {#each api.statuses as s, i}
                                  <div
                                    class="chip {s.status} {i === dayIndex
                                      ? s.status
                                      : ''}"
                                  ></div>
                                {/each}
                              </div>
                              <div class="timeline">
                                <span class="label15">15 days ago</span>
                                <span class="label30">30 days ago</span>
                                <span class="label60">60 days ago</span>
                                <span class="label90">90 days ago</span>
                                <span>Today</span>
                              </div>
                            </section>
                          {/each}
                        {/if}

                        <div class="status-page">
                          <div class="left">
                            <h3>System status</h3>
                            {#if !miniMonitors && monitors.length === 0}
                              No monitors available.
                            {:else if submitMonitors && monitors.length === 0}
                              <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="24"
                                height="24"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  fill="none"
                                  stroke="currentColor"
                                  stroke-linecap="round"
                                  stroke-width="2"
                                  d="M12 6.99998C9.1747 6.99987 6.99997 9.24998 7 12C7.00003 14.55 9.02119 17 12 17C14.7712 17 17 14.75 17 12"
                                  ><animateTransform
                                    attributeName="transform"
                                    attributeType="XML"
                                    dur="560ms"
                                    from="0,12,12"
                                    repeatCount="indefinite"
                                    to="360,12,12"
                                    type="rotate"
                                  />
                                </path>
                              </svg>
                            {:else if !submitMonitors}
                              No monitors available.
                            {:else}
                              {#each monitors as status}
                                <div class="status-card">
                                  <div
                                    style="display: flex; align-items: center; gap: 10px;"
                                  >
                                    {#if status.status === "up"}
                                      <svg
                                        class="w-10 h-10 inline-block"
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 24 24"
                                      >
                                        <path
                                          fill="#21ba45"
                                          fill-rule="evenodd"
                                          d="M22 12c0 5.523-4.477 10-10 10S2 17.523 2 12S6.477 2 12 2s10 4.477 10 10m-5.97-3.03a.75.75 0 0 1 0 1.06l-5 5a.75.75 0 0 1-1.06 0l-2-2a.75.75 0 1 1 1.06-1.06l1.47 1.47l2.235-2.235L14.97 8.97a.75.75 0 0 1 1.06 0"
                                          clip-rule="evenodd"
                                        />
                                      </svg>
                                    {:else if status.status === "warn"}
                                      <svg
                                        class="w-10 h-10 inline-block"
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 20 20"
                                      >
                                        <path
                                          fill="#d97706"
                                          d="M10 2c4.42 0 8 3.58 8 8s-3.58 8-8 8s-8-3.58-8-8s3.58-8 8-8m1.13 9.38l.35-6.46H8.52l.35 6.46zm-.09 3.36c.24-.23.37-.55.37-.96c0-.42-.12-.74-.36-.97s-.59-.35-1.06-.35s-.82.12-1.07.35s-.37.55-.37.97c0 .41.13.73.38.96c.26.23.61.34 1.06.34s.8-.11 1.05-.34"
                                        />
                                      </svg>
                                    {:else if status.status === "down"}
                                      <svg
                                        class="w-10 h-10 inline-block"
                                        xmlns="http://www.w3.org/2000/svg"
                                        viewBox="0 0 24 24"
                                      >
                                        <path
                                          fill="#db2828"
                                          d="M12 2c5.53 0 10 4.47 10 10s-4.47 10-10 10S2 17.53 2 12S6.47 2 12 2m3.59 5L12 10.59L8.41 7L7 8.41L10.59 12L7 15.59L8.41 17L12 13.41L15.59 17L17 15.59L13.41 12L17 8.41z"
                                        />
                                      </svg>
                                    {/if}
                                    <div>
                                      <strong>{status.title}</strong>
                                      <p class="monitors" style="color: #666;">
                                        {status.description}
                                      </p>
                                    </div>
                                  </div>
                                </div>
                              {/each}
                            {/if}
                          </div>

                          <div class="right">
                            <h3>Recent incidents</h3>
                            {#if incidents.every( (incident) => incident.entries.some((entry) => entry.status === Indicators.Resolved), )}
                              No incidents to display
                            {:else}
                              {#each incidents as incident}
                                {#if !incident.entries.some((entry) => entry.status === Indicators.Resolved)}
                                  <div class="incident-card">
                                    {#each incident.entries as entry}
                                      <div class="status-entry">
                                        <span class="time font-bold"
                                          >{entry.time}</span
                                        >
                                        <span
                                          class="badge mt-1 {entry.status
                                            .badge}"
                                        >
                                          {entry.status.statusLabel}
                                        </span>
                                        <p
                                          class="mt-2 text-gray-600"
                                          style="font-size: 16px"
                                        >
                                          {entry.description}
                                        </p>
                                      </div>
                                    {/each}
                                  </div>
                                {/if}
                              {/each}
                            {/if}
                            <h3>Maintenance</h3>
                            <div class="maintenance-list">
                              {#if maintenances.length === 0}
                                No maintenance windows available.
                              {/if}
                              {#each maintenances as maintenance}
                                {#if !maintenance.entries.some((entry) => entry.status === Indicators.Completed)}
                                  {#each maintenance.entries as entry}
                                    <div
                                      class="flex justify-between items-center p-3 gap-4"
                                    >
                                      <span
                                        class="inline-flex items-center px-2.5 badge2 py-1 rounded-full text-xs font-medium {entry
                                          .status.badge}"
                                      >
                                        {entry.status.statusLabel}
                                      </span>
                                      <div
                                        class="flex flex-col text-left leading-tight"
                                      >
                                        <span
                                          class="text-base font-semibold text-[var(--inactive-service)]"
                                        >
                                          {maintenance.service}
                                        </span>
                                        <time
                                          class="text-base text-[var(--inactive)]"
                                        >
                                          {entry.time}
                                        </time>
                                      </div>
                                    </div>
                                  {/each}
                                {/if}
                              {/each}
                            </div>
                          </div>
                        </div>
                      {:else if i === 1}
                        {#if incidents.length === 0 && maintenances.length === 0}
                          <p
                            class="p-10 text-gray-600 text-center"
                            style="font-size: 16px"
                          >
                            No historical incidents available.
                          </p>
                        {/if}
                        {#each incidents as incident}
                          <div class="incident-card mt-10">
                            <h3>{incident.title}</h3>
                            {#each incident.entries as entry}
                              <div class="status-entry">
                                <span class="time font-bold">{entry.time}</span>
                                <span class="badge mt-1 {entry.status.badge}">
                                  {entry.status.statusLabel}
                                </span>
                                <p
                                  class="mt-2 text-gray-600"
                                  style="font-size: 16px"
                                >
                                  {entry.description}
                                </p>
                              </div>
                            {/each}
                          </div>
                        {/each}

                        {#each maintenances as maintenance}
                          <div class="incident-card mt-10">
                            <h3>
                              Scheduled maintenance for {maintenance.service}
                            </h3>
                            {#each maintenance.entries as entry}
                              <div class="status-entry">
                                <span class="time font-bold">{entry.time}</span>
                                <span class="badge mt-1 {entry.status.badge}">
                                  {entry.status.statusLabel}
                                </span>
                                <p
                                  class="mt-2 text-gray-600"
                                  style="font-size: 16px"
                                >
                                  System affected: {maintenance.service}
                                </p>
                              </div>
                            {/each}
                          </div>
                        {/each}
                      {/if}
                    </TabsContent>
                  {/each}
                </Tabs>
              </div>
            </div>
          </article>
          <div aria-hidden="true" class="h-[150px]"></div>
        </div>
      </main>

      <!-- RIGHT ToC -->
      <aside
        class="xl:sticky xl:h-[calc(100vh-56px)] top-14 px-6 py-8 border-l border-black/5 z-30"
      >
        <div class="text-xs uppercase tracking-wider text-gray-500 mb-4"></div>
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

  <Footer />
{/if}
