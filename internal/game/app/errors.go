package app

import "errors"

var (
	ErrUnknownQuestion = errors.New("unknown question")
	ErrUnknownPerson   = errors.New("unknown person")
	ErrClosed          = errors.New("puzzle closed")
)
