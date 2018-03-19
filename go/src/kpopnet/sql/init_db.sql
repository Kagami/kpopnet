CREATE TABLE IF NOT EXISTS bands (
  id uuid PRIMARY KEY,
  data jsonb NOT NULL,
  CHECK (NOT(data ? 'id') AND data ? 'name')
);

CREATE TABLE IF NOT EXISTS idols (
  id uuid PRIMARY KEY,
  band_id uuid NOT NULL REFERENCES bands ON DELETE CASCADE,
  data jsonb NOT NULL,
  CHECK (NOT(data ? 'id') AND NOT(data ? 'band_id') AND data ? 'name')
);