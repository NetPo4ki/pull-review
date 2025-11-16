-- name: UpsertUser :one
INSERT INTO users (user_id, username, team_name, is_active)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET username  = EXCLUDED.username,
    team_name = EXCLUDED.team_name,
    is_active = EXCLUDED.is_active
RETURNING user_id, username, team_name, is_active;

-- name: SetIsActive :one
UPDATE users
SET is_active = $2
WHERE user_id = $1
RETURNING user_id, username, team_name, is_active;

-- name: GetUserByID :one
SELECT user_id, username, team_name, is_active
FROM users
WHERE user_id = $1;