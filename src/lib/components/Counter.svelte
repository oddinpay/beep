<script lang="ts">
  import { browser } from '$app/environment';
  import { derived, readable, writable } from 'svelte/store';
  import { onMount } from 'svelte';
  import { localStore } from '$lib/storage';

  const logOutput = writable<string | null>(null);

  // Stores
  const count   = localStore<number>('count', 0);
  const fruits  = localStore<string[]>('fruits', []);
  const session = localStore<string[]>('session', ['🍎', '🍌'], 5000); // 5s
  const note    = localStore<string>('note', 'hello', 15000);          // 15s

  // Expired flags
  const sessionExpired = session.expired;
  const noteExpired    = note.expired;

  // 1s ticker
  const now = readable<number>(0, (set) => {
    if (!browser) return () => {};
    set(Date.now());
    const id = setInterval(() => set(Date.now()), 1000);
    return () => clearInterval(id);
  });

  // Remaining ms helper (−1 no TTL, 0 expired)
  function remainingMsFor(key: string, ttlMs?: number) {
    if (!ttlMs) return derived(now, () => -1);
    return derived(now, () => {
      if (!browser) return -1;
      const raw = localStorage.getItem(key);
      if (!raw) return 0;
      try {
        const { e } = JSON.parse(raw) as { e?: number | null };
        if (!e) return -1;
        const ms = e - Date.now();
        return ms > 0 ? ms : 0;
      } catch { return -1; }
    });
  }

  const sessionRemaining = remainingMsFor('session', 5000);
  const noteRemaining    = remainingMsFor('note', 15000);
  const s = (ms: number) => (ms < 0 ? '—' : Math.ceil(ms / 1000) + 's');

  // Actions
  const logAll = () => {
  const banner = '🧭 localStore + TTL Demo';
  const line = '──────────────────────────────';
  const summary = [
    `Count: ${count.get()}`,
    `Fruits: ${JSON.stringify(fruits.get())}`,
    `Session: ${JSON.stringify(session.get())} (expired: ${session.isExpired()})`,
    `Note: ${JSON.stringify(note.get())} (expired: ${note.isExpired()})`
  ];

  const out = [banner, line, ...summary, line].join('\n');

  console.clear();
  console.log(`%c${banner}`, 'color:white;background:#2563eb;padding:6px 12px;border-radius:6px;font-weight:bold;');
  console.log(line);
  console.log(`%c${summary[0]}`, 'color:lightgreen;font-weight:bold;');
  summary.slice(1).forEach((text) => console.log(text));

  logOutput.set(out);
};


  const setSessionList = () => session.set(['🍇', '🍉']);
  const setSessionPreserveExpiry = () => session.set(['🍍', '🥝'], { preserveExpiry: true });
  const updateSessionPush = () => session.update((a) => [...(a ?? []), '🍓']);
  const updateSessionPreserveExpiry = () => session.update((a) => [...(a ?? []), '🍑'], { preserveExpiry: true });
  const getSession = () => alert(session.get() ? `session: ${JSON.stringify(session.get())}` : 'session: null');
  const deleteSession = () => { session.delete(); alert('Deleted "session"'); };
  const checkSessionExpired = () => alert(`session.isExpired(): ${session.isExpired()}`);

  const setNote = () => note.set('new note (TTL refreshed)');
  const setNotePreserveExpiry = () => note.set('new note (expiry preserved)', { preserveExpiry: true });
  const updateNoteAppend = () => note.update((t) => (t ?? '') + ' + appended');
  const deleteNote = () => note.delete();
  const checkNoteExpired = () => alert(`note.isExpired(): ${note.isExpired()}`);

  const addApple = () => fruits.update((a) => [...a, '🍎']);
  const clearFruits = () => fruits.set([]);
  const incCount = () => count.update((n) => (n ?? 0) + 1);
  const setCount5 = () => count.set(5);
  const deleteCount = () => count.delete();

  let ready = false; onMount(() => ready = true);
</script>

{#if ready}
  <main class="wrap">
    <h1>localStore + TTL Demo</h1>
    <button class="ghost" onclick={logAll}>Console log</button>
    {#if $logOutput}
      <pre class="console-output">{$logOutput}</pre>
    {/if}

    <!-- Counter -->
    <section>
      <h2>Counter (no TTL)</h2>
      <p class="big">{$count}</p>
      <div class="grid">
        <button onclick={setCount5}>Set to 5</button>
        <button onclick={incCount}>Increment</button>
        <button class="danger" onclick={deleteCount}>Delete</button>
      </div>
    </section>

    <!-- Fruits -->
    <section>
    <h2>Set</h2>
    <p class="wrap-json">
      {#if $fruits.length}
        {JSON.stringify($fruits)}
      {:else}
        <span class="muted">Empty</span>
      {/if}
    </p>
    <div class="grid wrap-break-word">
      <button onclick={addApple}>Add 🍎</button>
      <button onclick={clearFruits}>Clear</button>
    </div>
    </section>

    <!-- Session -->
    <section>
      <h2>TTL & Set / Get / Del</h2>
      <ul>
      <li><strong>Value:</strong> <span class="wrap-json">{JSON.stringify($session)}</span></li>
      <li><strong>Expired:</strong> {$sessionExpired ? 'Yes' : 'No'}</li>
      <li><strong>Remaining:</strong> {s($sessionRemaining)}</li>
      </ul>
      <div class="grid">
      <button onclick={setSessionList}>set(['🍇','🍉'])</button>
      <button class="ghost" onclick={setSessionPreserveExpiry}>set preserveExpiry</button>
      <button onclick={updateSessionPush}>update push 🍓</button>
      <button class="ghost" onclick={updateSessionPreserveExpiry}>update preserveExpiry</button>
      <button class="subtle" onclick={getSession}>get()</button>
      <button class="subtle" onclick={checkSessionExpired}>isExpired()</button>
      <button class="danger" onclick={deleteSession}>delete()</button>
      </div>
    </section>

    <!-- Note -->
    <section>
      <h2>Long TTL & SET / GET</h2>
      <ul>
        <li><strong>Value:</strong> {$note}</li>
        <li><strong>Expired:</strong> {$noteExpired ? 'Yes' : 'No'}</li>
        <li><strong>Remaining:</strong> {s($noteRemaining)}</li>
      </ul>
      <div class="grid">
        <button onclick={setNote}>set('new note')</button>
        <button class="ghost" onclick={setNotePreserveExpiry}>set preserveExpiry</button>
        <button onclick={updateNoteAppend}>update append</button>
        <button class="subtle" onclick={checkNoteExpired}>isExpired()</button>
        <button class="danger" onclick={deleteNote}>delete</button>
      </div>
    </section>
  </main>
{/if}

<style>
  :root{ --border:#e5e7eb; --muted:#6b7280; --bg:#f8fafc; --card:#fff; }
  .wrap{ max-width: 760px; margin: 2rem auto; padding: 0 1rem 3rem; }
  h1{ margin:0 0 1rem; font-size:1.25rem; }
  h2{ margin:.5rem 0; font-size:1rem; }
  section{ background:var(--card); border:1px solid var(--border); border-radius:12px; padding:1rem; margin:1rem 0; }
  ul{ margin:.25rem 0 .5rem; padding-left:1.1rem; }
  .grid{ display:grid; gap:.5rem; grid-template-columns: repeat(auto-fit, minmax(180px,1fr)); }
  .big{ font-size:2rem; font-weight:700; margin:.25rem 0 .5rem; }
  .muted{ color:var(--muted); }
  button{ padding:.5rem .75rem; border:1px solid var(--border); border-radius:10px; background:#fff; cursor:pointer; }
  button:hover{ background:#f9fafb; }
  .ghost{ background:transparent; }
  .subtle{ background:#f3f4f6; }
  .danger{ background:#fef2f2; border-color:#fecaca; color:#b91c1c; }
  .console-output {
    margin-top: 0.75rem;
    background: #e6e6e9;
    color: #000000;
    padding: 1rem;
    border-radius: 8px;
    white-space: pre-wrap;
    line-height: 1.4;
  }
  .wrap-json {
    word-wrap: break-word;   
    overflow-wrap: anywhere; 
    white-space: normal;     
  }
</style>

