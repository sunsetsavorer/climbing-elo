package usecase

import (
	"context"
	"time"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type ClimberRepository interface {
	Save(ctx context.Context, climber *domain.Climber) error
	Update(ctx context.Context, climber *domain.Climber) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.Climber, error)
}

type AscentRepository interface {
	Save(ctx context.Context, ascent *domain.Ascent) error
}

type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type IDGenerator interface {
	Generate() string
}

type Clock interface {
	Now() time.Time
}
