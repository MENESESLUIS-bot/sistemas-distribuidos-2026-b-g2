// Loads configuration from environment variables — the only source of
// environment-specific configuration (library-docs/05-architecture/cross-cutting.md).

function getEnv(key, fallback) {
  const v = process.env[key];
  return v !== undefined && v !== '' ? v : fallback;
}

export function loadConfig() {
  const port = getEnv('PORT', '8080');

  const dbHost = getEnv('DB_HOST', 'localhost');
  const dbPort = getEnv('DB_PORT', '5432');
  const dbUser = getEnv('DB_USER', 'lms_user');
  const dbPassword = getEnv('DB_PASSWORD', 'lms_password');
  const dbName = getEnv('DB_NAME', 'lms_db');

  const jwtPrivateKeyPEM = getEnv('JWT_PRIVATE_KEY', '');
  if (!jwtPrivateKeyPEM) {
    throw new Error('JWT_PRIVATE_KEY must be set');
  }
  // .env stores a PEM as one line with literal \n — rules/2-anexos/C-api-hexagonal.md,
  // numeral 5.3.7: "el PEM en una sola línea, con \n escritos literalmente".
  const jwtPrivateKey = jwtPrivateKeyPEM.replace(/\\n/g, '\n');

  const jwtExpiry = getEnv('JWT_EXPIRY', '1h');

  return {
    port,
    dbHost,
    dbPort,
    dbUser,
    dbPassword,
    dbName,
    dsn: `postgres://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}`,
    // JWTPrivateKey signs every Administrator session token this service issues
    // (RS256). Only this service holds it — every other -api gets the matching
    // public key instead (JWT_PUBLIC_KEY), never this one
    // (rules/2-anexos/C-api-hexagonal.md, numeral 5.3.7).
    jwtPrivateKey,
    jwtExpiry,
    logLevel: getEnv('LOG_LEVEL', 'info'),
    corsOrigin: getEnv('CORS_ORIGIN', '*'),
  };
}
