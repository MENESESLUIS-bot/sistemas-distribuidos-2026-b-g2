# ADR-C3-01: One RDS for PostgreSQL HA instance with a schema per domain

## Status

Proposed

## Context

Section 2 of the [architecture design note](architecture-design-note.md) (row Cost) shows the problem. `06-data/migration-strategy.md` in `library-docs` says that every domain owns "its own database instance". On Huawei Cloud, each production-grade RDS for PostgreSQL instance is a primary/standby pair across two AZs. Applied literally, LMS-Library would pay for four pairs (`access_db`, `membership_db`, `catalog_db` and the `workflow_db` that ADR-008 adds) for a few thousand rows of administrators, students and books (`01-context/scope.md`, assumption #4). Each pair also needs its own backups, alarms, parameter group and failover test. A three-person team cannot run all of that. The team's own infrastructure already points the other way: `lms-infra-postgres/compose.yml` runs one PostgreSQL instance where every domain has its own schema and its own login user, and no `-db` repository defines a container anymore.

## Decision

Production and qa each get one RDS for PostgreSQL 16 instance in primary/standby mode (primary in AZ1, synchronous standby in AZ2). Inside it, every domain keeps the boundary that database per service needs: its own schema (`access`, `membership`, `catalog`, `workflow`), its own login user that can only read and write its schema, its own Liquibase changelog in its `-db` repository, and its own changelog tables. No API gets the credentials of another domain, and cross-schema foreign keys or joins are forbidden, as in `06-data/migration-strategy.md`. CSMS stores one secret per domain user. `loan_db` stays on DDS (ADR-005), so the system has two managed data services, not five.

## Alternatives considered

- One RDS HA instance per domain (the literal reading of the docs): rejected because it multiplies the database bill by four and the operating work (backups, alarms, failover drills, minor version upgrades) by four, while the data volume does not need it. The isolation it adds over separate schemas and users is only noisy-neighbour isolation, and at this load that does not matter.
- One shared schema with all tables in `public`, used by every service: rejected because it breaks database per service. Any API could join or write another domain's tables, and the course rule (`rules/2-anexos/A-db-postgres.md`, rule 9) forbids `public`.
- Self-managed PostgreSQL in a container on CCE with a persistent volume: rejected because the team would have to build replication, failover and point-in-time recovery itself, which RDS already provides.

## Consequences

- Positive: the cost and the operating work drop to one HA pair, and one backup policy (log backup every 5 minutes, 7-day retention) covers four domains. Ownership stays enforced by permissions and by separate Liquibase changelogs, not only by convention.
- Positive: the cloud design matches what `lms-infra-postgres` already runs locally, so the same migrations and connection settings work in every environment.
- Negative: the four domains share one failure domain and one maintenance window. An RDS failover or a heavy query in `catalog` affects `access` at the same time, and a domain cannot pick a different PostgreSQL version alone.
- Negative: if one domain later needs its own instance, the team has to move its schema with a dump and restore and change its secret. That migration is not free, so the decision should be reviewed if any schema passes 10 GB or needs a different SLA.
