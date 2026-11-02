package repository

import (
	"github.com/sunsetsavorer/climbing-elo/internal/adapter/postgres/models"
	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

func toClimberModel(c *domain.Climber) *models.Climber {
	var hardestGrade *string
	if c.HardestGrade != nil {
		s := string(*c.HardestGrade)
		hardestGrade = &s
	}
	return &models.Climber{
		ID:           c.ID,
		TelegramID:   c.TelegramID,
		Nickname:     c.Nickname,
		Rating:       c.Rating,
		AscentsCount: c.AscentsCount,
		HardestGrade: hardestGrade,
		CreatedAt:    c.CreatedAt,
	}
}

func fromClimberModel(m *models.Climber) *domain.Climber {
	var hardestGrade *domain.Grade
	if m.HardestGrade != nil {
		g := domain.Grade(*m.HardestGrade)
		hardestGrade = &g
	}
	return &domain.Climber{
		ID:           m.ID,
		TelegramID:   m.TelegramID,
		Nickname:     m.Nickname,
		Rating:       m.Rating,
		AscentsCount: m.AscentsCount,
		HardestGrade: hardestGrade,
		CreatedAt:    m.CreatedAt,
	}
}

func toAscentModel(a *domain.Ascent) *models.Ascent {
	return &models.Ascent{
		ID:          a.ID,
		ClimberID:   a.ClimberID,
		Grade:       string(a.Grade),
		Outcome:     string(a.Outcome),
		Delta:       a.Delta,
		RatingAfter: a.RatingAfter,
		CreatedAt:   a.CreatedAt,
	}
}
