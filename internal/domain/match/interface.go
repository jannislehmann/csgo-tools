package match

import (
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	"github.com/Cludch/csgo-tools/pkg/share_code"
)

type Repository interface {
	Create(*Match) error

	Find(entity.ID) (*Match, error)
	FindByFilename(filename string) (*Match, error)
	FindByFaceitId(entity.ID) (*Match, error)
	FindByValveId(uint64) (*Match, error)
	FindByValveOutcomeId(uint64) (*Match, error)

	List() ([]*Match, error)
	ListDownloadedMatches() ([]*Match, error)
	ListDownloadableMatches() ([]*Match, error)
	ListParsedMatches() ([]*Match, error)
	ListValveMatchesMissingDownloadUrl() ([]*Match, error)

	UpdateResult(*Match) error
	UpdateDownloadInformation(*Match) error
	UpdateStatus(*Match) error
	UpdateDownloaded(*Match) error

	Delete(entity.ID) error
}

type UseCase interface {
	CreateMatchFromManualUpload(filename string, matchTime time.Time) (*Match, error)
	CreateMatchFromSharecode(*share_code.ShareCodeData) (*Match, error)

	GetAll() ([]*Match, error)
	GetAllParsed() ([]*Match, error)
	GetMatch(entity.ID) (*Match, error)
	GetMatchByFilename(filename string) (*Match, error)
	GetMatchByValveId(uint64) (*Match, error)
	GetMatchByValveOutcomeId(uint64) (*Match, error)
	GetMatchByFaceitId(entity.ID) (*Match, error)
	GetDownloadableMatches() ([]*Match, error)
	GetValveMatchesMissingDownloadUrl() ([]*Match, error)
	GetParseableMatches(parserVersion byte) ([]*Match, error)

	UpdateStatus(*Match, Status) error
	UpdateResult(m *Match, r *MatchResult, parserVersion byte) error
	UpdateDownloadInformationForOutcomeId(matchId uint64, matchTime time.Time, url string) error
	SetDownloaded(m *Match, status Status, filename string) error
}
