package faceitapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// MatchResponse contains information about the latest match.
type PlayerMatchHistoryResponse struct {
	Result *[]PlayerMatchHistoryEntry `json:"items"`
}

type PlayerMatchHistoryEntry struct {
	MatchId string `json:"match_id"`
}

// InvalidFaceitApiCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type InvalidFaceitApiCredentials struct{}

func (e *InvalidFaceitApiCredentials) Error() string {
	return "Invalid faceit api credentials"
}

// InvalidFaceitApiCredentials is used to notify when the supplied credentials are not valid / cannot be used with the api.
type FaceitApiConnectionIssues struct{}

func (e *FaceitApiConnectionIssues) Error() string {
	return "Too many requests or unavailable api"
}

// GetPlayerMatchHistory returns the match history for a given player.
func GetPlayerMatchHistory(faceitAPIKey string, playerId uuid.UUID) (*PlayerMatchHistoryResponse, error) {
	u, err := url.Parse(fmt.Sprintf("https://open.faceit.com/data/v4/players/%s/history", playerId))
	if err != nil {
		return nil, errors.New("faceitapi: unable to parse url")
	}

	// Build query
	q := u.Query()
	q.Set("game", "csgo")
	q.Set("offset", "0")
	q.Set("limit", "20")
	// Query the last 24 hours only
	q.Set("from", fmt.Sprint(time.Now().AddDate(0, -1, 0).Unix()))
	u.RawQuery = q.Encode()

	matchResponse := &PlayerMatchHistoryResponse{}

	// Request player match history.
	client := &http.Client{}
	request, _ := http.NewRequest("GET", u.String(), nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", faceitAPIKey))
	r, rErr := client.Do(request)
	if rErr != nil {
		return nil, rErr
	}

	if r.StatusCode == http.StatusTooManyRequests || r.StatusCode == http.StatusServiceUnavailable {
		r.Body.Close()
		return nil, &FaceitApiConnectionIssues{}
	}

	if r.StatusCode == http.StatusUnauthorized || r.StatusCode == http.StatusForbidden {
		r.Body.Close()
		return nil, &InvalidFaceitApiCredentials{}
	}

	if err = json.NewDecoder(r.Body).Decode(matchResponse); err != nil {
		r.Body.Close()
		return nil, err
	}

	defer r.Body.Close()

	return matchResponse, nil
}
