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
  markStatus("2025-09-19", "warn");

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

  function getStartLabelForDays(days: number): string {
    const tempStart = new Date(end);
    tempStart.setDate(end.getDate() - (days));
    return tempStart.toLocaleString(undefined, { month: "short", day: "numeric" });
  }

  const startLabel90 = getStartLabelForDays(90);
  const startLabel60 = getStartLabelForDays(60);
  const startLabel30 = getStartLabelForDays(30);
  const startLabel15 = getStartLabelForDays(15);

  const endLabel = end.toLocaleString(undefined, { month: "short", day: "numeric" });
</script>

<div class="layout">
  <section class="card border">
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
      <span class="label15">{startLabel15}</span>
      <span class="label30">{startLabel30}</span>
      <span class="label60">{startLabel60}</span>
      <span class="label90">{startLabel90}</span>
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
    max-width: 950px;
    margin: 40px auto;
    padding: 40px 40px ;
    background: var(--bg);
    border-radius: 10px;
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
    grid-gap: 2px; 
  }

  .chip {
    border-radius: var(--chip-radius);
    background: var(--default);
    height: 24px;
    transition: transform 0.2s ease-in-out;
  }

  .chip:hover {
    transform: scale(1.2);
  }

  /* Hide all by default */
  .timeline .label15,
  .timeline .label30,
  .timeline .label60 {
    display: none;
  }
  .timeline .label90 { display: inline; }

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
    .timeline .label90 { display: none; }
    .timeline .label15 { display: inline; }
  }

  @media (min-width: 311px) and (max-width: 600px) {
    .bar { grid-template-columns: repeat(30, 1fr); }
    .chip:nth-child(n+31) { display: none; }
    .uptimes .uptime90 { display: none; }
    .uptimes .uptime30 { display: inline; }
    .timeline .label90 { display: none; }
    .timeline .label30 { display: inline; }
  }

  @media (min-width: 601px) and (max-width: 900px) {
    .bar { grid-template-columns: repeat(60, 1fr); }
    .chip:nth-child(n+61) { display: none; }
    .uptimes .uptime90 { display: none; }
    .uptimes .uptime60 { display: inline; }
    .timeline .label90 { display: none; }
    .timeline .label60 { display: inline; }
  }

  @media (min-width: 901px) {
    .bar { grid-template-columns: repeat(90, 1fr); }
    .uptime90 { display: inline; }
    .timeline .label90 { display: inline; }
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

  @media (max-width: 165px) {
   .bar {
    width: 90px;     
    height: 20px;    
  }
 }
</style>

