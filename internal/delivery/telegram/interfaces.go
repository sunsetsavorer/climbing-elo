package telegram

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/usecase"
)

type GetOrCreateClimberUseCase interface {
	Exec(ctx context.Context, in usecase.GetOrCreateClimberInput) (*usecase.GetOrCreateClimberOutput, error)
}

type RecordAscentUseCase interface {
	Exec(ctx context.Context, in usecase.RecordAscentInput) (*usecase.RecordAscentOutput, error)
}

type GetLeaderboardUseCase interface {
	Exec(ctx context.Context, in usecase.GetLeaderboardInput) (*usecase.GetLeaderboardOutput, error)
}
