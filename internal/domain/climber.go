package domain

import "time"

const StartingRating = 1000

type Climber struct {
	ID           string
	TelegramID   int64
	Nickname     string
	Rating       int
	AscentsCount int
	BestFlash    *Grade
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
		ClimberID:   c.ID,
		Grade:       grade,
		Outcome:     outcome,
		Delta:       delta,
		RatingAfter: c.Rating + delta,
		CreatedAt:   now,
	}

	c.Rating += delta
	c.AscentsCount++

	if outcome == OutcomeFlash {
		c.updateBestFlash(grade, gradeRating)
	}

	return ascent, nil
}

func (c *Climber) updateBestFlash(grade Grade, gradeRating int) {
	if c.BestFlash == nil || gradeRating > mustRating(*c.BestFlash) {
		c.BestFlash = &grade
	}
}

func mustRating(g Grade) int {
	r, _ := g.Rating()
	return r
}
