package user

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/share_code"
)

type User struct {
	ID        entity.ID   `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time   `json:"-" bson:"createdAt"`
	Steam     *SteamUser  `json:"steam" bson:"steam,omitempty"`
	Faceit    *FaceitUser `json:"faceit" bson:"faceit,omitempty"`
}

// SteamUser contains the information about the linked steam account api authentication code to fetch the latest match.
type SteamUser struct {
	ID            uint64 `json:"id" bson:"id"`
	Nickname      string `json:"nickname" bson:"nickname"`
	AuthCode      string `json:"authCode" bson:"authCode"`
	LastShareCode string `json:"lastShareCode" bson:"lastShareCode"`
	APIEnabled    bool   `json:"apiEnabled" bson:"apiEnabled"`
}

// FaceitUser contains the information about the linked faceit account.
type FaceitUser struct {
	ID       entity.ID `json:"id" bson:"id"`
	Nickname string    `json:"nickname" bson:"nickname"`
}

func newUser() (*User, error) {
	u := &User{
		ID:        entity.NewID(),
		CreatedAt: time.Now(),
	}

	return u, nil
}

func NewUserUsingFaceit(id entity.ID, nickname string) (*User, error) {
	user, err := newUser()
	user.Faceit = &FaceitUser{ID: id, Nickname: nickname}
	return user, err
}

func NewUserUsingSteam(id uint64, nickname string) (*User, error) {
	user, err := newUser()
	user.Steam = &SteamUser{ID: id, Nickname: nickname, APIEnabled: false}
	return user, err
}

func (u *User) UpdateLastShareCode(sc *share_code.ShareCodeData) error {
	u.Steam.LastShareCode = sc.Encoded
	return nil
}

func (u *User) AddSteamMatchHistoryAuthenticationCode(authenticationCode string, sc *share_code.ShareCodeData) error {
	u.Steam.AuthCode = authenticationCode
	return u.UpdateLastShareCode(sc)
}
