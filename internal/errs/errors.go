package errs

import "fmt"

type StatusError struct {
	Code    int
	Message string
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func (e StatusError) StatusCode() int {
	return e.Code
}

func BadRequest(message string) error {
	return StatusError{Code: 400, Message: message}
}

func Unauthorized(message string) error {
	return StatusError{Code: 401, Message: message}
}

func NotFound(message string) error {
	return StatusError{Code: 404, Message: message}
}

func Conflict(message string) error {
	return StatusError{Code: 409, Message: message}
}
