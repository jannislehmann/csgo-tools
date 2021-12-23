package player

import (
	"errors"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) CreatePlayer(id uint64) (*Player, error) {
	p, _ := NewPlayer(id)
	return p, s.repo.Create(p)
}

func (s *Service) GetAll() ([]*Player, error) {
	return s.repo.List()
}

func (s *Service) GetPlayer(id uint64) (*Player, error) {
	p, err := s.repo.Find(id)
	if p == nil {
		return s.CreatePlayer(id)
	}

	return p, err
}

func (s *Service) GetResult(p *Player, matchId entity.ID) (*PlayerResult, error) {
	for _, result := range p.Results {
		if result.MatchID == matchId {
			return result, nil
		}
	}

	return nil, entity.ErrNotFound
}

func (s *Service) AddResult(p *Player, r *PlayerResult) error {
	matchId := r.MatchID

	// Delete old result.
	dbResult, err := s.GetResult(p, matchId)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return err
	}

	if dbResult != nil {
		err = s.DeleteResult(p, matchId)
		if err != nil {
			return err
		}
	}

	p.Results = append(p.Results, r)
	return s.repo.AddResult(p, r)
}

func (s *Service) DeleteResult(p *Player, matchId entity.ID) error {
	return s.repo.DeleteResult(p, matchId)
}
