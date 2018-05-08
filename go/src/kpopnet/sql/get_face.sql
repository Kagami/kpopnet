SELECT rectangle, idol_id, idol_confirmed FROM faces
WHERE image_id = $1
LIMIT 1
