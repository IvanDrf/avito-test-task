package errs

import "net/http"

func ErrInvalidRequestBody() error {
	return Error{
		Code: http.StatusBadRequest,
		Msg:  "invalid body request",
	}
}
