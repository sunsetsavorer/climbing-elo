package usecase

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type RecordAscent struct {
	idGenerator       IDGenerator
	clock             Clock
	climberRepository ClimberRepository
	ascentRepository  AscentRepository
	txManager         TxManager
}

func NewRecordAscent(
	idGenerator IDGenerator,
	clock Clock,
	climberRepository ClimberRepository,
	ascentRepository AscentRepository,
	txManager TxManager,
) (*RecordAscent, error) {
	if idGenerator == nil {
		return nil, ErrIDGeneratorRequired
	}
	if clock == nil {
		return nil, ErrClockRequired
	}
	if climberRepository == nil {
		return nil, ErrClimberRepositoryRequired
	}
	if ascentRepository == nil {
		return nil, ErrAscentRepositoryRequired
	}
	if txManager == nil {
		return nil, ErrTxManagerRequired
	}
	return &RecordAscent{
		idGenerator:       idGenerator,
		clock:             clock,
		climberRepository: climberRepository,
		ascentRepository:  ascentRepository,
		txManager:         txManager,
	}, nil
}

type RecordAscentInput struct {
	TelegramID int64
	Grade      domain.Grade
	Outcome    domain.Outcome
}

type RecordAscentOutput struct {
	Ascent  *domain.Ascent
	Climber *domain.Climber
}

func (uc *RecordAscent) Exec(ctx context.Context, in RecordAscentInput) (*RecordAscentOutput, error) {
	climber, err := uc.climberRepository.GetByTelegramID(ctx, in.TelegramID)
	if err != nil {
		return nil, err
	}
	if climber == nil {
		return nil, domain.ErrClimberNotFound
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
