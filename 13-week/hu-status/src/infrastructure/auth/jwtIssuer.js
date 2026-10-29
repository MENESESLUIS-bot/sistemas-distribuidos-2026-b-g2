// Implements the TokenIssuer driven port with JWT (library-docs/07-api/authentication.md).

import jwt from 'jsonwebtoken';

// Parses a Go-style duration string ("1h", "30m", "15s") into seconds — the
// same subset time.ParseDuration supported for JWT_EXPIRY in the original
// Go implementation.
function parseDurationSeconds(value) {
  const match = /^(\d+)(s|m|h)$/.exec(value);
  if (!match) {
    throw new Error(`invalid duration: ${value}`);
  }
  const amount = Number(match[1]);
  const unit = match[2];
  const unitSeconds = { s: 1, m: 60, h: 3600 }[unit];
  return amount * unitSeconds;
}

// JWTIssuer signs Administrator session tokens with RS256, per
// rules/2-anexos/C-api-hexagonal.md numeral 5.3.7: "Ningún servicio posee la
// clave privada" — only lms-access-api does. Every other service validates
// with the matching public key (JWT_PUBLIC_KEY), never this one.
export class JWTIssuer {
  constructor({ privateKey, expiry }) {
    this.privateKey = privateKey;
    this.expirySeconds = parseDurationSeconds(expiry);
  }

  // Issues a JWT with the Administrator's ID as `sub` — there is no `role`
  // claim, since v1 has a single role and no RBAC (02-domain/domain-map.md).
  async issue(administratorId) {
    const token = jwt.sign(
      { sub: administratorId },
      this.privateKey,
      { algorithm: 'RS256', expiresIn: this.expirySeconds },
    );
    return { token, expiresIn: this.expirySeconds };
  }
}
