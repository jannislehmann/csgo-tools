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

func (s *Service) CreateUserUsingSteam(id uint64, nickname string) (*User, error) {
	u, err := NewUserUsingSteam(id, nickname)
	if err != nil {
		return nil, err
	}

	err = s.createUser(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) CreateUserUsingFaceit(id entity.ID, nickname string) (*User, error) {
	u, err := NewUserUsingFaceit(id, nickname)
	if err != nil {
		return nil, err
	}

	err = s.createUser(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) createUser(u *User) error {
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

	err := user.AddSteamMatchHistoryAuthenticationCode(authCode, shareCode)
	if err != nil {
		return err
	}

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
	log.Debugf("Attempting sign in using steam for user %d (%s)", id, nickname)
	user, err := s.repo.FindBySteamId(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			log.Debugf("no user with id %d found. creating a new one..", id)
			return s.CreateUserUsingSteam(id, nickname)
		}

		return nil, err
	}

	return user, nil
}

func (s *Service) SigninUsingFaceit(id entity.ID, nickname string) (*User, error) {
	log.Debugf("Attempting sign in using faceit for user %v (%s)", id, nickname)
	user, err := s.repo.FindByFaceitId(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			log.Debugf("no user with id %v found. creating a new one..", id)
			return s.CreateUserUsingFaceit(id, nickname)
		}

		return nil, err
	}

	return user, nil
}

func (s *Service) QueryLatestShareCode(u *User) (*share_code.ShareCodeData, error) {
	if !u.Steam.APIEnabled {
		return nil, errors.New("user: api usage is disabled")
	}

	steamID := u.Steam.ID
	shareCode, err := valveapi.GetNextMatch(s.configurationService.GetConfig().Steam.SteamAPIKey, steamID, u.Steam.AuthCode, u.Steam.LastShareCode)

	// Disable user on error.
	if err != nil {
		if os.IsTimeout(err) {
			return nil, errors.New("user: lost connection while querying the steam api for the latest sharecode")
		} else if s.configurationService.IsDebug() {
			const msg = "user.service: unable to query next valve match: %s"
			log.Errorf(msg, err)
		}

		/*
			// TODO Issue #64
			updateErr := s.UpdateSteamAPIUsage(u, false)
			if updateErr != nil {
				const msg = "disabled csgo user %d due to an error (%t) in fetching the share code"
				log.Warnf(msg, steamID, err)
			}
		*/
		return nil, err
	}

	// No new match.
	if shareCode == "" {
		const msg = "no new match found for %d"
		log.Debugf(msg, steamID)
		return nil, nil
	}

	const msg = "found match share code %v for %d"
	log.Infof(msg, shareCode, u.Steam.ID)

	sc, err := share_code.Decode(shareCode)
	if err != nil {
		const msg = "invalid share code %s"
		return nil, fmt.Errorf(msg, sc.Encoded)
	}

	return sc, nil
}
