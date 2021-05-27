package valveapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// MatchResponse contains information about the latest match.
type MatchResponse struct {
	Result struct {
		Nextcode string `json:"nextcode"`
	} `json:"result"`
}

// InvalidMatchHistoryCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type InvalidMatchHistoryCredentials struct {
	SteamID string
}

func (e *InvalidMatchHistoryCredentials) Error() string {
	return fmt.Sprintf("Invalid match history credentials for steam id %v.", e.SteamID)
}

// GetNextMatch returns the next match's share code.
// It uses the saved share codes as the current one.
func GetNextMatch(steamAPIKey string, steamID uint64, historyAuthenticationCode string, lastShareCode string) (string, error) {
	// Get latest match
	u, err := url.Parse("https://api.steampowered.com/ICSGOPlayers_730/GetNextMatchSharingCode/v1")
	if err != nil {
		log.Error(err)
	}

	steamIDString := strconv.FormatUint(steamID, 10)

	// Build query
	q := u.Query()
	q.Set("key", steamAPIKey)
	q.Set("steamid", steamIDString)
	q.Set("steamidkey", historyAuthenticationCode)
	q.Set("knowncode", lastShareCode)
	u.RawQuery = q.Encode()

	matchResponse := &MatchResponse{}

	// Request match code.
	r, err := http.Get(u.String())
	if err != nil {
		log.Error(err)
		return "", err
	}

	// Forbidden = wrong api keys.
	// Precondition Failed = Know match code or steam id wrong.
	if r.StatusCode == http.StatusForbidden || r.StatusCode == http.StatusPreconditionFailed {
		r.Body.Close()
		return "", &InvalidMatchHistoryCredentials{SteamID: steamIDString}
	}

	// Accepted means that there is no recent match code available.
	if r.StatusCode == http.StatusAccepted {
		r.Body.Close()
		return "", nil
	}

	err = json.NewDecoder(r.Body).Decode(matchResponse)

	if err != nil {
		r.Body.Close()
		log.Error(err)
		return "", err
	}

	defer r.Body.Close()

	return matchResponse.Result.Nextcode, nil
}
