package serror

import "fmt"

type Error struct {
	StatusCode int
	ErrMessage string
}

func (e Error) Error() string {
	return fmt.Sprintf(`{"statusCode": %d, "errMsaage": %s}`, e.StatusCode, e.ErrMessage)
}

func NewError(statusCode int, errMessage string) error {
	return Error{
		StatusCode: statusCode,
		ErrMessage: errMessage,
	}
}
