package errs

import "net/http"

func ErrCantStartTransaction() error {
	return Error{
		Code: http.StatusInternalServerError,
		Msg:  "cant start transaction",
	}
}
