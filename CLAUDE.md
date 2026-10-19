# Beacon

Frontend monitoring platform: TS SDK → Go ingestion API → Postgres → Vue dashboard.

## Rules
- I write the code; explain and review, don't write whole features unless asked.
- Go: stdlib net/http, pgx v5, sqlc, goose, slog. No ORMs, no web frameworks.
- Vue: <script setup> + TS, composables in src/composables, TanStack Query for server state.
- Every performance change must be measured before/after and logged in docs/performance.md.

## Commands
- `docker compose up -d`: Postgres
- `cd apps/api && go run ./cmd/api`: API on :8080
- `pnpm --filter dashboard dev`: dashboard on :5173
