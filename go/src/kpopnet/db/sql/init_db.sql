CREATE TABLE IF NOT EXISTS bands (
  id uuid PRIMARY KEY,
  data jsonb NOT NULL
);

CREATE TABLE IF NOT EXISTS idols (
  id uuid PRIMARY KEY,
  band_id uuid NOT NULL REFERENCES bands ON DELETE CASCADE,
  data jsonb NOT NULL
);
