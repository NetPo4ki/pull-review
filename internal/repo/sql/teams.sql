-- name: CreateTeam :exec
INSERT INTO teams (team_name)
VALUES ($1);

-- name: GetTeamWithMembers :many
SELECT
  t.team_name,
  u.user_id,
  u.username,
  u.is_active
FROM teams t
LEFT JOIN users u ON u.team_name = t.team_name
WHERE t.team_name = $1
ORDER BY u.user_id;