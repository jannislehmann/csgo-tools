package player

type PlayerList struct {
	Players []*PlayerListEntry `json:"players"`
}

type PlayerListEntry struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Games int    `json:"games"`
}

// PlayerGameStats describes average stats across all matches.
type PlayerGameStats struct {
	Games                      int    `json:"games"`
	SteamID                    uint64 `json:"id"`
	Name                       string `json:"name"`
	KillsPerRound              int    `json:"killsPerRound"`
	EntryKillsPerRound         int    `json:"entryKillsPerRound"`
	OpeningDuelAttempsPerRound int    `json:"openingDuelAttemptsPerRound"`
	HeadshotsPerRound          int    `json:"headshotsPerRound"`
	AssistsPerRound            int    `json:"assistsPerRound"`
	DeathsPerRound             int    `json:"deathsPerRound"`
	MVPsPerRound               int    `json:"mvpsPerRound"`
	Won1v3                     int    `json:"won1v3"`
	Won1v4                     int    `json:"won1v4"`
	Won1v5                     int    `json:"won1v5"`
	RoundsWith3K               int    `json:"roundsWith3k"`
	RoundsWith4K               int    `json:"roundsWith4k"`
	RoundsWith5K               int    `json:"roundsWith5k"`
}
