// Orchestrates the domain for HU-01 — it contains no business logic itself
// (that lives in the Aggregate), only coordination
// (library-docs/05-architecture/hexagonal-architecture.md).

import { ErrAdministratorNotFound, ErrInvalidCredentials } from '../../domain/access/errors.js';

// Login implements HU-01's acceptance criteria. Any failure — administrator
// not found, or wrong password — maps to the exact same ErrInvalidCredentials,
// per INV-002 on Administrator.
export class Login {
  constructor({ administrators, tokens }) {
    this.administrators = administrators;
    this.tokens = tokens;
  }

  async execute(username, password) {
    let admin;
    try {
      admin = await this.administrators.findByUsername(username);
    } catch (err) {
      if (err instanceof ErrAdministratorNotFound) {
        throw new ErrInvalidCredentials();
      }
      throw err;
    }

    admin.authenticate(password); // throws ErrInvalidCredentials on mismatch

    const { token, expiresIn } = await this.tokens.issue(admin.id);
    return { accessToken: token, expiresIn };
  }
}
