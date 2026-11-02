package repository

import (
	"context"

	"github.com/sunsetsavorer/climbing-elo/internal/domain"
	"gorm.io/gorm"
)

type AscentRepository struct {
	db *gorm.DB
}

func NewAscentRepository(db *gorm.DB) *AscentRepository {
	return &AscentRepository{db: db}
}

func (r *AscentRepository) Save(ctx context.Context, ascent *domain.Ascent) error {
	return conn(ctx, r.db).Save(toAscentModel(ascent)).Error
}
