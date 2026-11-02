package models

import "time"

type Ascent struct {
	ID          string    `gorm:"primaryKey"`
	ClimberID   string    `gorm:"index;not null"`
	Grade       string    `gorm:"not null"`
	Outcome     string    `gorm:"not null"`
	Delta       int       `gorm:"not null"`
	RatingAfter int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (Ascent) TableName() string { return "ascents" }
