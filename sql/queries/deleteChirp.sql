-- name: DeleteChirpByUserID :one
DELETE FROM chirps
WHERE 
    id = $1
    AND user_id = $2
RETURNING id;