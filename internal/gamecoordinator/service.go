package gamecoordinator

import (
	"github.com/Cludch/csgo-tools/internal/domain/match"
)

type Service struct {
	matchService match.UseCase
	gc           *GC
}

func NewService(m match.UseCase) *Service {
	return &Service{
		matchService: m,
	}
}
