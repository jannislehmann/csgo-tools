package user_test

import (
	"testing"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/internal/domain/user"
	share_code_data "github.com/Cludch/csgo-tools/pkg/share_code"
	"github.com/stretchr/testify/assert"
)

const validShareCode = "CSGO-12345-12345-12345"

func TestNewUserUsingFaceit(t *testing.T) {
	id := uint64(1)
	name := "faceit"
	u, err := user.NewUserUsingSteam(id, name)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.Nil(t, u.Steam)
	assert.NotNil(t, u.Faceit)
	assert.Equal(t, u.Faceit.ID, id)
	assert.Equal(t, u.Faceit.Nickname, name)
}

func TestNewUserUsingSteam(t *testing.T) {
	id := entity.NewID()
	name := "faceit"
	u, err := user.NewUserUsingFaceit(id, name)
	assert.Nil(t, err)
	assert.NotNil(t, u.ID)
	assert.Nil(t, u.Steam)
	assert.NotNil(t, u.Faceit)
	assert.Equal(t, u.Steam.ID, id)
	assert.Equal(t, u.Faceit.Nickname, name)
}

func TestUpdateLastShareCode(t *testing.T) {
	u, _ := user.NewUserUsingSteam(1, "steam")
	shareCode, _ := share_code_data.Decode(validShareCode)
	err := u.UpdateLastShareCode(shareCode)
	assert.Nil(t, err)
	assert.NotNil(t, u.Steam.LastShareCode)
}

func TestAddSteamMatchHistoryAuthenticationCode(t *testing.T) {
	u, _ := user.NewUserUsingSteam(1, "steam")
	shareCode, _ := share_code_data.Decode(validShareCode)
	authenticationCode := "test"
	err := u.AddSteamMatchHistoryAuthenticationCode(authenticationCode, shareCode)
	assert.Nil(t, err)
	assert.NotNil(t, u.Steam.AuthCode)
}
