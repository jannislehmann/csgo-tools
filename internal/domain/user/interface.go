package user

import (
	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/share_code"
)

// Repository defines repository functions for user entities.
type Repository interface {
	Create(*User) error

	Find(entity.ID) (*User, error)
	FindUsersContainingAuthenticationCode() ([]*User, error)
	FindBySteamId(uint64) (*User, error)
	FindByFaceitId(entity.ID) (*User, error)

	List() ([]*User, error)

	UpdateLatestShareCode(*User) error
	UpdateMatchAuthCode(u *User) error
	UpdateSteamAPIUsage(*User) error

	Delete(entity.ID) error
}

// UseCase defines the user service functions.
type UseCase interface {
	CreateUser(*User) error

	GetUser(entity.ID) (*User, error)
	GetUserBySteamId(uint64) (*User, error)
	GetUserByFaceitId(entity.ID) (*User, error)
	GetUsersWithAuthenticationCode() ([]*User, error)

	AddSteamMatchHistoryAuthenticationCode(user *User, authCode string, sc string) error
	UpdateSteamAPIUsage(*User, bool) error
	UpdateLatestShareCode(*User, *share_code.ShareCodeData) error

	SigninUsingSteam(uint64, string) (*User, error)
	SigninUsingFaceit(entity.ID, string) (*User, error)

	QueryLatestShareCode(*User) (*share_code.ShareCodeData, error)
}
