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
  image_id char(40) UNIQUE NOT NULL REFERENCES images
);

-- Don't reference images to be able to preload descriptors.
-- TODO(Kagami): Indexes!
CREATE TABLE IF NOT EXISTS faces (
  id bigserial PRIMARY KEY,
  rectangle box NOT NULL,
  descriptor bytea NOT NULL CHECK (octet_length(descriptor) = 512),
  image_id char(40) NOT NULL,
  idol_id uuid NOT NULL REFERENCES idols ON DELETE CASCADE,
  idol_confirmed boolean NOT NULL DEFAULT FALSE,
  source varchar(100) NOT NULL,
  UNIQUE (image_id, idol_id)
);

CREATE INDEX IF NOT EXISTS faces_idol_id on faces (idol_id);
