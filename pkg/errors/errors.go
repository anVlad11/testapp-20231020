package errors

import "errors"

var (
	RequestBodyIsMalformed = errors.New("error.request_body_is_malformed")
	InternalError          = errors.New("error.internal_error")
)

func New(text string) error {
	return errors.New(text)
}
