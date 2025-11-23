package pullrequests

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/IvanDrf/avito-test-task/internal/errs"
	"github.com/IvanDrf/avito-test-task/internal/repo"
	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type pullRequestsRepo struct {
	db *sql.DB
}

func NewPullRequestsRepo(db *sql.DB) *pullRequestsRepo {
	return &pullRequestsRepo{
		db: db,
	}
}

func (r *pullRequestsRepo) CreatePullRequest(ctx context.Context, id string, pullRequest *api.PullRequest) (*api.PullRequest, error) {
	tr, err := r.db.Begin()
	if err != nil {
		return nil, errs.ErrCantStartTransaction()
	}

	query := fmt.Sprintf("INSERT INTO %s (id, name, author_id, status, created_at) VALUES (?, ?, ?, ?, ?)", repo.PullRequestsTable)

	now := time.Now()
	_, err = tr.ExecContext(ctx, query,
		id,
		pullRequest.PullRequestName,
		pullRequest.AuthorId,
		"OPEN",
		now,
	)
	if err != nil {
		tr.Rollback()
		return nil, errs.ErrCantCreatePullRequest()
	}

	err = insertReviewers(ctx, tr, id, pullRequest)
	if err != nil {
		tr.Rollback()
		return nil, errs.ErrCantCreateReviewers()
	}

	pullRequest.Status = api.PullRequestStatusOPEN

	return &api.PullRequest{
		PullRequestId:     id,
		PullRequestName:   pullRequest.PullRequestName,
		AuthorId:          pullRequest.AuthorId,
		Status:            api.PullRequestStatusOPEN,
		CreatedAt:         &now,
		MergedAt:          nil,
		AssignedReviewers: pullRequest.AssignedReviewers,
	}, tr.Commit()
}

func insertReviewers(ctx context.Context, tr *sql.Tx, id string, pullRequest *api.PullRequest) error {
	if len(pullRequest.AssignedReviewers) == 0 {
		return nil
	}

	query := fmt.Sprintf("INSERT INTO %s (pr_id, first_reviewer_id, second_reviewer_id) VALUES (?, ?, ?)", repo.ReviewersTable)

	var first, second any = nil, nil
	if len(pullRequest.AssignedReviewers) >= 1 {
		first = pullRequest.AssignedReviewers[0]
	}

	if len(pullRequest.AssignedReviewers) >= 2 {
		second = pullRequest.AssignedReviewers[1]
	}

	_, err := tr.ExecContext(ctx, query, id, first, second)
	return err
}

func (r *pullRequestsRepo) FindPullRequestByUserID(ctx context.Context, userID string) ([]api.PullRequestShort, error) {
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.author_id, p.status
		FROM %s p
		JOIN %s r ON p.id = r.pr_id
		WHERE r.first_reviewer_id = ? OR r.second_reviewer_id = ?`,
		repo.PullRequestsTable, repo.ReviewersTable)

	rows, err := r.db.QueryContext(ctx, query, userID, userID)
	if err != nil {
		return nil, errs.ErrCantFindPullRequestAuthor()
	}
	defer rows.Close()

	pullRequests := []api.PullRequestShort{}

	for rows.Next() {
		pullRequest := api.PullRequestShort{}

		err = rows.Scan(&pullRequest.PullRequestId, &pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.Status)
		if err != nil {
			return nil, errs.ErrCantValidatePullRequest()
		}

		pullRequests = append(pullRequests, pullRequest)
	}

	return pullRequests, nil
}

func (r *pullRequestsRepo) FindPullRequestByID(ctx context.Context, pullRequestID string) (*api.PullRequest, error) {
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.author_id, p.status, p.created_at, p.merged_at,
		       r.first_reviewer_id, r.second_reviewer_id
		FROM %s p
		LEFT JOIN %s r ON p.id = r.pr_id
		WHERE p.id = ?`,
		repo.PullRequestsTable, repo.ReviewersTable)

	pullRequest := api.PullRequest{}
	createdAt := time.Time{}
	var mergedAt *time.Time = nil
	var firstReviewer, secondReviewer *string = nil, nil

	err := r.db.QueryRowContext(ctx, query, pullRequestID).Scan(
		&pullRequest.PullRequestId,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorId,
		&pullRequest.Status,
		&createdAt,
		&mergedAt,
		&firstReviewer,
		&secondReviewer,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrCantFindPullRequest()
		}
		return nil, err
	}

	pullRequest.CreatedAt = &createdAt
	pullRequest.MergedAt = mergedAt

	pullRequest.AssignedReviewers = []string{}
	if firstReviewer != nil {
		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, *firstReviewer)
	}
	if secondReviewer != nil {
		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, *secondReviewer)
	}

	return &pullRequest, nil
}

func (r *pullRequestsRepo) ChangePullRequestStatus(ctx context.Context, pullRequestID string, status string) error {
	var mergedAt any = nil
	if status == string(api.PullRequestShortStatusMERGED) {
		mergedAt = time.Now()
	}

	query := fmt.Sprintf(`
		UPDATE %s 
		SET status = ?, merged_at = ? 
		WHERE id = ?`,
		repo.PullRequestsTable)

	_, err := r.db.ExecContext(ctx, query, status, mergedAt, pullRequestID)
	if err != nil {
		return errs.ErrCantChangePullRequestStatus()
	}

	return nil
}

func (r *pullRequestsRepo) ChangeReviewer(ctx context.Context, pullRequestID string, oldUserID string, newUserID string) (*api.PullRequest, error) {
	tr, err := r.db.Begin()
	if err != nil {
		return nil, errs.ErrCantStartTransaction()
	}

	firstReviewer, secondReviewer, err := findReviewers(ctx, tr, pullRequestID)
	if err != nil {
		tr.Rollback()
		return nil, err
	}

	if !isUserReviewer(firstReviewer, secondReviewer, oldUserID) {
		tr.Rollback()
		return nil, errs.ErrCantAssignUserForPullRequest()
	}

	newFirst, newSecond := swapReviewers(firstReviewer, secondReviewer, oldUserID, newUserID)

	err = updateReviewers(ctx, tr, newFirst, newSecond, pullRequestID)
	if err != nil {
		tr.Rollback()
		return nil, errs.ErrCantUpdateReviewers()
	}

	pullRequest, err := findPullRequest(ctx, tr, newFirst, newSecond, pullRequestID)
	if err != nil {
		tr.Rollback()
		return nil, errs.ErrCantFindPullRequest()
	}

	return pullRequest, tr.Commit()
}

func findReviewers(ctx context.Context, tr *sql.Tx, pullRequestID string) (*string, *string, error) {
	var firstReviewer, secondReviewer *string = nil, nil

	query := fmt.Sprintf("SELECT first_reviewer_id, second_reviewer_id FROM %s WHERE pr_id = ?", repo.ReviewersTable)
	err := tr.QueryRowContext(ctx, query, pullRequestID).Scan(&firstReviewer, &secondReviewer)

	return firstReviewer, secondReviewer, err
}

func isUserReviewer(first, second *string, userID string) bool {
	return (first != nil && *first == userID) || (second != nil && *second == userID)
}

func swapReviewers(first, second *string, oldUserID, newUserID string) (any, any) {
	var newFirst, newSecond any = nil, nil
	if first != nil && *first == oldUserID {
		newFirst = newUserID
	} else {
		newFirst = first
	}

	if second != nil && *second == oldUserID {
		newSecond = newUserID
	} else {
		newSecond = second
	}

	return newFirst, newSecond
}

func updateReviewers(ctx context.Context, tr *sql.Tx, newFirst, newSecond any, pullRequestID string) error {
	query := fmt.Sprintf("UPDATE %s SET first_reviewer_id = ?, second_reviewer_id = ? WHERE pr_id = ?", repo.ReviewersTable)
	_, err := tr.ExecContext(ctx, query, newFirst, newSecond, pullRequestID)

	return err
}

func findPullRequest(ctx context.Context, tr *sql.Tx, newFirstReviewer, newSecondReviewer any, pullRequestID string) (*api.PullRequest, error) {
	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.author_id, p.status, p.created_at, p.merged_at
		FROM %s p
		WHERE p.id = ?`,
		repo.PullRequestsTable)

	pullRequest := api.PullRequest{}
	createdAt := time.Time{}
	mergedAt := &time.Time{}
	err := tr.QueryRowContext(ctx, query, pullRequestID).Scan(&pullRequest.PullRequestId, &pullRequest.PullRequestName, &pullRequest.AuthorId, &pullRequest.Status, &createdAt, &mergedAt)
	if err != nil {
		return nil, err
	}

	pullRequest.CreatedAt = &createdAt
	pullRequest.MergedAt = mergedAt

	pullRequest.AssignedReviewers = []string{}
	if newFirstReviewer != nil {
		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, newFirstReviewer.(string))
	}
	if newSecondReviewer != nil {
		pullRequest.AssignedReviewers = append(pullRequest.AssignedReviewers, newSecondReviewer.(string))
	}

	return &pullRequest, nil
}

func (r *pullRequestsRepo) IsPullRequestExists(ctx context.Context, pullRequestID string) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = ?)", repo.PullRequestsTable)

	exists := false
	err := r.db.QueryRowContext(ctx, query, pullRequestID).Scan(&exists)
	if err != nil {
		return exists, errs.ErrCantFindPullRequest()
	}

	return exists, nil
}
