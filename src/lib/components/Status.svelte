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

  const apiNames = ["API 1", "API 2", "API 3"];

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



  const Indicators = {
    Resolved: {
      badge: "resolved",
      statusLabel: "Resolved",
    },
    Investigating: {
      badge: "investigating",
      statusLabel: "Investigating",
    },
    Scheduled: {
      badge: "scheduled",
      statusLabel: "Scheduled",
    },
    Inprogress: {
      badge: "inprogress",
      statusLabel: "In Progress",
    },
    Completed: {
      badge: "completed",
      statusLabel: "Completed",
    },
  } as const;

  // Each value inside Indicators
  type Indicator = typeof Indicators[keyof typeof Indicators];

  interface IncidentEntry {
    time: string;
    description: string;
    status: Indicator; 
  }

  interface Incident {
    title: string;
    entries: IncidentEntry[];
  }

  interface Maintenance {
        status: Indicator;
        service: string;
        time: string;
      }


  let incidents: Incident[] = [
      {
        title: "Elevated iDeal errors",
        entries: [
        {
          time: "Sep 22, 2025 20:14 UTC",
          status: Indicators.Resolved,
          description: "From 13:05–19:15 UTC, we saw elevated errors on iDeal payments. This is now resolved.",
        },
        {
          time: "Sep 22, 2025 13:05 UTC",
          status: Indicators.Investigating,
          description: "We are investigating reports of increased errors on iDeal payments.",
        },
        {
          time: "Sep 22, 2025 12:45 UTC",
          status: Indicators.Inprogress,
          description: "We are investigating reports of increased errors on iDeal payments.",
        },
      ],
    },
    
  ];


  let maintenances: Maintenance[] = [
      {
        status: Indicators.Completed,
        service: "API",
        time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00",
      },
      {
        status: Indicators.Scheduled,
        service: "API",
        time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00",
      },
      {
        status: Indicators.Scheduled,
        service: "API",
        time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00",
      },
      {
        status: Indicators.Inprogress,
        service: "PayPal",
        time: "Sep 25, 2025 05:00 — Sep 25, 2025 07:00",
      },
  ];

    
  let statuses = [
        { title: "Global payments", description: "Checkout" },
        { title: "Revenue automation", description: "Billing" },
        { title: "Custom store", description: "Domain" },
        { title: "Core components", description: "Dashboard" },
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
        <div style="display: flex; align-items: center; gap: 10px;">
          
          <div>
        <strong>{status.title}</strong>
        <p style="color: #666;">{status.description}</p>
          </div>
        </div>
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
            <span class="time font-bold">{entry.time}</span>
            <span class="badge mt-1 {entry.status.badge}">{entry.status.statusLabel}</span>
            <p class="mt-2 text-gray-600" style="font-size: 16px">{entry.description}</p>
          </div>
        {/each}
      </div>
    {/each}

    <h3>Maintenance</h3>
      <div class="maintenance-list">
        {#each maintenances as maintenance}
        <div class="flex justify-between items-center p-3 gap-4 ">
        <!-- Badge -->
        <span
          class="inline-flex items-center px-2.5 badge2 py-1 rounded-full text-xs font-medium
                {maintenance.status.badge}">
          {maintenance.status.statusLabel}
        </span>

        <!-- Service + Time stacked -->
        <div class="flex flex-col text-left leading-tight">
          <span class="text-base font-semibold text-gray-900">
            {maintenance.service}
          </span>
          <time class="text-base text-gray-500">
            {maintenance.time}
          </time>
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
    --inactive: #6b7280;
    --inactive-service: #6b7280;
    --active-service: black;
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
    background:  #d7f7c2;
    color:  #006908;
    font-weight: 600;
    border: 1px solid  #a6eb84;
  }

  .badge.investigating {
    background: #ebeef1;
    color: #545969;
    font-weight: 600;
    border: 1px solid #d5dbe1;
  }


  .badge.inprogress {
    background: #fff4e5;
    color: #b45309;
    font-weight: 600;
    border: 1px solid #ffddb3;
  }


  .badge2.completed {
    background:  #d7f7c2;
    color:  #006908;
    font-weight: 600;
    border: 1px solid  #a6eb84;
  }

  .badge2.inprogress {
    background: #fff4e5;
    color: #b45309;
    font-weight: 600;
    border: 1px solid #ffddb3;
  }

  .badge2.scheduled {
    background: white;
    color: #4b5563;
    font-weight: 600;
    border: 1px solid #cecece;
  }

  .badge {
    font-size: 0.8rem;
    padding: 0px 6px;
    display: inline-block;    
    white-space: nowrap; 
    border-radius: 4px;
  }
  
  .badge2 {
    font-size: 0.8rem;
    padding: 2px 6px;
    display: inline-block;    
    white-space: nowrap; 
    border-radius: 4px;
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
    color: var(--inactive-service);
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