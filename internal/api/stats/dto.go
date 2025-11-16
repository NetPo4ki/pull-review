package stats

type UserAssignmentDTO struct {
	UserID        string `json:"user_id"`
	AssignedCount int32  `json:"assigned_count"`
}

type PRAssignmentDTO struct {
	PullRequestID string `json:"pull_request_id"`
	ReviewerCount int32  `json:"reviewer_count"`
}

type StatsResponse struct {
	Users        []UserAssignmentDTO `json:"users"`
	PullRequests []PRAssignmentDTO   `json:"pull_requests"`
}
