import { writable } from "svelte/store";

const CACHE_KEY = "probeMapCache";

function readCache(): Record<number, any> | null {
    if (typeof localStorage === "undefined") return null;
    try {
        const raw = localStorage.getItem(CACHE_KEY);
        return raw ? JSON.parse(raw) : null;
    } catch {
        return null;
    }
}

function writeCache(value: Record<number, any>) {
    if (typeof localStorage === "undefined") return;
    try {
        localStorage.setItem(CACHE_KEY, JSON.stringify(value));
    } catch (error) {
        console.warn("Failed to cache probe map:", error);
    }
}

function createProbeMapStore() {
    const cached = readCache() ?? {};
    const baseStore = writable<Record<number, any>>(cached);
    const loading = writable(Object.keys(cached).length === 0);
    const { subscribe, set: baseSet, update: baseUpdate } = baseStore;

    const set = (value: Record<number, any>) => {
        writeCache(value);
        baseSet(value);
    };

    const update = (fn: (value: Record<number, any>) => Record<number, any>) => {
        baseUpdate((current) => {
            const next = fn(current);
            writeCache(next);
            return next;
        });
    };

    (async () => {
        try {
            const res = await fetch("http://127.0.0.1:8976/v1/status");
            const data = await res.json();

            const entries = (data as any[])
                .map(
                    (item) =>
                        [
                            Number(item.index),
                            { ...(item.payload?.probe ?? {}), ...(item.payload?.sla ?? {}) },
                        ] as const,
                )
                .sort((a, b) => a[0] - b[0]);

            set(Object.fromEntries(entries));
        } catch (error) {
            console.error("Error fetching data:", error);
        } finally {
            loading.set(false);
        }
    })();

    return { subscribe, set, update, loading };
}

export const probeMap = createProbeMapStore();
