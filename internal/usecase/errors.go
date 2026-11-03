package usecase

import "errors"

var (
	ErrIDGeneratorRequired       = errors.New("id generator required")
	ErrClockRequired             = errors.New("clock required")
	ErrClimberRepositoryRequired = errors.New("climber repository required")
	ErrAscentRepositoryRequired  = errors.New("ascent repository required")
	ErrTxManagerRequired         = errors.New("tx manager required")

	ErrInvalidLimit = errors.New("limit must be positive")
)
