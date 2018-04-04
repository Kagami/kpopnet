INSERT INTO faces (rectangle, descriptor, image_id, idol_id, idol_confirmed)
VALUES            ($1,        $2,         $3,       $4,      TRUE)
ON CONFLICT (image_id, idol_id) DO NOTHING
