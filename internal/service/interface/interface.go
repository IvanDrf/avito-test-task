package repo_interface

import (
	"context"

	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type TeamsRepo interface {
	CreateTeam(ctx context.Context, teamdID string, team *api.Team) error
	FindTeam(ctx context.Context, teamName string) (*api.Team, error)
}

type UsersRepo interface {
	GetUser(ctx context.Context, userID string) (*api.User, error)
	ChangeUserActivity(ctx context.Context, userID string, isActive bool) error
}

type PullRequestsRepo interface {
	CreatePullRequest(ctx context.Context, id string, pullRequest *api.PullRequest) (*api.PullRequest, error)
	IsPullRequestExists(ctx context.Context, pullRequestID string) (bool, error)

	FindPullRequestByUserID(ctx context.Context, userID string) ([]api.PullRequestShort, error)
	FindPullRequestByID(ctx context.Context, pullRequestID string) (*api.PullRequest, error)

	ChangePullRequestStatus(ctx context.Context, pullRequestID string, status string) error
	ChangeReviewer(ctx context.Context, pullRequestID string, oldUserID string, newUserID string) (*api.PullRequest, error)
}
