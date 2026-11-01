package domain

import "testing"

const (
	SignMinus = "-"
	SignNone  = ""
	SignPlus  = "+"
)

func TestCalcDelta(t *testing.T) {
	tests := []struct {
		name          string
		climberRating int
		gradeRating   int
		ascentsCount  int
		outcome       Outcome
		wantSign      string
	}{
		{
			name:          "flash on equal rated grade gives positive delta",
			climberRating: 1200,
			gradeRating:   1200,
			ascentsCount:  0,
			outcome:       OutcomeFlash,
			wantSign:      SignPlus,
		},
		{
			name:          "fail on equal rated grade gives negative delta",
			climberRating: 1200,
			gradeRating:   1200,
			ascentsCount:  0,
			outcome:       OutcomeFail,
			wantSign:      SignMinus,
		},
		{
			name:          "flash on much higher grade gives large positive delta",
			climberRating: 1000,
			gradeRating:   2000,
			ascentsCount:  0,
			outcome:       OutcomeFlash,
			wantSign:      SignPlus,
		},
		{
			name:          "fail on much lower grade gives small negative delta",
			climberRating: 2000,
			gradeRating:   1000,
			ascentsCount:  0,
			outcome:       OutcomeFail,
			wantSign:      SignMinus,
		},
		{
			name:          "flash on much easier grade gives zero delta",
			climberRating: 2000,
			gradeRating:   1000,
			ascentsCount:  30,
			outcome:       OutcomeFlash,
			wantSign:      SignNone,
		},
		{
			name:          "fail on much harder grade gives zero delta",
			climberRating: 1000,
			gradeRating:   2000,
			ascentsCount:  0,
			outcome:       OutcomeFail,
			wantSign:      SignNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := CalcDelta(tt.climberRating, tt.gradeRating, tt.ascentsCount, tt.outcome)
			got := sign(delta)
			if got != tt.wantSign {
				t.Errorf("CalcDelta() = %d (sign %q), want sign %q", delta, got, tt.wantSign)
			}
		})
	}
}

func TestCalcDeltaKFactor(t *testing.T) {
	newDelta := CalcDelta(1200, 1200, 0, OutcomeFlash)
	estDelta := CalcDelta(1200, 1200, 30, OutcomeFlash)
	if newDelta <= estDelta {
		t.Errorf("new climber delta (%d) should be greater than established (%d)", newDelta, estDelta)
	}
}

func sign(n int) string {
	if n > 0 {
		return SignPlus
	}
	if n < 0 {
		return SignMinus
	}
	return SignNone
}
