package usecase

import (
	"context"
	"errors"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type (
	GetOrCreateClimber struct {
		idGenerator       IDGenerator
		clock             Clock
		climberRepository ClimberRepository
	}

	GetOrCreateClimberDeps struct {
		IDGenerator       IDGenerator
		Clock             Clock
		ClimberRepository ClimberRepository
	}

	GetOrCreateClimberInput struct {
		Nickname   string
		TelegramID int64
	}

	GetOrCreateClimberOutput struct {
		Climber *domain.Climber
		Created bool
	}
)

func NewGetOrCreateClimber(deps GetOrCreateClimberDeps) (*GetOrCreateClimber, error) {
	if deps.IDGenerator == nil {
		return nil, ErrIDGeneratorRequired
	}
	if deps.Clock == nil {
		return nil, ErrClockRequired
	}
	if deps.ClimberRepository == nil {
		return nil, ErrClimberRepositoryRequired
	}
	return &GetOrCreateClimber{
		idGenerator:       deps.IDGenerator,
		clock:             deps.Clock,
		climberRepository: deps.ClimberRepository,
	}, nil
}

func (uc *GetOrCreateClimber) Exec(ctx context.Context, in GetOrCreateClimberInput) (*GetOrCreateClimberOutput, error) {
	existing, err := uc.climberRepository.GetByTelegramID(ctx, in.TelegramID)
	if err != nil && !errors.Is(err, domain.ErrClimberNotFound) {
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
