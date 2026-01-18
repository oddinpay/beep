import { browser } from '$app/environment';
import { writable, readable, type Readable, type Writable } from 'svelte/store';

type Stored<T> = { v: T; e?: number | null };

export type LocalStore<T> = Writable<T> & {
  get(): T | null;
  set(value: T, opts?: { preserveExpiry?: boolean }): void;
  update(updater: (value: T) => T, opts?: { preserveExpiry?: boolean }): void;
  delete(): void;
  expired: Readable<boolean>;
  /** Imperative expiry check. */
  isExpired(): boolean;
};

function now() {
  return Date.now();
}

function pack<T>(v: T, expiry?: number | null): string {
  const payload: Stored<T> = { v, e: expiry ?? null };
  return JSON.stringify(payload);
}

function unpack<T>(raw: string | null): Stored<T> | null {
  if (!raw) return null;
  try {
    const parsed = JSON.parse(raw) as Stored<T>;
    return parsed && typeof parsed === 'object' ? parsed : null;
  } catch {
    return null;
  }
}

export function localStore<T>(key: string, initial: T, ttlMs?: number): LocalStore<T> {
  const hasTTL = typeof ttlMs === 'number' && ttlMs > 0;

  const load = (): { value: T; expired: boolean } => {
    if (!browser) return { value: initial, expired: false };
    const stored = unpack<T>(localStorage.getItem(key));
    if (!stored) return { value: initial, expired: false };
    const expired = !!stored.e && now() > (stored.e as number);
    if (expired) {
      try {
        localStorage.removeItem(key);
      } catch {}
      return { value: initial, expired: true };
    }
    return { value: stored.v, expired: false };
  };

  const { value: startValue, expired: startExpired } = load();
  const inner = writable<T>(startValue);

  const expired = readable<boolean>(startExpired, (set) => {
    if (!browser || !hasTTL) {
      set(false);
      return () => {};
    }
    let timer: number | undefined;
    const tick = () => {
      const s = unpack<T>(localStorage.getItem(key));
      const isExp = !!(s?.e) && now() > (s!.e as number);
      set(isExp);
      timer = window.setTimeout(tick, 250);
    };
    tick();
    return () => {
      if (timer) clearTimeout(timer);
    };
  });

  const write = (value: T, opts?: { preserveExpiry?: boolean }) => {
    inner.set(value);
    if (!browser) return;
    try {
      let expiry: number | null = null;
      if (hasTTL) {
        if (opts?.preserveExpiry) {
          const existing = unpack<T>(localStorage.getItem(key));
          expiry = (existing?.e as number | null) ?? (now() + (ttlMs as number));
        } else {
          expiry = now() + (ttlMs as number);
        }
      }
      localStorage.setItem(key, pack(value, expiry));
    } catch {
     
    }
  };

  const api: LocalStore<T> = {
    subscribe: inner.subscribe,
    set(value: T, opts?: { preserveExpiry?: boolean }) {
      write(value, opts);
    },
    update(updater: (v: T) => T, opts?: { preserveExpiry?: boolean }) {
      let current: T;
      const unsub = inner.subscribe((v) => (current = v as T));
      unsub!();
      write(updater(current!), opts);
    },
    get(): T | null {
      const loaded = load();
      return loaded.expired ? null : loaded.value;
    },
    delete() {
      if (browser) {
        try {
          localStorage.removeItem(key);
        } catch {}
      }
      inner.set(initial);
    },
    expired,
    isExpired(): boolean {
      if (!browser) return false;
      const s = unpack<T>(localStorage.getItem(key));
      return !!(s?.e) && now() > (s!.e as number);
    }
  };

  if (browser) {
    const existing = unpack<T>(localStorage.getItem(key));
    if (!existing) {
      const expiry = hasTTL ? now() + (ttlMs as number) : null;
      try {
        localStorage.setItem(key, pack(startValue, expiry));
      } catch {}
    }
  }

  return api;
}

