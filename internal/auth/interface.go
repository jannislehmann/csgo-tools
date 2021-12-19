package auth

import (
	"github.com/markbates/goth"
)

type UseCase interface {
	HandleAuth(user goth.User) (string, error)
	ValidateToken(encodedToken string) (*Claims, error)
}
