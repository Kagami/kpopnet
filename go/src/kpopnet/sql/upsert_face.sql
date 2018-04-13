INSERT INTO faces (rectangle, descriptor, image_id, idol_id, idol_confirmed, source)
VALUES            ($1,        $2,         $3,       $4,      $5,             $6)
ON CONFLICT (image_id, idol_id) DO NOTHING
