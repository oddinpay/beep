import { writable } from "svelte/store";

const CACHE_KEY = "probeMapCache";

const loadCache = () => {
    if (typeof localStorage === "undefined") return null;
    try {
        const cached = localStorage.getItem(CACHE_KEY);
        if (!cached) return null;
        const parsed = JSON.parse(cached);
        return parsed && typeof parsed === "object" ? parsed as Record<number, any> : null;
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

function createProbeMapStore() {
    const { subscribe, set, update } = writable<Record<number, any>>({});
    const loading = writable(true);

    const applyCacheOrClear = () => {
        const cached = loadCache();
        if (cached) {
            set(cached);
            return true;
        }
        set({});
        return false;
    };

    applyCacheOrClear();

    (async () => {
        try {
            const res = await fetch("http://127.0.0.1:8976/v1/status");
            const body = await res.text();

            if (!res.ok) {
                console.error(`Failed to fetch probe status (${res.status} ${res.statusText}):`, body);
                applyCacheOrClear();
                return;
            }

            let parsed: unknown;
            try {
                parsed = JSON.parse(body);
            } catch {
                console.error("Unexpected non-JSON response from /v1/status:", body);
                applyCacheOrClear();
                return;
            }

            if (!Array.isArray(parsed)) {
                console.error("Expected an array response from /v1/status but received:", parsed);
                applyCacheOrClear();
                return;
            }

            const parsedList = parsed as any[];

            const entries = parsedList
                .map((item) => [Number(item.index), { ...(item.payload?.probe ?? {}), ...(item.payload?.sla ?? {}) }] as const)
                .sort((a, b) => a[0] - b[0]);

            const next = Object.fromEntries(entries);
            set(next);
            saveCache(next);
        } catch (error) {
            console.error("Error fetching data:", error);
            applyCacheOrClear();
        } finally {
            loading.set(false);
        }
    })();

    return { subscribe, set, update, loading };
}

export const probeMap = createProbeMapStore();
