INSERT INTO idols (id, band_id, data) VALUES ($1, $2, $3)
ON CONFLICT (id) DO
  UPDATE SET band_id = $2, data = $3
