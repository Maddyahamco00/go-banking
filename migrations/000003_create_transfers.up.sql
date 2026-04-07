CREATE TABLE transfers (
  id              BIGSERIAL PRIMARY KEY,
  from_account_id BIGINT NOT NULL REFERENCES accounts(id),
  to_account_id   BIGINT NOT NULL REFERENCES accounts(id),
  amount          BIGINT NOT NULL CHECK (amount > 0),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON transfers (from_account_id, to_account_id);