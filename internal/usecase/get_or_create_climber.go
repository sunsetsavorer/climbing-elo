package usecase

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type GetOrCreateClimber struct {
	idGenerator       IDGenerator
	clock             Clock
	climberRepository ClimberRepository
}

func NewGetOrCreateClimber(
	idGenerator IDGenerator,
	clock Clock,
	climberRepository ClimberRepository,
) (*GetOrCreateClimber, error) {
	if idGenerator == nil {
		return nil, ErrIDGeneratorRequired
	}
	if clock == nil {
		return nil, ErrClockRequired
	}
	if climberRepository == nil {
		return nil, ErrClimberRepositoryRequired
	}
	return &GetOrCreateClimber{
		idGenerator:       idGenerator,
		clock:             clock,
		climberRepository: climberRepository,
	}, nil
}

type GetOrCreateClimberInput struct {
	Nickname   string
	TelegramID int64
}

type GetOrCreateClimberOutput struct {
	Climber *domain.Climber
	Created bool
}

func (uc *GetOrCreateClimber) Exec(ctx context.Context, in GetOrCreateClimberInput) (*GetOrCreateClimberOutput, error) {
	existing, err := uc.climberRepository.GetByTelegramID(ctx, in.TelegramID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return &GetOrCreateClimberOutput{Climber: existing, Created: false}, nil
	}

	climber := domain.NewClimber(
		uc.idGenerator.Generate(),
		in.TelegramID,
		in.Nickname,
		uc.clock.Now(),
	)

	if err := uc.climberRepository.Save(ctx, climber); err != nil {
		return nil, err
	}

	return &GetOrCreateClimberOutput{Climber: climber, Created: true}, nil
}
