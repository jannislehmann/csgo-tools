package user

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/Cludch/csgo-tools/pkg/valveapi"
)

type Service struct {
	repo                 Repository
	configurationService config.UseCase
}

func NewService(r Repository, c config.UseCase) *Service {
	return &Service{
		repo:                 r,
		configurationService: c,
	}
}

func (s *Service) GetUser(id entity.ID) (*User, error) {
	return s.repo.Find(id)
}

func (s *Service) GetUserBySteamId(id uint64) (*User, error) {
	return s.repo.FindBySteamId(id)
}
func (s *Service) GetUserByFaceitId(id entity.ID) (*User, error) {
	return s.repo.FindByFaceitId(id)
}

func (s *Service) GetUsersWithAuthenticationCode() ([]*User, error) {
	return s.repo.FindUsersContainingAuthenticationCode()
}

func (s *Service) CreateUser(u *User) error {
	dbUser, _ := s.GetUser(u.ID)
	if dbUser != nil {
		return errors.New("user with id already exists")
	}

	if u.Steam != nil {
		steamUser, _ := s.GetUserBySteamId(u.Steam.ID)
		if steamUser != nil {
			return errors.New("user with steam id already exists")
		}
	}

	if u.Faceit != nil {
		faceitUser, _ := s.GetUserByFaceitId(u.Faceit.ID)
		if faceitUser != nil {
			return errors.New("user with faceit id already exists")
		}
	}

	return s.repo.Create(u)
}

func (s *Service) AddSteamMatchHistoryAuthenticationCode(user *User, authCode string, sc string) error {
	shareCode, _ := share_code.Decode(sc)

	// Test credentials
	_, errTest := valveapi.GetNextMatch(s.configurationService.GetConfig().Steam.SteamAPIKey, user.Steam.ID, authCode, sc)
	if errTest != nil {
		return errors.New("invalid authentication code or last share code")
	}

	user.AddSteamMatchHistoryAuthenticationCode(authCode, shareCode)
	errAdd := s.repo.UpdateMatchAuthCode(user)
	if errAdd != nil {
		return errAdd
	}

	errSc := s.repo.UpdateLatestShareCode(user)
	if errSc != nil {
		return errSc
	}

	errApi := s.UpdateSteamAPIUsage(user, true)
	if errApi != nil {
		return errApi
	}

	return nil
}

func (s *Service) UpdateSteamAPIUsage(u *User, active bool) error {
	if active && u.Steam.AuthCode == "" {
		return errors.New("missing steam api auth code")
	}

	u.Steam.APIEnabled = active
	return s.repo.UpdateSteamAPIUsage(u)
}
func (s *Service) UpdateLatestShareCode(u *User, sc *share_code.ShareCodeData) error {
	u.Steam.LastShareCode = sc.Encoded
	return s.repo.UpdateLatestShareCode(u)
}

func (s *Service) SigninUsingSteam(id uint64, nickname string) (*User, error) {
	u, err := s.repo.FindBySteamId(id)
	if u == nil {
		u, e := NewUserUsingSteam(id, nickname)
		if e != nil {
			return nil, e
		}
		return u, s.CreateUser(u)
	}

	return u, err
}

func (s *Service) SigninUsingFaceit(id entity.ID, nickname string) (*User, error) {
	u, err := s.repo.FindByFaceitId(id)
	if u == nil {
		u, e := NewUserUsingFaceit(id, nickname)
		if e != nil {
			return nil, e
		}
		return u, s.CreateUser(u)
	}

	return u, err
}

func (s *Service) QueryLatestShareCode(u *User) (*share_code.ShareCodeData, error) {
	if !u.Steam.APIEnabled {
		return nil, errors.New("user: api usage is disabled")
	}

	steamID := u.Steam.ID
	shareCode, err := valveapi.GetNextMatch(s.configurationService.GetConfig().Steam.SteamAPIKey, steamID, u.Steam.AuthCode, u.Steam.LastShareCode)

	// Disable user on error
	if err != nil {
		if os.IsTimeout(err) {
			return nil, errors.New("user: lost connection while querying the steam api for the latest sharecode")
		}
		s.UpdateSteamAPIUsage(u, false)
		log.Warnf("disabled csgo user %d due to an error in fetching the share code", steamID)
		return nil, err
	}

	// No new match.
	if shareCode == "" {
		log.Debugf("no new match found for %d", steamID)
		return nil, nil
	}

	log.Infof("found match share code %v for %d", shareCode, u.Steam.ID)

	sc, err := share_code.Decode(shareCode)
	if err != nil {
		return nil, fmt.Errorf("invalid share code %s", sc.Encoded)
	}

	return sc, nil
}
