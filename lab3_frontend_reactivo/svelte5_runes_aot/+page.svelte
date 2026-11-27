<script lang="ts">
  interface CellTelemetryDTO {
    cellId: string;
    latencyMs: number;
    packetLossPct: number;
    throughputMbps: number;
  }

  const WINDOW_SIZE = 20;

  let packets = $state<CellTelemetryDTO[]>([
    { cellId: "5G-NR-VE-01", latencyMs: 14.2, packetLossPct: 0.05, throughputMbps: 920.5 },
    { cellId: "5G-NR-VE-02", latencyMs: 62.8, packetLossPct: 1.40, throughputMbps: 410.0 },
    { cellId: "5G-NR-VE-03", latencyMs: 21.5, packetLossPct: 0.12, throughputMbps: 845.2 },
    { cellId: "5G-NR-VE-04", latencyMs: 55.1, packetLossPct: 0.85, throughputMbps: 512.8 }
  ]);

  let lastUpdateMs = $state<number>(0);
  let burstCount = $state<number>(0);

  let metrics = $derived.by(() => {
    if (packets.length === 0) {
      return { avgLatencyMs: 0, criticalCells: 0, avgThroughputMbps: 0 };
    }
    let sumLatency = 0;
    let sumThroughput = 0;
    let criticalCells = 0;

    for (const pkt of packets) {
      sumLatency += pkt.latencyMs;
      sumThroughput += pkt.throughputMbps;
      if (pkt.latencyMs > 50) {
        criticalCells += 1;
      }
    }

    return {
      avgLatencyMs: Number((sumLatency / packets.length).toFixed(2)),
      criticalCells,
      avgThroughputMbps: Number((sumThroughput / packets.length).toFixed(1))
    };
  });

  function injectBurst(size = 10) {
    const start = performance.now();
    const nextBatch: CellTelemetryDTO[] = Array.from({ length: size }, (_, idx) => {
      const seq = burstCount * size + idx + 1;
      const isSpike = seq % 3 === 0;
      return {
        cellId: `5G-NR-VE-${String((seq % 16) + 1).padStart(2, "0")}`,
        latencyMs: Number((isSpike ? 52 + (seq % 45) : 10 + (seq % 35)).toFixed(1)),
        packetLossPct: Number((isSpike ? 1.8 : 0.08).toFixed(2)),
        throughputMbps: Number((isSpike ? 360 + (seq % 180) : 780 + (seq % 320)).toFixed(1))
      };
    });

    packets = [...packets, ...nextBatch].slice(-WINDOW_SIZE);
    burstCount += 1;

    requestAnimationFrame(() => {
      lastUpdateMs = Number((performance.now() - start).toFixed(3));
    });
  }
</script>

<section aria-label="Telemetría 5G Svelte 5">
  <header>
    <h2>Monitor de Celdas 5G — Svelte 5 Runes (Compilación AOT)</h2>
    <button type="button" onclick={() => injectBurst(10)}>
      Inyectar ráfaga (10 paquetes)
    </button>
  </header>

  <div>
    <div>
      <span>Latencia promedio</span>
      <strong>{metrics.avgLatencyMs} ms</strong>
    </div>
    <div>
      <span>Celdas críticas (&gt; 50 ms)</span>
      <strong>{metrics.criticalCells}</strong>
    </div>
    <div>
      <span>Throughput prom.</span>
      <strong>{metrics.avgThroughputMbps} Mbps</strong>
    </div>
    <div>
      <span>Tiempo actualización (performance.now)</span>
      <strong>{lastUpdateMs} ms</strong>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th>Celda</th>
        <th>Latencia (ms)</th>
        <th>Pérdida (%)</th>
        <th>Throughput (Mbps)</th>
        <th>Estado</th>
      </tr>
    </thead>
    <tbody>
      {#each packets as pkt, i (pkt.cellId + "-" + i)}
        {@const critical = pkt.latencyMs > 50}
        <tr>
          <td>{pkt.cellId}</td>
          <td>{pkt.latencyMs}</td>
          <td>{pkt.packetLossPct}</td>
          <td>{pkt.throughputMbps}</td>
          <td>{critical ? "CRÍTICO" : "NORMAL"}</td>
        </tr>
      {/each}
    </tbody>
  </table>
</section>