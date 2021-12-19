package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Cludch/csgo-tools/internal/config"
	"github.com/Cludch/csgo-tools/internal/domain/user"
	"github.com/golang-jwt/jwt"
	"github.com/markbates/goth"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	configService config.UseCase
	userService   user.UseCase
}

// Claims JWT claims
type Claims struct {
	SteamId uint64 `json:"steamId"`
	jwt.StandardClaims
}

func NewService(c config.UseCase, u user.UseCase) *Service {
	return &Service{
		configService: c,
		userService:   u,
	}
}

func (s *Service) HandleAuth(gothUser goth.User) (string, error) {
	userId := gothUser.UserID
	provider := gothUser.Provider

	var dbUser *user.User
	var err error

	if provider == "steam" {
		steamId, _ := strconv.ParseUint(userId, 10, 64)
		dbUser, err = s.userService.SigninUsingSteam(steamId, gothUser.Name)
	} else {
		return "", fmt.Errorf("unknown authentication provider: %s", provider)
	}

	if err != nil {
		return "", nil
	}

	return s.generateToken(dbUser, provider), nil
}

func (s *Service) ValidateToken(encodedToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.configService.GetConfig().Auth.Secret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s *Service) generateToken(dbUser *user.User, provider string) string {
	claims := &Claims{
		dbUser.Steam.ID,
		jwt.StandardClaims{
			Id:        dbUser.ID.String(),
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    provider,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.configService.GetConfig().Auth.Secret))
	if err != nil {
		log.Error(err)
	}
	return t
}
