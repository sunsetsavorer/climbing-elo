package idgenerator

import "github.com/google/uuid"

type UUIDGenerator struct{}

func New() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) Generate() string {
	return uuid.New().String()
}
