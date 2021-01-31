package database

import (
	"time"

	"gorm.io/gorm"
)

// Match holds the central information about a csgo match
type Match struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	MatchID     uint64         `gorm:"primaryKey"`
	ShareCode   ShareCode
	MatchTime   time.Time
	DownloadURL string
	Downloaded  bool
}

// CreateMatch creates a match in the database. The match will be downloaded later.
func CreateMatch(shareCode *ShareCode) *Match {
	match := &Match{MatchID: shareCode.MatchID, ShareCode: *shareCode}
	db.FirstOrCreate(match)
	return match
}

// CreateDownloadedMatchFromMatchID creates a match in the database. The match will be marked as downloaded.
func CreateDownloadedMatchFromMatchID(matchID uint64) *Match {
	match := &Match{MatchID: matchID, Downloaded: true}
	db.FirstOrCreate(match)
	return match
}
