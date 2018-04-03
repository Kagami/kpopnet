INSERT INTO faces (descriptor, image_id, idol_id) VALUES ($1, $2, $3)
ON CONFLICT (image_id, idol_id) DO NOTHING
