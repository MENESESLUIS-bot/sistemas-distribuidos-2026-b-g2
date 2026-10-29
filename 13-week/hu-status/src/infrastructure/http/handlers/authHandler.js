// Implements HU-01 (library-docs/07-api/contracts/openapi/library-api.yaml, /auth/login).
// A primary (driving) adapter — it depends on the application layer, never the
// other way around (library-docs/05-architecture/hexagonal-architecture.md).

import { ErrInvalidCredentials } from '../../../domain/access/errors.js';
import { sendError, sendJSON } from '../response.js';

export function createAuthHandler(loginUseCase) {
  return {
    async login(req, res) {
      const { username, password } = req.body ?? {};
      if (typeof username !== 'string' || typeof password !== 'string') {
        sendError(res, 400, 'VALIDATION_ERROR', 'Invalid request body', req.correlationId);
        return;
      }

      try {
        const result = await loginUseCase.execute(username, password);
        sendJSON(res, 200, { accessToken: result.accessToken, expiresIn: result.expiresIn });
      } catch (err) {
        if (err instanceof ErrInvalidCredentials) {
          // INV-002: same generic error whether the username or the password was wrong.
          sendError(res, 401, 'INVALID_CREDENTIALS', 'Incorrect username or password', req.correlationId);
          return;
        }
        sendError(res, 500, 'INTERNAL_ERROR', 'Internal server error', req.correlationId);
      }
    },
  };
}
