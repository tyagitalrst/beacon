# 1. Postgres as the single datastore
Status: Accepted
## Context
We store relational data (users, projects, issues) and high-volume events.
## Decision
Use Postgres for both, with time-partitioned event tables and hourly rollups.
## Alternatives
- ClickHouse: better for analytics at scale, but no permanent free hosting and extra ops.
- Redis for rate limiting/queues: extra service; in-memory is enough at our scale.
## Consequences
Need partition maintenance + retention jobs. Revisit if > ~50M events/day.
