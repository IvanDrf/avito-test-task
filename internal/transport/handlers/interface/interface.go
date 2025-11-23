package service_interface

import (
	"context"

	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type Service interface {
	TeamService
	UserService
	PullRequestService
}

type TeamService interface {
	CreateTeam(ctx context.Context, team *api.Team) error
	GetTeam(ctx context.Context, teamName string) (*api.Team, error)
}

type UserService interface {
	ChangeUserActivity(ctx context.Context, user *api.PostUsersSetIsActiveJSONBody) error
}

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, req *api.PostPullRequestCreateJSONBody) (*api.PullRequest, error)
	GetUserPullRequests(ctx context.Context, userID string) ([]api.PullRequestShort, error)
	MergePullRequest(ctx context.Context, req *api.PostPullRequestMergeJSONBody) error
	ReassignReviewer(ctx context.Context, req *api.PostPullRequestReassignJSONBody) (*api.PullRequest, error)
}
