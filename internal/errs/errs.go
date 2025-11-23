package errs

import (
	"fmt"
	"net/http"

	"github.com/IvanDrf/avito-test-task/pkg/api"
)

type Error struct {
	Code    int
	ApiCode api.ErrorResponseErrorCode
	Msg     string
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %v, Api: %s Msg: %s", e.Code, e.ApiCode, e.Msg)
}

func ParseError(err error) (api.ErrorResponseErrorCode, string, int) {
	apiErr, ok := err.(Error)
	if !ok {
		return api.ErrorResponseErrorCode("Bad request"), "Bad request", http.StatusBadRequest
	}

	return apiErr.ApiCode, apiErr.Msg, apiErr.Code

}
