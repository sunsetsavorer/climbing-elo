package repository

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/adapter/postgre/models"
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
	return r.db.WithContext(ctx).Save(toClimberModel(climber)).Error
}

func (r *ClimberRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.Climber, error) {
	var m models.Climber

	err := r.db.WithContext(ctx).
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
