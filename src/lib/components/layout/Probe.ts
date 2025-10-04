import { browser } from "$app/environment";
import { readable, writable } from "svelte/store";

type ProbeMap = Record<number, any>;

function createProbeMapStore() {
    const state = writable<ProbeMap>({});
    const loading = writable(browser);

    const load = async () => {
        if (!browser) return;
        loading.set(true);

        try {
            const res = await fetch("http://127.0.0.1:8976/v1/status");
            const data = (await res.json()) as any[];

            const entries = data
                .map((item) => [Number(item.index), { ...(item.payload?.probe ?? {}), ...(item.payload?.sla ?? {}) }] as const)
                .sort((a, b) => a[0] - b[0]);

            state.set(Object.fromEntries(entries));
        } catch (error) {
            console.error("Error fetching data:", error);
            state.set({});
        } finally {
            loading.set(false);
        }
    };

    load();

    return {
        subscribe: state.subscribe,
        refresh: load,
        loading: readable(false, (set) => loading.subscribe(set)),
    };
}

export const probeMap = createProbeMapStore();
