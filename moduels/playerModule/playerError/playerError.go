package playerError

import "errors"

var (
	ErrDuplicateUsername = errors.New("username has already exists")
	ErrDuplicateEmail    = errors.New("email has already exists")
)
