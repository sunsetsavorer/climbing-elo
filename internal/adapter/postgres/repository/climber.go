package repository

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/adapter/postgres/models"
	"github.com/sunsetsavorer/climbing-elo/internal/domain"
	"gorm.io/gorm"
)

type ClimberRepository struct {
	db *gorm.DB
}

func NewClimberRepository(db *gorm.DB) *ClimberRepository {
	return &ClimberRepository{db: db}
}

func (r *ClimberRepository) Save(ctx context.Context, climber *domain.Climber) error {
	return conn(ctx, r.db).Save(toClimberModel(climber)).Error
}

func (r *ClimberRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.Climber, error) {
	var m models.Climber

	err := conn(ctx, r.db).
		Where("telegram_id = ?", telegramID).
		First(&m).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, domain.ErrClimberNotFound
	}

	if err != nil {
		return nil, err
	}
	return fromClimberModel(&m), nil
}

// ListTopByRating returns climbers by rating descending; on a tie the one
// who registered earlier ranks higher.
func (r *ClimberRepository) ListTopByRating(ctx context.Context, limit int) ([]*domain.Climber, error) {
	var ms []models.Climber

	err := conn(ctx, r.db).
		Order("rating DESC").
		Order("created_at ASC").
		Limit(limit).
		Find(&ms).
		Error

	if err != nil {
		return nil, err
	}

	climbers := make([]*domain.Climber, len(ms))
	for i := range ms {
		climbers[i] = fromClimberModel(&ms[i])
	}
	return climbers, nil
}
