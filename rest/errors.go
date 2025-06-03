package rest

import (
	"errors"
	"fmt"
)

var (
	ErrAPIError = errors.New("api error")
	ErrClient   = fmt.Errorf("%w: client error", ErrAPIError)
	ErrServer   = fmt.Errorf("%w: server error", ErrAPIError)
)

var (
	ErrBadRequest = fmt.Errorf("%w: validation failed", ErrClient)
	ErrConflict   = fmt.Errorf("%w: conflict", ErrClient)
)

func IsAPIError(err error) bool {
	return errors.Is(err, ErrAPIError)
}

func IsClientError(err error) bool {
	return errors.Is(err, ErrClient)
}

func IsServerError(err error) bool {
	return errors.Is(err, ErrServer)
}

func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}
