// Structured JSON logging — the Node equivalent of the Go service's zap logger.

import pino from 'pino';

export function createLogger(level) {
  return pino({ level });
}
