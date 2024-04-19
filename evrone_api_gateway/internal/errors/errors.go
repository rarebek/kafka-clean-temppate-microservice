package errors

import (
	"errors"
)

var (
	ErrorConflict       = NewErrConflict("object")
	ErrorNotFound       = NewErrNotFound("object")
	ErrorInvalidOTPCode = errors.New("code is invalid")
	ErrorOTPExpired     = errors.New("one time password has expired")
)

// error not found
type ErrNotFound struct {
	name string
}

func (e *ErrNotFound) Error() string {
	return e.name + " not found"
}

func NewErrNotFound(text string) *ErrNotFound {
	return &ErrNotFound{text}
}

// error conflict
type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return e.name + " already exist"
}

func NewErrConflict(text string) *ErrConflict {
	return &ErrConflict{text}
}

// error validation
type ErrValidation struct {
	Err    error
	Errors map[string]string
}

func (e ErrValidation) Error() string {
	return e.Err.Error()
}

func NewErrValidation() *ErrValidation {
	return &ErrValidation{Errors: make(map[string]string)}
}

// error bad request
type ErrBadRequest struct {
	Err error
}

func (e ErrBadRequest) Error() string {
	return e.Err.Error()
}

func NewErrBadRequest(err error) *ErrBadRequest {
	return &ErrBadRequest{err}
}
