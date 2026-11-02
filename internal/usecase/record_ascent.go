package usecase

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type (
	RecordAscent struct {
		idGenerator       IDGenerator
		clock             Clock
		climberRepository ClimberRepository
		ascentRepository  AscentRepository
		txManager         TxManager
	}

	RecordAscentDeps struct {
		IDGenerator       IDGenerator
		Clock             Clock
		ClimberRepository ClimberRepository
		AscentRepository  AscentRepository
		TxManager         TxManager
	}

	RecordAscentInput struct {
		TelegramID int64
		Grade      domain.Grade
		Outcome    domain.Outcome
	}

	RecordAscentOutput struct {
		Ascent  *domain.Ascent
		Climber *domain.Climber
	}
)

func NewRecordAscent(deps RecordAscentDeps) (*RecordAscent, error) {
	if deps.IDGenerator == nil {
		return nil, ErrIDGeneratorRequired
	}
	if deps.Clock == nil {
		return nil, ErrClockRequired
	}
	if deps.ClimberRepository == nil {
		return nil, ErrClimberRepositoryRequired
	}
	if deps.AscentRepository == nil {
		return nil, ErrAscentRepositoryRequired
	}
	if deps.TxManager == nil {
		return nil, ErrTxManagerRequired
	}
	return &RecordAscent{
		idGenerator:       deps.IDGenerator,
		clock:             deps.Clock,
		climberRepository: deps.ClimberRepository,
		ascentRepository:  deps.AscentRepository,
		txManager:         deps.TxManager,
	}, nil
}

func (uc *RecordAscent) Exec(ctx context.Context, in RecordAscentInput) (*RecordAscentOutput, error) {
	climber, err := uc.climberRepository.GetByTelegramID(ctx, in.TelegramID)
	if err != nil {
		return nil, err
	}

	ascent, err := climber.RecordAscent(in.Grade, in.Outcome, uc.clock.Now())
	if err != nil {
		return nil, err
	}
	ascent.ID = uc.idGenerator.Generate()

	if err := uc.txManager.RunInTx(ctx, func(ctx context.Context) error {
		if err := uc.ascentRepository.Save(ctx, ascent); err != nil {
			return err
		}
		return uc.climberRepository.Save(ctx, climber)
	}); err != nil {
		return nil, err
	}

	return &RecordAscentOutput{Ascent: ascent, Climber: climber}, nil
}
