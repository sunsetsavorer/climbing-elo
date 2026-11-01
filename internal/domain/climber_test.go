package domain

import (
	"errors"
	"testing"
	"time"
)

var now = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func newClimber() *Climber {
	return &Climber{
		ID:     "test-id",
		Rating: startingRating,
	}
}

func TestRecordAscentInvalidGrade(t *testing.T) {
	c := newClimber()
	_, err := c.RecordAscent(Grade("9A"), OutcomeFlash, now)
	if !errors.Is(err, ErrUnknownGrade) {
		t.Errorf("expected ErrUnknownGrade, got %v", err)
	}
}

func TestRecordAscentInvalidOutcome(t *testing.T) {
	c := newClimber()
	_, err := c.RecordAscent(Grade6A, Outcome("send"), now)
	if !errors.Is(err, ErrUnknownOutcome) {
		t.Errorf("expected ErrUnknownOutcome, got %v", err)
	}
}

func TestRecordAscentFlashUpdatesState(t *testing.T) {
	c := newClimber()
	before := c.Rating

	ascent, err := c.RecordAscent(Grade6A, OutcomeFlash, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ascent.Delta <= 0 {
		t.Errorf("flash delta should be positive, got %d", ascent.Delta)
	}
	if c.Rating != before+ascent.Delta {
		t.Errorf("Rating not updated correctly")
	}
	if c.AscentsCount != 1 {
		t.Errorf("AscentsCount = %d, want 1", c.AscentsCount)
	}
	if c.HardestGrade == nil || *c.HardestGrade != Grade6A {
		t.Errorf("HardestGrade not updated after flash")
	}
	if ascent.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", ascent.CreatedAt, now)
	}
}

func TestRecordAscentFailDoesNotUpdateHardestGrade(t *testing.T) {
	c := newClimber()

	_, err := c.RecordAscent(Grade6A, OutcomeFail, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.HardestGrade != nil {
		t.Errorf("HardestGrade should remain nil after fail, got %v", c.HardestGrade)
	}
}

func TestRecordAscentHardestGradeUpdatedOnHarderGrade(t *testing.T) {
	c := newClimber()

	c.RecordAscent(Grade6A, OutcomeFlash, now)
	c.RecordAscent(Grade6B, OutcomeFlash, now)

	if c.HardestGrade == nil || *c.HardestGrade != Grade6B {
		t.Errorf("HardestGrade should be 6B, got %v", c.HardestGrade)
	}
}

func TestRecordAscentRatingFloor(t *testing.T) {
	c := &Climber{ID: "test-id", Rating: minRating + 1}

	ascent, err := c.RecordAscent(Grade5C, OutcomeFail, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Rating != minRating {
		t.Errorf("rating = %d, want %d (clamped to floor)", c.Rating, minRating)
	}
	if ascent.RatingAfter != minRating {
		t.Errorf("ascent.RatingAfter = %d, want %d", ascent.RatingAfter, minRating)
	}
}

func TestRecordAscentHardestGradeNotDowngraded(t *testing.T) {
	c := newClimber()

	c.RecordAscent(Grade6B, OutcomeFlash, now)
	c.RecordAscent(Grade6A, OutcomeFlash, now)

	if c.HardestGrade == nil || *c.HardestGrade != Grade6B {
		t.Errorf("HardestGrade should stay 6B, got %v", c.HardestGrade)
	}
}
