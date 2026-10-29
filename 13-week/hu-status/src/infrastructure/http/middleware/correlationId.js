// Assigns a correlation ID to every request (reusing an inbound one from nginx
// if present) and propagates it via req + the response header, so it can be
// included in every log line for that request.

import { randomUUID } from 'node:crypto';

export const CORRELATION_ID_HEADER = 'X-Correlation-Id';

export function correlationId(req, res, next) {
  const id = req.header(CORRELATION_ID_HEADER) || randomUUID();
  res.setHeader(CORRELATION_ID_HEADER, id);
  req.correlationId = id;
  next();
}
