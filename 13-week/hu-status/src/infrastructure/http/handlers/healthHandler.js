// Implements the /health and /health/ready endpoints
// (library-docs/07-api/contracts/openapi/library-api.yaml).

export function createHealthHandler(pool) {
  return {
    // Liveness — GET /health — responds 200 if the process is alive.
    liveness(_req, res) {
      res.status(200).json({ status: 'ok' });
    },

    // Readiness — GET /health/ready — responds 200 only if the database is
    // reachable within a 2s bound.
    async readiness(_req, res) {
      const ping = (async () => {
        const client = await pool.connect();
        try {
          await client.query('SELECT 1');
        } finally {
          client.release();
        }
      })();
      const timeout = new Promise((_, reject) =>
        setTimeout(() => reject(new Error('readiness timeout')), 2000),
      );

      try {
        await Promise.race([ping, timeout]);
        res.status(200).json({ status: 'ok', database: 'connected' });
      } catch {
        res.status(503).json({ status: 'unavailable', database: 'disconnected' });
      }
    },
  };
}
