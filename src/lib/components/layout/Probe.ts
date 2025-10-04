import { writable } from "svelte/store";

const CACHE_KEY = "probeMapCache";
const POLL_INTERVAL_MS = 2000;

const PLACEHOLDER_PROBES: Record<number, any> = {
    0: {
        name: "Probe unavailable",
        status: "unknown",
        description: "Probe status could not be retrieved.",
    },
};

const delay = (ms: number) => new Promise<void>((resolve) => setTimeout(resolve, ms));

const loadCache = () => {
    if (typeof localStorage === "undefined") return null;
    try {
        const cached = localStorage.getItem(CACHE_KEY);
        if (!cached) return null;
        const parsed = JSON.parse(cached);
        return parsed && typeof parsed === "object" ? (parsed as Record<number, any>) : null;
    } catch (error) {
        console.warn("Failed to load probe cache:", error);
        return null;
    }
};

const saveCache = (value: Record<number, any>) => {
    if (typeof localStorage === "undefined") return;
    try {
        localStorage.setItem(CACHE_KEY, JSON.stringify(value));
    } catch (error) {
        console.warn("Failed to persist probe cache:", error);
    }
};

const fetchProbeStatus = async (): Promise<Record<number, any> | null> => {
    try {
        const res = await fetch("http://127.0.0.1:8976/v1/status");
        const body = await res.text();

        if (!res.ok) {
            console.error(`Failed to fetch probe status (${res.status} ${res.statusText}):`, body);
            return null;
        }

        let parsed: unknown;
        try {
            parsed = JSON.parse(body);
        } catch {
            console.error("Unexpected non-JSON response from /v1/status:", body);
            return null;
        }

        if (!Array.isArray(parsed)) {
            console.error("Expected an array response from /v1/status but received:", parsed);
            return null;
        }

        const parsedList = parsed as any[];
        if (!parsedList.length) {
            console.warn("Empty probe list received from /v1/status.");
            return null;
        }

        const entries = parsedList
            .map((item) => [Number(item.index), { ...(item.payload?.probe ?? {}), ...(item.payload?.sla ?? {}) }] as const)
            .sort((a, b) => a[0] - b[0]);

        return Object.fromEntries(entries);
    } catch (error) {
        console.error("Error fetching data:", error);
        return null;
    }
};

function createProbeMapStore() {
    const { subscribe, set, update } = writable<Record<number, any>>({});
    const loading = writable(true);

    let readyResolve: (() => void) | null = null;
    let readyReject: ((reason?: unknown) => void) | null = null;
    let readySettled = false;

    const ready = new Promise<void>((resolve, reject) => {
        readyResolve = resolve;
        readyReject = reject;
    });

    const resolveReady = () => {
        if (!readySettled) {
            readySettled = true;
            readyResolve?.();
        }
    };

    const rejectReady = (reason: unknown) => {
        if (!readySettled) {
            readySettled = true;
            readyReject?.(reason);
        }
    };

    const useCacheOrPlaceholder = () => {
        const cached = loadCache();
        if (cached) {
            set(cached);
            return true;
        }
        set(PLACEHOLDER_PROBES);
        return false;
    };

    useCacheOrPlaceholder();

    const poll = async () => {
        while (true) {
            const next = await fetchProbeStatus();
            if (next) {
                set(next);
                saveCache(next);
                loading.set(false);
                resolveReady();
                return;
            }
            await delay(POLL_INTERVAL_MS);
        }
    };

    poll().catch((error) => {
        console.error("Error during probe polling:", error);
        rejectReady(error);
        loading.set(false);
    });

    return { subscribe, set, update, loading, ready };
}

export const probeMap = createProbeMapStore();

