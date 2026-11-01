package domain

import "time"

type Outcome string

const (
	OutcomeFlash Outcome = "flash"
	OutcomeFail  Outcome = "fail"
)

func (o Outcome) Valid() bool {
	return o == OutcomeFlash || o == OutcomeFail
}

type Ascent struct {
	ID          string
	ClimberID   string
	Grade       Grade
	Outcome     Outcome
	Delta       int
	RatingAfter int
	CreatedAt   time.Time
}
