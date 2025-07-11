-- name: GetUserByID :one
SELECT 
    id, 
    created_at, 
    updated_at, 
    email, 
    hashed_password,
    is_chirpy_red
FROM 
    users
WHERE 
    id = $1;