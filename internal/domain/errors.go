package domain

import "errors"

var (
	ErrUnknownGrade   = errors.New("unknown grade")
	ErrUnknownOutcome = errors.New("unknown outcome")
)
