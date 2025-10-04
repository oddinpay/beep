import { writable } from "svelte/store";

function createProbeMapStore() {
    const { subscribe, set, update } = writable<Record<number, any>>({});
    const loading = writable(true);

    (async () => {
        try {
            const res = await fetch("http://127.0.0.1:8976/v1/status");
            const data = await res.json();

            const entries = (data as any[])
                .map((item) => [Number(item.index), { ...(item.payload?.probe ?? {}), ...(item.payload?.sla ?? {}) }] as const)
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

