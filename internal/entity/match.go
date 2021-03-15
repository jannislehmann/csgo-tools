package entity

import (
	"time"

	"gorm.io/gorm"
)

// Match holds the central information about a csgo match
type Match struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	ID            uint64         `gorm:"primaryKey"`
	ShareCode     ShareCode
	Filename      string
	MatchTime     time.Time
	DownloadURL   string
	Downloaded    bool
	ParserVersion byte
}

// CreateMatch creates a match in the database. The match will be downloaded later.
func CreateMatch(shareCode *ShareCode) *Match {
	match := &Match{ID: shareCode.OutcomeID, ShareCode: *shareCode}
	db.FirstOrCreate(match)
	return match
}

// CreateDownloadedMatchFromMatchID creates a match in the database. The match will be marked as downloaded.
func CreateDownloadedMatchFromMatchID(matchID uint64, fileName string, matchDate time.Time) *Match {
	match := &Match{ID: matchID, Filename: fileName, Downloaded: true, MatchTime: matchDate}
	db.FirstOrCreate(match)
	return match
}
