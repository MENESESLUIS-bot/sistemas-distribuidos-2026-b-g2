// Implements the Access bounded context (library-docs/02-domain/domain-map.md):
// the single authenticated role that operates the system.

import bcrypt from 'bcryptjs';
import { ErrInvalidCredentials } from './errors.js';

// BCRYPT_COST matches NFR-004 (library-docs/04-requirements/non-functional.md): cost >= 12.
export const BCRYPT_COST = 12;

// Administrator is the Access bounded context's Aggregate Root.
// See library-docs/02-domain/entities-and-rules.md, Entity: Administrator.
export class Administrator {
  constructor({ id, username, passwordHash, createdAt, updatedAt }) {
    this.id = id;
    this.username = username;
    this.passwordHash = passwordHash;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;
  }

  // Verifies a plaintext password against the stored hash. Never distinguishes
  // "wrong username" from "wrong password" to the caller — callers must map any
  // error from this method (and a not-found lookup) to the same
  // ErrInvalidCredentials response.
  authenticate(plaintextPassword) {
    if (!bcrypt.compareSync(plaintextPassword, this.passwordHash)) {
      throw new ErrInvalidCredentials();
    }
  }
}

// Hashes the plaintext password (INV-001: password is always hashed; the
// plaintext is never persisted, logged, or returned).
export function newAdministrator(username, plaintextPassword) {
  if (!username) {
    throw new Error('username must not be empty');
  }
  if (!plaintextPassword) {
    throw new Error('password must not be empty');
  }

  const passwordHash = bcrypt.hashSync(plaintextPassword, BCRYPT_COST);
  const now = new Date();

  return new Administrator({
    id: undefined,
    username,
    passwordHash,
    createdAt: now,
    updatedAt: now,
  });
}
