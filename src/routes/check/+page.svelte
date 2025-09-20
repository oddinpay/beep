<script lang="ts">
  const TOTAL_DAYS = 90;
  const today = new Date();
  const end = new Date(today);
  const start = new Date(end);
  start.setDate(end.getDate() - (TOTAL_DAYS - 1));

  type StatusType = "up" | "down" | "warn" | "default";
  interface StatusEntry {
    date: Date;
    status: StatusType;
  }

  const statuses: StatusEntry[] = Array.from({ length: TOTAL_DAYS }, (_, i) => {
    const tempDate = new Date(start);
    tempDate.setDate(tempDate.getDate() + i);
    return { date: tempDate, status: "default" };
  });

  function markStatus(dateString: string, status: StatusType) {
    const date = new Date(dateString);
    const index = statuses.findIndex(
      (entry) => entry.date.toDateString() === date.toDateString()
    );
    if (index !== -1) {
      statuses[index].status = status;
    }
  }

  // Example marks
  markStatus("2025-06-23", "down");
  markStatus("2025-09-20", "warn");

  // --- precompute uptime values ---
  function uptimeForLast(n: number): string {
    const recent = statuses.slice(-n);
    const upDays = recent.filter((s) => s.status !== "down").length;
    return ((upDays / n) * 100).toFixed(3);
  }

  const uptime15 = uptimeForLast(15);
  const uptime30 = uptimeForLast(30);
  const uptime60 = uptimeForLast(60);
  const uptime90 = uptimeForLast(90);

  const dayIndex = Math.floor(
    (today.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)
  );

  const startLabel = start.toLocaleString(undefined, { month: "short", day: "numeric" });
  const endLabel = end.toLocaleString(undefined, { month: "short", day: "numeric" });
</script>

<div class="layout">
  <section class="card">
    <div class="card-header">
      <div>API</div>
      <div class="uptimes">
        <span class="uptime15">{uptime15}% uptime</span>
        <span class="uptime30">{uptime30}% uptime</span>
        <span class="uptime60">{uptime60}% uptime</span>
        <span class="uptime90">{uptime90}% uptime</span>
      </div>
    </div>

    <div class="bar">
      {#each statuses as s, i}
        <div class="chip {s.status} {i === dayIndex ? s.status : ''}"></div>
      {/each}
    </div>

    <div class="timeline">
      <span>{startLabel}</span>
      <span>{endLabel}</span>
    </div>
  </section>
</div>

<style>
  :root {
    --bg: #0e2a3f;
    --text: #e9f2fb;
    --up: #4ce04c;
    --warn: #f2a900;
    --down: #f05d5e;
    --default: #ffffff;
    --chip-radius: 1px;
  }

  .layout { padding: 5px; }

  .card {
    max-width: 1150px;
    margin: 40px auto;
    padding: 20px 32px;
    background: var(--bg);
    border-radius: 12px;
    box-shadow: 0 0 10px rgba(0,0,0,0.25);
    color: var(--text);
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
    display: grid;
    gap: 2px;
  }

  .chip {
    border-radius: var(--chip-radius);
    background: var(--default);
    height: 22px;
  }

  /* Hide all uptime spans by default */
  .uptimes .uptime15,
  .uptimes .uptime30,
  .uptimes .uptime60 {
    display: none;
  }
  .uptimes .uptime90 {
    display: inline;
  }

  /* Show only the right one per breakpoint */
  @media (max-width: 310px) {
    .bar { grid-template-columns: repeat(15, 1fr); }
    .chip:nth-child(n+16) { display: none; }
    .uptimes .uptime90 { display: none; }
    .uptimes .uptime15 { display: inline; }
    .chip { height: 14px; }
  }

  @media (min-width: 311px) and (max-width: 600px) {
    .bar { grid-template-columns: repeat(30, 1fr); }
    .chip:nth-child(n+31) { display: none; }
    .uptimes .uptime90 { display: none; }
    .uptimes .uptime30 { display: inline; }
    .chip { height: 18px; }
  }

  @media (min-width: 601px) and (max-width: 900px) {
    .bar { grid-template-columns: repeat(60, 1fr); }
    .chip:nth-child(n+61) { display: none; }
    .uptimes .uptime90 { display: none; }
    .uptimes .uptime60 { display: inline; }
    .chip { height: 20px; }
  }

  @media (min-width: 901px) {
    .bar { grid-template-columns: repeat(90, 1fr); }
    .uptime90 { display: inline; }
    .chip { height: 22px; }
  }

  .chip.warn { background: var(--warn); }
  .chip.down { background: var(--down); }
  .chip.up   { background: var(--up); }

  .timeline {
    display: flex;
    justify-content: space-between;
    margin-top: 8px;
    font-size: 0.85rem;
    color: var(--muted);
  }
</style>

