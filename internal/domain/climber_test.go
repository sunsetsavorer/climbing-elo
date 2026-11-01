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
	if c.BestFlash == nil || *c.BestFlash != Grade6A {
		t.Errorf("BestFlash not updated after flash")
	}
	if ascent.CreatedAt != now {
		t.Errorf("CreatedAt = %v, want %v", ascent.CreatedAt, now)
	}
}

func TestRecordAscentFailDoesNotUpdateBestFlash(t *testing.T) {
	c := newClimber()

	_, err := c.RecordAscent(Grade6A, OutcomeFail, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if c.BestFlash != nil {
		t.Errorf("BestFlash should remain nil after fail, got %v", c.BestFlash)
	}
}

func TestRecordAscentBestFlashUpdatedOnHarderGrade(t *testing.T) {
	c := newClimber()

	c.RecordAscent(Grade6A, OutcomeFlash, now)
	c.RecordAscent(Grade6B, OutcomeFlash, now)

	if c.BestFlash == nil || *c.BestFlash != Grade6B {
		t.Errorf("BestFlash should be 6B, got %v", c.BestFlash)
	}
}

func TestRecordAscentBestFlashNotDowngraded(t *testing.T) {
	c := newClimber()

	c.RecordAscent(Grade6B, OutcomeFlash, now)
	c.RecordAscent(Grade6A, OutcomeFlash, now)

	if c.BestFlash == nil || *c.BestFlash != Grade6B {
		t.Errorf("BestFlash should stay 6B, got %v", c.BestFlash)
	}
}
