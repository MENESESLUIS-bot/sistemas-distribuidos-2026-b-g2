// Express error-handling middleware — maps unexpected failures (including a
// malformed JSON body from express.json()) to the standard error envelope.
// Domain/application errors expected during normal operation (e.g.
// ErrInvalidCredentials) are handled inline in each handler, not here.

import { sendError } from '../response.js';

export function errorHandler(err, req, res, _next) {
  if (err.type === 'entity.parse.failed') {
    sendError(res, 400, 'VALIDATION_ERROR', 'Invalid request body', req.correlationId);
    return;
  }

  sendError(res, 500, 'INTERNAL_ERROR', 'Internal server error', req.correlationId);
}
