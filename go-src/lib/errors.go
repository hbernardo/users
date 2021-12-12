package lib

import "errors"

var (
	// ErrNotFound error
	ErrNotFound = errors.New("not found")
	// ErrPreconditionFailed error
	ErrPreconditionFailed = errors.New("precondition failed")
)
