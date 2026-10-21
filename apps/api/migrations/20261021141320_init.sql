-- +goose Up
CREATE TABLE users (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  github_id   bigint UNIQUE,
  login       text NOT NULL,
  is_demo     boolean NOT NULL DEFAULT false,
  created_at  timestamptz NOT NULL DEFAULT now()
);


CREATE TABLE projects (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id    uuid NOT NULL REFERENCES users(id),
  name        text NOT NULL,
  public_key  text NOT NULL UNIQUE,         -- goes into the SDK, like Sentry's DSN
  created_at  timestamptz NOT NULL DEFAULT now()
);


CREATE TABLE issues (
  id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  project_id   uuid NOT NULL REFERENCES projects(id),
  fingerprint  text NOT NULL,
  title        text NOT NULL,
  first_seen   timestamptz NOT NULL,
  last_seen    timestamptz NOT NULL,
  event_count  bigint NOT NULL DEFAULT 0,
  status       text NOT NULL DEFAULT 'unresolved',
  UNIQUE (project_id, fingerprint)
);


-- High-volume table, partitioned by day
CREATE TABLE events (
  id           bigint GENERATED ALWAYS AS IDENTITY,
  project_id   uuid NOT NULL,
  issue_id     bigint,
  type         text NOT NULL,              -- 'error' | 'vital'
  name         text,                       -- vital name: LCP, INP, CLS...
  value        double precision,           -- vital value
  url_path     text,
  release      text,
  browser      text,
  payload      jsonb NOT NULL,
  received_at  timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (id, received_at)
) PARTITION BY RANGE (received_at);

-- +goose Down

DROP TABLE events; DROP TABLE issues; DROP TABLE projects; DROP TABLE users;
