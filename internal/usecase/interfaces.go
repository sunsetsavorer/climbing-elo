package usecase

import (
	"context"
	"time"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type ClimberRepository interface {
	Save(ctx context.Context, climber *domain.Climber) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.Climber, error)
}

type IDGenerator interface {
	Generate() string
}

type Clock interface {
	Now() time.Time
}
