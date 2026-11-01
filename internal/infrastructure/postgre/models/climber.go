package models

import "time"

type Climber struct {
	ID           string `gorm:"primaryKey"`
	TelegramID   int64  `gorm:"uniqueIndex;not null"`
	Nickname     string `gorm:"not null"`
	Rating       int    `gorm:"not null"`
	AscentsCount int    `gorm:"not null"`
	HardestGrade *string
	CreatedAt    time.Time `gorm:"not null"`
}

func (Climber) TableName() string { return "climbers" }
