package telegram

import "errors"

var (
	ErrBotRequired                = errors.New("bot required")
	ErrGetOrCreateClimberRequired = errors.New("get or create climber use case required")
	ErrRecordAscentRequired       = errors.New("record ascent use case required")
)
