# 2. Modular monolith over microservices

Status: Accepted

## Context

Beacon is built by one developer with limited time, and it must run at $0/month
(Cloud Run + Neon free tiers). All data lives in one Postgres database (ADR 0001).
The backend has several parts: ingestion, dashboard API, auth, live stream (SSE),
and scheduled jobs. The dashboard is served by the API on the same origin so login
cookies work without a paid domain (ADR 0003).

## Decision

Build Beacon as a modular monolith: one Go module and one deployable API program
(cmd/api) that serves ingestion, the dashboard API, auth, the live stream and the
dashboard itself. Inside, code is split by domain (ingest, issues, vitals, auth,
live) plus platform for shared code. Domains only use each other's exported
functions, never each other's internals. Scheduled work runs as a separate
program (cmd/jobs).

## Alternatives

- Microservices (ingestion, dashboard API and auth as separate services):
  - Good: each part can scale and deploy on its own.
  - Not now: more Cloud Run services and DB connections to keep inside the free tier, network calls that can fail, tracing across services, service-to-service auth, and either a shared database (still coupled) or one database per service (data sync through queues). Login cookies would also need a gateway or a paid domain. These solve team-scaling problems a one-developer project doesn't have.

## Consequences

- Easier: one program to build, deploy and monitor; local dev is a single
  `go run ./cmd/api`; modules call each other as normal function calls, with no
  network failures in between.
- Harder: one failure domain. A crash takes down everything, and a traffic spike
  on ingestion can slow down the dashboard (noisy neighbor).
- When to split: move ingestion into its own service (same image, new cmd/ingest)
  when ingest load measurably hurts dashboard p95 latency, or when ingestion needs
  different scaling or limits. The domain layout keeps that a small change.
