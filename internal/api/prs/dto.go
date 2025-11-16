package prs

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type PRDTO struct {
	PullRequestID   string   `json:"pull_request_id"`
	PullRequestName string   `json:"pull_request_name"`
	AuthorID        string   `json:"author_id"`
	Status          string   `json:"status"`
	Assigned        []string `json:"assigned_reviewers"`
	CreatedAt       *string  `json:"createdAt,omitempty"`
	MergedAt        *string  `json:"mergedAt,omitempty"`
}
type CreatePRResponse struct {
	PR PRDTO `json:"pr"`
}

type MergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}
type MergeResponse struct {
	PR PRDTO `json:"pr"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}
type ReassignResponse struct {
	PR         PRDTO  `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}
