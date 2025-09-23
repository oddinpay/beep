<script lang="ts">

  const TOTAL_DAYS = 90;
  const today = new Date();
  const end = new Date(today); 
  const start = new Date(end);
  start.setDate(end.getDate() - 89);

  type StatusType = "up" | "down" | "warn" | "default";

  interface StatusEntry {
    date: Date;
    status: StatusType;
  }

  interface ApiData {
    name: string;
    statuses: StatusEntry[];
    uptime15: string;
    uptime30: string;
    uptime60: string;
    uptime90: string;
  }

  function calculateUptime(api: ApiData, n: number): string {
    const recent = api.statuses.slice(-n);
    const upDays = recent.filter((s) => s.status !== "down").length;
    return ((upDays / n) * 100).toFixed(3);
  }

  const apiNames = ["API 1", "API 2"];

  const mockData: ApiData[] = apiNames.map((name) => ({
    name,
    statuses: Array.from({ length: TOTAL_DAYS }, (_, i) => {
      new Date(start).toLocaleDateString() === new Date().toLocaleDateString() ? new Date(start) : new Date(start);
      new Date (end ).toLocaleDateString() === new Date().toLocaleDateString() ? new Date(end) : new Date(end);
      const tempDate = new Date(start);
      tempDate.setDate(start.getDate() + i);
      console.log("tempDate", tempDate.toLocaleDateString());
      return { date: tempDate, status: "default" };
    }),
    uptime15: "0.000",
    uptime30: "0.000",
    uptime60: "0.000",
    uptime90: "0.000",
  }));

  function updateUptime(api: ApiData) {
    api.uptime15 = calculateUptime(api, 15);
    api.uptime30 = calculateUptime(api, 30);
    api.uptime60 = calculateUptime(api, 60);
    api.uptime90 = calculateUptime(api, 90);
  }


  function updateApiStatus(apiName: string, date: string, status: StatusType) {
    const api = mockData.find((a) => a.name === apiName);
    if (!api) {
      console.log(`API ${apiName} not found`);
      return;
    }

    const entry = api.statuses.find(
      (s) => s.date.toLocaleDateString() === date
    );

    if (entry) {
      entry.status = status;
      updateUptime(api);
      console.log(`✅ Updated ${apiName} on ${date} -> ${status}`);
    } else {
      console.log(`❌ Date ${date} not found for ${apiName}`);
    }
  }


  // Example usage of the function
  updateApiStatus("API 1", "21/09/2025", "up");
  updateApiStatus("API 1", "22/09/2025", "up");

  const dayIndex = Math.floor(
    (today.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)
  );

  function getStartLabelForDays(days: number): string {
    const tempStart = new Date(end);
    tempStart.setDate(end.getDate() - days);
    return tempStart.toLocaleString(undefined, { month: "short", day: "numeric" });
  }

  const startLabel90 = getStartLabelForDays(89);
  const startLabel60 = getStartLabelForDays(59);
  const startLabel30 = getStartLabelForDays(29);
  const startLabel15 = getStartLabelForDays(14);
  const endLabel = end.toLocaleString(undefined, { month: "short", day: "numeric" });


  interface IncidentEntry {
        time: string;
        status: string;
        statusLabel: string;
        description: string;
      }

      interface Incident {
        title: string;
        entries: IncidentEntry[];
      }

      let incidents: Incident[] = [
        {
          title: "Elevated iDeal errors",
          entries: [
            {
              time: "Sep 22, 2025 20:14 UTC",
              status: "resolved",
              statusLabel: "Resolved",
              description: "From 13:05–19:15 UTC, we saw elevated errors on iDeal payments. This is now resolved.",
            },
            {
              time: "Sep 22, 2025 13:05 UTC",
              status: "investigating",
              statusLabel: "Investigating",
              description: "We are investigating reports of increased errors on iDeal payments.",
            },
          ],
        },
    
      ];


 interface Maintenance {
        status: string;
        statusLabel: string;
        service: string;
        time: string;
      }

      let maintenances: Maintenance[] = [
        {
      status: "inprogress",
      statusLabel: "In progress",
      service: "API",
      time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00 UTC",
        },
        {
      status: "inprogress",
      statusLabel: "In progress",
      service: "API",
      time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00 UTC",
        },
        {
      status: "inprogress",
      statusLabel: "In progress",
      service: "PayPal",
      time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00 UTC",
        },
      ];

    
      let statuses = [
        { title: "Global payments", description: "Checkout" },
        { title: "Revenue and finance automation", description: "Billing" },
        { title: "Store", description: "Domain" },
        { title: "oddin core components", description: "Dashboard, support, payouts, and webhooks" },
        { title: "Acquirers and payment methods", description: "Banks, card networks, and local payments" }
      ];

</script>

<div class="layout">
  {#each mockData as api, index}
    <section class="card" style="margin-bottom: 2px;">
      <div class="card-header">
        <div>{api.name}</div>
        <div class="uptimes">
          <span class="uptime15">{api.uptime15}% uptime</span>
          <span class="uptime30">{api.uptime30}% uptime</span>
          <span class="uptime60">{api.uptime60}% uptime</span>
          <span class="uptime90">{api.uptime90}% uptime</span>
        </div>
      </div>

      <div class="bar">
        {#each api.statuses as s, i}
          <div class="chip {s.status} {i === dayIndex ? s.status : ''}"></div>
        {/each}
      </div>

      <div class="timeline">
        <span class="label15">15 days ago</span>
        <span class="label30">30 days ago</span>
        <span class="label60">60 days ago</span>
        <span class="label90">90 days ago</span>
        <span>Today</span>
      </div>
    </section>
  {/each}

<div class="status-page">
  <!-- Left column -->
  <div class="left">
    <h3>System status</h3>

    {#each [
      ...statuses
    ] as status}
      <div class="status-card">
        <strong>{status.title}</strong>
        <p style="color: #666;">{status.description}</p>
      </div>
    {/each}
  </div>

  <!-- Right column -->
  <div class="right">
    <h3>Incidents</h3>
    {#each incidents as incident}
      <div class="incident-card">
        <h3>{incident.title}</h3>
        {#each incident.entries as entry}
          <div class="status-entry">
            <span class="time">{entry.time}</span>
            <span class="badge {entry.status}">{entry.statusLabel}</span>
            <p class="mt-2" style="font-size: 16px">{entry.description}</p>
          </div>
        {/each}
      </div>
    {/each}

    <h3>Maintenance</h3>
      <div class="maintenance-list">
        {#each maintenances as maintenance}
          <div class="maintenance-card">
        <div class="header">
          <span class="badge {maintenance.status}">{maintenance.statusLabel}</span>
          <span class="service">{maintenance.service}</span>
        </div>
        <div class="time lg:text-right">
          <time>{maintenance.time}</time>
        </div>
          </div>
        {/each}
      </div>
 
    </div>
    </div>
</div>


<style>
  :root{
    --bg:#FFFFFF;
    --text:#000000;
    --up:#4ce04c;
    --warn:#f2a900;
    --down:#f05d5e;
    --default:#e5e7eb;
    --chip-radius:1px;
    --today-ring:rgba(0,0,0,.25);
  }


  .status-page {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 50px;
    max-width: 1200px;
    margin: 40px auto;
    padding: 0 20px;
  }

  .status-card {
    padding: 15px;
    border-radius: 8px;
    margin-bottom: 15px;
    background: #fff;
    color: #222;
    box-shadow: rgba(0,0,0,0.05) 0 6px 24px, rgba(0,0,0,0.08) 0 0 0 1px;

  }

  .incident-card {
    background: #fdfdfd;
    color: black;
    padding: 30px;
    border-radius: 8px;
    margin-bottom: 20px;
    box-shadow: rgba(0, 0, 0, 0.05) 0px 0px 0px 1px;
  }


  .status-entry {
    margin-bottom: 25px;
  }

  .time {
    font-size: 0.9rem;
    color: #666;
    margin-right: 10px;
  }


  .badge.resolved {
    background: #4CAF50;
  }

  .badge.investigating {
    background: #607d8b;
  }

  .badge.inprogress {
    background: #1976d2;
  }

  .maintenance-list {
    display: flex;
    flex-direction: column;
    gap: 15px; 
  }

  .maintenance-card {
    padding: 15px;
    border-radius: 6px;
    color: black;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 5px;
  }

  .service {
    font-weight: 600;
  }


  .badge {
    font-size: 0.8rem;
    padding: 3px 8px;
    border-radius: 4px;
    color: white;
  }


  @media (max-width: 750px) {
      .status-page {
        grid-template-columns: 1fr;
      }
    }
 

  .layout{ padding: 10px; }

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
  .bar {
    display: flex;
    gap: 2px;           
    justify-content: space-between; 
  }
 
  .chip {
    flex: 1 1 0;           
    display: inline-block; 
    height: 24px;
    background: var(--default);
    border-radius: var(--chip-radius);
  
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