// Holds the secondary (driven) adapters that implement the Access bounded
// context's repository port against PostgreSQL
// (library-docs/05-architecture/decisions/records/ADR-002-hexagonal-modular-monolith.md).

import pg from 'pg';

const { Pool } = pg;

// Opens a connection pool and verifies connectivity with a bounded timeout.
export async function newPool(connectionString) {
  const pool = new Pool({ connectionString });

  const client = await pool.connect();
  try {
    await client.query('SELECT 1');
  } finally {
    client.release();
  }

  return pool;
}
