-- name: AssignmentsPerUser :many
SELECT
  user_id,
  COUNT(*)::int AS assigned_count
FROM pr_reviewers
GROUP BY user_id
ORDER BY assigned_count DESC, user_id;

-- name: AssignmentsPerPR :many
SELECT
  pr_id,
  COUNT(*)::int AS reviewer_count
FROM pr_reviewers
GROUP BY pr_id
ORDER BY pr_id;


