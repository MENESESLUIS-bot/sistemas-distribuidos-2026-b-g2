// Standard JSON envelope helpers used by every HTTP handler — see
// library-docs/07-api/guidelines.md ("Error format") and
// library-docs/05-architecture/cross-cutting.md.

// Matches library-docs/07-api/contracts/openapi/_shared.yaml#/components/schemas/ErrorResponse.
export function sendJSON(res, status, body) {
  res.status(status).json(body);
}

export function sendError(res, status, code, message, correlationId) {
  sendJSON(res, status, { error: code, message, correlationId });
}
