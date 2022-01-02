package player

type PlayerList struct {
	Players []*PlayerListEntry `json:"players"`
}

type PlayerListEntry struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Games int    `json:"games"`
	Wins  int    `json:"wins"`
	Rank  int    `json:"rank"`
}

// PlayerGameStats describes average stats across all matches.
type PlayerGameStats struct {
	Games                      int     `json:"games"`
	SteamID                    uint64  `json:"id"`
	Name                       string  `json:"name"`
	KillsPerRound              float32 `json:"killsPerRound"`
	EntryKillsPerRound         float32 `json:"entryKillsPerRound"`
	OpeningDuelAttempsPerRound float32 `json:"openingDuelAttemptsPerRound"`
	HeadshotsPerRound          float32 `json:"headshotsPerRound"`
	AssistsPerRound            float32 `json:"assistsPerRound"`
	DeathsPerRound             float32 `json:"deathsPerRound"`
	MVPsPerRound               float32 `json:"mvpsPerRound"`
	DamagePerRound             float32 `json:"damagePerRound"`
	Won1v3                     int     `json:"won1v3"`
	Won1v4                     int     `json:"won1v4"`
	Won1v5                     int     `json:"won1v5"`
	RoundsWith3K               int     `json:"roundsWith3k"`
	RoundsWith4K               int     `json:"roundsWith4k"`
	RoundsWith5K               int     `json:"roundsWith5k"`
}
