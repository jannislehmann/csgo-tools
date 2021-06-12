package user

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/go-playground/validator"
)

var validate = validator.New()

type User struct {
	ID        entity.ID   `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time   `json:"-" bson:"createdAt"`
	Steam     *SteamUser  `json:"steam" bson:"steam,omitempty"`
	Faceit    *FaceitUser `json:"faceit" bson:"faceit,omitempty"`
}

// SteamUser contains the information about the linked steam account api authentication code to fetch the latest match.
type SteamUser struct {
	ID            uint64 `json:"id" bson:"id" validate:"required"`
	Nickname      string `json:"nickname" bson:"nickname" validate:"required"`
	AuthCode      string `json:"authCode" bson:"authCode"`
	LastShareCode string `json:"lastShareCode" bson:"lastShareCode"`
	APIEnabled    bool   `json:"apiEnabled" bson:"apiEnabled" validate:"required"`
}

// FaceitUser contains the information about the linked faceit account.
type FaceitUser struct {
	ID       entity.ID `json:"id" bson:"id" validate:"required"`
	Nickname string    `json:"nickname" bson:"nickname" validate:"required"`
}

func newUser() (*User, error) {
	u := &User{
		ID:        entity.NewID(),
		CreatedAt: time.Now(),
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}
	return u, nil
}

func NewUserUsingFaceit(id entity.ID, nickname string) (*User, error) {
	u, err := newUser()
	if err != nil {
		return nil, err
	}

	u.Faceit = &FaceitUser{ID: id, Nickname: nickname}
	return u, u.Validate()
}

func NewUserUsingSteam(id uint64, nickname string) (*User, error) {
	u, err := newUser()
	if err != nil {
		return nil, err
	}

	u.Steam = &SteamUser{ID: id, Nickname: nickname, APIEnabled: false}
	return u, u.Validate()
}

func (u *User) UpdateLastShareCode(sc *share_code.ShareCodeData) error {
	u.Steam.LastShareCode = sc.Encoded
	return u.Validate()
}

func (u *User) AddSteamMatchHistoryAuthenticationCode(authenticationCode string, sc *share_code.ShareCodeData) error {
	u.Steam.AuthCode = authenticationCode
	if err := u.Validate(); err != nil {
		return err
	}

	return u.UpdateLastShareCode(sc)
}

func (u *User) Validate() error {
	err := validate.Struct(u)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	return nil
}
