package errs

import (
	"net/http"

	"github.com/IvanDrf/avito-test-task/pkg/api"
)

func ErrCantCreatePullRequest() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant create new pull request",
	}
}

func ErrCantCreateReviewers() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant assign reviewers for pull request",
	}
}

func ErrCantFindPullRequestAuthor() error {
	return Error{
		Code:    http.StatusBadRequest,
		ApiCode: api.NOTFOUND,
		Msg:     "cant find pull request author",
	}
}

func ErrCantFindPullRequest() error {
	return Error{
		Code:    http.StatusBadRequest,
		ApiCode: api.NOTFOUND,
		Msg:     "cant find pull request",
	}
}

func ErrCantValidatePullRequest() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant validate pull request",
	}
}

func ErrCantChangePullRequestStatus() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant change pull request status",
	}
}

func ErrCantAssignUserForPullRequest() error {
	return Error{
		Code:    http.StatusInternalServerError,
		ApiCode: api.NOTASSIGNED,
		Msg:     "cant assign user for pull request",
	}
}

func ErrCantUpdateReviewers() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant update reviewers for pull request",
	}
}

func ErrCantFindReviewers() error {
	return Error{
		Code:    http.StatusConflict,
		ApiCode: api.NOCANDIDATE,
		Msg:     "cant find candidates for review",
	}
}

func ErrPullRequestAlreadyMerged() error {
	return Error{
		Code:    http.StatusConflict,
		ApiCode: api.ErrorResponseErrorCode(api.PullRequestShortStatusMERGED),
		Msg:     "pull request already merged",
	}
}

func ErrPullRequestAlreadyExists() error {
	return Error{
		Code:    http.StatusConflict,
		ApiCode: api.PREXISTS,
		Msg:     "pull request already exists",
	}
}

func ErrCantCheckPullRequestExistence() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant check for pull request existence",
	}
}
