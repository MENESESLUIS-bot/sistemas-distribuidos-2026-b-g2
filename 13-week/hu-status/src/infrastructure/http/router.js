// Builds the Express router with the base middleware stack, health endpoints,
// and the Access bounded context's public /auth/login route.
//
// No JWT-verification middleware here — this service issues tokens, it never
// validates one of its own (there is no protected route in this repo; a
// RequireAuth-style middleware belongs to the other services, not this one).

import express from 'express';
import { correlationId } from './middleware/correlationId.js';
import { cors } from './middleware/cors.js';
import { errorHandler } from './middleware/errorHandler.js';
import { createHealthHandler } from './handlers/healthHandler.js';

export function createRouter({ pool, corsOrigin, authHandler }) {
  const app = express();

  app.use(correlationId);
  app.use(cors(corsOrigin));
  app.use(express.json());

  const health = createHealthHandler(pool);
  app.get('/health', health.liveness);
  app.get('/health/ready', health.readiness);

  app.post('/api/v1/auth/login', authHandler.login); // HU-01 — public

  app.use(errorHandler);

  return app;
}
