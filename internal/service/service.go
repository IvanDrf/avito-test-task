package service

import "github.com/IvanDrf/avito-test-task/pkg/api"

type TeamService interface {
	AddTeam(team *api.Team) error
	GetTeam(teamName string) (*api.Team, error)
}

type UserService interface {
	ChangeUserActivity(user *api.PostUsersSetIsActiveJSONBody) error
}

type PullRequestService interface {
	CreatePullRequest(pullRequest *api.PostPullRequestCreateJSONBody) (*api.PullRequest, error)
	GetPullRequest(userID string) ([]api.PullRequestShort, error)
	MergedPullRequest(pullRequestID string) error
	ReassignReviewer(pullRequestID string, oldUserID string) (*api.PullRequest, error)
}
