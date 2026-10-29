// Implements the AdministratorRepository port against PostgreSQL — this is the
// only code allowed to touch the `access.administrator` table
// (library-docs/09-microservices/service-boundary-rules.md).
//
// Table/schema per `lms-access-db` (ADR-011 migration note): `access.administrator`,
// not the bare `administrators` in `public` the prior Go implementation used —
// corrected to match rules/A-db-postgres.md (rule 1: singular name; rule 9:
// schema named after the domain, nothing in `public`).

import { Administrator } from '../../domain/access/administrator.js';
import { ErrAdministratorNotFound } from '../../domain/access/errors.js';

export class AdministratorRepository {
  constructor(pool) {
    this.pool = pool;
  }

  async findByUsername(username) {
    const query = `
      SELECT id, username, password_hash, created_at, updated_at
      FROM access.administrator
      WHERE username = $1`;

    const { rows } = await this.pool.query(query, [username]);
    if (rows.length === 0) {
      throw new ErrAdministratorNotFound();
    }

    const row = rows[0];
    return new Administrator({
      id: row.id,
      username: row.username,
      passwordHash: row.password_hash,
      createdAt: row.created_at,
      updatedAt: row.updated_at,
    });
  }

  async save(administrator) {
    const query = `
      INSERT INTO access.administrator (id, username, password_hash, created_at, updated_at)
      VALUES (COALESCE(NULLIF($1, ''), gen_random_uuid()::text)::uuid, $2, $3, $4, $5)
      ON CONFLICT (id) DO UPDATE SET
        password_hash = EXCLUDED.password_hash,
        updated_at = EXCLUDED.updated_at`;

    await this.pool.query(query, [
      administrator.id ?? '',
      administrator.username,
      administrator.passwordHash,
      administrator.createdAt,
      administrator.updatedAt,
    ]);
  }
}
