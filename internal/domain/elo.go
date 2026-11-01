package domain

import "math"

const (
	kFactorNew      = 40
	kFactorEstabled = 20
	kThreshold      = 30
)

func kFactor(ascentsCount int) float64 {
	if ascentsCount < kThreshold {
		return kFactorNew
	}
	return kFactorEstabled
}

func expected(climberRating, gradeRating int) float64 {
	return 1 / (1 + math.Pow(10, float64(gradeRating-climberRating)/400))
}

func score(outcome Outcome) float64 {
	if outcome == OutcomeFlash {
		return 1
	}
	return 0
}

// CalcDelta returns the rating change for an ascent attempt.
func CalcDelta(climberRating, gradeRating, ascentsCount int, outcome Outcome) int {
	e := expected(climberRating, gradeRating)
	s := score(outcome)
	k := kFactor(ascentsCount)
	return int(math.Round(k * (s - e)))
}
