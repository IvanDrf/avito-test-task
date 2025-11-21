package repo

import "github.com/IvanDrf/avito-test-task/pkg/api"

type TeamRepo interface {
	CreateTeam(team *api.Team) error
	FindTeam(teamName string) (*api.Team, error)
}

type UserRepo interface {
	ChangeUserActivity(userID string, isActive bool) error
}

type PullRequestRepo interface {
	CreatePullRequest(authorID string, pullRequestID string, name string) (*api.PullRequest, error)
	FindPullRequest(userID string) ([]api.PullRequestShort, error)
	ChangePullRequestStatus(pullRequestID string, status string) error
	ChangeReviewer(pullRequestID string, oldUserID string, newUserID string) (*api.PullRequest, error)
}
