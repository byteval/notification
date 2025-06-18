package notification

import "errors"

var (
	ErrNotFound  = errors.New("notification not found")
	ErrEmptyID   = errors.New("notification ID cannot be empty")
	ErrInvalidID = errors.New("invalid notification ID format")
)
