package usecase

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
)

type (
	GetLeaderboard struct {
		climberRepository ClimberRepository
	}

	GetLeaderboardDeps struct {
		ClimberRepository ClimberRepository
	}

	GetLeaderboardInput struct {
		Limit int
	}

	GetLeaderboardOutput struct {
		Climbers []*domain.Climber // ordered by rank, best first
	}
)

func NewGetLeaderboard(deps GetLeaderboardDeps) (*GetLeaderboard, error) {
	if deps.ClimberRepository == nil {
		return nil, ErrClimberRepositoryRequired
	}
	return &GetLeaderboard{
		climberRepository: deps.ClimberRepository,
	}, nil
}

func (uc *GetLeaderboard) Exec(ctx context.Context, in GetLeaderboardInput) (*GetLeaderboardOutput, error) {
	if in.Limit <= 0 {
		return nil, ErrInvalidLimit
	}

	climbers, err := uc.climberRepository.ListTopByRating(ctx, in.Limit)
	if err != nil {
		return nil, err
	}

	return &GetLeaderboardOutput{Climbers: climbers}, nil
}
