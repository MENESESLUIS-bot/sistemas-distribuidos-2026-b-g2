// Allows the configured origin (the shell/portal) to call the API directly
// during local development. In production, nginx serves both, so this is
// mostly a local-dev convenience (library-docs/05-architecture/cross-cutting.md).

export function cors(allowedOrigin) {
  return (req, res, next) => {
    res.setHeader('Access-Control-Allow-Origin', allowedOrigin);
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PATCH, DELETE, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');

    if (req.method === 'OPTIONS') {
      res.status(204).end();
      return;
    }
    next();
  };
}
