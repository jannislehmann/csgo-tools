package entity

import (
	"strconv"
	"time"

	"github.com/Cludch/csgo-tools/internal/config"
	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

// CSGOUser holds information about a csgo user whose match history should be watched
type CSGOUser struct {
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	DeletedAt                      gorm.DeletedAt `gorm:"index"`
	SteamID                        uint64         `gorm:"primaryKey;autoIncrement:false"`
	MatchHistoryAuthenticationCode string
	ShareCode                      ShareCode
	ShareCodeID                    string
	Disabled                       bool
}

// UpdateLatestShareCode sets the latest fetched share code
func (c *CSGOUser) UpdateLatestShareCode(shareCode *ShareCode) {
	db.Model(&c).Update("ShareCode", shareCode)
}

// AddConfigUsers takes csgo users from the config file and turns them into database entities
func AddConfigUsers(users []*config.CSGOConfig) {
	for _, user := range users {
		steamID, err := strconv.ParseUint(user.SteamID, 10, 64)
		if err != nil {
			panic(err)
		}
		log.Debugf("add %d from config file to csgo_user database\n", steamID)

		// Create share code
		shareCode := user.KnownMatchCode
		sc := CreateShareCodeFromEncoded(shareCode)

		// Create user
		csgoUser := &CSGOUser{SteamID: steamID, MatchHistoryAuthenticationCode: user.HistoryAPIKey, ShareCode: *sc}
		db.FirstOrCreate(csgoUser)

		// Create match
		CreateMatch(sc)
	}
}
