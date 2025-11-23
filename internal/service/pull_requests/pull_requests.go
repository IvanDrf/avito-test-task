package pullrequests

import (
	"context"
	"database/sql"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	pullrequests "github.com/IvanDrf/avito-test-task/internal/repo/pull_requests"
	"github.com/IvanDrf/avito-test-task/internal/repo/teams"
	"github.com/IvanDrf/avito-test-task/internal/repo/users"
	repo_interface "github.com/IvanDrf/avito-test-task/internal/service/interface"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type PullRequestService struct {
	usersRepo        repo_interface.UsersRepo
	pullRequestsRepo repo_interface.PullRequestsRepo
	teamsRepo        repo_interface.TeamsRepo
}

func NewPullRequestService(db *sql.DB) *PullRequestService {
	return &PullRequestService{
		usersRepo:        users.NewUsersRepo(db),
		pullRequestsRepo: pullrequests.NewPullRequestsRepo(db),
		teamsRepo:        teams.NewTeamsRepo(db),
	}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, req *api.PostPullRequestCreateJSONBody) (*api.PullRequest, error) {
	alreadyExists, err := s.pullRequestsRepo.IsPullRequestExists(ctx, req.PullRequestId)
	if err != nil {
		return nil, errs.ErrCantCheckPullRequestExistence()
	}

	if alreadyExists {
		return nil, errs.ErrPullRequestAlreadyExists()
	}

	author, err := s.usersRepo.GetUser(ctx, req.AuthorId)
	if err != nil {
		return nil, errs.ErrCantFindPullRequestAuthor()
	}

	reviewers, err := s.findReviewersForAuthor(ctx, author.UserId, author.TeamName)
	if err != nil {
		return nil, err
	}

	pullRequest := &api.PullRequest{
		PullRequestId:     req.PullRequestId,
		PullRequestName:   req.PullRequestName,
		AuthorId:          req.AuthorId,
		AssignedReviewers: reviewers,
	}

	return s.pullRequestsRepo.CreatePullRequest(ctx, req.PullRequestId, pullRequest)
}

func (s *PullRequestService) findReviewersForAuthor(ctx context.Context, authorID, teamName string) ([]string, error) {
	team, err := s.teamsRepo.FindTeam(ctx, teamName)
	if err != nil {
		return nil, errs.ErrCantFindTeam()
	}

	reviewers := findReviewers(team.Members, authorID)
	return reviewers, nil
}

func findReviewers(members []api.TeamMember, authorID string) []string {
	reviewers := make([]string, 0, 2)

	for _, member := range members {
		if authorID != member.UserId && member.IsActive && len(reviewers) < 2 {
			reviewers = append(reviewers, member.UserId)
		}
	}

	return reviewers
}

func (s *PullRequestService) MergePullRequest(ctx context.Context, req *api.PostPullRequestMergeJSONBody) error {
	pullRequest, err := s.pullRequestsRepo.FindPullRequestByID(ctx, req.PullRequestId)
	if err != nil {
		return errs.ErrCantFindPullRequest()
	}

	if pullRequest.Status == api.PullRequestStatusMERGED {
		return nil
	}

	return s.pullRequestsRepo.ChangePullRequestStatus(ctx, req.PullRequestId, string(api.PullRequestShortStatusMERGED))
}

func (s *PullRequestService) GetUserPullRequests(ctx context.Context, userID string) ([]api.PullRequestShort, error) {
	return s.pullRequestsRepo.FindPullRequestByUserID(ctx, userID)
}

func (s *PullRequestService) ReassignReviewer(ctx context.Context, req *api.PostPullRequestReassignJSONBody) (*api.PullRequest, error) {
	pullRequest, err := s.pullRequestsRepo.FindPullRequestByID(ctx, req.PullRequestId)
	if err != nil {
		return nil, err
	}

	if pullRequest.Status == api.PullRequestStatusMERGED {
		return nil, errs.ErrPullRequestAlreadyMerged()
	}

	oldReviewer, err := s.usersRepo.GetUser(ctx, req.OldUserId)
	if err != nil {
		return nil, err
	}

	newReviewer, err := s.findNewReviewer(ctx, oldReviewer.TeamName, pullRequest.AuthorId, pullRequest.AssignedReviewers)
	if err != nil {
		return nil, err
	}

	return s.pullRequestsRepo.ChangeReviewer(ctx, req.PullRequestId, req.OldUserId, newReviewer)
}

func (s *PullRequestService) findNewReviewer(ctx context.Context, teamName, authorID string, assignedReviewers []string) (string, error) {
	team, err := s.teamsRepo.FindTeam(ctx, teamName)
	if err != nil {
		return "", err
	}

	ignoredReviewers := map[string]struct{}{
		authorID: {},
	}

	for i := range assignedReviewers {
		ignoredReviewers[assignedReviewers[i]] = struct{}{}
	}

	for i := range team.Members {
		if isReviewerNotBusy(team.Members[i], ignoredReviewers) {
			return team.Members[i].UserId, nil
		}
	}

	return "", errs.ErrCantFindReviewers()
}

func isReviewerNotBusy(reviewer api.TeamMember, ignoredReviewers map[string]struct{}) bool {
	_, ok := ignoredReviewers[reviewer.UserId]
	return !ok && reviewer.IsActive
}
