# ADR-C3-02: Three-AZ DDS replica set with point-in-time backups for loan_db

## Status

Proposed

## Context

Section 4 of the [architecture design note](architecture-design-note.md) names `loan_db` as the most critical data store of LMS-Library. It is the only record of which student has which copy and until when. ADR-008 and ADR-009 compute overdue penalties and suspensions from it, and `available_copies` in `catalog` can be checked against it, but nothing can rebuild it. Today `library-docs` promises only a daily backup, with RPO ≤ 24 hours and RTO under 30 minutes (`04-requirements/non-functional.md`, NFR-008). The MongoDB that `lms-circulation-db/deploy/compose.yml` migrates is a replica set `rs0` on one host. Losing a day of loans would leave books "on loan" with no due date, and students suspended or not suspended by mistake. The note sets RTO = 15 minutes and RPO = 5 minutes for this store.

## Decision

`loan_db` runs on Huawei Cloud DDS as a three-node replica set, with one node each in AZ1, AZ2 and AZ3, in the private data subnets. `lms-circulation-api` writes with `writeConcern: majority` and reads from the primary, so a loan is confirmed only when it is on nodes in at least two AZs. DDS automated backups stay on: one full snapshot per day plus incremental oplog backups every 5 minutes, kept 7 days in OBS. CBR copies the daily snapshot to a second region for a full regional loss, which is outside the 15/5-minute target. A runbook in `lms-circulation-db` and a drill in `lms-qa` before each release check both paths: automatic failover, and a point-in-time restore into a new instance followed by a CSMS secret update.

## Alternatives considered

- Keep the daily backup of NFR-008 (backup and restore only): rejected because a daily backup gives an RPO of up to 1,440 minutes, which can never meet the 5-minute target, and a full restore plus redeploy does not fit in 15 minutes either.
- Two-AZ replica set (primary and secondary in AZ1 and AZ2, no third zone): rejected because if the AZ holding two of the three votes fails, the node left has no majority, cannot elect itself primary, and `loan_db` stays read-only until someone acts, which breaks the RTO.
- Cross-region active-active with two writable clusters: rejected because MongoDB replica sets have a single primary, so two writers would need conflict resolution for the max-2-active-loans rule (INV-004), and the cost and complexity are far beyond a single-campus library.

## Consequences

- Positive: an AZ failure costs no confirmed loan (RPO 0 in that case) and less than one minute of write downtime, with no human action. Logical damage, for example a bad release deleting documents, can be undone to any point in the last 7 days with at most 5 minutes lost.
- Positive: the drill turns the RTO and RPO into numbers the team has measured, not only promised, and it reuses the chaos testing practice from week 15.
- Negative: three DDS nodes plus backup storage cost more than one container with a volume, and the AZ3 node exists mainly to vote.
- Negative: `writeConcern: majority` adds the cross-AZ round trip (a few milliseconds) to every `POST /api/v1/loans` and return. The p95 target of 300 ms still holds, but the team has to watch it in AOM.
