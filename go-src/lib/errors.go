package lib

import "errors"

var (
	// ErrNotFound represents not found error
	ErrNotFound = errors.New("not found")
	// ErrPreconditionFailed represents precondition failed error
	ErrPreconditionFailed = errors.New("precondition failed")
)
