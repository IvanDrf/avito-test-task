package errs

import (
	"net/http"

	"github.com/IvanDrf/avito-test-task/pkg/api"
)

func ErrCantCreateNewTeam() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant create new team",
	}
}

func ErrCantCreateTeamMembers() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant create team members",
	}
}

func ErrCantFindTeamMembers() error {
	return Error{
		Code:    http.StatusInternalServerError,
		ApiCode: api.NOTFOUND,
		Msg:     "cant find team members",
	}
}

func ErrCantFindTeam() error {
	return Error{
		Code:    http.StatusBadRequest,
		ApiCode: api.NOTFOUND,
		Msg:     "cant find team",
	}
}
