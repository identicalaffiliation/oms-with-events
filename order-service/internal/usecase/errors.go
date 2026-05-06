package usecase

import "errors"

var (
	ErrInternal      = errors.New("internal server error")
	ErrInvalidUserId = errors.New("invalid user id")
)
