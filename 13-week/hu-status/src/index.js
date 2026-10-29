// Entry point for access-service — see
// library-docs/09-microservices/services/02-access-service/README.md.

import { loadConfig } from './config/config.js';
import { createLogger } from './infrastructure/logger/logger.js';
import { newPool } from './infrastructure/postgres/db.js';
import { AdministratorRepository } from './infrastructure/postgres/administratorRepository.js';
import { JWTIssuer } from './infrastructure/auth/jwtIssuer.js';
import { Login } from './application/usecase/login.js';
import { createAuthHandler } from './infrastructure/http/handlers/authHandler.js';
import { createRouter } from './infrastructure/http/router.js';

async function run() {
  const cfg = loadConfig();
  const log = createLogger(cfg.logLevel);

  const pool = await newPool(cfg.dsn);

  const administrators = new AdministratorRepository(pool);
  const tokens = new JWTIssuer({ privateKey: cfg.jwtPrivateKey, expiry: cfg.jwtExpiry });
  const login = new Login({ administrators, tokens });
  const authHandler = createAuthHandler(login);

  const app = createRouter({ pool, corsOrigin: cfg.corsOrigin, authHandler });

  const server = app.listen(cfg.port, () => {
    log.info({ port: cfg.port }, 'access-service listening');
  });

  const shutdown = async (signal) => {
    log.info({ signal }, 'shutting down');
    server.close(async () => {
      await pool.end();
      process.exit(0);
    });
    // Force-exit if graceful shutdown doesn't finish in time.
    setTimeout(() => process.exit(1), 10_000).unref();
  };

  process.on('SIGINT', () => shutdown('SIGINT'));
  process.on('SIGTERM', () => shutdown('SIGTERM'));
}

run().catch((err) => {
  console.error('config error:', err.message);
  process.exit(1);
});
