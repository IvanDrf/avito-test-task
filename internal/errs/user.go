package errs

import (
	"net/http"

	"github.com/IvanDrf/avito-test-task/pkg/api"
)

func ErrCantFindUser() error {
	return Error{
		Code:    http.StatusBadRequest,
		ApiCode: api.NOTFOUND,
		Msg:     "user not found",
	}
}

func ErrCantChangeUserActivity() error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "cant change user activity",
	}
}
