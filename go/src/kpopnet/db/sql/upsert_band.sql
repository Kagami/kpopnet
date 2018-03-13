INSERT INTO bands (id, data) VALUES ($1, $2)
ON CONFLICT (id) DO
  UPDATE SET data = $2
