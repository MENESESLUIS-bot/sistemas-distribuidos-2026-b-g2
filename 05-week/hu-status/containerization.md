# MVP 1 Containerization

## Runtime topology

`docker-compose.yml` runs the React frontend, the Access, Membership, and Catalog Go services, their PostgreSQL databases, and migration steps on the shared `lms-network` bridge network. Each PostgreSQL database uses a named volume so application restarts do not remove data.

## Required artifacts

| Artifact | Location | Purpose |
| --- | --- | --- |
| Access service Dockerfile | `MVP/lms-library-1.0.0/lms-library-1.0.0/access-service/Dockerfile` | Builds the Access service image. |
| Membership service Dockerfile | `MVP/lms-library-1.0.0/lms-library-1.0.0/membership-service/Dockerfile` | Builds the Membership service image. |
| Catalog service Dockerfile | `MVP/lms-library-1.0.0/lms-library-1.0.0/catalog-service/Dockerfile` | Builds the Catalog service image. |
| Frontend multi-stage Dockerfile | `MVP/lms-library-1.0.0/lms-library-1.0.0/frontend/Dockerfile` | Builds the React bundle and serves it with NGINX. |
| Compose topology | `MVP/lms-library-1.0.0/lms-library-1.0.0/docker-compose.yml` | Starts all services, database, health checks, network, and volume. |
| Build exclusions | Each service build context contains `.dockerignore` | Prevents local files, secrets, dependencies, and build output entering images. |

## Configuration and data

- Environment-specific values come from `.env` and `.env.example`; real secrets are never committed.
- Containers communicate by service name, never by `localhost`.
- Only the frontend port (`3000`) and NGINX gateway port (`8080`) are published to the host; service and database ports remain internal.
- PostgreSQL data is stored in the `access-db-data`, `membership-db-data`, and `catalog-db-data` named volumes.
- Health checks control startup order for the database and API.

## Verification

```bash
docker compose config
docker compose up -d --build
docker compose ps
curl http://localhost:8080/health
```

The remaining release evidence is a clean-environment build, a successful health check, and confirmation that data remains after restarting application containers.
