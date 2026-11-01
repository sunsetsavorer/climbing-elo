package domain

import "time"

const (
	startingRating = 1000
	minRating      = 600
)

func NewClimber(
	id string,
	telegramID int64,
	nickname string,
	now time.Time,
) *Climber {
	return &Climber{
		ID:         id,
		TelegramID: telegramID,
		Nickname:   nickname,
		Rating:     startingRating,
		CreatedAt:  now,
	}
}

type Climber struct {
	ID           string
	TelegramID   int64
	Nickname     string
	Rating       int
	AscentsCount int
	HardestGrade    *Grade
	CreatedAt    time.Time
}

func (c *Climber) RecordAscent(grade Grade, outcome Outcome, now time.Time) (*Ascent, error) {
	if !grade.Valid() {
		return nil, ErrUnknownGrade
	}
	if !outcome.Valid() {
		return nil, ErrUnknownOutcome
	}

	gradeRating, _ := grade.Rating()
	delta := CalcDelta(c.Rating, gradeRating, c.AscentsCount, outcome)

	ascent := &Ascent{
		ClimberID: c.ID,
		Grade:     grade,
		Outcome:   outcome,
		Delta:     delta,
		CreatedAt: now,
	}

	newRating := c.Rating + delta
	if newRating < minRating {
		newRating = minRating
	}

	ascent.RatingAfter = newRating
	c.Rating = newRating
	c.AscentsCount++

	if outcome == OutcomeFlash {
		c.updateHardestGrade(grade, gradeRating)
	}

	return ascent, nil
}

func (c *Climber) updateHardestGrade(grade Grade, gradeRating int) {
	if c.HardestGrade == nil || gradeRating > mustRating(*c.HardestGrade) {
		c.HardestGrade = &grade
	}
}

func mustRating(g Grade) int {
	r, _ := g.Rating()
	return r
}
