CREATE TABLE accounts (
  id         BIGSERIAL PRIMARY KEY,
  owner      VARCHAR NOT NULL,
  balance    BIGINT NOT NULL,
  currency   VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON accounts (owner);
