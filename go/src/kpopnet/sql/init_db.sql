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

-- XXX(Kagami): Part of cutechan.
CREATE TABLE IF NOT EXISTS images (
  sha1 char(40) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS idol_previews (
  id uuid PRIMARY KEY REFERENCES idols ON DELETE CASCADE,
  image_id char(40) NOT NULL REFERENCES images ON DELETE CASCADE
);
