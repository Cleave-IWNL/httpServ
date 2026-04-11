package apperror

import "net/http"

type ErrorResp struct {
	Status  int
	Message string
}

func (e ErrorResp) Error() string {
	return e.Message
}

var (
	ErrNotFound = ErrorResp{
		Status:  http.StatusNotFound,
		Message: "платеж не найден",
	}
	ErrAlreadyExists = ErrorResp{
		Status:  http.StatusForbidden,
		Message: "already exist",
	}
)
