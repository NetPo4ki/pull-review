-- name: CreatePR :one
INSERT INTO pull_requests (pr_id, name, author_id)
VALUES ($1, $2, $3)
RETURNING pr_id, name, author_id, status, created_at, merged_at;

-- name: InsertPRReviewer :exec
INSERT INTO pr_reviewers (pr_id, user_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetPR :one
SELECT pr_id, name, author_id, status, created_at, merged_at
FROM pull_requests
WHERE pr_id = $1;

-- name: GetPRReviewers :many
SELECT user_id
FROM pr_reviewers
WHERE pr_id = $1
ORDER BY user_id;

-- name: UpdatePRStatusIfOpen :one
UPDATE pull_requests
SET status = 'MERGED',
    merged_at = COALESCE(merged_at, now())
WHERE pr_id = $1 AND status = 'OPEN'
RETURNING pr_id, name, author_id, status, created_at, merged_at;

-- name: CandidatesForCreate :many
SELECT u.user_id
FROM users u
JOIN users a ON a.user_id = sqlc.arg(author_id)
WHERE u.team_name = a.team_name
  AND u.is_active = TRUE
  AND u.user_id <> a.user_id
ORDER BY random()
LIMIT sqlc.arg(limit_count)::int;

-- name: CandidateForReassign :one
SELECT u.user_id
FROM users u
JOIN users oldr ON oldr.user_id = sqlc.arg(old_user_id)
JOIN pull_requests pr ON pr.pr_id = sqlc.arg(pr_id)
WHERE u.team_name = oldr.team_name
  AND u.is_active = TRUE
  AND u.user_id <> pr.author_id
  AND u.user_id NOT IN (SELECT r.user_id FROM pr_reviewers r WHERE r.pr_id = pr.pr_id)
ORDER BY random()
LIMIT 1;

-- name: GetPRsForReviewer :many
SELECT pr_id, name, author_id, status, created_at, merged_at
FROM pull_requests
WHERE pr_id IN (
  SELECT pr_id FROM pr_reviewers WHERE user_id = $1
)
ORDER BY created_at DESC;