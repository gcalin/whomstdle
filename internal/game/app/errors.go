package app

import "errors"

// The errors the use cases return when they refuse a request.
var (
	ErrUnknownQuestion = errors.New("unknown question")
	ErrUnknownPerson   = errors.New("unknown person")
	ErrClosed          = errors.New("puzzle closed")
)
