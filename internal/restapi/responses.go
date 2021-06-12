package restapi

// PlayerGameStats describes average stats across all matches.
type PlayerGameStats struct {
	Games             int    `json:"games"`
	SteamID           uint64 `json:"steamId"`
	KillsPerGame      int    `json:"killsPerGame"`
	EntryKillsPerGame int    `json:"entryKillsPerGame"`
	HeadshotsPerGame  int    `json:"headshotsPerGame"`
	AssistsPerGame    int    `json:"assistsPerGame"`
	DeathsPerGame     int    `json:"deathsPerGame"`
	MVPsPerGame       int    `json:"mvpsPerGame"`
	Won1v3            int    `json:"won1v3"`
	Won1v4            int    `json:"won1v4"`
	Won1v5            int    `json:"won1v5"`
	RoundsWith3K      int    `json:"3k"`
	RoundsWith4K      int    `json:"4k"`
	RoundsWith5K      int    `json:"5k"`
}
