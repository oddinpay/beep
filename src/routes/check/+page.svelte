<script lang="ts">

  const TOTAL_DAYS = 90;
  const today = new Date();
  const end = new Date(today);
  const start = new Date(end);
  start.setDate(end.getDate() - (TOTAL_DAYS));

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
      <span class="label15">{startLabel15}</span>
      <span class="label30">{startLabel30}</span>
      <span class="label60">{startLabel60}</span>
      <span class="label90">{startLabel90}</span>
      <span>{endLabel}</span>
    </div>
  </section>
</div>


<style>
  :root{
    --bg:#FFFFFF;
    --text:#000000;
    --up:#4ce04c;
    --warn:#f2a900;
    --down:#f05d5e;
    --default:#e5e7eb;
    --chip-radius:2px;
    --today-ring:rgba(0,0,0,.25);
  }

  .layout{ padding:10px; }

  .card{
    max-width:950px;
    margin:40px auto;
    padding:40px;
    background:var(--bg);
    border-radius:10px;
    color:var(--text);
    box-shadow: rgba(0,0,0,0.05) 0 6px 24px, rgba(0,0,0,0.08) 0 0 0 1px;
  }

  .card-header{
    display:flex; justify-content:space-between; margin-bottom:12px;
    font-weight:600; flex-wrap:wrap; gap:6px;
  }

  /* Use Grid for crisp equal-width chips */
  .bar{ display:grid; grid-auto-flow:column; grid-auto-columns:1fr; gap:2px; }

  .chip{
    background:var(--default);
    height:24px;
    border-radius:var(--chip-radius);
  }
  .chip:hover{ transform:scaleY(1.15) }

  .chip.up{ background:var(--up); }
  .chip.warn{ background:var(--warn); }
  .chip.down{ background:var(--down); }
  .chip.today{ box-shadow: inset 0 0 0 2px var(--today-ring); }
  .chip{ display:none; }

  /* >=901px: last 90 (all) */
  @media (min-width: 901px){
    .chip:nth-last-child(-n+90){ display:block; }
    .uptimes .uptime90{ display:inline; }
    .uptimes .uptime60, .uptimes .uptime30, .uptimes .uptime15{ display:none; }
    .timeline .label90{ display:inline; }
    .timeline .label60, .timeline .label30, .timeline .label15{ display:none; }
  }

  /* 601–900px: last 60 */
  @media (min-width: 601px) and (max-width: 900px){
    .chip:nth-last-child(-n+60){ display:block; }
    .uptimes .uptime60{ display:inline; }
    .uptimes .uptime90, .uptimes .uptime30, .uptimes .uptime15{ display:none; }
    .timeline .label60{ display:inline; }
    .timeline .label90, .timeline .label30, .timeline .label15{ display:none; }
  }

  /* 311–600px: last 30 */
  @media (min-width: 311px) and (max-width: 600px){
    .chip:nth-last-child(-n+30){ display:block; }
    .uptimes .uptime30{ display:inline; }
    .uptimes .uptime90, .uptimes .uptime60, .uptimes .uptime15{ display:none; }
    .timeline .label30{ display:inline; }
    .timeline .label90, .timeline .label60, .timeline .label15{ display:none; }
  }

  /* <=310px: last 15 */
  @media (max-width: 310px){
    .chip:nth-last-child(-n+15){ display:block; }
    .uptimes .uptime15{ display:inline; }
    .uptimes .uptime90, .uptimes .uptime60, .uptimes .uptime30{ display:none; }
    .timeline .label15{ display:inline; }
    .timeline .label90, .timeline .label60, .timeline .label30{ display:none; }
  }

   @media (max-width: 165px) {
   .bar {
    width: 90px;     
    height: 20px;    
  }
  } 
  .timeline{
    display:flex; justify-content:space-between; margin-top:8px;
    font-size:0.85rem; color:#9ea0a3;
  }
</style>