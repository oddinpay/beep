<script lang="ts">

  const TOTAL_DAYS = 90;
  const today = new Date();

  // End is today
  const end = new Date(today);

  // Start is 90 days before end
  const start = new Date(end);
  start.setDate(end.getDate() - (TOTAL_DAYS - 1 ));
  
  // Types
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



  // Function to mark a specific date with a status
  function markStatus(dateString: string, status: StatusType) {
    const date = new Date(dateString);
    const index = statuses.findIndex((entry) => entry.date.toDateString() === date.toDateString());
    if (index !== -1) {
      statuses[index].status = status;
    }
  }

  // Example: mark specific days

  markStatus("2025-09-20", "warn");
  markStatus("2025-06-23", "down");



  let width: number = $state(0);

  let days: number = $derived(
    width === 0 ? 90 : width <= 310 ? 15 : width <= 600 ? 30 : width <= 900 ? 60 : 90
 );

  const dayIndex: number = Math.floor(
    (today.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)
  );

  let recent: StatusEntry[] = $derived(statuses.slice(-Math.min(days, TOTAL_DAYS)));

  let upDays: number = $derived(recent.filter((s) => s.status !== "down").length);

  let pct: string = $derived(((upDays / days) * 100).toFixed(3));

  let startLabel: string = $derived(
    start.toLocaleString(undefined, { month: "short", day: "numeric" })
  );

  let endLabel: string = $derived(
    end.toLocaleString(undefined, { month: "short", day: "numeric" })
  );
</script>

<div class="layout">
  <section class="card {days === 15 ? 'tiny' : ''}">
    <div class="card-header">
      <div>API</div>
      <div>{pct}% uptime</div>
    </div>

    <div class="bar">
      {#each recent as s, i}
      <div
        class="chip {s.status} {i === dayIndex ? s.status : ''}"
      ></div>
      {/each}
    </div>
    
    <div class="timeline">
      <span>{startLabel}</span>
      <span>{endLabel}</span>
    </div>
  </section>
</div>

<svelte:window bind:innerWidth={width} />

<style>
  :root {
    --bg: #0e2a3f;
    --text: #e9f2fb;
    --up: #4ce04c;
    --warn: #f2a900;
    --down: #f05d5e;
    --default: #ffffff; 
    --chip-radius: 2px;
  }

  .layout {
    padding: 5px;
  }

  .card {
    max-width: 1150px;
    margin: 40px auto;
    padding: 20px 32px;
    background: var(--bg);
    border-radius: 12px;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.25);
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

  .card-header > div:first-child {
    flex-shrink: 0; 
  }
  .card-header > div:last-child {
    white-space: nowrap; 
    flex-shrink: 0;
    font-variant-numeric: tabular-nums; 
    min-width: 140px; 
    text-align: right;
  }
  @media (max-width: 400px) {
    .card-header {
      flex-direction: column;
      align-items: flex-start;
    }
    .card-header > div:last-child {
      text-align: left;
      min-width: auto;
    }
  }
  /* --- END FIX --- */

  .bar {
    display: grid;
    grid-auto-flow: column;
    gap: 2px;
 }

  .chip {
    height: 25px;
    border-radius: var(--chip-radius);
    background: var(--default); 
  }

  .chip.warn {
    background: var(--warn);
  }
  .chip.down {
    background: var(--down);
  }
  .chip.up {
    background: var(--up); 
  }

  .timeline {
    display: flex;
    justify-content: space-between;
    margin-top: 8px;
    font-size: 0.85rem;
    color: var(--muted);
  }

  /* very small screens (15 days) */
  @media (max-width: 310px) {
    .card-header {
      font-size: 0.8rem;
    }
    .timeline {
      font-size: 0.7rem;
    }
    .bar {
      gap: 2px;
    }
    .chip {
      height: 14px;
    }
  }

  .card.tiny .card-header {
    font-size: 0.8rem;
  }
  .card.tiny .timeline {
    font-size: 0.7rem;
  }
  .card.tiny .bar {
    gap: 2px;
  }
  .card.tiny .chip {
    height: 14px;
  }
</style>

