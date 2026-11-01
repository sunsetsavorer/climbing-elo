package domain

import (
	"errors"
	"testing"
)

func TestGradeValid(t *testing.T) {
	tests := []struct {
		grade Grade
		want  bool
	}{
		{Grade6A, true},
		{Grade7B, true},
		{Grade("8A"), false},
		{Grade(""), false},
	}

	for _, tt := range tests {
		if got := tt.grade.Valid(); got != tt.want {
			t.Errorf("Grade(%q).Valid() = %v, want %v", tt.grade, got, tt.want)
		}
	}
}

func TestGradeRating(t *testing.T) {
	tests := []struct {
		grade Grade
		want  int
	}{
		{Grade5C, 1000},
		{Grade6A, 1200},
		{Grade7B, 2000},
	}

	for _, tt := range tests {
		got, err := tt.grade.Rating()
		if err != nil {
			t.Fatalf("Grade(%q).Rating() unexpected error: %v", tt.grade, err)
		}
		if got != tt.want {
			t.Errorf("Grade(%q).Rating() = %d, want %d", tt.grade, got, tt.want)
		}
	}
}

func TestGradeRatingUnknown(t *testing.T) {
	_, err := Grade("9A").Rating()
	if !errors.Is(err, ErrUnknownGrade) {
		t.Errorf("expected ErrUnknownGrade, got %v", err)
	}
}
