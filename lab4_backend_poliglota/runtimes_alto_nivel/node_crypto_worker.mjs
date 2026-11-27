import { Worker, isMainThread, parentPort, workerData } from "node:worker_threads";
import crypto from "node:crypto";
import { fileURLToPath } from "node:url";

if (!isMainThread) {
  const { payload, salt, iterations } = workerData;
  const derivedKey = crypto.pbkdf2Sync(payload, salt, iterations, 32, "sha256");
  parentPort.postMessage({
    hashHex: derivedKey.toString("hex"),
    iterations,
  });
} else {
  function runCryptoWorker(payload, salt, iterations) {
    return new Promise((resolve, reject) => {
      const start = performance.now();
      const worker = new Worker(fileURLToPath(import.meta.url), {
        workerData: { payload, salt, iterations },
      });

      worker.on("message", (msg) => {
        const elapsedMs = Number((performance.now() - start).toFixed(2));
        resolve({ ...msg, elapsedMs });
      });
      worker.on("error", reject);
      worker.on("exit", (code) => {
        if (code !== 0) {
          reject(new Error(`Worker finalizó con código ${code}`));
        }
      });
    });
  }

  let loopTicks = 0;
  const heartbeat = setInterval(() => {
    loopTicks += 1;
  }, 10);

  const jobs = [
    runCryptoWorker("tx-audit-01:+584121112233:5000", "salt-unefa-1", 150_000),
    runCryptoWorker("tx-audit-02:+584141112233:9500", "salt-unefa-2", 150_000),
  ];

  const results = await Promise.all(jobs);
  clearInterval(heartbeat);

  results.forEach((res, idx) => {
    console.log(
      `Job ${idx + 1} -> hash=${res.hashHex.slice(0, 16)}... | iteraciones=${res.iterations} | tiempo=${res.elapsedMs} ms`
    );
  });
  console.log(`Event loop libuv activo durante cálculos: ${loopTicks} ticks`);
}